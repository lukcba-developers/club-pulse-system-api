"use client"

import * as React from "react"
import { cn } from "@/lib/utils"

interface CheckboxProps extends React.InputHTMLAttributes<HTMLInputElement> {
    onCheckedChange?: (checked: boolean) => void
}

const Checkbox = React.forwardRef<HTMLInputElement, CheckboxProps>(
    ({ className, onCheckedChange, ...props }, ref) => {
        const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
            onCheckedChange?.(e.target.checked)
            props.onChange?.(e)
        }

        return (
            <input
                type="checkbox"
                ref={ref}
                className={cn(
                    "h-4 w-4 shrink-0 rounded border-gray-300 text-brand-600 focus:ring-brand-500 focus:ring-2 focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50",
                    className
                )}
                onChange={handleChange}
                {...props}
            />
        )
    }
)
Checkbox.displayName = "Checkbox"

export { Checkbox }
