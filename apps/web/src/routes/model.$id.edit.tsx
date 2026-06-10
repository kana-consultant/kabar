import { createFileRoute } from '@tanstack/react-router'
import EditModelPage from '@/pages/ai-management/Model/[id]/edit'
export const Route = createFileRoute('/model/$id/edit')({
  component: RouteComponent,
})

function RouteComponent() {
  return EditModelPage
}
