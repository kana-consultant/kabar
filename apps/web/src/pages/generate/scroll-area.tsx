// components/ui/scroll-area.tsx
import * as React from "react";
import { cn } from "@/lib/utils";

interface ScrollAreaProps extends React.HTMLAttributes<HTMLDivElement> {
    className?: string;
    children?: React.ReactNode;
    orientation?: "vertical" | "horizontal" | "both";
    maxHeight?: string | number;
    maxWidth?: string | number;
    showScrollbar?: boolean;
}

export const ScrollArea = React.forwardRef<HTMLDivElement, ScrollAreaProps>(
    ({ 
        className, 
        children, 
        orientation = "vertical",
        maxHeight = "100%",
        maxWidth = "100%",
        showScrollbar = true,
        ...props 
    }, ref) => {
        const [isHovered, setIsHovered] = React.useState(false);
        const [scrollbarVisible, setScrollbarVisible] = React.useState(false);
        const containerRef = React.useRef<HTMLDivElement>(null);

        React.useEffect(() => {
            const container = containerRef.current;
            if (!container) return;

            const checkScroll = () => {
                const hasVerticalScroll = container.scrollHeight > container.clientHeight;
                const hasHorizontalScroll = container.scrollWidth > container.clientWidth;
                setScrollbarVisible(hasVerticalScroll || hasHorizontalScroll);
            };

            checkScroll();
            window.addEventListener('resize', checkScroll);
            
            if (typeof ResizeObserver !== 'undefined') {
                const resizeObserver = new ResizeObserver(checkScroll);
                resizeObserver.observe(container);
                return () => resizeObserver.disconnect();
            }

            return () => window.removeEventListener('resize', checkScroll);
        }, [children]);

        // Build scrollbar classes based on state
        const getScrollbarClasses = () => {
            const baseClasses = [
                "h-full w-full overflow-auto scroll-smooth",
            ];

            if (showScrollbar) {
                baseClasses.push(
                    // Webkit scrollbar styles
                    "[&::-webkit-scrollbar]:w-1.5",
                    "[&::-webkit-scrollbar]:h-1.5",
                    "[&::-webkit-scrollbar-track]:bg-transparent",
                    "[&::-webkit-scrollbar-thumb]:rounded-full",
                );

                if (isHovered) {
                    baseClasses.push(
                        "[&::-webkit-scrollbar-thumb]:bg-slate-300",
                        "dark:[&::-webkit-scrollbar-thumb]:bg-slate-600",
                        "hover:[&::-webkit-scrollbar-thumb]:bg-slate-400",
                        "dark:hover:[&::-webkit-scrollbar-thumb]:bg-slate-500",
                        // Firefox
                        "scrollbar-width:thin",
                        "scrollbar-color:slate-300 transparent",
                        "dark:scrollbar-color:slate-600 transparent",
                        "hover:scrollbar-color:slate-400 transparent",
                        "dark:hover:scrollbar-color:slate-500 transparent",
                    );
                } else {
                    baseClasses.push(
                        "[&::-webkit-scrollbar-thumb]:bg-transparent",
                        "dark:[&::-webkit-scrollbar-thumb]:bg-transparent",
                        // Firefox
                        "scrollbar-width:thin",
                        "scrollbar-color:transparent transparent",
                    );
                }
            }

            return baseClasses.join(' ');
        };

        return (
            <div
                ref={ref}
                className={cn(
                    "relative overflow-hidden",
                    `${className}`
                )}
                style={{
                    maxHeight: typeof maxHeight === 'number' ? `${maxHeight}px` : maxHeight,
                    maxWidth: typeof maxWidth === 'number' ? `${maxWidth}px` : maxWidth,
                }}
                onMouseEnter={() => setIsHovered(true)}
                onMouseLeave={() => setIsHovered(false)}
                {...props}
            >
                <div
                    ref={containerRef}
                    className={getScrollbarClasses()}
                    style={{
                        scrollBehavior: 'smooth',
                    }}
                >
                    {children}
                </div>
            </div>
        );
    }
);

ScrollArea.displayName = "ScrollArea";