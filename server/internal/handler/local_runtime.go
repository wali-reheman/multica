package handler

// MULTICA-LOCAL: Stage 4 — Local agent runtime management, task execution, and skills.

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"os/exec"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/multica-ai/multica/server/internal/service"
	db "github.com/multica-ai/multica/server/pkg/db/generated"
	"github.com/multica-ai/multica/server/pkg/protocol"
)

// ---------------------------------------------------------------------------
// 4.1 — Local Agent Runtime Manager
// ---------------------------------------------------------------------------

// DetectLocalAgents scans the system for available agent CLIs and returns results.
// POST /api/local/agents/detect
func (h *Handler) DetectLocalAgents(w http.ResponseWriter, r *http.Request) {
	workspaceID := resolveWorkspaceID(r)
	if workspaceID == "" {
		writeError(w, http.StatusBadRequest, "workspace_id is required")
		return
	}

	svc := service.NewLocalRuntimeService(h.Queries)
	agents, err := svc.DetectAgents(r.Context(), workspaceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "detection failed: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"agents": agents})
}

// ListLocalAgents returns previously detected agent CLI configurations.
// GET /api/local/agents
func (h *Handler) ListLocalAgents(w http.ResponseWriter, r *http.Request) {
	workspaceID := resolveWorkspaceID(r)
	if workspaceID == "" {
		writeError(w, http.StatusBadRequest, "workspace_id is required")
		return
	}

	configs, err := h.Queries.ListLocalAgentConfigs(r.Context(), workspaceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "list failed")
		return
	}

	result := make([]service.DetectedAgent, 0, len(configs))
	for _, cfg := range configs {
		da := service.DetectedAgent{
			Provider: cfg.Provider,
			Path:     cfg.CliPath,
			Version:  cfg.Version,
			Status:   cfg.Status,
			IsCustom: cfg.IsCustomPath,
		}
		if cfg.HealthError.Valid {
			da.Error = cfg.HealthError.String
		}
		result = append(result, da)
	}

	writeJSON(w, http.StatusOK, map[string]any{"agents": result})
}

// SetLocalAgentPath updates the CLI path for a specific provider.
// PUT /api/local/agents/{provider}/path
func (h *Handler) SetLocalAgentPath(w http.ResponseWriter, r *http.Request) {
	workspaceID := resolveWorkspaceID(r)
	if workspaceID == "" {
		writeError(w, http.StatusBadRequest, "workspace_id is required")
		return
	}

	provider := chi.URLParam(r, "provider")
	if provider == "" {
		writeError(w, http.StatusBadRequest, "provider is required")
		return
	}

	var req struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.Path = strings.TrimSpace(req.Path)
	if req.Path == "" {
		writeError(w, http.StatusBadRequest, "path is required")
		return
	}

	if _, err := exec.LookPath(req.Path); err != nil {
		writeError(w, http.StatusBadRequest, "path not found or not executable")
		return
	}

	svc := service.NewLocalRuntimeService(h.Queries)
	agent, err := svc.SetCustomPath(r.Context(), workspaceID, provider, req.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, agent)
}

// HealthCheckLocalAgents runs health checks on all configured agents.
// POST /api/local/agents/health-check
func (h *Handler) HealthCheckLocalAgents(w http.ResponseWriter, r *http.Request) {
	workspaceID := resolveWorkspaceID(r)
	if workspaceID == "" {
		writeError(w, http.StatusBadRequest, "workspace_id is required")
		return
	}

	svc := service.NewLocalRuntimeService(h.Queries)
	agents, err := svc.HealthCheckAll(r.Context(), workspaceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "health check failed")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"agents": agents})
}

// ---------------------------------------------------------------------------
// 4.2 — Local Task Execution
// ---------------------------------------------------------------------------

