"use client";

import { FilePlus2, FileEdit, FileX2, FileSymlink } from "lucide-react";
import type { DiffEntry } from "@/shared/types";
import { cn } from "@/lib/utils";

const changeIcons: Record<DiffEntry["change"], typeof FileEdit> = {
  add: FilePlus2,
  modify: FileEdit,
  delete: FileX2,
  rename: FileSymlink,
};

const changeColors: Record<DiffEntry["change"], string> = {
  add: "text-success",
  modify: "text-warning",
  delete: "text-destructive",
  rename: "text-info",
};

const changeBgColors: Record<DiffEntry["change"], string> = {
  add: "bg-success/10",
  modify: "bg-warning/10",
  delete: "bg-destructive/10",
  rename: "bg-info/10",
};

function PatchLine({ line }: { line: string }) {
  let lineClass = "text-muted-foreground";
  if (line.startsWith("+") && !line.startsWith("+++")) {
    lineClass = "text-success bg-success/5";
  } else if (line.startsWith("-") && !line.startsWith("---")) {
    lineClass = "text-destructive bg-destructive/5";
  } else if (line.startsWith("@@")) {
    lineClass = "text-info bg-info/5";
  }

  return (
    <div className={cn("px-3 font-mono text-[11px] leading-5 whitespace-pre", lineClass)}>
      {line}
    </div>
  );
}

export function DiffViewer({ diff }: { diff: DiffEntry }) {
  const Icon = changeIcons[diff.change];
  const color = changeColors[diff.change];
  const bgColor = changeBgColors[diff.change];

  const lines = diff.patch?.split("\n") ?? [];

  return (
    <div className="rounded-lg border overflow-hidden">
      <div className={cn("flex items-center gap-2 px-3 py-2", bgColor)}>
        <Icon className={cn("h-3.5 w-3.5", color)} />
        <span className="text-xs font-medium truncate">{diff.path}</span>
        {diff.old_path && diff.old_path !== diff.path && (
          <span className="text-xs text-muted-foreground truncate">
            (from {diff.old_path})
          </span>
        )}
      </div>
      {lines.length > 0 && (
        <div className="max-h-64 overflow-y-auto border-t bg-background">
          {lines.map((line, i) => (
            <PatchLine key={i} line={line} />
          ))}
        </div>
      )}
    </div>
  );
}
