// components/SitemapManager.tsx
import { useState, useEffect } from "react";
import { Button } from "@kana-consultant/ui-kit";
import { Input } from "@kana-consultant/ui-kit";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@kana-consultant/ui-kit";
import { Download, Copy, Check } from "lucide-react";

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
        lastmod?: string; // ISO 8601: YYYY-MM-DD
    }[];
    dynamicConfig: {
        urlPattern: string;       // contoh: /p/{id}
        prioritySource: string;   // field key untuk priority (0.0–1.0)
        changefreqSource: string; // field key untuk changefreq
        imageSource: string;      // field key untuk gambar (harus absolute URL)
        lastmodSource: string;    // field key untuk last modified (ISO 8601)
    };
}

const CHANGEFREQ_OPTIONS = ["always", "hourly", "daily", "weekly", "monthly", "yearly", "never"] as const;

const clampPriority = (val: number) => Math.min(1.0, Math.max(0.0, val));

export function SitemapManager({ fieldMapping, baseUrl = "https://domainanda.com", onChange }: SitemapManagerProps) {
    const [config, setConfig] = useState<SitemapConfig>({
        enabled: true,
        staticUrls: [
            { loc: "/", priority: 1.0, changefreq: "daily", lastmod: new Date().toISOString().split("T")[0] },
            { loc: "/privacy", priority: 0.3, changefreq: "yearly" },
            { loc: "/terms", priority: 0.3, changefreq: "yearly" },
        ],
        dynamicConfig: {
            urlPattern: "/p/{id}",
            prioritySource: "sitemap_priority",
            changefreqSource: "sitemap_changefreq",
            imageSource: "image_url",       // expected: absolute URL dari field
            lastmodSource: "updated_at",    // expected: ISO 8601 dari field
        },
    });

    const [sitemapXml, setSitemapXml] = useState<string>("");
    const [copied, setCopied] = useState(false);
    const [availableFields, setAvailableFields] = useState<string[]>([]);

    // Extract available fields dari fieldMapping
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

    // Generate sitemap XML preview sesuai standar sitemaps.org
    useEffect(() => {
        const today = new Date().toISOString().split("T")[0]; // YYYY-MM-DD

        const staticXml = config.staticUrls.map(url => {
            const lastmod = url.lastmod || today;
            return `  <url>
    <loc>${baseUrl}${url.loc}</loc>
    <lastmod>${lastmod}</lastmod>
    <changefreq>${url.changefreq}</changefreq>
    <priority>${url.priority.toFixed(1)}</priority>
  </url>`;
        }).join("\n");

        // Dynamic URL: placeholder menggunakan notasi {field} sebagai keterangan,
        // bukan XML literal — ini adalah preview template, bukan XML aktual
        const dynamicXml = `  <!-- Dynamic URL (di-generate per konten, contoh satu entry) -->
  <url>
    <loc>${baseUrl}${config.dynamicConfig.urlPattern.replace("{id}", "123")}</loc>
    <lastmod>{${config.dynamicConfig.lastmodSource}}</lastmod>
    <changefreq>{${config.dynamicConfig.changefreqSource}}</changefreq>
    <priority>{${config.dynamicConfig.prioritySource}}</priority>
    <image:image>
      <!-- image:loc harus berupa absolute URL -->
      <image:loc>{${config.dynamicConfig.imageSource}}</image:loc>
    </image:image>
  </url>`;

        const fullXml = `<?xml version="1.0" encoding="UTF-8"?>
<!-- 
  CATATAN: Bagian dynamic URL menggunakan placeholder {field_name}.
  Nilai aktual akan di-resolve saat runtime dari data konten Anda.
  Pastikan field image mengandung absolute URL (https://...).
  Pastikan field lastmod mengandung format ISO 8601 (YYYY-MM-DD).
-->
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
                staticUrls: [
                    ...config.staticUrls,
                    {
                        loc: newUrl,
                        priority: 0.5,
                        changefreq: "monthly",
                        lastmod: new Date().toISOString().split("T")[0],
                    },
                ],
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
        // Validasi priority agar tetap dalam range 0.0–1.0
        if (field === "priority") {
            value = clampPriority(parseFloat(value) || 0);
        }
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
                    <label className="flex items-center gap-2">
                        <input
                            type="checkbox"
                            checked={config.enabled}
                            onChange={(e: any) => updateConfig({ enabled: e.target.checked })}
                            className="w-4 h-4"
                        />
                        <span className="text-sm">Aktifkan Sitemap</span>
                    </label>

                    {/* Static URLs */}
                    <div className="border rounded-lg p-3 space-y-2">
                        <div className="flex justify-between items-center">
                            <h4 className="text-sm font-medium">Static URLs</h4>
                            <Button size="sm" variant="outline" onClick={addStaticUrl} className="h-6 text-xs">
                                + Tambah URL
                            </Button>
                        </div>
                        {/* Header kolom */}
                        <div className="grid grid-cols-[1fr_64px_100px_80px_32px] gap-2 text-xs text-slate-400 px-1">
                            <span>Path</span>
                            <span>Priority</span>
                            <span>Changefreq</span>
                            <span>Lastmod</span>
                            <span></span>
                        </div>
                        <div className="space-y-2">
                            {config.staticUrls.map((url, idx) => (
                                <div key={idx} className="grid grid-cols-[1fr_64px_100px_80px_32px] gap-2 items-center text-xs">
                                    <Input
                                        value={url.loc}
                                        onChange={(e: any) => updateStaticUrl(idx, "loc", e.target.value)}
                                        className="h-7"
                                        placeholder="/path"
                                    />
                                    <Input
                                        type="number"
                                        step="0.1"
                                        min="0"
                                        max="1"
                                        value={url.priority}
                                        onChange={(e: any) => updateStaticUrl(idx, "priority", e.target.value)}
                                        className="h-7"
                                    />
                                    <select
                                        value={url.changefreq}
                                        onChange={(e: any) => updateStaticUrl(idx, "changefreq", e.target.value)}
                                        className="h-7 rounded border px-1 text-xs"
                                    >
                                        {CHANGEFREQ_OPTIONS.map(opt => (
                                            <option key={opt} value={opt}>{opt}</option>
                                        ))}
                                    </select>
                                    {/* lastmod: ISO 8601 date picker */}
                                    <Input
                                        type="date"
                                        value={url.lastmod || ""}
                                        onChange={(e: any) => updateStaticUrl(idx, "lastmod", e.target.value)}
                                        className="h-7 text-xs"
                                    />
                                    <Button size="sm" variant="ghost" onClick={() => removeStaticUrl(idx)} className="h-7 px-2">
                                        ❌
                                    </Button>
                                </div>
                            ))}
                        </div>
                    </div>

                    {/* Dynamic URL Config */}
                    <div className="border rounded-lg p-3 space-y-3">
                        <h4 className="text-sm font-medium">Dynamic URLs (per Konten)</h4>

                        <div>
                            <label className="text-xs text-slate-500">URL Pattern</label>
                            <Input
                                value={config.dynamicConfig.urlPattern}
                                onChange={(e: any) => updateConfig({
                                    dynamicConfig: { ...config.dynamicConfig, urlPattern: e.target.value },
                                })}
                                className="h-8 text-xs"
                                placeholder="/p/{id} atau /article/{slug}"
                            />
                            <p className="text-xs text-slate-400 mt-1">
                                Gunakan {"{id}"}, {"{slug}"}, atau field lain dari JSON
                            </p>
                        </div>

                        <div className="grid grid-cols-2 gap-2">
                            {/* Priority Source */}
                            <div>
                                <label className="text-xs text-slate-500">
                                    Priority Source
                                    <span className="ml-1 text-slate-400">(nilai: 0.0–1.0)</span>
                                </label>
                                <select
                                    value={config.dynamicConfig.prioritySource}
                                    onChange={(e: any) => updateConfig({
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

                            {/* Changefreq Source */}
                            <div>
                                <label className="text-xs text-slate-500">
                                    Changefreq Source
                                </label>
                                <select
                                    value={config.dynamicConfig.changefreqSource}
                                    onChange={(e: any) => updateConfig({
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

                            {/* Image Source */}
                            <div>
                                <label className="text-xs text-slate-500">
                                    Image Source
                                    <span className="ml-1 text-slate-400">(harus absolute URL)</span>
                                </label>
                                <select
                                    value={config.dynamicConfig.imageSource}
                                    onChange={(e: any) => updateConfig({
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

                            {/* Lastmod Source */}
                            <div>
                                <label className="text-xs text-slate-500">
                                    Lastmod Source
                                    <span className="ml-1 text-slate-400">(format: ISO 8601)</span>
                                </label>
                                <select
                                    value={config.dynamicConfig.lastmodSource}
                                    onChange={(e: any) => updateConfig({
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

                        {/* Info box standar sitemap */}
                        <div className="bg-blue-50 dark:bg-blue-950 rounded p-2 text-xs text-blue-700 dark:text-blue-300 space-y-1">
                            <p>📌 <strong>Standar sitemap.org:</strong></p>
                            <ul className="list-disc list-inside space-y-0.5 text-blue-600 dark:text-blue-400">
                                <li><code>lastmod</code> harus ISO 8601: <code>YYYY-MM-DD</code> atau <code>YYYY-MM-DDTHH:MM:SSZ</code></li>
                                <li><code>priority</code> range <code>0.0</code> hingga <code>1.0</code></li>
                                <li><code>image:loc</code> harus berupa absolute URL (https://...)</li>
                                <li>Max 50.000 URL per file sitemap</li>
                            </ul>
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
                        <h4 className="text-xs font-medium mb-2">Field tersedia dari JSON Builder:</h4>
                        <div className="flex flex-wrap gap-1">
                            {availableFields.length === 0 && (
                                <span className="text-xs text-slate-400">Tidak ada field tersedia.</span>
                            )}
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