"use client";

<<<<<<< HEAD
import { useState } from "react";
import { FileText, FilePlus, FileX, FilePenLine, FileQuestion, GitCommit } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Skeleton } from "@/components/ui/skeleton";
import type { FileStatusEntry } from "@/shared/types";
import { useProjectStore } from "../store";

interface FileStatusProps {
  projectId: string;
  onCommit: () => void;
}

export function FileStatus({ projectId, onCommit }: FileStatusProps) {
  const status = useProjectStore((s) => s.status);
  const loading = useProjectStore((s) => s.statusLoading);
  const [selectedFiles, setSelectedFiles] = useState<Set<string>>(new Set());

  if (loading) {
    return (
      <div className="flex flex-col gap-3">
        {Array.from({ length: 3 }).map((_, i) => (
          <Skeleton key={i} className="h-10 w-full" />
=======
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
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
        ))}
      </div>
    );
  }

<<<<<<< HEAD
  if (!status) {
    return (
      <div className="flex flex-col items-center justify-center gap-2 py-12 text-muted-foreground">
        <FileText className="size-8 opacity-40" />
        <p className="text-sm">Unable to load file status</p>
      </div>
    );
  }

  const allFiles = [
    ...status.modified.map((f) => ({ ...f, category: "modified" as const })),
    ...status.added.map((f) => ({ ...f, category: "added" as const })),
    ...status.deleted.map((f) => ({ ...f, category: "deleted" as const })),
    ...status.untracked.map((f) => ({ ...f, category: "untracked" as const })),
  ];

  if (allFiles.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center gap-2 py-12 text-muted-foreground">
        <GitCommit className="size-8 opacity-40" />
        <p className="text-sm">Working tree clean</p>
        <p className="text-xs">No changes to commit</p>
=======
  if (allFiles.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-16 text-muted-foreground">
        <FileEdit className="h-8 w-8 text-muted-foreground/30" />
        <p className="mt-3 text-sm">Working tree clean</p>
        <p className="mt-1 text-xs">No uncommitted changes</p>
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
      </div>
    );
  }

<<<<<<< HEAD
  const toggleFile = (path: string) => {
    setSelectedFiles((prev) => {
      const next = new Set(prev);
      if (next.has(path)) {
        next.delete(path);
      } else {
        next.add(path);
      }
      return next;
    });
  };

  const toggleAll = () => {
    if (selectedFiles.size === allFiles.length) {
      setSelectedFiles(new Set());
    } else {
      setSelectedFiles(new Set(allFiles.map((f) => f.path)));
    }
  };

  return (
    <div className="flex flex-col gap-4">
      {/* Summary */}
      <div className="flex items-center gap-4 text-xs text-muted-foreground">
        {status.modified.length > 0 && (
          <span className="flex items-center gap-1">
            <FilePenLine className="size-3 text-yellow-600" />
            {status.modified.length} modified
          </span>
        )}
        {status.added.length > 0 && (
          <span className="flex items-center gap-1">
            <FilePlus className="size-3 text-green-600" />
            {status.added.length} added
          </span>
        )}
        {status.deleted.length > 0 && (
          <span className="flex items-center gap-1">
            <FileX className="size-3 text-red-600" />
            {status.deleted.length} deleted
          </span>
        )}
        {status.untracked.length > 0 && (
          <span className="flex items-center gap-1">
            <FileQuestion className="size-3 text-muted-foreground" />
            {status.untracked.length} untracked
          </span>
        )}
      </div>

      {/* File list with checkboxes */}
      <div className="rounded border">
        <div className="flex items-center gap-2 border-b bg-muted/50 px-3 py-1.5">
          <Checkbox
            checked={selectedFiles.size === allFiles.length}
            onCheckedChange={toggleAll}
          />
          <span className="text-xs text-muted-foreground">
            {selectedFiles.size > 0
              ? `${selectedFiles.size} of ${allFiles.length} selected`
              : `${allFiles.length} changed file${allFiles.length !== 1 ? "s" : ""}`}
          </span>
          {selectedFiles.size > 0 && (
            <Button
              size="sm"
              variant="outline"
              className="ml-auto h-6 text-xs"
              onClick={onCommit}
            >
              <GitCommit className="mr-1 size-3" />
              Commit selected
            </Button>
          )}
        </div>
        {allFiles.map((file) => (
          <FileRow
            key={file.path}
            file={file}
            category={file.category}
            selected={selectedFiles.has(file.path)}
            onToggle={() => toggleFile(file.path)}
          />
        ))}
      </div>
    </div>
  );
}

function FileRow({
  file,
  category,
  selected,
  onToggle,
}: {
  file: FileStatusEntry;
  category: "modified" | "added" | "deleted" | "untracked";
  selected: boolean;
  onToggle: () => void;
}) {
  const icons = {
    modified: <FilePenLine className="size-3.5 text-yellow-600 dark:text-yellow-400" />,
    added: <FilePlus className="size-3.5 text-green-600 dark:text-green-400" />,
    deleted: <FileX className="size-3.5 text-red-600 dark:text-red-400" />,
    untracked: <FileQuestion className="size-3.5 text-muted-foreground" />,
  };

  return (
    <div className="flex items-center gap-2 border-b last:border-0 px-3 py-1.5 hover:bg-accent/30">
      <Checkbox checked={selected} onCheckedChange={onToggle} />
      {icons[category]}
      <span className="flex-1 truncate text-xs font-mono">{file.path}</span>
      <span className="text-[10px] uppercase text-muted-foreground">{category}</span>
    </div>
=======
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
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
  );
}
