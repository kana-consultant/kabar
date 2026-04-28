// src/components/settings/ProvidersTab.tsx
import { useState, useEffect } from 'react';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Badge } from '@/components/ui/badge';
import { toast } from 'sonner';
import { getProviders, createProvider, updateProvider, deleteProvider, type APIProvider, type CreateProviderRequest } from '@/services/modelProviderService';
import { Loader2, Plus, Edit2, Trash2, Key, Globe, Code, Eye, EyeOff } from 'lucide-react';

export function ProvidersTab() {
    const [providers, setProviders] = useState<APIProvider[]>([]);
    const [loading, setLoading] = useState(true);
    const [showDialog, setShowDialog] = useState(false);
    const [editingProvider, setEditingProvider] = useState<APIProvider | null>(null);
    const [saving, setSaving] = useState(false);
    const [showJsonPreview, setShowJsonPreview] = useState(false);

    // Form states
    const [formName, setFormName] = useState('');
    const [formDisplayName, setFormDisplayName] = useState('');
    const [formDescription, setFormDescription] = useState('');
    const [formBaseUrl, setFormBaseUrl] = useState('');
    const [formAuthType, setFormAuthType] = useState('bearer');
    const [formAuthHeader, setFormAuthHeader] = useState('Authorization');
    const [formAuthPrefix, setFormAuthPrefix] = useState('Bearer');
    const [formTextEndpoint, setFormTextEndpoint] = useState('');
    const [formImageEndpoint, setFormImageEndpoint] = useState('');
    const [formDefaultHeaders, setFormDefaultHeaders] = useState('{}');
    const [formRequestTemplate, setFormRequestTemplate] = useState('{}');
    const [formResponseTextPath, setFormResponseTextPath] = useState('');
    const [formResponseImagePath, setFormResponseImagePath] = useState('');

    useEffect(() => {
        loadProviders();
    }, []);

    const loadProviders = async () => {
        setLoading(true);
        try {
            const data = await getProviders();
            setProviders(data);
        } catch (error) {
            toast.error('Gagal memuat providers');
        } finally {
            setLoading(false);
        }
    };

    const resetForm = () => {
        setEditingProvider(null);
        setFormName('');
        setFormDisplayName('');
        setFormDescription('');
        setFormBaseUrl('');
        setFormAuthType('bearer');
        setFormAuthHeader('Authorization');
        setFormAuthPrefix('Bearer');
        setFormTextEndpoint('');
        setFormImageEndpoint('');
        setFormDefaultHeaders('{}');
        setFormRequestTemplate('{}');
        setFormResponseTextPath('');
        setFormResponseImagePath('');
    };

    const handleEdit = (provider: APIProvider) => {
        setEditingProvider(provider);
        setFormName(provider.name);
        setFormDisplayName(provider.displayName);
        setFormDescription(provider.description || '');
        setFormBaseUrl(provider.baseUrl);
        setFormAuthType(provider.authType);
        setFormAuthHeader(provider.authHeader);
        setFormAuthPrefix(provider.authPrefix);
        setFormTextEndpoint(provider.textEndpoint);
        setFormImageEndpoint(provider.imageEndpoint || '');
        setFormDefaultHeaders(JSON.stringify(provider.defaultHeaders, null, 2));
        setFormRequestTemplate(JSON.stringify(provider.requestTemplate, null, 2));
        setFormResponseTextPath(provider.responseTextPath);
        setFormResponseImagePath(provider.responseImagePath || '');
        setShowDialog(true);
    };

    const handleSave = async () => {
        if (!formName || !formDisplayName || !formBaseUrl || !formTextEndpoint || !formResponseTextPath) {
            toast.error('Isi semua field yang diperlukan');
            return;
        }

        // Validate JSON
        try {
            JSON.parse(formDefaultHeaders);
            JSON.parse(formRequestTemplate);
        } catch (e) {
            toast.error('Format JSON tidak valid');
            return;
        }

        setSaving(true);
        try {
            const data: CreateProviderRequest = {
                name: formName,
                displayName: formDisplayName,
                description: formDescription,
                baseUrl: formBaseUrl,
                authType: formAuthType,
                authHeader: formAuthHeader,
                authPrefix: formAuthPrefix,
                textEndpoint: formTextEndpoint,
                imageEndpoint: formImageEndpoint || undefined,
                defaultHeaders: JSON.parse(formDefaultHeaders),
                requestTemplate: JSON.parse(formRequestTemplate),
                responseTextPath: formResponseTextPath,
                responseImagePath: formResponseImagePath || undefined,
            };

            if (editingProvider) {
                await updateProvider(editingProvider.id, data);
                toast.success('Provider updated');
            } else {
                await createProvider(data);
                toast.success('Provider created');
            }
            await loadProviders();
            setShowDialog(false);
            resetForm();
        } catch (error) {
            toast.error('Failed to save provider');
        } finally {
            setSaving(false);
        }
    };

    const handleDelete = async (id: string, name: string) => {
        if (confirm(`Hapus provider "${name}"?`)) {
            try {
                await deleteProvider(id);
                toast.success('Provider deleted');
                await loadProviders();
            } catch (error) {
                toast.error('Failed to delete provider');
            }
        }
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
                    <h2 className="text-lg font-semibold">API Providers</h2>
                    <p className="text-sm text-slate-500">Kelola konfigurasi provider AI (Gemini, OpenAI, Claude, dll)</p>
                </div>
                <Button onClick={() => { resetForm(); setShowDialog(true); }}>
                    <Plus className="mr-2 h-4 w-4" />
                    Add Provider
                </Button>
            </div>

            <div className="grid gap-4">
                {providers.map((provider) => (
                    <Card key={provider.id}>
                        <CardContent className="p-4">
                            <div className="flex items-start justify-between">
                                <div className="flex-1">
                                    <div className="flex items-center gap-2">
                                        <h3 className="font-semibold">{provider.displayName}</h3>
                                        <Badge variant="outline">{provider.name}</Badge>
                                        {!provider.isActive && <Badge variant="destructive">Inactive</Badge>}
                                    </div>
                                    <p className="text-sm text-slate-500 mt-1">{provider.description}</p>
                                    <div className="grid grid-cols-2 gap-2 mt-2 text-xs text-slate-400">
                                        <div className="flex items-center gap-1">
                                            <Globe className="h-3 w-3" />
                                            {provider.baseUrl}
                                        </div>
                                        <div className="flex items-center gap-1">
                                            <Key className="h-3 w-3" />
                                            Auth: {provider.authType}
                                        </div>
                                        <div className="flex items-center gap-1">
                                            <Code className="h-3 w-3" />
                                            Text: {provider.textEndpoint}
                                        </div>
                                        {provider.imageEndpoint && (
                                            <div className="flex items-center gap-1">
                                                🖼️ Image: {provider.imageEndpoint}
                                            </div>
                                        )}
                                    </div>
                                </div>
                                <div className="flex gap-1">
                                    <Button variant="ghost" size="sm" onClick={() => handleEdit(provider)}>
                                        <Edit2 className="h-4 w-4" />
                                    </Button>
                                    <Button variant="ghost" size="sm" onClick={() => handleDelete(provider.id, provider.displayName)}>
                                        <Trash2 className="h-4 w-4 text-red-500" />
                                    </Button>
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                ))}
            </div>

            {/* Add/Edit Provider Dialog */}
            <Dialog open={showDialog} onOpenChange={setShowDialog}>
                <DialogContent className="max-w-3xl max-h-[90vh] overflow-y-auto">
                    <DialogHeader>
                        <DialogTitle>{editingProvider ? 'Edit Provider' : 'Add New Provider'}</DialogTitle>
                        <DialogDescription>
                            Konfigurasi API provider untuk AI model
                        </DialogDescription>
                    </DialogHeader>

                    <Tabs defaultValue="basic" className="mt-4">
                        <TabsList className="grid w-full grid-cols-4">
                            <TabsTrigger value="basic">Basic</TabsTrigger>
                            <TabsTrigger value="auth">Authentication</TabsTrigger>
                            <TabsTrigger value="endpoints">Endpoints</TabsTrigger>
                            <TabsTrigger value="templates">Templates</TabsTrigger>
                        </TabsList>

                        <TabsContent value="basic" className="space-y-4">
                            <div>
                                <Label>Provider Name *</Label>
                                <Input
                                    placeholder="gemini, openai, anthropic"
                                    value={formName}
                                    onChange={(e) => setFormName(e.target.value)}
                                />
                                <p className="text-xs text-slate-500 mt-1">Unique identifier (used in API calls)</p>
                            </div>
                            <div>
                                <Label>Display Name *</Label>
                                <Input
                                    placeholder="Google Gemini, OpenAI GPT"
                                    value={formDisplayName}
                                    onChange={(e) => setFormDisplayName(e.target.value)}
                                />
                            </div>
                            <div>
                                <Label>Description</Label>
                                <Textarea
                                    placeholder="Description of this provider"
                                    value={formDescription}
                                    onChange={(e) => setFormDescription(e.target.value)}
                                    rows={2}
                                />
                            </div>
                            <div>
                                <Label>Base URL *</Label>
                                <Input
                                    placeholder="https://generativelanguage.googleapis.com/v1"
                                    value={formBaseUrl}
                                    onChange={(e) => setFormBaseUrl(e.target.value)}
                                />
                            </div>
                        </TabsContent>

                        <TabsContent value="auth" className="space-y-4">
                            <div>
                                <Label>Auth Type *</Label>
                                <Select value={formAuthType} onValueChange={setFormAuthType}>
                                    <SelectTrigger>
                                        <SelectValue />
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectItem value="bearer">Bearer Token</SelectItem>
                                        <SelectItem value="api_key">API Key (Header)</SelectItem>
                                        <SelectItem value="x-api-key">X-API-Key Header</SelectItem>
                                    </SelectContent>
                                </Select>
                            </div>
                            <div>
                                <Label>Auth Header *</Label>
                                <Input
                                    placeholder="Authorization"
                                    value={formAuthHeader}
                                    onChange={(e) => setFormAuthHeader(e.target.value)}
                                />
                            </div>
                            <div>
                                <Label>Auth Prefix</Label>
                                <Input
                                    placeholder="Bearer"
                                    value={formAuthPrefix}
                                    onChange={(e) => setFormAuthPrefix(e.target.value)}
                                />
                            </div>
                        </TabsContent>

                        <TabsContent value="endpoints" className="space-y-4">
                            <div>
                                <Label>Text Endpoint *</Label>
                                <Input
                                    placeholder="/models/{model}:generateContent"
                                    value={formTextEndpoint}
                                    onChange={(e) => setFormTextEndpoint(e.target.value)}
                                />
                                <p className="text-xs text-slate-500 mt-1">Use {'{model}'} as placeholder for model name</p>
                            </div>
                            <div>
                                <Label>Image Endpoint (Optional)</Label>
                                <Input
                                    placeholder="/models/{model}:generateContent"
                                    value={formImageEndpoint}
                                    onChange={(e) => setFormImageEndpoint(e.target.value)}
                                />
                            </div>
                            <div>
                                <Label>Response Text Path *</Label>
                                <Input
                                    placeholder="$.candidates[0].content.parts[0].text"
                                    value={formResponseTextPath}
                                    onChange={(e) => setFormResponseTextPath(e.target.value)}
                                />
                                <p className="text-xs text-slate-500 mt-1">JSON path to extract text from response</p>
                            </div>
                            <div>
                                <Label>Response Image Path (Optional)</Label>
                                <Input
                                    placeholder="$.candidates[0].content.parts[0].inlineData.data"
                                    value={formResponseImagePath}
                                    onChange={(e) => setFormResponseImagePath(e.target.value)}
                                />
                            </div>
                        </TabsContent>

                        <TabsContent value="templates" className="space-y-4">
                            <div>
                                <Label>Default Headers (JSON)</Label>
                                <div className="relative">
                                    <Textarea
                                        value={formDefaultHeaders}
                                        onChange={(e) => setFormDefaultHeaders(e.target.value)}
                                        rows={4}
                                        className="font-mono text-sm"
                                    />
                                    <button
                                        type="button"
                                        onClick={() => setShowJsonPreview(!showJsonPreview)}
                                        className="absolute top-2 right-2 text-slate-400 hover:text-slate-600"
                                    >
                                        {showJsonPreview ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                                    </button>
                                </div>
                                <p className="text-xs text-slate-500 mt-1">Example: {"{\"Content-Type\": \"application/json\"}"}</p>
                            </div>
                            <div>
                                <Label>Request Template * (JSON)</Label>
                                <Textarea
                                    value={formRequestTemplate}
                                    onChange={(e) => setFormRequestTemplate(e.target.value)}
                                    rows={8}
                                    className="font-mono text-sm"
                                />
                                <p className="text-xs text-slate-500 mt-1">
                                    Use placeholders: {'{model}'}, {'{system_prompt}'}, {'{prompt}'}
                                </p>
                                <p className="text-xs text-slate-500">
                                    Example for Gemini: {"{\"contents\": [{\"parts\": [{\"text\": \"{prompt}\"}]}]}"}
                                </p>
                            </div>
                        </TabsContent>
                    </Tabs>

                    <DialogFooter>
                        <Button variant="outline" onClick={() => setShowDialog(false)}>Cancel</Button>
                        <Button onClick={handleSave} disabled={saving}>
                            {saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                            {editingProvider ? 'Update' : 'Create'}
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </div>
    );
}