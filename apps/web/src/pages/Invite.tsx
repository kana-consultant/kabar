// invite.tsx
import { acceptInvite, verifyInviteToken } from "@/services/team";
import { useEffect, useState } from "react";
import { useNavigate } from "@tanstack/react-router";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Loader2, CheckCircle, XCircle, Mail, Lock, Eye, EyeOff, Users, Shield, AlertCircle, User } from "lucide-react";
import type { TeamInvite } from "@/services/team/types";

interface InviteData extends TeamInvite {
  TeamName: string;
  role: string;
  email: string;
}

export default function InvitePage() {
    const navigate = useNavigate();
    const searchParams = new URLSearchParams(location.search);
    const token = searchParams.get("token") as string;

    const [isLoading, setIsLoading] = useState<boolean>(true);
    const [inviteData, setInviteData] = useState<InviteData | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [isAccepted, setIsAccepted] = useState<boolean>(false);

    // Form state untuk registrasi
    const [showPassword, setShowPassword] = useState<boolean>(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState<boolean>(false);
    const [password, setPassword] = useState<string>("");
    const [confirmPassword, setConfirmPassword] = useState<string>("");
    const [fullName, setFullName] = useState<string>("");
    const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
    const [passwordError, setPasswordError] = useState<string | null>(null);

    // Step 1: Verify invite token
    const verifyToken = async (token: string) => {
        try {
            console.log("Verifying token:", token);
            const result = await verifyInviteToken(token);
            setInviteData(result as InviteData);
            setError(null);
        } catch (err: any) {
            setError(err?.message || "Failed to verify invitation");
        } finally {
            setIsLoading(false);
        }
    };

    // Step 2: Accept invite and redirect
    const acceptInviteAndRedirect = async (token: string, password: string) => {
        setIsSubmitting(true);
        try {
            const response = await acceptInvite(fullName, token, password);
            if (response) {
                setIsAccepted(true);
                setTimeout(() => {
                    navigate({ to: "/login" });
                }, 2000);
                return true;
            }
        } catch (err: any) {
            setError(err?.message || "Failed to accept invitation");
        } finally {
            setIsSubmitting(false);
        }
    };

    // Handle form submit untuk registrasi
    const handleRegisterAndAccept = async (e: React.FormEvent) => {
        e.preventDefault();
        e.stopPropagation();
        setPasswordError(null);

        // Validasi password
        if (password.length < 6) {
            setPasswordError("Password must be at least 6 characters");
            return;
        }

        if (password !== confirmPassword) {
            setPasswordError("Passwords do not match");
            return;
        }

        if (!fullName.trim()) {
            setPasswordError("Full name is required");
            return;
        }

        await acceptInviteAndRedirect(token, password);
    };

    useEffect(() => {
        if (!token) {
            setError("No invitation token provided");
            setIsLoading(false);
            return;
        }

        verifyToken(token);
    }, [token]);

    // Loading state
    if (isLoading) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-950 dark:to-slate-900 p-4">
                <Card className="w-full max-w-md">
                    <CardContent className="flex flex-col items-center justify-center py-12">
                        <Loader2 className="h-8 w-8 animate-spin text-cyan-500" />
                        <p className="mt-4 text-slate-600 dark:text-zinc-400">Verifying invitation...</p>
                    </CardContent>
                </Card>
            </div>
        );
    }

    // Error state
    if (error) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-950 dark:to-slate-900 p-4">
                <Card className="w-full max-w-md">
                    <CardHeader>
                        <div className="flex justify-center mb-4">
                            <div className="rounded-full bg-red-500/10 p-3">
                                <XCircle className="h-12 w-12 text-red-500" />
                            </div>
                        </div>
                        <CardTitle className="text-center text-red-600 dark:text-red-400">
                            Invalid Invitation
                        </CardTitle>
                        <CardDescription className="text-center text-slate-600 dark:text-zinc-400">
                            {error}
                        </CardDescription>
                    </CardHeader>
                    <CardFooter className="flex justify-center">
                        <Button 
                            onClick={() => navigate({ to: "/" })}
                            className="bg-gradient-to-r from-cyan-500 to-cyan-600 hover:from-cyan-600 hover:to-cyan-700"
                        >
                            Go to Home
                        </Button>
                    </CardFooter>
                </Card>
            </div>
        );
    }

    // Success state
    if (isAccepted) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-950 dark:to-slate-900 p-4">
                <Card className="w-full max-w-md">
                    <CardHeader>
                        <div className="flex justify-center mb-4">
                            <div className="rounded-full bg-green-500/10 p-3">
                                <CheckCircle className="h-12 w-12 text-green-500" />
                            </div>
                        </div>
                        <CardTitle className="text-center text-green-600 dark:text-green-400">
                            Successfully Joined!
                        </CardTitle>
                        <CardDescription className="text-center text-slate-600 dark:text-zinc-400">
                            You have successfully joined {inviteData?.TeamName || "the team"}.
                            Redirecting to login page...
                        </CardDescription>
                    </CardHeader>
                    <CardContent className="flex justify-center">
                        <Loader2 className="h-6 w-6 animate-spin text-cyan-500" />
                    </CardContent>
                </Card>
            </div>
        );
    }

    // Registration form for new users
    if (inviteData) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-950 dark:to-slate-900 p-4">
                <Card className="w-full max-w-md">
                    <CardHeader className="space-y-1 text-center">
                        <div className="flex justify-center mb-4">
                            <div className="h-12 w-12 rounded-xl bg-gradient-to-r from-cyan-500 to-cyan-600 flex items-center justify-center shadow-lg shadow-cyan-500/25">
                                <Users className="h-6 w-6 text-white" />
                            </div>
                        </div>
                        <CardTitle className="text-2xl font-bold tracking-tight dark:bg-gradient-to-r dark:from-white dark:to-zinc-400 dark:bg-clip-text dark:text-transparent">
                            Join {inviteData.TeamName}
                        </CardTitle>
                        <CardDescription className="text-slate-600 dark:text-zinc-400">
                            Complete your registration to join the team
                        </CardDescription>
                    </CardHeader>

                    <form onSubmit={handleRegisterAndAccept} noValidate>
                        <CardContent className="space-y-4">
                            {passwordError && (
                                <div className="flex items-start gap-3 rounded-md bg-red-500/10 border border-red-500/20 p-3 text-sm text-red-600 dark:text-red-400">
                                    <AlertCircle className="h-4 w-4 mt-0.5 flex-shrink-0" />
                                    <span>{passwordError}</span>
                                </div>
                            )}

                            <div className="space-y-2">
                                <Label htmlFor="fullName" className="text-slate-700 dark:text-zinc-400">
                                    Full Name
                                </Label>
                                <div className="relative">
                                    <User className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-500 dark:text-zinc-500" />
                                    <Input
                                        id="fullName"
                                        type="text"
                                        placeholder="John Doe"
                                        value={fullName}
                                        onChange={(e) => setFullName(e.target.value)}
                                        className="pl-10 bg-white dark:bg-zinc-950/50 border-slate-200 dark:border-zinc-800 focus:border-cyan-500"
                                        disabled={isSubmitting}
                                        required
                                    />
                                </div>
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="email" className="text-slate-700 dark:text-zinc-400">
                                    Email Address
                                </Label>
                                <div className="relative">
                                    <Mail className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-500 dark:text-zinc-500" />
                                    <Input
                                        id="email"
                                        type="email"
                                        value={inviteData.email}
                                        disabled
                                        className="pl-10 bg-slate-50 dark:bg-zinc-900/50 border-slate-200 dark:border-zinc-800 text-slate-600 dark:text-zinc-400"
                                    />
                                </div>
                                <p className="text-xs text-slate-500 dark:text-zinc-500">
                                    This email will be used for your account
                                </p>
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="password" className="text-slate-700 dark:text-zinc-400">
                                    Password
                                </Label>
                                <div className="relative">
                                    <Lock className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-500 dark:text-zinc-500" />
                                    <Input
                                        id="password"
                                        type={showPassword ? 'text' : 'password'}
                                        placeholder="Create a strong password"
                                        value={password}
                                        onChange={(e) => setPassword(e.target.value)}
                                        className="pl-10 pr-10 bg-white dark:bg-zinc-950/50 border-slate-200 dark:border-zinc-800 focus:border-cyan-500"
                                        disabled={isSubmitting}
                                        required
                                    />
                                    <button
                                        type="button"
                                        onClick={() => setShowPassword(!showPassword)}
                                        className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-500 hover:text-slate-700 dark:text-zinc-500 dark:hover:text-zinc-300 transition-colors"
                                        disabled={isSubmitting}
                                    >
                                        {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                                    </button>
                                </div>
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="confirmPassword" className="text-slate-700 dark:text-zinc-400">
                                    Confirm Password
                                </Label>
                                <div className="relative">
                                    <Lock className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-500 dark:text-zinc-500" />
                                    <Input
                                        id="confirmPassword"
                                        type={showConfirmPassword ? 'text' : 'password'}
                                        placeholder="Confirm your password"
                                        value={confirmPassword}
                                        onChange={(e) => setConfirmPassword(e.target.value)}
                                        className="pl-10 pr-10 bg-white dark:bg-zinc-950/50 border-slate-200 dark:border-zinc-800 focus:border-cyan-500"
                                        disabled={isSubmitting}
                                        required
                                    />
                                    <button
                                        type="button"
                                        onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                                        className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-500 hover:text-slate-700 dark:text-zinc-500 dark:hover:text-zinc-300 transition-colors"
                                        disabled={isSubmitting}
                                    >
                                        {showConfirmPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                                    </button>
                                </div>
                            </div>

                            <div className="rounded-lg border border-slate-200 dark:border-zinc-800 bg-slate-50 dark:bg-zinc-900/50 p-4 space-y-2">
                                <div className="flex items-center gap-2">
                                    <Shield className="h-4 w-4 text-slate-500 dark:text-zinc-500" />
                                    <p className="text-sm text-slate-700 dark:text-zinc-300">
                                        <span className="font-medium">Role:</span> {inviteData.role}
                                    </p>
                                </div>
                                <p className="text-xs text-slate-500 dark:text-zinc-500">
                                    You will be added to the team with {inviteData.role} permissions.
                                </p>
                            </div>
                        </CardContent>

                        <CardFooter className="flex flex-col gap-4">
                            <div className="flex gap-3 w-full">
                                <Button
                                    type="button"
                                    variant="outline"
                                    onClick={() => navigate({ to: "/" })}
                                    disabled={isSubmitting}
                                    className="flex-1 border-slate-200 dark:border-zinc-800"
                                >
                                    Cancel
                                </Button>
                                <Button
                                    type="submit"
                                    disabled={isSubmitting}
                                    className="flex-1 bg-gradient-to-r from-cyan-500 to-cyan-600 hover:from-cyan-600 hover:to-cyan-700 text-white font-medium shadow-lg shadow-cyan-500/25"
                                >
                                    {isSubmitting ? (
                                        <>
                                            <div className="mr-2 h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent" />
                                            Creating Account...
                                        </>
                                    ) : (
                                        "Create Account & Join Team"
                                    )}
                                </Button>
                            </div>
                        </CardFooter>
                    </form>
                </Card>
            </div>
        );
    }

    return null;
}