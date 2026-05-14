import type { StatusData, ActionData } from "./types";

export function formatDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString("id-ID", {
        day: "numeric",
        month: "long",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit",
    });
}

// Hanya data mentah, bukan JSX!
export function getStatusData(status: string): StatusData {
    switch (status) {
        case "published":
            return { label: "Berhasil", icon: "✅", color: "green" };
        case "failed":
            return { label: "Gagal", icon: "❌", color: "red" };
        default:
            return { label: "Pending", icon: "⏳", color: "yellow" };
    }
}

export function getActionData(action: string): ActionData {
    switch (action) {
        case "published":
            return { label: "Publikasi", icon: "🚀" };
        case "scheduled":
            return { label: "Terjadwal", icon: "📅" };
        default:
            return { label: "Draft", icon: "📝" };
    }
}