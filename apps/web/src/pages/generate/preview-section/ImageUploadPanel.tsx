import { useState, useRef } from "react";
import { ImageIcon, Upload, X, Camera } from "lucide-react";
import { cn } from "@/lib/utils";

interface ImageUploadPanelProps {
    hasImage: boolean;
    uploadedImage: string | null;
    onUploadedImageChange: (image: string | null) => void;
    onImageUpload?: (base64Image: string) => void;
}

const VALID_IMAGE_TYPES = ['image/jpeg', 'image/png', 'image/gif', 'image/webp', 'image/svg+xml'];
const MAX_IMAGE_SIZE = 5 * 1024 * 1024;

const convertToBase64 = (file: File): Promise<string> => {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = () => resolve(reader.result as string);
        reader.onerror = (error) => reject(error);
    });
};

export function ImageUploadPanel({
    hasImage,
    uploadedImage,
    onUploadedImageChange,
    onImageUpload,
}: ImageUploadPanelProps) {
    const [isUploading, setIsUploading] = useState(false);
    const [dragActive, setDragActive] = useState(false);
    const fileInputRef = useRef<HTMLInputElement>(null);

    const handleFileUpload = async (file: File) => {
        if (!VALID_IMAGE_TYPES.includes(file.type)) {
            alert('Format file tidak didukung. Gunakan JPG, PNG, GIF, WEBP, atau SVG.');
            return;
        }

        if (file.size > MAX_IMAGE_SIZE) {
            alert('Ukuran file terlalu besar. Maksimal 5MB.');
            return;
        }

        setIsUploading(true);
        try {
            const base64 = await convertToBase64(file);
            onUploadedImageChange(base64);
            onImageUpload?.(base64);
        } catch (error) {
            console.error('Error uploading image:', error);
            alert('Gagal mengupload gambar. Silakan coba lagi.');
        } finally {
            setIsUploading(false);
        }
    };

    const handleDrag = (e: React.DragEvent) => {
        e.preventDefault();
        e.stopPropagation();
        if (e.type === 'dragenter' || e.type === 'dragover') {
            setDragActive(true);
        } else if (e.type === 'dragleave') {
            setDragActive(false);
        }
    };

    const handleDrop = (e: React.DragEvent) => {
        e.preventDefault();
        e.stopPropagation();
        setDragActive(false);

        const files = e.dataTransfer.files;
        if (files && files.length > 0) {
            handleFileUpload(files[0]);
        }
    };

    const handleFileInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const files = e.target.files;
        if (files && files.length > 0) {
            handleFileUpload(files[0]);
        }
        if (fileInputRef.current) {
            fileInputRef.current.value = '';
        }
    };

    const handleRemoveImage = () => {
        onUploadedImageChange(null);
        onImageUpload?.('');
    };

    const triggerFileInput = () => {
        fileInputRef.current?.click();
    };

    return (
        <div
            className={cn(
                "rounded-xl border p-5 min-h-[200px]",
                "bg-white border-slate-200/80",
                "dark:bg-[#0f0d1a] dark:border-white/[0.06]"
            )}
        >
            <div className="flex gap-3 mb-4">
                <button
                    onClick={triggerFileInput}
                    disabled={isUploading}
                    className={cn(
                        "inline-flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-colors",
                        "bg-blue-600 hover:bg-blue-700 text-white",
                        "dark:bg-blue-500 dark:hover:bg-blue-600",
                        "disabled:opacity-50 disabled:cursor-not-allowed"
                    )}
                >
                    <Camera className="h-4 w-4" />
                    {isUploading ? "Mengupload..." : "Pilih Gambar"}
                </button>

                <div
                    className={cn(
                        "flex-1 flex items-center justify-center px-4 py-2 rounded-lg border-2 border-dashed transition-colors",
                        dragActive ? "border-blue-500 bg-blue-50/50 dark:bg-blue-950/30" : "border-slate-300 dark:border-slate-700",
                        "hover:border-slate-400 dark:hover:border-slate-600"
                    )}
                    onDragEnter={handleDrag}
                    onDragLeave={handleDrag}
                    onDragOver={handleDrag}
                    onDrop={handleDrop}
                >
                    <div className="flex items-center gap-2 text-sm text-slate-500 dark:text-slate-400">
                        <Upload className="h-4 w-4" />
                        <span>atau drag & drop gambar di sini</span>
                    </div>
                </div>
            </div>

            <input
                ref={fileInputRef}
                type="file"
                accept={VALID_IMAGE_TYPES.join(',')}
                onChange={handleFileInputChange}
                className="hidden"
                disabled={isUploading}
            />

            <div className="mb-4 text-xs text-slate-400 dark:text-slate-600">
                Supported formats: JPG, PNG, GIF, WEBP, SVG (max 5MB)
            </div>

            {uploadedImage ? (
                <div className="space-y-3">
                    <div className="relative">
                        <img
                            src={uploadedImage}
                            alt="Preview"
                            className="w-full rounded-lg border border-slate-200/80 object-cover max-h-[400px] dark:border-white/[0.06]"
                            onError={(e) => {
                                (e.target as HTMLImageElement).src =
                                    "https://placehold.co/800x400?text=Image+Failed";
                            }}
                        />
                        <button
                            onClick={handleRemoveImage}
                            className="absolute top-2 right-2 p-1.5 rounded-full bg-red-500 hover:bg-red-600 text-white transition-colors shadow-lg"
                            title="Hapus gambar"
                        >
                            <X className="h-4 w-4" />
                        </button>
                    </div>
                    <div className="flex items-center justify-between">
                        <p className="text-xs text-slate-400 dark:text-slate-600">
                            Klik kanan → Save Image As untuk menyimpan
                        </p>
                        <span className="text-xs text-emerald-600 dark:text-emerald-400">
                            ✓ Gambar tersedia
                        </span>
                    </div>
                </div>
            ) : (
                <div className="flex flex-col items-center justify-center py-12 text-slate-400">
                    <ImageIcon className="h-8 w-8 mb-2 opacity-30" />
                    <p className="text-sm">
                        {hasImage ? "Generate gambar terlebih dahulu atau upload gambar" : "Belum ada gambar. Upload gambar di atas"}
                    </p>
                </div>
            )}
        </div>
    );
}
