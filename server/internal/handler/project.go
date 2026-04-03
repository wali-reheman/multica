package handler

// MULTICA-LOCAL: Local project management with git version history.

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
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
	FileCount     int64   `json:"file_count"`
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
		IsGitRepo:     isGitRepo,
	}
}

// ListProjects returns all projects in the workspace.
func (h *Handler) ListProjects(w http.ResponseWriter, r *http.Request) {
	workspaceID := resolveWorkspaceID(r)
	limit := int64(50)
	offset := int64(0)
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.ParseInt(l, 10, 64); err == nil {
			limit = v
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.ParseInt(o, 10, 64); err == nil {
			offset = v
		}
	}

	projects, err := h.Queries.ListLocalProjects(r.Context(), db.ListLocalProjectsParams{
		WorkspaceID: workspaceID,
		Limit:       limit,
		Offset:      offset,
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
	h.Queries.UpdateLocalProjectLastOpened(r.Context(), project.ID)
	writeJSON(w, http.StatusOK, projectToResponse(project, h.GitService))
}

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
		WorkspaceID: workspaceID,
		LocalPath:   req.LocalPath,
	})
	if err == nil {
		writeError(w, http.StatusConflict, "project at this path is already tracked")
		return
	}

	name := req.Name
	if name == "" {
		name = info.Name()
	}

	defaultBranch := "main"
	if req.InitGit && h.GitService != nil {
		h.GitService.Init(req.LocalPath)
	}
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
		ID:            newUUID(),
		WorkspaceID:   workspaceID,
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

	if h.WatcherService != nil {
		h.WatcherService.Watch(project.ID, workspaceID, req.LocalPath)
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
	workspaceID := project.WorkspaceID

	var raw map[string]json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	params := db.UpdateLocalProjectParams{ID: project.ID}
	if v, ok := raw["name"]; ok {
		var name string
		json.Unmarshal(v, &name)
		params.Name = sql.NullString{String: name, Valid: true}
	}
	if v, ok := raw["default_branch"]; ok {
		var branch string
		json.Unmarshal(v, &branch)
		params.DefaultBranch = sql.NullString{String: branch, Valid: true}
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

// DeleteProject removes a project from tracking.
func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}
	project, ok := h.loadProjectForUser(w, r)
	if !ok {
		return
	}
	workspaceID := project.WorkspaceID
	projectID := project.ID

	if err := h.Queries.DeleteLocalProject(r.Context(), project.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete project")
		return
	}
	if h.WatcherService != nil {
		h.WatcherService.Unwatch(projectID)
	}
	actorType, actorID := h.resolveActor(r, userID, workspaceID)
	h.publish(protocol.EventProjectDeleted, workspaceID, actorType, actorID,
		map[string]any{"project_id": projectID})
	w.WriteHeader(http.StatusNoContent)
}

// --- Git endpoints ---

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
	writeJSON(w, http.StatusOK, map[string]any{"commits": commits, "total": len(commits)})
}

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

type CreateProjectCommitRequest struct {
	Message string   `json:"message"`
	Files   []string `json:"files"`
}

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
	if len(req.Files) > 0 {
		if err := h.GitService.Add(project.LocalPath, req.Files); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to stage files")
			return
		}
	}

	userID := requestUserID(r)
	authorName := "Multica User"
	authorEmail := "user@multica.ai"
	if userID != "" {
		user, err := h.Queries.GetUser(r.Context(), userID)
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
	writeJSON(w, http.StatusCreated, map[string]string{"hash": hash})
}

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
	writeJSON(w, http.StatusOK, map[string]any{"branches": branches})
}

type CreateProjectBranchRequest struct {
	Name string `json:"name"`
}

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
	writeJSON(w, http.StatusCreated, map[string]string{"name": req.Name})
}

type CheckoutProjectBranchRequest struct {
	Branch string `json:"branch"`
}

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
	writeJSON(w, http.StatusOK, map[string]string{"branch": req.Branch})
}

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
	writeJSON(w, http.StatusOK, map[string]any{"diffs": diffs})
}

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
	writeJSON(w, http.StatusOK, map[string]any{"initialized": initialized})
}

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
		ID:          projectID,
		WorkspaceID: workspaceID,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, "project not found")
		return db.LocalProject{}, false
	}
	return project, true
}
