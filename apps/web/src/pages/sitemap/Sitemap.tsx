// pages/sitemap/index.tsx
import { useState, useEffect } from "react";
import { useSitemap } from "@/hooks/useSitemap";
import {
    Button,
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
    Input,
    Label,
    Switch,
    Skeleton,
    Alert,
    AlertDescription,
} from "@kana-consultant/ui-kit";
import {
    Download,
    Copy,
    Check,
    Loader2,
    AlertCircle,
    FileText,
    Globe,
    Image,
    RefreshCw,
    Clock,
    Link,
} from "lucide-react";
import { useToast } from "@/hooks/use-toast";
import { cn } from "@/lib/utils";

export default function SitemapPage() {
    const {
        isLoading,
        error,
        sitemapData,
        history,
        isHistoryLoading,
        generateSitemap,
        downloadSitemap,
        fetchHistory,
        clearError,
    } = useSitemap();

    const toast = useToast();

    const [baseURL, setBaseURL] = useState("");
    const [includeImages, setIncludeImages] = useState(true);
    const [limit, setLimit] = useState<string>("0");
    const [copied, setCopied] = useState(false);

    // Load history on mount
    useEffect(() => {
        fetchHistory();
    }, [fetchHistory]);

    // Handle generate
    const handleGenerate = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!baseURL.trim()) {
            toast.error("Base URL is required");
            return;
        }

        await generateSitemap({
            baseURL: baseURL.trim(),
            includeImages,
            limit: parseInt(limit) || 0,
        });

        toast.success("Sitemap generated successfully!");
    };

    // Handle copy URL
    const handleCopyURL = () => {
        if (!sitemapData) return;

        const sitemapURL = `${window.location.origin}/sitemap?base_url=${encodeURIComponent(baseURL)}`;
        navigator.clipboard.writeText(sitemapURL);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
        toast.success("Sitemap URL copied!");
    };

    // Handle download
    const handleDownload = () => {
        downloadSitemap();
        toast.success("Sitemap downloaded!");
    };

    // Handle regenerate
    const handleRegenerate = () => {
        handleGenerate(new Event("submit") as any);
    };

    return (
        <div className="container mx-auto py-6 space-y-6 max-w-5xl">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-2xl font-bold text-slate-900 dark:text-white">
                        Sitemap Generator
                    </h1>
                    <p className="text-sm text-slate-500 dark:text-slate-400 mt-1">
                        Generate and manage sitemap XML for your website
                    </p>
                </div>
            </div>

            {/* Error Alert */}
            {error && (
                <Alert tone="warning" className="mb-4">
                    <AlertCircle className="h-4 w-4" />
                    <AlertDescription>{error}</AlertDescription>
                    <button
                        onClick={clearError}
                        className="ml-auto text-sm font-medium underline-offset-2 hover:underline"
                    >
                        Dismiss
                    </button>
                </Alert>
            )}

            {/* Form Card */}
            <Card>
                <CardHeader>
                    <CardTitle className="text-lg flex items-center gap-2">
                        <Globe className="h-5 w-5 text-green-600 dark:text-green-400" />
                        Sitemap Configuration
                    </CardTitle>
                    <CardDescription>
                        Enter your website base URL and configure sitemap options
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <form onSubmit={handleGenerate} className="space-y-4">
                        {/* Base URL */}
                        <div className="space-y-2">
                            <Label htmlFor="baseURL">Base URL</Label>
                            <div className="relative">
                                <Globe className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400" />
                                <Input
                                    id="baseURL"
                                    type="url"
                                    placeholder="https://your-website.com"
                                    value={baseURL}
                                    onChange={(e) => setBaseURL(e.target.value)}
                                    className="pl-9"
                                    required
                                />
                            </div>
                            <p className="text-xs text-slate-400">
                                Example: https://client.com (without trailing slash)
                            </p>
                        </div>

                        {/* Options Row */}
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            {/* Include Images */}
                            <div className="flex items-center space-x-3 rounded-lg border border-slate-200 dark:border-slate-700 p-3">
                                <Switch
                                    id="includeImages"
                                    checked={includeImages}
                                    onCheckedChange={setIncludeImages}
                                />
                                <div className="flex items-center gap-2">
                                    <Image className="h-4 w-4 text-slate-400" />
                                    <Label htmlFor="includeImages" className="cursor-pointer">
                                        Include Images
                                    </Label>
                                </div>
                            </div>

                            {/* Limit */}
                            <div className="flex items-center space-x-3 rounded-lg border border-slate-200 dark:border-slate-700 p-3">
                                <div className="flex items-center gap-2 flex-1">
                                    <FileText className="h-4 w-4 text-slate-400" />
                                    <Label htmlFor="limit" className="cursor-pointer">
                                        Limit
                                    </Label>
                                </div>
                                <Input
                                    id="limit"
                                    type="number"
                                    min={0}
                                    placeholder="0 = all"
                                    value={limit}
                                    onChange={(e) => setLimit(e.target.value)}
                                    className="w-20"
                                />
                            </div>
                        </div>

                        {/* Action Buttons */}
                        <div className="flex flex-wrap gap-3 pt-2">
                            <Button
                                type="submit"
                                disabled={isLoading}
                                className="bg-green-600 hover:bg-green-700 dark:bg-purple-600 dark:hover:bg-purple-700"
                            >
                                {isLoading ? (
                                    <>
                                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                        Generating...
                                    </>
                                ) : (
                                    <>
                                        <RefreshCw className="mr-2 h-4 w-4" />
                                        Generate Sitemap
                                    </>
                                )}
                            </Button>
                        </div>
                    </form>
                </CardContent>
            </Card>

            {/* Result Card */}
            {sitemapData && (
                <Card>
                    <CardHeader>
                        <CardTitle className="text-lg flex items-center gap-2">
                            <FileText className="h-5 w-5 text-green-600 dark:text-green-400" />
                            Sitemap Generated
                        </CardTitle>
                        <CardDescription>
                            Your sitemap has been generated successfully
                        </CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        {/* Stats */}
                        <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
                            <div className="rounded-lg bg-slate-50 dark:bg-white/[0.03] p-4 text-center">
                                <p className="text-2xl font-bold text-slate-900 dark:text-white">
                                    {sitemapData.totalURLs}
                                </p>
                                <p className="text-xs text-slate-500 dark:text-slate-400">Total URLs</p>
                            </div>
                            <div className="rounded-lg bg-slate-50 dark:bg-white/[0.03] p-4 text-center">
                                <p className="text-2xl font-bold text-slate-900 dark:text-white">
                                    {new Date(sitemapData.generatedAt).toLocaleDateString()}
                                </p>
                                <p className="text-xs text-slate-500 dark:text-slate-400">Generated At</p>
                            </div>
                            <div className="rounded-lg bg-slate-50 dark:bg-white/[0.03] p-4 text-center">
                                <p className="text-2xl font-bold text-slate-900 dark:text-white">
                                    {includeImages ? "✅ Yes" : "❌ No"}
                                </p>
                                <p className="text-xs text-slate-500 dark:text-slate-400">Images Included</p>
                            </div>
                        </div>

                        {/* Actions */}
                        <div className="flex flex-wrap gap-3 pt-2 border-t border-slate-200 dark:border-slate-700">
                            <Button
                                variant="outline"
                                onClick={handleDownload}
                                className="gap-2"
                            >
                                <Download className="h-4 w-4" />
                                Download Sitemap
                            </Button>

                            <Button
                                variant="outline"
                                onClick={handleCopyURL}
                                className="gap-2"
                            >
                                {copied ? (
                                    <>
                                        <Check className="h-4 w-4 text-green-600" />
                                        Copied!
                                    </>
                                ) : (
                                    <>
                                        <Copy className="h-4 w-4" />
                                        Copy URL
                                    </>
                                )}
                            </Button>

                            <Button
                                variant="outline"
                                onClick={handleRegenerate}
                                disabled={isLoading}
                                className="gap-2"
                            >
                                <RefreshCw className={cn("h-4 w-4", `${isLoading && "animate-spin"}`)} />
                                Regenerate
                            </Button>
                        </div>

                        {/* Preview */}
                        <div className="mt-4">
                            <p className="text-sm font-medium text-slate-700 dark:text-slate-300 mb-2">
                                Sitemap URL:
                            </p>
                            <code className="block w-full rounded-lg bg-slate-100 dark:bg-slate-800 p-3 text-xs text-slate-600 dark:text-slate-300 break-all">
                                {`${window.location.origin}/sitemap?base_url=${encodeURIComponent(baseURL)}&include_images=${includeImages}`}
                            </code>
                        </div>
                    </CardContent>
                </Card>
            )}

            {/* History Card */}
            <Card>
                <CardHeader>
                    <CardTitle className="text-lg flex items-center gap-2">
                        <Clock className="h-5 w-5 text-slate-500" />
                        Sitemap History
                    </CardTitle>
                    <CardDescription>
                        Recent sitemap generation history
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    {isHistoryLoading ? (
                        <div className="space-y-3">
                            <Skeleton className="h-12 w-full" />
                            <Skeleton className="h-12 w-full" />
                            <Skeleton className="h-12 w-full" />
                        </div>
                    ) : history.length === 0 ? (
                        <div className="text-center py-8 text-slate-500">
                            <FileText className="h-12 w-12 mx-auto text-slate-300 mb-3" />
                            <p>No sitemap history yet</p>
                            <p className="text-sm">Generate your first sitemap above</p>
                        </div>
                    ) : (
                        <div className="overflow-x-auto">
                            <table className="w-full text-sm">
                                <thead>
                                    <tr className="border-b border-slate-200 dark:border-slate-700">
                                        <th className="text-left py-3 px-4 font-medium text-slate-500">Date</th>
                                        <th className="text-left py-3 px-4 font-medium text-slate-500">URLs</th>
                                        <th className="text-left py-3 px-4 font-medium text-slate-500">Status</th>
                                        <th className="text-right py-3 px-4 font-medium text-slate-500">Action</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {history.map((item) => (
                                        <tr
                                            key={item.id}
                                            className="border-b border-slate-100 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-white/[0.02]"
                                        >
                                            <td className="py-3 px-4 text-slate-700 dark:text-slate-300">
                                                {new Date(item.createdAt).toLocaleString()}
                                            </td>
                                            <td className="py-3 px-4 text-slate-700 dark:text-slate-300">
                                                {item.totalURLs}
                                            </td>
                                            <td className="py-3 px-4">
                                                <span
                                                    className={cn(
                                                        "inline-flex items-center px-2 py-1 rounded-full text-xs font-medium",
                                                        item.status === "success"
                                                            ? "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400"
                                                            : item.status === "failed"
                                                                ? "bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400"
                                                                : "bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400"
                                                    )}
                                                >
                                                    {item.status === "success" && <Check className="h-3 w-3 mr-1" />}
                                                    {item.status === "failed" && <AlertCircle className="h-3 w-3 mr-1" />}
                                                    {item.status}
                                                </span>
                                            </td>
                                            <td className="py-3 px-4 text-right">
                                                {item.sitemapURL && (
                                                    <Button
                                                        variant="ghost"
                                                        size="sm"
                                                        onClick={() => window.open(item.sitemapURL, "_blank")}
                                                        className="gap-1"
                                                    >
                                                        <Link className="h-3 w-3" />
                                                        View
                                                    </Button>
                                                )}
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                    )}
                </CardContent>
            </Card>
        </div>
    );
}