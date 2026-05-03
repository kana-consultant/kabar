import { toast } from "sonner";
import { generateArticle, generateImage, type GenerateArticleRequest } from "@/services/generateService";
import type { Tone, ArticleLength, Language } from "./types";

interface UseArticleGenerationProps {
    topic: string;
    selectedModelId: string;
    tone: Tone;
    articleLength: ArticleLength;
    keywords: string[];
    language: Language;
    autoGenerateImage: boolean;
    setArticle: (article: string) => void;
    setArticleResponse: (response: any) => void;
    setImageUrl: (url: string) => void;
    setSeoScore: (score: number) => void;
    setReadabilityScore: (score: number) => void;
    setWordCount: (count: number) => void;
    setLoadingArticle: (loading: boolean) => void;
    setLoadingImage: (loading: boolean) => void;
}

export function useArticleGeneration(props: UseArticleGenerationProps) {
    const generateArticleContent = async () => {
        if (!props.topic) {
            toast.error("Masukkan topik terlebih dahulu");
            return;
        }
        if (!props.selectedModelId) {
            toast.error("Pilih model AI terlebih dahulu");
            return;
        }

        props.setLoadingArticle(true);
        try {
            const request: GenerateArticleRequest = {
                topic: props.topic,
                modelId: props.selectedModelId,
                tone: props.tone,
                length: props.articleLength,
                keywords: props.keywords.length > 0 ? props.keywords : undefined,
                language: props.language,
                autoGenerateImage: props.autoGenerateImage,
            };

            const response = await generateArticle(request);
            props.setArticleResponse(response);
            props.setArticle(response.content);
            props.setSeoScore(response.seoScore);
            props.setReadabilityScore(response.readabilityScore);
            props.setWordCount(response.wordCount);

            if (props.autoGenerateImage && response.imageUrl) {
                props.setImageUrl(response.imageUrl);
            }

            toast.success("Artikel berhasil di-generate!", {
                description: `Topik: ${response.title} | ${response.wordCount} kata`,
            });
        } catch (error) {
            console.error("Error generating article:", error);
            toast.error("Gagal mengenerate artikel", {
                description: error instanceof Error ? error.message : "Terjadi kesalahan pada server",
            });
        } finally {
            props.setLoadingArticle(false);
        }
    };

    const generateImageManually = async (articleResponse: any, topic: string) => {
        if (!articleResponse?.imagePrompt && !topic) {
            toast.error("Generate artikel terlebih dahulu");
            return;
        }

        props.setLoadingImage(true);
        const prompt = articleResponse?.imagePrompt || `Ilustrasi tentang ${topic}`;

        try {
            const response = await generateImage({
                prompt: prompt,
                modelId: props.selectedModelId
            });
            props.setImageUrl(response.imageUrl);
            toast.success("Gambar berhasil di-generate!");
        } catch (error) {
            toast.error("Gagal mengenerate gambar");
        } finally {
            props.setLoadingImage(false);
        }
    };

    return {
        generateArticleContent,
        generateImageManually
    };
}