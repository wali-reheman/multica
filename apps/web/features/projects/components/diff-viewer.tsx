"use client";

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
    </div>
  );
}
