import { create } from "zustand";
import { toast } from "sonner";
import { api } from "@/shared/api";
import type { Channel, ChannelMessage } from "@/shared/types";

interface ChannelState {
  channels: Channel[];
  loading: boolean;
  activeChannelId: string | null;
  messages: Record<string, ChannelMessage[]>;
  messagesLoading: Record<string, boolean>;

  fetch: () => Promise<void>;
  setChannels: (channels: Channel[]) => void;
  addChannel: (channel: Channel) => void;
  updateChannel: (id: string, updates: Partial<Channel>) => void;
  removeChannel: (id: string) => void;
  setActiveChannel: (id: string | null) => void;

  fetchMessages: (channelId: string) => Promise<void>;
  addMessage: (channelId: string, message: ChannelMessage) => void;
  updateMessage: (channelId: string, messageId: string, updates: Partial<ChannelMessage>) => void;
  removeMessage: (channelId: string, messageId: string) => void;
}

export const useChannelStore = create<ChannelState>((set, get) => ({
  channels: [],
  loading: true,
  activeChannelId: null,
  messages: {},
  messagesLoading: {},

  fetch: async () => {
    const isInitialLoad = get().channels.length === 0;
    if (isInitialLoad) set({ loading: true });
    try {
      const channels = await api.listChannels();
      set({ channels, loading: false });
    } catch {
      toast.error("Failed to load channels");
      if (isInitialLoad) set({ loading: false });
    }
  },

  setChannels: (channels) => set({ channels }),

  addChannel: (channel) =>
    set((s) => ({
      channels: s.channels.some((c) => c.id === channel.id)
        ? s.channels
        : [...s.channels, channel],
    })),

  updateChannel: (id, updates) =>
    set((s) => ({
      channels: s.channels.map((c) => (c.id === id ? { ...c, ...updates } : c)),
    })),

  removeChannel: (id) =>
    set((s) => ({
      channels: s.channels.filter((c) => c.id !== id),
      activeChannelId: s.activeChannelId === id ? null : s.activeChannelId,
    })),

  setActiveChannel: (id) => set({ activeChannelId: id }),

  fetchMessages: async (channelId: string) => {
    set((s) => ({ messagesLoading: { ...s.messagesLoading, [channelId]: true } }));
    try {
      const messages = await api.listChannelMessages(channelId, { limit: 100 });
      set((s) => ({
        messages: { ...s.messages, [channelId]: messages },
        messagesLoading: { ...s.messagesLoading, [channelId]: false },
      }));
    } catch {
      toast.error("Failed to load messages");
      set((s) => ({ messagesLoading: { ...s.messagesLoading, [channelId]: false } }));
    }
  },

  addMessage: (channelId, message) =>
    set((s) => {
      const existing = s.messages[channelId] ?? [];
      if (existing.some((m) => m.id === message.id)) return s;
      return {
        messages: { ...s.messages, [channelId]: [...existing, message] },
      };
    }),

  updateMessage: (channelId, messageId, updates) =>
    set((s) => {
      const existing = s.messages[channelId] ?? [];
      return {
        messages: {
          ...s.messages,
          [channelId]: existing.map((m) =>
            m.id === messageId ? { ...m, ...updates } : m,
          ),
        },
      };
    }),

  removeMessage: (channelId, messageId) =>
    set((s) => {
      const existing = s.messages[channelId] ?? [];
      return {
        messages: {
          ...s.messages,
          [channelId]: existing.filter((m) => m.id !== messageId),
        },
      };
    }),
}));
