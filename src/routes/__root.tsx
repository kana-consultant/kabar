import { Layout } from '@/components/layout/Layout'
import { createRootRoute, } from '@tanstack/react-router'

export const Route = createRootRoute({
    component: () => (
        <>
            <Layout />
        </>
    ),
})