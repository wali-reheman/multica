"use client";

import { useState } from "react";
import { GitCommit, ChevronRight, FileText } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";
import type { CommitDetail } from "@/shared/types";
import { api } from "@/shared/api";
import { useProjectStore } from "../store";
import { DiffViewer } from "./diff-viewer";

export function CommitHistory({ projectId }: { projectId: string }) {
  const commits = useProjectStore((s) => s.commits);
  const loading = useProjectStore((s) => s.commitsLoading);
  const [expandedCommit, setExpandedCommit] = useState<string | null>(null);
  const [commitDetail, setCommitDetail] = useState<CommitDetail | null>(null);
  const [detailLoading, setDetailLoading] = useState(false);

  const handleExpand = async (sha: string) => {
    if (expandedCommit === sha) {
      setExpandedCommit(null);
      setCommitDetail(null);
      return;
    }
    setExpandedCommit(sha);
    setDetailLoading(true);
    try {
      const detail = await api.getProjectCommitDetail(projectId, sha);
      setCommitDetail(detail);
    } catch {
      setCommitDetail(null);
    } finally {
      setDetailLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex flex-col gap-3">
        {Array.from({ length: 5 }).map((_, i) => (
          <Skeleton key={i} className="h-14 w-full" />
        ))}
      </div>
    );
  }

  if (commits.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center gap-2 py-12 text-muted-foreground">
        <GitCommit className="size-8 opacity-40" />
        <p className="text-sm">No commits yet</p>
      </div>
    );
  }

  return (
    <div className="flex flex-col">
      {commits.map((commit) => (
        <div key={commit.hash}>
          <button
            type="button"
            className="flex w-full items-start gap-3 rounded-lg px-3 py-2.5 text-left hover:bg-accent/50 transition-colors"
            onClick={() => handleExpand(commit.hash)}
          >
            <div className="mt-1 flex flex-col items-center">
              <div className="flex size-6 items-center justify-center rounded-full bg-muted">
                <GitCommit className="size-3 text-muted-foreground" />
              </div>
            </div>
            <div className="min-w-0 flex-1">
              <p className="truncate text-sm font-medium">{commit.message}</p>
              <div className="mt-0.5 flex items-center gap-2 text-xs text-muted-foreground">
                <span className="font-mono">{commit.short_hash}</span>
                <span>&middot;</span>
                <span>{commit.author}</span>
                <span>&middot;</span>
                <span>{formatRelativeDate(commit.date)}</span>
                {commit.files_changed > 0 && (
                  <>
                    <span>&middot;</span>
                    <span className="flex items-center gap-0.5">
                      <FileText className="size-3" />
                      {commit.files_changed}
                    </span>
                  </>
                )}
              </div>
            </div>
            <ChevronRight
              className={`mt-1 size-4 text-muted-foreground transition-transform ${expandedCommit === commit.hash ? "rotate-90" : ""}`}
            />
          </button>

          {expandedCommit === commit.hash && (
            <div className="ml-9 border-l pl-4 pb-3">
              {detailLoading ? (
                <Skeleton className="h-32 w-full" />
              ) : commitDetail ? (
                <DiffViewer diffs={commitDetail.diffs} />
              ) : (
                <p className="text-xs text-muted-foreground">Failed to load diff</p>
              )}
            </div>
          )}
        </div>
      ))}
    </div>
  );
}

function formatRelativeDate(dateStr: string): string {
  const date = new Date(dateStr);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffMins = Math.floor(diffMs / 60_000);
  const diffHours = Math.floor(diffMs / 3_600_000);
  const diffDays = Math.floor(diffMs / 86_400_000);

  if (diffMins < 1) return "just now";
  if (diffMins < 60) return `${diffMins}m ago`;
  if (diffHours < 24) return `${diffHours}h ago`;
  if (diffDays < 7) return `${diffDays}d ago`;
  return date.toLocaleDateString();
}
