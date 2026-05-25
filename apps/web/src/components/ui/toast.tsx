import { useState } from "react"
import { 
  Toast,  
  ToastTitle, 
  ToastDescription,
  ToastClose,
  ToastAction
} from "@kana-consultant/ui-kit"

export function MyComponent() {
  const [open, setOpen] = useState(false)

  return (
    <>
      <button onClick={() => setOpen(true)}>Show Toast</button>
      
      <Toast open={open} onOpenChange={setOpen} tone="success">
        <ToastTitle>Sukses!</ToastTitle>
        <ToastDescription>Data berhasil disimpan</ToastDescription>
        <ToastClose />
        <ToastAction altText="Undo">Undo</ToastAction>
      </Toast>
    </>
  )
}