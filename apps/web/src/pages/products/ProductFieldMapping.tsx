import { useState, useEffect, useCallback, useMemo } from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@kana-consultant/ui-kit";
import { Button } from "@kana-consultant/ui-kit";
import { Input } from "@kana-consultant/ui-kit";
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
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@kana-consultant/ui-kit";
import {Copy, Plus, Trash2, Check } from "lucide-react";

// ==================== TYPES ====================
interface ProductFieldMappingProps {
    domain: string | undefined;
    metaConfig?: MetaConfig;
    sitemapConfig?: any;
    onChange: (metaConfig?: string, sitemapConfig?: string) => void;
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
    enabled: false,
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
    enabled: false,
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

// ==================== SITEMAP MANAGER (DENGAN SELECT PLACEHOLDER) ====================
function SitemapManager({
    initialConfig,
    baseUrl = "https://domainanda.com",
    onChange,
}: {
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
        const defaults = ["0.7", "0.8", "0.9", "1.0", "weekly", "daily", "monthly", "image_url", "updated_at", "created_at"];
        return defaults.sort();
    }, []);

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
                            <Switch checked={config.enabled} onCheckedChange={(v: any) => updateConfig({ enabled: v })} />
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
                                        <Input placeholder="/about" value={newStaticUrl} onChange={(e: any) => setNewStaticUrl(e.target.value)} />
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
                                    <Input value={url.loc} onChange={(e: any) => updateStaticUrl(idx, "loc", e.target.value)} className="flex-1" />
                                    <Input type="number" step="0.1" value={url.priority} onChange={(e: any) => updateStaticUrl(idx, "priority", parseFloat(e.target.value))} className="w-24" />
                                    <Select value={url.changefreq} onValueChange={(v: any) => updateStaticUrl(idx, "changefreq", v)}>
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
                            <Input value={config.dynamicConfig.urlPattern} onChange={(e: any) => updateConfig({ dynamicConfig: { ...config.dynamicConfig, urlPattern: e.target.value } })} placeholder="/p/{id}" />
                        </div>
                        <div className="grid grid-cols-2 gap-3">
                            {(["prioritySource", "changefreqSource", "imageSource", "lastmodSource"] as const).map((key) => (
                                <div key={key}>
                                    <Label className="text-xs">{key.replace("Source", "")}</Label>
                                    <Select
                                        value={config.dynamicConfig[key]}
                                        onValueChange={(v: any) => updateConfig({ dynamicConfig: { ...config.dynamicConfig, [key]: v } })}
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
    domain,
    metaConfig: initialMetaConfig,
    sitemapConfig: initialSitemapConfig,
    onChange,
}: ProductFieldMappingProps) {
    const [activeTab, setActiveTab] = useState<"meta" | "sitemap" | "payload">("meta");

    const [metaConfig, setMetaConfig] = useState<MetaConfig>(() => ({ ...defaultMetaConfig, ...getObjectValue(initialMetaConfig) }));
    const [sitemapConfig, setSitemapConfig] = useState<SitemapConfig>(() => ({ ...defaultSitemapConfig, ...getObjectValue(initialSitemapConfig) }));

    const [userInput, setUserInput] = useState({ title: "", topic: "" });

    const handleSitemapChange = (newSitemap: SitemapConfig) => {
        setSitemapConfig(newSitemap);
        onChange(JSON.stringify(metaConfig, null, 2), JSON.stringify(newSitemap, null, 2));
    };

    return (
        <div className="border rounded-xl p-5 space-y-6 mt-2">
            <Tabs value={activeTab} onValueChange={(v: any) => setActiveTab(v as any)}>
                <TabsList className="grid w-full grid-cols-2">
                    <TabsTrigger value="sitemap">🗺️ Sitemap</TabsTrigger>
                    <TabsTrigger value="payload">🚀 Payload</TabsTrigger>
                </TabsList>
                <TabsContent value="sitemap" className="pt-4">
                    <SitemapManager
                        baseUrl={domain}
                        onChange={handleSitemapChange}
                    />
                </TabsContent>

                <TabsContent value="payload" className="pt-4">
                    <PayloadPreview metaConfig={metaConfig} sitemapConfig={sitemapConfig} userInput={userInput} />
                </TabsContent>
            </Tabs>
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