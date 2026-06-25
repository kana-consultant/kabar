import { useEffect, useState, useRef, useCallback } from "react";
import { Input } from "@kana-consultant/ui-kit";
import { Label } from "@kana-consultant/ui-kit";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@kana-consultant/ui-kit";
import {
    Card,
    CardContent,
    CardHeader,
    CardTitle
} from "@kana-consultant/ui-kit";
import { Button } from "@kana-consultant/ui-kit";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from "@kana-consultant/ui-kit";

import type { AdapterConfig, Product } from "@/services/product";

interface ProductApiConfigProps {
    config: Partial<AdapterConfig>;
    onUpdate: (updates: Partial<AdapterConfig>) => void;
    onUpdateProduct: (updates: Partial<Product>) => void;
}

interface HeaderItem {
    key: string;
    value: string;
}

type AuthType = 'none' | 'apiKey';
type AuthPrefix = 'Bearer' | 'ApiKey' | 'Basic' | 'custom';

const PLACEHOLDER_OPTIONS = [
    { label: '{{id}}', value: '{{id}}' },
    { label: '{{name}}', value: '{{name}}' },
    { label: '{{sku}}', value: '{{sku}}' },
    { label: '{{price}}', value: '{{price}}' },
    { label: '{{stock}}', value: '{{stock}}' },
    { label: '{{api_key}}', value: '{{api_key}}' },
];

const AUTH_PREFIX_OPTIONS = [
    { label: 'Bearer', value: 'Bearer' },
    { label: 'ApiKey', value: 'ApiKey' },
    { label: 'Basic', value: 'Basic' },
    { label: 'Custom', value: 'custom' },
];

