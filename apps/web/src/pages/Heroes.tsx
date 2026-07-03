import { BarChart3, Calendar, Globe, Rocket, Sparkles, Zap } from "lucide-react";
import AnimatedBackground from "./AnimatedBackground";

export default function Heroes() {
    const features = [
        {
            icon: Zap,
            title: 'AI-Powered Content',
            description: 'Generate artikel dan gambar berkualitas tinggi dengan AI.'
        },
        {
            icon: Globe,
            title: 'Multi-Platform Publish',
            description: 'Terbitkan ke semua platform dalam satu klik.'
        },
        {
            icon: Calendar,
            title: 'Smart Scheduling',
            description: 'Atur waktu terbit paling strategis secara otomatis.'
        },
        {
            icon: BarChart3,
            title: 'SEO & Analytics',
            description: 'Optimasi SEO otomatis dan pantau performa konten.'
        }
    ];
    return (
        <>
            <div className="hidden lg:flex lg:w-[45%] xl:w-[42%] relative bg-[#090514] overflow-hidden flex-shrink-0 border-r border-purple-950/40">
                {/* Background Base Effects */}
                <div className="absolute inset-0 bg-[radial-gradient(circle_at_30%_20%,#3b1578_0%,transparent_50%)] opacity-40" />
                <div className="absolute inset-0 bg-[radial-gradient(circle_at_80%_80%,#1e0b4a_0%,transparent_50%)] opacity-60" />

                {/* Interactive Animated Canvas */}
                <AnimatedBackground />

                {/* Cyber Dot Grid Overlay */}
                <div className="absolute inset-0 bg-[linear-gradient(to_right,#ffffff03_1px,transparent_1px),linear-gradient(to_bottom,#ffffff03_1px,transparent_1px)] bg-[size:24px_24px]" />

                {/* Main Content Wrapper */}
                <div className="relative z-10 flex flex-col justify-between h-full p-10 xl:p-14 w-full overflow-y-auto select-none">

                    {/* Top Branding Header */}
                    <div className="flex items-center gap-3.5">
                        <div className="relative group">
                            <div className="absolute inset-0 bg-purple-500/40 rounded-xl blur-lg transition-transform group-hover:scale-110 duration-500" />
                            <div className="relative h-11 w-11 rounded-xl bg-gradient-to-b from-white/10 to-white/0 backdrop-blur-xl flex items-center justify-center border border-white/20 shadow-xl">
                                <Rocket className="h-5 w-5 text-purple-400 animate-pulse" />
                            </div>
                        </div>
                        <div>
                            <div className="flex items-center gap-2">
                                <h1 className="text-xl font-bold text-white tracking-wider uppercase bg-clip-text text-transparent bg-gradient-to-r from-white via-slate-200 to-purple-300">Kabar</h1>
                                <span className="px-1.5 py-0.5 text-[10px] tracking-normal font-semibold text-purple-300 bg-purple-500/10 border border-purple-500/20 rounded-md">v2.0</span>
                            </div>
                            <p className="text-xs text-purple-300/60 font-medium tracking-wide">Next-Gen Content Engine</p>
                        </div>
                    </div>

                    {/* Middle Features Display */}
                    <div className="my-auto py-3">
                        <div className="mb-2">
                            <h2 className="text-2xl xl:text-3xl font-extrabold text-white tracking-tight leading-tight">
                                Akselerasi Manajemen <br />
                                <span className="bg-clip-text text-transparent bg-gradient-to-r from-purple-400 via-violet-300 to-indigo-300">
                                    Konten Media Sosial Anda
                                </span>
                            </h2>
                        </div>

                        {/* Interactive Glassmorphism Grid */}
                        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
                            {features.map((feature, index) => (
                                <div
                                    key={index}
                                    className="group relative bg-white/[0.02] hover:bg-white/[0.05] backdrop-blur-md rounded-xl p-4.5 border border-white/[0.05] hover:border-purple-500/30 transition-all duration-300 shadow-xl"
                                >
                                    <div className="h-8 w-8 rounded-lg bg-purple-500/10 border border-purple-500/20 flex items-center justify-center mb-3 group-hover:bg-purple-500/20 transition-all">
                                        <feature.icon className="h-4 w-4 text-purple-300 group-hover:text-white transition-colors" />
                                    </div>
                                    <h3 className="text-sm font-semibold text-slate-200 mb-1">
                                        {feature.title}
                                    </h3>
                                    <p className="text-xs text-slate-400 leading-relaxed">
                                        {feature.description}
                                    </p>
                                </div>
                            ))}
                        </div>
                    </div>

                    {/* Bottom Footer Quote */}
                    <div className="pt-1 border-t border-white/[0.06]">
                        <p className="text-xs text-purple-200/50 italic leading-relaxed">
                            "Orchestrate your workflow, automated by intelligence, tailored perfectly for modern creators and visionary businesses."
                        </p>
                    </div>
                </div>
            </div>
        </>
    )
}