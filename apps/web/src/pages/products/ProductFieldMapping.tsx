import { useState, useEffect, useRef, useCallback, useMemo } from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@kana-consultant/ui-kit";
import { Button } from "@kana-consultant/ui-kit";
import { Input } from "@kana-consultant/ui-kit";
import { Textarea } from "@kana-consultant/ui-kit";
import { Label } from "@kana-consultant/ui-kit";
import { Switch } from "@kana-consultant/ui-kit";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@kana-consultant/ui-kit";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@kana-consultant/ui-kit";
import { SimpleJsonBuilder } from "./SimpleJsonBuilder";
import { Eye, Send, Copy, Plus, Trash2, Check } from "lucide-react";

// ==================== TYPES ====================
interface ProductFieldMappingProps {
    fieldMapping: any;
    metaConfig?: any;
    sitemapConfig?: any;
    onChange: (fieldMapping: string, metaConfig?: string, sitemapConfig?: string) => void;
}

interface MetaConfig {
    enabled: boolean;
    defaultTags: {
        charset: string;
        viewport: string;
        robots: string;
        generator: string;
    };
    dynamicTags: {
        titleSource: string;
        descriptionSource: string;
        imageSource: string;
        customTags: Record<string, string>;
    };
}

interface SitemapConfig {
    enabled: boolean;
    staticUrls: {
        loc: string;
        priority: number;
        changefreq: "always" | "hourly" | "daily" | "weekly" | "monthly" | "yearly" | "never";
    }[];
    dynamicConfig: {
        urlPattern: string;
        prioritySource: string;
        changefreqSource: string;
        imageSource: string;
        lastmodSource: string;
    };
}

// ==================== UTILITIES ====================
const extractFields = (obj: any, prefix = ""): string[] => {
    if (!obj || typeof obj !== "object") return [];
    let fields: string[] = [];
    for (const [key, value] of Object.entries(obj)) {
        const fullPath = prefix ? `${prefix}.${key}` : key;
        fields.push(fullPath);
        if (typeof value === "object" && value !== null && !Array.isArray(value)) {
            fields = [...fields, ...extractFields(value, fullPath)];
        }
    }
    return fields;
};

const getObjectValue = (input: any): any => {
    if (typeof input === 'string') {
        try { return JSON.parse(input || "{}"); } catch { return {}; }
    }
    return input || {};
};

// Default Configs
const defaultMetaConfig: MetaConfig = {
    enabled: true,
    defaultTags: {
        charset: "UTF-8",
        viewport: "width=device-width, initial-scale=1.0",
        robots: "index, follow",
        generator: "AI Content Generator v1.0",
    },
    dynamicTags: {
        titleSource: "title",
        descriptionSource: "excerpt",
        imageSource: "image_url",
        customTags: {},
    },
};

const defaultSitemapConfig: SitemapConfig = {
    enabled: true,
    staticUrls: [
        { loc: "/", priority: 1.0, changefreq: "daily" },
        { loc: "/privacy", priority: 0.3, changefreq: "yearly" },
        { loc: "/terms", priority: 0.3, changefreq: "yearly" },
    ],
    dynamicConfig: {
        urlPattern: "/p/{id}",
        prioritySource: "0.7",
        changefreqSource: "weekly",
        imageSource: "image_url",
        lastmodSource: "updated_at",
    },
};

