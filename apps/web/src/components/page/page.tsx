// components/ui/page.tsx

import type { ReactNode } from "react";

interface PageProps {
    children: ReactNode;
    className?: string;
}

export function Page({ children, className = "" }: PageProps) {
    return (
        <div className={`flex flex-col h-full ${className}`}>
            {children}
        </div>
    );
}

interface PageHeaderProps {
    children: ReactNode;
    className?: string;
}

export function PageHeader({ children, className = "" }: PageHeaderProps) {
    return (
        <div className={`border-b pb-4 mb-6 ${className}`}>
            {children}
        </div>
    );
}

interface PageTitleProps {
    children: ReactNode;
    className?: string;
}

export function PageTitle({ children, className = "" }: PageTitleProps) {
    return (
        <h1 className={`text-3xl font-bold tracking-tight ${className}`}>
            {children}
        </h1>
    );
}

interface PageDescriptionProps {
    children: ReactNode;
    className?: string;
}

export function PageDescription({ children, className = "" }: PageDescriptionProps) {
    return (
        <p className={`text-sm text-gray-500 dark:text-gray-400 mt-1 ${className}`}>
            {children}
        </p>
    );
}

interface PageContentProps {
    children: ReactNode;
    className?: string;
}

export function PageContent({ children, className = "" }: PageContentProps) {
    return (
        <div className={`flex-1 ${className}`}>
            {children}
        </div>
    );
}