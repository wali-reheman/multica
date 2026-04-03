"use client";

import { create } from "zustand";
import type { Label } from "@/shared/types";
import { api } from "@/shared/api";
import { toast } from "sonner";

interface LabelState {
  labels: Label[];
  loading: boolean;
  fetch: () => Promise<void>;
  addLabel: (label: Label) => void;
  updateLabel: (id: string, updates: Partial<Label>) => void;
  removeLabel: (id: string) => void;
}

export const useLabelStore = create<LabelState>((set, get) => ({
  labels: [],
  loading: false,

  fetch: async () => {
    if (get().labels.length === 0) set({ loading: true });
    try {
      const labels = await api.listLabels();
      set({ labels, loading: false });
    } catch {
      toast.error("Failed to load labels");
      set({ loading: false });
    }
  },

  addLabel: (label) =>
    set((s) => ({
      labels: s.labels.some((l) => l.id === label.id)
        ? s.labels
        : [...s.labels, label],
    })),

  updateLabel: (id, updates) =>
    set((s) => ({
      labels: s.labels.map((l) => (l.id === id ? { ...l, ...updates } : l)),
    })),

  removeLabel: (id) =>
    set((s) => ({ labels: s.labels.filter((l) => l.id !== id) })),
}));
