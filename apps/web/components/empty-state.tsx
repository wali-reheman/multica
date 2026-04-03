"use client";

import type { ReactNode } from "react";
import type { LucideIcon } from "lucide-react";

interface EmptyStateProps {
  icon: LucideIcon;
  title: string;
  description?: string;
  action?: ReactNode;
}

export function EmptyState({ icon: Icon, title, description, action }: EmptyStateProps) {
  return (
    <div className="flex flex-1 min-h-0 flex-col items-center justify-center gap-3 p-8">
      <div className="flex items-center justify-center size-12 rounded-full bg-muted">
        <Icon className="size-6 text-muted-foreground" />
      </div>
      <div className="text-center space-y-1">
        <h3 className="text-sm font-medium">{title}</h3>
        {description && (
          <p className="text-sm text-muted-foreground max-w-sm">{description}</p>
        )}
      </div>
      {action}
    </div>
  );
}
