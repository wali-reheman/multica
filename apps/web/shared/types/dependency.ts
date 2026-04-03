export type DependencyType = "blocks" | "blocked_by" | "related";

export interface IssueDependency {
  id: string;
  issue_id: string;
  depends_on_issue_id: string;
  type: DependencyType;
}
