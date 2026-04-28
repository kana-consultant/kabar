import { useState, useEffect } from "react";
import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Loader2, Sparkles } from "lucide-react";
import { useGenerate } from "@/hooks/useGenerate";
import { useNavigate } from "@tanstack/react-router";

export function QuickGenerate() {

    const [loading, setLoading] = useState(false);

    const navigate = useNavigate();

    const {
        products,
        topic,
        setSelectedProducts,
        selectedProducts,
        setTopic,
        quickGenerate
    } = useGenerate();


    // auto pilih produk pertama
    useEffect(() => {
        if (
            products.length > 0 &&
            selectedProducts.length === 0
        ) {
            setSelectedProducts([
                products[0].id
            ]);
        }
    }, [products]);


    const handleGenerate = async () => {

        setLoading(true);

        try {

            const draftId =
                await quickGenerate();

            if (draftId) {
                navigate({
                  to:`/history`
                })
            }

        } finally {
            setLoading(false);
        }
    };

    return (
        <Card>
            <CardHeader>
                <CardTitle>
                    Quick Generate
                </CardTitle>

                <CardDescription>
                    Masukkan topik, sistem akan
                    generate artikel + gambar otomatis
                </CardDescription>

            </CardHeader>


            <CardContent className="space-y-4">

                <div>
                    <label className="mb-2 block text-sm font-medium">
                        Topik / Keyword
                    </label>

                    <Input
                        value={topic}
                        placeholder="Contoh: Cara Memilih Sepatu Gunung untuk Pemula"
                        onChange={(e) =>
                            setTopic(
                                e.target.value
                            )
                        }
                    />
                </div>


                <div>
                    <label className="mb-2 block text-sm font-medium">
                        Target Produk
                    </label>

                    <select
                        className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"

                        value={
                            selectedProducts[0] || ""
                        }

                        onChange={(e) =>
                            setSelectedProducts([
                                e.target.value
                            ])
                        }
                    >

                        {products.map((item) => (
                            <option
                                key={item.id}
                                value={item.id}
                            >
                                {item.name}
                            </option>
                        ))}

                    </select>

                </div>

            </CardContent>


            <CardFooter>

                <Button
                    onClick={handleGenerate}
                    disabled={
                        loading ||
                        !topic
                    }
                    className="w-full"
                >

                    {loading ? (
                        <>
                            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                            Generating...
                        </>
                    ) : (
                        <>
                            <Sparkles className="mr-2 h-4 w-4" />
                            Generate Konten
                        </>
                    )}

                </Button>

            </CardFooter>
        </Card>
    );
}