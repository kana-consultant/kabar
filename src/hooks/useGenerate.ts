// src/hooks/useGenerate.ts
import { useState, useEffect, useCallback } from "react";
import { useLocation } from "@tanstack/react-router";
import { toast } from "sonner";
import { createDraft, updateDraft, getDraftById, publishDraft, draftSchedule, publishDraftInstant, type PublishResponse, type CreateDraftRequest } from "@/services/draftService";
import { generateArticle, generateImage, type GenerateArticleRequest, type GenerateArticleResponse } from "@/services/generateService";
import { getProducts } from "@/services/productService";
import { getModelsFromAPIKeys, type AIModel } from "@/services/modelService";
import { type Product } from "@/types/product";
import { getHistoryById } from "@/services/historyService";
import type { ScheduleRequest } from "@/types/schedule";
import { getAPIKeys, type APIKey } from "@/services/apiKeyService";


export function useGenerate() {
    const location = useLocation();
    const searchParams = new URLSearchParams(location.search);
    const editId = searchParams.get("edit") || undefined;
    const TopicId = searchParams.get("topic") || undefined;

    // Basic states
    const [topic, setTopic] = useState("");
    const [article, setArticle] = useState("");
    const [articleResponse, setArticleResponse] = useState<GenerateArticleResponse | null>(null);
    const [imageUrl, setImageUrl] = useState("");
    const [loadingArticle, setLoadingArticle] = useState(false);
    const [loadingImage, setLoadingImage] = useState(false);
    const [selectedProducts, setSelectedProducts] = useState<string[]>([]);
    const [postMode, setPostMode] = useState<"instant" | "scheduled" | "draft">("instant");
    const [scheduleTime, setScheduleTime] = useState("09:00");
    const [scheduleDate, setScheduleDate] = useState(() => {
        const tomorrow = new Date();
        tomorrow.setDate(tomorrow.getDate() + 1);
        return tomorrow.toISOString().split("T")[0];
    });
    const [dailySchedule, setDailySchedule] = useState(false);
    const [dailyTime, setDailyTime] = useState("08:00");
    const [autoGenerateImage, setAutoGenerateImage] = useState(false);
    const [postToAll, setPostToAll] = useState(false);
    const [currentDraftId, setCurrentDraftId] = useState<string | null>(null);
    const [publishing, setPublishing] = useState(false);

    // Result dialog states
    const [showResultDialog, setShowResultDialog] = useState(false);
    const [publishResults, setPublishResults] = useState<PublishResponse | null>(null);

    // Model states
    const [models, setModels] = useState<AIModel[]>([]);
    const [selectedModelId, setSelectedModelId] = useState("");
    const [loadingModels, setLoadingModels] = useState(true);

    // Article generation options
    const [tone, setTone] = useState<"professional" | "casual" | "friendly" | "formal">("professional");
    const [articleLength, setArticleLength] = useState<"short" | "medium" | "long">("medium");
    const [keywords, setKeywords] = useState<string[]>([]);
    const [keywordInput, setKeywordInput] = useState("");
    const [language, setLanguage] = useState<"id" | "en">("id");

    // SEO metrics
    const [seoScore, setSeoScore] = useState<number | null>(null);
    const [readabilityScore, setReadabilityScore] = useState<number | null>(null);
    const [wordCount, setWordCount] = useState<number | null>(null);

    // Products state
    const [products, setProducts] = useState<Product[]>([]);
    const [productNames, setProductNames] = useState<string[]>([]);
    const [productsLoading, setProductsLoading] = useState(true);
    const [productsError, setProductsError] = useState<string | null>(null);

    // Fetch products
    const fetchProducts = useCallback(async () => {
        setProductsLoading(true);
        setProductsError(null);
        try {
            const productsData = await getProducts();
            setProducts(productsData || []);
            setProductNames(productsData ? productsData.map(p => p.name) : []);
        } catch (error) {
            console.error("Failed to fetch products:", error);
            setProductsError("Gagal memuat data produk");
            setProductNames(["TrekkingID", "CampingMart", "OutdoorGear"]);
        } finally {
            setProductsLoading(false);
        }
    }, []);

    // Load models from API keys
    const loadModels = useCallback(async () => {
        setLoadingModels(true);
        try {
            const modelsData = await getModelsFromAPIKeys();
            setModels(modelsData as any);
            if (modelsData.length > 0 && !selectedModelId) {
                setSelectedModelId(modelsData[0].id);
            }
        } catch (error) {
            console.error("Failed to load models:", error);
            toast.error("Gagal memuat model AI");
        } finally {
            setLoadingModels(false);
        }
    }, []);

    useEffect(() => {
        fetchProducts();
        loadModels();
    }, []);

    // Load draft if edit mode
    useEffect(() => {
        const LoadDraft = async () => {
            if (editId) {
                try {
                    const draft = await getDraftById(editId);
                    if (draft) {
                        setTopic(draft.topic);
                        setArticle(draft.article);
                        setImageUrl(draft.imageUrl || "");
                        setSelectedProducts(draft.targetProducts || []);
                        setCurrentDraftId(draft.id);
                        toast.info("Memuat draft", {
                            description: `"${draft.title}" siap diedit`,
                        });
                    } else {
                        toast.error("Draft tidak ditemukan");
                    }
                } catch (error) {
                    console.error("Failed to load draft:", error);
                    toast.error("Gagal memuat draft");
                }
            } else if (TopicId) {
                try {
                    const draft = await getHistoryById(TopicId);
                    if (draft) {
                        setTopic(draft.topic);
                        setArticle(draft.content);
                        setImageUrl(draft.imageUrl || "");
                        setSelectedProducts(draft.targetProducts || []);
                        // setCurrentDraftId(draft.id);
                        toast.info("Memuat draft", {
                            description: `"${draft.title}" siap diedit`,
                        });
                    } else {
                        toast.error("Draft tidak ditemukan");
                    }
                } catch (error) {
                    console.error("Failed to load draft:", error);
                    toast.error("Gagal memuat draft");
                }
            }
        };
        LoadDraft();
    }, [editId]);

    const generateArticleContent = async () => {
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
                keywords: keywords.length > 0 ? keywords : undefined,
                language: language,
                autoGenerateImage: autoGenerateImage,
            };

            const response = await generateArticle(request);
            setArticleResponse(response);
            setArticle(response.content);
            setSeoScore(response.seoScore);
            setReadabilityScore(response.readabilityScore);
            setWordCount(response.wordCount);

            if (autoGenerateImage && response.imageUrl) {
                setImageUrl(response.imageUrl);
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
            setLoadingArticle(false);
        }
    };

    const generateImageManually = async () => {
        if (!articleResponse?.imagePrompt && !topic) {
            toast.error("Generate artikel terlebih dahulu");
            return;
        }

        setLoadingImage(true);
        const prompt = articleResponse?.imagePrompt || `Ilustrasi tentang ${topic}`;

        try {
            const response = await generateImage({
                prompt: prompt,
                modelId: selectedModelId
            });
            setImageUrl(response.imageUrl);
            toast.success("Gambar berhasil di-generate!");
        } catch (error) {
            toast.error("Gagal mengenerate gambar");
        } finally {
            setLoadingImage(false);
        }
    };

    const handleAddKeyword = () => {
        if (keywordInput.trim() && !keywords.includes(keywordInput.trim())) {
            setKeywords([...keywords, keywordInput.trim()]);
            setKeywordInput("");
        }
    };

    const handleRemoveKeyword = (keyword: string) => {
        setKeywords(keywords.filter(k => k !== keyword));
    };

    const handleProductToggle = (product: string) => {
        if (postToAll) return;
        setSelectedProducts(prev =>
            prev.includes(product) ? prev.filter(p => p !== product) : [...prev, product]
        );
    };

    const handleSelectAll = () => {
        if (postToAll) {
            setPostToAll(false);
            setSelectedProducts([]);
            toast.info("Semua produk dibatalkan");
        } else {
            setPostToAll(true);
            setSelectedProducts([...productNames]);
            toast.success("Semua produk dipilih");
        }
    };

    const closeResultDialog = () => {
        setShowResultDialog(false);
        setPublishResults(null);
    };

    const saveAsDraft = async () => {
        if (!article) {
            toast.error("Generate artikel terlebih dahulu");
            return false;
        }

        const draftData = {
            Title: topic,
            Topic: topic,
            Article: article,
            ImageUrl: imageUrl || undefined,
            ImagePrompt: topic,
            TargetProducts: selectedProducts,
            hasImage: !!imageUrl,
        };

        try {
            if (currentDraftId) {
                await updateDraft(currentDraftId, { ...draftData, status: 'draft' });
                toast.success("Draft diperbarui!");
            } else {
                const response = await createDraft(draftData);
                setCurrentDraftId(response.id);
                toast.success("Draft tersimpan!");
            }
            return true;
        } catch (error) {
            toast.error("Gagal menyimpan draft");
            return false;
        }
    };

    const saveAsSchedule = async () => {
        if (!article) {
            toast.error("Generate artikel terlebih dahulu");
            return false;
        }
        if (selectedProducts.length === 0) {
            toast.error("Pilih minimal 1 produk");
            return false;
        }

        let scheduledFor: string;
        if (dailySchedule) {
            scheduledFor = `daily:${dailyTime}`;
        } else {
            if (!scheduleDate) {
                toast.error("Pilih tanggal jadwal");
                return false;
            }
            scheduledFor = `${scheduleDate}T${scheduleTime}:00`;
        }

        setPublishing(true);
        try {
            let response: PublishResponse;

            if (currentDraftId) {
                response = await publishDraft(currentDraftId, scheduledFor);
            } else {
                const draftData: ScheduleRequest = {
                    Title: topic,
                    Topic: topic,
                    Article: article,
                    ImageURL: imageUrl || undefined,
                    ImagePrompt: topic,
                    TargetProducts: selectedProducts,
                    ScheduledFor: scheduledFor,
                    HasImage: !!imageUrl,
                };
                response = await draftSchedule(draftData);
            }

            const hasErrors = response.results?.some(r => !r.success);

            if (hasErrors) {
                setPublishResults(response);
                setShowResultDialog(true);
                toast.error("Publikasi sebagian gagal", {
                    description: "Beberapa produk tidak dapat dijangkau. Lihat detail untuk info lebih lanjut."
                });
            } else {
                toast.success("Berhasil dijadwalkan!", {
                    description: dailySchedule
                        ? `"${topic}" akan diposting setiap hari jam ${dailyTime}`
                        : `"${topic}" dijadwalkan pada ${scheduleDate} jam ${scheduleTime}`,
                });
                resetForm();
            }
            return true;
        } catch (error: any) {
            console.error("Failed to save schedule:", error);

            if (error?.results) {
                setPublishResults({
                    message: error.message || "Publikasi gagal",
                    results: error.results,
                    status: "failed"
                });
                setShowResultDialog(true);
            } else {
                toast.error("Gagal menjadwalkan", {
                    description: error?.message || "Terjadi kesalahan pada server",
                });
            }
            return false;
        } finally {
            setPublishing(false);
        }
    };

    useEffect(() => {
        if (
            products.length > 0 &&
            !selectedProducts
        ) {
            setSelectedProducts(
                products[0].id as any
            );
        }
    }, [products]);


    // Buat fungsi helper di luar komponen
    const delay = (ms) => new Promise(resolve => setTimeout(resolve, ms));

    const quickGenerate = async () => {
        if (!topic) {
            toast.error("Generate artikel terlebih dahulu");
            return;
        }

        if (selectedProducts.length === 0) {
            toast.error("Pilih minimal 1 produk");
            return;
        }

        try {
            const data = await getAPIKeys();
            setModels(data as any);

            const availableModels = data
            console.log(availableModels);

            if (availableModels.length < 2) {
                toast.error("Model belum cukup");
                return;
            }

            const articleModel = availableModels.find(model => model.service === 'text');
            const imageModel = availableModels.find(model => model.service === 'image');

            if (!articleModel) {
                toast.error("Model untuk generate artikel tidak ditemukan");
                return;
            }

            if (!imageModel) {
                toast.error("Model untuk generate gambar tidak ditemukan");
                return;
            }

            // STEP 1: Generate article content
            const articleResult = await generateArticleContent(articleModel.id);
            const generatedArticle = articleResult?.content || article;
            const imagePrompt = `Ilustrasi tentang ${topic}`;

            // DELAY 1: Delay 2 detik setelah generate artikel
            await delay(2000); // 2 detik

            // STEP 2: Generate image
            const requestImage = {
                prompt: imagePrompt,
                modelId: imageModel.id
            };

            const imageResult = await generateImage(requestImage);
            const generatedImageUrl = imageResult?.imageUrl || imageUrl;

            // DELAY 2: Delay 1 detik setelah generate gambar
            await delay(1000); // 1 detik

            // STEP 3: Create draft
            const draftData: CreateDraftRequest = {
                Title: topic,
                Topic: topic,
                Article: generatedArticle,
                ImageURL: generatedImageUrl || undefined,
                ImagePrompt: imagePrompt,
                TargetProducts: selectedProducts,
            };

            const draft = await createDraft(draftData);

            // DELAY 3: Delay 1.5 detik sebelum publish
            await delay(1500); // 1.5 detik

            // STEP 4: Publish
            await publishDraft(draft.id);

            toast.success("Publish berhasil");
            return draft.id;

        } catch (error: any) {
            console.error(error);
            // ... error handling tetap sama
        }
    };


    const postInstant = async () => {
        if (!article) {
            toast.error("Generate artikel terlebih dahulu");
            return false;
        }
        if (selectedProducts.length === 0) {
            toast.error("Pilih minimal 1 produk");
            return false;
        }

        setPublishing(true);
        try {
            let response: PublishResponse;

            if (currentDraftId) {
                response = await publishDraft(currentDraftId);
            } else {
                const draftData = {
                    Title: topic,
                    Topic: topic,
                    Article: article,
                    ImageUrl: imageUrl || undefined,
                    ImagePrompt: topic,
                    TargetProducts: selectedProducts,
                };

                response = await publishDraftInstant(draftData);
            }

            const hasErrors = response.results?.some(r => !r.success);

            if (hasErrors) {
                setPublishResults(response);
                setShowResultDialog(true);
                toast.error("Publikasi sebagian gagal", {
                    description: "Beberapa produk tidak dapat dijangkau. Lihat detail untuk info lebih lanjut."
                });
            } else {
                toast.success("Berhasil diposting!", {
                    description: `"${topic}" telah diposting ke ${selectedProducts.length} produk`,
                });
                resetForm();
            }
            return true;
        } catch (error: any) {
            console.error("Failed to post:", error);

            if (error?.results) {
                setPublishResults({
                    message: error.message || "Publikasi gagal",
                    results: error.results,
                    status: "failed"
                });
                setShowResultDialog(true);
            } else {
                toast.error("Gagal memposting", {
                    description: error?.message || "Terjadi kesalahan pada server",
                });
            }
            return false;
        } finally {
            setPublishing(false);
        }
    };

    const resetForm = () => {
        setTopic("");
        setArticle("");
        setArticleResponse(null);
        setImageUrl("");
        setSelectedProducts([]);
        setPostMode("instant");
        setCurrentDraftId(null);
        setKeywords([]);
        setKeywordInput("");
        setSeoScore(null);
        setReadabilityScore(null);
        setWordCount(null);
        setTone("professional");
        setArticleLength("medium");
        setLanguage("id");
    };

    const [isPosting, setIsPosting] = useState(false);

    const handlePost = async () => {
        try {
            if (selectedProducts.length === 0) {
                toast.error("Pilih minimal 1 produk terlebih dahulu");
                return;
            }
            if (!article) {
                toast.error("Generate artikel terlebih dahulu");
                return;
            }

            console.log(postMode)

            setIsPosting(true)
            if (postMode === "draft") {
                await saveAsDraft();
            } else if (postMode === "scheduled") {
                await saveAsSchedule();
            } else {
                await postInstant();
            }
        } finally {
            setIsPosting(false)
        }
    };

    return {
        topic, setTopic,
        article, setArticle,
        articleResponse,
        imageUrl, setImageUrl,
        loadingArticle, loadingImage,
        publishing,
        selectedProducts, setSelectedProducts,
        postMode, setPostMode,
        scheduleTime, setScheduleTime,
        scheduleDate, setScheduleDate,
        dailySchedule, setDailySchedule,
        dailyTime, setDailyTime,
        autoGenerateImage, setAutoGenerateImage,
        postToAll, setPostToAll,
        currentDraftId,
        models,
        selectedModelId, setSelectedModelId,
        loadingModels,
        products, productNames, productsLoading, productsError,
        tone, setTone,
        articleLength, setArticleLength,
        keywords, keywordInput, setKeywordInput,
        language, setLanguage,
        seoScore, readabilityScore, wordCount,
        showResultDialog, publishResults, closeResultDialog,
        generateArticle: generateArticleContent,
        generateImage: generateImageManually,
        handleProductToggle, handleSelectAll, handlePost,
        saveAsDraft, resetForm,
        handleAddKeyword, handleRemoveKeyword,
        isPosting, setIsPosting, quickGenerate
    };
}