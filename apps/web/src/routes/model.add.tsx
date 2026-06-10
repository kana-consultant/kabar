import { createFileRoute } from '@tanstack/react-router'
import CreateModelPage from '@/pages/ai-management/Model/create'
export const Route = createFileRoute('/model/add')({
    component: RouteComponent,
})

function RouteComponent() {
    return CreateModelPage
}
