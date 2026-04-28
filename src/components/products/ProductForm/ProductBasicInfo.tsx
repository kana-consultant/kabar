import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { Product } from "@/types/product";

interface ProductBasicInfoProps {
    product: Partial<Product>;
    onUpdate: (updates: Partial<Product>) => void;
    onTestConnection: () => void;
    isTesting: boolean;
}

export function ProductBasicInfo({ product, onUpdate, onTestConnection, isTesting }: ProductBasicInfoProps) {
    return (
        <Card className="shadow-sm">
            <CardHeader className="pb-3">
                <CardTitle className="text-lg">Informasi Dasar</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
                <div className="space-y-2">
                    <Label htmlFor="name" className="text-sm font-medium">
                        Nama Produk <span className="text-red-500">*</span>
                    </Label>
                    <Input
                        id="name"
                        value={product.name}
                        onChange={(e) => onUpdate({ name: e.target.value })}
                        placeholder="Contoh: Toko Saya"
                        className="h-10"
                    />
                </div>

                <div className="space-y-2">
                    <Label htmlFor="platform" className="text-sm font-medium">
                        Platform
                    </Label>
                    <Select
                        value={product.platform}
                        onValueChange={(v) => onUpdate({ platform: v as any })}
                    >
                        <SelectTrigger className="h-10">
                            <SelectValue placeholder="Pilih platform" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="wordpress">📝 WordPress</SelectItem>
                            <SelectItem value="shopify">🛍️ Shopify</SelectItem>
                            <SelectItem value="custom">🔧 Custom API</SelectItem>
                        </SelectContent>
                    </Select>
                </div>

                <div className="space-y-2">
                    <Label htmlFor="apiEndpoint" className="text-sm font-medium">
                        API Endpoint <span className="text-red-500">*</span>
                    </Label>
                    <Input
                        id="apiEndpoint"
                        value={product.apiEndpoint}
                        onChange={(e) => onUpdate({ apiEndpoint: e.target.value })}
                        placeholder="https://domain.com/wp-json/wp/v2/posts"
                        className="h-10 font-mono text-sm"
                    />
                </div>

                <div className="space-y-2">
                    <Label htmlFor="apiKey" className="text-sm font-medium">
                        API Key / Token
                    </Label>
                    <Input
                        id="apiKey"
                        type="password"
                        value={product.apiKey}
                        onChange={(e) => onUpdate({ apiKey: e.target.value })}
                        placeholder="Masukkan API Key"
                        className="h-10"
                    />
                    <p className="text-xs text-slate-400">
                        ⚠️ Key akan dienkripsi sebelum disimpan ke database
                    </p>
                </div>
            </CardContent>
        </Card>
    );
}