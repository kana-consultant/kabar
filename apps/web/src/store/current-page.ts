// store/current-page.ts
let currentPageCache = '/';

export const pageTracker = {
    get: (): string => {
        return currentPageCache; 
    },
    set: (path: string): void => {
        currentPageCache = path;
    }
};