// RunAgentOnIssue triggers a local agent task execution for an issue.
// POST /api/issues/{id}/run-agent
func (h *Handler) RunAgentOnIssue(w http.ResponseWriter, r *http.Request) {
	issueID := chi.URLParam(r, "id")
	issue, ok := h.loadIssueForUser(w, r, issueID)
	if !ok {
		return
	}

	var req struct {
		AgentID  string `json:"agent_id"`
		Provider string `json:"provider"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	agentID := req.AgentID
	if agentID == "" {
		if !issue.AssigneeID.Valid || issue.AssigneeType.String != "agent" {
			writeError(w, http.StatusBadRequest, "issue has no agent assignee — specify agent_id")
			return
		}
		agentID = issue.AssigneeID.String
	}

	agentRow, err := h.Queries.GetAgent(r.Context(), agentID)
	if err != nil {
		writeError(w, http.StatusNotFound, "agent not found")
		return
	}
	if agentRow.ArchivedAt.Valid {
		writeError(w, http.StatusBadRequest, "agent is archived")
		return
	}

	hasActive, _ := h.Queries.HasActiveTaskForIssue(r.Context(), issue.ID)
	if hasActive {
		writeError(w, http.StatusConflict, "a task is already active for this issue")
		return
	}

	if agentRow.RuntimeID == "" {
		writeError(w, http.StatusBadRequest, "agent has no runtime configured")
		return
	}

	task, err := h.TaskService.EnqueueTaskForIssue(r.Context(), issue)
	if err != nil {
		slog.Error("run agent: enqueue failed", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to create task: "+err.Error())
		return
	}

	executor := service.NewLocalExecutor(
		service.LocalExecutorConfig{MaxConcurrentTasks: 3},
		h.Queries, h.Hub, h.Bus, h.TaskService,
	)
	if err := executor.ExecuteTask(r.Context(), task.ID); err != nil {
		slog.Error("run agent: execute failed", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to execute task: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"task_id":  task.ID,
		"agent_id": agentID,
		"issue_id": issueID,
		"status":   "running",
	})
}

// GetIssueDiff returns the git diff of changes made by an agent for an issue.
// GET /api/issues/{id}/agent-diff
func (h *Handler) GetIssueDiff(w http.ResponseWriter, r *http.Request) {
	issueID := chi.URLParam(r, "id")
	issue, ok := h.loadIssueForUser(w, r, issueID)
	if !ok {
		return
	}

	tasks, err := h.Queries.ListTasksByIssue(r.Context(), issue.ID)
	if err != nil || len(tasks) == 0 {
		writeJSON(w, http.StatusOK, map[string]any{"diff": "", "has_changes": false})
		return
	}

	var workDir string
	for _, t := range tasks {
		if t.Status == "completed" && t.WorkDir.Valid {
			workDir = t.WorkDir.String
			break
		}
	}

	if workDir == "" {
		writeJSON(w, http.StatusOK, map[string]any{"diff": "", "has_changes": false})
		return
	}

	diff := getGitDiff(workDir)
	writeJSON(w, http.StatusOK, map[string]any{
		"diff":        diff,
		"has_changes": diff != "",
		"work_dir":    workDir,
	})
}

func getGitDiff(dir string) string {
	cmd := exec.Command("git", "diff", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		cmd2 := exec.Command("git", "diff")
		cmd2.Dir = dir
		out2, _ := cmd2.Output()
		return string(out2)
	}
	return string(out)
}

// CommitAgentChanges commits changes made by an agent with an auto-generated message.
// POST /api/issues/{id}/agent-commit
func (h *Handler) CommitAgentChanges(w http.ResponseWriter, r *http.Request) {
	issueID := chi.URLParam(r, "id")
	issue, ok := h.loadIssueForUser(w, r, issueID)
	if !ok {
		return
	}

	var req struct {
		Message string `json:"message"`
		WorkDir string `json:"work_dir"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.WorkDir == "" {
		tasks, err := h.Queries.ListTasksByIssue(r.Context(), issue.ID)
		if err != nil || len(tasks) == 0 {
			writeError(w, http.StatusBadRequest, "no completed tasks found")
			return
		}
		for _, t := range tasks {
			if t.Status == "completed" && t.WorkDir.Valid {
				req.WorkDir = t.WorkDir.String
				break
			}
		}
	}

	if req.WorkDir == "" {
		writeError(w, http.StatusBadRequest, "work_dir not found")
		return
	}

	message := req.Message
	if message == "" {
		message = "feat: agent changes for " + issue.Title
	}

	addCmd := exec.Command("git", "add", "-A")
	addCmd.Dir = req.WorkDir
	if out, err := addCmd.CombinedOutput(); err != nil {
		writeError(w, http.StatusInternalServerError, "git add failed: "+string(out))
		return
	}

	commitCmd := exec.Command("git", "commit", "-m", message)
	commitCmd.Dir = req.WorkDir
	commitOut, err := commitCmd.CombinedOutput()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "git commit failed: "+string(commitOut))
		return
	}

	h.publish(protocol.EventIssueUpdated, issue.WorkspaceID, "system", "",
		map[string]any{
			"issue_id": issueID,
			"message":  "Agent changes committed",
		})

	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "committed",
		"message": message,
		"output":  string(commitOut),
	})
}

