import { useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";

export function PreferencesTab() {
    const [darkMode, setDarkMode] = useState(false);
    const [emailNotifications, setEmailNotifications] = useState(true);
    const [autoSave, setAutoSave] = useState(true);

    return (
        <Card>
            <CardHeader>
                <CardTitle>Preferences</CardTitle>
                <CardDescription>Pengaturan aplikasi</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                    <div>
                        <Label>Dark Mode</Label>
                        <p className="text-xs text-slate-500">Tampilan gelap untuk aplikasi</p>
                    </div>
                    <Switch checked={darkMode} onCheckedChange={setDarkMode} />
                </div>
                <div className="flex items-center justify-between">
                    <div>
                        <Label>Email Notifications</Label>
                        <p className="text-xs text-slate-500">Terima notifikasi via email</p>
                    </div>
                    <Switch checked={emailNotifications} onCheckedChange={setEmailNotifications} />
                </div>
                <div className="flex items-center justify-between">
                    <div>
                        <Label>Auto-save Draft</Label>
                        <p className="text-xs text-slate-500">Simpan draft secara otomatis</p>
                    </div>
                    <Switch checked={autoSave} onCheckedChange={setAutoSave} />
                </div>
            </CardContent>
        </Card>
    );
}