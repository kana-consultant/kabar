// routes/provider/index.tsx
import { createFileRoute } from "@tanstack/react-router";
import AIManagementPage from "@/pages/ai-management/Provider/index";

export const Route = createFileRoute("/provider")({
    component: AIManagementPage,
});