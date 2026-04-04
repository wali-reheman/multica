"use client";

import { useState, useRef, useCallback } from "react";
import { SendHorizontal, ListPlus } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { api } from "@/shared/api";
import { useChannelStore } from "../store";
import { useIssueStore } from "@/features/issues";
import { CreateIssueFromChatDialog } from "./create-issue-dialog";

interface MessageComposerProps {
  channelId: string;
  channelName: string;
}

export function MessageComposer({ channelId, channelName }: MessageComposerProps) {
  const [content, setContent] = useState("");
  const [sending, setSending] = useState(false);
  const [showCreateIssue, setShowCreateIssue] = useState(false);
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
      <div className="border-t p-4">
        <div className="flex items-end gap-2">
          <DropdownMenu>
            <DropdownMenuTrigger
              render={
                <Button variant="ghost" size="icon" className="h-8 w-8 shrink-0">
                  <ListPlus className="h-4 w-4" />
                </Button>
              }
            />
            <DropdownMenuContent align="start">
              <DropdownMenuItem onClick={() => setShowCreateIssue(true)}>
                Create issue from chat
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
          <Textarea
            ref={textareaRef}
            placeholder={`Message #${channelName}`}
            value={content}
            onChange={(e) => setContent(e.target.value)}
            onKeyDown={handleKeyDown}
            rows={1}
            className="min-h-[40px] max-h-[200px] resize-none"
          />
          <Button
            size="icon"
            className="h-8 w-8 shrink-0"
            onClick={handleSend}
            disabled={!content.trim() || sending}
          >
            <SendHorizontal className="h-4 w-4" />
          </Button>
        </div>
        <p className="mt-1 text-xs text-muted-foreground">
          Press Enter to send, Shift+Enter for new line. Use @agent-name to mention agents.
        </p>
      </div>
      <CreateIssueFromChatDialog
        channelId={channelId}
        open={showCreateIssue}
        onOpenChange={setShowCreateIssue}
      />
    </>
  );
}
