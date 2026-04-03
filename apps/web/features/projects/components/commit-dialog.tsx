"use client";

import { useState, useCallback } from "react";
import { Loader2 } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { api } from "@/shared/api";
import { toast } from "sonner";

export function CommitDialog({
  open,
  onOpenChange,
  projectId,
  files,
  onCommitted,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  projectId: string;
  files: string[];
  onCommitted: () => void;
}) {
  const [message, setMessage] = useState("");
  const [description, setDescription] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = useCallback(
    async (e: React.FormEvent) => {
      e.preventDefault();
      const fullMessage = description
        ? `${message}\n\n${description}`
        : message;
      if (!fullMessage.trim()) return;

      setSubmitting(true);
      try {
        await api.createProjectCommit(projectId, {
          message: fullMessage,
          files: files.length > 0 ? files : undefined,
        });
        toast.success("Commit created");
        setMessage("");
        setDescription("");
        onCommitted();
      } catch {
        toast.error("Failed to create commit");
      } finally {
        setSubmitting(false);
      }
    },
    [projectId, files, message, description, onCommitted],
  );

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create Commit</DialogTitle>
          <DialogDescription>
            {files.length} file{files.length !== 1 ? "s" : ""} selected
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="commit-message">Commit message</Label>
            <Input
              id="commit-message"
              placeholder="Short summary of changes"
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              autoFocus
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="commit-description">Description (optional)</Label>
            <Textarea
              id="commit-description"
              placeholder="Detailed description..."
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={3}
            />
          </div>
          <DialogFooter>
            <Button
              type="submit"
              disabled={!message.trim() || submitting}
            >
              {submitting && <Loader2 className="h-4 w-4 animate-spin" />}
              Commit
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
