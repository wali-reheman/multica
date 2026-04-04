"use client";

import { useState } from "react";
import { Hash, Plus, Users } from "lucide-react";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { cn } from "@/lib/utils";
import { useChannelStore } from "../store";
import { CreateChannelDialog } from "./create-channel-dialog";

export function ChannelList() {
  const channels = useChannelStore((s) => s.channels);
  const activeChannelId = useChannelStore((s) => s.activeChannelId);
  const setActiveChannel = useChannelStore((s) => s.setActiveChannel);
  const [showCreate, setShowCreate] = useState(false);

  return (
    <div className="flex h-full flex-col">
      <div className="flex items-center justify-between border-b px-4 py-3">
        <h2 className="text-sm font-semibold">Channels</h2>
        <Button variant="ghost" size="icon" className="h-6 w-6" onClick={() => setShowCreate(true)}>
          <Plus className="h-4 w-4" />
        </Button>
      </div>
      <ScrollArea className="flex-1">
        <div className="space-y-0.5 p-2">
          {channels.map((channel) => (
            <button
              key={channel.id}
              className={cn(
                "flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-sm transition-colors",
                "hover:bg-accent hover:text-accent-foreground",
                activeChannelId === channel.id && "bg-accent text-accent-foreground font-medium",
              )}
              onClick={() => setActiveChannel(channel.id)}
            >
              {channel.type === "direct" ? (
                <Users className="h-4 w-4 shrink-0 text-muted-foreground" />
              ) : (
                <Hash className="h-4 w-4 shrink-0 text-muted-foreground" />
              )}
              <span className="truncate">{channel.name}</span>
              {channel.members && (
                <span className="ml-auto text-xs text-muted-foreground">
                  {channel.members.length}
                </span>
              )}
            </button>
          ))}
          {channels.length === 0 && (
            <div className="px-2 py-8 text-center text-sm text-muted-foreground">
              No channels yet. Create one to start chatting with your agents.
            </div>
          )}
        </div>
      </ScrollArea>
      <CreateChannelDialog open={showCreate} onOpenChange={setShowCreate} />
    </div>
  );
}
