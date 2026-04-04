package handler

// MULTICA-LOCAL: Slock Phase 2 — Task suggestions with approval flow.
// Members and agents can propose tasks in channels. Suggestions can be
// approved (creating an issue + enqueuing agent task) or dismissed.

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/multica-ai/multica/server/internal/logger"
	db "github.com/multica-ai/multica/server/pkg/db/generated"
	"github.com/multica-ai/multica/server/pkg/protocol"
)

// --- Response types ---

type TaskSuggestionResponse struct {
	ID              string  `json:"id"`
	ChannelID       string  `json:"channel_id"`
	WorkspaceID     string  `json:"workspace_id"`
	MessageID       *string `json:"message_id"`
	SuggestedByType string  `json:"suggested_by_type"`
	SuggestedByID   string  `json:"suggested_by_id"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	Priority        string  `json:"priority"`
	AssigneeType    *string `json:"assignee_type"`
	AssigneeID      *string `json:"assignee_id"`
	Status          string  `json:"status"`
	ResolvedByType  *string `json:"resolved_by_type"`
	ResolvedByID    *string `json:"resolved_by_id"`
	IssueID         *string `json:"issue_id"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

func suggestionToResponse(s db.TaskSuggestion) TaskSuggestionResponse {
	return TaskSuggestionResponse{
		ID:              s.ID,
		ChannelID:       s.ChannelID,
		WorkspaceID:     s.WorkspaceID,
		MessageID:       nullStringToPtr(s.MessageID),
		SuggestedByType: s.SuggestedByType,
		SuggestedByID:   s.SuggestedByID,
		Title:           s.Title,
		Description:     s.Description,
		Priority:        s.Priority,
		AssigneeType:    nullStringToPtr(s.AssigneeType),
		AssigneeID:      nullStringToPtr(s.AssigneeID),
		Status:          s.Status,
		ResolvedByType:  nullStringToPtr(s.ResolvedByType),
		ResolvedByID:    nullStringToPtr(s.ResolvedByID),
		IssueID:         nullStringToPtr(s.IssueID),
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
	}
}

// --- Create suggestion ---

type CreateSuggestionRequest struct {
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Priority     string  `json:"priority"`
	AssigneeType *string `json:"assignee_type"`
	AssigneeID   *string `json:"assignee_id"`
}

func (h *Handler) CreateSuggestion(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelId")
	workspaceID := resolveWorkspaceID(r)
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	_, err := h.Queries.GetChannelInWorkspace(r.Context(), db.GetChannelInWorkspaceParams{
		ID:          channelID,
		WorkspaceID: workspaceID,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, "channel not found")
		return
	}

	var req CreateSuggestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	if req.Priority == "" {
		req.Priority = "none"
	}

	actorType, actorID := h.resolveActor(r, userID, workspaceID)

	suggestion, err := h.Queries.CreateTaskSuggestion(r.Context(), db.CreateTaskSuggestionParams{
		ID:              newUUID(),
		ChannelID:       channelID,
		WorkspaceID:     workspaceID,
		SuggestedByType: actorType,
		SuggestedByID:   actorID,
		Title:           req.Title,
		Description:     req.Description,
		Priority:        req.Priority,
		AssigneeType:    ptrToNullString(req.AssigneeType),
		AssigneeID:      ptrToNullString(req.AssigneeID),
	})
	if err != nil {
		slog.Warn("create suggestion failed", append(logger.RequestAttrs(r), "error", err)...)
		writeError(w, http.StatusInternalServerError, "failed to create suggestion")
		return
	}

	// Post a system message in the channel about the suggestion.
	assigneeDesc := ""
	if req.AssigneeID != nil {
		assigneeDesc = " and assigned to " + *req.AssigneeID
	}
	sysMsg, _ := h.Queries.CreateChannelMessage(r.Context(), db.CreateChannelMessageParams{
		ID:          newUUID(),
		ChannelID:   channelID,
		WorkspaceID: workspaceID,
		AuthorType:  actorType,
		AuthorID:    actorID,
		Content:     "Suggested task: **" + req.Title + "**" + assigneeDesc + " (pending approval)",
		Type:        "system",
	})
	// Link the message to the suggestion.
	if sysMsg.ID != "" {
		h.Queries.UpdateTaskSuggestionMessage(r.Context(), db.UpdateTaskSuggestionMessageParams{
			MessageID: sql.NullString{String: sysMsg.ID, Valid: true},
			ID:        suggestion.ID,
		})
	}
	h.Queries.TouchChannel(r.Context(), channelID)

	resp := suggestionToResponse(suggestion)
	h.publish(protocol.EventSuggestionCreated, workspaceID, actorType, actorID, map[string]any{
		"suggestion": resp,
		"channel_id": channelID,
	})
	if sysMsg.ID != "" {
		h.publish(protocol.EventChannelMessageCreated, workspaceID, actorType, actorID, map[string]any{
			"message": channelMessageToResponse(sysMsg),
		})
	}

	writeJSON(w, http.StatusCreated, resp)
}

// --- List suggestions ---

func (h *Handler) ListSuggestions(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelId")
	workspaceID := resolveWorkspaceID(r)

	statusFilter := r.URL.Query().Get("status")

	var suggestions []db.TaskSuggestion
	var err error

	if statusFilter == "pending" {
		suggestions, err = h.Queries.ListPendingTaskSuggestions(r.Context(), db.ListPendingTaskSuggestionsParams{
			ChannelID:   channelID,
			WorkspaceID: workspaceID,
		})
	} else {
		limit := int64(50)
		offset := int64(0)
		if v := r.URL.Query().Get("limit"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n > 0 {
				limit = int64(n)
			}
		}
		if v := r.URL.Query().Get("offset"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n >= 0 {
				offset = int64(n)
			}
		}
		suggestions, err = h.Queries.ListTaskSuggestions(r.Context(), db.ListTaskSuggestionsParams{
			ChannelID:   channelID,
			WorkspaceID: workspaceID,
			Limit:       limit,
			Offset:      offset,
		})
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list suggestions")
		return
	}

	resp := make([]TaskSuggestionResponse, len(suggestions))
	for i, s := range suggestions {
		resp[i] = suggestionToResponse(s)
	}
	writeJSON(w, http.StatusOK, resp)
}

// --- Approve suggestion ---

func (h *Handler) ApproveSuggestion(w http.ResponseWriter, r *http.Request) {
	suggestionID := chi.URLParam(r, "suggestionId")
	workspaceID := resolveWorkspaceID(r)
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	suggestion, err := h.Queries.GetTaskSuggestionInWorkspace(r.Context(), db.GetTaskSuggestionInWorkspaceParams{
		ID:          suggestionID,
		WorkspaceID: workspaceID,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, "suggestion not found")
		return
	}
	if suggestion.Status != "pending" {
		writeError(w, http.StatusConflict, "suggestion already "+suggestion.Status)
		return
	}

	actorType, actorID := h.resolveActor(r, userID, workspaceID)

	// Create issue from suggestion.
	issueNumber, err := h.Queries.IncrementIssueCounter(r.Context(), workspaceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to allocate issue number")
		return
	}

	issue, err := h.Queries.CreateIssue(r.Context(), db.CreateIssueParams{
		ID:           newUUID(),
		WorkspaceID:  workspaceID,
		Title:        suggestion.Title,
		Description:  strToNullString(suggestion.Description),
		Status:       "todo",
		Priority:     suggestion.Priority,
		AssigneeType: suggestion.AssigneeType,
		AssigneeID:   suggestion.AssigneeID,
		CreatorType:  suggestion.SuggestedByType,
		CreatorID:    suggestion.SuggestedByID,
		Number:       issueNumber,
	})
	if err != nil {
		slog.Warn("approve suggestion: create issue failed", append(logger.RequestAttrs(r), "error", err)...)
		writeError(w, http.StatusInternalServerError, "failed to create issue")
		return
	}

	// Mark suggestion as approved.
	updated, err := h.Queries.ApproveTaskSuggestion(r.Context(), db.ApproveTaskSuggestionParams{
		ResolvedByType: sql.NullString{String: actorType, Valid: true},
		ResolvedByID:   sql.NullString{String: actorID, Valid: true},
		IssueID:        sql.NullString{String: issue.ID, Valid: true},
		ID:             suggestionID,
	})
	if err != nil {
		slog.Warn("approve suggestion: update failed", append(logger.RequestAttrs(r), "error", err)...)
		writeError(w, http.StatusInternalServerError, "failed to approve suggestion")
		return
	}

	// Post system message about approval.
	prefix := h.getIssuePrefix(r.Context(), workspaceID)
	identifier := prefix + "-" + strconv.Itoa(int(issue.Number))
	sysMsg, _ := h.Queries.CreateChannelMessage(r.Context(), db.CreateChannelMessageParams{
		ID:          newUUID(),
		ChannelID:   suggestion.ChannelID,
		WorkspaceID: workspaceID,
		AuthorType:  actorType,
		AuthorID:    actorID,
		Content:     "Approved task suggestion → [" + identifier + "](mention://issue/" + issue.ID + "): " + issue.Title,
		Type:        "issue_created",
		IssueID:     strToNullString(issue.ID),
	})
	h.Queries.TouchChannel(r.Context(), suggestion.ChannelID)

	// Publish events.
	issueResp := issueToResponse(issue, prefix)
	h.publish(protocol.EventIssueCreated, workspaceID, actorType, actorID, map[string]any{"issue": issueResp})

	resp := suggestionToResponse(updated)
	h.publish(protocol.EventSuggestionApproved, workspaceID, actorType, actorID, map[string]any{
		"suggestion": resp,
		"channel_id": suggestion.ChannelID,
		"issue":      issueResp,
	})
	if sysMsg.ID != "" {
		h.publish(protocol.EventChannelMessageCreated, workspaceID, actorType, actorID, map[string]any{
			"message": channelMessageToResponse(sysMsg),
		})
	}

	// If assigned to an agent, enqueue task.
	if issue.AssigneeType.Valid && issue.AssigneeType.String == "agent" {
		if _, err := h.TaskService.EnqueueTaskForIssue(r.Context(), issue, sql.NullString{}); err != nil {
			slog.Warn("enqueue agent task from approved suggestion failed", "issue_id", issue.ID, "error", err)
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"suggestion": resp,
		"issue":      issueResp,
	})
}

// --- Dismiss suggestion ---

func (h *Handler) DismissSuggestion(w http.ResponseWriter, r *http.Request) {
	suggestionID := chi.URLParam(r, "suggestionId")
	workspaceID := resolveWorkspaceID(r)
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	suggestion, err := h.Queries.GetTaskSuggestionInWorkspace(r.Context(), db.GetTaskSuggestionInWorkspaceParams{
		ID:          suggestionID,
		WorkspaceID: workspaceID,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, "suggestion not found")
		return
	}
	if suggestion.Status != "pending" {
		writeError(w, http.StatusConflict, "suggestion already "+suggestion.Status)
		return
	}

	actorType, actorID := h.resolveActor(r, userID, workspaceID)

	updated, err := h.Queries.DismissTaskSuggestion(r.Context(), db.DismissTaskSuggestionParams{
		ResolvedByType: sql.NullString{String: actorType, Valid: true},
		ResolvedByID:   sql.NullString{String: actorID, Valid: true},
		ID:             suggestionID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to dismiss suggestion")
		return
	}

	// Post system message.
	h.Queries.CreateChannelMessage(r.Context(), db.CreateChannelMessageParams{
		ID:          newUUID(),
		ChannelID:   suggestion.ChannelID,
		WorkspaceID: workspaceID,
		AuthorType:  actorType,
		AuthorID:    actorID,
		Content:     "Dismissed task suggestion: " + suggestion.Title,
		Type:        "system",
	})
	h.Queries.TouchChannel(r.Context(), suggestion.ChannelID)

	resp := suggestionToResponse(updated)
	h.publish(protocol.EventSuggestionDismissed, workspaceID, actorType, actorID, map[string]any{
		"suggestion": resp,
		"channel_id": suggestion.ChannelID,
	})

	writeJSON(w, http.StatusOK, resp)
}
