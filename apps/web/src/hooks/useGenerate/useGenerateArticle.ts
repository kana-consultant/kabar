import { toast } from "sonner";
import { generateArticle, type GenerateArticleRequest } from "@/services/generate";

export async function generateArticleContent(
    topic: string,
    selectedModelId: string,
    tone: "professional" | "casual" | "friendly" | "formal",
    articleLength: "short" | "medium" | "long", 
    language: "id" | "en",
    setLoadingArticle: (val: boolean) => void,
    setArticleResponse: (val: any) => void,
    setArticle: (val: string) => void,
    setSeoScore: (val: number | null) => void,
    setReadabilityScore: (val: number | null) => void,
    setWordCount: (val: number | null) => void,
) {
    if (!topic) {
        toast.error("Masukkan topik terlebih dahulu");
        return;
    }
    if (!selectedModelId) {
        toast.error("Pilih model AI terlebih dahulu");
        return;
    }

    setLoadingArticle(true);
    try {
        const request: GenerateArticleRequest = {
            topic: topic,
            modelId: selectedModelId,
            tone: tone,
            length: articleLength,
            language: language,
        };

        const response = await generateArticle(request);
        setArticleResponse(response);
        setArticle(response.content);
        setSeoScore(response.seoScore);
        setReadabilityScore(response.readabilityScore);
        setWordCount(response.wordCount);

        toast.success("Artikel berhasil di-generate!", {
            description: `Topik: ${response.title} | ${response.wordCount} kata`,
        });
    } catch (error) {
        console.error("Error generating article:", error);
        toast.error("Gagal mengenerate artikel", {
            description: error instanceof Error ? error.message : "Terjadi kesalahan pada server",
        });
    } finally {
        setLoadingArticle(false);
    }
}