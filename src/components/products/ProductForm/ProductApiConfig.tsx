import { useEffect, useState } from "react";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import {
    Card,
    CardContent,
    CardHeader,
    CardTitle
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";

import type { AdapterConfig } from "@/types/product";

interface ProductApiConfigProps {
    config: Partial<AdapterConfig>;
    onUpdate: (updates: Partial<AdapterConfig>) => void;
}

interface HeaderItem {
    key: string;
    value: string;
}

type AuthType = 'none' | 'xapiKey' | 'bearer';

export function ProductApiConfig({
    config,
    onUpdate,
}: ProductApiConfigProps) {

    const [headers, setHeaders] = useState<HeaderItem[]>([]);
    const [authType, setAuthType] = useState<AuthType>('none');

    // convert object -> list
    useEffect(() => {
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

        if (hasBearer) setAuthType('bearer');
        else if (hasxApiKey) setAuthType('xapiKey');
        else setAuthType('none');

        // filter auth headers from list
        const regularHeaders: Record<string, string> = {};

        Object.entries(obj).forEach(([key, value]) => {
            const lowerKey = key.toLowerCase();

            if (
                !lowerKey.includes('api-key') &&
                lowerKey !== 'x-api-key' &&
                !(lowerKey === 'authorization')
            ) {
                regularHeaders[key] = value;
            }
        });

        const mapped: HeaderItem[] = Object.entries(regularHeaders).map(
            ([key, value]) => ({ key, value: String(value) })
        );

        setHeaders(mapped.length ? mapped : [{ key: "", value: "" }]);

    }, [config.customHeaders]);

    const updateFullConfig = (
        currentHeaders: HeaderItem[],
        currentAuthType: AuthType
    ) => {

        const allHeaders: Record<string, string> = {};

        // custom headers
        currentHeaders.forEach(item => {
            if (item.key.trim()) {
                allHeaders[item.key.trim()] = item.value;
            }
        });

        // auth headers WITH PLACEHOLDER (INI INTINYA FIX)
        if (currentAuthType === 'xapiKey') {
            allHeaders['X-API-Key'] = '{{api_key}}';
        }

        if (currentAuthType === 'bearer') {
            allHeaders['Authorization'] = 'Bearer {{api_key}}';
        }

        onUpdate({
            customHeaders: JSON.stringify(allHeaders)
        });
    };

    const updateHeaders = (updated: HeaderItem[]) => {
        setHeaders(updated);
        updateFullConfig(updated, authType);
    };

    const updateAuthType = (newType: AuthType) => {
        setAuthType(newType);
        updateFullConfig(headers, newType);
    };

    const addHeader = () => {
        updateHeaders([...headers, { key: "", value: "" }]);
    };

    const removeHeader = (index: number) => {
        const updated = headers.filter((_, i) => i !== index);
        updateHeaders(updated.length ? updated : [{ key: "", value: "" }]);
    };

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
                        onChange={(e) => onUpdate({ endpointPath: e.target.value })}
                        placeholder="/api/v1/products"
                    />
                </div>

                {/* METHOD */}
                <div className="space-y-2">
                    <Label>HTTP Method</Label>
                    <Select
                        value={config.httpMethod || "POST"}
                        onValueChange={(v) => onUpdate({ httpMethod: v as any })}
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
                        onValueChange={(value) => updateAuthType(value as AuthType)}
                        className="flex flex-col space-y-1"
                    >
                        <div className="flex items-center space-x-2">
                            <RadioGroupItem value="none" id="none" />
                            <Label htmlFor="none">Tanpa Autentikasi</Label>
                        </div>

                        <div className="flex items-center space-x-2">
                            <RadioGroupItem value="xapiKey" id="xapiKey" />
                            <Label htmlFor="xapiKey">API Key</Label>
                        </div>

                        <div className="flex items-center space-x-2">
                            <RadioGroupItem value="bearer" id="bearer" />
                            <Label htmlFor="bearer">Bearer Token</Label>
                        </div>
                    </RadioGroup>
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
                                    onChange={(e) => {
                                        const updated = [...headers];
                                        updated[index].key = e.target.value;
                                        updateHeaders(updated);
                                    }}
                                />

                                <Input
                                    placeholder="Header Value"
                                    value={header.value}
                                    onChange={(e) => {
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