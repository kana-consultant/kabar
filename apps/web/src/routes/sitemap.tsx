// routes/sitemap.tsx
import { createFileRoute } from "@tanstack/react-router";
import SitemapPage from "@/pages/sitemap/Sitemap";

export const Route = createFileRoute("/sitemap")({
    component: SitemapPage,
});