// src/components/ui/LoadingCard.tsx
import { Skeleton } from "@/components/ui/skeleton";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Loader2 } from "lucide-react";

interface LoadingCardProps {
    count?: number;
    className?: string;
    variant?: "grid" | "list" | "spinner";
    message?: string;
}

export function LoadingCard({
    count = 3,
    className = "",
    variant = "grid",
    message = "Memuat data..."
}: LoadingCardProps) {

    // Variant Spinner (full page / center)
    if (variant === "spinner") {
        return (
            <div className="flex flex-col items-center justify-center min-h-[400px] gap-4">
                <div className="relative">
                    <div className="absolute inset-0 rounded-full bg-blue-500/20 animate-ping" />
                    <Loader2 className="h-10 w-10 animate-spin text-blue-500" />
                </div>
                <p className="text-sm text-slate-500 animate-pulse">{message}</p>
            </div>
        );
    }

    // Variant List (horizontal items)
    if (variant === "list") {
        return (
            <div className="space-y-4">
                {Array.from({ length: count }).map((_, i) => (
                    <div
                        key={i}
                        className="flex items-center justify-between rounded-lg border p-4 animate-pulse"
                        style={{ animationDelay: `${i * 0.1}s` }}
                    >
                        <div className="flex-1 space-y-2">
                            <div className="flex items-center gap-3">
                                <Skeleton className="h-5 w-5 rounded-full" />
                                <Skeleton className="h-5 w-40" />
                                <Skeleton className="h-5 w-20 rounded-full" />
                            </div>
                            <Skeleton className="h-4 w-full max-w-md" />
                            <div className="flex gap-4">
                                <Skeleton className="h-3 w-32" />
                                <Skeleton className="h-3 w-24" />
                                <Skeleton className="h-3 w-20" />
                            </div>
                        </div>
                        <div className="flex gap-2">
                            <Skeleton className="h-8 w-8" />
                            <Skeleton className="h-8 w-8" />
                            <Skeleton className="h-8 w-8" />
                            <Skeleton className="h-8 w-24" />
                        </div>
                    </div>
                ))}
            </div>
        );
    }

    // Variant Grid (default)
    const renderCard = (index: number) => (
        <Card
            className="overflow-hidden transition-all duration-300 animate-pulse"
            style={{ animationDelay: `${index * 0.05}s` }}
        >
            <CardHeader className="pb-3">
                <div className="flex items-start justify-between">
                    <div className="flex items-center gap-2">
                        <Skeleton className="h-5 w-5 rounded-full" />
                        <Skeleton className="h-5 w-32" />
                    </div>
                    <Skeleton className="h-6 w-16 rounded-full" />
                </div>
                <Skeleton className="mt-2 h-4 w-24" />
            </CardHeader>
            <CardContent className="space-y-3">
                <Skeleton className="h-10 w-full rounded-lg" />
                <div className="flex justify-between">
                    <Skeleton className="h-4 w-20" />
                    <Skeleton className="h-4 w-16" />
                </div>
                <div className="flex gap-2 pt-2">
                    <Skeleton className="h-8 flex-1 rounded-md" />
                    <Skeleton className="h-8 flex-1 rounded-md" />
                    <Skeleton className="h-8 flex-1 rounded-md" />
                </div>
            </CardContent>
        </Card>
    );

    // Determine grid columns based on count
    const getGridCols = () => {
        if (count === 1) return "grid-cols-1";
        if (count === 2) return "sm:grid-cols-2";
        if (count >= 3) return "sm:grid-cols-2 lg:grid-cols-3";
        return "sm:grid-cols-2 lg:grid-cols-3";
    };

    return (
        <div className={`w-full ${className}`}>
            {/* Loading Indicator Bar */}
            <div className="mb-6 flex items-center justify-between">
                <div className="flex items-center gap-2">
                    <div className="h-4 w-4 rounded-full bg-blue-500 animate-pulse" />
                    <span className="text-sm text-slate-400">{message}</span>
                </div>
                <div className="flex gap-1">
                    <div className="h-1.5 w-1.5 rounded-full bg-blue-500 animate-bounce" style={{ animationDelay: "0s" }} />
                    <div className="h-1.5 w-1.5 rounded-full bg-blue-500 animate-bounce" style={{ animationDelay: "0.2s" }} />
                    <div className="h-1.5 w-1.5 rounded-full bg-blue-500 animate-bounce" style={{ animationDelay: "0.4s" }} />
                </div>
            </div>

            {/* Cards Grid */}
            <div className={`grid gap-4 ${getGridCols()}`}>
                {Array.from({ length: count }).map((_, i) => (
                    <div key={i} className="animate-fade-in">
                        {renderCard(i)}
                    </div>
                ))}
            </div>
        </div>
    );
}