import { createFileRoute } from '@tanstack/react-router'
export const Route = createFileRoute('/model/$id')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/model/$id"!</div>
}
