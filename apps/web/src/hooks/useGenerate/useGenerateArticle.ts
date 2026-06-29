// Hapus import Toast langsung
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
    setExcerpt: (val: string | null) => void,
    toast: ToastContextType, //   Tambahkan parameter toast
    autoGenerateImage?: boolean,   // 👈 BARU: apakah generate gambar inline diminta
    imageModelId?: string          // 👈 BARU: model gambar yang dipilih user
) {
    if (!topic) {
        toast.error("Masukkan topik terlebih dahulu"); //   Ganti toast.error dengan toast.error
        return;
    }
    if (!selectedModelId) {
        toast.error("Pilih model AI terlebih dahulu"); //   Ganti toast.error dengan toast.error
        return;
    }
    // 👇 BARU: pengaman ganda — kalau toggle aktif tapi model gambar belum dipilih, jangan lanjut
    if (autoGenerateImage && !imageModelId) {
        toast.error("Pilih model gambar terlebih dahulu");
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
            // 👇 BARU: diteruskan ke backend supaya artikel ini juga di-generate
            // dengan beberapa gambar inline, menggunakan model gambar yang dipilih.
            generateImages: autoGenerateImage ?? false,
            imageModelId: autoGenerateImage ? imageModelId : undefined,
        };

        const response = await generateArticle(request);
        setArticleResponse(response);
        setArticle(response.content);
        setSeoScore(response.seoScore);
        setReadabilityScore(response.readabilityScore);
        setWordCount(response.wordCount);
        setSlug(response.slug);
        setKeywords(response.keywords);
        setExcerpt(response.excerpt)
        console.log(response.slug);

        toast.success("Artikel berhasil di-generate!", { //   Ganti toast.success dengan toast.success
            description: autoGenerateImage
                ? `Topik: ${response.title} | ${response.wordCount} kata | gambar disertakan`
                : `Topik: ${response.title} | ${response.wordCount} kata`,
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