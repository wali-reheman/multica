package handler

<<<<<<< HEAD
import (
=======
// MULTICA-LOCAL: Local project management with git version history.

import (
	"database/sql"
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
<<<<<<< HEAD
	"github.com/jackc/pgx/v5/pgtype"
=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
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
<<<<<<< HEAD
	FileCount     int32   `json:"file_count"`
=======
	FileCount     int64   `json:"file_count"`
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
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
<<<<<<< HEAD
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
=======
		ID:            p.ID,
		WorkspaceID:   p.WorkspaceID,
		Name:          p.Name,
		LocalPath:     p.LocalPath,
		DefaultBranch: p.DefaultBranch,
		Language:      nullStringToPtr(p.Language),
		FileCount:     p.FileCount,
		SizeBytes:     p.SizeBytes,
		LastOpenedAt:  nullStringToPtr(p.LastOpenedAt),
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
		IsGitRepo:     isGitRepo,
	}
}

// ListProjects returns all projects in the workspace.
func (h *Handler) ListProjects(w http.ResponseWriter, r *http.Request) {
	workspaceID := resolveWorkspaceID(r)
<<<<<<< HEAD

	limit := 50
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
=======
	limit := int64(50)
	offset := int64(0)
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.ParseInt(l, 10, 64); err == nil {
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
			limit = v
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
<<<<<<< HEAD
		if v, err := strconv.Atoi(o); err == nil {
=======
		if v, err := strconv.ParseInt(o, 10, 64); err == nil {
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
			offset = v
		}
	}

	projects, err := h.Queries.ListLocalProjects(r.Context(), db.ListLocalProjectsParams{
<<<<<<< HEAD
		WorkspaceID: parseUUID(workspaceID),
		Limit:       int32(limit),
		Offset:      int32(offset),
=======
		WorkspaceID: workspaceID,
		Limit:       limit,
		Offset:      offset,
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list projects")
		return
	}

	resp := make([]ProjectResponse, len(projects))
	for i, p := range projects {
		resp[i] = projectToResponse(p, h.GitService)
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
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
<<<<<<< HEAD

	// Update last opened timestamp
	h.Queries.UpdateLocalProjectLastOpened(r.Context(), project.ID)

	writeJSON(w, http.StatusOK, projectToResponse(project, h.GitService))
}

// CreateProjectRequest is the request body for creating a project.
=======
	h.Queries.UpdateLocalProjectLastOpened(r.Context(), project.ID)
	writeJSON(w, http.StatusOK, projectToResponse(project, h.GitService))
}

>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
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
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	workspaceID := resolveWorkspaceID(r)

	var req CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	if req.LocalPath == "" {
		writeError(w, http.StatusBadRequest, "local_path is required")
		return
	}

<<<<<<< HEAD
	// Validate path exists and is a directory
=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
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
<<<<<<< HEAD
		WorkspaceID: parseUUID(workspaceID),
=======
		WorkspaceID: workspaceID,
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
		LocalPath:   req.LocalPath,
	})
	if err == nil {
		writeError(w, http.StatusConflict, "project at this path is already tracked")
		return
	}

<<<<<<< HEAD
	// Auto-detect name if not provided
=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	name := req.Name
	if name == "" {
		name = info.Name()
	}

<<<<<<< HEAD
	// Init git if requested
=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	defaultBranch := "main"
	if req.InitGit && h.GitService != nil {
		h.GitService.Init(req.LocalPath)
	}
<<<<<<< HEAD

	// Detect default branch if it's already a repo
=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
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
<<<<<<< HEAD
		WorkspaceID:   parseUUID(workspaceID),
=======
		ID:            newUUID(),
		WorkspaceID:   workspaceID,
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
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

<<<<<<< HEAD
	// Start watching if watcher is available
	if h.WatcherService != nil {
		h.WatcherService.Watch(uuidToString(project.ID), workspaceID, req.LocalPath)
=======
	if h.WatcherService != nil {
		h.WatcherService.Watch(project.ID, workspaceID, req.LocalPath)
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	}

	writeJSON(w, http.StatusCreated, resp)
}

// UpdateProject updates a project's metadata.
func (h *Handler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	project, ok := h.loadProjectForUser(w, r)
	if !ok {
		return
	}
<<<<<<< HEAD
	workspaceID := uuidToString(project.WorkspaceID)
=======
	workspaceID := project.WorkspaceID
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609

	var raw map[string]json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	params := db.UpdateLocalProjectParams{ID: project.ID}
	if v, ok := raw["name"]; ok {
		var name string
		json.Unmarshal(v, &name)
<<<<<<< HEAD
		params.Name = pgtype.Text{String: name, Valid: true}
=======
		params.Name = sql.NullString{String: name, Valid: true}
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	}
	if v, ok := raw["default_branch"]; ok {
		var branch string
		json.Unmarshal(v, &branch)
<<<<<<< HEAD
		params.DefaultBranch = pgtype.Text{String: branch, Valid: true}
=======
		params.DefaultBranch = sql.NullString{String: branch, Valid: true}
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
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
<<<<<<< HEAD

	writeJSON(w, http.StatusOK, resp)
}

// DeleteProject removes a project from tracking (does not delete files).
=======
	writeJSON(w, http.StatusOK, resp)
}

