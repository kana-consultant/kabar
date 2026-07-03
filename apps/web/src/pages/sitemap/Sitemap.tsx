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
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
    Badge,
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger,
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
    ExternalLink,
    Package,
    Pencil,
    Info,
    HelpCircle,
} from "lucide-react";
import { useToast } from "@/hooks/use-toast";
import { cn } from "@/lib/utils";

// Placeholder options
const PLACEHOLDER_OPTIONS = [
    { value: "{slug}", label: "{slug} - URL-friendly product name" },
    { value: "{title}", label: "{title} - Product title" },
    { value: "{id}", label: "{id} - Product ID" },
    { value: "{category}", label: "{category} - Product category" },
    { value: "{sku}", label: "{sku} - Product SKU" },
    { value: "{date}", label: "{date} - Current date (YYYY-MM-DD)" },
    { value: "{timestamp}", label: "{timestamp} - Unix timestamp" },
];

export default function SitemapPage() {
    const {
        isLoading,
        isProductsLoading,
        error,
        sitemapData,
        products,
        generateSitemap,
        downloadSitemap,
        clearError,
    } = useSitemap();

    const toast = useToast();

    const [selectedProductId, setSelectedProductId] = useState<string>("");
    const [baseURL, setBaseURL] = useState("");
    const [includeImages, setIncludeImages] = useState(true);
    const [limit, setLimit] = useState<string>("0");
    const [copied, setCopied] = useState(false);
    const [isBaseURLManuallyEdited, setIsBaseURLManuallyEdited] = useState(false);
    const [showPlaceholderHelp, setShowPlaceholderHelp] = useState(false);


    useEffect(() => {
        if (selectedProductId && products.length > 0 && !isBaseURLManuallyEdited) {
            const selected = products.find(p => p.id === selectedProductId);
            if (selected) {
                // Gunakan template default dengan {slug}
                const url = selected.api_endpoint || 
                           selected.domain || 
                           `https://${selected.name.toLowerCase().replace(/\s+/g, '-')}.com/{slug}`;
                setBaseURL(url);
            }
        }
    }, [selectedProductId, products, isBaseURLManuallyEdited]);

    useEffect(() => {
        setIsBaseURLManuallyEdited(false);
    }, [selectedProductId]);

    const handleGenerate = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!selectedProductId) {
            toast.error("Please select a product");
            return;
        }

        if (!baseURL.trim()) {
            toast.error("Base URL is required");
            return;
        }

        // Validasi placeholder
        const hasPlaceholder = PLACEHOLDER_OPTIONS.some(opt => baseURL.includes(opt.value));
        if (!hasPlaceholder) {
            toast.warning("No placeholder detected. Consider using {slug} or {title} for dynamic URLs");
        }

        await generateSitemap({
            productId: selectedProductId,
            baseURL: baseURL.trim(),
            includeImages,
            limit: parseInt(limit) || 0,
        });

        toast.success("Sitemap generated successfully!");
    };

    const handleBaseURLChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setBaseURL(e.target.value);
        setIsBaseURLManuallyEdited(true);
    };

    const handleCopyURL = () => {
        if (!sitemapData) return;

        const sitemapURL = `${window.location.origin}/sitemap?product_id=${selectedProductId}&base_url=${encodeURIComponent(baseURL)}&include_images=${includeImages}`;
        navigator.clipboard.writeText(sitemapURL);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
        toast.success("Sitemap URL copied!");
    };

    const handleDownload = () => {
        downloadSitemap();
        toast.success("Sitemap downloaded!");
    };

    const handleRegenerate = () => {
        handleGenerate(new Event("submit") as any);
    };

    const handleResetToDefault = () => {
        setIsBaseURLManuallyEdited(false);
        if (selectedProductId) {
            const selected = products.find(p => p.id === selectedProductId);
            if (selected) {
                const url = selected.api_endpoint || 
                           selected.domain || 
                           `https://${selected.name.toLowerCase().replace(/\s+/g, '-')}.com/{slug}`;
                setBaseURL(url);
                toast.info("Base URL reset to product default");
            }
        }
    };

    const insertPlaceholder = (placeholder: string) => {
        const input = document.getElementById('baseURL') as HTMLInputElement;
        if (input) {
            const start = input.selectionStart || 0;
            const end = input.selectionEnd || 0;
            const value = baseURL;
            const newValue = value.substring(0, start) + placeholder + value.substring(end);
            setBaseURL(newValue);
            setIsBaseURLManuallyEdited(true);
            
            // Set cursor position after inserted placeholder
            setTimeout(() => {
                input.focus();
                const newPosition = start + placeholder.length;
                input.setSelectionRange(newPosition, newPosition);
            }, 0);
        } else {
            // Fallback: append to end
            setBaseURL(baseURL + placeholder);
            setIsBaseURLManuallyEdited(true);
        }
    };

    const getPlaceholderPreview = () => {
        if (!baseURL) return null;
        
        const preview = baseURL
            .replace(/{slug}/g, "product-name")
            .replace(/{title}/g, "Product Title")
            .replace(/{id}/g, "123")
            .replace(/{category}/g, "electronics")
            .replace(/{sku}/g, "SKU-001")
            .replace(/{date}/g, new Date().toISOString().split('T')[0])
            .replace(/{timestamp}/g, Date.now().toString());
        
        return preview;
    };

    const getStatusColor = (status: string) => {
        switch (status) {
            case "success":
                return "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400";
            case "failed":
                return "bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400";
            default:
                return "bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400";
        }
    };

    const getStatusIcon = (status: string) => {
        switch (status) {
            case "success":
                return <Check className="h-3 w-3 mr-1" />;
            case "failed":
                return <AlertCircle className="h-3 w-3 mr-1" />;
            default:
                return <Loader2 className="h-3 w-3 mr-1 animate-spin" />;
        }
    };

    return (
        <div className="container space-y-6 max-w-full">
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
                        Select a product and configure sitemap options
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <form onSubmit={handleGenerate} className="space-y-4">
                        {/* Select Product */}
                        <div className="space-y-2">
                            <Label htmlFor="product">Select Product</Label>
                            {isProductsLoading ? (
                                <Skeleton className="h-10 w-full" />
                            ) : (
                                <Select
                                    value={selectedProductId}
                                    onValueChange={setSelectedProductId}
                                >
                                    <SelectTrigger id="product" className="w-full">
                                        <SelectValue placeholder="Select a product..." />
                                    </SelectTrigger>
                                    <SelectContent>
                                        {products.map((product) => (
                                            <SelectItem key={product.id} value={product.id}>
                                                <div className="flex items-center gap-2">
                                                    <Package className="h-4 w-4" />
                                                    <span>{product.name}</span>
                                                    {product.domain && (
                                                        <span className="text-xs text-slate-400">
                                                            ({product.domain})
                                                        </span>
                                                    )}
                                                </div>
                                            </SelectItem>
                                        ))}
                                    </SelectContent>
                                </Select>
                            )}
                            <p className="text-xs text-slate-400">
                                Select a product to generate sitemap for its content
                            </p>
                        </div>

                        {/* Base URL */}
                        <div className="space-y-2">
                            <div className="flex items-center justify-between">
                                <Label htmlFor="baseURL" className="flex items-center gap-2">
                                    Base URL Template
                                    <TooltipProvider>
                                        <Tooltip>
                                            <TooltipTrigger asChild>
                                                <HelpCircle className="h-4 w-4 text-slate-400 cursor-help" />
                                            </TooltipTrigger>
                                            <TooltipContent className="max-w-xs">
                                                <p>Use placeholders for dynamic URLs. Example: https://domain.com/{'{slug}'}</p>
                                            </TooltipContent>
                                        </Tooltip>
                                    </TooltipProvider>
                                </Label>
                                {isBaseURLManuallyEdited && selectedProductId && (
                                    <Button
                                        type="button"
                                        variant="ghost"
                                        size="sm"
                                        onClick={handleResetToDefault}
                                        className="h-6 text-xs gap-1 text-slate-500 hover:text-slate-700"
                                    >
                                        <RefreshCw className="h-3 w-3" />
                                        Reset to default
                                    </Button>
                                )}
                            </div>
                            <div className="relative">
                                <Globe className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400" />
                                <Input
                                    id="baseURL"
                                    type="text"
                                    placeholder="https://your-website.com/{slug}"
                                    value={baseURL}
                                    onChange={handleBaseURLChange}
                                    className={cn(
                                        "pl-9 font-mono text-sm",
                                        `${isBaseURLManuallyEdited && "border-amber-400 dark:border-amber-500"}`
                                    )}
                                    required
                                />
                                {isBaseURLManuallyEdited && (
                                    <div className="absolute right-3 top-1/2 -translate-y-1/2">
                                        <Pencil className="h-4 w-4 text-amber-500" />
                                    </div>
                                )}
                            </div>
                            
                            {/* Placeholder Preview */}
                            {baseURL && (
                                <div className="mt-2 p-3 rounded-lg bg-slate-50 dark:bg-white/[0.03] border border-slate-200 dark:border-slate-700">
                                    <div className="flex items-center justify-between mb-2">
                                        <span className="text-xs font-medium text-slate-600 dark:text-slate-400 flex items-center gap-1">
                                            <Info className="h-3 w-3" />
                                            Preview
                                        </span>
                                        <Button
                                            type="button"
                                            variant="ghost"
                                            size="sm"
                                            onClick={() => setShowPlaceholderHelp(!showPlaceholderHelp)}
                                            className="h-6 text-xs"
                                        >
                                            {showPlaceholderHelp ? "Hide placeholders" : "Show placeholders"}
                                        </Button>
                                    </div>
                                    <code className="text-xs text-slate-700 dark:text-slate-300 break-all">
                                        {getPlaceholderPreview()}
                                    </code>
                                </div>
                            )}

                            {/* Placeholder Buttons */}
                            <div className="flex flex-wrap gap-2 mt-2">
                                {PLACEHOLDER_OPTIONS.map((opt) => (
                                    <Button
                                        key={opt.value}
                                        type="button"
                                        variant="outline"
                                        size="sm"
                                        onClick={() => insertPlaceholder(opt.value)}
                                        className="h-7 text-xs font-mono gap-1"
                                    >
                                        {opt.value}
                                        <Badge tone="neutral" className="text-[10px] h-4">
                                            {opt.label.split(' - ')[1] || 'dynamic'}
                                        </Badge>
                                    </Button>
                                ))}
                            </div>

                            {/* Placeholder Help */}
                            {showPlaceholderHelp && (
                                <div className="mt-2 p-3 rounded-lg bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800">
                                    <h4 className="text-sm font-medium text-blue-800 dark:text-blue-300 mb-2">
                                        Available Placeholders
                                    </h4>
                                    <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
                                        {PLACEHOLDER_OPTIONS.map((opt) => (
                                            <div key={opt.value} className="flex items-start gap-2 text-xs">
                                                <code className="px-1.5 py-0.5 bg-blue-100 dark:bg-blue-800 rounded text-blue-800 dark:text-blue-300 font-mono">
                                                    {opt.value}
                                                </code>
                                                <span className="text-blue-700 dark:text-blue-400">
                                                    {opt.label.split(' - ')[1] || opt.label}
                                                </span>
                                            </div>
                                        ))}
                                    </div>
                                    <p className="text-xs text-blue-700 dark:text-blue-400 mt-2">
                                        ℹ️ Placeholders will be replaced with actual data when generating the sitemap
                                    </p>
                                </div>
                            )}

                            <p className="text-xs text-slate-400">
                                {selectedProductId && !isBaseURLManuallyEdited
                                    ? "Auto-filled from selected product (you can edit and use placeholders)"
                                    : isBaseURLManuallyEdited
                                        ? "✏️ Manually edited - use placeholders for dynamic content"
                                        : "Use placeholders like {slug} or {title} for dynamic URLs"}
                            </p>
                        </div>

                        {/* Options */}
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
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

                        {/* Generate Button */}
                        <div className="flex flex-wrap gap-3 pt-2">
                            <Button
                                type="submit"
                                disabled={isLoading || !selectedProductId}
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

                        <div className="rounded-lg bg-slate-50 dark:bg-white/[0.03] p-3 text-sm">
                            <div className="flex items-center gap-2">
                                <Package className="h-4 w-4 text-slate-400" />
                                <span className="font-medium text-slate-700 dark:text-slate-300">Product ID:</span>
                                <span className="text-slate-600 dark:text-slate-400">
                                    {sitemapData.productId || selectedProductId}
                                </span>
                            </div>
                            <div className="flex items-center gap-2 mt-1">
                                <Globe className="h-4 w-4 text-slate-400" />
                                <span className="font-medium text-slate-700 dark:text-slate-300">Base URL Template:</span>
                                <span className="text-slate-600 dark:text-slate-400 break-all font-mono">
                                    {sitemapData.baseURL || baseURL}
                                </span>
                            </div>
                            {baseURL && baseURL.includes('{') && (
                                <div className="flex items-center gap-2 mt-1 text-xs text-slate-500">
                                    <Info className="h-3 w-3" />
                                    <span>Placeholders will be replaced with actual data</span>
                                </div>
                            )}
                        </div>

                        <div className="flex flex-wrap gap-3 pt-2 border-t border-slate-200 dark:border-slate-700">
                            <Button variant="outline" onClick={handleDownload} className="gap-2">
                                <Download className="h-4 w-4" />
                                Download Sitemap
                            </Button>
                            <Button variant="outline" onClick={handleCopyURL} className="gap-2">
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
                            <Button variant="outline" onClick={handleRegenerate} disabled={isLoading} className="gap-2">
                                <RefreshCw className={cn("h-4 w-4",`${ isLoading && "animate-spin"}`)} />
                                Regenerate
                            </Button>
                        </div>

                        <div className="mt-4">
                            <p className="text-sm font-medium text-slate-700 dark:text-slate-300 mb-2">Sitemap URL:</p>
                            <code className="block w-full rounded-lg bg-slate-100 dark:bg-slate-800 p-3 text-xs text-slate-600 dark:text-slate-300 break-all">
                                {`${window.location.origin}/sitemap?product_id=${selectedProductId}&base_url=${encodeURIComponent(baseURL)}&include_images=${includeImages}`}
                            </code>
                        </div>
                    </CardContent>
                </Card>
            )}
        </div>
    );
}