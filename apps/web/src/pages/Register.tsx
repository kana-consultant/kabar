import { useState } from 'react';
import { useNavigate } from '@tanstack/react-router';
import { Rocket, Mail, Lock, User, Eye, EyeOff, CheckCircle, AlertCircle, ArrowRight, Zap, Globe, Users, Calendar, BarChart3, Sparkles } from 'lucide-react';
import { Button, Input, Label } from '@kana-consultant/ui-kit';
import { useToast } from '@/hooks/use-toast';
import { register } from '@/services/user';
import Heroes from './Heroes';

export default function Register() {
    const toast = useToast();
    const navigate = useNavigate();
    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [showPassword, setShowPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const [errors, setErrors] = useState<{ [key: string]: string }>({});

    const validateForm = () => {
        const newErrors: { [key: string]: string } = {};

        if (!name.trim()) {
            newErrors.name = 'Nama wajib diisi';
        } else if (name.length < 3) {
            newErrors.name = 'Nama minimal 3 karakter';
        }

        if (!email.trim()) {
            newErrors.email = 'Email wajib diisi';
        } else if (!/\S+@\S+\.\S+/.test(email)) {
            newErrors.email = 'Format email tidak valid';
        }

        if (!password) {
            newErrors.password = 'Password wajib diisi';
        } else if (password.length < 6) {
            newErrors.password = 'Password minimal 6 karakter';
        } else if (password.length > 50) {
            newErrors.password = 'Password maksimal 50 karakter';
        }

        if (!confirmPassword) {
            newErrors.confirmPassword = 'Konfirmasi password wajib diisi';
        } else if (password !== confirmPassword) {
            newErrors.confirmPassword = 'Password tidak cocok';
        }

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!validateForm()) {
            return;
        }

        setIsLoading(true);
        try {
            await register({ email, name, password });
            toast.success('Registrasi berhasil!', {
                description: 'Silakan login dengan akun Anda',
            });
            navigate({ to: '/login' });
        } catch (err: any) {
            const message = err.response?.data?.message || err.message || 'Registrasi gagal';
            toast.error('Registrasi gagal', { description: message });
            if (message.toLowerCase().includes('already exists')) {
                setErrors({ email: 'Email sudah terdaftar' });
            }
        } finally {
            setIsLoading(false);
        }
    };

    const goToLogin = () => {
        navigate({ to: '/login' });
    };

    const getPasswordStrength = (pass: string) => {
        if (pass.length === 0) return 0;
        let strength = 0;
        if (pass.length >= 6) strength++;
        if (pass.length >= 10) strength++;
        if (/[A-Z]/.test(pass)) strength++;
        if (/[0-9]/.test(pass)) strength++;
        if (/[^A-Za-z0-9]/.test(pass)) strength++;
        return strength;
    };

    const passwordStrength = getPasswordStrength(password);
    const strengthColors = ['bg-red-500', 'bg-orange-500', 'bg-yellow-500', 'bg-lime-500', 'bg-green-500'];
    const strengthLabels = ['', 'Lemah', 'Sedang', 'Kuat', 'Sangat Kuat'];

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
        <div className="h-screen flex flex-col lg:flex-row overflow-hidden">
            {/* Hero Section - Left Side (Hidden on mobile) */}
            <Heroes />
            <div className="flex-1 flex items-center justify-center p-6 sm:p-8 lg:p-10 xl:p-12 overflow-y-auto bg-slate-50 dark:bg-[#080612]">
                <div className="w-full max-w-md">
                    {/* Mobile Logo */}
                    <div className="lg:hidden flex justify-center mb-8">
                        <div className="flex items-center gap-3">
                            <div className="relative">
                                <div className="absolute inset-0 bg-purple-500/20 rounded-xl blur-lg" />
                                <div className="relative h-11 w-11 rounded-xl bg-gradient-to-br from-purple-500 to-violet-600 flex items-center justify-center shadow-xl shadow-purple-500/20">
                                    <Rocket className="h-5 w-5 text-white" />
                                </div>
                            </div>
                            <div>
                                <h1 className="text-xl font-bold text-slate-900 dark:text-white">Kabar</h1>
                                <p className="text-sm text-slate-500 dark:text-slate-400">AI Content Management</p>
                            </div>
                        </div>
                    </div>

                    {/* Form Header */}
                    <div className="text-center mt-11 mb-8">
                        <h2 className="text-2xl sm:text-3xl font-bold tracking-tight mb-2">
                            <span className="bg-gradient-to-br from-slate-900 to-slate-700 dark:from-white dark:to-slate-300 bg-clip-text text-transparent">
                                Buat Akun Baru
                            </span>
                        </h2>
                        <p className="text-sm text-slate-500 dark:text-slate-400 max-w-xs mx-auto">
                            Mulai kelola konten Anda dengan AI
                        </p>
                    </div>

                    <form onSubmit={handleSubmit} noValidate className="space-y-5">
                        {/* Name Field */}
                        <div className="space-y-2">
                            <Label htmlFor="name" className="text-sm font-semibold text-slate-700 dark:text-slate-300">
                                Nama Lengkap
                            </Label>
                            <div className="relative group">
                                <div className="absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400 group-focus-within:text-purple-500 transition-colors duration-200">
                                    <User className="h-4 w-4" />
                                </div>
                                <Input
                                    id="name"
                                    type="text"
                                    placeholder="Masukkan nama lengkap"
                                    value={name}
                                    onChange={(e: any) => {
                                        setName(e.target.value);
                                        if (errors.name) setErrors({ ...errors, name: '' });
                                    }}
                                    className={`pl-11 h-11 text-sm bg-white dark:bg-slate-800/50 border-2 rounded-xl transition-all duration-200
                            ${errors.name
                                            ? 'border-red-300 dark:border-red-800 focus:border-red-500 focus:ring-red-500/20'
                                            : 'border-slate-200 dark:border-slate-700 focus:border-purple-500 focus:ring-4 focus:ring-purple-500/10'
                                        }
                            ${name.length >= 3 && !errors.name ? 'border-green-300 dark:border-green-800' : ''}
                        `}
                                    disabled={isLoading}
                                    autoComplete="name"
                                />
                                {name.length >= 3 && !errors.name && (
                                    <div className="absolute right-3.5 top-1/2 -translate-y-1/2">
                                        <CheckCircle className="h-4 w-4 text-green-500" />
                                    </div>
                                )}
                            </div>
                            {errors.name && (
                                <div className="flex items-center gap-2 text-xs text-red-500 bg-red-50 dark:bg-red-950/30 px-3 py-1.5 rounded-lg">
                                    <AlertCircle className="h-3.5 w-3.5 flex-shrink-0" />
                                    <span>{errors.name}</span>
                                </div>
                            )}
                        </div>

                        {/* Email Field */}
                        <div className="space-y-2">
                            <Label htmlFor="email" className="text-sm font-semibold text-slate-700 dark:text-slate-300">
                                Alamat Email
                            </Label>
                            <div className="relative group">
                                <div className="absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400 group-focus-within:text-purple-500 transition-colors duration-200">
                                    <Mail className="h-4 w-4" />
                                </div>
                                <Input
                                    id="email"
                                    type="email"
                                    placeholder="nama@email.com"
                                    value={email}
                                    onChange={(e: any) => {
                                        setEmail(e.target.value);
                                        if (errors.email) setErrors({ ...errors, email: '' });
                                    }}
                                    className={`pl-11 h-11 text-sm bg-white dark:bg-slate-800/50 border-2 rounded-xl transition-all duration-200
                            ${errors.email
                                            ? 'border-red-300 dark:border-red-800 focus:border-red-500 focus:ring-red-500/20'
                                            : 'border-slate-200 dark:border-slate-700 focus:border-purple-500 focus:ring-4 focus:ring-purple-500/10'
                                        }
                            ${email && /\S+@\S+\.\S+/.test(email) && !errors.email ? 'border-green-300 dark:border-green-800' : ''}
                        `}
                                    disabled={isLoading}
                                    autoComplete="email"
                                />
                                {email && /\S+@\S+\.\S+/.test(email) && !errors.email && (
                                    <div className="absolute right-3.5 top-1/2 -translate-y-1/2">
                                        <CheckCircle className="h-4 w-4 text-green-500" />
                                    </div>
                                )}
                            </div>
                            {errors.email && (
                                <div className="flex items-center gap-2 text-xs text-red-500 bg-red-50 dark:bg-red-950/30 px-3 py-1.5 rounded-lg">
                                    <AlertCircle className="h-3.5 w-3.5 flex-shrink-0" />
                                    <span>{errors.email}</span>
                                </div>
                            )}
                        </div>

                        {/* Password Field */}
                        <div className="space-y-2">
                            <Label htmlFor="password" className="text-sm font-semibold text-slate-700 dark:text-slate-300">
                                Password
                            </Label>
                            <div className="relative group">
                                <div className="absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400 group-focus-within:text-purple-500 transition-colors duration-200">
                                    <Lock className="h-4 w-4" />
                                </div>
                                <Input
                                    id="password"
                                    type={showPassword ? 'text' : 'password'}
                                    placeholder="Minimal 6 karakter"
                                    value={password}
                                    onChange={(e: any) => {
                                        setPassword(e.target.value);
                                        if (errors.password) setErrors({ ...errors, password: '' });
                                    }}
                                    className={`pl-11 pr-12 h-11 text-sm bg-white dark:bg-slate-800/50 border-2 rounded-xl transition-all duration-200
                            ${errors.password
                                            ? 'border-red-300 dark:border-red-800 focus:border-red-500 focus:ring-red-500/20'
                                            : 'border-slate-200 dark:border-slate-700 focus:border-purple-500 focus:ring-4 focus:ring-purple-500/10'
                                        }
                        `}
                                    disabled={isLoading}
                                    autoComplete="new-password"
                                />
                                <button
                                    type="button"
                                    onClick={() => setShowPassword(!showPassword)}
                                    className="absolute right-3.5 top-1/2 -translate-y-1/2 p-1.5 text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-700 transition-all"
                                    tabIndex={-1}
                                >
                                    {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                                </button>
                            </div>
                            {password && (
                                <div className="space-y-2 bg-slate-50 dark:bg-slate-800/30 p-3 rounded-lg">
                                    <div className="flex gap-1.5">
                                        {[1, 2, 3, 4, 5].map((level) => (
                                            <div
                                                key={level}
                                                className={`h-1.5 flex-1 rounded-full transition-all duration-300 ${level <= passwordStrength
                                                        ? strengthColors[passwordStrength - 1] + ' shadow-lg'
                                                        : 'bg-slate-200 dark:bg-slate-700'
                                                    }`}
                                            />
                                        ))}
                                    </div>
                                    <p className={`text-xs font-medium ${passwordStrength <= 1 ? 'text-red-600' :
                                            passwordStrength === 2 ? 'text-orange-600' :
                                                passwordStrength === 3 ? 'text-yellow-600' :
                                                    'text-green-600'
                                        }`}>
                                        {strengthLabels[passwordStrength]}
                                        {passwordStrength >= 4 && <CheckCircle className="inline h-3 w-3 ml-1" />}
                                    </p>
                                </div>
                            )}
                            {errors.password && (
                                <div className="flex items-center gap-2 text-xs text-red-500 bg-red-50 dark:bg-red-950/30 px-3 py-1.5 rounded-lg">
                                    <AlertCircle className="h-3.5 w-3.5 flex-shrink-0" />
                                    <span>{errors.password}</span>
                                </div>
                            )}
                        </div>

                        {/* Confirm Password Field */}
                        <div className="space-y-2">
                            <Label htmlFor="confirmPassword" className="text-sm font-semibold text-slate-700 dark:text-slate-300">
                                Konfirmasi Password
                            </Label>
                            <div className="relative group">
                                <div className="absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400 group-focus-within:text-purple-500 transition-colors duration-200">
                                    <Lock className="h-4 w-4" />
                                </div>
                                <Input
                                    id="confirmPassword"
                                    type={showConfirmPassword ? 'text' : 'password'}
                                    placeholder="Ulangi password"
                                    value={confirmPassword}
                                    onChange={(e: any) => {
                                        setConfirmPassword(e.target.value);
                                        if (errors.confirmPassword) setErrors({ ...errors, confirmPassword: '' });
                                    }}
                                    className={`pl-11 pr-12 h-11 text-sm bg-white dark:bg-slate-800/50 border-2 rounded-xl transition-all duration-200
                            ${errors.confirmPassword
                                            ? 'border-red-300 dark:border-red-800 focus:border-red-500 focus:ring-red-500/20'
                                            : 'border-slate-200 dark:border-slate-700 focus:border-purple-500 focus:ring-4 focus:ring-purple-500/10'
                                        }
                            ${password && confirmPassword && password === confirmPassword ? 'border-green-300 dark:border-green-800' : ''}
                        `}
                                    disabled={isLoading}
                                    autoComplete="new-password"
                                />
                                <button
                                    type="button"
                                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                                    className="absolute right-3.5 top-1/2 -translate-y-1/2 p-1.5 text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-700 transition-all"
                                    tabIndex={-1}
                                >
                                    {showConfirmPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                                </button>
                            </div>
                            {password && confirmPassword && password === confirmPassword && (
                                <div className="flex items-center gap-2 text-xs text-green-600 bg-green-50 dark:bg-green-950/30 px-3 py-1.5 rounded-lg">
                                    <CheckCircle className="h-3.5 w-3.5" />
                                    <span>Password cocok</span>
                                </div>
                            )}
                            {errors.confirmPassword && (
                                <div className="flex items-center gap-2 text-xs text-red-500 bg-red-50 dark:bg-red-950/30 px-3 py-1.5 rounded-lg">
                                    <AlertCircle className="h-3.5 w-3.5 flex-shrink-0" />
                                    <span>{errors.confirmPassword}</span>
                                </div>
                            )}
                        </div>

                        {/* Submit Button & Login Link */}
                        <div className="pt-2 space-y-4">
                            <Button
                                type="submit"
                                className="w-full h-12 bg-gradient-to-r from-purple-600 to-violet-600 hover:from-purple-700 hover:to-violet-700 text-white font-semibold rounded-xl shadow-xl shadow-purple-500/25 hover:shadow-2xl hover:shadow-purple-500/30 transition-all duration-300 hover:scale-[1.02] active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100"
                                disabled={isLoading}
                            >
                                {isLoading ? (
                                    <div className="flex items-center justify-center gap-3">
                                        <div className="h-5 w-5 animate-spin rounded-full border-2 border-white/30 border-t-white" />
                                        <span className="font-medium">Memproses...</span>
                                    </div>
                                ) : (
                                    <div className="flex items-center justify-center gap-2">
                                        <span className="font-medium">Daftar Sekarang</span>
                                        <ArrowRight className="h-5 w-5" />
                                    </div>
                                )}
                            </Button>

                            <p className="text-sm text-slate-600 dark:text-slate-400 text-center">
                                Sudah punya akun?{' '}
                                <button
                                    type="button"
                                    onClick={goToLogin}
                                    className="text-purple-600 hover:text-purple-700 dark:text-purple-400 dark:hover:text-purple-300 font-semibold hover:underline decoration-2 underline-offset-2 transition-all"
                                    disabled={isLoading}
                                >
                                    Masuk sekarang
                                </button>
                            </p>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
}