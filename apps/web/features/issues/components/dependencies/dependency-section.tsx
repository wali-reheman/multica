"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { ArrowRight, ArrowLeft, Link2, Plus, X } from "lucide-react";
import { toast } from "sonner";
import type { IssueDependency, DependencyType, Issue } from "@/shared/types";
import { api } from "@/shared/api";
import { useIssueStore } from "@/features/issues/store";
import { StatusIcon } from "@/features/issues/components/status-icon";
import {
  Popover,
  PopoverTrigger,
  PopoverContent,
} from "@/components/ui/popover";
import { Button } from "@/components/ui/button";

const DEP_CONFIG: Record<DependencyType, { label: string; icon: typeof ArrowRight }> = {
  blocks: { label: "Blocks", icon: ArrowRight },
  blocked_by: { label: "Blocked by", icon: ArrowLeft },
  related: { label: "Related to", icon: Link2 },
};

interface DependencySectionProps {
  issueId: string;
}

export function DependencySection({ issueId }: DependencySectionProps) {
  const router = useRouter();
  const [deps, setDeps] = useState<IssueDependency[]>([]);
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState("");
  const [depType, setDepType] = useState<DependencyType>("blocks");
  const allIssues = useIssueStore((s) => s.issues);

  useEffect(() => {
    api.listIssueDependencies(issueId).then(setDeps).catch(() => {});
  }, [issueId]);

  const relatedIssueIds = useMemo(
    () => new Set([issueId, ...deps.map((d) => d.issue_id), ...deps.map((d) => d.depends_on_issue_id)]),
    [issueId, deps],
  );

  const filteredIssues = useMemo(() => {
    const q = search.toLowerCase();
    return allIssues
      .filter((i) => !relatedIssueIds.has(i.id))
      .filter((i) => !q || i.title.toLowerCase().includes(q) || i.identifier.toLowerCase().includes(q))
      .slice(0, 8);
  }, [allIssues, relatedIssueIds, search]);

  const handleAdd = useCallback(
    async (targetIssueId: string) => {
      setOpen(false);
      setSearch("");
      try {
        const dep = await api.createIssueDependency(issueId, targetIssueId, depType);
        setDeps((prev) => [...prev, dep]);
      } catch {
        toast.error("Failed to add dependency");
      }
    },
    [issueId, depType],
  );

  const handleRemove = useCallback(
    async (depId: string) => {
      setDeps((prev) => prev.filter((d) => d.id !== depId));
      try {
        await api.deleteIssueDependency(issueId, depId);
      } catch {
        toast.error("Failed to remove dependency");
        api.listIssueDependencies(issueId).then(setDeps).catch(() => {});
      }
    },
    [issueId],
  );

  const getIssue = useCallback(
    (id: string): Issue | undefined => allIssues.find((i) => i.id === id),
    [allIssues],
  );

  if (deps.length === 0) {
    return (
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger className="flex items-center gap-1.5 cursor-pointer rounded px-1 -mx-1 hover:bg-accent/30 transition-colors text-sm text-muted-foreground">
          <Link2 className="size-3.5" />
          Add dependency
        </PopoverTrigger>
        <PopoverContent align="start" className="w-64 p-0">
          <DependencyAdder
            search={search}
            onSearchChange={setSearch}
            depType={depType}
            onDepTypeChange={setDepType}
            issues={filteredIssues}
            onSelect={handleAdd}
          />
        </PopoverContent>
      </Popover>
    );
  }

  return (
    <div className="space-y-1">
      {deps.map((dep) => {
        const isSource = dep.issue_id === issueId;
        const otherIssueId = isSource ? dep.depends_on_issue_id : dep.issue_id;
        const otherIssue = getIssue(otherIssueId);
        const config = DEP_CONFIG[dep.type];
        const Icon = config.icon;

        return (
          <div
            key={dep.id}
            className="group flex items-center gap-1.5 rounded px-1 -mx-1 hover:bg-accent/30 transition-colors"
          >
            <Icon className="size-3 text-muted-foreground shrink-0" />
            <button
              type="button"
              className="flex items-center gap-1 min-w-0 text-xs truncate hover:underline"
              onClick={() => router.push(`/issues/${otherIssueId}`)}
            >
              {otherIssue && <StatusIcon status={otherIssue.status} className="size-3 shrink-0" />}
              <span className="text-muted-foreground shrink-0">{otherIssue?.identifier}</span>
              <span className="truncate">{otherIssue?.title ?? otherIssueId}</span>
            </button>
            <button
              type="button"
              onClick={() => handleRemove(dep.id)}
              className="ml-auto opacity-0 group-hover:opacity-100 p-0.5 rounded hover:bg-accent transition-opacity"
            >
              <X className="size-3 text-muted-foreground" />
            </button>
          </div>
        );
      })}
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger className="flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground transition-colors mt-1">
          <Plus className="size-3" />
          Add
        </PopoverTrigger>
        <PopoverContent align="start" className="w-64 p-0">
          <DependencyAdder
            search={search}
            onSearchChange={setSearch}
            depType={depType}
            onDepTypeChange={setDepType}
            issues={filteredIssues}
            onSelect={handleAdd}
          />
        </PopoverContent>
      </Popover>
    </div>
  );
}

function DependencyAdder({
  search,
  onSearchChange,
  depType,
  onDepTypeChange,
  issues,
  onSelect,
}: {
  search: string;
  onSearchChange: (v: string) => void;
  depType: DependencyType;
  onDepTypeChange: (v: DependencyType) => void;
  issues: Issue[];
  onSelect: (issueId: string) => void;
}) {
  return (
    <div>
      <div className="flex gap-1 p-2 border-b">
        {(["blocks", "blocked_by", "related"] as const).map((t) => (
          <button
            key={t}
            type="button"
            onClick={() => onDepTypeChange(t)}
            className={`rounded px-2 py-1 text-xs transition-colors ${
              depType === t ? "bg-accent font-medium" : "text-muted-foreground hover:bg-accent/50"
            }`}
          >
            {DEP_CONFIG[t].label}
          </button>
        ))}
      </div>
      <div className="p-2 border-b">
        <input
          type="text"
          value={search}
          onChange={(e) => onSearchChange(e.target.value)}
          placeholder="Search issues..."
          className="w-full bg-transparent text-sm placeholder:text-muted-foreground outline-none"
          autoFocus
        />
      </div>
      <div className="max-h-48 overflow-y-auto p-1">
        {issues.length === 0 && (
          <p className="px-2 py-3 text-center text-sm text-muted-foreground">
            No matching issues
          </p>
        )}
        {issues.map((issue) => (
          <button
            key={issue.id}
            type="button"
            onClick={() => onSelect(issue.id)}
            className="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-sm hover:bg-accent transition-colors"
          >
            <StatusIcon status={issue.status} className="size-3.5 shrink-0" />
            <span className="text-muted-foreground text-xs shrink-0">{issue.identifier}</span>
            <span className="truncate">{issue.title}</span>
          </button>
        ))}
      </div>
    </div>
  );
}
