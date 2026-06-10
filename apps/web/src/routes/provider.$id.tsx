// routes/provider/$id.tsx
import { createFileRoute } from "@tanstack/react-router";
import ProviderDetailPage from "@/pages/ai-management/Provider/[id]/index";

export const Route = createFileRoute("/provider/$id")({
    component: ProviderDetailPage,
});