import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { FileText } from "lucide-react";
import { useNavigate } from "@tanstack/react-router";
import { DraftItem } from "./DraftItem";
import type { Draft } from "@/types/draft";

interface DraftListProps {
    drafts: Draft[];
    onView: (draft: Draft) => void;
    onEdit: (draft: Draft) => void;
    onSchedule: (draft: Draft) => void;
    onPublishNow: (draft: Draft) => void;
    onDelete: (draft: Draft) => void;
    formatDate: (date: string) => string;
}

export function DraftList({
    drafts,
    onView,
    onEdit,
    onSchedule,
    onPublishNow,
    onDelete,
    formatDate,
}: DraftListProps) {
    const navigate = useNavigate();

    if (drafts.length === 0) {
        return (
            <Card>
                <CardContent className="py-12 text-center">
                    <FileText className="mx-auto h-12 w-12 text-slate-300" />
                    <p className="mt-2 text-slate-500">Belum ada draft</p>
                    <Button
                        variant="outline"
                        className="mt-4"
                        onClick={() => navigate({ to: "/generate" })}
                    >
                        Buat Draft Baru
                    </Button>
                </CardContent>
            </Card>
        );
    }

    return (
        <Card>
            <CardHeader>
                <CardTitle>Daftar Draft</CardTitle>
                <CardDescription>{drafts.length} draft ditemukan</CardDescription>
            </CardHeader>
            <CardContent>
                <div className="space-y-4">
                    {drafts.map((draft) => (
                        <DraftItem
                            key={draft.id}
                            draft={draft}
                            onView={onView}
                            onEdit={onEdit}
                            onSchedule={onSchedule}
                            onPublishNow={onPublishNow}
                            onDelete={onDelete}
                            formatDate={formatDate}
                        />
                    ))}
                </div>
            </CardContent>
        </Card>
    );
}