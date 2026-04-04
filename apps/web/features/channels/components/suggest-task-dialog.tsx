"use client";

import { useState } from "react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { api } from "@/shared/api";
import { useWorkspaceStore } from "@/features/workspace";
import { useChannelStore } from "../store";

interface SuggestTaskDialogProps {
  channelId: string;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function SuggestTaskDialog({
  channelId,
  open,
  onOpenChange,
}: SuggestTaskDialogProps) {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [priority, setPriority] = useState("none");
  const [assigneeId, setAssigneeId] = useState<string>("");
  const [assigneeType, setAssigneeType] = useState<string>("");
  const [submitting, setSubmitting] = useState(false);

  const agents = useWorkspaceStore((s) => s.agents).filter((a) => !a.archived_at);
  const members = useWorkspaceStore((s) => s.members);
  const addSuggestion = useChannelStore((s) => s.addSuggestion);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;

    setSubmitting(true);
    try {
      const suggestion = await api.createSuggestion(channelId, {
        title: title.trim(),
        description: description || undefined,
        priority,
        assignee_type: assigneeType || undefined,
        assignee_id: assigneeId || undefined,
      });
      addSuggestion(channelId, suggestion);
      onOpenChange(false);
      setTitle("");
      setDescription("");
      setPriority("none");
      setAssigneeId("");
      setAssigneeType("");
      toast.success("Task suggested — awaiting approval");
    } catch {
      toast.error("Failed to suggest task");
    } finally {
      setSubmitting(false);
    }
  };

  const handleAssigneeChange = (v: string | null) => {
    if (!v) {
      setAssigneeId("");
      setAssigneeType("");
      return;
    }
    const parts = v.split(":");
    setAssigneeType(parts[0] ?? "");
    setAssigneeId(parts[1] ?? "");
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>Suggest a Task</DialogTitle>
          </DialogHeader>
          <div className="space-y-4 py-4">
            <p className="text-sm text-muted-foreground">
              Propose a task for the team. It will appear in the channel for approval before
              an issue is created and work begins.
            </p>
            <div className="space-y-2">
              <Label htmlFor="suggest-title">Title</Label>
              <Input
                id="suggest-title"
                placeholder="What needs to be done?"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                autoFocus
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="suggest-desc">Description</Label>
              <Textarea
                id="suggest-desc"
                placeholder="Optional details"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                rows={3}
              />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>Priority</Label>
                <Select value={priority} onValueChange={(v) => { if (v) setPriority(v); }}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="none">None</SelectItem>
                    <SelectItem value="low">Low</SelectItem>
                    <SelectItem value="medium">Medium</SelectItem>
                    <SelectItem value="high">High</SelectItem>
                    <SelectItem value="urgent">Urgent</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label>Assign to</Label>
                <Select
                  value={assigneeId ? `${assigneeType}:${assigneeId}` : ""}
                  onValueChange={handleAssigneeChange}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Unassigned" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="">Unassigned</SelectItem>
                    {agents.map((agent) => (
                      <SelectItem key={agent.id} value={`agent:${agent.id}`}>
                        {agent.name} (Agent)
                      </SelectItem>
                    ))}
                    {members.map((member) => (
                      <SelectItem key={member.user_id} value={`member:${member.user_id}`}>
                        {member.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>
          </div>
          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={!title.trim() || submitting}>
              {submitting ? "Suggesting..." : "Suggest Task"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