// DeleteProject removes a project from tracking.
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	project, ok := h.loadProjectForUser(w, r)
	if !ok {
		return
	}
<<<<<<< HEAD
	workspaceID := uuidToString(project.WorkspaceID)
	projectID := uuidToString(project.ID)
=======
	workspaceID := project.WorkspaceID
	projectID := project.ID
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609

	if err := h.Queries.DeleteLocalProject(r.Context(), project.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete project")
		return
	}
<<<<<<< HEAD

	// Stop watching
	if h.WatcherService != nil {
		h.WatcherService.Unwatch(projectID)
	}

	actorType, actorID := h.resolveActor(r, userID, workspaceID)
	h.publish(protocol.EventProjectDeleted, workspaceID, actorType, actorID,
		map[string]any{"project_id": projectID})

=======
	if h.WatcherService != nil {
		h.WatcherService.Unwatch(projectID)
	}
	actorType, actorID := h.resolveActor(r, userID, workspaceID)
	h.publish(protocol.EventProjectDeleted, workspaceID, actorType, actorID,
		map[string]any{"project_id": projectID})
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	w.WriteHeader(http.StatusNoContent)
}

// --- Git endpoints ---

<<<<<<< HEAD
// GetProjectCommits returns paginated commit history.
=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
func (h *Handler) GetProjectCommits(w http.ResponseWriter, r *http.Request) {
	project, ok := h.loadProjectForUser(w, r)
	if !ok {
		return
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	if h.GitService == nil {
		writeError(w, http.StatusServiceUnavailable, "git service not available")
		return
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
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
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	commits, err := h.GitService.Log(project.LocalPath, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get commit history")
		return
	}
<<<<<<< HEAD

	writeJSON(w, http.StatusOK, map[string]any{
		"commits": commits,
		"total":   len(commits),
	})
}

// GetProjectCommitDetail returns full commit details with diffs.
=======
	writeJSON(w, http.StatusOK, map[string]any{"commits": commits, "total": len(commits)})
}

>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
func (h *Handler) GetProjectCommitDetail(w http.ResponseWriter, r *http.Request) {
	project, ok := h.loadProjectForUser(w, r)
	if !ok {
		return
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	if h.GitService == nil {
		writeError(w, http.StatusServiceUnavailable, "git service not available")
		return
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	sha := chi.URLParam(r, "sha")
	if sha == "" {
		writeError(w, http.StatusBadRequest, "commit sha is required")
		return
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	detail, err := h.GitService.Show(project.LocalPath, sha)
	if err != nil {
		writeError(w, http.StatusNotFound, "commit not found")
		return
	}
<<<<<<< HEAD

	writeJSON(w, http.StatusOK, detail)
}

// GetProjectStatus returns the working tree status.
=======
	writeJSON(w, http.StatusOK, detail)
}

>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
func (h *Handler) GetProjectStatus(w http.ResponseWriter, r *http.Request) {
	project, ok := h.loadProjectForUser(w, r)
	if !ok {
		return
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	if h.GitService == nil {
		writeError(w, http.StatusServiceUnavailable, "git service not available")
		return
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	status, err := h.GitService.Status(project.LocalPath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get status")
		return
	}
<<<<<<< HEAD

	writeJSON(w, http.StatusOK, status)
}

// CreateProjectCommitRequest is the request for creating a commit.
=======
	writeJSON(w, http.StatusOK, status)
}

>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
type CreateProjectCommitRequest struct {
	Message string   `json:"message"`
	Files   []string `json:"files"`
}

<<<<<<< HEAD
// CreateProjectCommit stages files and creates a commit.
=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
func (h *Handler) CreateProjectCommit(w http.ResponseWriter, r *http.Request) {
	project, ok := h.loadProjectForUser(w, r)
	if !ok {
		return
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	if h.GitService == nil {
		writeError(w, http.StatusServiceUnavailable, "git service not available")
		return
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	var req CreateProjectCommitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	if req.Message == "" {
		writeError(w, http.StatusBadRequest, "message is required")
		return
	}
<<<<<<< HEAD

	// Stage files
=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	if len(req.Files) > 0 {
		if err := h.GitService.Add(project.LocalPath, req.Files); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to stage files")
			return
		}
	}

<<<<<<< HEAD
	// Get user info for author
=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	userID := requestUserID(r)
	authorName := "Multica User"
	authorEmail := "user@multica.ai"
	if userID != "" {
<<<<<<< HEAD
		user, err := h.Queries.GetUser(r.Context(), parseUUID(userID))
=======
		user, err := h.Queries.GetUser(r.Context(), userID)
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
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
<<<<<<< HEAD

	writeJSON(w, http.StatusCreated, map[string]string{
		"hash": hash,
	})
}

// GetProjectBranches returns all branches.
=======
	writeJSON(w, http.StatusCreated, map[string]string{"hash": hash})
}

>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
func (h *Handler) GetProjectBranches(w http.ResponseWriter, r *http.Request) {
	project, ok := h.loadProjectForUser(w, r)
	if !ok {
		return
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	if h.GitService == nil {
		writeError(w, http.StatusServiceUnavailable, "git service not available")
		return
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	branches, err := h.GitService.Branches(project.LocalPath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list branches")
		return
	}
<<<<<<< HEAD

	writeJSON(w, http.StatusOK, map[string]any{
		"branches": branches,
	})
}

// CreateProjectBranchRequest is the request for creating a branch.
=======
	writeJSON(w, http.StatusOK, map[string]any{"branches": branches})
}

>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
type CreateProjectBranchRequest struct {
	Name string `json:"name"`
}

<<<<<<< HEAD
// CreateProjectBranch creates a new branch.
=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
func (h *Handler) CreateProjectBranch(w http.ResponseWriter, r *http.Request) {
	project, ok := h.loadProjectForUser(w, r)
	if !ok {
		return
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	if h.GitService == nil {
		writeError(w, http.StatusServiceUnavailable, "git service not available")
		return
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	var req CreateProjectBranchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	if err := h.GitService.CreateBranch(project.LocalPath, req.Name); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create branch")
		return
	}
<<<<<<< HEAD

	writeJSON(w, http.StatusCreated, map[string]string{
		"name": req.Name,
	})
}

// CheckoutProjectBranchRequest is the request for switching branches.
=======
	writeJSON(w, http.StatusCreated, map[string]string{"name": req.Name})
}

>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
type CheckoutProjectBranchRequest struct {
	Branch string `json:"branch"`
}

<<<<<<< HEAD
// CheckoutProjectBranch switches to a branch.
=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
func (h *Handler) CheckoutProjectBranch(w http.ResponseWriter, r *http.Request) {
	project, ok := h.loadProjectForUser(w, r)
	if !ok {
		return
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	if h.GitService == nil {
		writeError(w, http.StatusServiceUnavailable, "git service not available")
		return
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	var req CheckoutProjectBranchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	if req.Branch == "" {
		writeError(w, http.StatusBadRequest, "branch is required")
		return
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	if err := h.GitService.Checkout(project.LocalPath, req.Branch); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to checkout branch")
		return
	}
<<<<<<< HEAD

	writeJSON(w, http.StatusOK, map[string]string{
		"branch": req.Branch,
	})
}

// GetProjectDiff returns the working tree diff.
=======
	writeJSON(w, http.StatusOK, map[string]string{"branch": req.Branch})
}

>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
func (h *Handler) GetProjectDiff(w http.ResponseWriter, r *http.Request) {
	project, ok := h.loadProjectForUser(w, r)
	if !ok {
		return
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	if h.GitService == nil {
		writeError(w, http.StatusServiceUnavailable, "git service not available")
		return
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	diffs, err := h.GitService.Diff(project.LocalPath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get diff")
		return
	}
<<<<<<< HEAD

	writeJSON(w, http.StatusOK, map[string]any{
		"diffs": diffs,
	})
}

// InitProjectGit initializes a git repo for the project.
=======
	writeJSON(w, http.StatusOK, map[string]any{"diffs": diffs})
}

>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
func (h *Handler) InitProjectGit(w http.ResponseWriter, r *http.Request) {
	project, ok := h.loadProjectForUser(w, r)
	if !ok {
		return
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	if h.GitService == nil {
		writeError(w, http.StatusServiceUnavailable, "git service not available")
		return
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	initialized, err := h.GitService.Init(project.LocalPath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to initialize git repo")
		return
	}
<<<<<<< HEAD

	writeJSON(w, http.StatusOK, map[string]any{
		"initialized": initialized,
	})
}

// loadProjectForUser loads a project and validates workspace membership.
=======
	writeJSON(w, http.StatusOK, map[string]any{"initialized": initialized})
}

>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
func (h *Handler) loadProjectForUser(w http.ResponseWriter, r *http.Request) (db.LocalProject, bool) {
	if _, ok := requireUserID(w, r); !ok {
		return db.LocalProject{}, false
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	workspaceID := resolveWorkspaceID(r)
	if workspaceID == "" {
		writeError(w, http.StatusBadRequest, "workspace_id is required")
		return db.LocalProject{}, false
	}
<<<<<<< HEAD

	projectID := chi.URLParam(r, "projectId")
	project, err := h.Queries.GetLocalProjectInWorkspace(r.Context(), db.GetLocalProjectInWorkspaceParams{
		ID:          parseUUID(projectID),
		WorkspaceID: parseUUID(workspaceID),
=======
	projectID := chi.URLParam(r, "projectId")
	project, err := h.Queries.GetLocalProjectInWorkspace(r.Context(), db.GetLocalProjectInWorkspaceParams{
		ID:          projectID,
		WorkspaceID: workspaceID,
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	})
	if err != nil {
		writeError(w, http.StatusNotFound, "project not found")
		return db.LocalProject{}, false
	}
	return project, true
}
