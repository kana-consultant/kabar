import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/landing")({
    component: LandingPage,
});

export function LandingPage() {
    return (
        <div>
            Landing Page
        </div>
    );
}