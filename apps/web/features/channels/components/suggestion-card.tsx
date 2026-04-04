"use client";

import { useState } from "react";
import { Check, X, Lightbulb } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { api } from "@/shared/api";
import { useActorName } from "@/features/workspace";
import { useIssueStore } from "@/features/issues";
import { useChannelStore } from "../store";
import type { TaskSuggestion } from "@/shared/types";

interface SuggestionCardProps {
  suggestion: TaskSuggestion;
}

export function SuggestionCard({ suggestion }: SuggestionCardProps) {
  const { getActorName } = useActorName();
  const [loading, setLoading] = useState(false);

  const isPending = suggestion.status === "pending";
  const suggestedBy = getActorName(suggestion.suggested_by_type, suggestion.suggested_by_id);
  const assigneeName = suggestion.assignee_id
    ? getActorName(suggestion.assignee_type ?? "member", suggestion.assignee_id)
    : null;

  const handleApprove = async () => {
    setLoading(true);
    try {
      const result = await api.approveSuggestion(suggestion.channel_id, suggestion.id);
      useChannelStore.getState().updateSuggestion(suggestion.channel_id, suggestion.id, result.suggestion);
      useIssueStore.getState().addIssue(result.issue);
      toast.success(`Issue ${result.issue.identifier} created`);
    } catch {
      toast.error("Failed to approve suggestion");
    } finally {
      setLoading(false);
    }
  };

  const handleDismiss = async () => {
    setLoading(true);
    try {
      const updated = await api.dismissSuggestion(suggestion.channel_id, suggestion.id);
      useChannelStore.getState().updateSuggestion(suggestion.channel_id, suggestion.id, updated);
    } catch {
      toast.error("Failed to dismiss suggestion");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      className={cn(
        "mx-2 my-2 rounded-lg border p-3",
        isPending && "border-amber-300 bg-amber-50 dark:border-amber-700 dark:bg-amber-950/30",
        suggestion.status === "approved" && "border-green-300 bg-green-50 dark:border-green-700 dark:bg-green-950/30",
        suggestion.status === "dismissed" && "border-muted bg-muted/30 opacity-60",
      )}
    >
      <div className="flex items-start gap-2">
        <Lightbulb className={cn(
          "mt-0.5 h-4 w-4 shrink-0",
          isPending ? "text-amber-600" : "text-muted-foreground",
        )} />
        <div className="min-w-0 flex-1">
          <div className="flex items-baseline gap-2">
            <span className="text-sm font-semibold">{suggestion.title}</span>
            {suggestion.priority !== "none" && (
              <span className="text-xs text-muted-foreground capitalize">
                {suggestion.priority}
              </span>
            )}
          </div>
          {suggestion.description && (
            <p className="mt-0.5 text-sm text-muted-foreground">{suggestion.description}</p>
          )}
          <div className="mt-1 flex items-center gap-2 text-xs text-muted-foreground">
            <span>Suggested by {suggestedBy}</span>
            {assigneeName && <span>· Assign to {assigneeName}</span>}
            {suggestion.status !== "pending" && (
              <span className="capitalize">· {suggestion.status}</span>
            )}
          </div>
        </div>
        {isPending && (
          <div className="flex gap-1 shrink-0">
            <Button
              variant="ghost"
              size="icon"
              className="h-7 w-7 text-green-600 hover:bg-green-100 hover:text-green-700"
              onClick={handleApprove}
              disabled={loading}
              title="Approve — creates issue and starts work"
            >
              <Check className="h-4 w-4" />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="h-7 w-7 text-muted-foreground hover:bg-destructive/10 hover:text-destructive"
              onClick={handleDismiss}
              disabled={loading}
              title="Dismiss suggestion"
            >
              <X className="h-4 w-4" />
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}
