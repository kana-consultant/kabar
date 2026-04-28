import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import type { User } from "@/types/user";

const roleLabels: Record<string, string> = {
    admin: "Administrator",
    manager: "Manager",
    editor: "Editor",
    viewer: "Viewer",
};

const roleColors: Record<string, string> = {
    admin: "text-red-600 bg-red-50 dark:bg-red-950",
    manager: "text-blue-600 bg-blue-50 dark:bg-blue-950",
    editor: "text-green-600 bg-green-50 dark:bg-green-950",
    viewer: "text-gray-600 bg-gray-50 dark:bg-gray-950",
};

interface ProfileTabProps {
    currentUser: User | null;
}

export function ProfileTab({ currentUser }: ProfileTabProps) {
    return (
        <Card>
            <CardHeader>
                <CardTitle>Profile Information</CardTitle>
                <CardDescription>Informasi akun Anda</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
                <div className="flex items-center gap-4">
                    <div className="flex h-20 w-20 items-center justify-center rounded-full bg-blue-600 text-3xl font-bold text-white">
                        {currentUser?.name.charAt(0)}
                    </div>
                    <div>
                        <p className="text-lg font-semibold">{currentUser?.name}</p>
                        <p className="text-sm text-slate-500">{currentUser?.email}</p>
                        <span className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs ${roleColors[currentUser?.role || "viewer"]}`}>
                            {roleLabels[currentUser?.role || "viewer"]}
                        </span>
                    </div>
                </div>
                <div className="grid gap-4 sm:grid-cols-2">
                    <div>
                        <Label>Name</Label>
                        <Input defaultValue={currentUser?.name} className="mt-1" disabled />
                    </div>
                    <div>
                        <Label>Email</Label>
                        <Input defaultValue={currentUser?.email} className="mt-1" disabled />
                    </div>
                </div>
                {/* <Button>Update Profile</Button> */}
            </CardContent>
        </Card>
    );
}