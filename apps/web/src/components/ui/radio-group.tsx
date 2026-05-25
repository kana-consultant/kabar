// components/ui/radio-group.tsx
import * as React from "react"
import { cn } from "@/lib/utils"

interface RadioGroupProps extends React.HTMLAttributes<HTMLDivElement> {
    value?: string
    onValueChange?: (value: string) => void
    defaultValue?: string
    disabled?: boolean
    name?: string
}

interface RadioGroupItemProps extends Omit<React.InputHTMLAttributes<HTMLInputElement>, "type"> {
    value: string
    id?: string
}

const RadioGroupContext = React.createContext<{
    value?: string
    onValueChange?: (value: string) => void
    disabled?: boolean
    name?: string
}>({})

export function RadioGroup({
    children,
    value,
    onValueChange,
    defaultValue,
    disabled = false,
    name,
    className,
    ...props
}: RadioGroupProps) {
    const [selectedValue, setSelectedValue] = React.useState(defaultValue || value || "")
    
    const handleValueChange = React.useCallback((newValue: string) => {
        if (disabled) return
        setSelectedValue(newValue)
        onValueChange?.(newValue)
    }, [disabled, onValueChange])
    
    const contextValue = React.useMemo(() => ({
        value: value !== undefined ? value : selectedValue,
        onValueChange: handleValueChange,
        disabled,
        name,
    }), [value, selectedValue, handleValueChange, disabled, name])
    
    return (
        <RadioGroupContext.Provider value={contextValue}>
            <div
                role="radiogroup"
                className={cn("grid w-full gap-3", className as string)}
                {...props}
            >
                {children}
            </div>
        </RadioGroupContext.Provider>
    )
}

export function RadioGroupItem({
    value,
    id,
    className,
    disabled,
    ...props
}: RadioGroupItemProps) {
    const context = React.useContext(RadioGroupContext)
    const isChecked = context.value === value
    const isDisabled = disabled || context.disabled
    const itemId = id || `radio-${value}`
    
    return (
        <div className="flex items-center gap-2">
            <button
                type="button"
                role="radio"
                aria-checked={isChecked}
                aria-disabled={isDisabled}
                data-state={isChecked ? "checked" : "unchecked"}
                data-disabled={isDisabled ? "disabled" : undefined}
                value={value}
                id={itemId}
                disabled={isDisabled}
                onClick={() => !isDisabled && context.onValueChange?.(value)}
                className={cn(
                    "group/radio-item peer relative flex aspect-square size-4 shrink-0 rounded-full border border-input outline-none",
                    "focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50",
                    "disabled:cursor-not-allowed disabled:opacity-50",
                    "aria-invalid:border-destructive aria-invalid:ring-3 aria-invalid:ring-destructive/20",
                    "data-[state=checked]:border-primary data-[state=checked]:bg-primary",
                    "dark:bg-input/30",
                    "dark:data-[state=checked]:bg-primary",
                    className as string
                )}
            >
                {isChecked && (
                    <span className="absolute top-1/2 left-1/2 size-2 -translate-x-1/2 -translate-y-1/2 rounded-full bg-primary-foreground" />
                )}
            </button>
            <input
                type="radio"
                name={context.name}
                value={value}
                checked={isChecked}
                disabled={isDisabled}
                onChange={() => {}}
                className="sr-only"
                {...props}
            />
        </div>
    )
}

interface RadioLabelProps extends React.LabelHTMLAttributes<HTMLLabelElement> {
    htmlFor?: string
}

export function RadioLabel({ children, className, htmlFor, ...props }: RadioLabelProps) {
    return (
        <label
            htmlFor={htmlFor}
            className={cn(
                "text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70",
                className as string
            )}
            {...props}
        >
            {children}
        </label>
    )
}