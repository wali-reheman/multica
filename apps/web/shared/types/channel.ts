export interface Channel {
  id: string;
  workspace_id: string;
  name: string;
  description: string;
  type: "group" | "direct";
  created_by_type: "member" | "agent";
  created_by_id: string;
  archived_at: string | null;
  created_at: string;
  updated_at: string;
  members?: ChannelMember[];
}

export interface ChannelMember {
  channel_id: string;
  member_type: "member" | "agent";
  member_id: string;
  role: "owner" | "member";
  joined_at: string;
}

export type ChannelMessageType = "message" | "system" | "issue_created";

export interface ChannelMessage {
  id: string;
  channel_id: string;
  author_type: "member" | "agent";
  author_id: string;
  content: string;
  type: ChannelMessageType;
  parent_id: string | null;
  issue_id: string | null;
  created_at: string;
  updated_at: string;
}

export interface CreateChannelRequest {
  name: string;
  description?: string;
  type?: "group" | "direct";
  member_ids?: string[];
  agent_ids?: string[];
}

export interface CreateChannelMessageRequest {
  content: string;
  parent_id?: string;
}

export interface CreateIssueFromChannelRequest {
  title: string;
  description?: string;
  priority?: string;
  assignee_id?: string;
}
