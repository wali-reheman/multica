package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/multica-ai/multica/server/pkg/db/generated"
	"github.com/multica-ai/multica/server/internal/service"
	"github.com/multica-ai/multica/server/pkg/protocol"
)

// ProjectResponse is the JSON response for a local project.
type ProjectResponse struct {
	ID            string  `json:"id"`
	WorkspaceID   string  `json:"workspace_id"`
	Name          string  `json:"name"`
	LocalPath     string  `json:"local_path"`
	DefaultBranch string  `json:"default_branch"`
	Language      *string `json:"language"`
	FileCount     int32   `json:"file_count"`
	SizeBytes     int64   `json:"size_bytes"`
	LastOpenedAt  *string `json:"last_opened_at"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	IsGitRepo     bool    `json:"is_git_repo"`
}

func projectToResponse(p db.LocalProject, gitSvc *service.GitService) ProjectResponse {
	isGitRepo := false
	if gitSvc != nil {
		isGitRepo = gitSvc.IsRepo(p.LocalPath)
	}
	return ProjectResponse{
		ID:            uuidToString(p.ID),
		WorkspaceID:   uuidToString(p.WorkspaceID),
		Name:          p.Name,
		LocalPath:     p.LocalPath,
		DefaultBranch: p.DefaultBranch,
		Language:      textToPtr(p.Language),
		FileCount:     p.FileCount,
		SizeBytes:     p.SizeBytes,
		LastOpenedAt:  timestampToPtr(p.LastOpenedAt),
		CreatedAt:     timestampToString(p.CreatedAt),
		UpdatedAt:     timestampToString(p.UpdatedAt),
		IsGitRepo:     isGitRepo,
	}
}

// ListProjects returns all projects in the workspace.
func (h *Handler) ListProjects(w http.ResponseWriter, r *http.Request) {
	workspaceID := resolveWorkspaceID(r)

	limit := 50
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			limit = v
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil {
			offset = v
		}
	}

	projects, err := h.Queries.ListLocalProjects(r.Context(), db.ListLocalProjectsParams{
		WorkspaceID: parseUUID(workspaceID),
		Limit:       int32(limit),
		Offset:      int32(offset),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list projects")
		return
	}

	resp := make([]ProjectResponse, len(projects))
	for i, p := range projects {
		resp[i] = projectToResponse(p, h.GitService)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"projects": resp,
		"total":    len(resp),
	})
}

// GetProject returns a single project.
func (h *Handler) GetProject(w http.ResponseWriter, r *http.Request) {
	project, ok := h.loadProjectForUser(w, r)
	if !ok {
		return
	}

	// Update last opened timestamp
	h.Queries.UpdateLocalProjectLastOpened(r.Context(), project.ID)

	writeJSON(w, http.StatusOK, projectToResponse(project, h.GitService))
}

// CreateProjectRequest is the request body for creating a project.
type CreateProjectRequest struct {
	Name      string `json:"name"`
	LocalPath string `json:"local_path"`
	InitGit   bool   `json:"init_git"`
}

// CreateProject adds a new local project.
func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	workspaceID := resolveWorkspaceID(r)

	var req CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.LocalPath == "" {
		writeError(w, http.StatusBadRequest, "local_path is required")
		return
	}

	// Validate path exists and is a directory
	info, err := os.Stat(req.LocalPath)
	if err != nil {
		writeError(w, http.StatusBadRequest, "path does not exist or is not accessible")
		return
	}
	if !info.IsDir() {
		writeError(w, http.StatusBadRequest, "path is not a directory")
		return
	}

	// Check not already tracked
	_, err = h.Queries.GetLocalProjectByPath(r.Context(), db.GetLocalProjectByPathParams{
		WorkspaceID: parseUUID(workspaceID),
		LocalPath:   req.LocalPath,
	})
	if err == nil {
		writeError(w, http.StatusConflict, "project at this path is already tracked")
		return
	}

	// Auto-detect name if not provided
	name := req.Name
	if name == "" {
		name = info.Name()
	}

	// Init git if requested
	defaultBranch := "main"
	if req.InitGit && h.GitService != nil {
		h.GitService.Init(req.LocalPath)
	}

	// Detect default branch if it's already a repo
	if h.GitService != nil && h.GitService.IsRepo(req.LocalPath) {
		branches, err := h.GitService.Branches(req.LocalPath)
		if err == nil {
			for _, b := range branches {
				if b.IsHead && !b.IsRemote {
					defaultBranch = b.Name
					break
				}
			}
		}
	}

	project, err := h.Queries.CreateLocalProject(r.Context(), db.CreateLocalProjectParams{
		WorkspaceID:   parseUUID(workspaceID),
		Name:          name,
		LocalPath:     req.LocalPath,
		DefaultBranch: defaultBranch,
	})
	if err != nil {
		if isUniqueViolation(err) {
			writeError(w, http.StatusConflict, "project at this path is already tracked")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create project")
		return
	}

	actorType, actorID := h.resolveActor(r, userID, workspaceID)
	resp := projectToResponse(project, h.GitService)
	h.publish(protocol.EventProjectCreated, workspaceID, actorType, actorID,
		map[string]any{"project": resp})

	// Start watching if watcher is available
	if h.WatcherService != nil {
		h.WatcherService.Watch(uuidToString(project.ID), workspaceID, req.LocalPath)
	}

	writeJSON(w, http.StatusCreated, resp)
}

// UpdateProject updates a project's metadata.
func (h *Handler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	project, ok := h.loadProjectForUser(w, r)
	if !ok {
		return
	}
	workspaceID := uuidToString(project.WorkspaceID)

	var raw map[string]json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	params := db.UpdateLocalProjectParams{ID: project.ID}
	if v, ok := raw["name"]; ok {
		var name string
		json.Unmarshal(v, &name)
		params.Name = pgtype.Text{String: name, Valid: true}
	}
	if v, ok := raw["default_branch"]; ok {
		var branch string
		json.Unmarshal(v, &branch)
		params.DefaultBranch = pgtype.Text{String: branch, Valid: true}
	}

	updated, err := h.Queries.UpdateLocalProject(r.Context(), params)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update project")
		return
	}

	actorType, actorID := h.resolveActor(r, userID, workspaceID)
	resp := projectToResponse(updated, h.GitService)
	h.publish(protocol.EventProjectUpdated, workspaceID, actorType, actorID,
		map[string]any{"project": resp})

	writeJSON(w, http.StatusOK, resp)
}

// DeleteProject removes a project from tracking (does not delete files).
func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	project, ok := h.loadProjectForUser(w, r)
	if !ok {
		return
	}
	workspaceID := uuidToString(project.WorkspaceID)
	projectID := uuidToString(project.ID)

	if err := h.Queries.DeleteLocalProject(r.Context(), project.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete project")
		return
	}

	// Stop watching
	if h.WatcherService != nil {
		h.WatcherService.Unwatch(projectID)
	}

	actorType, actorID := h.resolveActor(r, userID, workspaceID)
	h.publish(protocol.EventProjectDeleted, workspaceID, actorType, actorID,
		map[string]any{"project_id": projectID})

	w.WriteHeader(http.StatusNoContent)
}

// --- Git endpoints ---

// GetProjectCommits returns paginated commit history.
func (h *Handler) GetProjectCommits(w http.ResponseWriter, r *http.Request) {
	project, ok := h.loadProjectForUser(w, r)
	if !ok {
		return
	}

	if h.GitService == nil {
		writeError(w, http.StatusServiceUnavailable, "git service not available")
		return
	}

	limit := 50
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			limit = v
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil {
			offset = v
		}
	}

	commits, err := h.GitService.Log(project.LocalPath, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get commit history")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"commits": commits,
		"total":   len(commits),
	})
}

// GetProjectCommitDetail returns full commit details with diffs.
func (h *Handler) GetProjectCommitDetail(w http.ResponseWriter, r *http.Request) {
	project, ok := h.loadProjectForUser(w, r)
	if !ok {
		return
	}

	if h.GitService == nil {
		writeError(w, http.StatusServiceUnavailable, "git service not available")
		return
	}

	sha := chi.URLParam(r, "sha")
	if sha == "" {
		writeError(w, http.StatusBadRequest, "commit sha is required")
		return
	}

	detail, err := h.GitService.Show(project.LocalPath, sha)
	if err != nil {
		writeError(w, http.StatusNotFound, "commit not found")
		return
	}

	writeJSON(w, http.StatusOK, detail)
}

// GetProjectStatus returns the working tree status.
func (h *Handler) GetProjectStatus(w http.ResponseWriter, r *http.Request) {
	project, ok := h.loadProjectForUser(w, r)
	if !ok {
		return
	}

	if h.GitService == nil {
		writeError(w, http.StatusServiceUnavailable, "git service not available")
		return
	}

	status, err := h.GitService.Status(project.LocalPath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get status")
		return
	}

	writeJSON(w, http.StatusOK, status)
}

// CreateProjectCommitRequest is the request for creating a commit.
type CreateProjectCommitRequest struct {
	Message string   `json:"message"`
	Files   []string `json:"files"`
}

// CreateProjectCommit stages files and creates a commit.
func (h *Handler) CreateProjectCommit(w http.ResponseWriter, r *http.Request) {
	project, ok := h.loadProjectForUser(w, r)
	if !ok {
		return
	}

	if h.GitService == nil {
		writeError(w, http.StatusServiceUnavailable, "git service not available")
		return
	}

	var req CreateProjectCommitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Message == "" {
		writeError(w, http.StatusBadRequest, "message is required")
		return
	}

	// Stage files
	if len(req.Files) > 0 {
		if err := h.GitService.Add(project.LocalPath, req.Files); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to stage files")
			return
		}
	}

	// Get user info for author
	userID := requestUserID(r)
	authorName := "Multica User"
	authorEmail := "user@multica.ai"
	if userID != "" {
		user, err := h.Queries.GetUser(r.Context(), parseUUID(userID))
		if err == nil {
			authorName = user.Name
			authorEmail = user.Email
		}
	}

	hash, err := h.GitService.Commit(project.LocalPath, req.Message, authorName, authorEmail)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create commit")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"hash": hash,
	})
}

// GetProjectBranches returns all branches.
func (h *Handler) GetProjectBranches(w http.ResponseWriter, r *http.Request) {
	project, ok := h.loadProjectForUser(w, r)
	if !ok {
		return
	}

	if h.GitService == nil {
		writeError(w, http.StatusServiceUnavailable, "git service not available")
		return
	}

	branches, err := h.GitService.Branches(project.LocalPath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list branches")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"branches": branches,
	})
}

// CreateProjectBranchRequest is the request for creating a branch.
type CreateProjectBranchRequest struct {
	Name string `json:"name"`
}

// CreateProjectBranch creates a new branch.
func (h *Handler) CreateProjectBranch(w http.ResponseWriter, r *http.Request) {
	project, ok := h.loadProjectForUser(w, r)
	if !ok {
		return
	}

	if h.GitService == nil {
		writeError(w, http.StatusServiceUnavailable, "git service not available")
		return
	}

	var req CreateProjectBranchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	if err := h.GitService.CreateBranch(project.LocalPath, req.Name); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create branch")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"name": req.Name,
	})
}

// CheckoutProjectBranchRequest is the request for switching branches.
type CheckoutProjectBranchRequest struct {
	Branch string `json:"branch"`
}

// CheckoutProjectBranch switches to a branch.
func (h *Handler) CheckoutProjectBranch(w http.ResponseWriter, r *http.Request) {
	project, ok := h.loadProjectForUser(w, r)
	if !ok {
		return
	}

	if h.GitService == nil {
		writeError(w, http.StatusServiceUnavailable, "git service not available")
		return
	}

	var req CheckoutProjectBranchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Branch == "" {
		writeError(w, http.StatusBadRequest, "branch is required")
		return
	}

	if err := h.GitService.Checkout(project.LocalPath, req.Branch); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to checkout branch")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"branch": req.Branch,
	})
}

// GetProjectDiff returns the working tree diff.
func (h *Handler) GetProjectDiff(w http.ResponseWriter, r *http.Request) {
	project, ok := h.loadProjectForUser(w, r)
	if !ok {
		return
	}

	if h.GitService == nil {
		writeError(w, http.StatusServiceUnavailable, "git service not available")
		return
	}

	diffs, err := h.GitService.Diff(project.LocalPath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get diff")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"diffs": diffs,
	})
}

// InitProjectGit initializes a git repo for the project.
func (h *Handler) InitProjectGit(w http.ResponseWriter, r *http.Request) {
	project, ok := h.loadProjectForUser(w, r)
	if !ok {
		return
	}

	if h.GitService == nil {
		writeError(w, http.StatusServiceUnavailable, "git service not available")
		return
	}

	initialized, err := h.GitService.Init(project.LocalPath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to initialize git repo")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"initialized": initialized,
	})
}

// loadProjectForUser loads a project and validates workspace membership.
func (h *Handler) loadProjectForUser(w http.ResponseWriter, r *http.Request) (db.LocalProject, bool) {
	if _, ok := requireUserID(w, r); !ok {
		return db.LocalProject{}, false
	}

	workspaceID := resolveWorkspaceID(r)
	if workspaceID == "" {
		writeError(w, http.StatusBadRequest, "workspace_id is required")
		return db.LocalProject{}, false
	}

	projectID := chi.URLParam(r, "projectId")
	project, err := h.Queries.GetLocalProjectInWorkspace(r.Context(), db.GetLocalProjectInWorkspaceParams{
		ID:          parseUUID(projectID),
		WorkspaceID: parseUUID(workspaceID),
	})
	if err != nil {
		writeError(w, http.StatusNotFound, "project not found")
		return db.LocalProject{}, false
	}
	return project, true
}
