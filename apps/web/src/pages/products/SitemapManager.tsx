// components/SitemapManager.tsx
import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { MapPin, Download, Copy, Check } from "lucide-react";

interface SitemapManagerProps {
    fieldMapping: any;
    baseUrl?: string;
    onChange?: (sitemapConfig: SitemapConfig) => void;
}

export interface SitemapConfig {
    enabled: boolean;
    staticUrls: {
        loc: string;
        priority: number;
        changefreq: "always" | "hourly" | "daily" | "weekly" | "monthly" | "yearly" | "never";
    }[];
    dynamicConfig: {
        urlPattern: string;      // contoh: /p/{id}
        prioritySource: string;   // field key untuk priority
        changefreqSource: string; // field key untuk changefreq
        imageSource: string;      // field key untuk gambar
        lastmodSource: string;    // field key untuk last modified
    };
}

export function SitemapManager({ fieldMapping, baseUrl = "https://domainanda.com", onChange }: SitemapManagerProps) {
    const [config, setConfig] = useState<SitemapConfig>({
        enabled: true,
        staticUrls: [
            { loc: "/", priority: 1.0, changefreq: "daily" },
            { loc: "/privacy", priority: 0.3, changefreq: "yearly" },
            { loc: "/terms", priority: 0.3, changefreq: "yearly" },
        ],
        dynamicConfig: {
            urlPattern: "/p/{id}",
            prioritySource: "sitemap_priority",
            changefreqSource: "sitemap_changefreq",
            imageSource: "image_url",
            lastmodSource: "updated_at",
        },
    });

    const [sitemapXml, setSitemapXml] = useState<string>("");
    const [copied, setCopied] = useState(false);
    const [availableFields, setAvailableFields] = useState<string[]>([]);

    // Extract available fields
    useEffect(() => {
        const extractFields = (obj: any, prefix = ""): string[] => {
            let fields: string[] = [];
            if (!obj || typeof obj !== "object") return fields;
            
            for (const [key, value] of Object.entries(obj)) {
                const fullPath = prefix ? `${prefix}.${key}` : key;
                fields.push(fullPath);
                if (typeof value === "object" && value !== null && !Array.isArray(value)) {
                    fields = [...fields, ...extractFields(value, fullPath)];
                }
            }
            return fields;
        };
        
        const mappingObj = typeof fieldMapping === "string" ? JSON.parse(fieldMapping || "{}") : fieldMapping;
        setAvailableFields(extractFields(mappingObj));
    }, [fieldMapping]);

    // Generate sitemap XML preview
    useEffect(() => {
        const staticXml = config.staticUrls.map(url => `
  <url>
    <loc>${baseUrl}${url.loc}</loc>
    <lastmod>{timestamp}</lastmod>
    <changefreq>${url.changefreq}</changefreq>
    <priority>${url.priority}</priority>
  </url>`).join("");

        const dynamicXml = `
  <!-- Dynamic URLs (akan di-generate per konten user) -->
  <url>
    <loc>${baseUrl}${config.dynamicConfig.urlPattern}</loc>
    <lastmod>{${config.dynamicConfig.lastmodSource || "timestamp"}}</lastmod>
    <changefreq>{${config.dynamicConfig.changefreqSource || "weekly"}}</changefreq>
    <priority>{${config.dynamicConfig.prioritySource || "0.7"}}</priority>
    <image:image>
      <image:loc>${baseUrl}/{${config.dynamicConfig.imageSource || "image_url"}}</image:loc>
    </image:image>
  </url>`;

        const fullXml = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
        xmlns:image="http://www.google.com/schemas/sitemap-image/1.1">
${staticXml}
${dynamicXml}
</urlset>`;

        setSitemapXml(fullXml);
    }, [config, baseUrl]);

    const updateConfig = (updates: Partial<SitemapConfig>) => {
        const newConfig = { ...config, ...updates };
        setConfig(newConfig);
        onChange?.(newConfig);
    };

    const addStaticUrl = () => {
        const newUrl = prompt("Masukkan path URL (contoh: /about)");
        if (newUrl) {
            updateConfig({
                staticUrls: [...config.staticUrls, { loc: newUrl, priority: 0.5, changefreq: "monthly" }],
            });
        }
    };

    const removeStaticUrl = (index: number) => {
        const newUrls = [...config.staticUrls];
        newUrls.splice(index, 1);
        updateConfig({ staticUrls: newUrls });
    };

    const updateStaticUrl = (index: number, field: string, value: any) => {
        const newUrls = [...config.staticUrls];
        newUrls[index] = { ...newUrls[index], [field]: value };
        updateConfig({ staticUrls: newUrls });
    };

    const copyToClipboard = () => {
        navigator.clipboard.writeText(sitemapXml);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
    };

    const downloadSitemap = () => {
        const blob = new Blob([sitemapXml], { type: "application/xml" });
        const url = URL.createObjectURL(blob);
        const a = document.createElement("a");
        a.href = url;
        a.download = "sitemap.xml";
        a.click();
        URL.revokeObjectURL(url);
    };

    return (
        <div className="space-y-4">
            <Tabs defaultValue="config" className="w-full">
                <TabsList className="grid w-full grid-cols-3">
                    <TabsTrigger value="config">⚙️ Konfigurasi Sitemap</TabsTrigger>
                    <TabsTrigger value="preview">👁️ Preview XML</TabsTrigger>
                    <TabsTrigger value="fields">📋 Available Fields</TabsTrigger>
                </TabsList>

                {/* Tab Konfigurasi */}
                <TabsContent value="config" className="space-y-4 pt-4">
                    <div className="flex items-center justify-between">
                        <label className="flex items-center gap-2">
                            <input
                                type="checkbox"
                                checked={config.enabled}
                                onChange={(e) => updateConfig({ enabled: e.target.checked })}
                                className="w-4 h-4"
                            />
                            <span className="text-sm">Aktifkan Sitemap</span>
                        </label>
                    </div>

                    <div className="border rounded-lg p-3 space-y-2">
                        <div className="flex justify-between items-center">
                            <h4 className="text-sm font-medium">Static URLs (Tidak Bervariasi)</h4>
                            <Button size="sm" variant="outline" onClick={addStaticUrl} className="h-6 text-xs">
                                + Tambah URL
                            </Button>
                        </div>
                        <div className="space-y-2">
                            {config.staticUrls.map((url, idx) => (
                                <div key={idx} className="flex gap-2 items-center text-xs">
                                    <Input
                                        value={url.loc}
                                        onChange={(e) => updateStaticUrl(idx, "loc", e.target.value)}
                                        className="h-7 flex-1"
                                        placeholder="/path"
                                    />
                                    <Input
                                        type="number"
                                        step="0.1"
                                        value={url.priority}
                                        onChange={(e) => updateStaticUrl(idx, "priority", parseFloat(e.target.value))}
                                        className="h-7 w-16"
                                    />
                                    <select
                                        value={url.changefreq}
                                        onChange={(e) => updateStaticUrl(idx, "changefreq", e.target.value)}
                                        className="h-7 rounded border px-1 text-xs"
                                    >
                                        <option value="always">always</option>
                                        <option value="hourly">hourly</option>
                                        <option value="daily">daily</option>
                                        <option value="weekly">weekly</option>
                                        <option value="monthly">monthly</option>
                                        <option value="yearly">yearly</option>
                                        <option value="never">never</option>
                                    </select>
                                    <Button size="sm" variant="ghost" onClick={() => removeStaticUrl(idx)} className="h-7 px-2">
                                        ❌
                                    </Button>
                                </div>
                            ))}
                        </div>
                    </div>

                    <div className="border rounded-lg p-3 space-y-2">
                        <h4 className="text-sm font-medium">Dynamic URLs (Bervariasi per Konten)</h4>
                        <div className="space-y-2">
                            <div>
                                <label className="text-xs text-slate-500">URL Pattern</label>
                                <Input
                                    value={config.dynamicConfig.urlPattern}
                                    onChange={(e) => updateConfig({
                                        dynamicConfig: { ...config.dynamicConfig, urlPattern: e.target.value },
                                    })}
                                    className="h-8 text-xs"
                                    placeholder="/p/{id} atau /article/{slug}"
                                />
                                <p className="text-xs text-slate-400 mt-1">
                                    Gunakan {'{id}'}, {'{slug}'}, atau field lain dari JSON
                                </p>
                            </div>
                            <div className="grid grid-cols-2 gap-2">
                                <div>
                                    <label className="text-xs text-slate-500">Priority Source</label>
                                    <select
                                        value={config.dynamicConfig.prioritySource}
                                        onChange={(e) => updateConfig({
                                            dynamicConfig: { ...config.dynamicConfig, prioritySource: e.target.value },
                                        })}
                                        className="w-full h-8 rounded border px-2 text-xs"
                                    >
                                        <option value="0.7">Default: 0.7</option>
                                        {availableFields.map(f => (
                                            <option key={f} value={f}>Field: {f}</option>
                                        ))}
                                    </select>
                                </div>
                                <div>
                                    <label className="text-xs text-slate-500">Changefreq Source</label>
                                    <select
                                        value={config.dynamicConfig.changefreqSource}
                                        onChange={(e) => updateConfig({
                                            dynamicConfig: { ...config.dynamicConfig, changefreqSource: e.target.value },
                                        })}
                                        className="w-full h-8 rounded border px-2 text-xs"
                                    >
                                        <option value="weekly">Default: weekly</option>
                                        {availableFields.map(f => (
                                            <option key={f} value={f}>Field: {f}</option>
                                        ))}
                                    </select>
                                </div>
                                <div>
                                    <label className="text-xs text-slate-500">Image Source</label>
                                    <select
                                        value={config.dynamicConfig.imageSource}
                                        onChange={(e) => updateConfig({
                                            dynamicConfig: { ...config.dynamicConfig, imageSource: e.target.value },
                                        })}
                                        className="w-full h-8 rounded border px-2 text-xs"
                                    >
                                        <option value="image_url">Default: image_url</option>
                                        {availableFields.map(f => (
                                            <option key={f} value={f}>Field: {f}</option>
                                        ))}
                                    </select>
                                </div>
                                <div>
                                    <label className="text-xs text-slate-500">Lastmod Source</label>
                                    <select
                                        value={config.dynamicConfig.lastmodSource}
                                        onChange={(e) => updateConfig({
                                            dynamicConfig: { ...config.dynamicConfig, lastmodSource: e.target.value },
                                        })}
                                        className="w-full h-8 rounded border px-2 text-xs"
                                    >
                                        <option value="updated_at">Default: updated_at</option>
                                        <option value="created_at">created_at</option>
                                        {availableFields.map(f => (
                                            <option key={f} value={f}>Field: {f}</option>
                                        ))}
                                    </select>
                                </div>
                            </div>
                        </div>
                    </div>
                </TabsContent>

                {/* Tab Preview XML */}
                <TabsContent value="preview" className="pt-4">
                    <div className="bg-slate-900 rounded-lg p-3 overflow-auto max-h-96">
                        <pre className="text-green-400 text-xs font-mono whitespace-pre-wrap">{sitemapXml}</pre>
                    </div>
                    <div className="flex gap-2 mt-3">
                        <Button size="sm" variant="outline" onClick={copyToClipboard} className="flex-1">
                            {copied ? <Check className="w-3 h-3 mr-1" /> : <Copy className="w-3 h-3 mr-1" />}
                            {copied ? "Tersalin!" : "Salin XML"}
                        </Button>
                        <Button size="sm" variant="outline" onClick={downloadSitemap} className="flex-1">
                            <Download className="w-3 h-3 mr-1" />
                            Download sitemap.xml
                        </Button>
                    </div>
                </TabsContent>

                {/* Tab Available Fields */}
                <TabsContent value="fields" className="pt-4">
                    <div className="bg-slate-100 dark:bg-slate-800 rounded-lg p-3">
                        <h4 className="text-xs font-medium mb-2">Field yang tersedia dari JSON Builder:</h4>
                        <div className="flex flex-wrap gap-1">
                            {availableFields.map(f => (
                                <span key={f} className="bg-blue-100 dark:bg-blue-900 px-2 py-1 rounded text-xs font-mono">
                                    {f}
                                </span>
                            ))}
                        </div>
                    </div>
                </TabsContent>
            </Tabs>
        </div>
    );
}