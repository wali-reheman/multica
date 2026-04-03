"use client";

import { X } from "lucide-react";

interface LabelBadgeProps {
  name: string;
  color: string;
  onRemove?: () => void;
}

export function LabelBadge({ name, color, onRemove }: LabelBadgeProps) {
  return (
    <span
      className="inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs font-medium"
      style={{
        backgroundColor: `${color}20`,
        color: color,
        border: `1px solid ${color}40`,
      }}
    >
      <span
        className="size-2 rounded-full shrink-0"
        style={{ backgroundColor: color }}
      />
      {name}
      {onRemove && (
        <button
          type="button"
          onClick={(e) => {
            e.stopPropagation();
            onRemove();
          }}
          className="ml-0.5 rounded-full p-0.5 hover:bg-black/10 dark:hover:bg-white/10"
        >
          <X className="size-2.5" />
        </button>
      )}
    </span>
  );
}
