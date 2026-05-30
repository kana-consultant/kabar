import { useState, useEffect } from "react";
import { Button, Input, Label, Textarea, Switch, Select, SelectContent, SelectItem, SelectTrigger, SelectValue, Table, TableBody, TableCell, TableHead, TableHeader, TableRow, Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@kana-consultant/ui-kit";
import { useToast } from "@/hooks/use-toast";
import { Loader2, Eye, EyeOff, Trash2, Plus, Key, CheckCircle, XCircle, Edit } from "lucide-react";
import { getAPIKeys, createAPIKey, updateAPIKey, deleteAPIKey, type APIKey } from "@/services/apiKey";
import { getModels, type AIModel } from "@/services/model";
import { type APIProvider } from "@/services/modelProvider/types";
import { getProviders } from "@/services/modelProvider/modelQueries";
import { cn } from "@/lib/utils";
import { Can } from "@/components/ui/Can";

export function ApiKeysTab() {
    const toast = useToast();

    // ── state ──────────────────────────────────────────────────
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [apiKeys, setApiKeys] = useState<APIKey[]>([]);
    const [models, setModels] = useState<AIModel[]>([]);
    const [providers, setProviders] = useState<APIProvider[]>([]);
    const [showDialog, setShowDialog] = useState(false);
    const [editingKey, setEditingKey] = useState<APIKey | null>(null);
    const [dialogService, setDialogService] = useState<'text' | 'image'>('text');
    const [showApiKey, setShowApiKey] = useState(false);
    const [formProviderId, setFormProviderId] = useState("");
    const [formModelId, setFormModelId] = useState("");
    const [formApiKey, setFormApiKey] = useState("");
    const [formSystemPrompt, setFormSystemPrompt] = useState("");
    const [formIsActive, setFormIsActive] = useState(true);
    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
    const [deleting, setDeleting] = useState(false);
    const [selectedDeleteKey, setSelectedDeleteKey] = useState<APIKey | null>(null);

    const loadData = async () => {
        setLoading(true);
        try {
            const [keysData, modelsData, providersData] = await Promise.all([
                getAPIKeys(), getModels(), getProviders()
            ]);
            setApiKeys(keysData || []);
            setModels(modelsData || []);
            setProviders(providersData || []);
        } catch {
            toast.error('Gagal memuat data');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => { loadData(); }, []);

    const openAddDialog = (service: 'text' | 'image') => {
        setDialogService(service); setEditingKey(null);
        setFormProviderId(""); setFormModelId("");
        setFormApiKey(""); setFormSystemPrompt("");
        setFormIsActive(true); setShowApiKey(false);
        setShowDialog(true);
    };

    const openEditDialog = (key: APIKey) => {
        setEditingKey(key); setDialogService(key.service as 'text' | 'image');
        setFormProviderId(key.providerId || ""); setFormModelId(key.modelId || "");
        setFormApiKey(""); setFormSystemPrompt(key.systemPrompt || "");
        setFormIsActive(key.isActive); setShowApiKey(false);
        setShowDialog(true);
    };

    const handleSave = async () => {
        if (!formProviderId) {
            toast.error('Pilih provider terlebih dahulu');
            return;
        }
        if (!formModelId) {
            toast.error('Pilih model terlebih dahulu');
            return;
        }
        if (!formApiKey && !editingKey) {
            toast.error('Masukkan API Key');
            return;
        }

        setSaving(true);
        try {
            if (editingKey) {
                const updateData: any = {
                    providerId: formProviderId,
                    modelId: formModelId,
                    systemPrompt: formSystemPrompt || undefined,
                    isActive: formIsActive
                };
                if (formApiKey) updateData.key = formApiKey;
                await updateAPIKey(editingKey.id, updateData);
                toast.success('API Key updated');
            } else {
                await createAPIKey({
                    service: dialogService,
                    providerId: formProviderId,
                    modelId: formModelId,
                    key: formApiKey,
                    systemPrompt: formSystemPrompt
                });
                toast.success('API Key created');
            }
            setShowDialog(false);
            await loadData();
        } catch {
            toast.error('Failed to save API Key');
        } finally {
            setSaving(false);
        }
    };

    const handleDelete = (key: APIKey) => {
        setSelectedDeleteKey(key);
        setDeleteDialogOpen(true);
    };

    const confirmDelete = async () => {
        if (!selectedDeleteKey) return;
        setDeleting(true);
        try {
            await deleteAPIKey(selectedDeleteKey.id);
            toast.success("API Key deleted");
            setDeleteDialogOpen(false);
            setSelectedDeleteKey(null);
            await loadData();
        } catch {
            toast.error("Failed to delete API Key");
        } finally {
            setDeleting(false);
        }
    };

    const getProviderName = (id: string) => providers.find(p => p.id === id)?.displayName || id;
    const getModelName = (id: string) => models.find(m => m.id === id)?.displayName || id;

    // ── helpers ────────────────────────────────────────────────
    const textKeys = apiKeys.filter(k => k.service === 'text');
    const imageKeys = apiKeys.filter(k => k.service === 'image');

    const inputCls = cn(
        "h-8 text-sm border-slate-200/80 bg-white placeholder:text-slate-400",
        "dark:border-white/[0.08] dark:bg-white/[0.03] dark:placeholder:text-slate-600",
        "focus-visible:ring-1 focus-visible:ring-green-500/40 focus-visible:border-green-400/60",
        "dark:focus-visible:ring-purple-500/40 dark:focus-visible:border-purple-400/40"
    );

    if (loading) return (
        <div className="flex items-center justify-center py-20">
            <Loader2 className="h-6 w-6 animate-spin text-slate-400" />
        </div>
    );

    return (
        <div className="space-y-4">
            {/* Main card */}
            <div className={cn(
                "overflow-hidden rounded-2xl border",
                "bg-white border-slate-200/80",
                "dark:bg-[#0f0d1a] dark:border-white/[0.06]"
            )}>
                {/* Header */}
                <div className={cn(
                    "flex items-center justify-between px-6 py-4 border-b",
                    "border-slate-100 bg-slate-50/60",
                    "dark:border-white/[0.05] dark:bg-white/[0.02]"
                )}>
                    <div>
                        <p className="text-sm font-medium text-slate-800 dark:text-slate-100">
                            API Keys Management
                        </p>
                        <p className="text-xs text-slate-400 dark:text-slate-600 mt-0.5">
                            Kelola API Key untuk berbagai provider AI
                        </p>
                    </div>
                    <div className="flex gap-2">
                        <Can permission="api_keys:create:team">
                            <Button
                                size="sm"
                                className="h-8 gap-1.5 px-3 text-xs bg-green-600 hover:bg-green-700 text-white dark:bg-purple-600 dark:hover:bg-purple-700"
                                onClick={() => openAddDialog('text')}
                            >
                                <Plus className="h-3.5 w-3.5" /> Add Text Key
                            </Button>
                        </Can>
                        <Can permission="api_keys:create:team">
                            <Button
                                variant="outline" size="sm"
                                className="h-8 gap-1.5 px-3 text-xs border-slate-200/80 dark:border-white/[0.08]"
                                onClick={() => openAddDialog('image')}
                            >
                                <Plus className="h-3.5 w-3.5" /> Add Image Key
                            </Button>
                        </Can>
                    </div>
                </div>

                <div className="p-5 space-y-5">
                    {/* Mini stats */}
                    <div className="grid grid-cols-2 gap-3">
                        {[
                            { label: "Text Generation", count: textKeys.length, active: textKeys.filter(k => k.isActive).length, iconCls: "bg-blue-50 text-blue-600 ring-blue-200/60 dark:bg-blue-500/10 dark:text-blue-400 dark:ring-blue-500/20" },
                            { label: "Image Generation", count: imageKeys.length, active: imageKeys.filter(k => k.isActive).length, iconCls: "bg-violet-50 text-violet-600 ring-violet-200/60 dark:bg-violet-500/10 dark:text-violet-400 dark:ring-violet-500/20" },
                        ].map(({ label, count, active, iconCls }) => (
                            <div key={label} className={cn(
                                "flex items-center justify-between rounded-xl border p-4",
                                "bg-slate-50/60 border-slate-200/60",
                                "dark:bg-white/[0.02] dark:border-white/[0.05]"
                            )}>
                                <div>
                                    <p className="text-xs text-slate-400 dark:text-slate-600">{label}</p>
                                    <p className="text-2xl font-semibold tracking-tight text-slate-900 dark:text-white tabular-nums mt-1">
                                        {count}
                                    </p>
                                    <p className="text-[11px] text-slate-400 dark:text-slate-600 mt-1.5 flex items-center gap-1">
                                        <span className="h-1.5 w-1.5 rounded-full bg-green-500 inline-block" />
                                        {active} aktif
                                    </p>
                                </div>
                                <div className={cn("flex h-9 w-9 items-center justify-center rounded-xl ring-1", iconCls)}>
                                    <Key className="h-4 w-4" />
                                </div>
                            </div>
                        ))}
                    </div>

                    {/* Table */}
                    {apiKeys.length === 0 ? (
                        <div className="flex flex-col items-center py-12">
                            <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-slate-100 text-slate-400 dark:bg-white/[0.04] dark:text-slate-600">
                                <Key className="h-5 w-5" />
                            </div>
                            <p className="mt-3 text-sm font-medium text-slate-600 dark:text-slate-400">
                                Belum ada API Key
                            </p>
                            <p className="mt-1 text-xs text-slate-400 dark:text-slate-600">
                                Klik "Add Text Key" atau "Add Image Key" untuk menambahkan
                            </p>
                        </div>
                    ) : (
                        <div className="overflow-hidden rounded-xl border border-slate-200/80 dark:border-white/[0.06]">
                            <Table>
                                <TableHeader>
                                    <TableRow className="bg-slate-50/80 dark:bg-white/[0.02] hover:bg-slate-50/80">
                                        {["Service", "Provider", "Model", "Status", "System Prompt", ""].map(h => (
                                            <TableHead key={h} className={cn(
                                                "text-[10px] font-medium uppercase tracking-wide text-slate-400 dark:text-slate-600 py-2.5",
                                                h === "" ? "text-right" : ""
                                            )}>
                                                {h}
                                            </TableHead>
                                        ))}
                                    </TableRow>
                                </TableHeader>
                                <TableBody>
                                    {apiKeys.map((key) => (
                                        <TableRow key={key.id} className="border-slate-100 dark:border-white/[0.04] hover:bg-slate-50/50 dark:hover:bg-white/[0.01]">
                                            <TableCell>
                                                <span className={cn(
                                                    "inline-flex items-center gap-1 rounded-md border px-1.5 py-0.5 text-[10px] font-medium",
                                                    key.service === 'text'
                                                        ? "bg-blue-50 text-blue-700 border-blue-200/60 dark:bg-blue-500/10 dark:text-blue-400 dark:border-blue-500/20"
                                                        : "bg-violet-50 text-violet-700 border-violet-200/60 dark:bg-violet-500/10 dark:text-violet-400 dark:border-violet-500/20"
                                                )}>
                                                    {key.service === 'text' ? 'Text Gen' : 'Image Gen'}
                                                </span>
                                            </TableCell>
                                            <TableCell className="text-sm font-medium text-slate-800 dark:text-slate-200">
                                                {getProviderName(key.providerId)}
                                            </TableCell>
                                            <TableCell className="text-xs text-slate-500 dark:text-slate-500">
                                                {getModelName(key.modelId)}
                                            </TableCell>
                                            <TableCell>
                                                <span className={cn(
                                                    "inline-flex items-center gap-1 rounded-md border px-1.5 py-0.5 text-[10px] font-medium",
                                                    key.isActive
                                                        ? "bg-green-50 text-green-700 border-green-200/60 dark:bg-green-500/10 dark:text-green-400 dark:border-green-500/20"
                                                        : "bg-red-50 text-red-700 border-red-200/60 dark:bg-red-500/10 dark:text-red-400 dark:border-red-500/20"
                                                )}>
                                                    {key.isActive
                                                        ? <><CheckCircle className="h-2.5 w-2.5" /> Active</>
                                                        : <><XCircle className="h-2.5 w-2.5" /> Inactive</>
                                                    }
                                                </span>
                                            </TableCell>
                                            <TableCell className="max-w-[180px]">
                                                {key.systemPrompt
                                                    ? <span className="text-xs text-slate-400 truncate block" title={key.systemPrompt}>{key.systemPrompt.substring(0, 50)}…</span>
                                                    : <span className="text-xs text-slate-300 dark:text-slate-700">—</span>
                                                }
                                            </TableCell>
                                            <TableCell className="text-right">
                                                <div className="flex justify-end gap-1">
                                                    <Can permission="api_key:edit:team">
                                                        <Button variant="ghost" size="icon"
                                                            className="h-7 w-7 rounded-lg text-slate-400 hover:text-green-600 hover:bg-green-50 dark:hover:text-purple-400 dark:hover:bg-purple-500/10"
                                                            onClick={() => openEditDialog(key)}
                                                        >
                                                            <Edit className="h-3.5 w-3.5" />
                                                        </Button>
                                                    </Can>
                                                    <Can permission="api_keys:delete:team">
                                                        <Button variant="ghost" size="icon"
                                                            className="h-7 w-7 rounded-lg text-slate-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-500/10"
                                                            onClick={() => handleDelete(key)}
                                                        >
                                                            <Trash2 className="h-3.5 w-3.5" />
                                                        </Button>
                                                    </Can>
                                                </div>
                                            </TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </div>
                    )}
                </div>
            </div>

            {/* Add/Edit Dialog */}
            <Dialog open={showDialog} onOpenChange={setShowDialog}>
                <DialogContent className="sm:max-w-2xl">
                    <DialogHeader>
                        <DialogTitle className="text-base">
                            {editingKey ? 'Edit API Key' : `Add ${dialogService === 'text' ? 'Text Generation' : 'Image Generation'} API Key`}
                        </DialogTitle>
                        <DialogDescription className="text-xs">
                            Konfigurasi API key untuk {dialogService === 'text' ? 'generate artikel' : 'generate gambar'}
                        </DialogDescription>
                    </DialogHeader>
                    <div className="space-y-4 py-2">
                        <div className="space-y-1.5">
                            <Label className="text-xs font-medium uppercase tracking-wide text-slate-500">Provider *</Label>
                            <Select value={formProviderId} onValueChange={setFormProviderId}>
                                <SelectTrigger className={inputCls}><SelectValue placeholder="Pilih Provider" /></SelectTrigger>
                                <SelectContent>
                                    {providers.map(p => <SelectItem key={p.id} value={p.id}>{p.displayName} ({p.name})</SelectItem>)}
                                </SelectContent>
                            </Select>
                        </div>

                        {formProviderId && (
                            <div className="space-y-1.5">
                                <Label className="text-xs font-medium uppercase tracking-wide text-slate-500">
                                    Model *
                                </Label>

                                <Select
                                    value={formModelId}
                                    onValueChange={setFormModelId}
                                    required
                                >
                                    <SelectTrigger
                                        className={`${inputCls} ${!formModelId ? "border-red-500" : ""
                                            }`}
                                    >
                                        <SelectValue placeholder="Pilih Model" />
                                    </SelectTrigger>

                                    <SelectContent>
                                        {models
                                            .filter((m) => m.providerId === formProviderId)
                                            .map((m) => (
                                                <SelectItem key={m.id} value={m.id}>
                                                    {m.displayName} ({m.name})
                                                </SelectItem>
                                            ))}
                                    </SelectContent>
                                </Select>

                                {!formModelId && (
                                    <p className="text-xs text-red-500">
                                        Model wajib dipilih
                                    </p>
                                )}
                            </div>
                        )}

                        <div className="space-y-1.5">
                            <Label className="text-xs font-medium uppercase tracking-wide text-slate-500">API Key *</Label>
                            <div className="relative">
                                <Input
                                    type={showApiKey ? "text" : "password"}
                                    placeholder="Masukkan API Key"
                                    value={formApiKey}
                                    onChange={(e: any) => setFormApiKey(e.target.value)}
                                    className={cn(inputCls, "pr-9")}
                                />
                                <button
                                    type="button"
                                    onClick={() => setShowApiKey(!showApiKey)}
                                    className="absolute right-2.5 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600"
                                >
                                    {showApiKey ? <EyeOff className="h-3.5 w-3.5" /> : <Eye className="h-3.5 w-3.5" />}
                                </button>
                            </div>
                            {editingKey && (
                                <p className="text-xs text-green-600 dark:text-green-400 flex items-center gap-1">
                                    <Key className="h-3 w-3" />
                                    Kosongkan jika tidak ingin mengubah key yang sudah ada
                                </p>
                            )}
                        </div>

                        <div className="space-y-1.5">
                            <Label className="text-xs font-medium uppercase tracking-wide text-slate-500">System Prompt (Opsional)</Label>
                            <Textarea
                                placeholder="Custom system prompt untuk AI..."
                                value={formSystemPrompt}
                                onChange={(e: any) => setFormSystemPrompt(e.target.value)}
                                rows={3}
                                className={cn(
                                    inputCls,
                                    "h-auto py-2 resize-none",
                                    "max-h-[150px] overflow-y-auto"
                                )}
                            />
                        </div>

                        <div className={cn(
                            "flex items-center justify-between rounded-xl border p-4",
                            "bg-slate-50/60 border-slate-200/60",
                            "dark:bg-white/[0.02] dark:border-white/[0.05]"
                        )}>
                            <div>
                                <p className="text-sm font-medium text-slate-800 dark:text-slate-200">Active</p>
                                <p className="text-xs text-slate-400 dark:text-slate-600 mt-0.5">
                                    Nonaktifkan jika tidak ingin menggunakan key ini
                                </p>
                            </div>
                            <Switch
                                checked={formIsActive}
                                onCheckedChange={setFormIsActive}
                                className="data-[state=checked]:bg-green-600 dark:data-[state=checked]:bg-purple-600"
                            />
                        </div>
                    </div>
                    <DialogFooter>
                        <Button variant="outline" size="sm" className="h-8 text-xs" onClick={() => setShowDialog(false)}>
                            Cancel
                        </Button>
                        <Button size="sm" className="h-8 text-xs bg-green-600 hover:bg-green-700 dark:bg-purple-600 dark:hover:bg-purple-700" onClick={handleSave} disabled={saving}>
                            {saving && <Loader2 className="mr-1.5 h-3.5 w-3.5 animate-spin" />}
                            {editingKey ? 'Update' : 'Create'}
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>

            {/* Delete Dialog */}
            <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
                <DialogContent className="sm:max-w-sm">
                    <DialogHeader>
                        <DialogTitle className="flex items-center gap-2 text-sm text-red-600">
                            <Trash2 className="h-4 w-4" /> Hapus API Key
                        </DialogTitle>
                        <DialogDescription className="text-xs">
                            Tindakan ini tidak bisa dibatalkan.
                        </DialogDescription>
                    </DialogHeader>
                    <div className={cn(
                        "rounded-xl border p-4 space-y-3",
                        "bg-slate-50/60 border-slate-200/60",
                        "dark:bg-white/[0.02] dark:border-white/[0.05]"
                    )}>
                        {[
                            { label: "Provider", value: selectedDeleteKey && getProviderName(selectedDeleteKey.providerId) },
                            { label: "Model", value: selectedDeleteKey && getModelName(selectedDeleteKey.modelId) },
                        ].map(({ label, value }) => (
                            <div key={label}>
                                <p className="text-[10px] uppercase tracking-wide text-slate-400">{label}</p>
                                <p className="text-sm font-medium text-slate-800 dark:text-slate-200 mt-0.5">{value}</p>
                            </div>
                        ))}
                        {selectedDeleteKey && (
                            <div>
                                <p className="text-[10px] uppercase tracking-wide text-slate-400">Service</p>
                                <span className={cn(
                                    "mt-1 inline-flex items-center gap-1 rounded-md border px-1.5 py-0.5 text-[10px] font-medium",
                                    selectedDeleteKey.service === 'text'
                                        ? "bg-blue-50 text-blue-700 border-blue-200/60"
                                        : "bg-violet-50 text-violet-700 border-violet-200/60"
                                )}>
                                    {selectedDeleteKey.service === 'text' ? 'Text Gen' : 'Image Gen'}
                                </span>
                            </div>
                        )}
                    </div>
                    <DialogFooter>
                        <Button variant="outline" size="sm" className="h-8 text-xs" onClick={() => setDeleteDialogOpen(false)} disabled={deleting}>
                            Cancel
                        </Button>
                        <Button size="sm" className="h-8 text-xs bg-red-600 hover:bg-red-700 text-white" onClick={confirmDelete} disabled={deleting}>
                            {deleting && <Loader2 className="mr-1.5 h-3.5 w-3.5 animate-spin" />}
                            Hapus
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </div>
    );
}