// ---------------------------------------------------------------------------
// 4.5 — Local Skills Management
// ---------------------------------------------------------------------------

// ListLocalSkills returns local skills (global + workspace + project).
// GET /api/local/skills
func (h *Handler) ListLocalSkills(w http.ResponseWriter, r *http.Request) {
	workspaceID := resolveWorkspaceID(r)
	if workspaceID == "" {
		writeError(w, http.StatusBadRequest, "workspace_id is required")
		return
	}

	projectPath := r.URL.Query().Get("project_path")

	skills, err := h.Queries.ListLocalSkills(r.Context(), db.ListLocalSkillsParams{
		WorkspaceID: workspaceID,
		ProjectPath: sql.NullString{String: projectPath, Valid: projectPath != ""},
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "list failed")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"skills": skills})
}

// CreateLocalSkill creates a new local skill.
// POST /api/local/skills
func (h *Handler) CreateLocalSkill(w http.ResponseWriter, r *http.Request) {
	workspaceID := resolveWorkspaceID(r)

	var req struct {
		ProjectPath string `json:"project_path"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Content     string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	skill, err := h.Queries.CreateLocalSkill(r.Context(), db.CreateLocalSkillParams{
		ID:          newUUID(),
		WorkspaceID: sql.NullString{String: workspaceID, Valid: workspaceID != ""},
		ProjectPath: sql.NullString{String: req.ProjectPath, Valid: req.ProjectPath != ""},
		Name:        req.Name,
		Description: req.Description,
		Content:     req.Content,
		IsDefault:   false,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "create failed")
		return
	}

	writeJSON(w, http.StatusCreated, skill)
}

// UpdateLocalSkill updates a local skill.
// PUT /api/local/skills/{id}
func (h *Handler) UpdateLocalSkill(w http.ResponseWriter, r *http.Request) {
	skillID := chi.URLParam(r, "id")

	var req struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Content     *string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	skill, err := h.Queries.UpdateLocalSkill(r.Context(), db.UpdateLocalSkillParams{
		ID:          skillID,
		Name:        ptrToNullString(req.Name),
		Description: ptrToNullString(req.Description),
		Content:     ptrToNullString(req.Content),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "update failed")
		return
	}

	writeJSON(w, http.StatusOK, skill)
}

// DeleteLocalSkill removes a local skill.
// DELETE /api/local/skills/{id}
func (h *Handler) DeleteLocalSkill(w http.ResponseWriter, r *http.Request) {
	skillID := chi.URLParam(r, "id")

	if err := h.Queries.DeleteLocalSkill(r.Context(), skillID); err != nil {
		writeError(w, http.StatusInternalServerError, "delete failed")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
