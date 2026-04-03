"use client";

import { useState, useCallback } from "react";
import {
  FilePlus2,
  FileEdit,
  FileX2,
  FileQuestion,
  Loader2,
} from "lucide-react";
import { Checkbox } from "@/components/ui/checkbox";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import type { FileStatusEntry } from "@/shared/types";
import { cn } from "@/lib/utils";
import { useProjectStore } from "../store";
import { CommitDialog } from "./commit-dialog";

type FileCategory = "modified" | "added" | "deleted" | "untracked";

const categoryConfig: Record<
  FileCategory,
  { label: string; icon: typeof FileEdit; color: string }
> = {
  modified: { label: "Modified", icon: FileEdit, color: "text-warning" },
  added: { label: "Added", icon: FilePlus2, color: "text-success" },
  deleted: { label: "Deleted", icon: FileX2, color: "text-destructive" },
  untracked: { label: "Untracked", icon: FileQuestion, color: "text-muted-foreground" },
};

function FileRow({
  entry,
  category,
  checked,
  onToggle,
}: {
  entry: FileStatusEntry;
  category: FileCategory;
  checked: boolean;
  onToggle: () => void;
}) {
  const config = categoryConfig[category];
  const Icon = config.icon;

  return (
    <label className="flex items-center gap-3 px-4 py-2 hover:bg-accent/50 transition-colors cursor-pointer">
      <Checkbox checked={checked} onCheckedChange={onToggle} />
      <Icon className={cn("h-4 w-4 shrink-0", config.color)} />
      <span className="text-sm truncate flex-1">{entry.path}</span>
      <span className="text-xs text-muted-foreground shrink-0">{config.label}</span>
    </label>
  );
}

export function FileStatus({ projectId }: { projectId: string }) {
  const status = useProjectStore((s) => s.status);
  const statusLoading = useProjectStore((s) => s.statusLoading);
  const fetchStatus = useProjectStore((s) => s.fetchStatus);
  const [selected, setSelected] = useState<Set<string>>(new Set());
  const [commitOpen, setCommitOpen] = useState(false);

  const allFiles: { entry: FileStatusEntry; category: FileCategory }[] = [];
  if (status) {
    for (const entry of status.modified) allFiles.push({ entry, category: "modified" });
    for (const entry of status.added) allFiles.push({ entry, category: "added" });
    for (const entry of status.deleted) allFiles.push({ entry, category: "deleted" });
    for (const entry of status.untracked) allFiles.push({ entry, category: "untracked" });
  }

  const toggleFile = useCallback((path: string) => {
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(path)) next.delete(path);
      else next.add(path);
      return next;
    });
  }, []);

  const toggleAll = useCallback(() => {
    if (selected.size === allFiles.length) {
      setSelected(new Set());
    } else {
      setSelected(new Set(allFiles.map((f) => f.entry.path)));
    }
  }, [selected.size, allFiles]);

  const handleCommitted = useCallback(() => {
    setSelected(new Set());
    setCommitOpen(false);
    fetchStatus(projectId);
  }, [fetchStatus, projectId]);

  if (statusLoading) {
    return (
      <div className="p-4 space-y-3">
        {Array.from({ length: 4 }).map((_, i) => (
          <div key={i} className="flex items-center gap-3">
            <Skeleton className="h-4 w-4 rounded" />
            <Skeleton className="h-4 w-4 rounded" />
            <Skeleton className="h-4 w-48" />
          </div>
        ))}
      </div>
    );
  }

  if (allFiles.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-16 text-muted-foreground">
        <FileEdit className="h-8 w-8 text-muted-foreground/30" />
        <p className="mt-3 text-sm">Working tree clean</p>
        <p className="mt-1 text-xs">No uncommitted changes</p>
      </div>
    );
  }

  return (
    <>
      <div className="flex flex-col h-full">
        {/* Toolbar */}
        <div className="flex items-center justify-between border-b px-4 py-2">
          <label className="flex items-center gap-2 cursor-pointer">
            <Checkbox
              checked={selected.size === allFiles.length}
              onCheckedChange={toggleAll}
            />
            <span className="text-xs text-muted-foreground">
              {selected.size} of {allFiles.length} selected
            </span>
          </label>
          <Button
            size="sm"
            disabled={selected.size === 0}
            onClick={() => setCommitOpen(true)}
          >
            Commit
          </Button>
        </div>

        {/* File list */}
        <div className="flex-1 overflow-y-auto divide-y">
          {allFiles.map(({ entry, category }) => (
            <FileRow
              key={entry.path}
              entry={entry}
              category={category}
              checked={selected.has(entry.path)}
              onToggle={() => toggleFile(entry.path)}
            />
          ))}
        </div>
      </div>

      <CommitDialog
        open={commitOpen}
        onOpenChange={setCommitOpen}
        projectId={projectId}
        files={Array.from(selected)}
        onCommitted={handleCommitted}
      />
    </>
  );
}
