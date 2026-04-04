"use client";

import { useState, useRef, useCallback } from "react";
import { SendHorizontal, Plus, Lightbulb, ListPlus } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { api } from "@/shared/api";
import { useChannelStore } from "../store";
import { CreateIssueFromChatDialog } from "./create-issue-dialog";
import { SuggestTaskDialog } from "./suggest-task-dialog";

interface MessageComposerProps {
  channelId: string;
  channelName: string;
}

export function MessageComposer({ channelId, channelName }: MessageComposerProps) {
  const [content, setContent] = useState("");
  const [sending, setSending] = useState(false);
  const [showCreateIssue, setShowCreateIssue] = useState(false);
  const [showSuggestTask, setShowSuggestTask] = useState(false);
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const addMessage = useChannelStore((s) => s.addMessage);

  const handleSend = useCallback(async () => {
    const trimmed = content.trim();
    if (!trimmed || sending) return;

    setSending(true);
    try {
      const message = await api.sendChannelMessage(channelId, { content: trimmed });
      addMessage(channelId, message);
      setContent("");
      textareaRef.current?.focus();
    } catch {
      toast.error("Failed to send message");
    } finally {
      setSending(false);
    }
  }, [content, sending, channelId, addMessage]);

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  return (
    <>
      <div className="border-t bg-background px-5 py-3">
        <div className="flex items-end gap-2 rounded-lg border bg-muted/30 p-1.5 focus-within:ring-1 focus-within:ring-ring">
          <DropdownMenu>
            <DropdownMenuTrigger
              render={
                <Button variant="ghost" size="icon" className="h-7 w-7 shrink-0 text-muted-foreground hover:text-foreground">
                  <Plus className="h-4 w-4" />
                </Button>
              }
            />
            <DropdownMenuContent align="start" sideOffset={8}>
              <DropdownMenuItem onClick={() => setShowSuggestTask(true)}>
                <Lightbulb className="mr-2 h-4 w-4" />
                Suggest a task
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={() => setShowCreateIssue(true)}>
                <ListPlus className="mr-2 h-4 w-4" />
                Create issue directly
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
          <Textarea
            ref={textareaRef}
            placeholder={`Message #${channelName}...`}
            value={content}
            onChange={(e) => setContent(e.target.value)}
            onKeyDown={handleKeyDown}
            rows={1}
            className="min-h-[36px] max-h-[160px] resize-none border-0 bg-transparent shadow-none focus-visible:ring-0 text-sm px-1"
          />
          <Button
            size="icon"
            className="h-7 w-7 shrink-0"
            onClick={handleSend}
            disabled={!content.trim() || sending}
          >
            <SendHorizontal className="h-3.5 w-3.5" />
          </Button>
        </div>
        <p className="mt-1.5 text-[11px] text-muted-foreground px-1">
          <kbd className="rounded border bg-muted px-1 py-0.5 text-[10px]">Enter</kbd> to send
          <span className="mx-1.5">·</span>
          <kbd className="rounded border bg-muted px-1 py-0.5 text-[10px]">Shift+Enter</kbd> for new line
        </p>
      </div>
      <CreateIssueFromChatDialog
        channelId={channelId}
        open={showCreateIssue}
        onOpenChange={setShowCreateIssue}
      />
      <SuggestTaskDialog
        channelId={channelId}
        open={showSuggestTask}
        onOpenChange={setShowSuggestTask}
      />
    </>
  );
}