// ==================== META TAG MANAGER ====================
function MetaTagManager({
    initialConfig,
    onChange,
}: {
    fieldMapping?: any;
    initialConfig?: MetaConfig;
    onChange: (config: MetaConfig) => void;
}) {
    const [config, setConfig] = useState<MetaConfig>(() => ({
        ...defaultMetaConfig,
        ...getObjectValue(initialConfig),
    }));

    const isMounted = useRef(false);

    const [customDialogOpen, setCustomDialogOpen] = useState(false);
    const [newTagKey, setNewTagKey] = useState("");
    const [newTagValue, setNewTagValue] = useState("");
    const [editingCustomTag, setEditingCustomTag] = useState<{ oldKey: string; key: string; value: string } | null>(null);

    // LIST UNTUK DETECTION - SEMUA DALAM FORMAT {}
    const PLACEHOLDER_VALUES = ["{title}", "{topic}", "{content}", "{excerpt}", "{image_url}", "{scheduled_for}"];

    // DEFAULT VALUES UNTUK SELECT (jika data kosong)
    const DEFAULT_DYNAMIC_TAGS = {
        titleSource: "{title}",
        descriptionSource: "{excerpt}",
        imageSource: "{image_url}",
    };

    const previewMeta = useMemo(() => {
        if (!config.enabled) return ["<!-- Meta tags disabled -->"];
        const tags = [
            `<!-- ========== DEFAULT TAGS ========== -->`,
            `<meta charset="${config.defaultTags.charset}">`,
            `<meta name="viewport" content="${config.defaultTags.viewport}">`,
            `<meta name="robots" content="${config.defaultTags.robots}">`,
            `<meta name="generator" content="${config.defaultTags.generator}">`,
            ``,
            `<!-- ========== DYNAMIC TAGS ========== -->`,
            `<title>${config.dynamicTags.titleSource || DEFAULT_DYNAMIC_TAGS.titleSource}</title>`,
            `<meta name="description" content="${config.dynamicTags.descriptionSource || DEFAULT_DYNAMIC_TAGS.descriptionSource}">`,
            `<meta property="og:title" content="${config.dynamicTags.titleSource || DEFAULT_DYNAMIC_TAGS.titleSource}">`,
            `<meta property="og:description" content="${config.dynamicTags.descriptionSource || DEFAULT_DYNAMIC_TAGS.descriptionSource}">`,
            `<meta property="og:image" content="${config.dynamicTags.imageSource || DEFAULT_DYNAMIC_TAGS.imageSource}">`,
            `<meta property="og:type" content="article">`,
            `<meta name="twitter:card" content="summary_large_image">`,
            `<meta name="twitter:title" content="${config.dynamicTags.titleSource || DEFAULT_DYNAMIC_TAGS.titleSource}">`,
            `<meta name="twitter:description" content="${config.dynamicTags.descriptionSource || DEFAULT_DYNAMIC_TAGS.descriptionSource}">`,
            `<meta name="twitter:image" content="${config.dynamicTags.imageSource || DEFAULT_DYNAMIC_TAGS.imageSource}">`,
        ];

        if (Object.keys(config.dynamicTags.customTags).length > 0) {
            tags.push(``);
            tags.push(`<!-- ========== CUSTOM TAGS ========== -->`);
            Object.entries(config.dynamicTags.customTags).forEach(([k, v]) => tags.push(`<meta name="${k}" content="${v}">`));
        }
        return tags;
    }, [config]);

    useEffect(() => {
        if (!isMounted.current) { isMounted.current = true; return; }
        const parsed = getObjectValue(initialConfig);
        if (Object.keys(parsed).length > 0) {
            setConfig({ ...defaultMetaConfig, ...parsed });
        } else {
            // Jika tidak ada data (create new), set default values untuk dynamicTags
            setConfig(prev => ({
                ...prev,
                dynamicTags: {
                    ...prev.dynamicTags,
                    ...DEFAULT_DYNAMIC_TAGS,
                }
            }));
        }
    }, [initialConfig]);

    const updateConfig = useCallback((updates: Partial<MetaConfig>) => {
        setConfig(prev => {
            const newConfig = { ...prev, ...updates };
            onChange(newConfig);
            return newConfig;
        });
    }, [onChange]);

    // Format value ke dalam bentuk {}
    const formatToPlaceholder = (value: string): string => {
        const trimmed = value.trim();
        if (!trimmed) return "";
        // Jika sudah dalam format {xxx}
        if (/^\{[a-z_]+\}$/.test(trimmed)) return trimmed;
        // Jika sudah mengandung {} di dalamnya (contoh: "judul: {title}")
        if (trimmed.includes("{") && trimmed.includes("}")) return trimmed;
        // Jika tidak ada {}, bungkus dengan {}
        return `{${trimmed}}`;
    };

    const addCustomTag = () => {
        if (!newTagKey.trim()) return;
        
        const finalValue = formatToPlaceholder(newTagValue);
        
        updateConfig({
            dynamicTags: {
                ...config.dynamicTags,
                customTags: { ...config.dynamicTags.customTags, [newTagKey.trim()]: finalValue }
            }
        });
        setNewTagKey(""); 
        setNewTagValue(""); 
        setCustomDialogOpen(false);
    };

    const openEditCustomTag = (key: string, value: string) => {
        setEditingCustomTag({ oldKey: key, key: key, value: value });
    };

    const saveEditCustomTag = () => {
        if (!editingCustomTag) return;
        const { oldKey, key, value } = editingCustomTag;
        if (!key.trim()) return;
        
        const newTags = { ...config.dynamicTags.customTags };
        if (oldKey !== key) delete newTags[oldKey];
        newTags[key.trim()] = formatToPlaceholder(value);
        
        updateConfig({ dynamicTags: { ...config.dynamicTags, customTags: newTags } });
        setEditingCustomTag(null);
    };

    const removeCustomTag = (key: string) => {
        const newTags = { ...config.dynamicTags.customTags };
        delete newTags[key];
        updateConfig({ dynamicTags: { ...config.dynamicTags, customTags: newTags } });
    };

    const handleDynamicTagChange = (field: string, value: string) => {
        updateConfig({ dynamicTags: { ...config.dynamicTags, [field]: value } });
    };

    // Get value untuk select, fallback ke default jika kosong
    const getSelectValue = (field: string) => {
        const value = config.dynamicTags[field as keyof typeof config.dynamicTags];
        if (value && PLACEHOLDER_VALUES.includes(value as string)) return value;
        return DEFAULT_DYNAMIC_TAGS[field as keyof typeof DEFAULT_DYNAMIC_TAGS];
    };

    const [copied, setCopied] = useState(false);
    const handleCopy = () => {
        navigator.clipboard.writeText(previewMeta.join('\n'));
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
    };

    return (
        <div className="space-y-4">
            <Tabs defaultValue="config">
                <TabsList className="grid w-full grid-cols-2">
                    <TabsTrigger value="config">⚙️ Konfigurasi</TabsTrigger>
                    <TabsTrigger value="preview">👁️ Preview</TabsTrigger>
                </TabsList>

                <TabsContent value="config" className="pt-4 space-y-4">
                    {/* Enable */}
                    <div className="flex items-center justify-between p-3 bg-slate-50 dark:bg-slate-800 rounded-lg">
                        <div className="flex items-center gap-2">
                            <Switch checked={config.enabled} onCheckedChange={(v : any) => updateConfig({ enabled: v })} />
                            <Label>Aktifkan Meta Tags</Label>
                        </div>
                    </div>

                    {/* Default Tags */}
                    <div className="border rounded-lg p-4">
                        <h4 className="font-medium mb-3">Default Tags (Tidak Bervariasi)</h4>
                        <div className="grid grid-cols-2 gap-3">
                            {Object.entries(config.defaultTags).map(([key, val]) => (
                                <div key={key}>
                                    <Label className="text-xs capitalize">{key}</Label>
                                    <Input 
                                        value={val} 
                                        onChange={(e : any) => updateConfig({ defaultTags: { ...config.defaultTags, [key]: e.target.value } })} 
                                    />
                                </div>
                            ))}
                        </div>
                    </div>

                    {/* Dynamic Tags */}
                    <div className="border rounded-lg p-4 space-y-4">
                        <h4 className="font-medium mb-3">Dynamic Tags (Bervariasi per Konten)</h4>
                        <div className="text-xs text-muted-foreground mb-2">
                            💡 Pilih placeholder dari daftar (disimpan dalam format {'{...}'})
                        </div>
                        {(["titleSource", "descriptionSource", "imageSource"] as const).map((field) => {
                            const selectValue = getSelectValue(field);
                            
                            return (
                                <div key={field}>
                                    <Label className="text-xs mb-1 block">{field.replace("Source", "")} Source</Label>
                                    <Select value={selectValue as string} onValueChange={(v : any) => handleDynamicTagChange(field, v)}>
                                        <SelectTrigger>
                                            <SelectValue placeholder="Pilih placeholder..." />
                                        </SelectTrigger>
                                        <SelectContent>
                                            {PLACEHOLDER_VALUES.map(opt => (
                                                <SelectItem key={opt} value={opt}>{opt}</SelectItem>
                                            ))}
                                        </SelectContent>
                                    </Select>
                                    <p className="text-xs text-muted-foreground mt-1">
                                        Nilai tersimpan: <code className="bg-slate-100 px-1 rounded">{config.dynamicTags[field] || DEFAULT_DYNAMIC_TAGS[field]}</code>
                                    </p>
                                </div>
                            );
                        })}
                    </div>

                    {/* Custom Tags */}
                    <div className="border rounded-lg p-4">
                        <div className="flex justify-between items-center mb-3">
                            <h4 className="font-medium">Custom Meta Tags</h4>
                            <Dialog open={customDialogOpen} onOpenChange={setCustomDialogOpen}>
                                <DialogTrigger asChild>
                                    <Button size="sm"><Plus className="w-4 h-4 mr-1" /> Tambah</Button>
                                </DialogTrigger>
                                <DialogContent>
                                    <DialogHeader>
                                        <DialogTitle>Tambah Custom Meta Tag</DialogTitle>
                                        <DialogDescription>Masukkan nama dan nilai meta tag</DialogDescription>
                                    </DialogHeader>
                                    <div className="space-y-4 py-4">
                                        <div>
                                            <Label>Nama Tag</Label>
                                            <Input placeholder="author" value={newTagKey} onChange={(e : any) => setNewTagKey(e.target.value)} />
                                        </div>
                                        <div>
                                            <Label>Nilai</Label>
                                            <Input 
                                                placeholder="{title} - John Doe" 
                                                value={newTagValue} 
                                                onChange={(e : any) => setNewTagValue(e.target.value)} 
                                            />
                                            <p className="text-xs text-muted-foreground mt-1">
                                                Gunakan {'{title}'}, {'{excerpt}'}, dll untuk placeholder
                                            </p>
                                        </div>
                                    </div>
                                    <DialogFooter>
                                        <Button variant="outline" onClick={() => setCustomDialogOpen(false)}>Batal</Button>
                                        <Button onClick={addCustomTag}>Tambahkan</Button>
                                    </DialogFooter>
                                </DialogContent>
                            </Dialog>
                        </div>

                        <div className="text-xs text-muted-foreground mb-3">
                            💡 Tips: Nilai akan otomatis dibungkus {'{}'} jika belum ada
                        </div>

                        {Object.entries(config.dynamicTags.customTags).map(([key, value]) => (
                            <div key={key} className="flex gap-2 mb-3 items-center">
                                {editingCustomTag?.oldKey === key ? (
                                    <>
                                        <Input 
                                            value={editingCustomTag.key} 
                                            onChange={(e : any) => setEditingCustomTag({ ...editingCustomTag, key: e.target.value })} 
                                            className="w-48" 
                                            autoFocus
                                        />
                                        <Input 
                                            value={editingCustomTag.value} 
                                            onChange={(e : any) => setEditingCustomTag({ ...editingCustomTag, value: e.target.value })} 
                                            className="flex-1" 
                                        />
                                        <Button size="sm" onClick={saveEditCustomTag}>💾</Button>
                                        <Button size="sm" variant="ghost" onClick={() => setEditingCustomTag(null)}>✖️</Button>
                                    </>
                                ) : (
                                    <>
                                        <div className="w-48 px-3 py-2 bg-slate-100 dark:bg-slate-800 rounded text-sm font-mono">
                                            {key}
                                        </div>
                                        <div className="flex-1 px-3 py-2 bg-slate-50 dark:bg-slate-900 rounded text-sm font-mono text-blue-600">
                                            {value}
                                        </div>
                                        <Button variant="ghost" size="icon" onClick={() => openEditCustomTag(key, value)}>
                                            ✏️
                                        </Button>
                                        <Button variant="ghost" size="icon" onClick={() => removeCustomTag(key)}>
                                            <Trash2 className="w-4 h-4 text-red-500" />
                                        </Button>
                                    </>
                                )}
                            </div>
                        ))}
                    </div>
                </TabsContent>

                <TabsContent value="preview" className="pt-4">
                    <pre className="bg-slate-900 p-4 rounded-lg max-h-96 overflow-auto text-emerald-400 text-sm font-mono whitespace-pre-wrap">
                        {previewMeta.join('\n')}
                    </pre>
                    <Button onClick={handleCopy} className="mt-3 w-full" variant="outline">
                        {copied ? <Check className="mr-2 w-4 h-4" /> : <Copy className="mr-2 w-4 h-4" />}
                        {copied ? "Tersalin!" : "Salin Preview"}
                    </Button>
                </TabsContent>
            </Tabs>
        </div>
    );
}

