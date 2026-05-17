import { apiClient } from '../api';

export async function changePasswordApi(OldPassword: string, NewPassword: string): Promise<void> {
    const payload = {OldPassword, NewPassword}
    return await apiClient.post("/auth/change-password",payload)
}