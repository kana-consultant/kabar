// components/ModelDetailModal.tsx
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@kana-consultant/ui-kit";
import { Badge } from "@kana-consultant/ui-kit";
import { Calendar, Zap, Users, Key, Check, X } from "lucide-react";
import type { AIModel } from "@/types/provider.types";

interface ModelDetailModalProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    model: AIModel | null;
}

export function ModelDetailModal({ open, onOpenChange, model }: ModelDetailModalProps) {
    if (!model) return null;

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="max-w-2xl">
                <DialogHeader>
                    <DialogTitle className="flex items-center justify-between">
                        <span>{model.display_name}</span>
                        <div className="flex gap-2">
                            {model.is_active ? (
                                <Badge tone="success" className="flex items-center gap-1">
                                    <Check className="w-3 h-3" /> Active
                                </Badge>
                            ) : (
                                <Badge tone="neutral" className="flex items-center gap-1">
                                    <X className="w-3 h-3" /> Inactive
                                </Badge>
                            )}
                            {model.is_default && (
                                <Badge tone="primary">Default</Badge>
                            )}
                        </div>
                    </DialogTitle>
                </DialogHeader>

                <div className="space-y-4">
                    {/* Basic Info */}
                    <div className="grid grid-cols-2 gap-4">
                        <InfoItem
                            label="Model Name"
                            value={model.name}
                        />
                        <InfoItem
                            label="Provider"
                            value={model.provider}
                        />
                        <InfoItem
                            label="Provider ID"
                            value={model.provider_id}
                        />
                    </div>

                    {/* Team Info */}
                    {model.team_id && (
                        <div className="border-t pt-4">
                            <InfoItem
                                label="Team ID"
                                value={model.team_id}
                                icon={<Users className="w-4 h-4" />}
                            />
                        </div>
                    )}

                    {/* Timestamps */}
                    <div className="border-t pt-4">
                        <div className="grid grid-cols-2 gap-4">
                            <InfoItem
                                label="Created At"
                                value={formatDate(model.created_at)}
                                icon={<Calendar className="w-4 h-4" />}
                            />
                            <InfoItem
                                label="Updated At"
                                value={formatDate(model.updated_at)}
                                icon={<Calendar className="w-4 h-4" />}
                            />
                        </div>
                    </div>
                </div>
            </DialogContent>
        </Dialog>
    );
}

interface InfoItemProps {
    label: string;
    value: string;
    icon?: React.ReactNode;
}

function InfoItem({ label, value, icon }: InfoItemProps) {
    return (
        <div className="space-y-1">
            <p className="text-sm text-muted-foreground flex items-center gap-1">
                {icon}
                {label}
            </p>
            <p className="font-medium">{value || "-"}</p>
        </div>
    );
}

function formatDate(dateString: string): string {
    if (!dateString) return "-";
    return new Date(dateString).toLocaleString("id-ID", {
        dateStyle: "medium",
        timeStyle: "short",
    });
}