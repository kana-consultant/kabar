import { Layout } from "@/pages/layout/Layout'
import { createRootRoute, } from '@tanstack/react-router'

export const Route = createRootRoute({
    component: () => (
        <>
            <Layout />
        </>
    ),
})