// Hapus import Toast langsung
import type { ToastContextType } from "@/hooks/use-toast";
import { generateArticle, type GenerateArticleRequest } from "@/services/generate";
import { createDraft, type CreateDraftRequest } from "@/services/draft"; // 👈 Import createDraft

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
    toast: ToastContextType,
    autoGenerateImage?: boolean,
    imageModelId?: string,
    // 👇 Parameter opsional untuk draft
    teamId?: string,
    targetProducts?: string[]
) {
    if (!topic) {
        toast.error("Masukkan topik terlebih dahulu");
        return;
    }
    if (!selectedModelId) {
        toast.error("Pilih model AI terlebih dahulu");
        return;
    }
    if (autoGenerateImage && !imageModelId) {
        toast.error("Pilih model gambar terlebih dahulu");
        return;
    }

    setLoadingArticle(true);
    try {
        // 👇 Cek current path - gunakan includes untuk fleksibilitas
        const currentPath = window.location.pathname;
        const isDraft = !currentPath.includes("/generate");

        const request: GenerateArticleRequest = {
            topic: topic,
            modelId: selectedModelId,
            tone: tone,
            length: articleLength,
            language: language,
            autoGenerateImage: autoGenerateImage ?? false,
            imageModelId: autoGenerateImage ? imageModelId : undefined,
        };

        const response = await generateArticle(request);

        // 👇 Jika bukan di halaman generate, simpan sebagai draft
        if (isDraft) {
            try {
                // 👇 Sesuaikan dengan tipe CreateDraftRequest
                const draftRequest: CreateDraftRequest = {
                    title: response.title || topic, // Fallback ke topic jika title tidak ada
                    topic: topic,
                    article: response.content,
                    slug: response.slug || "",
                    keywords: response.keywords || [],
                    excerpt: response.excerpt || "",
                    target_products: targetProducts || [],
                    team_id: teamId,
                    image_url: response.imageUrl, // Jika ada dari response
                    image_prompt: response.imagePrompt, // Jika ada dari response
                };

                await createDraft(draftRequest);

                toast.success("Artikel berhasil disimpan sebagai draft!", {
                    description: `Topik: ${response.title || topic} | ${response.wordCount} kata`,
                });

                // Update state untuk konsistensi UI
                setArticleResponse(response);
                setArticle(response.content);
                setSeoScore(response.seoScore);
                setReadabilityScore(response.readabilityScore);
                setWordCount(response.wordCount);
                setSlug(response.slug);
                setKeywords(response.keywords || []);
                setExcerpt(response.excerpt);

                return; // Hentikan eksekusi
            } catch (draftError) {
                console.error("Error saving draft:", draftError);
                toast.error("Gagal menyimpan draft", {
                    description: "Artikel berhasil dibuat tapi gagal disimpan sebagai draft",
                });
                // Tetap lanjutkan untuk menampilkan artikel di editor
            }
        }

        // 👇 Jika di halaman generate, tampilkan seperti biasa
        setArticleResponse(response);
        setArticle(response.content);
        setSeoScore(response.seoScore);
        setReadabilityScore(response.readabilityScore);
        setWordCount(response.wordCount);
        setSlug(response.slug);
        setKeywords(response.keywords || []);
        setExcerpt(response.excerpt);
        console.log(response.slug);

        toast.success("Artikel berhasil di-generate!", {
            description: autoGenerateImage
                ? `Topik: ${response.title} | ${response.wordCount} kata | gambar disertakan`
                : `Topik: ${response.title} | ${response.wordCount} kata`,
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