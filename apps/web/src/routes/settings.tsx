import { createFileRoute } from '@tanstack/react-router';
import Settings from "@/pages/settings/Settings";

export const Route = createFileRoute("/settings")({
  component: Settings,
});