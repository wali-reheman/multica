"use client";

import { useState, useCallback } from "react";
import { GitBranch, Check, Plus, Loader2 } from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuLabel,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import type { BranchInfo } from "@/shared/types";
import { api } from "@/shared/api";
import { toast } from "sonner";
import { cn } from "@/lib/utils";
import { useProjectStore } from "../store";

export function BranchSelector({
  projectId,
  branches,
  currentBranch,
  loading,
}: {
  projectId: string;
  branches: BranchInfo[];
  currentBranch: string;
  loading: boolean;
}) {
  const fetchBranches = useProjectStore((s) => s.fetchBranches);
  const fetchCommits = useProjectStore((s) => s.fetchCommits);
  const fetchStatus = useProjectStore((s) => s.fetchStatus);
  const [switching, setSwitching] = useState(false);

  const handleCheckout = useCallback(
    async (branch: string) => {
      if (branch === currentBranch) return;
      setSwitching(true);
      try {
        await api.checkoutProjectBranch(projectId, { branch });
        await Promise.all([
          fetchBranches(projectId),
          fetchCommits(projectId),
          fetchStatus(projectId),
        ]);
        toast.success(`Switched to ${branch}`);
      } catch {
        toast.error("Failed to switch branch");
      } finally {
        setSwitching(false);
      }
    },
    [projectId, currentBranch, fetchBranches, fetchCommits, fetchStatus],
  );

  const localBranches = branches.filter((b) => !b.is_remote);

  return (
    <DropdownMenu>
      <DropdownMenuTrigger
        render={
          <Button variant="outline" size="sm" disabled={loading || switching}>
            {switching ? (
              <Loader2 className="h-3 w-3 animate-spin" />
            ) : (
              <GitBranch className="h-3 w-3" />
            )}
            <span className="max-w-[120px] truncate">{currentBranch}</span>
          </Button>
        }
      />
      <DropdownMenuContent align="end">
        <DropdownMenuLabel>Branches</DropdownMenuLabel>
        {localBranches.map((branch) => (
          <DropdownMenuItem
            key={branch.name}
            onClick={() => handleCheckout(branch.name)}
          >
            <Check
              className={cn(
                "h-3.5 w-3.5",
                branch.is_head ? "opacity-100" : "opacity-0",
              )}
            />
            <span className="truncate">{branch.name}</span>
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
