"use client";

import { useState, useCallback } from "react";
import { Loader2 } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import { api } from "@/shared/api";
import { toast } from "sonner";
import { useProjectStore } from "../store";

export function AddProjectDialog({
  open,
  onOpenChange,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  const addProject = useProjectStore((s) => s.addProject);
  const setActiveProject = useProjectStore((s) => s.setActiveProject);
  const [name, setName] = useState("");
  const [localPath, setLocalPath] = useState("");
  const [initGit, setInitGit] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = useCallback(
    async (e: React.FormEvent) => {
      e.preventDefault();
      if (!localPath.trim()) return;

      setSubmitting(true);
      try {
        const project = await api.createProject({
          name: name.trim() || undefined,
          local_path: localPath.trim(),
          init_git: initGit,
        });
        addProject(project);
        setActiveProject(project.id);
        toast.success("Project added");
        setName("");
        setLocalPath("");
        setInitGit(false);
        onOpenChange(false);
      } catch {
        toast.error("Failed to add project");
      } finally {
        setSubmitting(false);
      }
    },
    [name, localPath, initGit, addProject, setActiveProject, onOpenChange],
  );

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add Project</DialogTitle>
          <DialogDescription>
            Register a local directory as a project
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="project-path">Local path</Label>
            <Input
              id="project-path"
              placeholder="/path/to/project"
              value={localPath}
              onChange={(e) => setLocalPath(e.target.value)}
              autoFocus
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="project-name">Name (optional)</Label>
            <Input
              id="project-name"
              placeholder="Derived from directory name if empty"
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
          </div>
          <label className="flex items-center gap-2 cursor-pointer">
            <Checkbox
              checked={initGit}
              onCheckedChange={() => setInitGit(!initGit)}
            />
            <span className="text-sm">Initialize git repository</span>
          </label>
          <DialogFooter>
            <Button
              type="submit"
              disabled={!localPath.trim() || submitting}
            >
              {submitting && <Loader2 className="h-4 w-4 animate-spin" />}
              Add Project
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
