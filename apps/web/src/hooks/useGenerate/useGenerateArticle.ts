// Hapus import Toast langsung
// import { Toast } from "@kana-consultant/ui-kit";
import type { ToastContextType } from "@/hooks/use-toast"; // Import type untuk typing
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
    setSlug: (val: string | null) => void,
    setKeywords: (val: string[]) => void,
    toast: ToastContextType //   Tambahkan parameter toast
) {
    if (!topic) {
        toast.error("Masukkan topik terlebih dahulu"); //   Ganti toast.error dengan toast.error
        return;
    }
    if (!selectedModelId) {
        toast.error("Pilih model AI terlebih dahulu"); //   Ganti toast.error dengan toast.error
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
        setSlug(response.slug);
        setKeywords(response.keywords);
        console.log(response.slug);

        toast.success("Artikel berhasil di-generate!", { //   Ganti toast.success dengan toast.success
            description: `Topik: ${response.title} | ${response.wordCount} kata`,
        });
    } catch (error) {
        console.error("Error generating article:", error);
        toast.error("Gagal mengenerate artikel", { //   Ganti toast.error dengan toast.error
            description: error instanceof Error ? error.message : "Terjadi kesalahan pada server",
        });
    } finally {
        setLoadingArticle(false);
    }
}