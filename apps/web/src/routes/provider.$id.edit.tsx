// routes/provider/$id.edit.tsx
import { createFileRoute } from "@tanstack/react-router";
import EditProviderPage from "@/pages/ai-management/Provider/[id]/edit";

export const Route = createFileRoute("/provider/$id/edit")({
    component: EditProviderPage,
});