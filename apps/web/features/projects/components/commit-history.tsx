"use client";

<<<<<<< HEAD
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
=======
import { useState, useCallback } from "react";
import { GitCommitHorizontal, ChevronDown, ChevronRight, Loader2 } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";
import type { CommitInfo, CommitDetail } from "@/shared/types";
import { api } from "@/shared/api";
import { cn } from "@/lib/utils";
import { useProjectStore } from "../store";
import { DiffViewer } from "./diff-viewer";

function CommitItem({
  commit,
  projectId,
}: {
  commit: CommitInfo;
  projectId: string;
}) {
  const [expanded, setExpanded] = useState(false);
  const [detail, setDetail] = useState<CommitDetail | null>(null);
  const [loading, setLoading] = useState(false);

  const toggle = useCallback(async () => {
    if (expanded) {
      setExpanded(false);
      return;
    }
    if (!detail) {
      setLoading(true);
      try {
        const data = await api.getProjectCommitDetail(projectId, commit.hash);
        setDetail(data);
      } catch {
        // silently fail - user can retry
      } finally {
        setLoading(false);
      }
    }
    setExpanded(true);
  }, [expanded, detail, projectId, commit.hash]);

  const Chevron = expanded ? ChevronDown : ChevronRight;

  return (
    <div className="border-b last:border-b-0">
      <button
        onClick={toggle}
        className="flex w-full items-start gap-3 px-4 py-3 text-left hover:bg-accent/50 transition-colors"
      >
        <div className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-muted mt-0.5">
          <GitCommitHorizontal className="h-3.5 w-3.5 text-muted-foreground" />
        </div>
        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-2">
            <span className="text-sm font-medium truncate">{commit.message}</span>
            {loading && <Loader2 className="h-3 w-3 animate-spin text-muted-foreground" />}
          </div>
          <div className="mt-1 flex items-center gap-2 text-xs text-muted-foreground">
            <code className="rounded bg-muted px-1 py-0.5 font-mono text-[11px]">
              {commit.short_hash}
            </code>
            <span>{commit.author}</span>
            <span>&middot;</span>
            <span>{new Date(commit.date).toLocaleDateString()}</span>
            {commit.files_changed > 0 && (
              <>
                <span>&middot;</span>
                <span>{commit.files_changed} file{commit.files_changed !== 1 ? "s" : ""}</span>
              </>
            )}
          </div>
        </div>
        <Chevron className="h-4 w-4 shrink-0 text-muted-foreground mt-1" />
      </button>

      {expanded && detail && (
        <div className="border-t bg-muted/20 px-4 py-3">
          {detail.diffs.length === 0 ? (
            <p className="text-xs text-muted-foreground">No file changes</p>
          ) : (
            <div className="space-y-2">
              {detail.diffs.map((diff) => (
                <DiffViewer key={diff.path} diff={diff} />
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  );
}

export function CommitHistory({ projectId }: { projectId: string }) {
  const commits = useProjectStore((s) => s.commits);
  const commitsLoading = useProjectStore((s) => s.commitsLoading);

  if (commitsLoading) {
    return (
      <div className="p-4 space-y-3">
        {Array.from({ length: 5 }).map((_, i) => (
          <div key={i} className="flex items-start gap-3">
            <Skeleton className="h-6 w-6 rounded-full" />
            <div className="flex-1 space-y-1.5">
              <Skeleton className="h-4 w-48" />
              <Skeleton className="h-3 w-32" />
            </div>
          </div>
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
        ))}
      </div>
    );
  }

  if (commits.length === 0) {
    return (
<<<<<<< HEAD
      <div className="flex flex-col items-center justify-center gap-2 py-12 text-muted-foreground">
        <GitCommit className="size-8 opacity-40" />
        <p className="text-sm">No commits yet</p>
=======
      <div className="flex flex-col items-center justify-center py-16 text-muted-foreground">
        <GitCommitHorizontal className="h-8 w-8 text-muted-foreground/30" />
        <p className="mt-3 text-sm">No commits yet</p>
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
      </div>
    );
  }

  return (
<<<<<<< HEAD
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
=======
    <div className={cn("divide-y-0")}>
      {commits.map((commit) => (
        <CommitItem key={commit.hash} commit={commit} projectId={projectId} />
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
      ))}
    </div>
  );
}
<<<<<<< HEAD

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
=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
