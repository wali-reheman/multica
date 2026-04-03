"use client";

import { useState } from "react";
import { toast } from "sonner";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import { api } from "@/shared/api";
import { useProjectStore } from "../store";

interface AddProjectDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function AddProjectDialog({ open, onOpenChange }: AddProjectDialogProps) {
  const [localPath, setLocalPath] = useState("");
  const [name, setName] = useState("");
  const [initGit, setInitGit] = useState(true);
  const [loading, setLoading] = useState(false);
  const addProject = useProjectStore((s) => s.addProject);
  const setActiveProject = useProjectStore((s) => s.setActiveProject);

  const handleSubmit = async () => {
    if (!localPath.trim()) return;
    setLoading(true);
    try {
      const project = await api.createProject({
        local_path: localPath.trim(),
        name: name.trim() || undefined,
        init_git: initGit,
      });
      addProject(project);
      setActiveProject(project.id);
      toast.success(`Project "${project.name}" added`);
      setLocalPath("");
      setName("");
      onOpenChange(false);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to add project");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Add Project</DialogTitle>
        </DialogHeader>
        <div className="flex flex-col gap-4">
          <div className="flex flex-col gap-1.5">
            <Label htmlFor="local-path" className="text-xs">
              Local Path
            </Label>
            <Input
              id="local-path"
              placeholder="/path/to/your/project"
              value={localPath}
              onChange={(e) => setLocalPath(e.target.value)}
              autoFocus
            />
          </div>
          <div className="flex flex-col gap-1.5">
            <Label htmlFor="project-name" className="text-xs">
              Name (optional)
            </Label>
            <Input
              id="project-name"
              placeholder="Auto-detected from folder name"
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
          </div>
          <div className="flex items-center gap-2">
            <Checkbox
              id="init-git"
              checked={initGit}
              onCheckedChange={(checked) => setInitGit(!!checked)}
            />
            <Label htmlFor="init-git" className="text-xs">
              Initialize git if not a repository
            </Label>
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button onClick={handleSubmit} disabled={!localPath.trim() || loading}>
            {loading ? "Adding..." : "Add Project"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
