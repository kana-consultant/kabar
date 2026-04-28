// LocalStorage wrapper - database sementara

const STORAGE_KEYS = {
    PRODUCTS: 'seo_products',
    DRAFTS: 'seo_drafts',
    HISTORY: 'seo_history',
    USERS: 'seo_users',
    TEAMS: 'seo_teams',
    CURRENT_USER: 'seo_current_user',
    SETTINGS: 'seo_settings',
    API_KEYS: 'seo_api_keys',
} as const;

type BaseEntity = {
    id: string
    createdAt: string
    updatedAt: string
}

export class LocalStorageDB {

    // Get all items
    static get<T>(key: string): T[] {
        const data = localStorage.getItem(key);
        return data ? JSON.parse(data) : [];
    }

    // Get single item by id
    static getById<T extends BaseEntity>(key: string, id: string): T | null {
        const items = this.get<T>(key);
        return items.find(item => item.id === id) || null;
    }

    // Add item
    static add<T extends BaseEntity>(
        key: string,
        item: Omit<T, keyof BaseEntity>
    ): T {

        const items = this.get<T>(key);

        const newItem: T = {
            ...item,
            id: crypto.randomUUID(),
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
        } as T;

        items.push(newItem);

        localStorage.setItem(key, JSON.stringify(items));

        return newItem;
    }

    // Update item
    static update<T extends { id: string; updatedAt?: string }>(key: string, id: string, updates: Partial<T>): T | null {
        const items = this.get<T>(key);
        const index = items.findIndex(item => item.id === id);
        if (index === -1) return null;

        items[index] = {
            ...items[index],
            ...updates,
            updatedAt: new Date().toISOString(),
        };
        localStorage.setItem(key, JSON.stringify(items));
        return items[index];
    }

    // Delete item
    static delete<T extends BaseEntity>(
        key: string,
        id: string
    ): boolean {

        const items = this.get<T>(key);

        const filtered = items.filter(item => item.id !== id);

        if (filtered.length === items.length) return false;

        localStorage.setItem(key, JSON.stringify(filtered));

        return true;
    }

    // Clear all
    static clear(key: string): void {
        localStorage.setItem(key, JSON.stringify([]));
    }

    // Seed default data
    static seed<T>(key: string, data: T[]): void {
        const existing = this.get<T>(key);
        if (existing.length === 0) {
            localStorage.setItem(key, JSON.stringify(data));
        }
    }
}

export { STORAGE_KEYS };