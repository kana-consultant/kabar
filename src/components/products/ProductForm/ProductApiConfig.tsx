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

import type { AdapterConfig } from "@/types/product";

interface ProductApiConfigProps {
    config: Partial<AdapterConfig>;
    onUpdate: (updates: Partial<AdapterConfig>) => void;
}

export function ProductApiConfig({
    config,
    onUpdate,
}: ProductApiConfigProps) {

    const [headersText, setHeadersText] = useState(
        JSON.stringify(config.customHeaders || {}, null, 2)
    );

    useEffect(() => {
        setHeadersText(
            JSON.stringify(config.customHeaders || {}, null, 2)
        );
    }, [config.customHeaders]);

    return (
        <Card className="shadow-sm">
            <CardHeader className="pb-3">
                <CardTitle className="text-lg">
                    Konfigurasi API
                </CardTitle>
            </CardHeader>

            <CardContent className="space-y-4">

                <div className="space-y-2">
                    <Label htmlFor="endpointPath">
                        Endpoint Path
                    </Label>

                    <Input
                        id="endpointPath"
                        value={config.endpointPath || ""}
                        onChange={(e) =>
                            onUpdate({
                                endpointPath: e.target.value
                            })
                        }
                    />
                </div>


                <div className="space-y-2">
                    <Label>
                        HTTP Method
                    </Label>

                    <Select
                        value={config.httpMethod || "POST"}
                        onValueChange={(v) =>
                            onUpdate({
                                httpMethod: v as any
                            })
                        }
                    >
                        <SelectTrigger>
                            <SelectValue />
                        </SelectTrigger>

                        <SelectContent>
                            <SelectItem value="POST">
                                POST
                            </SelectItem>

                            <SelectItem value="PUT">
                                PUT
                            </SelectItem>

                            <SelectItem value="PATCH">
                                PATCH
                            </SelectItem>

                        </SelectContent>
                    </Select>
                </div>



                <div className="space-y-2">
                    <Label htmlFor="headers">
                        Custom Headers
                    </Label>

                    <textarea
                        id="headers"
                        rows={6}
                        value={headersText}
                        onChange={(e) => {
                            const val = e.target.value;

                            // selalu biarkan user mengetik
                            setHeadersText(val);

                            // hanya update parent kalau valid json
                            try {
                                const parsed = JSON.parse(val);

                                onUpdate({
                                    customHeaders: parsed
                                });

                            } catch {
                                // biarkan user lanjut ngetik
                            }
                        }}
                        className="w-full rounded-md border p-3 font-mono text-sm"
                    />

                    <p className="text-xs text-slate-400">
                        Format JSON
                    </p>

                </div>

            </CardContent>
        </Card>
    )
}