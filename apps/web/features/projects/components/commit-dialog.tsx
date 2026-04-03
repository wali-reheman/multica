"use client";

import { useState } from "react";
import { toast } from "sonner";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { api } from "@/shared/api";
import { useProjectStore } from "../store";

interface CommitDialogProps {
  projectId: string;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function CommitDialog({ projectId, open, onOpenChange }: CommitDialogProps) {
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);
  const fetchCommits = useProjectStore((s) => s.fetchCommits);
  const fetchStatus = useProjectStore((s) => s.fetchStatus);

  const handleCommit = async () => {
    if (!message.trim()) return;
    setLoading(true);
    try {
      await api.createProjectCommit(projectId, { message: message.trim() });
      toast.success("Commit created");
      setMessage("");
      onOpenChange(false);
      // Refresh
      await Promise.all([fetchCommits(projectId), fetchStatus(projectId)]);
    } catch {
      toast.error("Failed to create commit");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Create Commit</DialogTitle>
        </DialogHeader>
        <div className="flex flex-col gap-3">
          <Input
            placeholder="Commit message summary"
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter" && !e.shiftKey) {
                e.preventDefault();
                handleCommit();
              }
            }}
            autoFocus
          />
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button onClick={handleCommit} disabled={!message.trim() || loading}>
            {loading ? "Committing..." : "Commit"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
