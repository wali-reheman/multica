package handler

// MULTICA-LOCAL: Slock — Slack-like channels for group chat between agents and members.

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

type ChannelResponse struct {
	ID            string                 `json:"id"`
	WorkspaceID   string                 `json:"workspace_id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Type          string                 `json:"type"`
	CreatedByType string                 `json:"created_by_type"`
	CreatedByID   string                 `json:"created_by_id"`
	ArchivedAt    *string                `json:"archived_at"`
	CreatedAt     string                 `json:"created_at"`
	UpdatedAt     string                 `json:"updated_at"`
	Members       []ChannelMemberResponse `json:"members,omitempty"`
}

type ChannelMemberResponse struct {
	ChannelID  string `json:"channel_id"`
	MemberType string `json:"member_type"`
	MemberID   string `json:"member_id"`
	Role       string `json:"role"`
	JoinedAt   string `json:"joined_at"`
}

type ChannelMessageResponse struct {
	ID         string  `json:"id"`
	ChannelID  string  `json:"channel_id"`
	AuthorType string  `json:"author_type"`
	AuthorID   string  `json:"author_id"`
	Content    string  `json:"content"`
	Type       string  `json:"type"`
	ParentID   *string `json:"parent_id"`
	IssueID    *string `json:"issue_id"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

func channelToResponse(c db.Channel) ChannelResponse {
	return ChannelResponse{
		ID:            c.ID,
		WorkspaceID:   c.WorkspaceID,
		Name:          c.Name,
		Description:   c.Description,
		Type:          c.Type,
		CreatedByType: c.CreatedByType,
		CreatedByID:   c.CreatedByID,
		ArchivedAt:    nullStringToPtr(c.ArchivedAt),
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
	}
}

func channelMemberToResponse(m db.ChannelMember) ChannelMemberResponse {
	return ChannelMemberResponse{
		ChannelID:  m.ChannelID,
		MemberType: m.MemberType,
		MemberID:   m.MemberID,
		Role:       m.Role,
		JoinedAt:   m.JoinedAt,
	}
}

func channelMessageToResponse(m db.ChannelMessage) ChannelMessageResponse {
	return ChannelMessageResponse{
		ID:         m.ID,
		ChannelID:  m.ChannelID,
		AuthorType: m.AuthorType,
		AuthorID:   m.AuthorID,
		Content:    m.Content,
		Type:       m.Type,
		ParentID:   nullStringToPtr(m.ParentID),
		IssueID:    nullStringToPtr(m.IssueID),
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

// --- Channel CRUD ---

type CreateChannelRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	MemberIDs   []string `json:"member_ids"`
	AgentIDs    []string `json:"agent_ids"`
}

func (h *Handler) CreateChannel(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}
	workspaceID := resolveWorkspaceID(r)

	var req CreateChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if req.Type == "" {
		req.Type = "group"
	}

	actorType, actorID := h.resolveActor(r, userID, workspaceID)

	channel, err := h.Queries.CreateChannel(r.Context(), db.CreateChannelParams{
		ID:            newUUID(),
		WorkspaceID:   workspaceID,
		Name:          req.Name,
		Description:   req.Description,
		Type:          req.Type,
		CreatedByType: actorType,
		CreatedByID:   actorID,
	})
	if err != nil {
		if isUniqueViolation(err) {
			writeError(w, http.StatusConflict, "channel name already exists")
			return
		}
		slog.Warn("create channel failed", append(logger.RequestAttrs(r), "error", err)...)
		writeError(w, http.StatusInternalServerError, "failed to create channel")
		return
	}

	// Add creator as channel owner.
	h.Queries.AddChannelMember(r.Context(), db.AddChannelMemberParams{
		ChannelID:  channel.ID,
		MemberType: actorType,
		MemberID:   actorID,
		Role:       "owner",
	})

	// Add requested members.
	for _, mid := range req.MemberIDs {
		h.Queries.AddChannelMember(r.Context(), db.AddChannelMemberParams{
			ChannelID:  channel.ID,
			MemberType: "member",
			MemberID:   mid,
			Role:       "member",
		})
	}
	// Add requested agents.
	for _, aid := range req.AgentIDs {
		h.Queries.AddChannelMember(r.Context(), db.AddChannelMemberParams{
			ChannelID:  channel.ID,
			MemberType: "agent",
			MemberID:   aid,
			Role:       "member",
		})
	}

	// Fetch members for response.
	members, _ := h.Queries.ListChannelMembers(r.Context(), channel.ID)
	resp := channelToResponse(channel)
	resp.Members = make([]ChannelMemberResponse, len(members))
	for i, m := range members {
		resp.Members[i] = channelMemberToResponse(m)
	}

	slog.Info("channel created", append(logger.RequestAttrs(r), "channel_id", channel.ID)...)
	h.publish(protocol.EventChannelCreated, workspaceID, actorType, actorID, map[string]any{"channel": resp})
	writeJSON(w, http.StatusCreated, resp)
}

func (h *Handler) ListChannels(w http.ResponseWriter, r *http.Request) {
	workspaceID := resolveWorkspaceID(r)

	channels, err := h.Queries.ListChannels(r.Context(), workspaceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list channels")
		return
	}

	resp := make([]ChannelResponse, len(channels))
	for i, c := range channels {
		resp[i] = channelToResponse(c)
		members, _ := h.Queries.ListChannelMembers(r.Context(), c.ID)
		resp[i].Members = make([]ChannelMemberResponse, len(members))
		for j, m := range members {
			resp[i].Members[j] = channelMemberToResponse(m)
		}
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetChannel(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelId")
	workspaceID := resolveWorkspaceID(r)

	channel, err := h.Queries.GetChannelInWorkspace(r.Context(), db.GetChannelInWorkspaceParams{
		ID:          channelID,
		WorkspaceID: workspaceID,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, "channel not found")
		return
	}

	resp := channelToResponse(channel)
	members, _ := h.Queries.ListChannelMembers(r.Context(), channel.ID)
	resp.Members = make([]ChannelMemberResponse, len(members))
	for i, m := range members {
		resp.Members[i] = channelMemberToResponse(m)
	}

	writeJSON(w, http.StatusOK, resp)
}

type UpdateChannelRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func (h *Handler) UpdateChannel(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelId")
	workspaceID := resolveWorkspaceID(r)
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	channel, err := h.Queries.GetChannelInWorkspace(r.Context(), db.GetChannelInWorkspaceParams{
		ID:          channelID,
		WorkspaceID: workspaceID,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, "channel not found")
		return
	}

	var req UpdateChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	actorType, actorID := h.resolveActor(r, userID, workspaceID)

	updated, err := h.Queries.UpdateChannel(r.Context(), db.UpdateChannelParams{
		Name:        ptrToNullString(req.Name),
		Description: ptrToNullString(req.Description),
		ID:          channel.ID,
	})
	if err != nil {
		slog.Warn("update channel failed", append(logger.RequestAttrs(r), "error", err)...)
		writeError(w, http.StatusInternalServerError, "failed to update channel")
		return
	}

	resp := channelToResponse(updated)
	h.publish(protocol.EventChannelUpdated, workspaceID, actorType, actorID, map[string]any{"channel": resp})
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) DeleteChannel(w http.ResponseWriter, r *http.Request) {
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

	actorType, actorID := h.resolveActor(r, userID, workspaceID)

	if err := h.Queries.DeleteChannel(r.Context(), channelID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete channel")
		return
	}

	h.publish(protocol.EventChannelDeleted, workspaceID, actorType, actorID, map[string]any{"channel_id": channelID})
	w.WriteHeader(http.StatusNoContent)
}

// --- Channel members ---

type AddChannelMemberRequest struct {
	MemberType string `json:"member_type"`
	MemberID   string `json:"member_id"`
}

func (h *Handler) AddChannelMember(w http.ResponseWriter, r *http.Request) {
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

	var req AddChannelMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.MemberType == "" || req.MemberID == "" {
		writeError(w, http.StatusBadRequest, "member_type and member_id are required")
		return
	}

	actorType, actorID := h.resolveActor(r, userID, workspaceID)

	if err := h.Queries.AddChannelMember(r.Context(), db.AddChannelMemberParams{
		ChannelID:  channelID,
		MemberType: req.MemberType,
		MemberID:   req.MemberID,
		Role:       "member",
	}); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to add member")
		return
	}

	// Post system message.
	h.Queries.CreateChannelMessage(r.Context(), db.CreateChannelMessageParams{
		ID:          newUUID(),
		ChannelID:   channelID,
		WorkspaceID: workspaceID,
		AuthorType:  actorType,
		AuthorID:    actorID,
		Content:     req.MemberType + " joined the channel",
		Type:        "system",
	})
	h.Queries.TouchChannel(r.Context(), channelID)

	h.publish(protocol.EventChannelMemberAdded, workspaceID, actorType, actorID, map[string]any{
		"channel_id":  channelID,
		"member_type": req.MemberType,
		"member_id":   req.MemberID,
	})
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RemoveChannelMember(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelId")
	workspaceID := resolveWorkspaceID(r)
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	var req AddChannelMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	actorType, actorID := h.resolveActor(r, userID, workspaceID)

	if err := h.Queries.RemoveChannelMember(r.Context(), db.RemoveChannelMemberParams{
		ChannelID:  channelID,
		MemberType: req.MemberType,
		MemberID:   req.MemberID,
	}); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to remove member")
		return
	}

	h.publish(protocol.EventChannelMemberRemoved, workspaceID, actorType, actorID, map[string]any{
		"channel_id":  channelID,
		"member_type": req.MemberType,
		"member_id":   req.MemberID,
	})
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListChannelMembers(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelId")

	members, err := h.Queries.ListChannelMembers(r.Context(), channelID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list members")
		return
	}

	resp := make([]ChannelMemberResponse, len(members))
	for i, m := range members {
		resp[i] = channelMemberToResponse(m)
	}
	writeJSON(w, http.StatusOK, resp)
}

// --- Channel messages ---

type CreateChannelMessageRequest struct {
	Content  string  `json:"content"`
	ParentID *string `json:"parent_id"`
}

func (h *Handler) CreateChannelMessage(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelId")
	workspaceID := resolveWorkspaceID(r)
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	channel, err := h.Queries.GetChannelInWorkspace(r.Context(), db.GetChannelInWorkspaceParams{
		ID:          channelID,
		WorkspaceID: workspaceID,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, "channel not found")
		return
	}

	var req CreateChannelMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Content == "" {
		writeError(w, http.StatusBadRequest, "content is required")
		return
	}

	actorType, actorID := h.resolveActor(r, userID, workspaceID)

	msg, err := h.Queries.CreateChannelMessage(r.Context(), db.CreateChannelMessageParams{
		ID:          newUUID(),
		ChannelID:   channel.ID,
		WorkspaceID: workspaceID,
		AuthorType:  actorType,
		AuthorID:    actorID,
		Content:     req.Content,
		Type:        "message",
		ParentID:    ptrToNullString(req.ParentID),
	})
	if err != nil {
		slog.Warn("create channel message failed", append(logger.RequestAttrs(r), "error", err)...)
		writeError(w, http.StatusInternalServerError, "failed to send message")
		return
	}

	// Touch channel to update ordering.
	h.Queries.TouchChannel(r.Context(), channel.ID)

	resp := channelMessageToResponse(msg)
	h.publish(protocol.EventChannelMessageCreated, workspaceID, actorType, actorID, map[string]any{
		"message":      resp,
		"channel_name": channel.Name,
	})

	// Trigger @mentioned agents in the channel message.
	h.enqueueChannelMentionedAgentTasks(r, channel, msg, actorType, actorID)

	writeJSON(w, http.StatusCreated, resp)
}

func (h *Handler) ListChannelMessages(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelId")
	workspaceID := resolveWorkspaceID(r)

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

	var messages []db.ChannelMessage
	var err error

	if since := r.URL.Query().Get("since"); since != "" {
		messages, err = h.Queries.ListChannelMessagesSince(r.Context(), db.ListChannelMessagesSinceParams{
			ChannelID:   channelID,
			WorkspaceID: workspaceID,
			CreatedAt:   since,
			Limit:       limit,
		})
	} else {
		messages, err = h.Queries.ListChannelMessages(r.Context(), db.ListChannelMessagesParams{
			ChannelID:   channelID,
			WorkspaceID: workspaceID,
			Limit:       limit,
			Offset:      offset,
		})
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list messages")
		return
	}

	resp := make([]ChannelMessageResponse, len(messages))
	for i, m := range messages {
		resp[i] = channelMessageToResponse(m)
	}

	// Include total count header.
	total, countErr := h.Queries.CountChannelMessages(r.Context(), db.CountChannelMessagesParams{
		ChannelID:   channelID,
		WorkspaceID: workspaceID,
	})
	if countErr == nil {
		w.Header().Set("X-Total-Count", strconv.FormatInt(total, 10))
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) UpdateChannelMessage(w http.ResponseWriter, r *http.Request) {
	messageID := chi.URLParam(r, "messageId")
	workspaceID := resolveWorkspaceID(r)
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	existing, err := h.Queries.GetChannelMessage(r.Context(), messageID)
	if err != nil {
		writeError(w, http.StatusNotFound, "message not found")
		return
	}

	actorType, actorID := h.resolveActor(r, userID, workspaceID)
	if existing.AuthorType != actorType || existing.AuthorID != actorID {
		writeError(w, http.StatusForbidden, "only author can edit message")
		return
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Content == "" {
		writeError(w, http.StatusBadRequest, "content is required")
		return
	}

	updated, err := h.Queries.UpdateChannelMessage(r.Context(), db.UpdateChannelMessageParams{
		Content: req.Content,
		ID:      messageID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update message")
		return
	}

	resp := channelMessageToResponse(updated)
	h.publish(protocol.EventChannelMessageUpdated, workspaceID, actorType, actorID, map[string]any{"message": resp})
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) DeleteChannelMessage(w http.ResponseWriter, r *http.Request) {
	messageID := chi.URLParam(r, "messageId")
	workspaceID := resolveWorkspaceID(r)
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	existing, err := h.Queries.GetChannelMessage(r.Context(), messageID)
	if err != nil {
		writeError(w, http.StatusNotFound, "message not found")
		return
	}

	actorType, actorID := h.resolveActor(r, userID, workspaceID)
	if existing.AuthorType != actorType || existing.AuthorID != actorID {
		member, mok := h.workspaceMember(w, r, workspaceID)
		if !mok {
			return
		}
		if !roleAllowed(member.Role, "owner", "admin") {
			writeError(w, http.StatusForbidden, "only author or admin can delete message")
			return
		}
	}

	if err := h.Queries.DeleteChannelMessage(r.Context(), messageID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete message")
		return
	}

	h.publish(protocol.EventChannelMessageDeleted, workspaceID, actorType, actorID, map[string]any{
		"message_id": messageID,
		"channel_id": existing.ChannelID,
	})
	w.WriteHeader(http.StatusNoContent)
}

// --- Create issue from channel message ---

type CreateIssueFromMessageRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Priority    string  `json:"priority"`
	AssigneeID  *string `json:"assignee_id"`
}

func (h *Handler) CreateIssueFromChannelMessage(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelId")
	workspaceID := resolveWorkspaceID(r)
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	channel, err := h.Queries.GetChannelInWorkspace(r.Context(), db.GetChannelInWorkspaceParams{
		ID:          channelID,
		WorkspaceID: workspaceID,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, "channel not found")
		return
	}

	var req CreateIssueFromMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	actorType, actorID := h.resolveActor(r, userID, workspaceID)

	if req.Priority == "" {
		req.Priority = "none"
	}

	// Allocate issue number.
	issueNumber, err := h.Queries.IncrementIssueCounter(r.Context(), workspaceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to allocate issue number")
		return
	}

	var assigneeType, assigneeID sql.NullString
	if req.AssigneeID != nil {
		// Determine if assignee is an agent or member.
		if _, agentErr := h.Queries.GetAgent(r.Context(), *req.AssigneeID); agentErr == nil {
			assigneeType = strToNullString("agent")
		} else {
			assigneeType = strToNullString("member")
		}
		assigneeID = strToNullString(*req.AssigneeID)
	}

	issue, err := h.Queries.CreateIssue(r.Context(), db.CreateIssueParams{
		ID:           newUUID(),
		WorkspaceID:  workspaceID,
		Title:        req.Title,
		Description:  ptrToNullString(req.Description),
		Status:       "todo",
		Priority:     req.Priority,
		AssigneeType: assigneeType,
		AssigneeID:   assigneeID,
		CreatorType:  actorType,
		CreatorID:    actorID,
		Number:       issueNumber,
	})
	if err != nil {
		slog.Warn("create issue from channel failed", append(logger.RequestAttrs(r), "error", err)...)
		writeError(w, http.StatusInternalServerError, "failed to create issue")
		return
	}

	// Post a system message in the channel linking to the new issue.
	prefix := h.getIssuePrefix(r.Context(), workspaceID)
	identifier := prefix + "-" + strconv.Itoa(int(issue.Number))
	sysMsg, _ := h.Queries.CreateChannelMessage(r.Context(), db.CreateChannelMessageParams{
		ID:          newUUID(),
		ChannelID:   channel.ID,
		WorkspaceID: workspaceID,
		AuthorType:  actorType,
		AuthorID:    actorID,
		Content:     "Created issue [" + identifier + "](mention://issue/" + issue.ID + "): " + issue.Title,
		Type:        "issue_created",
		IssueID:     strToNullString(issue.ID),
	})
	h.Queries.TouchChannel(r.Context(), channel.ID)

	// Publish issue created event.
	issueResp := issueToResponse(issue, prefix)
	h.publish(protocol.EventIssueCreated, workspaceID, actorType, actorID, map[string]any{"issue": issueResp})

	// Publish channel message event.
	if sysMsg.ID != "" {
		h.publish(protocol.EventChannelMessageCreated, workspaceID, actorType, actorID, map[string]any{
			"message":      channelMessageToResponse(sysMsg),
			"channel_name": channel.Name,
		})
	}

	// If assigned to an agent, enqueue task.
	if issue.AssigneeType.Valid && issue.AssigneeType.String == "agent" {
		if _, err := h.TaskService.EnqueueTaskForIssue(r.Context(), issue, sql.NullString{}); err != nil {
			slog.Warn("enqueue agent task from channel failed", "issue_id", issue.ID, "error", err)
		}
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"issue":   issueResp,
		"message": channelMessageToResponse(sysMsg),
	})
}

// enqueueChannelMentionedAgentTasks checks for @agent mentions in channel messages
// and enqueues tasks for mentioned agents. Unlike issue comments, channel messages
// don't have an "assigned agent" so all mentioned agents are triggered.
func (h *Handler) enqueueChannelMentionedAgentTasks(r *http.Request, channel db.Channel, msg db.ChannelMessage, authorType, authorID string) {
	// Parse @mentions from message content.
	mentions := parseAgentMentionsFromContent(msg.Content)
	for _, agentID := range mentions {
		if authorType == "agent" && authorID == agentID {
			continue // skip self-mention
		}

		agent, err := h.Queries.GetAgent(r.Context(), agentID)
		if err != nil || agent.RuntimeID == "" || agent.ArchivedAt.Valid {
			continue
		}
		if !agentHasTriggerEnabled(agent.Triggers, "on_mention") {
			continue
		}

		// Create or find an issue to run the agent task against.
		// For now, we skip issue-less agent execution since the task system
		// requires an issue_id. Users can create issues from chat to trigger agents.
		slog.Debug("channel mention detected but skipping task enqueue (no issue context)",
			"agent_id", agentID, "channel_id", channel.ID, "message_id", msg.ID)
	}
}

// parseAgentMentionsFromContent extracts agent IDs from mention://agent/<id> links.
func parseAgentMentionsFromContent(content string) []string {
	const prefix = "mention://agent/"
	var ids []string
	for i := 0; i < len(content); i++ {
		idx := indexOf(content[i:], prefix)
		if idx == -1 {
			break
		}
		start := i + idx + len(prefix)
		end := start
		for end < len(content) && content[end] != ')' && content[end] != ' ' && content[end] != '\n' {
			end++
		}
		if end > start {
			ids = append(ids, content[start:end])
		}
		i = end
	}
	return ids
}

func indexOf(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