// ==================== SITEMAP MANAGER (DENGAN SELECT PLACEHOLDER) ====================
function SitemapManager({
    fieldMapping,
    initialConfig,
    baseUrl = "https://domainanda.com",
    onChange,
}: {
    fieldMapping?: any;
    initialConfig?: SitemapConfig;
    baseUrl?: string;
    onChange: (config: SitemapConfig) => void;
}) {
    const [config, setConfig] = useState<SitemapConfig>(() => ({
        ...defaultSitemapConfig,
        ...getObjectValue(initialConfig),
    }));

    const [staticDialogOpen, setStaticDialogOpen] = useState(false);
    const [newStaticUrl, setNewStaticUrl] = useState("");
    const [copied, setCopied] = useState(false);

    const selectOptions = useMemo(() => {
        const fields = extractFields(getObjectValue(fieldMapping));
        const defaults = ["0.7", "0.8", "0.9", "1.0", "weekly", "daily", "monthly", "image_url", "updated_at", "created_at"];
        return [...new Set([...defaults, ...fields])].sort();
    }, [fieldMapping]);

    const sitemapXml = useMemo(() => {
        const staticXml = config.staticUrls.map(url => `
        <url>
            <loc>${baseUrl}${url.loc}</loc>
            <lastmod>{timestamp}</lastmod>
            <changefreq>${url.changefreq}</changefreq>
            <priority>${url.priority}</priority>
        </url>`).join("");

        const dynamicXml = `
        <!-- Dynamic URLs -->
        <url>
            <loc>${baseUrl}${config.dynamicConfig.urlPattern}</loc>
            <lastmod>{${config.dynamicConfig.lastmodSource}}</lastmod>
            <changefreq>{${config.dynamicConfig.changefreqSource}}</changefreq>
            <priority>{${config.dynamicConfig.prioritySource}}</priority>
            <image:image>
            <image:loc>${baseUrl}/{${config.dynamicConfig.imageSource}}</image:loc>
            </image:image>
        </url>`;

        return `<?xml version="1.0" encoding="UTF-8"?>
        <urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
                xmlns:image="http://www.google.com/schemas/sitemap-image/1.1">
            ${staticXml}
            ${dynamicXml}
        </urlset>`;

    }, [config, baseUrl]);
    useEffect(() => {
        const parsed = getObjectValue(initialConfig);
        if (Object.keys(parsed).length > 0) setConfig({ ...defaultSitemapConfig, ...parsed });
    }, [initialConfig]);

    const updateConfig = useCallback((updates: Partial<SitemapConfig>) => {
        setConfig(prev => {
            const newConfig = { ...prev, ...updates };
            onChange(newConfig);
            return newConfig;
        });
    }, [onChange]);

    const addStaticUrl = () => {
        if (!newStaticUrl.trim()) return;
        updateConfig({ staticUrls: [...config.staticUrls, { loc: newStaticUrl.trim(), priority: 0.5, changefreq: "monthly" }] });
        setNewStaticUrl("");
        setStaticDialogOpen(false);
    };

    const updateStaticUrl = (index: number, field: string, value: any) => {
        const newUrls = [...config.staticUrls];
        newUrls[index] = { ...newUrls[index], [field]: value };
        updateConfig({ staticUrls: newUrls });
    };

    const removeStaticUrl = (index: number) => {
        updateConfig({ staticUrls: config.staticUrls.filter((_, i) => i !== index) });
    };

    return (
        <div className="space-y-4">
            <Tabs defaultValue="config">
                <TabsList className="grid w-full grid-cols-2">
                    <TabsTrigger value="config">⚙️ Konfigurasi</TabsTrigger>
                    <TabsTrigger value="preview">👁️ Preview XML</TabsTrigger>
                </TabsList>

                <TabsContent value="config" className="pt-4 space-y-4">
                    {/* Enable Switch */}
                    <div className="flex items-center justify-between p-3 bg-slate-50 dark:bg-slate-800 rounded-lg">
                        <div className="flex items-center gap-2">
                            <Switch checked={config.enabled} onCheckedChange={(v : any) => updateConfig({ enabled: v })} />
                            <Label>Aktifkan Sitemap</Label>
                        </div>
                    </div>

                    {/* Static URLs */}
                    <div className="border rounded-lg p-4">
                        <div className="flex justify-between mb-3">
                            <h4 className="font-medium">Static URLs</h4>
                            <Dialog open={staticDialogOpen} onOpenChange={setStaticDialogOpen}>
                                <DialogTrigger asChild>
                                    <Button size="sm"><Plus className="w-4 h-4 mr-1" /> Tambah URL</Button>
                                </DialogTrigger>
                                <DialogContent>
                                    <DialogHeader>
                                        <DialogTitle>Tambah Static URL</DialogTitle>
                                    </DialogHeader>
                                    <div className="py-4">
                                        <Label>URL Path</Label>
                                        <Input placeholder="/about" value={newStaticUrl} onChange={(e : any) => setNewStaticUrl(e.target.value)} />
                                    </div>
                                    <DialogFooter>
                                        <Button variant="outline" onClick={() => setStaticDialogOpen(false)}>Batal</Button>
                                        <Button onClick={addStaticUrl}>Tambahkan</Button>
                                    </DialogFooter>
                                </DialogContent>
                            </Dialog>
                        </div>
                        {/* Static list tetap sama */}
                        <div className="space-y-2">
                            {config.staticUrls.map((url, idx) => (
                                <div key={idx} className="flex gap-2 items-center">
                                    <Input value={url.loc} onChange={(e : any) => updateStaticUrl(idx, "loc", e.target.value)} className="flex-1" />
                                    <Input type="number" step="0.1" value={url.priority} onChange={(e : any) => updateStaticUrl(idx, "priority", parseFloat(e.target.value))} className="w-24" />
                                    <Select value={url.changefreq} onValueChange={(v : any) => updateStaticUrl(idx, "changefreq", v)}>
                                        <SelectTrigger className="w-32"><SelectValue /></SelectTrigger>
                                        <SelectContent>
                                            {["always", "hourly", "daily", "weekly", "monthly", "yearly", "never"].map(v => <SelectItem key={v} value={v}>{v}</SelectItem>)}
                                        </SelectContent>
                                    </Select>
                                    <Button variant="ghost" size="icon" onClick={() => removeStaticUrl(idx)}><Trash2 className="w-4 h-4 text-red-500" /></Button>
                                </div>
                            ))}
                        </div>
                    </div>

                    {/* Dynamic Config - BISA PILIH PLACEHOLDER */}
                    <div className="border rounded-lg p-4 space-y-3">
                        <h4 className="font-medium">Dynamic URLs</h4>
                        <div>
                            <Label>URL Pattern</Label>
                            <Input value={config.dynamicConfig.urlPattern} onChange={(e : any) => updateConfig({ dynamicConfig: { ...config.dynamicConfig, urlPattern: e.target.value } })} placeholder="/p/{id}" />
                        </div>
                        <div className="grid grid-cols-2 gap-3">
                            {(["prioritySource", "changefreqSource", "imageSource", "lastmodSource"] as const).map((key) => (
                                <div key={key}>
                                    <Label className="text-xs">{key.replace("Source", "")}</Label>
                                    <Select
                                        value={config.dynamicConfig[key]}
                                        onValueChange={(v : any) => updateConfig({ dynamicConfig: { ...config.dynamicConfig, [key]: v } })}
                                    >
                                        <SelectTrigger>
                                            <SelectValue />
                                        </SelectTrigger>
                                        <SelectContent>
                                            {selectOptions.map(opt => (
                                                <SelectItem key={opt} value={opt}>{opt}</SelectItem>
                                            ))}
                                        </SelectContent>
                                    </Select>
                                </div>
                            ))}
                        </div>
                    </div>
                </TabsContent>

                <TabsContent value="preview" className="pt-4">
                    <pre className="bg-slate-900 p-4 rounded-lg max-h-96 overflow-auto text-xs text-emerald-400 font-mono">{sitemapXml}</pre>
                    <Button onClick={() => { navigator.clipboard.writeText(sitemapXml as string); setCopied(true); setTimeout(() => setCopied(false), 2000); }} className="mt-3 w-full" variant="outline">
                        {copied ? <Check className="mr-2" /> : <Copy className="mr-2" />} {copied ? "Tersalin!" : "Salin XML"}
                    </Button>
                </TabsContent>
            </Tabs>
        </div>
    );
}
// ==================== MAIN COMPONENT ====================
export function ProductFieldMapping({
    fieldMapping,
    metaConfig: initialMetaConfig,
    sitemapConfig: initialSitemapConfig,
    onChange,
}: ProductFieldMappingProps) {
    const [activeTab, setActiveTab] = useState<"form" | "raw" | "meta" | "sitemap" | "payload">("form");
    const [rawJson, setRawJson] = useState(() => JSON.stringify(getObjectValue(fieldMapping), null, 2));
    const [previewObject, setPreviewObject] = useState(getObjectValue(fieldMapping));

    const [metaConfig, setMetaConfig] = useState<MetaConfig>(() => ({ ...defaultMetaConfig, ...getObjectValue(initialMetaConfig) }));
    const [sitemapConfig, setSitemapConfig] = useState<SitemapConfig>(() => ({ ...defaultSitemapConfig, ...getObjectValue(initialSitemapConfig) }));

    const [userInput, setUserInput] = useState({ title: "", topic: "" });

    useEffect(() => {
        const obj = getObjectValue(fieldMapping);
        setRawJson(JSON.stringify(obj, null, 2));
        setPreviewObject(obj);
    }, [fieldMapping]);

    const handleBuilderChange = (newValue: any) => {
        setPreviewObject(newValue);
        onChange(
            JSON.stringify(newValue, null, 2),
            JSON.stringify(metaConfig, null, 2),
            JSON.stringify(sitemapConfig, null, 2)
        );
    };

    const handleMetaChange = (newMeta: MetaConfig) => {
        setMetaConfig(newMeta);
        onChange(JSON.stringify(previewObject, null, 2), JSON.stringify(newMeta, null, 2), JSON.stringify(sitemapConfig, null, 2));
    };

    const handleSitemapChange = (newSitemap: SitemapConfig) => {
        setSitemapConfig(newSitemap);
        onChange(JSON.stringify(previewObject, null, 2), JSON.stringify(metaConfig, null, 2), JSON.stringify(newSitemap, null, 2));
    };

    const saveRawJson = () => {
        try {
            const parsed = JSON.parse(rawJson);
            setPreviewObject(parsed);
            onChange(rawJson, JSON.stringify(metaConfig, null, 2), JSON.stringify(sitemapConfig, null, 2));
            setActiveTab("form");
        } catch {
            alert("JSON tidak valid!");
        }
    };

    return (
        <div className="border rounded-xl p-5 space-y-6">
            <Tabs value={activeTab} onValueChange={(v : any) => setActiveTab(v as any)}>
                <TabsList className="grid w-full grid-cols-5">
                    <TabsTrigger value="form">📝 Drag & Drop</TabsTrigger>
                    <TabsTrigger value="raw">📄 Raw JSON</TabsTrigger>
                    <TabsTrigger value="meta">🏷️ Meta Tags</TabsTrigger>
                    <TabsTrigger value="sitemap">🗺️ Sitemap</TabsTrigger>
                    <TabsTrigger value="payload">🚀 Payload</TabsTrigger>
                </TabsList>

                <TabsContent value="form" className="pt-4">
                    <SimpleJsonBuilder value={previewObject} onChange={handleBuilderChange} />
                </TabsContent>

                <TabsContent value="raw" className="pt-4 space-y-3">
                    <Textarea value={rawJson} onChange={(e : any) => setRawJson(e.target.value)} className="h-96 font-mono text-sm" />
                    <Button onClick={saveRawJson} className="w-full">💾 Simpan JSON</Button>
                </TabsContent>

                <TabsContent value="meta" className="pt-4">
                    <MetaTagManager fieldMapping={fieldMapping} initialConfig={metaConfig} onChange={handleMetaChange} />
                </TabsContent>

                <TabsContent value="sitemap" className="pt-4">
                    <SitemapManager initialConfig={sitemapConfig} onChange={handleSitemapChange} />
                </TabsContent>

                <TabsContent value="payload" className="pt-4">
                    <PayloadPreview fieldMapping={fieldMapping} metaConfig={metaConfig} sitemapConfig={sitemapConfig} userInput={userInput} />
                </TabsContent>
            </Tabs>

            <div className="border-t pt-4">
                <div className="flex items-center gap-2 mb-2">
                    <Eye className="w-4 h-4" />
                    <h4 className="font-medium">Preview JSON Structure</h4>
                </div>
                <pre className="bg-slate-100 dark:bg-slate-900 p-4 rounded-lg text-xs overflow-auto max-h-60 font-mono">
                    {JSON.stringify(previewObject, null, 2)}
                </pre>
            </div>
        </div>
    );
}

// PayloadPreview Component
function PayloadPreview({ fieldMapping, metaConfig, sitemapConfig, userInput = {} }: any) {
    const [copied, setCopied] = useState(false);
    const payloadStr = JSON.stringify({
        api_payload: getObjectValue(fieldMapping),
        meta_config: metaConfig,
        sitemap_config: sitemapConfig,
        user_input: userInput,
    }, null, 2);

    return (
        <div className="space-y-3">
            <div className="flex justify-between">
                <div className="font-medium">JSON Payload</div>
                <Button size="sm" variant="outline" onClick={() => { navigator.clipboard.writeText(payloadStr); setCopied(true); setTimeout(() => setCopied(false), 2000); }}>
                    {copied ? <Check className="mr-2" /> : <Copy className="mr-2" />} {copied ? "Tersalin!" : "Salin"}
                </Button>
            </div>
            <pre className="bg-slate-900 p-4 rounded-lg overflow-auto max-h-96 text-xs text-emerald-400 font-mono">{payloadStr}</pre>
        </div>
    );
}