import { useState } from 'react';
import { useNavigate } from '@tanstack/react-router';
import { Rocket, Mail, Lock, Eye, EyeOff, AlertCircle } from 'lucide-react';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { toast } from 'sonner';
import { useAuth } from '@/hooks/auth/useAuth';

export default function Login() {
    const navigate = useNavigate();
    const { login, isLoading: authLoading } = useAuth();

    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [showPassword, setShowPassword] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState('');

    // Handler submit yang benar
    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();        // ← Ini yang paling penting
        e.stopPropagation();       // Tambahan pencegahan

        setError('');

        if (!email || !password) {
            setError('Email dan password wajib diisi');
            return;
        }

        setIsLoading(true);

        try {
            const success = await login(email, password);

            if (success) {
                toast.success('Berhasil masuk!', {
                    description: 'Selamat datang kembali di KABAR',
                    duration: 3000,
                });
                navigate({ to: '/' });
            }
        } catch (err: any) {
            const status = err.response?.status;
            let errorMessage = 'Login gagal';

            if (status === 401) {
                errorMessage = 'Email atau password salah';
                toast.error('Login gagal', {
                    description: 'Email atau password yang Anda masukkan tidak sesuai.',
                    duration: 4500,
                });
            } else if (status === 403) {
                errorMessage = 'Akses ditolak. Akun Anda mungkin belum aktif atau diblokir.';
                toast.error('Akses Ditolak', { description: errorMessage });
            } else {
                errorMessage = err.response?.data?.message 
                    || err.message 
                    || 'Terjadi kesalahan saat login. Silakan coba lagi.';
                toast.error('Login gagal', { description: errorMessage });
            }

            setError(errorMessage);
            // JANGAN navigate di sini!
        } finally {
            setIsLoading(false);
        }
    };

    const goToRegister = (e: React.MouseEvent) => {
        e.preventDefault();
        navigate({ to: '/register' });
    };

    // const goToForgotPassword = (e: React.MouseEvent) => {
    //     e.preventDefault();
    //     navigate({ to: "/forgot-password" });
    // };

    const isLoadingState = isLoading || authLoading;

    return (
        <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-950 dark:to-slate-900 p-4">
            <Card className="w-full max-w-md">
                <CardHeader className="space-y-1 text-center">
                    <div className="flex justify-center mb-4">
                        <div className="h-12 w-12 rounded-xl bg-gradient-to-r from-cyan-500 to-cyan-600 flex items-center justify-center shadow-lg shadow-cyan-500/25">
                            <Rocket className="h-6 w-6 text-white" />
                        </div>
                    </div>
                    <CardTitle className="text-2xl font-bold tracking-tight dark:bg-gradient-to-r dark:from-white dark:to-zinc-400 dark:bg-clip-text dark:text-transparent">
                        Welcome Back
                    </CardTitle>
                    <CardDescription>
                        Masuk ke akun Anda untuk melanjutkan
                    </CardDescription>
                </CardHeader>

                {/* FORM DENGAN onSubmit */}
                <form onSubmit={handleSubmit} noValidate>
                    <CardContent className="space-y-4">
                        {error && (
                            <div className="flex items-start gap-3 rounded-md bg-red-500/10 border border-red-500/20 p-3 text-sm text-red-600 dark:text-red-400">
                                <AlertCircle className="h-4 w-4 mt-0.5 flex-shrink-0" />
                                <span>{error}</span>
                            </div>
                        )}

                        <div className="space-y-2">
                            <Label htmlFor="email" className="text-slate-700 dark:text-zinc-400">Email</Label>
                            <div className="relative">
                                <Mail className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-500 dark:text-zinc-500" />
                                <Input
                                    id="email"
                                    type="email"
                                    placeholder="admin@seo.com"
                                    value={email}
                                    onChange={(e) => setEmail(e.target.value)}
                                    className="pl-10 bg-white dark:bg-zinc-950/50 border-slate-200 dark:border-zinc-800 focus:border-cyan-500"
                                    disabled={isLoadingState}
                                    autoComplete="email"
                                />
                            </div>
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="password" className="text-slate-700 dark:text-zinc-400">Password</Label>
                            <div className="relative">
                                <Lock className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-500 dark:text-zinc-500" />
                                <Input
                                    id="password"
                                    type={showPassword ? 'text' : 'password'}
                                    placeholder="••••••••"
                                    value={password}
                                    onChange={(e) => setPassword(e.target.value)}
                                    className="pl-10 pr-10 bg-white dark:bg-zinc-950/50 border-slate-200 dark:border-zinc-800 focus:border-cyan-500"
                                    disabled={isLoadingState}
                                    autoComplete="current-password"
                                />
                                <button
                                    type="button"
                                    onClick={() => setShowPassword(!showPassword)}
                                    className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-500 hover:text-slate-700 dark:text-zinc-500 dark:hover:text-zinc-300 transition-colors"
                                    disabled={isLoadingState}
                                >
                                    {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                                </button>
                            </div>
                        </div>

                        <div className="flex items-center justify-between text-sm">
                            <label className="flex items-center gap-2 cursor-pointer">
                                <input
                                    type="checkbox"
                                    className="rounded border-slate-300 dark:border-zinc-700 accent-cyan-500"
                                    disabled={isLoadingState}
                                />
                                <span className="text-slate-600 dark:text-zinc-400">Ingat saya</span>
                            </label>
                            {/* <button
                                type="button"
                                onClick={goToForgotPassword}
                                className="text-cyan-600 hover:text-cyan-700 dark:text-cyan-400 dark:hover:text-cyan-300 hover:underline transition-colors"
                                disabled={isLoadingState}
                            >
                                Lupa password?
                            </button> */}
                        </div>
                    </CardContent>

                    <CardFooter className="flex flex-col gap-4">
                        <Button
                            type="submit"                     
                            className="w-full bg-gradient-to-r from-cyan-500 to-cyan-600 hover:from-cyan-600 hover:to-cyan-700 text-white font-medium shadow-lg shadow-cyan-500/25"
                            disabled={isLoadingState}
                        >
                            {isLoadingState ? (
                                <>
                                    <div className="mr-2 h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent" />
                                    Memproses...
                                </>
                            ) : (
                                'Masuk'
                            )}
                        </Button>

                        <p className="text-center text-sm text-slate-600 dark:text-zinc-400">
                            Belum punya akun?{' '}
                            <button
                                type="button"
                                onClick={goToRegister}
                                className="text-cyan-600 hover:text-cyan-700 dark:text-cyan-400 dark:hover:text-cyan-300 hover:underline font-medium transition-colors"
                                disabled={isLoadingState}
                            >
                                Daftar sekarang
                            </button>
                        </p>
                    </CardFooter>
                </form>
            </Card>
        </div>
    );
}