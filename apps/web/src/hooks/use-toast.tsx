// hooks/use-toast.tsx
import { createContext, useContext, useState,type ReactNode } from "react"
import { 
  Toast, 
  ToastTitle, 
  ToastDescription,
  ToastClose,
  ToastProvider
} from "@kana-consultant/ui-kit"

interface ToastMessage {
  id: string
  title: string
  description?: string
  tone?: "default" | "info" | "success" | "warning" | "danger"
}

export interface ToastContextType {
  showToast: (message: Omit<ToastMessage, "id">) => void
  success: (title: string, options?: { description?: string }) => void
  error: (title: string, options?: { description?: string }) => void
  info: (title: string, options?: { description?: string }) => void
  warning: (title: string, options?: { description?: string }) => void
}

const ToastContext = createContext<ToastContextType | undefined>(undefined)

export function ToastProviderWrapper({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<ToastMessage[]>([])

  const removeToast = (id: string) => {
    setToasts(prev => prev.filter(toast => toast.id !== id))
  }

  const showToast = (message: Omit<ToastMessage, "id">) => {
    const id = Math.random().toString(36).substring(7)
    setToasts(prev => [...prev, { ...message, id }])
    
    // Auto remove after 5 seconds
    setTimeout(() => removeToast(id), 5000)
  }

  const success = (title: string, options?: { description?: string }) => {
    showToast({ title, description: options?.description, tone: "success" })
  }

  const error = (title: string, options?: { description?: string }) => {
    showToast({ title, description: options?.description, tone: "danger" })
  }

  const info = (title: string, options?: { description?: string }) => {
    showToast({ title, description: options?.description, tone: "info" })
  }

  const warning = (title: string, options?: { description?: string }) => {
    showToast({ title, description: options?.description, tone: "warning" })
  }

  return (
    <ToastContext.Provider value={{ showToast, success, error, info, warning }}>
      <ToastProvider>
        {children}
        {/* Render all toasts */}
        {toasts.map(toast => (
          <Toast
            key={toast.id}
            open={true}
            onOpenChange={(open) => !open && removeToast(toast.id)}
            tone={toast.tone}
          >
            <ToastTitle>{toast.title}</ToastTitle>
            {toast.description && (
              <ToastDescription>{toast.description}</ToastDescription>
            )}
            <ToastClose />
          </Toast>
        ))}
      </ToastProvider>
    </ToastContext.Provider>
  )
}

export function useToast() {
  const context = useContext(ToastContext)
  if (!context) {
    throw new Error("useToast must be used within ToastProviderWrapper")
  }
  return context
}