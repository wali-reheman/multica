package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	db "github.com/multica-ai/multica/server/pkg/db/generated"
	"github.com/multica-ai/multica/server/pkg/protocol"
)

type LabelResponse struct {
	ID          string `json:"id"`
	WorkspaceID string `json:"workspace_id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
}

func labelToResponse(l db.IssueLabel) LabelResponse {
	return LabelResponse{
		ID:          uuidToString(l.ID),
		WorkspaceID: uuidToString(l.WorkspaceID),
		Name:        l.Name,
		Color:       l.Color,
	}
}

func (h *Handler) ListLabels(w http.ResponseWriter, r *http.Request) {
	workspaceID := resolveWorkspaceID(r)
	labels, err := h.Queries.ListLabelsByWorkspace(r.Context(), parseUUID(workspaceID))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list labels")
		return
	}
	resp := make([]LabelResponse, len(labels))
	for i, l := range labels {
		resp[i] = labelToResponse(l)
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) CreateLabel(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}
	workspaceID := resolveWorkspaceID(r)

	var body struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if body.Color == "" {
		body.Color = "#6b7280"
	}

	label, err := h.Queries.CreateLabel(r.Context(), db.CreateLabelParams{
		WorkspaceID: parseUUID(workspaceID),
		Name:        body.Name,
		Color:       body.Color,
	})
	if err != nil {
		if isUniqueViolation(err) {
			writeError(w, http.StatusConflict, "label with this name already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create label")
		return
	}

	resp := labelToResponse(label)
	h.publish(protocol.EventLabelCreated, workspaceID, "member", userID, map[string]any{"label": resp})
	writeJSON(w, http.StatusCreated, resp)
}

func (h *Handler) UpdateLabel(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}
	labelID := chi.URLParam(r, "labelId")
	workspaceID := resolveWorkspaceID(r)

	var body struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	label, err := h.Queries.UpdateLabel(r.Context(), db.UpdateLabelParams{
		ID:    parseUUID(labelID),
		Name:  body.Name,
		Color: body.Color,
	})
	if err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, "label not found")
			return
		}
		if isUniqueViolation(err) {
			writeError(w, http.StatusConflict, "label with this name already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update label")
		return
	}

	resp := labelToResponse(label)
	h.publish(protocol.EventLabelUpdated, workspaceID, "member", userID, map[string]any{"label": resp})
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) DeleteLabel(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}
	labelID := chi.URLParam(r, "labelId")
	workspaceID := resolveWorkspaceID(r)

	if err := h.Queries.DeleteLabel(r.Context(), parseUUID(labelID)); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete label")
		return
	}

	h.publish(protocol.EventLabelDeleted, workspaceID, "member", userID, map[string]any{"label_id": labelID})
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListIssueLabels(w http.ResponseWriter, r *http.Request) {
	issueID := chi.URLParam(r, "id")

	labels, err := h.Queries.ListIssueLabels(r.Context(), parseUUID(issueID))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list issue labels")
		return
	}
	resp := make([]LabelResponse, len(labels))
	for i, l := range labels {
		resp[i] = labelToResponse(l)
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) AddIssueLabel(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}
	issueID := chi.URLParam(r, "id")
	workspaceID := resolveWorkspaceID(r)

	var body struct {
		LabelID string `json:"label_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.LabelID == "" {
		writeError(w, http.StatusBadRequest, "label_id is required")
		return
	}

	if err := h.Queries.AddIssueLabel(r.Context(), db.AddIssueLabelParams{
		IssueID: parseUUID(issueID),
		LabelID: parseUUID(body.LabelID),
	}); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to add label")
		return
	}

	h.publish(protocol.EventIssueUpdated, workspaceID, "member", userID, map[string]any{"issue_id": issueID})
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RemoveIssueLabel(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}
	issueID := chi.URLParam(r, "id")
	workspaceID := resolveWorkspaceID(r)

	var body struct {
		LabelID string `json:"label_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.Queries.RemoveIssueLabel(r.Context(), db.RemoveIssueLabelParams{
		IssueID: parseUUID(issueID),
		LabelID: parseUUID(body.LabelID),
	}); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to remove label")
		return
	}

	h.publish(protocol.EventIssueUpdated, workspaceID, "member", userID, map[string]any{"issue_id": issueID})
	w.WriteHeader(http.StatusNoContent)
}
