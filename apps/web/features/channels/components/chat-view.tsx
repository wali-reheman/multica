"use client";

import { useEffect, useRef } from "react";
import { Hash, Users, MessageSquarePlus, Loader2 } from "lucide-react";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { useChannelStore } from "../store";
import { MessageItem } from "./message-item";
import { MessageComposer } from "./message-composer";
import { ChannelMembersPopover } from "./channel-members-popover";
import { SuggestionCard } from "./suggestion-card";

export function ChatView() {
  const activeChannelId = useChannelStore((s) => s.activeChannelId);
  const channels = useChannelStore((s) => s.channels);
  const messages = useChannelStore((s) => s.messages);
  const messagesLoading = useChannelStore((s) => s.messagesLoading);
  const fetchMessages = useChannelStore((s) => s.fetchMessages);
  const suggestions = useChannelStore((s) => s.suggestions);
  const fetchSuggestions = useChannelStore((s) => s.fetchSuggestions);

  const channel = channels.find((c) => c.id === activeChannelId);
  const channelMessages = activeChannelId ? (messages[activeChannelId] ?? []) : [];
  const isLoading = activeChannelId ? (messagesLoading[activeChannelId] ?? false) : false;
  const pendingSuggestions = activeChannelId
    ? (suggestions[activeChannelId] ?? []).filter((s) => s.status === "pending")
    : [];

  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (activeChannelId) {
      fetchMessages(activeChannelId);
      fetchSuggestions(activeChannelId);
    }
  }, [activeChannelId, fetchMessages, fetchSuggestions]);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [channelMessages.length]);

  if (!channel) {
    return (
      <div className="flex h-full items-center justify-center bg-background text-muted-foreground">
        <div className="text-center space-y-3">
          <div className="flex h-16 w-16 items-center justify-center rounded-full bg-muted mx-auto">
            <MessageSquarePlus className="h-8 w-8 text-muted-foreground/60" />
          </div>
          <div>
            <p className="text-lg font-medium text-foreground">Welcome to Chat</p>
            <p className="mt-1 text-sm max-w-[280px]">
              Select a channel from the sidebar, or create a new one to start collaborating with your team and agents.
            </p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="flex h-full flex-col bg-background">
      {/* Channel header */}
      <div className="flex items-center justify-between border-b px-5 py-3">
        <div className="flex items-center gap-2 min-w-0">
          {channel.type === "direct" ? (
            <Users className="h-4 w-4 shrink-0 text-muted-foreground" />
          ) : (
            <Hash className="h-5 w-5 shrink-0 text-muted-foreground" />
          )}
          <h2 className="font-semibold truncate">{channel.name}</h2>
          {channel.description && (
            <>
              <Separator orientation="vertical" className="h-4 mx-1" />
              <span className="text-sm text-muted-foreground truncate">
                {channel.description}
              </span>
            </>
          )}
        </div>
        <div className="flex items-center gap-1 shrink-0">
          <ChannelMembersPopover channelId={channel.id} members={channel.members ?? []} />
        </div>
      </div>

      {/* Pending suggestions */}
      {pendingSuggestions.length > 0 && (
        <div className="border-b bg-muted/30 px-4 py-3 space-y-2">
          <p className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
            Pending suggestions
          </p>
          {pendingSuggestions.map((s) => (
            <SuggestionCard key={s.id} suggestion={s} />
          ))}
        </div>
      )}

      {/* Messages */}
      <ScrollArea className="flex-1">
        <div className="px-5">
          {isLoading ? (
            <div className="flex items-center justify-center py-16">
              <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
            </div>
          ) : channelMessages.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-16 text-center">
              <div className="flex h-12 w-12 items-center justify-center rounded-full bg-muted mb-3">
                <Hash className="h-6 w-6 text-muted-foreground" />
              </div>
              <p className="font-medium text-foreground">
                This is the start of #{channel.name}
              </p>
              <p className="mt-1 text-sm text-muted-foreground max-w-[320px]">
                Send a message to start the conversation. @mention an agent to bring them into the discussion.
              </p>
            </div>
          ) : (
            <div className="space-y-0.5 py-4">
              {channelMessages.map((msg, i) => {
                const prev = i > 0 ? channelMessages[i - 1] : null;
                const isGrouped = prev
                  && prev.author_type === msg.author_type
                  && prev.author_id === msg.author_id
                  && prev.type === "message"
                  && msg.type === "message"
                  && new Date(msg.created_at).getTime() - new Date(prev.created_at).getTime() < 5 * 60 * 1000;
                return (
                  <MessageItem key={msg.id} message={msg} grouped={!!isGrouped} />
                );
              })}
              <div ref={bottomRef} />
            </div>
          )}
        </div>
      </ScrollArea>

      {/* Composer */}
      <MessageComposer channelId={channel.id} channelName={channel.name} />
    </div>
  );
}
