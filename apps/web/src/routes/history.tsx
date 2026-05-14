import { createFileRoute } from "@tanstack/react-router";
import History from "@/pages/history/History";

export const Route = createFileRoute("/history")({
    component: History,
});