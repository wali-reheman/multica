"use client";

<<<<<<< HEAD
import { FilePlus, FileX, FilePenLine, ArrowRightLeft } from "lucide-react";
import type { DiffEntry } from "@/shared/types";

const changeIcons = {
  add: FilePlus,
  delete: FileX,
  modify: FilePenLine,
  rename: ArrowRightLeft,
} as const;

const changeColors = {
  add: "text-green-600 dark:text-green-400",
  delete: "text-red-600 dark:text-red-400",
  modify: "text-yellow-600 dark:text-yellow-400",
  rename: "text-blue-600 dark:text-blue-400",
} as const;

export function DiffViewer({ diffs }: { diffs: DiffEntry[] }) {
  if (diffs.length === 0) {
    return <p className="text-xs text-muted-foreground">No file changes</p>;
  }

  return (
    <div className="flex flex-col gap-2">
      {diffs.map((diff, idx) => {
        const Icon = changeIcons[diff.change] ?? FilePenLine;
        const color = changeColors[diff.change] ?? "text-muted-foreground";

        return (
          <div key={`${diff.path}-${idx}`} className="rounded border">
            <div className="flex items-center gap-2 bg-muted/50 px-3 py-1.5">
              <Icon className={`size-3.5 ${color}`} />
              <span className="text-xs font-mono">
                {diff.old_path ? (
                  <>
                    <span className="text-muted-foreground line-through">{diff.old_path}</span>
                    {" → "}
                  </>
                ) : null}
                {diff.path}
              </span>
              <span className={`ml-auto text-[10px] uppercase ${color}`}>
                {diff.change}
              </span>
            </div>
            {diff.patch && (
              <pre className="max-h-64 overflow-auto bg-background p-3 text-xs font-mono leading-relaxed">
                {diff.patch.split("\n").map((line, i) => (
                  <div
                    key={i}
                    className={
                      line.startsWith("+")
                        ? "bg-green-500/10 text-green-700 dark:text-green-400"
                        : line.startsWith("-")
                          ? "bg-red-500/10 text-red-700 dark:text-red-400"
                          : line.startsWith("@@")
                            ? "text-blue-600 dark:text-blue-400"
                            : "text-foreground"
                    }
                  >
                    {line}
                  </div>
                ))}
              </pre>
            )}
          </div>
        );
      })}
=======
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
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
    </div>
  );
}
