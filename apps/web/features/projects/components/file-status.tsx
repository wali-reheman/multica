"use client";

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
        ))}
      </div>
    );
  }

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
      </div>
    );
  }

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
  );
}
