export interface Project {
  id: string;
  workspace_id: string;
  name: string;
  local_path: string;
  default_branch: string;
  language: string | null;
  file_count: number;
  size_bytes: number;
  last_opened_at: string | null;
  created_at: string;
  updated_at: string;
  is_git_repo: boolean;
}

export interface CommitInfo {
  hash: string;
  short_hash: string;
  message: string;
  author: string;
  author_email: string;
  date: string;
  files_changed: number;
}

export interface CommitDetail extends CommitInfo {
  diffs: DiffEntry[];
}

export interface DiffEntry {
  path: string;
  old_path?: string;
  change: "add" | "modify" | "delete" | "rename";
  patch?: string;
}

export interface BranchInfo {
  name: string;
  is_head: boolean;
  is_remote: boolean;
  hash: string;
}

export interface FileStatusEntry {
  path: string;
  staging: string;
  worktree: string;
}

export interface GitStatus {
  modified: FileStatusEntry[];
  added: FileStatusEntry[];
  deleted: FileStatusEntry[];
  untracked: FileStatusEntry[];
}

export interface CreateProjectRequest {
  name?: string;
  local_path: string;
  init_git?: boolean;
}

export interface UpdateProjectRequest {
  name?: string;
  default_branch?: string;
}

export interface CreateCommitRequest {
  message: string;
  files?: string[];
}

export interface CreateBranchRequest {
  name: string;
}

export interface CheckoutBranchRequest {
  branch: string;
}
