// services/generate-article.ts
import type { ToastContextType } from "@/hooks/use-toast";
import { generateArticle, type GenerateArticleRequest } from "@/services/generate";
import { createDraft, type CreateDraftRequest } from "@/services/draft";
import { pageTracker } from "@/store/current-page"

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
    setTopic: (val: string) => void,
    toast: ToastContextType,
    autoGenerateImage?: boolean,
    imageModelId?: string,
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

        // 👇 Cek dari cache, bukan window.location
        const currentPage = pageTracker.get();
        const userStillOnGeneratePage = currentPage.includes('/generate');

        if (!userStillOnGeneratePage) {
            // User sudah pindah halaman → save draft
            try {
                const draftRequest: CreateDraftRequest = {
                    title: response.title || topic,
                    topic: topic,
                    article: response.content,
                    slug: response.slug || "",
                    keywords: response.keywords || [],
                    excerpt: response.excerpt || "",
                    target_products: targetProducts || [],
                    team_id: teamId,
                    image_url: response.imageUrl,
                    image_prompt: response.imagePrompt,
                };

                await createDraft(draftRequest);

                toast.success("Artikel berhasil disimpan sebagai draft!", {
                    description: `${response.title || topic} | ${response.wordCount} kata`,
                });
            } catch (draftError) {
                console.error("Error saving draft:", draftError);
                toast.error("Gagal menyimpan draft", {
                    description: "Artikel berhasil dibuat tapi gagal disimpan",
                });
            }

            // Tetap update state (redundant tapi aman)
            setArticleResponse(response);
            setArticle(response.content);
            setTopic(response.title);
            setSeoScore(response.seoScore);
            setReadabilityScore(response.readabilityScore);
            setWordCount(response.wordCount);
            setSlug(response.slug);
            setKeywords(response.keywords || []);
            setExcerpt(response.excerpt);

        } else {
            // User masih di halaman generate → tampilkan editor
            setArticleResponse(response);
            setTopic(response.title);
            setArticle(response.content);
            setSeoScore(response.seoScore);
            setReadabilityScore(response.readabilityScore);
            setWordCount(response.wordCount);
            setSlug(response.slug);
            setKeywords(response.keywords || []);
            setExcerpt(response.excerpt);

            toast.success("Artikel berhasil di-generate!", {
                description: autoGenerateImage
                    ? `Topik: ${response.title} | ${response.wordCount} kata | gambar disertakan`
                    : `Topik: ${response.title} | ${response.wordCount} kata`,
            });
        }
    } catch (error) {
        console.error("Error generating article:", error);
        toast.error("Gagal mengenerate artikel", {
            description: error instanceof Error ? error.message : "Terjadi kesalahan pada server",
        });
    } finally {
        setLoadingArticle(false);
    }
}