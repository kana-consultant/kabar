// hooks/usePublishActions.ts
import type { ToastContextType } from "@/hooks/use-toast";
import { createDraft, updateDraft, publishDraft, draftSchedule, publishDraftInstant, type Draft, type CreateDraftRequest } from "@/services/draft";
import type { ScheduleRequest } from "@/types/schedule";

export async function saveAsDraft(
    article: string,
    topic: string,
    imageUrl: string,
    selectedProducts: string[],
    currentDraftId: string | null,
    setCurrentDraftId: (id: string | null) => void,
    slug: string,
    tags: string[] | null,
    excerpt: string,
    toast: ToastContextType
): Promise<{ success: boolean; result?: any }> {
    if (!article) {
        toast.error("Generate artikel terlebih dahulu");
        return { success: false };
    }

    const draftData = {
        title: topic,
        topic: topic,
        article: article,
        image_url: imageUrl || undefined,
        image_prompt: topic,
        target_products: selectedProducts,
        has_image: !!imageUrl,
        slug: slug,
        keywords: tags as string[],
        excerpt: excerpt
    };

    try {
        let response;
        if (currentDraftId) {
            await updateDraft(currentDraftId, { ...draftData, status: 'draft' });
            toast.success("Draft diperbarui!");
            return { success: true };
        } else {
            response = await createDraft(draftData);
            setCurrentDraftId(response.id);
            toast.success("Draft tersimpan!");
            return { success: true, result: response };
        }
    } catch (error) {
    
        toast.error("Gagal menyimpan draft");
        return { success: false };
    }
}

export async function saveAsSchedule(
    article: string,
    topic: string,
    imageUrl: string,
    selectedProducts: string[],
    currentDraftId: string | null,
    dailySchedule: boolean,
    dailyTime: string,
    scheduleDate: string,
    scheduleTime: string,
    setPublishing: (val: boolean) => void,
    setPublishResults: (val: any) => void,
    setShowResultDialog: (val: boolean) => void,
    onSuccess: (result: any) => void, // Callback dengan parameter
    onReset: () => void, // Callback untuk reset
    slug: string,
    tags: string[] | null,
    toast: ToastContextType,
    excerpt : string
): Promise<{ success: boolean; result?: any }> {
    if (!article) {
        toast.error("Generate artikel terlebih dahulu");
        return { success: false };
    }
    if (selectedProducts.length === 0) {
        toast.error("Pilih minimal 1 produk");
        return { success: false };
    }

    let scheduledFor: string;

    if (dailySchedule) {
        scheduledFor = `daily:${dailyTime}`;
    } else {
        if (!scheduleDate) {
            toast.error("Pilih tanggal jadwal");
            return { success: false };
        }
        scheduledFor = `${scheduleDate}T${scheduleTime}:00`;
    }

    setPublishing(true);

    try {
        let response: any;

        const draftData: Draft = {
            title: topic,
            topic: topic,
            article: article,
            image_url: imageUrl || undefined,
            image_prompt: topic,
            target_products: selectedProducts,
            scheduled_for: scheduledFor,
            has_image: !!imageUrl,
            keywords: tags as string[],
            slug: slug as string,
            excerpt : excerpt
        };

        if (currentDraftId) {
            response = await publishDraft(currentDraftId, draftData);
        } else {
            const scheduleDraft: ScheduleRequest = {
                title: topic,
                topic: topic,
                article: article,
                image_url: imageUrl || undefined,
                image_prompt: topic,
                target_products: selectedProducts,
                scheduled_for: scheduledFor,
                has_image: !!imageUrl,
                slug: slug,
                keywords: tags as string[]
            };
            response = await draftSchedule(scheduleDraft);
        }

        // Set results
        setPublishResults(response);
        

        // Panggil onSuccess dengan response
        if (onSuccess) {
            onSuccess(response);
        }

        const hasErrors = response.results?.some((r: any) => !r.success);

        if (hasErrors) {
            toast.error("Publikasi sebagian gagal", {
                description: "Beberapa produk tidak dapat dijangkau. Lihat detail untuk info lebih lanjut."
            });
        } else {
            toast.success("Berhasil dijadwalkan!", {
                description: dailySchedule
                    ? `"${topic}" akan diposting setiap hari jam ${dailyTime}`
                    : `"${topic}" dijadwalkan pada ${scheduleDate} jam ${scheduleTime}`,
            });

            // Reset form setelah sukses
            if (onReset) {
                onReset();
            }
        }

        return { success: true, result: response };
    } catch (error: any) {
        
        if (error?.results) {
            const errorResponse = {
                message: error.message || "Publikasi gagal",
                results: error.results,
                status: "failed"
            };
            setPublishResults(errorResponse);
            

            // Panggil onSuccess dengan error response
            if (onSuccess) {
                onSuccess(errorResponse);
            }
        } else {
            toast.error("Gagal menjadwalkan", {
                description: error?.message || "Terjadi kesalahan pada server",
            });
        }

        return { success: false };
    } finally {
        setPublishing(false);
    }
}

export async function postInstant(
    article: string,
    topic: string,
    imageUrl: string,
    selectedProducts: string[],
    currentDraftId: string | null,
    setPublishing: (val: boolean) => void,
    setPublishResults: (val: any) => void,
    setShowResultDialog: (val: boolean) => void,
    onSuccess: (result: any) => void, // Callback dengan parameter
    onReset: () => void, // Callback untuk reset
    slug: string,
    tags: string[],
    excerpt: string,
    toast: ToastContextType
): Promise<{ success: boolean; result?: any }> {
    if (!article) {
        toast.error("Generate artikel terlebih dahulu");
        return { success: false };
    }

    if (selectedProducts.length === 0) {
        toast.error("Pilih minimal 1 produk");
        return { success: false };
    }

    setPublishing(true);

    try {
        const draftData: CreateDraftRequest = {
            title: topic,
            topic: topic,
            article: article,
            image_url: imageUrl || undefined,
            image_prompt: topic,
            target_products: selectedProducts,
            slug: slug,
            keywords: tags as string[],
            excerpt: excerpt
        };

     
        let response: any;
        if (currentDraftId) {
            response = await publishDraft(currentDraftId, draftData);
        } else {
            response = await publishDraftInstant(draftData);
        }

        // Set results
        setPublishResults(response.data);
       

       
        const hasErrors = response.data.results?.some((r: any) => !r.success);

        if (hasErrors) {
            toast.error("Publikasi sebagian gagal", {
                description: "Beberapa produk tidak dapat dijangkau. Lihat detail untuk info lebih lanjut."
            });
        } else {
            onSuccess(response);
            toast.success("Berhasil diposting!", {
                description: `"${topic}" telah diposting ke ${selectedProducts.length} produk`,
            });

            // Reset form setelah sukses
            if (onReset) {
                onReset();
            }
        }
        return { success: true, result: response };
    } catch (error: any) {

        const errorResponse = error.response
        setPublishResults(errorResponse);
        

        // Panggil onSuccess dengan error response
        if (onSuccess) {
            onSuccess(errorResponse);
        }

        throw error
    } finally {
        setPublishing(false);
    }
}