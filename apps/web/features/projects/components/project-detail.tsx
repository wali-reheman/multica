"use client";

import { useEffect, useState } from "react";
import { GitBranch, GitCommit, FileText, FolderGit2 } from "lucide-react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import type { Project } from "@/shared/types";
import { useProjectStore } from "../store";
import { CommitHistory } from "./commit-history";
import { FileStatus } from "./file-status";
import { BranchSelector } from "./branch-selector";
import { CommitDialog } from "./commit-dialog";
import { useWSEvent } from "@/features/realtime";

export function ProjectDetail({ project }: { project: Project }) {
  const fetchCommits = useProjectStore((s) => s.fetchCommits);
  const fetchBranches = useProjectStore((s) => s.fetchBranches);
  const fetchStatus = useProjectStore((s) => s.fetchStatus);
  const [commitDialogOpen, setCommitDialogOpen] = useState(false);

  useEffect(() => {
    if (project.is_git_repo) {
      fetchCommits(project.id);
      fetchBranches(project.id);
      fetchStatus(project.id);
    }
  }, [project.id, project.is_git_repo, fetchCommits, fetchBranches, fetchStatus]);

  // Auto-refresh status when files change
  useWSEvent("project:files_changed", (payload: unknown) => {
    const data = payload as { project_id?: string } | undefined;
    if (data?.project_id === project.id) {
      fetchStatus(project.id);
    }
  });

  if (!project.is_git_repo) {
    return (
      <div className="flex h-full flex-col items-center justify-center gap-3 p-6 text-muted-foreground">
        <FolderGit2 className="size-10 opacity-40" />
        <p className="text-sm">This project is not a git repository</p>
        <p className="text-xs">Initialize git to enable version history</p>
      </div>
    );
  }

  return (
    <div className="flex h-full flex-col">
      {/* Project header */}
      <div className="flex items-center justify-between border-b px-6 py-3">
        <div>
          <h2 className="text-sm font-medium">{project.name}</h2>
          <p className="text-xs text-muted-foreground">{project.local_path}</p>
        </div>
        <BranchSelector projectId={project.id} />
      </div>

      {/* Tabs */}
      <Tabs defaultValue="history" className="flex flex-1 flex-col overflow-hidden">
        <TabsList className="mx-6 mt-3 w-fit">
          <TabsTrigger value="history">
            <GitCommit className="mr-1.5 size-3.5" />
            History
          </TabsTrigger>
          <TabsTrigger value="changes">
            <FileText className="mr-1.5 size-3.5" />
            Changes
          </TabsTrigger>
        </TabsList>

        <TabsContent value="history" className="flex-1 overflow-y-auto px-6 py-3">
          <CommitHistory projectId={project.id} />
        </TabsContent>

        <TabsContent value="changes" className="flex-1 overflow-y-auto px-6 py-3">
          <FileStatus
            projectId={project.id}
            onCommit={() => setCommitDialogOpen(true)}
          />
        </TabsContent>
      </Tabs>

      <CommitDialog
        projectId={project.id}
        open={commitDialogOpen}
        onOpenChange={setCommitDialogOpen}
      />
    </div>
  );
}
