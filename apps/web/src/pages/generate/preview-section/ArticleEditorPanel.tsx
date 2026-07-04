import { useRef, useCallback } from "react";
import { FileText, Pencil, Save, Eye, Plus, Trash2 } from "lucide-react";
import { cn } from "@/lib/utils";
import ReactQuill from "react-quill";
import "react-quill/dist/quill.snow.css";
import { type Paragraph, wrapParagraph, getTypeStyles } from "./types";
import { getModulesForType, getFormatsForType, createImageHandler } from "./QuilConfig";
import { useArticleEditor } from "./useArticleEditor";

interface ArticleEditorPanelProps {
    article: string;
    onArticleUpdate?: (newArticle: string) => void;
}

const imageHandler = createImageHandler();

export function ArticleEditorPanel({ article, onArticleUpdate }: ArticleEditorPanelProps) {
    const {
        isEditing,
        paragraphs,
        isSaving,
        toggleEdit,
        handleParagraphChange,
        changeParagraphType,
        addParagraph,
        removeParagraph,
        handleSave,
    } = useArticleEditor({ article, onArticleUpdate });

    const editorRefs = useRef<Map<string, any>>(new Map());
    const styledEditorsRef = useRef<Set<string>>(new Set());

    // Setup Quill ref dengan custom image handler
    const setupQuillRef = useCallback((paragraphId: string) => (el: any) => {
        if (el) {
            const quill = el.getEditor();
            editorRefs.current.set(paragraphId, quill);

            const toolbar = quill.getModule('toolbar');
            if (toolbar) {
                toolbar.addHandler('image', () => {
                    imageHandler(quill);
                });
            }
        } else {
            editorRefs.current.delete(paragraphId);
            styledEditorsRef.current.delete(paragraphId);
        }
    }, []);

    const handleRemoveParagraph = (paragraphId: string) => {
        removeParagraph(paragraphId, (id) => styledEditorsRef.current.delete(id));
    };

    return (
        <div
            className={cn(
                "rounded-xl border p-5 min-h-[200px]",
                "bg-white border-slate-200/80",
                "dark:bg-[#0f0d1a] dark:border-white/[0.06]"
            )}
        >
            <div className="flex items-center justify-between mb-3">
                <span className="text-xs text-slate-400 dark:text-slate-600">
                    {isEditing ? "✏️ Mode Edit" : "👁️ Mode Baca"}
                </span>
                <div className="flex gap-2">
                    {isEditing ? (
                        <>
                            <button
                                onClick={handleSave}
                                disabled={isSaving}
                                className={cn(
                                    "inline-flex items-center gap-1.5 rounded-md px-3 py-1.5 text-xs font-medium transition-colors",
                                    "bg-emerald-600 hover:bg-emerald-700 text-white",
                                    "dark:bg-emerald-500 dark:hover:bg-emerald-600",
                                    `${isSaving && "opacity-50 cursor-not-allowed"}`
                                )}
                            >
                                <Save className="h-3 w-3" />
                                {isSaving ? "Menyimpan..." : "Simpan"}
                            </button>
                            <button
                                onClick={toggleEdit}
                                className="inline-flex items-center gap-1.5 rounded-md px-3 py-1.5 text-xs font-medium transition-colors bg-slate-100 hover:bg-slate-200 text-slate-700 dark:bg-white/[0.06] dark:hover:bg-white/[0.12] dark:text-slate-300"
                            >
                                <Eye className="h-3 w-3" />
                                Batal
                            </button>
                        </>
                    ) : (
                        <button
                            onClick={toggleEdit}
                            className="inline-flex items-center gap-1.5 rounded-md px-3 py-1.5 text-xs font-medium transition-colors bg-blue-600 hover:bg-blue-700 text-white dark:bg-blue-500 dark:hover:bg-blue-600"
                        >
                            <Pencil className="h-3 w-3" />
                            Edit Artikel
                        </button>
                    )}
                </div>
            </div>

            {paragraphs.length > 0 ? (
                <>
                    {isEditing ? (
                        <div className="space-y-3">
                            <style>{`
                                .paragraph-editor .ql-editor {
                                    word-wrap: break-word;
                                    overflow-wrap: break-word;
                                    padding: 8px 12px !important;
                                    min-height: 60px !important;
                                }
                                .paragraph-editor .ql-editor img {
                                    max-width: 100% !important;
                                    height: auto !important;
                                    border-radius: 8px !important;
                                    margin: 8px 0 !important;
                                }
                                .paragraph-editor[data-type="h1"] .ql-editor { font-size: 2.25rem !important; font-weight: 700 !important; line-height: 1.2 !important; }
                                .paragraph-editor[data-type="h2"] .ql-editor { font-size: 1.875rem !important; font-weight: 700 !important; line-height: 1.3 !important; }
                                .paragraph-editor[data-type="h3"] .ql-editor { font-size: 1.5rem !important; font-weight: 600 !important; line-height: 1.4 !important; }
                                .paragraph-editor[data-type="h4"] .ql-editor { font-size: 1.25rem !important; font-weight: 600 !important; line-height: 1.4 !important; }
                                .paragraph-editor[data-type="h5"] .ql-editor { font-size: 1.125rem !important; font-weight: 600 !important; line-height: 1.5 !important; }
                                .paragraph-editor[data-type="h6"] .ql-editor { font-size: 1rem !important; font-weight: 600 !important; line-height: 1.5 !important; }
                                .paragraph-editor[data-type="p"] .ql-editor { font-size: 1rem !important; font-weight: 400 !important; line-height: 1.6 !important; }
                                .paragraph-editor[data-type="blockquote"] .ql-editor { font-size: 1rem !important; font-weight: 400 !important; line-height: 1.6 !important; font-style: italic !important; border-left: 4px solid #e2e8f0 !important; padding-left: 16px !important; }
                                .paragraph-editor[data-type="ul"] .ql-editor,
                                .paragraph-editor[data-type="ol"] .ql-editor { font-size: 1rem !important; font-weight: 400 !important; line-height: 1.6 !important; }
                                .paragraph-editor[data-type="pre"] .ql-editor { font-size: 0.875rem !important; font-weight: 400 !important; line-height: 1.5 !important; font-family: 'Courier New', monospace !important; background: #f1f5f9 !important; }
                                .paragraph-editor .ql-toolbar { border-radius: 6px 6px 0 0; background: #f8fafc; border-color: #e2e8f0; }
                                .paragraph-editor .ql-container { border-radius: 0 0 6px 6px; border-color: #e2e8f0; }
                                .dark .paragraph-editor .ql-toolbar { background: rgba(255,255,255,0.05); border-color: rgba(255,255,255,0.1); }
                                .dark .paragraph-editor .ql-container { border-color: rgba(255,255,255,0.1); }
                                .dark .paragraph-editor .ql-editor { color: #e2e8f0 !important; }
                                .ql-snow .ql-toolbar button.ql-image:hover .ql-stroke,
                                .ql-snow .ql-toolbar button.ql-image.ql-active .ql-stroke {
                                    stroke: #3b82f6 !important;
                                }
                            `}</style>

                            {paragraphs.map((paragraph, index) => {
                                const styles = getTypeStyles(paragraph.type);
                                return (
                                    <div
                                        key={paragraph.id}
                                        className={cn(
                                            "paragraph-editor relative group rounded-lg border transition-colors",
                                            styles.border,
                                            styles.bg
                                        )}
                                        data-type={paragraph.type}
                                    >
                                        <div className="flex items-center justify-between px-3 py-2 border-b border-slate-200/60 dark:border-slate-700">
                                            <div className="flex items-center gap-2">
                                                <span className={cn("text-[10px] font-medium uppercase tracking-wider", styles.label)}>
                                                    {paragraph.label}
                                                </span>

                                                <select
                                                    value={paragraph.type}
                                                    onChange={(e) => changeParagraphType(paragraph.id, e.target.value as Paragraph['type'])}
                                                    className="text-[10px] border border-slate-200 dark:border-slate-600 rounded px-1.5 py-0.5 bg-white dark:bg-slate-800 text-slate-600 dark:text-slate-400 focus:outline-none focus:ring-1 focus:ring-blue-500 cursor-pointer"
                                                >
                                                    <option value="h1">H1 - Judul</option>
                                                    <option value="h2">H2 - Sub Judul</option>
                                                    <option value="h3">H3 - Sub Judul</option>
                                                    <option value="h4">H4</option>
                                                    <option value="h5">H5</option>
                                                    <option value="h6">H6</option>
                                                    <option value="p">Paragraf</option>
                                                    <option value="blockquote">Kutipan</option>
                                                    <option value="ul">List</option>
                                                    <option value="ol">List Bernomor</option>
                                                    <option value="pre">Code Block</option>
                                                </select>
                                            </div>

                                            <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                                                <button
                                                    onClick={() => addParagraph(index, 'p')}
                                                    className="p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-700 text-slate-600 dark:text-slate-400"
                                                    title="Tambah paragraf setelah ini"
                                                >
                                                    <Plus className="h-3 w-3" />
                                                </button>
                                                {paragraphs.length > 1 && (
                                                    <button
                                                        onClick={() => handleRemoveParagraph(paragraph.id)}
                                                        className="p-1 rounded hover:bg-red-100 dark:hover:bg-red-900/30 text-red-600 dark:text-red-400"
                                                        title="Hapus paragraf"
                                                    >
                                                        <Trash2 className="h-3 w-3" />
                                                    </button>
                                                )}
                                            </div>
                                        </div>

                                        <div className="p-3">
                                            <ReactQuill
                                                key={`${paragraph.id}::${paragraph.type}`}
                                                ref={setupQuillRef(paragraph.id)}
                                                theme="snow"
                                                defaultValue={paragraph.innerHTML}
                                                onChange={(value) => handleParagraphChange(paragraph.id, value)}
                                                modules={getModulesForType(paragraph.type)}
                                                formats={getFormatsForType(paragraph.type)}
                                                placeholder={`Tulis ${paragraph.label.toLowerCase()}...`}
                                                preserveWhitespace
                                            />
                                        </div>
                                    </div>
                                );
                            })}

                            <div className="flex gap-2 flex-wrap pt-2">
                                <button onClick={() => addParagraph(paragraphs.length - 1, 'h1')} className="px-3 py-1.5 border border-dashed border-slate-300 dark:border-slate-600 rounded-lg text-slate-500 dark:text-slate-400 hover:border-slate-400 dark:hover:border-slate-500 hover:text-slate-700 dark:hover:text-slate-300 text-xs font-medium transition-colors">
                                    + H1
                                </button>
                                <button onClick={() => addParagraph(paragraphs.length - 1, 'h2')} className="px-3 py-1.5 border border-dashed border-slate-300 dark:border-slate-600 rounded-lg text-slate-500 dark:text-slate-400 hover:border-slate-400 dark:hover:border-slate-500 hover:text-slate-700 dark:hover:text-slate-300 text-xs font-medium transition-colors">
                                    + H2
                                </button>
                                <button onClick={() => addParagraph(paragraphs.length - 1, 'h3')} className="px-3 py-1.5 border border-dashed border-slate-300 dark:border-slate-600 rounded-lg text-slate-500 dark:text-slate-400 hover:border-slate-400 dark:hover:border-slate-500 hover:text-slate-700 dark:hover:text-slate-300 text-xs font-medium transition-colors">
                                    + H3
                                </button>
                                <button onClick={() => addParagraph(paragraphs.length - 1, 'p')} className="px-3 py-1.5 border border-dashed border-slate-300 dark:border-slate-600 rounded-lg text-slate-500 dark:text-slate-400 hover:border-slate-400 dark:hover:border-slate-500 hover:text-slate-700 dark:hover:text-slate-300 text-xs font-medium transition-colors">
                                    + Paragraf
                                </button>
                                <button onClick={() => addParagraph(paragraphs.length - 1, 'blockquote')} className="px-3 py-1.5 border border-dashed border-slate-300 dark:border-slate-600 rounded-lg text-slate-500 dark:text-slate-400 hover:border-slate-400 dark:hover:border-slate-500 hover:text-slate-700 dark:hover:text-slate-300 text-xs font-medium transition-colors">
                                    + Kutipan
                                </button>
                                <button onClick={() => addParagraph(paragraphs.length - 1, 'ul')} className="px-3 py-1.5 border border-dashed border-slate-300 dark:border-slate-600 rounded-lg text-slate-500 dark:text-slate-400 hover:border-slate-400 dark:hover:border-slate-500 hover:text-slate-700 dark:hover:text-slate-300 text-xs font-medium transition-colors">
                                    + List
                                </button>
                            </div>
                        </div>
                    ) : (
                        <div
                            className="prose dark:prose-invert max-w-none prose-sm prose-headings:font-semibold prose-p:text-slate-600 dark:prose-p:text-slate-400"
                            dangerouslySetInnerHTML={{ __html: paragraphs.map(p => wrapParagraph(p)).join('\n') }}
                            style={{ wordWrap: 'break-word', overflowWrap: 'break-word', maxWidth: '100%' }}
                        />
                    )}
                </>
            ) : (
                <div className="flex flex-col items-center justify-center py-12 text-slate-400">
                    <FileText className="h-8 w-8 mb-2 opacity-30" />
                    <p className="text-sm">
                        Belum ada artikel. Klik "Generate Artikel" dulu.
                    </p>
                </div>
            )}
        </div>
    );
}
