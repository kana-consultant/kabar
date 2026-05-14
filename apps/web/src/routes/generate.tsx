import { createFileRoute } from "@tanstack/react-router";
import Generate from "@/pages/generate/generate";

export const Route = createFileRoute("/generate")({
    component: Generate,
});