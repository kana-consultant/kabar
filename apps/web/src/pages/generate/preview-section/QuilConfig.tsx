import { Quill } from "react-quill";
import type { Paragraph } from "./types";
import { isListType } from "./types";

// Register custom image handler ke Quill
const ImageFormat = Quill.import('formats/image');
class CustomImageBlot extends ImageFormat {
    static create(value: any) {
        const node = super.create(value);
        if (typeof value === 'string') {
            node.setAttribute('src', value);
        }
        return node;
    }
}
CustomImageBlot.blotName = 'image';
CustomImageBlot.tagName = 'img';
Quill.register(CustomImageBlot, true);

// Toolbar configuration tanpa format custom
const BASE_TOOLBAR = [
    ["bold", "italic", "underline", "strike"],
    ["link", "image"],
    ["clean"],
];

const LIST_TOOLBAR = [
    ["bold", "italic", "underline", "strike"],
    [{ list: "ordered" }, { list: "bullet" }],
    ["link", "image"],
    ["clean"],
];

const BASE_FORMATS = ["bold", "italic", "underline", "strike", "link", "image"];
const LIST_FORMATS = [...BASE_FORMATS, "list"];

export const getModulesForType = (type: Paragraph['type']) => ({
    toolbar: isListType(type) ? LIST_TOOLBAR : BASE_TOOLBAR,
    clipboard: {
        matchVisual: false
    }
});

export const getFormatsForType = (type: Paragraph['type']) => (isListType(type) ? LIST_FORMATS : BASE_FORMATS);

const VALID_IMAGE_TYPES = ['image/jpeg', 'image/png', 'image/gif', 'image/webp', 'image/svg+xml'];
const MAX_IMAGE_SIZE = 5 * 1024 * 1024;

// Handler untuk tombol image di toolbar Quill (dipakai per-editor paragraf)
export const createImageHandler = () => (quillInstance: any) => {
    const input = document.createElement('input');
    input.setAttribute('type', 'file');
    input.setAttribute('accept', VALID_IMAGE_TYPES.join(','));
    input.click();

    input.onchange = () => {
        const file = input.files?.[0];
        if (!file) return;

        if (!VALID_IMAGE_TYPES.includes(file.type)) {
            alert('Format file tidak didukung. Gunakan JPG, PNG, GIF, WEBP, atau SVG.');
            return;
        }

        if (file.size > MAX_IMAGE_SIZE) {
            alert('Ukuran file terlalu besar. Maksimal 5MB.');
            return;
        }

        const reader = new FileReader();
        reader.onload = (e) => {
            const base64 = e.target?.result as string;
            const range = quillInstance.getSelection(true);
            quillInstance.insertEmbed(range.index, 'image', base64);
            quillInstance.setSelection(range.index + 1);
        };
        reader.readAsDataURL(file);
    };
};