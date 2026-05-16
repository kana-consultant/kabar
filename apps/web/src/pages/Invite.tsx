// invite.tsx
import { acceptInvite, verifyInviteToken } from "@/services/team";
import { useEffect, useState } from "react";
import { useNavigate } from "@tanstack/react-router";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Loader2, CheckCircle, XCircle } from "lucide-react";
import type { TeamInvite } from "@/services/team/types";
import { getToken } from "@/services/auth";

export default function InvitePage() {
    const navigate = useNavigate();
    const searchParams = new URLSearchParams(location.search);
    const token = searchParams.get("token") as string;

    const [isLoading, setIsLoading] = useState<boolean>(true);
    const [isVerifying, setIsVerifying] = useState<boolean>(true);
    const [inviteData, setInviteData] = useState<TeamInvite | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [isAccepted, setIsAccepted] = useState<boolean>(false);

    // Form state untuk user yang belum punya account
    const [showPasswordForm, setShowPasswordForm] = useState<boolean>(false);
    const [password, setPassword] = useState<string>("");
    const [confirmPassword, setConfirmPassword] = useState<string>("");
    const [fullName, setFullName] = useState<string>("");
    const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
    const [passwordError, setPasswordError] = useState<string | null>(null);

    // Step 1: Verify invite token
    const verifyToken = async (token: string) => {
        try {
            console.log("HASIL TOKEN = ", token)
            const result = await verifyInviteToken(token);
            setInviteData(result as any);
            setError(null);

            // Cek apakah user sudah login atau belum
            const isLoggedIn = getToken(); // Implement sesuai auth system

            if (!isLoggedIn) {
                setShowPasswordForm(true);
            } else {
                // User sudah login, langsung accept
                await acceptInviteAndRedirect(token);
            }
        } catch (err: any) {
            setError(err?.message || "Failed to verify invitation");
        } finally {
            setIsVerifying(false);
            setIsLoading(false);
        }
    };

    // Step 2: Accept invite and redirect
    const acceptInviteAndRedirect = async (token: string) => {
        setIsSubmitting(true);
        try {
            const response = await acceptInvite(token);

            setIsAccepted(true);
            // Redirect after 2 seconds
            setTimeout(() => {
                navigate({
                    to: `/settings`,
                });
            }, 2000);
        } catch (err: any) {
            setError(err?.message || "Failed to accept invitation");
        } finally {
            setIsSubmitting(false);
        }
    };

    // Handle form submit untuk registrasi
    const handleRegisterAndAccept = async (e: React.FormEvent) => {
        e.preventDefault();
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

        await acceptInviteAndRedirect(token);
    };

    // Handle langsung accept (untuk user yang sudah login)
    const handleAccept = async () => {
        await acceptInviteAndRedirect(token);
    };

    useEffect(() => {
        if (!token) {
            setError("No invitation token provided");
            setIsLoading(false);
            setIsVerifying(false);
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
                        <Loader2 className="h-8 w-8 animate-spin text-primary" />
                        <p className="mt-4 text-muted-foreground">Verifying invitation...</p>
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
                            <XCircle className="h-12 w-12 text-red-500" />
                        </div>
                        <CardTitle className="text-center text-red-600">Invalid Invitation</CardTitle>
                        <CardDescription className="text-center">
                            {error}
                        </CardDescription>
                    </CardHeader>
                    <CardFooter className="flex justify-center">
                        <Button onClick={() => navigate({ to: "/" })}>
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
                            <CheckCircle className="h-12 w-12 text-green-500" />
                        </div>
                        <CardTitle className="text-center text-green-600">Successfully Joined!</CardTitle>
                        <CardDescription className="text-center">
                            You have successfully joined {inviteData?.TeamName || "the team"}.
                            Redirecting to team dashboard...
                        </CardDescription>
                    </CardHeader>
                    <CardContent className="flex justify-center">
                        <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                    </CardContent>
                </Card>
            </div>
        );
    }

    // Show password form for new users
    if (showPasswordForm && inviteData) {
        return (
             <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-950 dark:to-slate-900 p-4">
                <Card className="w-full max-w-md">
                    <CardHeader>
                        <CardTitle>Join {inviteData.TeamName}</CardTitle>
                        <CardDescription>
                            You've been invited to join {inviteData.TeamName} as {inviteData.role}.
                            Please complete your registration to continue.
                        </CardDescription>
                    </CardHeader>
                    <form onSubmit={handleRegisterAndAccept}>
                        <CardContent className="space-y-4">
                            <div className="space-y-2">
                                <Label htmlFor="fullName">Full Name</Label>
                                <Input
                                    id="fullName"
                                    type="text"
                                    placeholder="John Doe"
                                    value={fullName}
                                    onChange={(e) => setFullName(e.target.value)}
                                    disabled={isSubmitting}
                                    required
                                />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="email">Email</Label>
                                <Input
                                    id="email"
                                    type="email"
                                    value={inviteData.email}
                                    disabled
                                    className="bg-gray-100"
                                />
                                <p className="text-xs text-muted-foreground">
                                    This email will be used for your account
                                </p>
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="password">Password</Label>
                                <Input
                                    id="password"
                                    type="password"
                                    placeholder="Create a password"
                                    value={password}
                                    onChange={(e) => setPassword(e.target.value)}
                                    disabled={isSubmitting}
                                    required
                                />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="confirmPassword">Confirm Password</Label>
                                <Input
                                    id="confirmPassword"
                                    type="password"
                                    placeholder="Confirm your password"
                                    value={confirmPassword}
                                    onChange={(e) => setConfirmPassword(e.target.value)}
                                    disabled={isSubmitting}
                                    required
                                />
                            </div>
                            {passwordError && (
                                <Alert variant="destructive">
                                    <AlertDescription>{passwordError}</AlertDescription>
                                </Alert>
                            )}
                            <div className="rounded-md bg-gray-50 p-3">
                                <p className="text-sm">
                                    <span className="font-medium">Role:</span> {inviteData.role}
                                </p>
                                <p className="text-sm text-muted-foreground mt-1">
                                    You will be added to the team with {inviteData.role} permissions.
                                </p>
                            </div>
                        </CardContent>
                        <CardFooter className="flex gap-3">
                            <Button
                                type="button"
                                variant="outline"
                                onClick={() => navigate({ to: "/" })}
                                disabled={isSubmitting}
                            >
                                Cancel
                            </Button>
                            <Button type="submit" disabled={isSubmitting}>
                                {isSubmitting ? (
                                    <>
                                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                        Processing...
                                    </>
                                ) : (
                                    "Create Account & Join Team"
                                )}
                            </Button>
                        </CardFooter>
                    </form>
                </Card>
            </div>
        );
    }

    // Untuk user yang sudah login
    if (inviteData && !showPasswordForm) {
        return (
           <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-950 dark:to-slate-900 p-4">
                <Card className="w-full max-w-md">
                    <CardHeader>
                        <CardTitle>Join {inviteData.TeamName}</CardTitle>
                        <CardDescription>
                            You've been invited to join {inviteData.TeamName} as {inviteData.role}
                        </CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        <div className="rounded-md bg-gray-50 p-3">
                            <p className="text-sm">
                                <span className="font-medium">Team:</span> {inviteData.TeamName}
                            </p>
                            <p className="text-sm mt-1">
                                <span className="font-medium">Role:</span> {inviteData.role}
                            </p>
                            <p className="text-sm text-muted-foreground mt-2">
                                Click accept to join this team.
                            </p>
                        </div>
                    </CardContent>
                    <CardFooter className="flex gap-3">
                        <Button
                            variant="outline"
                            onClick={() => navigate({ to: "/" })}
                            disabled={isSubmitting}
                        >
                            Cancel
                        </Button>
                        <Button onClick={handleAccept} disabled={isSubmitting}>
                            {isSubmitting ? (
                                <>
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                    Joining...
                                </>
                            ) : (
                                "Accept Invitation"
                            )}
                        </Button>
                    </CardFooter>
                </Card>
            </div>
        );
    }

    return null;
}