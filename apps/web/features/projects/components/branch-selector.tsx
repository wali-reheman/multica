"use client";

import { GitBranch, Check, Plus } from "lucide-react";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";
import { toast } from "sonner";
import { api } from "@/shared/api";
import { useProjectStore } from "../store";

export function BranchSelector({ projectId }: { projectId: string }) {
  const branches = useProjectStore((s) => s.branches);
  const fetchBranches = useProjectStore((s) => s.fetchBranches);
  const fetchCommits = useProjectStore((s) => s.fetchCommits);
  const fetchStatus = useProjectStore((s) => s.fetchStatus);
  const [creating, setCreating] = useState(false);
  const [newBranchName, setNewBranchName] = useState("");

  const currentBranch = branches.find((b) => b.is_head && !b.is_remote);
  const localBranches = branches.filter((b) => !b.is_remote);

  const handleCheckout = async (branch: string) => {
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
    }
  };

  const handleCreate = async () => {
    if (!newBranchName.trim()) return;
    try {
      await api.createProjectBranch(projectId, { name: newBranchName.trim() });
      setNewBranchName("");
      setCreating(false);
      await fetchBranches(projectId);
      toast.success(`Branch ${newBranchName.trim()} created`);
    } catch {
      toast.error("Failed to create branch");
    }
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger
        render={
          <Button variant="outline" size="sm">
            <GitBranch className="mr-1.5 size-3.5" />
            {currentBranch?.name ?? "no branch"}
          </Button>
        }
      />
      <DropdownMenuContent align="end" className="w-56">
        <DropdownMenuLabel className="text-xs text-muted-foreground">
          Branches
        </DropdownMenuLabel>
        <DropdownMenuGroup>
          {localBranches.map((branch) => (
            <DropdownMenuItem
              key={branch.name}
              onClick={() => {
                if (!branch.is_head) handleCheckout(branch.name);
              }}
            >
              <GitBranch className="size-3.5 text-muted-foreground" />
              <span className="flex-1 truncate">{branch.name}</span>
              {branch.is_head && <Check className="size-3.5 text-primary" />}
            </DropdownMenuItem>
          ))}
        </DropdownMenuGroup>
        <DropdownMenuSeparator />
        {creating ? (
          <div className="flex items-center gap-1 p-1">
            <Input
              value={newBranchName}
              onChange={(e) => setNewBranchName(e.target.value)}
              placeholder="branch-name"
              className="h-7 text-xs"
              autoFocus
              onKeyDown={(e) => {
                if (e.key === "Enter") handleCreate();
                if (e.key === "Escape") setCreating(false);
              }}
            />
            <Button size="sm" variant="ghost" className="h-7 px-2" onClick={handleCreate}>
              <Check className="size-3.5" />
            </Button>
          </div>
        ) : (
          <DropdownMenuItem onClick={() => setCreating(true)}>
            <Plus className="size-3.5" />
            New branch
          </DropdownMenuItem>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
