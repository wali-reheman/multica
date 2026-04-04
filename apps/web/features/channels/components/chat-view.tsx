"use client";

import { useEffect, useRef } from "react";
import { Hash, Users, Settings, Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
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
      <div className="flex h-full items-center justify-center text-muted-foreground">
        <div className="text-center space-y-2">
          <Hash className="h-12 w-12 mx-auto text-muted-foreground/50" />
          <p className="text-lg font-medium">Select a channel</p>
          <p className="text-sm">Choose a channel from the sidebar or create a new one to start chatting.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex h-full flex-col">
      {/* Channel header */}
      <div className="flex items-center justify-between border-b px-4 py-3">
        <div className="flex items-center gap-2">
          {channel.type === "direct" ? (
            <Users className="h-4 w-4 text-muted-foreground" />
          ) : (
            <Hash className="h-4 w-4 text-muted-foreground" />
          )}
          <h2 className="font-semibold">{channel.name}</h2>
          {channel.description && (
            <span className="text-sm text-muted-foreground">
              — {channel.description}
            </span>
          )}
        </div>
        <div className="flex items-center gap-1">
          <ChannelMembersPopover channelId={channel.id} members={channel.members ?? []} />
        </div>
      </div>

      {/* Pending suggestions */}
      {pendingSuggestions.length > 0 && (
        <div className="border-b px-2 py-2 space-y-1">
          {pendingSuggestions.map((s) => (
            <SuggestionCard key={s.id} suggestion={s} />
          ))}
        </div>
      )}

      {/* Messages */}
      <ScrollArea className="flex-1 px-4">
        {isLoading ? (
          <div className="flex items-center justify-center py-12 text-muted-foreground">
            Loading messages...
          </div>
        ) : channelMessages.length === 0 ? (
          <div className="flex items-center justify-center py-12 text-muted-foreground">
            <div className="text-center space-y-1">
              <p className="font-medium">No messages yet</p>
              <p className="text-sm">Start the conversation! Agents in this channel will respond when @mentioned.</p>
            </div>
          </div>
        ) : (
          <div className="space-y-1 py-4">
            {channelMessages.map((msg) => (
              <MessageItem key={msg.id} message={msg} />
            ))}
            <div ref={bottomRef} />
          </div>
        )}
      </ScrollArea>

      {/* Composer */}
      <MessageComposer channelId={channel.id} channelName={channel.name} />
    </div>
  );
}
