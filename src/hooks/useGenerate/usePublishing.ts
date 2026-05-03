import { toast } from "sonner";
import { publishDraft, draftSchedule, publishDraftInstant, type PublishResponse } from "@/services/draft";
import type { ScheduleRequest } from "@/types/schedule";

interface UsePublishingProps {
    setShowResultDialog: (show: boolean) => void;
    setPublishResults: (results: PublishResponse) => void;
    resetForm: () => void;
    setPublishing: (publishing: boolean) => void;
}

export function usePublishing(props: UsePublishingProps) {
    const saveAsSchedule = async (
        topic: string,
        article: string,
        imageUrl: string,
        selectedProducts: string[],
        currentDraftId: string | null,
        scheduleDate: string,
        scheduleTime: string,
        dailySchedule: boolean,
        dailyTime: string
    ) => {
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

        props.setPublishing(true);
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
                props.setPublishResults(response);
                props.setShowResultDialog(true);
                toast.error("Publikasi sebagian gagal", {
                    description: "Beberapa produk tidak dapat dijangkau. Lihat detail untuk info lebih lanjut."
                });
            } else {
                toast.success("Berhasil dijadwalkan!", {
                    description: dailySchedule
                        ? `"${topic}" akan diposting setiap hari jam ${dailyTime}`
                        : `"${topic}" dijadwalkan pada ${scheduleDate} jam ${scheduleTime}`,
                });
                props.resetForm();
            }
            return true;
        } catch (error: any) {
            console.error("Failed to save schedule:", error);

            if (error?.results) {
                props.setPublishResults({
                    message: error.message || "Publikasi gagal",
                    results: error.results,
                    status: "failed"
                });
                props.setShowResultDialog(true);
            } else {
                toast.error("Gagal menjadwalkan", {
                    description: error?.message || "Terjadi kesalahan pada server",
                });
            }
            return false;
        } finally {
            props.setPublishing(false);
        }
    };

    const postInstant = async (
        topic: string,
        article: string,
        imageUrl: string,
        selectedProducts: string[],
        currentDraftId: string | null,
        resetForm: () => void
    ) => {
        if (!article) {
            toast.error("Generate artikel terlebih dahulu");
            return false;
        }
        if (selectedProducts.length === 0) {
            toast.error("Pilih minimal 1 produk");
            return false;
        }

        props.setPublishing(true);
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
                props.setPublishResults(response);
                props.setShowResultDialog(true);
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
                props.setPublishResults({
                    message: error.message || "Publikasi gagal",
                    results: error.results,
                    status: "failed"
                });
                props.setShowResultDialog(true);
            } else {
                toast.error("Gagal memposting", {
                    description: error?.message || "Terjadi kesalahan pada server",
                });
            }
            return false;
        } finally {
            props.setPublishing(false);
        }
    };

    const quickGenerate = async (
        topic: string,
        selectedProducts: string[],
        getTextAndImageModels: () => Promise<any>,
        generateArticleContent: (modelId: string) => Promise<any>,
        generateImage: (prompt: string, modelId: string) => Promise<any>,
        createDraft: (data: any) => Promise<any>,
        publishDraftFn: (id: string) => Promise<any>
    ) => {
        if (!topic) {
            toast.error("Generate artikel terlebih dahulu");
            return;
        }

        if (selectedProducts.length === 0) {
            toast.error("Pilih minimal 1 produk");
            return;
        }

        const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

        try {
            const { articleModel, imageModel, availableModels } = await getTextAndImageModels();

            if (availableModels.length < 2) {
                toast.error("Model belum cukup");
                return;
            }

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
            const generatedArticle = articleResult?.content || "";
            const imagePrompt = `Ilustrasi tentang ${topic}`;

            await delay(2000);

            // STEP 2: Generate image
            const imageResult = await generateImage(imagePrompt, imageModel.id);
            const generatedImageUrl = imageResult?.imageUrl || "";

            await delay(1000);

            // STEP 3: Create draft
            const draftData = {
                Title: topic,
                Topic: topic,
                Article: generatedArticle,
                ImageURL: generatedImageUrl || undefined,
                ImagePrompt: imagePrompt,
                TargetProducts: selectedProducts,
            };

            const draft = await createDraft(draftData);

            await delay(1500);

            // STEP 4: Publish
            await publishDraftFn(draft.id);

            toast.success("Publish berhasil");
            return draft.id;

        } catch (error: any) {
            console.error(error);
            toast.error("Gagal melakukan quick generate", {
                description: error?.message || "Terjadi kesalahan pada server",
            });
            return;
        }
    };

    return {
        saveAsSchedule,
        postInstant,
        quickGenerate
    };
}