"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useModalStore } from "@/features/modals";

function isInputFocused(): boolean {
  const el = document.activeElement;
  if (!el) return false;
  const tag = el.tagName.toLowerCase();
  if (tag === "input" || tag === "textarea" || tag === "select") return true;
  if ((el as HTMLElement).isContentEditable) return true;
  return false;
}

export function useGlobalShortcuts() {
  const router = useRouter();

  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      // Skip if user is typing in an input
      if (isInputFocused()) return;
      // Skip if modifier keys are held (except for specific combos)
      if (e.metaKey || e.ctrlKey || e.altKey) return;

      switch (e.key) {
        case "c":
          e.preventDefault();
          useModalStore.getState().open("create-issue");
          break;
        case "g":
          // Two-key combo: g then i/a/r/s/m
          {
            const onSecond = (e2: KeyboardEvent) => {
              if (isInputFocused()) return;
              document.removeEventListener("keydown", onSecond);
              clearTimeout(timeout);
              switch (e2.key) {
                case "i":
                  e2.preventDefault();
                  router.push("/issues");
                  break;
                case "a":
                  e2.preventDefault();
                  router.push("/agents");
                  break;
                case "r":
                  e2.preventDefault();
                  router.push("/runtimes");
                  break;
                case "s":
                  e2.preventDefault();
                  router.push("/settings");
                  break;
                case "m":
                  e2.preventDefault();
                  router.push("/my-issues");
                  break;
                case "n":
                  e2.preventDefault();
                  router.push("/inbox");
                  break;
              }
            };
            const timeout = setTimeout(() => {
              document.removeEventListener("keydown", onSecond);
            }, 1000);
            document.addEventListener("keydown", onSecond, { once: true });
          }
          break;
      }
    };

    document.addEventListener("keydown", handler);
    return () => document.removeEventListener("keydown", handler);
  }, [router]);
}
