"use client";

import { useState } from "react";
import { Hash, Plus, Users, MessageSquare, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Tooltip, TooltipTrigger, TooltipContent } from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";
import { useChannelStore } from "../store";
import { CreateChannelDialog } from "./create-channel-dialog";

export function ChannelList() {
  const channels = useChannelStore((s) => s.channels);
  const loading = useChannelStore((s) => s.loading);
  const activeChannelId = useChannelStore((s) => s.activeChannelId);
  const setActiveChannel = useChannelStore((s) => s.setActiveChannel);
  const [showCreate, setShowCreate] = useState(false);

  return (
    <div className="flex h-full flex-col bg-sidebar">
      {/* Header */}
      <div className="flex items-center justify-between border-b px-4 py-3">
        <div className="flex items-center gap-2">
          <MessageSquare className="h-4 w-4 text-muted-foreground" />
          <h2 className="text-sm font-semibold">Chat</h2>
        </div>
        <Tooltip>
          <TooltipTrigger
            className="flex h-6 w-6 items-center justify-center rounded-md hover:bg-accent"
            onClick={() => setShowCreate(true)}
          >
            <Plus className="h-4 w-4" />
          </TooltipTrigger>
          <TooltipContent side="bottom">New channel</TooltipContent>
        </Tooltip>
      </div>

      {/* Channel list */}
      <ScrollArea className="flex-1">
        {loading ? (
          <div className="flex items-center justify-center py-12">
            <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
          </div>
        ) : channels.length === 0 ? (
          <div className="flex flex-col items-center gap-3 px-4 py-12 text-center">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-muted">
              <Hash className="h-5 w-5 text-muted-foreground" />
            </div>
            <div>
              <p className="text-sm font-medium">No channels yet</p>
              <p className="mt-1 text-xs text-muted-foreground">
                Create a channel to start chatting with your team and agents.
              </p>
            </div>
            <Button
              size="sm"
              variant="outline"
              className="mt-1"
              onClick={() => setShowCreate(true)}
            >
              <Plus className="mr-1.5 h-3.5 w-3.5" />
              Create channel
            </Button>
          </div>
        ) : (
          <div className="space-y-0.5 p-2">
            {channels.map((channel) => {
              const isActive = activeChannelId === channel.id;
              const memberCount = channel.members?.length ?? 0;
              return (
                <button
                  key={channel.id}
                  className={cn(
                    "flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-sm transition-colors",
                    "text-muted-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground",
                    isActive && "bg-sidebar-accent text-sidebar-accent-foreground font-medium",
                  )}
                  onClick={() => setActiveChannel(channel.id)}
                >
                  {channel.type === "direct" ? (
                    <Users className="h-4 w-4 shrink-0" />
                  ) : (
                    <Hash className="h-4 w-4 shrink-0" />
                  )}
                  <span className="truncate">{channel.name}</span>
                  {memberCount > 0 && (
                    <span className="ml-auto text-[10px] tabular-nums text-muted-foreground">
                      {memberCount}
                    </span>
                  )}
                </button>
              );
            })}
          </div>
        )}
      </ScrollArea>
      <CreateChannelDialog open={showCreate} onOpenChange={setShowCreate} />
    </div>
  );
}
