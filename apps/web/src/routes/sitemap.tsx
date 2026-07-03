// routes/sitemap.tsx
import { createFileRoute } from "@tanstack/react-router";
import Sitemap from "@/pages/sitemap/Sitemap";

export const Route = createFileRoute("/sitemap")({
    component: Sitemap,
});