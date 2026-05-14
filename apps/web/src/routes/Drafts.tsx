import { createFileRoute } from "@tanstack/react-router";
import Drafts from "@/pages/draft/Drafts";

export const Route = createFileRoute("/Drafts")({
    component: Drafts,
});