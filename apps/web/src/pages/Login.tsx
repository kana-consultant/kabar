import { useState, useEffect, useRef } from 'react';
import { useNavigate } from '@tanstack/react-router';
import { Rocket, Mail, Lock, Eye, EyeOff, AlertCircle, ArrowRight, Zap, Globe, Calendar, BarChart3, Sparkles } from 'lucide-react';
import { Button, Input, Label } from '@kana-consultant/ui-kit';
import { useAuth } from '@/hooks/auth/useAuth';
import AnimatedBackground from './AnimatedBackground';
import Heroes from './Heroes';

// Component untuk background animasi futuristik 3D Interaktif
// Component untuk background animasi gerakan abstrak interaktif (Cosmic Flow)


export default function Login() {
    const navigate = useNavigate();
    const { login, isLoading: authLoading } = useAuth();

    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [showPassword, setShowPassword] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState('');

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        e.stopPropagation();
        setError('');

        if (!email || !password) {
            setError('Email dan password wajib diisi');
            return;
        }
        setIsLoading(true);
        try {
            const response = await login(email, password);
            if (response) {
                navigate({ to: '/dashboard' });
            }
        } catch (err: any) {
            setError(err?.message || 'Terjadi kesalahan saat login');
        }finally{
            setIsLoading(false);
        }
    };

    const goToRegister = (e: React.MouseEvent) => {
        e.preventDefault();
        navigate({ to: '/register' });
    };

    const isLoadingState = isLoading || authLoading;

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
        <div className="h-screen flex flex-col lg:flex-row overflow-hidden font-sans">
            {/* HERO SECTION - LEFT SIDE */}
            <Heroes />
            {/* LOGIN FORM - RIGHT SIDE */}
            <div className="flex-1 flex items-center justify-center p-6 sm:p-8 lg:p-10 xl:p-12 overflow-y-auto bg-slate-50 dark:bg-[#080612]">
                <div className="w-full max-w-md">
                    {/* Mobile Logo Layout */}
                    <div className="lg:hidden flex justify-center mb-8">
                        <div className="flex items-center gap-3">
                            <div className="relative">
                                <div className="absolute inset-0 bg-purple-500/20 rounded-xl blur-lg" />
                                <div className="relative h-11 w-11 rounded-xl bg-gradient-to-br from-purple-500 to-violet-600 flex items-center justify-center shadow-xl">
                                    <Rocket className="h-5 w-5 text-white" />
                                </div>
                            </div>
                            <div>
                                <h1 className="text-xl font-bold text-slate-900 dark:text-white tracking-wide">Kabar</h1>
                                <p className="text-xs text-slate-500 dark:text-slate-400">AI Content Workspace</p>
                            </div>
                        </div>
                    </div>

                    {/* Header */}
                    <div className="text-center mb-8">
                        <h2 className="text-2xl sm:text-3xl font-bold tracking-tight mb-1.5">
                            <span className="bg-gradient-to-br from-slate-900 to-slate-700 dark:from-white dark:to-slate-300 bg-clip-text text-transparent">
                                Welcome Back
                            </span>
                        </h2>
                        <p className="text-sm text-slate-500 dark:text-slate-400 max-w-xs mx-auto">
                            Masuk ke akun Anda untuk melanjutkan manajemen automasi
                        </p>
                    </div>

                    <form onSubmit={handleSubmit} noValidate className="space-y-4.5">
                        {error && (
                            <div className="flex items-start gap-3 rounded-lg bg-red-50 dark:bg-red-950/20 border border-red-200 dark:border-red-900/50 p-4">
                                <AlertCircle className="h-5 w-5 mt-0.5 flex-shrink-0 text-red-600 dark:text-red-400" />
                                <div className="flex-1">
                                    <p className="text-sm font-semibold text-red-800 dark:text-red-300">Gagal masuk</p>
                                    <p className="text-sm text-red-600 dark:text-red-400 mt-0.5">{error}</p>
                                </div>
                            </div>
                        )}

                        {/* Email Field */}
                        <div className="space-y-1.5">
                            <Label htmlFor="email" className="text-xs font-bold uppercase tracking-wider text-slate-600 dark:text-slate-400">
                                Alamat Email
                            </Label>
                            <div className="relative group">
                                <div className="absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400 group-focus-within:text-purple-500 transition-colors duration-200">
                                    <Mail className="h-4 w-4" />
                                </div>
                                <Input
                                    id="email"
                                    type="email"
                                    placeholder="nama@perusahaan.com"
                                    value={email}
                                    onChange={(e: any) => setEmail(e.target.value)}
                                    className="pl-11 h-11 text-sm bg-white dark:bg-slate-800/30 border border-slate-200 dark:border-slate-800 rounded-xl transition-all duration-200 focus:border-purple-500 focus:ring-4 focus:ring-purple-500/10"
                                    disabled={isLoadingState}
                                    autoComplete="email"
                                    required
                                />
                            </div>
                        </div>

                        {/* Password Field */}
                        <div className="space-y-1.5">
                            <Label htmlFor="password" className="text-xs font-bold uppercase tracking-wider text-slate-600 dark:text-slate-400">
                                Password
                            </Label>
                            <div className="relative group">
                                <div className="absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400 group-focus-within:text-purple-500 transition-colors duration-200">
                                    <Lock className="h-4 w-4" />
                                </div>
                                <Input
                                    id="password"
                                    type={showPassword ? 'text' : 'password'}
                                    placeholder="••••••••"
                                    value={password}
                                    onChange={(e: any) => setPassword(e.target.value)}
                                    className="pl-11 pr-12 h-11 text-sm bg-white dark:bg-slate-800/30 border border-slate-200 dark:border-slate-800 rounded-xl transition-all duration-200 focus:border-purple-500 focus:ring-4 focus:ring-purple-500/10"
                                    disabled={isLoadingState}
                                    autoComplete="current-password"
                                    required
                                />
                                <button
                                    type="button"
                                    onClick={() => setShowPassword(!showPassword)}
                                    className="absolute right-3.5 top-1/2 -translate-y-1/2 p-1.5 text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 transition-all"
                                    disabled={isLoadingState}
                                    tabIndex={-1}
                                >
                                    {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                                </button>
                            </div>
                        </div>

                        {/* Remember Me */}
                        <div className="flex items-center pt-0.5">
                            <label className="flex items-center gap-2.5 cursor-pointer group">
                                <input
                                    type="checkbox"
                                    className="h-4 w-4 rounded border-slate-300 dark:border-slate-800 text-purple-600 focus:ring-purple-500/20 transition-colors cursor-pointer"
                                    disabled={isLoadingState}
                                />
                                <span className="text-sm text-slate-600 dark:text-slate-400 group-hover:text-slate-900 dark:group-hover:text-slate-200 transition-colors">
                                    Ingat saya di perangkat ini
                                </span>
                            </label>
                        </div>

                        {/* Submit Actions */}
                        <div className="pt-2 space-y-4">
                            <Button
                                type="submit"
                                className="w-full h-11 bg-gradient-to-r from-purple-600 to-violet-600 hover:from-purple-700 hover:to-violet-700 text-white font-semibold rounded-xl shadow-xl shadow-purple-500/10 hover:shadow-purple-500/20 transition-all duration-300 hover:scale-[1.01] disabled:opacity-50"
                                disabled={isLoadingState}
                            >
                                {isLoadingState ? (
                                    <div className="flex items-center justify-center gap-2">
                                        <div className="h-4 w-4 animate-spin rounded-full border-2 border-white/30 border-t-white" />
                                        <span>Memproses...</span>
                                    </div>
                                ) : (
                                    <div className="flex items-center justify-center gap-1.5">
                                        <span>Masuk ke Workspace</span>
                                        <ArrowRight className="h-4 w-4" />
                                    </div>
                                )}
                            </Button>

                            <p className="text-sm text-slate-600 dark:text-slate-400 text-center">
                                Belum punya akses?{' '}
                                <button
                                    type="button"
                                    onClick={goToRegister}
                                    className="text-purple-600 hover:text-purple-700 dark:text-purple-400 dark:hover:text-purple-300 font-semibold hover:underline underline-offset-4 transition-all"
                                    disabled={isLoadingState}
                                >
                                    Daftar sekarang
                                </button>
                            </p>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
}