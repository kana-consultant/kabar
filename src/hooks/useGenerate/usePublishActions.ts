import { toast } from "sonner";
import { createDraft, updateDraft, publishDraft, draftSchedule, publishDraftInstant } from "@/services/draft";

export async function saveAsDraft(
    article: string,
    topic: string,
    imageUrl: string,
    selectedProducts: string[],
    currentDraftId: string | null,
    setCurrentDraftId: (id: string | null) => void
) {
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
    resetForm: () => void
) {
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
        let response: any;

        if (currentDraftId) {
            response = await publishDraft(currentDraftId, scheduledFor);
        } else {
            const draftData = {
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

        const hasErrors = response.results?.some((r: any) => !r.success);

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
    resetForm: () => void
) {
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
        let response: any;

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

        const hasErrors = response.results?.some((r: any) => !r.success);

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
}