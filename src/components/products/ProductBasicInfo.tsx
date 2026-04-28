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
    product: Product;
    onChange: (updates: Partial<Product>) => void;
}

export function ProductBasicInfo({ product, onChange }: ProductBasicInfoProps) {
    return (
        <Card>
            <CardHeader>
                <CardTitle>Informasi Dasar</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
                <div>
                    <Label>Nama Produk</Label>
                    <Input
                        value={product.name}
                        onChange={(e) => onChange({ name: e.target.value })}
                        placeholder="Contoh: Toko Saya"
                    />
                </div>
                <div>
                    <Label>Platform</Label>
                    <Select
                        value={product.platform}
                        onValueChange={(v) => onChange({ platform: v as any })}
                    >
                        <SelectTrigger>
                            <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="wordpress">WordPress</SelectItem>
                            <SelectItem value="shopify">Shopify</SelectItem>
                            <SelectItem value="custom">Custom API</SelectItem>
                        </SelectContent>
                    </Select>
                </div>
                <div>
                    <Label>API Endpoint</Label>
                    <Input
                        value={product.apiEndpoint}
                        onChange={(e) => onChange({ apiEndpoint: e.target.value })}
                        placeholder="https://domain.com/wp-json/wp/v2/posts"
                    />
                </div>
                <div>
                    <Label>API Key / Token</Label>
                    <Input
                        type="password"
                        value={product.apiKey || product.APIKeyEncrypted}
                        onChange={(e) => onChange({ apiKey: e.target.value })}
                        placeholder="Masukkan API Key"
                    />
                    <p className="mt-1 text-xs text-slate-500">
                        ⚠️ Key akan dienkripsi sebelum disimpan
                    </p>
                </div>
            </CardContent>
        </Card>
    );
}