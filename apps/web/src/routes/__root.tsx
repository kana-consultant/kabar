import { Layout } from "@/pages/layout/Layout"
import { createRootRoute, Outlet } from '@tanstack/react-router'

export const Route = createRootRoute({
    component: Root,
})

function Root() {
    const isLanding = window.location.pathname === '/landing'
    
    if (isLanding) {
        return <Outlet />
    }
    
    return <Layout />
}