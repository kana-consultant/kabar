import { useState, useRef } from "react";
import {
    Button, Input, Label, Textarea, Switch,
    Select, SelectContent, SelectItem, SelectTrigger, SelectValue,
    Dialog, DialogContent, DialogDescription, DialogFooter,
    DialogHeader, DialogTitle,
} from "@kana-consultant/ui-kit";
import { useToast } from "@/hooks/use-toast";
import {
    Loader2, Plus, GripVertical, X, ChevronRight,
    CheckCircle, Info,
} from "lucide-react";
import { cn } from "@/lib/utils";

// ─── Types ────────────────────────────────────────────────────────────────────

export interface SchemaField {
    key: string;
    type: "string" | "number" | "boolean" | "array" | "object" | "null";
    required: boolean;
}

export interface CreateModelPayload {
    providerId: string;
    modelName: string;
    displayName: string;
    description?: string;
    requestSchema: string; // JSON string
    responsePath?: string;
    imagePath?: string;
    isActive: boolean;
    supportsStreaming: boolean;
    maxTokens?: number;
    temperature?: number;
    // Optional: untuk image models
    modalities?: string[];
    imageConfig?: {
        aspect_ratio?: string;
        output_quality?: string;
    };
}

interface Props {
    open: boolean;
    onOpenChange: (v: boolean) => void;
    /** List of providers from getProviders() */
    providers: { id: string; name: string; displayName: string }[];
    /** Called with the payload after user clicks Save */
    onSave: (payload: CreateModelPayload) => Promise<void>;
}

// ─── Default schema fields ────────────────────────────────────────────────────

const DEFAULT_FIELDS: SchemaField[] = [
    { key: "model",       type: "string",  required: true  },
    { key: "messages",    type: "array",   required: true  },
    { key: "max_tokens",  type: "number",  required: false },
    { key: "temperature", type: "number",  required: false },
    { key: "stream",      type: "boolean", required: false },
];

const FIELD_TYPES = ["string", "number", "boolean", "array", "object", "null"] as const;

// ─── Small helpers ─────────────────────────────────────────────────────────────

const inputCls = cn(
    "h-8 text-sm border-slate-200/80 bg-white placeholder:text-slate-400",
    "dark:border-white/[0.08] dark:bg-white/[0.03] dark:placeholder:text-slate-600",
    "focus-visible:ring-1 focus-visible:ring-green-500/40 focus-visible:border-green-400/60",
    "dark:focus-visible:ring-purple-500/40 dark:focus-visible:border-purple-400/40",
);

const TABS = ["Provider & Model", "Schema Request", "Preview"] as const;
type Tab = 0 | 1 | 2;

// ─── Component ────────────────────────────────────────────────────────────────

