// routes/provider/add.tsx
import { createFileRoute } from "@tanstack/react-router";
import CreateProviderPage from "@/pages/ai-management/Provider/create";

export const Route = createFileRoute("/provider/add")({
    component: CreateProviderPage,
});