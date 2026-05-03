// src/components/settings/ModelsTab.tsx
import { useState, useEffect } from 'react';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Switch } from '@/components/ui/switch';
import { Textarea } from '@/components/ui/textarea';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { toast } from 'sonner';
import { type AIModel, type APIProvider } from '@/services/modelProvider/types';
import {createModel, updateModel, deleteModel} from "@/services/model";
import { getProviders } from '@/services/modelProvider/modelQueries';
import { Loader2, Plus, Edit2, Trash2, Star,  CheckCircle, XCircle,  } from 'lucide-react';
import { getModelsWithStatus } from '@/services/model';

// Interface untuk model dengan status API key
interface ModelWithStatus extends AIModel {
    hasApiKey: boolean;
    providerName?: string;
    providerDisplayName?: string;
}

export function ModelsTab() {
    const [models, setModels] = useState<ModelWithStatus[]>([]);
    const [providers, setProviders] = useState<APIProvider[]>([]);
    const [loading, setLoading] = useState(true);
    const [showAddDialog, setShowAddDialog] = useState(false);
    const [editingModel, setEditingModel] = useState<AIModel | null>(null);

    // Form states
    const [formName, setFormName] = useState('');
    const [formProviderId, setFormProviderId] = useState('');
    const [formDisplayName, setFormDisplayName] = useState('');
    const [formDescription, setFormDescription] = useState('');
    const [formMaxTokens, setFormMaxTokens] = useState(4096);
    const [formTemperature, setFormTemperature] = useState(0.7);
    const [formIsDefault, setFormIsDefault] = useState(false);
    const [saving, setSaving] = useState(false);

    useEffect(() => {
        loadModels();
        loadProviders();
    }, []);

    const loadModels = async () => {
        setLoading(true);
        try {
            const data = await getModelsWithStatus();
            setModels(data as any);
        } catch (error) {
            toast.error('Gagal memuat models');
        } finally {
            setLoading(false);
        }
    };

    const loadProviders = async () => {
        try {
            const data = await getProviders();
            setProviders(data);
        } catch (error) {
            console.error('Failed to load providers:', error);
        }
    };

    const handleSave = async () => {
        if (!formName || !formProviderId || !formDisplayName) {
            toast.error('Isi semua field yang diperlukan');
            return;
        }

        setSaving(true);
        try {
            if (editingModel) {
                await updateModel(editingModel.id, {
                    name: formName,
                    providerId: formProviderId,
                    displayName: formDisplayName,
                    description: formDescription,
                    maxTokens: formMaxTokens,
                    temperature: formTemperature,
                    isDefault: formIsDefault,
                });
                toast.success('Model updated');
            } else {
                await createModel({
                    name: formName,
                    providerId: formProviderId,
                    displayName: formDisplayName,
                    description: formDescription,
                    maxTokens: formMaxTokens,
                    temperature: formTemperature,
                    isDefault: formIsDefault,
                });
                toast.success('Model created');
            }
            resetForm();
            await loadModels();
        } catch (error) {
            toast.error('Failed to save model');
        } finally {
            setSaving(false);
        }
    };

    const resetForm = () => {
        setShowAddDialog(false);
        setEditingModel(null);
        setFormName('');
        setFormProviderId('');
        setFormDisplayName('');
        setFormDescription('');
        setFormMaxTokens(4096);
        setFormTemperature(0.7);
        setFormIsDefault(false);
    };

    const handleEdit = (model: AIModel) => {
        setEditingModel(model);
        setFormName(model.name);
        setFormProviderId(model.providerId);
        setFormDisplayName(model.displayName);
        setFormDescription(model.description || '');
        setFormMaxTokens(model.maxTokens);
        setFormTemperature(model.temperature);
        setFormIsDefault(model.isDefault);
        setShowAddDialog(true);
    };

    const handleDelete = async (id: string, name: string) => {
        if (confirm(`Hapus model "${name}"?`)) {
            try {
                await deleteModel(id);
                toast.success('Model deleted');
                await loadModels();
            } catch (error) {
                toast.error('Failed to delete model');
            }
        }
    };

    const handleSetDefault = async (id: string) => {
        try {
            await updateModel(id, { isDefault: true });
            toast.success('Default model updated');
            await loadModels();
        } catch (error) {
            toast.error('Failed to set default model');
        }
    };

    const getProviderName = (providerId: string) => {
        const provider = providers.find(p => p.id === providerId);
        return provider?.displayName || provider?.name || providerId;
    };

    const getProviderBadge = (providerId: string) => {
        const provider = providers.find(p => p.id === providerId);
        const colors: Record<string, string> = {
            google: 'bg-blue-100 text-blue-700',
            openai: 'bg-green-100 text-green-700',
            anthropic: 'bg-purple-100 text-purple-700',
            openrouter: 'bg-orange-100 text-orange-700',
        };
        return colors[provider?.name || ''] || 'bg-gray-100 text-gray-700';
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

    return (
        <div className="space-y-4">
            <div className="flex justify-between items-center">
                <div>
                    <h2 className="text-lg font-semibold">AI Models</h2>
                    <p className="text-sm text-slate-500">Kelola model AI yang tersedia untuk generate konten</p>
                </div>
                <Button onClick={() => setShowAddDialog(true)} size="sm">
                    <Plus className="mr-2 h-4 w-4" />
                    Add Model
                </Button>
            </div>

            <div className="grid gap-3">
                {models.length === 0 ? (
                    <Card>
                        <CardContent className="py-8 text-center text-slate-500">
                            Belum ada model. Klik "Add Model" untuk menambahkan.
                        </CardContent>
                    </Card>
                ) : (
                    models.map((model) => (
                        <Card key={model.id}>
                            <CardContent className="p-4">
                                <div className="flex items-start justify-between">
                                    <div className="flex-1">
                                        <div className="flex items-center gap-2 flex-wrap">
                                            <h3 className="font-medium">{model.displayName}</h3>
                                            <span className={`text-xs px-2 py-0.5 rounded-full ${getProviderBadge(model.providerId)}`}>
                                                {getProviderName(model.providerId)}
                                            </span>
                                            {model.isDefault && (
                                                <span className="text-xs text-blue-500 flex items-center gap-1">
                                                    <Star className="h-3 w-3 fill-blue-500" /> Default
                                                </span>
                                            )}
                                            {!model.isActive && (
                                                <span className="text-xs text-red-500">(Inactive)</span>
                                            )}
                                            {/* API Key Status Badge */}
                                            {model.hasApiKey ? (
                                                <span className="text-xs bg-green-100 text-green-700 px-2 py-0.5 rounded-full flex items-center gap-1">
                                                    <CheckCircle className="h-3 w-3" /> API Key Active
                                                </span>
                                            ) : (
                                                <span className="text-xs bg-red-100 text-red-700 px-2 py-0.5 rounded-full flex items-center gap-1">
                                                    <XCircle className="h-3 w-3" /> API Key Missing
                                                </span>
                                            )}
                                        </div>
                                        <p className="text-sm text-slate-500 mt-1">{model.description || 'Tidak ada deskripsi'}</p>
                                        <div className="flex flex-wrap gap-3 mt-2 text-xs text-slate-400">
                                            <span>Model ID: {model.name}</span>
                                            <span>Max tokens: {model.maxTokens}</span>
                                            <span>Temperature: {model.temperature}</span>
                                        </div>
                                    </div>
                                    <div className="flex gap-1">
                                        {!model.isDefault && (
                                            <Button
                                                variant="ghost"
                                                size="sm"
                                                onClick={() => handleSetDefault(model.id)}
                                                title="Set as default"
                                            >
                                                <Star className="h-4 w-4" />
                                            </Button>
                                        )}
                                        <Button variant="ghost" size="sm" onClick={() => handleEdit(model)}>
                                            <Edit2 className="h-4 w-4" />
                                        </Button>
                                        <Button variant="ghost" size="sm" onClick={() => handleDelete(model.id, model.displayName)}>
                                            <Trash2 className="h-4 w-4 text-red-500" />
                                        </Button>
                                    </div>
                                </div>
                            </CardContent>
                        </Card>
                    ))
                )}
            </div>

            {/* Add/Edit Dialog */}
            <Dialog open={showAddDialog} onOpenChange={setShowAddDialog}>
                <DialogContent className="sm:max-w-3xl">
                    <DialogHeader>
                        <DialogTitle>{editingModel ? 'Edit Model' : 'Add New Model'}</DialogTitle>
                        <DialogDescription>
                            Konfigurasi model AI untuk generate konten
                        </DialogDescription>
                    </DialogHeader>
                    <div className="space-y-4 py-4">
                        <div>
                            <Label>Provider *</Label>
                            <Select value={formProviderId} onValueChange={setFormProviderId}>
                                <SelectTrigger>
                                    <SelectValue placeholder="Pilih Provider" />
                                </SelectTrigger>
                                <SelectContent>
                                    {providers.length === 0 ? (
                                        <SelectItem value="" disabled>No providers available. Add provider first.</SelectItem>
                                    ) : (
                                        providers.map((p) => (
                                            <SelectItem key={p.id} value={p.id}>
                                                {p.displayName} ({p.name})
                                            </SelectItem>
                                        ))
                                    )}
                                </SelectContent>
                            </Select>
                            {providers.length === 0 && (
                                <p className="mt-1 text-xs text-amber-500">
                                    Belum ada provider. Tambahkan provider dulu di tab "Providers".
                                </p>
                            )}
                        </div>
                        <div>
                            <Label>Model Name *</Label>
                            <Input
                                placeholder="gemini-2.5-flash, gpt-4-turbo, claude-3-opus"
                                value={formName}
                                onChange={(e) => setFormName(e.target.value)}
                            />
                            <p className="mt-1 text-xs text-slate-500">
                                Nama model yang digunakan dalam API call
                            </p>
                        </div>
                        <div>
                            <Label>Display Name *</Label>
                            <Input
                                placeholder="Gemini 2.5 Flash"
                                value={formDisplayName}
                                onChange={(e) => setFormDisplayName(e.target.value)}
                            />
                        </div>
                        <div>
                            <Label>Description</Label>
                            <Textarea
                                placeholder="Deskripsi model ini..."
                                value={formDescription}
                                onChange={(e) => setFormDescription(e.target.value)}
                                rows={2}
                            />
                        </div>
                        <div className="grid grid-cols-2 gap-3">
                            <div>
                                <Label>Max Tokens</Label>
                                <Input
                                    type="number"
                                    value={formMaxTokens}
                                    onChange={(e) => setFormMaxTokens(parseInt(e.target.value))}
                                />
                                <p className="mt-1 text-xs text-slate-500">Maksimum token yang dihasilkan</p>
                            </div>
                            <div>
                                <Label>Temperature</Label>
                                <Input
                                    type="number"
                                    step="0.1"
                                    min="0"
                                    max="2"
                                    value={formTemperature}
                                    onChange={(e) => setFormTemperature(parseFloat(e.target.value))}
                                />
                                <p className="mt-1 text-xs text-slate-500">Kreativitas (0 = konsisten, 1 = kreatif)</p>
                            </div>
                        </div>
                        <div className="flex items-center justify-between">
                            <div>
                                <Label>Set as Default</Label>
                                <p className="text-xs text-slate-500">Model ini akan dipilih secara otomatis</p>
                            </div>
                            <Switch checked={formIsDefault} onCheckedChange={setFormIsDefault} />
                        </div>
                    </div>
                    <DialogFooter>
                        <Button variant="outline" onClick={resetForm}>Cancel</Button>
                        <Button onClick={handleSave} disabled={saving || providers.length === 0}>
                            {saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                            {editingModel ? 'Update' : 'Create'}
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </div>
    );
}