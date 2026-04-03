"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { Plus, Tag } from "lucide-react";
import { toast } from "sonner";
import type { Label } from "@/shared/types";
import { api } from "@/shared/api";
import { useLabelStore } from "@/features/labels/store";
import {
  PropertyPicker,
  PickerItem,
  PickerEmpty,
} from "@/features/issues/components/pickers";
import { LabelBadge } from "./label-badge";

const PRESET_COLORS = [
  "#ef4444", "#f97316", "#eab308", "#22c55e",
  "#06b6d4", "#3b82f6", "#8b5cf6", "#ec4899",
];

interface LabelPickerProps {
  issueId: string;
  issueLabels: Label[];
  onLabelsChange: (labels: Label[]) => void;
}

export function LabelPicker({ issueId, issueLabels, onLabelsChange }: LabelPickerProps) {
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState("");
  const [creating, setCreating] = useState(false);
  const allLabels = useLabelStore((s) => s.labels);
  const fetchLabels = useLabelStore((s) => s.fetch);

  useEffect(() => {
    if (allLabels.length === 0) fetchLabels();
  }, [allLabels.length, fetchLabels]);

  const selectedIds = useMemo(
    () => new Set(issueLabels.map((l) => l.id)),
    [issueLabels],
  );

  const filteredLabels = useMemo(() => {
    if (!search) return allLabels;
    const q = search.toLowerCase();
    return allLabels.filter((l) => l.name.toLowerCase().includes(q));
  }, [allLabels, search]);

  const toggleLabel = useCallback(
    async (label: Label) => {
      const isSelected = selectedIds.has(label.id);
      if (isSelected) {
        onLabelsChange(issueLabels.filter((l) => l.id !== label.id));
        try {
          await api.removeIssueLabel(issueId, label.id);
        } catch {
          onLabelsChange(issueLabels);
          toast.error("Failed to remove label");
        }
      } else {
        onLabelsChange([...issueLabels, label]);
        try {
          await api.addIssueLabel(issueId, label.id);
        } catch {
          onLabelsChange(issueLabels);
          toast.error("Failed to add label");
        }
      }
    },
    [issueId, issueLabels, selectedIds, onLabelsChange],
  );

  const createAndAdd = useCallback(async () => {
    if (!search.trim() || creating) return;
    setCreating(true);
    try {
      const color = PRESET_COLORS[Math.floor(Math.random() * PRESET_COLORS.length)] ?? "#6b7280";
      const label = await api.createLabel({ name: search.trim(), color });
      useLabelStore.getState().addLabel(label);
      onLabelsChange([...issueLabels, label]);
      await api.addIssueLabel(issueId, label.id);
      setSearch("");
    } catch {
      toast.error("Failed to create label");
    } finally {
      setCreating(false);
    }
  }, [search, creating, issueId, issueLabels, onLabelsChange]);

  const exactMatch = filteredLabels.some(
    (l) => l.name.toLowerCase() === search.toLowerCase(),
  );

  return (
    <PropertyPicker
      open={open}
      onOpenChange={setOpen}
      trigger={
        issueLabels.length > 0 ? (
          <div className="flex flex-wrap gap-1">
            {issueLabels.map((l) => (
              <LabelBadge key={l.id} name={l.name} color={l.color} />
            ))}
          </div>
        ) : (
          <span className="text-sm text-muted-foreground flex items-center gap-1">
            <Tag className="size-3.5" />
            Add label
          </span>
        )
      }
      width="w-56"
      searchable
      searchPlaceholder="Search or create label..."
      onSearchChange={setSearch}
    >
      {filteredLabels.map((label) => (
        <PickerItem
          key={label.id}
          selected={selectedIds.has(label.id)}
          onClick={() => toggleLabel(label)}
        >
          <span
            className="size-3 rounded-full shrink-0"
            style={{ backgroundColor: label.color }}
          />
          <span className="truncate">{label.name}</span>
        </PickerItem>
      ))}
      {filteredLabels.length === 0 && !search && <PickerEmpty />}
      {search && !exactMatch && (
        <button
          type="button"
          onClick={createAndAdd}
          disabled={creating}
          className="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-sm hover:bg-accent transition-colors text-muted-foreground"
        >
          <Plus className="size-3.5" />
          <span>Create &ldquo;{search}&rdquo;</span>
        </button>
      )}
    </PropertyPicker>
  );
}
