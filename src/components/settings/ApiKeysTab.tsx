// src/components/settings/ApiKeysTab.tsx
import { useState, useEffect } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Switch } from "@/components/ui/switch";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";
import { Loader2, Eye, EyeOff, Trash2, Plus, Key, CheckCircle, XCircle, Edit } from "lucide-react";
import { getAPIKeys, createAPIKey, updateAPIKey, deleteAPIKey, type APIKey } from "@/services/apiKey";
import { getModels, type AIModel } from "@/services/model";
import { type APIProvider } from "@/services/modelProvider/types";
import { getProviders } from "@/services/modelProvider/modelQueries";

export function ApiKeysTab() {
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [apiKeys, setApiKeys] = useState<APIKey[]>([]);
    const [models, setModels] = useState<AIModel[]>([]);
    const [providers, setProviders] = useState<APIProvider[]>([]);

    // Dialog states
    const [showDialog, setShowDialog] = useState(false);
    const [editingKey, setEditingKey] = useState<APIKey | null>(null);
    const [dialogService, setDialogService] = useState<'text' | 'image'>('text');
    const [showApiKey, setShowApiKey] = useState(false);

    // Form states
    const [formProviderId, setFormProviderId] = useState("");
    const [formModelId, setFormModelId] = useState("");
    const [formApiKey, setFormApiKey] = useState("");
    const [formSystemPrompt, setFormSystemPrompt] = useState("");
    const [formIsActive, setFormIsActive] = useState(true);

    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
    const [deleting, setDeleting] = useState(false);
    const [selectedDeleteKey, setSelectedDeleteKey] =
        useState<APIKey | null>(null);

    useEffect(() => {
        loadData();
    }, []);

    const loadData = async () => {
        setLoading(true);
        try {
            const [keysData, modelsData, providersData] = await Promise.all([
                getAPIKeys(),
                getModels(),
                getProviders()
            ]);
            setApiKeys(keysData || []);
            setModels(modelsData || []);
            setProviders(providersData || []);
        } catch (error) {
            console.error('Failed to load data:', error);
            toast.error('Gagal memuat data');
        } finally {
            setLoading(false);
        }
    };

    const openAddDialog = (service: 'text' | 'image') => {
        setDialogService(service);
        setEditingKey(null);
        setFormProviderId("");
        setFormModelId("");
        setFormApiKey("");
        setFormSystemPrompt("");
        setFormIsActive(true);
        setShowApiKey(false);
        setShowDialog(true);
    };

    const openEditDialog = (key: APIKey) => {
        setEditingKey(key);
        setDialogService(key.service as 'text' | 'image');
        setFormProviderId(key.providerId || "");
        setFormModelId(key.modelId || "");
        setFormApiKey("");
        setFormSystemPrompt(key.systemPrompt || "");
        setFormIsActive(key.isActive);
        setShowApiKey(false);
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
                    isActive: formIsActive,
                };
                if (formApiKey) {
                    updateData.key = formApiKey;
                }
                await updateAPIKey(editingKey.id, updateData);
                toast.success('API Key updated');
            } else {
                await createAPIKey({
                    service: dialogService,
                    providerId: formProviderId,
                    modelId: formModelId,
                    key: formApiKey,
                    systemPrompt: formSystemPrompt,
                });
                toast.success('API Key created');
            }
            setShowDialog(false);
            await loadData();
        } catch (error) {
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

            await deleteAPIKey(
                selectedDeleteKey.id
            );

            toast.success(
                "API Key deleted"
            );

            setDeleteDialogOpen(false);
            setSelectedDeleteKey(null);

            await loadData();

        } catch (error) {

            toast.error(
                "Failed to delete API Key"
            );

        } finally {

            setDeleting(false);
        }
    };

    const getProviderName = (providerId: string) => {
        const provider = providers.find(p => p.id === providerId);
        return provider?.displayName || provider?.name || providerId;
    };

    const getModelDisplayName = (modelId: string) => {
        const model = models.find(m => m.id === modelId);
        return model?.displayName || model?.name || modelId;
    };

    const getStatusBadge = (isActive: boolean) => {
        if (isActive) {
            return <Badge className="bg-green-100 text-green-700 hover:bg-green-100"><CheckCircle className="h-3 w-3 mr-1" /> Active</Badge>;
        }
        return <Badge variant="outline" className="text-red-500"><XCircle className="h-3 w-3 mr-1" /> Inactive</Badge>;
    };

    const getServiceBadge = (service: string) => {
        if (service === 'text') {
            return <Badge className="bg-blue-100 text-blue-700">Text Generation</Badge>;
        }
        return <Badge className="bg-purple-100 text-purple-700">Image Generation</Badge>;
    };

    if (loading) {
        return (
            <Card>
                <CardContent className="flex justify-center py-8">
                    <Loader2 className="h-8 w-8 animate-spin" />
                </CardContent>
            </Card>
        );
    }

    const textKeys = apiKeys.filter(k => k.service === 'text');
    const imageKeys = apiKeys.filter(k => k.service === 'image');

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex justify-between items-center">
                <div>
                    <h2 className="text-lg font-semibold">API Keys Management</h2>
                    <p className="text-sm text-slate-500">Kelola API Key untuk berbagai provider AI</p>
                </div>
                <div className="flex gap-2">
                    <Button onClick={() => openAddDialog('text')} size="sm">
                        <Plus className="mr-2 h-4 w-4" />
                        Add Text Key
                    </Button>
                    <Button onClick={() => openAddDialog('image')} size="sm" variant="outline">
                        <Plus className="mr-2 h-4 w-4" />
                        Add Image Key
                    </Button>
                </div>
            </div>

            {/* Info Cards */}
            <div className="grid grid-cols-2 gap-4">
                <Card>
                    <CardContent className="pt-4">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-sm text-slate-500">Text Generation</p>
                                <p className="text-2xl font-bold">{textKeys.length}</p>
                            </div>
                            <div className="h-10 w-10 rounded-full bg-blue-100 flex items-center justify-center">
                                <Key className="h-5 w-5 text-blue-600" />
                            </div>
                        </div>
                        <p className="text-xs text-slate-400 mt-2">
                            {textKeys.filter(k => k.isActive).length} aktif
                        </p>
                    </CardContent>
                </Card>
                <Card>
                    <CardContent className="pt-4">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-sm text-slate-500">Image Generation</p>
                                <p className="text-2xl font-bold">{imageKeys.length}</p>
                            </div>
                            <div className="h-10 w-10 rounded-full bg-purple-100 flex items-center justify-center">
                                <Key className="h-5 w-5 text-purple-600" />
                            </div>
                        </div>
                        <p className="text-xs text-slate-400 mt-2">
                            {imageKeys.filter(k => k.isActive).length} aktif
                        </p>
                    </CardContent>
                </Card>
            </div>

            {/* API Keys Table */}
            <Card>
                <CardHeader>
                    <CardTitle>Registered API Keys</CardTitle>
                    <CardDescription>Semua API key yang telah didaftarkan</CardDescription>
                </CardHeader>
                <CardContent>
                    {apiKeys.length === 0 ? (
                        <div className="text-center py-8 text-slate-500">
                            <Key className="h-12 w-12 mx-auto mb-3 text-slate-300" />
                            <p>Belum ada API Key yang didaftarkan</p>
                            <p className="text-sm">Klik tombol "Add Text Key" atau "Add Image Key" untuk menambahkan</p>
                        </div>
                    ) : (
                        <Table>
                            <TableHeader>
                                <TableRow>
                                    <TableHead>Service</TableHead>
                                    <TableHead>Provider</TableHead>
                                    <TableHead>Model</TableHead>
                                    <TableHead>Status</TableHead>
                                    <TableHead>System Prompt</TableHead>
                                    <TableHead className="text-right">Actions</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {apiKeys.map((key) => (
                                    <TableRow key={key.id}>
                                        <TableCell>{getServiceBadge(key.service)}</TableCell>
                                        <TableCell className="font-medium">{getProviderName(key.providerId)}</TableCell>
                                        <TableCell>{getModelDisplayName(key.modelId)}</TableCell>
                                        <TableCell>{getStatusBadge(key.isActive)}</TableCell>
                                        <TableCell className="max-w-[200px] truncate">
                                            {key.systemPrompt ? (
                                                <span className="text-xs text-slate-500" title={key.systemPrompt}>
                                                    {key.systemPrompt.substring(0, 50)}...
                                                </span>
                                            ) : (
                                                <span className="text-xs text-slate-400">-</span>
                                            )}
                                        </TableCell>
                                        <TableCell className="text-right">
                                            <div className="flex justify-end gap-1">
                                                <Button
                                                    variant="ghost"
                                                    size="sm"
                                                    onClick={() => openEditDialog(key)}
                                                >
                                                    <Edit className="h-4 w-4" />
                                                </Button>
                                                <Button
                                                    variant="ghost"
                                                    size="sm"
                                                    onClick={() => handleDelete(key)}
                                                    className="text-red-500 hover:text-red-700"
                                                >
                                                    <Trash2 className="h-4 w-4" />
                                                </Button>
                                            </div>
                                        </TableCell>
                                    </TableRow>
                                ))}
                            </TableBody>
                        </Table>
                    )}
                </CardContent>
            </Card>

            {/* Add/Edit Dialog */}
            <Dialog open={showDialog} onOpenChange={setShowDialog} >
                <DialogContent className="sm:max-w-3xl">
                    <DialogHeader>
                        <DialogTitle>
                            {editingKey ? 'Edit API Key' : `Add ${dialogService === 'text' ? 'Text Generation' : 'Image Generation'} API Key`}
                        </DialogTitle>
                        <DialogDescription>
                            Konfigurasi API key untuk {dialogService === 'text' ? 'generate artikel' : 'generate gambar'}
                        </DialogDescription>
                    </DialogHeader>
                    <div className="space-y-4 py-4">
                        {/* Provider */}
                        <div>
                            <Label>Provider *</Label>
                            <Select value={formProviderId} onValueChange={setFormProviderId}>
                                <SelectTrigger>
                                    <SelectValue placeholder="Pilih Provider" />
                                </SelectTrigger>

                                <SelectContent>
                                    {providers.map((p) => (
                                        <SelectItem key={p.id} value={p.id}>
                                            {p.displayName} ({p.name})
                                        </SelectItem>
                                    ))}
                                </SelectContent>


                            </Select>
                        </div>

                        {/* Model */}
                        {formProviderId && (
                            <div>
                                <Label>Model *</Label>

                                <Select
                                    value={formModelId}
                                    onValueChange={setFormModelId}
                                >
                                    <SelectTrigger>
                                        <SelectValue placeholder="Choose Model" />
                                    </SelectTrigger>

                                    <SelectContent>
                                        {models
                                            .filter((m) => m.providerId === formProviderId)
                                            .map((model) => (
                                                <SelectItem
                                                    key={model.id}
                                                    value={model.id}
                                                >
                                                    {model.displayName} ({model.name})
                                                </SelectItem>
                                            ))}
                                    </SelectContent>
                                </Select>
                            </div>
                        )}

                        {/* API Key */}
                        <div>
                            <Label>API Key *</Label>
                            <div className="relative mt-1">
                                <Input
                                    type={showApiKey ? "text" : "password"}
                                    placeholder="Masukkan API Key"
                                    value={formApiKey}
                                    onChange={(e) => setFormApiKey(e.target.value)}
                                    className="pr-10"
                                />
                                <button
                                    type="button"
                                    onClick={() => setShowApiKey(!showApiKey)}
                                    className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600"
                                >
                                    {showApiKey ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                                </button>
                            </div>
                            {editingKey && (
                                <p className="mt-1 text-xs text-green-600">
                                    <Key className="inline h-3 w-3 mr-1" />
                                    Kosongkan jika tidak ingin mengubah key yang sudah ada
                                </p>
                            )}
                        </div>

                        {/* System Prompt */}
                        <div>
                            <Label>System Prompt (Optional)</Label>
                            <Textarea
                                placeholder="Custom system prompt untuk AI..."
                                value={formSystemPrompt}
                                onChange={(e) => setFormSystemPrompt(e.target.value)}
                                rows={3}
                            />
                        </div>

                        {/* Active Switch */}
                        <div className="flex items-center justify-between">
                            <div>
                                <Label>Active</Label>
                                <p className="text-xs text-slate-500">Nonaktifkan jika tidak ingin menggunakan key ini</p>
                            </div>
                            <Switch checked={formIsActive} onCheckedChange={setFormIsActive} />
                        </div>
                    </div>
                    <DialogFooter>
                        <Button variant="outline" onClick={() => setShowDialog(false)}>Cancel</Button>
                        <Button onClick={handleSave} disabled={saving}>
                            {saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                            {editingKey ? 'Update' : 'Create'}
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>

            <Dialog
                open={deleteDialogOpen}
                onOpenChange={setDeleteDialogOpen}
            >
                <DialogContent className="sm:max-w-md">

                    <DialogHeader>
                        <DialogTitle className="flex items-center gap-2 text-red-600">
                            <Trash2 className="h-5 w-5" />
                            Hapus API Key
                        </DialogTitle>

                        <DialogDescription>
                            Apakah yakin ingin menghapus API Key ini?
                            Tindakan ini tidak bisa dibatalkan.
                        </DialogDescription>

                    </DialogHeader>

                    <div className="rounded-lg border p-4  space-y-2">
                        <div>
                            <p className="text-xs text-slate-500">
                                Provider
                            </p>

                            <p className="font-medium">
                                {
                                    selectedDeleteKey &&
                                    getProviderName(
                                        selectedDeleteKey.providerId
                                    )
                                }
                            </p>
                        </div>

                        <div>
                            <p className="text-xs text-slate-500">
                                Model
                            </p>

                            <p className="font-medium">
                                {
                                    selectedDeleteKey &&
                                    getModelDisplayName(
                                        selectedDeleteKey.modelId
                                    )
                                }
                            </p>
                        </div>

                        <div>
                            <p className="text-xs text-slate-500">
                                Service
                            </p>

                            <div className="mt-1">
                                {
                                    selectedDeleteKey &&
                                    getServiceBadge(
                                        selectedDeleteKey.service
                                    )
                                }
                            </div>

                        </div>
                    </div>

                    <DialogFooter>
                        <Button
                            variant="outline"
                            onClick={() => {
                                setDeleteDialogOpen(false)
                            }}
                            disabled={deleting}
                        >
                            Cancel
                        </Button>

                        <Button
                            variant="destructive"
                            onClick={confirmDelete}
                            disabled={deleting}
                        >
                            {deleting && (
                                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                            )}

                            Hapus
                        </Button>

                    </DialogFooter>

                </DialogContent>
            </Dialog>
        </div>
    );
}