export function ProductApiConfig({
    config,
    onUpdate,
    onUpdateProduct,
}: ProductApiConfigProps) {

    const [headers, setHeaders] = useState<HeaderItem[]>([]);
    const [authType, setAuthType] = useState<AuthType>('apiKey');
    const [authPrefix, setAuthPrefix] = useState<AuthPrefix>('ApiKey');
    const [apiKeyValue, setApiKeyValue] = useState<string>('');
    const [customAuthPrefix, setCustomAuthPrefix] = useState<string>('');

    const [openPopoverIndex, setOpenPopoverIndex] = useState<number | null>(null);
    const internalUpdateRef = useRef(false);
    const isInitializedRef = useRef(false);

    useEffect(() => {
        if (isInitializedRef.current) return;

        let obj: Record<string, string> = {};

        if (config.custom_headers) {
            if (typeof config.custom_headers === 'string') {
                try {
                    obj = JSON.parse(config.custom_headers);
                } catch (e) {
                    obj = {};
                }
            } else if (typeof config.custom_headers === 'object') {
                obj = config.custom_headers as Record<string, string>;
            }
        }

        const authHeader = obj['Authorization'] || obj['authorization'] || '';

        if (authHeader) {
            if (authHeader.toLowerCase().startsWith('bearer ')) {
                setAuthType('apiKey');
                setAuthPrefix('Bearer');
            } else if (authHeader.toLowerCase().startsWith('apikey ')) {
                setAuthType('apiKey');
                setAuthPrefix('ApiKey');
            } else if (authHeader.toLowerCase().startsWith('basic ')) {
                setAuthType('apiKey');
                setAuthPrefix('Basic');
            } else {
                const spaceIndex = authHeader.indexOf(' ');
                if (spaceIndex > 0) {
                    setAuthType('apiKey');
                    setAuthPrefix('custom');
                    const prefix = authHeader.substring(0, spaceIndex);
                    setCustomAuthPrefix(prefix);
                }
            }
        } else {
            setAuthType('apiKey');
            setAuthPrefix('ApiKey');
        }

     

        const regularHeaders: Record<string, string> = {};
        Object.entries(obj).forEach(([key, value]) => {
            if (key.toLowerCase() !== 'authorization') {
                regularHeaders[key] = value;
            }
        });

        const mapped: HeaderItem[] = Object.entries(regularHeaders).map(
            ([key, value]) => ({ key, value: String(value) })
        );

        setHeaders(mapped.length ? mapped : [{ key: "", value: "" }]);
        isInitializedRef.current = true;

    }, []);

    const buildAuthValue = useCallback((prefix: AuthPrefix, customPrefix: string) => {
        if (prefix === 'custom') {
            return `${customPrefix} {{api_key}}`;
        }
        return `${prefix} {{api_key}}`;
    }, []);

    const updateFullConfig = useCallback((
        currentHeaders: HeaderItem[],
        currentAuthType: AuthType,
        currentAuthPrefix: AuthPrefix,
        currentCustomPrefix: string,
        currentApiKey: string,
    ) => {
        internalUpdateRef.current = true;

        const allHeaders: Record<string, string> = {};

        // Authorization header selalu {{api_key}}
        if (currentAuthType !== 'none') {
            allHeaders['Authorization'] = buildAuthValue(currentAuthPrefix, currentCustomPrefix);
        }

        currentHeaders.forEach(item => {
            if (item.key.trim() && item.key.toLowerCase() !== 'authorization') {
                allHeaders[item.key.trim()] = item.value;
            }
        });

        // Kirim custom_headers dan api_key terpisah
        onUpdate({
            custom_headers: JSON.stringify(allHeaders),
        });

        onUpdateProduct({
            api_key: currentApiKey
        });

        setTimeout(() => {
            internalUpdateRef.current = false;
        }, 0);
    }, [buildAuthValue, onUpdate]);

    const updateHeaders = useCallback((updated: HeaderItem[]) => {
        setHeaders(updated);
        updateFullConfig(updated, authType, authPrefix, customAuthPrefix, apiKeyValue);
    }, [authType, authPrefix, customAuthPrefix, apiKeyValue, updateFullConfig]);

    const updateAuthType = useCallback((newType: AuthType) => {
        setAuthType(newType);
        updateFullConfig(headers, newType, authPrefix, customAuthPrefix, apiKeyValue);
    }, [headers, authPrefix, customAuthPrefix, apiKeyValue, updateFullConfig]);

    const updateAuthPrefix = useCallback((newPrefix: AuthPrefix) => {
        setAuthPrefix(newPrefix);
        if (newPrefix !== 'custom') {
            setCustomAuthPrefix('');
        }
        updateFullConfig(headers, authType, newPrefix, newPrefix === 'custom' ? customAuthPrefix : '', apiKeyValue);
    }, [headers, authType, customAuthPrefix, apiKeyValue, updateFullConfig]);

    const updateApiKeyValue = useCallback((value: string) => {
        setApiKeyValue(value);
        updateFullConfig(headers, authType, authPrefix, customAuthPrefix, value);
    }, [headers, authType, authPrefix, customAuthPrefix, updateFullConfig]);

    const updateCustomAuthPrefix = useCallback((prefix: string) => {
        setCustomAuthPrefix(prefix);
        updateFullConfig(headers, authType, 'custom', prefix, apiKeyValue);
    }, [headers, authType, apiKeyValue, updateFullConfig]);

    const addHeader = useCallback(() => {
        updateHeaders([...headers, { key: "", value: "" }]);
    }, [headers, updateHeaders]);

    const removeHeader = useCallback((index: number) => {
        const updated = headers.filter((_, i) => i !== index);
        updateHeaders(updated.length ? updated : [{ key: "", value: "" }]);
    }, [headers, updateHeaders]);

    const handlePlaceholderSelect = useCallback((index: number, placeholder: string) => {
        const updated = [...headers];
        updated[index].value = placeholder;
        updateHeaders(updated);
        setOpenPopoverIndex(null);
    }, [headers, updateHeaders]);

    return (
        <Card className="shadow-sm">
            <CardHeader className="pb-3">
                <CardTitle className="text-lg">
                    Konfigurasi API
                </CardTitle>
            </CardHeader>

            <CardContent className="space-y-4">

                {/* AUTH */}
                <div className="space-y-2">
                    <Label>Jenis Autentikasi</Label>

                    <RadioGroup
                        value={authType}
                        onValueChange={(v: string) => updateAuthType(v as AuthType)}
                        className="flex flex-col space-y-1"
                    >
                        <div className="flex items-center space-x-2">
                            <RadioGroupItem value="none" id="none" />
                            <Label htmlFor="none">Tanpa Autentikasi</Label>
                        </div>

                        <div className="flex items-center space-x-2">
                            <RadioGroupItem value="apiKey" id="apiKey" />
                            <Label htmlFor="apiKey">Authorization Header</Label>
                        </div>
                    </RadioGroup>

                    {authType === 'apiKey' && (
                        <div className="mt-3 p-3 border rounded-md bg-muted/40 space-y-3">
                            <div className="space-y-2">
                                <Label className="text-xs">Format Authorization</Label>
                                <Select
                                    value={authPrefix}
                                    onValueChange={(v: string) => updateAuthPrefix(v as AuthPrefix)}
                                >
                                    <SelectTrigger>
                                        <SelectValue />
                                    </SelectTrigger>
                                    <SelectContent>
                                        {AUTH_PREFIX_OPTIONS.map((option) => (
                                            <SelectItem key={option.value} value={option.value}>
                                                {option.label}
                                            </SelectItem>
                                        ))}
                                    </SelectContent>
                                </Select>
                            </div>

                            {authPrefix === 'custom' && (
                                <div className="space-y-2">
                                    <Label className="text-xs">Custom Prefix</Label>
                                    <Input
                                        placeholder="Masukkan prefix (mis: Token)"
                                        value={customAuthPrefix}
                                        onChange={(e: any) => updateCustomAuthPrefix(e.target.value)}
                                    />
                                </div>
                            )}


                            {/* INPUT API KEY */}
                            <div className="space-y-2 pt-3 border-t">
                                <Label className="text-xs">API Key / Token</Label>
                                <Input
                                    type="password"
                                    placeholder="Masukkan API Key atau pilih placeholder"
                                    value={apiKeyValue}
                                    onChange={(e: any) => updateApiKeyValue(e.target.value)}
                                    className="cursor-pointer"
                                />
                            </div>
                        </div>
                    )}
                </div>

                {/* HEADERS */}
                <div className="space-y-2">
                    <div className="flex justify-between items-center">
                        <Label>Custom Headers</Label>

                        <Button
                            type="button"
                            variant="outline"
                            size="sm"
                            onClick={addHeader}
                        >
                            Tambah Header
                        </Button>
                    </div>

                    <div className="space-y-3">
                        {headers.map((header, index) => (
                            <div key={index} className="flex gap-2 items-start">
                                <div className="flex-1">
                                    <Input
                                        placeholder="Header Key"
                                        value={header.key}
                                        onChange={(e: any) => {
                                            const updated = [...headers];
                                            updated[index].key = e.target.value;
                                            updateHeaders(updated);
                                        }}
                                    />
                                </div>

                                <div className="text-sm text-muted-foreground pt-2">=</div>

                                <div className="flex-1">
                                    <Popover
                                        open={openPopoverIndex === index}
                                        onOpenChange={(open) => {
                                            if (open) {
                                                setOpenPopoverIndex(index);
                                            } else {
                                                setOpenPopoverIndex(null);
                                            }
                                        }}
                                    >
                                        <PopoverTrigger asChild>
                                            <div className="relative">
                                                <Input
                                                    placeholder="Header Value"
                                                    value={header.value}
                                                    onChange={(e: any) => {
                                                        const updated = [...headers];
                                                        updated[index].value = e.target.value;
                                                        updateHeaders(updated);
                                                    }}
                                                    onClick={() => setOpenPopoverIndex(index)}
                                                    className="cursor-pointer"
                                                />
                                            </div>
                                        </PopoverTrigger>
                                        <PopoverContent className="w-48 p-0" align="start">
                                            <div className="p-2">
                                                <div className="text-xs font-medium mb-2 text-muted-foreground">
                                                    Pilih Placeholder
                                                </div>
                                                <div className="space-y-1">
                                                    {PLACEHOLDER_OPTIONS.map((option) => (
                                                        <button
                                                            key={option.value}
                                                            className="w-full text-left px-2 py-1.5 text-sm rounded hover:bg-accent hover:text-accent-foreground transition-colors"
                                                            onClick={() => handlePlaceholderSelect(index, option.value)}
                                                        >
                                                            {option.label}
                                                        </button>
                                                    ))}
                                                    <div className="border-t pt-2 mt-2">
                                                        <div className="text-xs font-medium mb-1 text-muted-foreground">
                                                            Custom
                                                        </div>
                                                        <Input
                                                            placeholder="Nama placeholder"
                                                            onClick={(e) => e.stopPropagation()}
                                                            onKeyDown={(e: any) => {
                                                                if (e.key === 'Enter') {
                                                                    const customValue = e.target.value.trim();
                                                                    if (customValue) {
                                                                        handlePlaceholderSelect(index, `{{${customValue}}}`);
                                                                    }
                                                                }
                                                            }}
                                                        />
                                                    </div>
                                                </div>
                                            </div>
                                        </PopoverContent>
                                    </Popover>
                                </div>

                                <Button
                                    type="button"
                                    variant="destructive"
                                    size="icon"
                                    onClick={() => removeHeader(index)}
                                    className="shrink-0 mt-0"
                                >
                                    ✕
                                </Button>
                            </div>
                        ))}
                    </div>
                </div>

            </CardContent>
        </Card>
    );
}