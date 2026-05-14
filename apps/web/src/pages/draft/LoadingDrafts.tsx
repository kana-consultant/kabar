import { Skeleton } from "@/components/ui/skeleton";
import { Card, CardContent } from "@/components/ui/card";

export function LoadingDrafts() {
    return (
        <div className="space-y-6">
            {/* Header Loading */}
            <div className="flex items-center justify-between">
                <div>
                    <Skeleton className="h-8 w-48" />
                    <Skeleton className="mt-2 h-4 w-64" />
                </div>
                <div className="flex gap-2">
                    <Skeleton className="h-10 w-64" />
                    <Skeleton className="h-10 w-32" />
                </div>
            </div>

            {/* Stats Loading */}
            <div className="grid gap-4 md:grid-cols-4">
                {[1, 2, 3, 4].map((i) => (
                    <Card key={i}>
                        <CardContent className="p-4">
                            <div className="flex items-center justify-between">
                                <div>
                                    <Skeleton className="h-4 w-20" />
                                    <Skeleton className="mt-1 h-8 w-12" />
                                </div>
                                <Skeleton className="h-8 w-8 rounded-full" />
                            </div>
                        </CardContent>
                    </Card>
                ))}
            </div>

            {/* List Loading */}
            <Card>
                <CardContent className="p-6">
                    <div className="space-y-4">
                        {[1, 2, 3, 4].map((i) => (
                            <div key={i} className="flex items-center justify-between rounded-lg border p-4">
                                <div className="flex-1 space-y-2">
                                    <div className="flex items-center gap-3">
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
                                    <Skeleton className="h-8 w-8" />
                                </div>
                            </div>
                        ))}
                    </div>
                </CardContent>
            </Card>
        </div>
    );
}