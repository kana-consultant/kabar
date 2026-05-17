import { createFileRoute } from "@tanstack/react-router";
import InvitePage from "@/pages/Invite";
export const Route = createFileRoute("/invite")({
    component: InvitePage,
});
