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

import type { AdapterConfig } from "@/services/product";

interface ProductApiConfigProps {
    config: Partial<AdapterConfig>;
    onUpdate: (updates: Partial<AdapterConfig>) => void;
}

interface HeaderItem {
    key: string;
    value: string;
}

type AuthType = 'none' | 'xapiKey' | 'bearer' | 'custom';

export function ProductApiConfig({
    config,
    onUpdate,
}: ProductApiConfigProps) {

    const [headers, setHeaders] = useState<HeaderItem[]>([]);
    const [authType, setAuthType] = useState<AuthType>('none');
    const [customAuthKey, setCustomAuthKey] = useState<string>('');
    const [initialized, setInitialized] = useState(false);
    
    // Gunakan ref untuk melacak apakah perubahan berasal dari internal
    const internalUpdateRef = useRef(false);

    // convert object -> list (hanya saat mount)
    useEffect(() => {
        if (initialized) return;

        let obj: Record<string, string> = {};

        if (config.customHeaders) {
            if (typeof config.customHeaders === 'string') {
                try {
                    obj = JSON.parse(config.customHeaders);
                } catch (e) {
                    obj = {};
                }
            } else if (typeof config.customHeaders === 'object') {
                obj = config.customHeaders as Record<string, string>;
            }
        }

        // detect auth type
        const hasxApiKey = Object.keys(obj).some(key =>
            key.toLowerCase().includes('api-key') ||
            key.toLowerCase() === 'x-api-key'
        );

        const hasBearer = Object.keys(obj).some(key =>
            key.toLowerCase() === 'authorization' &&
            String(obj[key]).toLowerCase().includes('bearer')
        );

        // detect custom auth
        const customAuthEntry = Object.entries(obj).find(([key, value]) => {
            const lowerKey = key.toLowerCase();
            return (
                String(value).includes('{{api_key}}') &&
                !lowerKey.includes('api-key') &&
                lowerKey !== 'x-api-key' &&
                lowerKey !== 'authorization'
            );
        });

        if (hasBearer) {
            setAuthType('bearer');
        } else if (hasxApiKey) {
            setAuthType('xapiKey');
        } else if (customAuthEntry) {
            setAuthType('custom');
            setCustomAuthKey(customAuthEntry[0]);
        } else {
            setAuthType('none');
        }

        // filter auth headers from list
        const regularHeaders: Record<string, string> = {};

        Object.entries(obj).forEach(([key, value]) => {
            const lowerKey = key.toLowerCase();
            const isAuthHeader =
                lowerKey.includes('api-key') ||
                lowerKey === 'x-api-key' ||
                lowerKey === 'authorization' ||
                (customAuthEntry && key === customAuthEntry[0]);

            if (!isAuthHeader) {
                regularHeaders[key] = value;
            }
        });

        const mapped: HeaderItem[] = Object.entries(regularHeaders).map(
            ([key, value]) => ({ key, value: String(value) })
        );

        setHeaders(mapped.length ? mapped : [{ key: "", value: "" }]);
        setInitialized(true);

    }, [config.customHeaders, initialized]);

    const updateFullConfig = useCallback((
        currentHeaders: HeaderItem[],
        currentAuthType: AuthType,
        currentCustomKey?: string,
    ) => {
        // Tandai bahwa ini adalah update internal
        internalUpdateRef.current = true;
        
        const allHeaders: Record<string, string> = {};

        // custom headers
        currentHeaders.forEach(item => {
            if (item.key.trim()) {
                allHeaders[item.key.trim()] = item.value;
            }
        });

        // auth headers
        if (currentAuthType === 'xapiKey') {
            allHeaders['X-API-Key'] = '{{api_key}}';
        } else if (currentAuthType === 'bearer') {
            allHeaders['Authorization'] = 'Bearer {{api_key}}';
        } else if (currentAuthType === 'custom') {
            const key = (currentCustomKey ?? customAuthKey).trim();
            if (key) {
                allHeaders[key] = '{{api_key}}';
            }
        }

        onUpdate({
            customHeaders: JSON.stringify(allHeaders)
        });
        
        // Reset flag setelah update
        setTimeout(() => {
            internalUpdateRef.current = false;
        }, 0);
    }, [customAuthKey, onUpdate]);

    const updateHeaders = useCallback((updated: HeaderItem[]) => {
        setHeaders(updated);
        updateFullConfig(updated, authType);
    }, [authType, updateFullConfig]);

    const updateAuthType = useCallback((newType: AuthType) => {
        setAuthType(newType);
        updateFullConfig(headers, newType);
    }, [headers, updateFullConfig]);

    const updateCustomAuthKey = useCallback((key: string) => {
        setCustomAuthKey(key);
        updateFullConfig(headers, authType, key);
    }, [headers, authType, updateFullConfig]);

    const addHeader = useCallback(() => {
        updateHeaders([...headers, { key: "", value: "" }]);
    }, [headers, updateHeaders]);

    const removeHeader = useCallback((index: number) => {
        const updated = headers.filter((_, i) => i !== index);
        updateHeaders(updated.length ? updated : [{ key: "", value: "" }]);
    }, [headers, updateHeaders]);

    return (
        <Card className="shadow-sm">
            <CardHeader className="pb-3">
                <CardTitle className="text-lg">
                    Konfigurasi API
                </CardTitle>
            </CardHeader>

            <CardContent className="space-y-4">

                {/* ENDPOINT */}
                <div className="space-y-2">
                    <Label>Endpoint Path</Label>
                    <Input
                        value={config.endpointPath || ""}
                        onChange={(e: any) => onUpdate({ endpointPath: e.target.value })}
                        placeholder="/api/v1/products"
                    />
                </div>

                {/* METHOD */}
                <div className="space-y-2">
                    <Label>HTTP Method</Label>
                    <Select
                        value={config.httpMethod || "POST"}
                        onValueChange={(v: any) => onUpdate({ httpMethod: v as any })}
                    >
                        <SelectTrigger>
                            <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="GET">GET</SelectItem>
                            <SelectItem value="POST">POST</SelectItem>
                            <SelectItem value="PUT">PUT</SelectItem>
                            <SelectItem value="PATCH">PATCH</SelectItem>
                        </SelectContent>
                    </Select>
                </div>

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
                            <RadioGroupItem value="xapiKey" id="xapiKey" />
                            <Label htmlFor="xapiKey">API Key (X-API-Key)</Label>
                        </div>

                        <div className="flex items-center space-x-2">
                            <RadioGroupItem value="bearer" id="bearer" />
                            <Label htmlFor="bearer">Bearer Token</Label>
                        </div>

                        <div className="flex items-center space-x-2">
                            <RadioGroupItem value="custom" id="custom" />
                            <Label htmlFor="custom">Custom Header</Label>
                        </div>
                    </RadioGroup>

                    {/* Custom Auth Fields */}
                    {authType === 'custom' && (
                        <div className="mt-3 p-3 border rounded-md bg-muted/40 space-y-3">
                            <div className="space-y-1">
                                <Label className="text-xs">Nama Header</Label>
                                <Input
                                    placeholder="Contoh: X-Custom-Token"
                                    value={customAuthKey}
                                    onChange={(e: any) => updateCustomAuthKey(e.target.value)}
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
                            <div key={index} className="flex gap-2">
                                <Input
                                    placeholder="Header Key"
                                    value={header.key}
                                    onChange={(e: any) => {
                                        const updated = [...headers];
                                        updated[index].key = e.target.value;
                                        updateHeaders(updated);
                                    }}
                                />

                                <Input
                                    placeholder="Header Value"
                                    value={header.value}
                                    onChange={(e: any) => {
                                        const updated = [...headers];
                                        updated[index].value = e.target.value;
                                        updateHeaders(updated);
                                    }}
                                />

                                <Button
                                    type="button"
                                    variant="destructive"
                                    onClick={() => removeHeader(index)}
                                >
                                    Hapus
                                </Button>
                            </div>
                        ))}
                    </div>
                </div>

            </CardContent>
        </Card>
    );
}