// routes/admin.ai-management.tsx
import { createFileRoute } from "@tanstack/react-router";
import AIManagement from "@/pages/ai-management/AIManagement";

export const Route = createFileRoute("/ai-management")({
    component: AIManagement,
});