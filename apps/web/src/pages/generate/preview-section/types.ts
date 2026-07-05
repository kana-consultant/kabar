export interface Paragraph {
    id: string;
    innerHTML: string;
    type: 'p' | 'h1' | 'h2' | 'h3' | 'h4' | 'h5' | 'h6' | 'blockquote' | 'ul' | 'ol' | 'pre';
    label: string;
}

export interface PreviewSectionProps {
    article: string;
    keywords : string[];
    excerpt : string;
    slug : string;
    imageUrl: string;
    hasImage: boolean;
    postMode: "instant" | "scheduled" | "draft";
    dailySchedule: boolean;
    dailyTime: string;
    scheduleDate: string;
    scheduleTime: string;
    selectedProductsCount: number;
    autoGenerateImage: boolean;
    onArticleUpdate?: (newArticle: string) => void;
    onImageUpload?: (base64Image: string) => void;
}

export const LABEL_MAP: Record<Paragraph['type'], string> = {
    h1: 'Judul H1',
    h2: 'Sub Judul H2',
    h3: 'Sub Judul H3',
    h4: 'Sub Judul H4',
    h5: 'Sub Judul H5',
    h6: 'Sub Judul H6',
    blockquote: 'Kutipan',
    ul: 'List',
    ol: 'List Bernomor',
    pre: 'Code Block',
    p: 'Paragraf',
};

let paragraphIdCounter = 0;
export const nextParagraphId = () => `para-${Date.now()}-${paragraphIdCounter++}`;

export const isListType = (type: Paragraph['type']) => type === 'ul' || type === 'ol';

export const wrapParagraph = (p: Pick<Paragraph, 'type' | 'innerHTML'>): string =>
    isListType(p.type) ? p.innerHTML : `<${p.type}>${p.innerHTML}</${p.type}>`;

export const stripQuillDefaultWrapper = (html: string): string => {
    const temp = document.createElement('div');
    temp.innerHTML = html;
    if (temp.children.length === 1 && temp.children[0].tagName === 'P') {
        return temp.children[0].innerHTML;
    }
    return html;
};

export const getTypeStyles = (type: Paragraph['type']) => {
    const styles = {
        p: { border: 'border-slate-200 dark:border-slate-700', bg: 'bg-white dark:bg-slate-900', label: 'text-slate-500 dark:text-slate-400' },
        h1: { border: 'border-blue-300 dark:border-blue-800', bg: 'bg-blue-50/50 dark:bg-blue-950/30', label: 'text-blue-600 dark:text-blue-400' },
        h2: { border: 'border-indigo-300 dark:border-indigo-800', bg: 'bg-indigo-50/50 dark:bg-indigo-950/30', label: 'text-indigo-600 dark:text-indigo-400' },
        h3: { border: 'border-violet-300 dark:border-violet-800', bg: 'bg-violet-50/50 dark:bg-violet-950/30', label: 'text-violet-600 dark:text-violet-400' },
        h4: { border: 'border-purple-300 dark:border-purple-800', bg: 'bg-purple-50/50 dark:bg-purple-950/30', label: 'text-purple-600 dark:text-purple-400' },
        h5: { border: 'border-pink-300 dark:border-pink-800', bg: 'bg-pink-50/50 dark:bg-pink-950/30', label: 'text-pink-600 dark:text-pink-400' },
        h6: { border: 'border-rose-300 dark:border-rose-800', bg: 'bg-rose-50/50 dark:bg-rose-950/30', label: 'text-rose-600 dark:text-rose-400' },
        blockquote: { border: 'border-amber-300 dark:border-amber-800 border-l-4', bg: 'bg-amber-50/50 dark:bg-amber-950/30', label: 'text-amber-600 dark:text-amber-400' },
        ul: { border: 'border-emerald-300 dark:border-emerald-800', bg: 'bg-emerald-50/50 dark:bg-emerald-950/30', label: 'text-emerald-600 dark:text-emerald-400' },
        ol: { border: 'border-teal-300 dark:border-teal-800', bg: 'bg-teal-50/50 dark:bg-teal-950/30', label: 'text-teal-600 dark:text-teal-400' },
        pre: { border: 'border-slate-400 dark:border-slate-600', bg: 'bg-slate-100 dark:bg-slate-800', label: 'text-slate-600 dark:text-slate-400' },
    };
    return styles[type] || styles.p;
};