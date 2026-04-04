"use client";

import { useEffect } from "react";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable";
import { useChannelStore } from "../store";
import { useWSEvent } from "@/features/realtime";
import type {
  ChannelCreatedPayload,
  ChannelUpdatedPayload,
  ChannelDeletedPayload,
  ChannelMessageCreatedPayload,
  ChannelMessageUpdatedPayload,
  ChannelMessageDeletedPayload,
} from "@/shared/types";
import { ChannelList } from "./channel-list";
import { ChatView } from "./chat-view";

export function ChatPage() {
  const fetch = useChannelStore((s) => s.fetch);

  useEffect(() => {
    fetch();
  }, [fetch]);

  // Real-time event handlers
  useWSEvent("channel:created", (payload) => {
    const { channel } = payload as ChannelCreatedPayload;
    useChannelStore.getState().addChannel(channel);
  });

  useWSEvent("channel:updated", (payload) => {
    const { channel } = payload as ChannelUpdatedPayload;
    useChannelStore.getState().updateChannel(channel.id, channel);
  });

  useWSEvent("channel:deleted", (payload) => {
    const { channel_id } = payload as ChannelDeletedPayload;
    useChannelStore.getState().removeChannel(channel_id);
  });

  useWSEvent("channel:message_created", (payload) => {
    const { message } = payload as ChannelMessageCreatedPayload;
    useChannelStore.getState().addMessage(message.channel_id, message);
  });

  useWSEvent("channel:message_updated", (payload) => {
    const { message } = payload as ChannelMessageUpdatedPayload;
    useChannelStore.getState().updateMessage(message.channel_id, message.id, message);
  });

  useWSEvent("channel:message_deleted", (payload) => {
    const { message_id, channel_id } = payload as ChannelMessageDeletedPayload;
    useChannelStore.getState().removeMessage(channel_id, message_id);
  });

  // Refetch channel list when members change
  useWSEvent("channel:member_added", () => {
    useChannelStore.getState().fetch();
  });

  useWSEvent("channel:member_removed", () => {
    useChannelStore.getState().fetch();
  });

  return (
    <div className="h-[calc(100vh-56px)]">
      <ResizablePanelGroup orientation="horizontal">
        <ResizablePanel defaultSize={25} minSize={15} maxSize={40}>
          <ChannelList />
        </ResizablePanel>
        <ResizableHandle withHandle />
        <ResizablePanel defaultSize={75}>
          <ChatView />
        </ResizablePanel>
      </ResizablePanelGroup>
    </div>
  );
}