export function AddModelDialog({ open, onOpenChange, providers, onSave }: Props) {
    const toast = useToast();

    // — tab —
    const [tab, setTab] = useState<Tab>(0);

    // — form: tab 0 —
    const [providerId,         setProviderId]         = useState("");
    const [modelName,          setModelName]          = useState("");
    const [displayName,        setDisplayName]        = useState("");
    const [description,        setDescription]        = useState("");
    const [responsePath,       setResponsePath]       = useState("choices[0].message.content");
    const [imagePath,          setImagePath]          = useState("");
    const [isActive,           setIsActive]           = useState(true);
    const [supportsStreaming,  setSupportsStreaming]  = useState(true);
    const [maxTokens,          setMaxTokens]          = useState("");
    const [temperature,        setTemperature]        = useState("0.7");
    
    // — image specific —
    const [modalities,         setModalities]         = useState<string[]>(["text"]);
    const [aspectRatio,        setAspectRatio]        = useState("16:9");
    const [outputQuality,      setOutputQuality]      = useState("standard");

    // — form: tab 1 (schema) —
    const [fields,       setFields]       = useState<SchemaField[]>(DEFAULT_FIELDS);
    const [newKey,       setNewKey]        = useState("");
    const [newType,      setNewType]       = useState<SchemaField["type"]>("string");
    const [newRequired,  setNewRequired]   = useState(true);

    // — drag state —
    const dragIndex = useRef<number | null>(null);

    // — saving —
    const [saving, setSaving] = useState(false);

    // ── reset on close ─────────────────────────────────────────────────────────
    const handleOpenChange = (v: boolean) => {
        if (!v) {
            setTab(0);
            setProviderId(""); setModelName(""); setDisplayName("");
            setDescription(""); setResponsePath("choices[0].message.content");
            setImagePath(""); setIsActive(true); setSupportsStreaming(true);
            setMaxTokens(""); setTemperature("0.7");
            setModalities(["text"]); setAspectRatio("16:9"); setOutputQuality("standard");
            setFields(DEFAULT_FIELDS);
            setNewKey(""); setNewType("string"); setNewRequired(true);
        }
        onOpenChange(v);
    };

    // ── drag helpers ───────────────────────────────────────────────────────────
    const onDragStart = (i: number) => { dragIndex.current = i; };
    const onDrop      = (i: number) => {
        const src = dragIndex.current;
        if (src === null || src === i) return;
        const next = [...fields];
        const [moved] = next.splice(src, 1);
        next.splice(i, 0, moved);
        setFields(next);
        dragIndex.current = null;
    };

    // ── add / remove field ─────────────────────────────────────────────────────
    const addField = () => {
        const key = newKey.trim();
        if (!key) { toast.error("Key field tidak boleh kosong"); return; }
        if (fields.find(f => f.key === key)) { toast.error(`Field "${key}" sudah ada`); return; }
        setFields(prev => [...prev, { key, type: newType, required: newRequired }]);
        setNewKey("");
    };

    const removeField = (i: number) => setFields(prev => prev.filter((_, idx) => idx !== i));

    // ── build request schema JSON ──────────────────────────────────────────────
    const buildRequestSchema = (): string => {
        const schema: Record<string, any> = {};
        
        // Add standard fields
        fields.forEach(field => {
            if (field.key === "model") {
                schema[field.key] = "{model}";
            } else if (field.key === "messages") {
                schema[field.key] = [{ role: "user", content: "{prompt}" }];
            } else {
                // For other fields, use placeholder or default value
                if (field.type === "string") schema[field.key] = `{${field.key}}`;
                else if (field.type === "number") schema[field.key] = 0;
                else if (field.type === "boolean") schema[field.key] = false;
                else if (field.type === "array") schema[field.key] = [];
                else if (field.type === "object") schema[field.key] = {};
                else schema[field.key] = null;
            }
        });
        
        // Add modalities if image model
        if (modalities.length > 0 && (modalities.includes("image") || modalities.length > 1)) {
            schema.modalities = modalities;
        }
        
        // Add image config if needed
        if (modalities.includes("image") && (aspectRatio !== "16:9" || outputQuality !== "standard")) {
            schema.image_config = {
                aspect_ratio: aspectRatio,
                output_quality: outputQuality,
            };
        }
        
        return JSON.stringify(schema, null, 2);
    };

    // ── save ───────────────────────────────────────────────────────────────────
    const handleSave = async () => {
        if (!providerId)  { toast.error("Pilih provider");           setTab(0); return; }
        if (!modelName)   { toast.error("Isi Model Name");           setTab(0); return; }
        if (!displayName) { toast.error("Isi Display Name");         setTab(0); return; }
        
        // Validate response path for text models or image path for image models
        const hasImageModality = modalities.includes("image");
        if (!hasImageModality && !responsePath) {
            toast.error("Isi Response Path untuk text model");
            setTab(0);
            return;
        }
        if (hasImageModality && !imagePath) {
            toast.error("Isi Image Path untuk image model");
            setTab(0);
            return;
        }

        setSaving(true);
        try {
            const requestSchema = buildRequestSchema();
            
            await onSave({
                providerId,
                modelName,
                displayName,
                description: description || undefined,
                requestSchema,
                responsePath: modalities.includes("image") ? undefined : responsePath,
                imagePath: modalities.includes("image") ? imagePath : undefined,
                isActive,
                supportsStreaming,
                maxTokens: maxTokens ? Number(maxTokens) : undefined,
                temperature: temperature ? Number(temperature) : undefined,
                modalities: modalities.length > 0 ? modalities : undefined,
                imageConfig: modalities.includes("image") && (aspectRatio !== "16:9" || outputQuality !== "standard") 
                    ? { aspect_ratio: aspectRatio, output_quality: outputQuality }
                    : undefined,
            });
            toast.success("Model berhasil disimpan");
            handleOpenChange(false);
        } catch (error) {
            console.error(error);
            toast.error("Gagal menyimpan model");
        } finally {
            setSaving(false);
        }
    };

    // ── JSON preview ───────────────────────────────────────────────────────────
    const jsonPreview = JSON.stringify(
        {
            provider_id: providerId,
            model_name: modelName || "<model_name>",
            display_name: displayName || "<display_name>",
            description: description || null,
            request_schema: JSON.parse(buildRequestSchema()),
            response_path: modalities.includes("image") ? null : (responsePath || null),
            image_path: modalities.includes("image") ? (imagePath || null) : null,
            is_active: isActive,
            supports_streaming: supportsStreaming,
            max_tokens: maxTokens ? Number(maxTokens) : null,
            temperature: temperature ? Number(temperature) : null,
            modalities: modalities.length > 0 ? modalities : null,
            image_config: modalities.includes("image") && (aspectRatio !== "16:9" || outputQuality !== "standard")
                ? { aspect_ratio: aspectRatio, output_quality: outputQuality }
                : null,
        },
        null, 2,
    );

    // ── render ─────────────────────────────────────────────────────────────────
    return (
        <Dialog open={open} onOpenChange={handleOpenChange}>
            <DialogContent className="sm:max-w-5xl max-h-[90vh] p-0 flex flex-col">
                {/* Header - fixed */}
                <div className="p-5 pb-0">
                    <DialogHeader>
                        <DialogTitle className="text-base">Tambah Model AI</DialogTitle>
                        <DialogDescription className="text-xs">
                            Daftarkan model baru beserta konfigurasi request/response API-nya
                        </DialogDescription>
                    </DialogHeader>
                </div>

                {/* Progress bar - fixed */}
                <div className="px-5 pt-4">
                    <div className="flex gap-1">
                        {TABS.map((_, i) => (
                            <div
                                key={i}
                                className={cn(
                                    "h-1 flex-1 rounded-full transition-colors",
                                    i < tab  ? "bg-green-500 dark:bg-purple-500"
                                    : i === tab ? "bg-green-300 dark:bg-purple-400"
                                    : "bg-slate-200 dark:bg-white/10",
                                )}
                            />
                        ))}
                    </div>
                </div>

                {/* Tab bar - fixed */}
                <div className="px-5 pt-3">
                    <div className={cn(
                        "flex rounded-lg border p-0.5 gap-0.5",
                        "bg-slate-50 border-slate-200/80",
                        "dark:bg-white/[0.02] dark:border-white/[0.06]",
                    )}>
                        {TABS.map((label, i) => (
                            <button
                                key={label}
                                type="button"
                                onClick={() => setTab(i as Tab)}
                                className={cn(
                                    "flex-1 rounded-md py-1.5 text-xs font-medium transition-all",
                                    tab === i
                                        ? "bg-white shadow-sm text-slate-800 dark:bg-white/[0.06] dark:text-slate-100"
                                        : "text-slate-400 hover:text-slate-600 dark:text-slate-600 dark:hover:text-slate-400",
                                )}
                            >
                                {label}
                            </button>
                        ))}
                    </div>
                </div>

                {/* Scrollable content area */}
                <div className="flex-1 min-h-0 overflow-y-auto px-5 py-4">
                    {/* ── TAB 0: Provider & Model ──────────────────────────────────── */}
                    {tab === 0 && (
                        <div className="space-y-3">
                            {/* Provider & Model Name */}
                            <div className="grid grid-cols-2 gap-3">
                                <div className="space-y-1.5">
                                    <Label className="text-xs font-medium uppercase tracking-wide text-slate-500">
                                        Provider *
                                    </Label>
                                    <Select value={providerId} onValueChange={setProviderId}>
                                        <SelectTrigger className={inputCls}>
                                            <SelectValue placeholder="Pilih Provider" />
                                        </SelectTrigger>
                                        <SelectContent>
                                            {providers.map(p => (
                                                <SelectItem key={p.id} value={p.id}>
                                                    {p.displayName} ({p.name})
                                                </SelectItem>
                                            ))}
                                        </SelectContent>
                                    </Select>
                                </div>
                                <div className="space-y-1.5">
                                    <Label className="text-xs font-medium uppercase tracking-wide text-slate-500">
                                        Model Name *
                                    </Label>
                                    <Input 
                                        className={inputCls} 
                                        placeholder="openai/gpt-4o"
                                        value={modelName} 
                                        onChange={e => setModelName(e.target.value)} 
                                    />
                                    <p className="text-[10px] text-slate-400">
                                        Format: provider/model-name
                                    </p>
                                </div>
                            </div>

                            {/* Display Name & Description */}
                            <div className="space-y-1.5">
                                <Label className="text-xs font-medium uppercase tracking-wide text-slate-500">
                                    Display Name *
                                </Label>
                                <Input 
                                    className={inputCls} 
                                    placeholder="GPT-4o"
                                    value={displayName} 
                                    onChange={e => setDisplayName(e.target.value)} 
                                />
                            </div>

                            <div className="space-y-1.5">
                                <Label className="text-xs font-medium uppercase tracking-wide text-slate-500">
                                    Deskripsi (Opsional)
                                </Label>
                                <Textarea
                                    placeholder="Deskripsi singkat model ini..."
                                    value={description} 
                                    onChange={e => setDescription(e.target.value)}
                                    rows={2}
                                    className={cn(inputCls, "h-auto py-2 resize-none")}
                                />
                            </div>

                            {/* Response Paths */}
                            <div className="space-y-1.5">
                                <Label className="text-xs font-medium uppercase tracking-wide text-slate-500">
                                    Response Path (Content)
                                </Label>
                                <Input 
                                    className={inputCls} 
                                    placeholder="choices[0].message.content"
                                    value={responsePath} 
                                    onChange={e => setResponsePath(e.target.value)} 
                                    disabled={modalities.includes("image")}
                                />
                                <p className="text-[10px] text-slate-400">
                                    JSONPath untuk mengekstrak teks dari response (untuk text model)
                                </p>
                            </div>

                            <div className="space-y-1.5">
                                <Label className="text-xs font-medium uppercase tracking-wide text-slate-500">
                                    Image Path (Optional)
                                </Label>
                                <Input 
                                    className={inputCls} 
                                    placeholder="choices[0].message.images[0].image_url.url"
                                    value={imagePath} 
                                    onChange={e => setImagePath(e.target.value)} 
                                />
                                <p className="text-[10px] text-slate-400">
                                    JSONPath untuk mengekstrak gambar dari response (untuk image model)
                                </p>
                            </div>

                            {/* Model Parameters */}
                            <div className="grid grid-cols-2 gap-3">
                                <div className="space-y-1.5">
                                    <Label className="text-xs font-medium uppercase tracking-wide text-slate-500">
                                        Max Tokens
                                    </Label>
                                    <Input 
                                        className={inputCls} 
                                        type="number" 
                                        placeholder="4096"
                                        value={maxTokens} 
                                        onChange={e => setMaxTokens(e.target.value)} 
                                    />
                                </div>
                                <div className="space-y-1.5">
                                    <Label className="text-xs font-medium uppercase tracking-wide text-slate-500">
                                        Temperature
                                    </Label>
                                    <Input 
                                        className={inputCls} 
                                        type="number" 
                                        step="0.1"
                                        min="0"
                                        max="2"
                                        placeholder="0.7"
                                        value={temperature} 
                                        onChange={e => setTemperature(e.target.value)} 
                                    />
                                </div>
                            </div>

                            {/* Modalities (Image specific) */}
                            <div className="space-y-1.5">
                                <Label className="text-xs font-medium uppercase tracking-wide text-slate-500">
                                    Modalities
                                </Label>
                                <div className="flex gap-4">
                                    <label className="flex items-center gap-2 cursor-pointer">
                                        <input 
                                            type="checkbox" 
                                            checked={modalities.includes("text")}
                                            onChange={(e) => {
                                                if(e.target.checked) setModalities([...modalities, "text"]);
                                                else setModalities(modalities.filter(m => m !== "text"));
                                            }}
                                            className="rounded border-slate-300"
                                        />
                                        <span className="text-sm">Text</span>
                                    </label>
                                    <label className="flex items-center gap-2 cursor-pointer">
                                        <input 
                                            type="checkbox" 
                                            checked={modalities.includes("image")}
                                            onChange={(e) => {
                                                if(e.target.checked) setModalities([...modalities, "image"]);
                                                else setModalities(modalities.filter(m => m !== "image"));
                                            }}
                                            className="rounded border-slate-300"
                                        />
                                        <span className="text-sm">Image</span>
                                    </label>
                                </div>
                                <p className="text-[10px] text-slate-400">
                                    Pilih modalitas yang didukung model (text, image, atau keduanya)
                                </p>
                            </div>

                            {/* Image Configuration (shown if image modality is selected) */}
                            {modalities.includes("image") && (
                                <div className="grid grid-cols-2 gap-3">
                                    <div className="space-y-1.5">
                                        <Label className="text-xs font-medium uppercase tracking-wide text-slate-500">
                                            Aspect Ratio
                                        </Label>
                                        <Select value={aspectRatio} onValueChange={setAspectRatio}>
                                            <SelectTrigger className={inputCls}>
                                                <SelectValue />
                                            </SelectTrigger>
                                            <SelectContent>
                                                <SelectItem value="1:1">1:1 (Square)</SelectItem>
                                                <SelectItem value="16:9">16:9 (Landscape)</SelectItem>
                                                <SelectItem value="9:16">9:16 (Portrait)</SelectItem>
                                                <SelectItem value="4:3">4:3 (Standard)</SelectItem>
                                            </SelectContent>
                                        </Select>
                                    </div>
                                    <div className="space-y-1.5">
                                        <Label className="text-xs font-medium uppercase tracking-wide text-slate-500">
                                            Output Quality
                                        </Label>
                                        <Select value={outputQuality} onValueChange={setOutputQuality}>
                                            <SelectTrigger className={inputCls}>
                                                <SelectValue />
                                            </SelectTrigger>
                                            <SelectContent>
                                                <SelectItem value="standard">Standard</SelectItem>
                                                <SelectItem value="high">High (2K)</SelectItem>
                                                <SelectItem value="4k">4K</SelectItem>
                                            </SelectContent>
                                        </Select>
                                    </div>
                                </div>
                            )}

                            {/* Switches */}
                            <div className="grid grid-cols-2 gap-3 pb-2">
                                <div className={cn(
                                    "flex items-center justify-between rounded-xl border p-3",
                                    "bg-slate-50/60 border-slate-200/60",
                                    "dark:bg-white/[0.02] dark:border-white/[0.05]",
                                )}>
                                    <div>
                                        <p className="text-sm font-medium text-slate-800 dark:text-slate-200">Active</p>
                                        <p className="text-xs text-slate-400 mt-0.5">
                                            Model langsung bisa digunakan
                                        </p>
                                    </div>
                                    <Switch
                                        checked={isActive} onCheckedChange={setIsActive}
                                        className="data-[state=checked]:bg-green-600 dark:data-[state=checked]:bg-purple-600"
                                    />
                                </div>
                                <div className={cn(
                                    "flex items-center justify-between rounded-xl border p-3",
                                    "bg-slate-50/60 border-slate-200/60",
                                    "dark:bg-white/[0.02] dark:border-white/[0.05]",
                                )}>
                                    <div>
                                        <p className="text-sm font-medium text-slate-800 dark:text-slate-200">Streaming</p>
                                        <p className="text-xs text-slate-400 mt-0.5">
                                            Support response streaming
                                        </p>
                                    </div>
                                    <Switch
                                        checked={supportsStreaming} onCheckedChange={setSupportsStreaming}
                                        className="data-[state=checked]:bg-green-600 dark:data-[state=checked]:bg-purple-600"
                                    />
                                </div>
                            </div>
                        </div>
                    )}

                    {/* ── TAB 1: Schema Request ─────────────────────────────────── */}
                    {tab === 1 && (
                        <div className="space-y-3">
                            <div className="flex items-center justify-between">
                                <p className="text-xs font-medium text-slate-600 dark:text-slate-400">
                                    Schema Fields
                                    <span className="ml-1.5 text-slate-400">({fields.length} fields)</span>
                                </p>
                                <p className="text-[11px] text-slate-400">Drag untuk susun ulang</p>
                            </div>

                            {/* Drag-and-drop list */}
                            <div className="space-y-1.5 max-h-52 overflow-y-auto pr-0.5 border rounded-lg p-2">
                                {fields.map((f, i) => (
                                    <div
                                        key={f.key + i}
                                        draggable
                                        onDragStart={() => onDragStart(i)}
                                        onDragOver={e => e.preventDefault()}
                                        onDrop={() => onDrop(i)}
                                        className={cn(
                                            "flex items-center gap-2 px-3 py-2 rounded-lg border cursor-grab active:cursor-grabbing",
                                            "bg-white border-slate-200/80",
                                            "dark:bg-white/[0.03] dark:border-white/[0.08]",
                                            "hover:border-slate-300 dark:hover:border-white/[0.12]",
                                            "transition-colors select-none",
                                        )}
                                    >
                                        <GripVertical className="h-3.5 w-3.5 text-slate-300 dark:text-slate-700 flex-shrink-0" />
                                        <span className="text-xs font-mono font-medium text-slate-800 dark:text-slate-200 flex-1">
                                            {f.key}
                                        </span>
                                        <span className={cn(
                                            "text-[10px] rounded border px-1.5 py-0.5",
                                            "bg-slate-50 border-slate-200/60 text-slate-500",
                                            "dark:bg-white/[0.04] dark:border-white/[0.08] dark:text-slate-500",
                                        )}>
                                            {f.type}
                                        </span>
                                        <span className={cn(
                                            "text-[10px]",
                                            f.required
                                                ? "text-red-500 dark:text-red-400"
                                                : "text-slate-400 dark:text-slate-600",
                                        )}>
                                            {f.required ? "* required" : "optional"}
                                        </span>
                                        <button
                                            type="button"
                                            onClick={() => removeField(i)}
                                            className="text-slate-300 hover:text-red-400 dark:text-slate-700 dark:hover:text-red-400 transition-colors"
                                            aria-label={`Hapus field ${f.key}`}
                                        >
                                            <X className="h-3.5 w-3.5" />
                                        </button>
                                    </div>
                                ))}
                                {fields.length === 0 && (
                                    <div className="text-center py-8 text-slate-400 text-xs">
                                        Belum ada field. Tambahkan field baru di bawah.
                                    </div>
                                )}
                            </div>

                            {/* Add field form */}
                            <div className={cn(
                                "rounded-xl border p-4 space-y-3",
                                "bg-slate-50/60 border-slate-200/60",
                                "dark:bg-white/[0.02] dark:border-white/[0.05]",
                            )}>
                                <p className="text-xs font-medium text-slate-600 dark:text-slate-400">
                                    Tambah Field Baru
                                </p>
                                <div className="grid grid-cols-[1fr_90px_80px_auto] gap-2 items-end">
                                    <div className="space-y-1">
                                        <Label className="text-[10px] font-medium uppercase tracking-wide text-slate-400">
                                            Key
                                        </Label>
                                        <Input className={inputCls} placeholder="messages"
                                            value={newKey} onChange={e => setNewKey(e.target.value)}
                                            onKeyDown={e => e.key === "Enter" && addField()} />
                                    </div>
                                    <div className="space-y-1">
                                        <Label className="text-[10px] font-medium uppercase tracking-wide text-slate-400">
                                            Type
                                        </Label>
                                        <Select value={newType} onValueChange={v => setNewType(v as SchemaField["type"])}>
                                            <SelectTrigger className={inputCls}>
                                                <SelectValue />
                                            </SelectTrigger>
                                            <SelectContent>
                                                {FIELD_TYPES.map(t => (
                                                    <SelectItem key={t} value={t}>{t}</SelectItem>
                                                ))}
                                            </SelectContent>
                                        </Select>
                                    </div>
                                    <div className="space-y-1">
                                        <Label className="text-[10px] font-medium uppercase tracking-wide text-slate-400">
                                            Required
                                        </Label>
                                        <Select value={newRequired ? "yes" : "no"} onValueChange={v => setNewRequired(v === "yes")}>
                                            <SelectTrigger className={inputCls}>
                                                <SelectValue />
                                            </SelectTrigger>
                                            <SelectContent>
                                                <SelectItem value="yes">Yes</SelectItem>
                                                <SelectItem value="no">No</SelectItem>
                                            </SelectContent>
                                        </Select>
                                    </div>
                                    <Button size="sm" className="h-8 text-xs bg-green-600 hover:bg-green-700 dark:bg-purple-600 dark:hover:bg-purple-700"
                                        onClick={addField}>
                                        <Plus className="h-3.5 w-3.5" />
                                    </Button>
                                </div>
                            </div>

                            {/* Info about placeholders */}
                            <div className={cn(
                                "rounded-lg border p-3 flex items-start gap-2",
                                "bg-blue-50/60 border-blue-200/60",
                                "dark:bg-blue-500/10 dark:border-blue-500/20",
                            )}>
                                <Info className="h-4 w-4 text-blue-500 flex-shrink-0 mt-0.5" />
                                <div className="text-xs text-slate-600 dark:text-slate-400">
                                    <p className="font-medium mb-1">Placeholder Variables:</p>
                                    <p><code className="text-[10px] bg-white/50 px-1 rounded">{`{model}`}</code> - Akan diganti dengan model name</p>
                                    <p><code className="text-[10px] bg-white/50 px-1 rounded">{`{prompt}`}</code> - Akan diganti dengan input user</p>
                                    <p><code className="text-[10px] bg-white/50 px-1 rounded">{`{temperature}`}</code> - Akan diganti dengan nilai temperature</p>
                                    <p><code className="text-[10px] bg-white/50 px-1 rounded">{`{max_tokens}`}</code> - Akan diganti dengan max tokens</p>
                                </div>
                            </div>
                        </div>
                    )}

                    {/* ── TAB 2: Preview ────────────────────────────────────────── */}
                    {tab === 2 && (
                        <div className="space-y-3 pb-2">
                            {/* Info summary */}
                            <div className={cn(
                                "rounded-xl border p-4 space-y-2",
                                "bg-slate-50/60 border-slate-200/60",
                                "dark:bg-white/[0.02] dark:border-white/[0.05]",
                            )}>
                                <div className="flex items-center justify-between">
                                    <p className="text-xs font-medium uppercase tracking-wide text-slate-400">
                                        Model Info
                                    </p>
                                    <button
                                        type="button" onClick={() => setTab(0)}
                                        className="text-[11px] text-slate-400 hover:text-slate-600 flex items-center gap-0.5"
                                    >
                                        Edit <ChevronRight className="h-3 w-3" />
                                    </button>
                                </div>
                                <div className="space-y-1.5 max-h-64 overflow-y-auto">
                                    {[
                                        ["Provider",    providers.find(p => p.id === providerId)?.displayName || "—"],
                                        ["Model Name",  modelName  || "—"],
                                        ["Display Name",displayName || "—"],
                                        ["Description", description || "—"],
                                        ["Response Path", modalities.includes("image") ? "N/A (Image Model)" : (responsePath || "—")],
                                        ["Image Path", modalities.includes("image") ? (imagePath || "—") : "N/A (Text Model)"],
                                        ["Max Tokens",  maxTokens  || "—"],
                                        ["Temperature", temperature || "—"],
                                        ["Modalities",  modalities.join(", ") || "—"],
                                        ...(modalities.includes("image") ? [
                                            ["Aspect Ratio", aspectRatio],
                                            ["Output Quality", outputQuality],
                                        ] : []),
                                        ["Streaming",   supportsStreaming ? "Yes" : "No"],
                                        ["Status",      isActive ? "Active" : "Inactive"],
                                    ].map(([k, v]) => (
                                        <div key={k} className="flex justify-between text-xs">
                                            <span className="text-slate-400">{k}</span>
                                            <span className="font-medium text-slate-800 dark:text-slate-200 text-right max-w-[260px] truncate">{v}</span>
                                        </div>
                                    ))}
                                </div>
                            </div>

                            {/* JSON Schema Preview */}
                            <div className={cn(
                                "rounded-xl border overflow-hidden",
                                "border-slate-200/60 dark:border-white/[0.05]",
                            )}>
                                <div className="flex items-center justify-between px-4 py-2.5 border-b border-slate-200/60 dark:border-white/[0.05] bg-slate-50/80 dark:bg-white/[0.02]">
                                    <p className="text-xs font-medium uppercase tracking-wide text-slate-400">
                                        Request Schema
                                    </p>
                                    <button
                                        type="button" onClick={() => setTab(1)}
                                        className="text-[11px] text-slate-400 hover:text-slate-600 flex items-center gap-0.5"
                                    >
                                        Edit <ChevronRight className="h-3 w-3" />
                                    </button>
                                </div>
                                <pre className={cn(
                                    "text-[11px] leading-relaxed p-4 overflow-auto max-h-52",
                                    "font-mono text-slate-500 dark:text-slate-500",
                                    "bg-white dark:bg-white/[0.01]",
                                )}>
                                    {buildRequestSchema()}
                                </pre>
                            </div>

                            {/* Final Payload Preview */}
                            <div className={cn(
                                "rounded-xl border overflow-hidden",
                                "border-slate-200/60 dark:border-white/[0.05]",
                            )}>
                                <div className="px-4 py-2.5 border-b border-slate-200/60 dark:border-white/[0.05] bg-slate-50/80 dark:bg-white/[0.02]">
                                    <p className="text-xs font-medium uppercase tracking-wide text-slate-400">
                                        Final Payload
                                    </p>
                                </div>
                                <pre className={cn(
                                    "text-[11px] leading-relaxed p-4 overflow-auto max-h-52",
                                    "font-mono text-slate-500 dark:text-slate-500",
                                    "bg-white dark:bg-white/[0.01]",
                                )}>
                                    {jsonPreview}
                                </pre>
                            </div>
                        </div>
                    )}
                </div>

                {/* Footer - fixed */}
                <div className="p-5 pt-0">
                    <DialogFooter>
                        {tab > 0 && (
                            <Button variant="outline" size="sm" className="h-8 text-xs mr-auto"
                                onClick={() => setTab(prev => (prev - 1) as Tab)}>
                                ← Sebelumnya
                            </Button>
                        )}
                        <Button variant="outline" size="sm" className="h-8 text-xs"
                            onClick={() => handleOpenChange(false)}>
                            Cancel
                        </Button>
                        {tab < 2 ? (
                            <Button size="sm" className="h-8 text-xs bg-green-600 hover:bg-green-700 dark:bg-purple-600 dark:hover:bg-purple-700"
                                onClick={() => setTab(prev => (prev + 1) as Tab)}>
                                Lanjut →
                            </Button>
                        ) : (
                            <Button size="sm" className="h-8 text-xs bg-green-600 hover:bg-green-700 dark:bg-purple-600 dark:hover:bg-purple-700"
                                onClick={handleSave} disabled={saving}>
                                {saving && <Loader2 className="mr-1.5 h-3.5 w-3.5 animate-spin" />}
                                Simpan Model
                            </Button>
                        )}
                    </DialogFooter>
                </div>
            </DialogContent>
        </Dialog>
    );
}