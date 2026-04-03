package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	db "github.com/multica-ai/multica/server/pkg/db/generated"
	"github.com/multica-ai/multica/server/pkg/protocol"
)

type DependencyResponse struct {
	ID               string `json:"id"`
	IssueID          string `json:"issue_id"`
	DependsOnIssueID string `json:"depends_on_issue_id"`
	Type             string `json:"type"`
}

func dependencyToResponse(d db.IssueDependency) DependencyResponse {
	return DependencyResponse{
		ID:               uuidToString(d.ID),
		IssueID:          uuidToString(d.IssueID),
		DependsOnIssueID: uuidToString(d.DependsOnIssueID),
		Type:             d.Type,
	}
}

func (h *Handler) ListIssueDependencies(w http.ResponseWriter, r *http.Request) {
	issueID := chi.URLParam(r, "id")

	deps, err := h.Queries.ListIssueDependencies(r.Context(), parseUUID(issueID))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list dependencies")
		return
	}
	resp := make([]DependencyResponse, len(deps))
	for i, d := range deps {
		resp[i] = dependencyToResponse(d)
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) CreateIssueDependency(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}
	issueID := chi.URLParam(r, "id")
	workspaceID := resolveWorkspaceID(r)

	var body struct {
		DependsOnIssueID string `json:"depends_on_issue_id"`
		Type             string `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.DependsOnIssueID == "" || body.Type == "" {
		writeError(w, http.StatusBadRequest, "depends_on_issue_id and type are required")
		return
	}
	if body.Type != "blocks" && body.Type != "blocked_by" && body.Type != "related" {
		writeError(w, http.StatusBadRequest, "type must be blocks, blocked_by, or related")
		return
	}

	dep, err := h.Queries.CreateIssueDependency(r.Context(), db.CreateIssueDependencyParams{
		IssueID:          parseUUID(issueID),
		DependsOnIssueID: parseUUID(body.DependsOnIssueID),
		Type:             body.Type,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create dependency")
		return
	}

	resp := dependencyToResponse(dep)
	h.publish(protocol.EventIssueUpdated, workspaceID, "member", userID, map[string]any{"issue_id": issueID})
	writeJSON(w, http.StatusCreated, resp)
}

func (h *Handler) DeleteIssueDependency(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}
	depID := chi.URLParam(r, "depId")
	workspaceID := resolveWorkspaceID(r)
	issueID := chi.URLParam(r, "id")

	if err := h.Queries.DeleteIssueDependency(r.Context(), parseUUID(depID)); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete dependency")
		return
	}

	h.publish(protocol.EventIssueUpdated, workspaceID, "member", userID, map[string]any{"issue_id": issueID})
	w.WriteHeader(http.StatusNoContent)
}
