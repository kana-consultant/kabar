import { StrictMode } from 'react'
import { Toaster } from "sonner";
import { createRoot } from 'react-dom/client'
import {
  createRouter,
  RouterProvider,
  createRootRoute,
  createRoute,
} from '@tanstack/react-router'
import { ThemeProvider } from './components/theme-provider'
import { AuthProvider } from './contexts/AuthContext'
import './index.css'

// Import Layouts
import { Layout } from './components/layout/Layout'           // Untuk protected routes
import { AuthLayout } from './components/layout/AuthLayout'   // Untuk login/register

// Import Pages
import { Dashboard } from './routes';
import { Generate } from './routes/generate';
import Products from './routes/products';
import ProductEdit from './routes/products.$id.edit';
import ProductAdd from './routes/products.add';
import History from './routes/history';
import Settings from './routes/settings';
import Drafts from './routes/Drafts'
import Schedule from './routes/schedule';
import Login from './pages/Login'
import Register from './pages/Register'

import { Outlet } from '@tanstack/react-router';

// 1. Buat root route (tanpa layout)
const rootRoute = createRootRoute({
  component: () => <Outlet />,  // Root hanya sebagai outlet
})

// 2. Buat auth routes (pakai AuthLayout)
const authLayoutRoute = createRoute({
  getParentRoute: () => rootRoute,
  id: 'auth',
  component: AuthLayout,
})

const loginRoute = createRoute({
  getParentRoute: () => authLayoutRoute,
  path: '/login',
  component: Login,
})

const registerRoute = createRoute({
  getParentRoute: () => authLayoutRoute,
  path: '/register',
  component: Register,
})

// 3. Buat protected routes (pakai Layout dengan sidebar)
const protectedLayoutRoute = createRoute({
  getParentRoute: () => rootRoute,
  id: 'protected',
  component: Layout,
})

const indexRoute = createRoute({
  getParentRoute: () => protectedLayoutRoute,
  path: '/',
  component: Dashboard,
})

const generateRoute = createRoute({
  getParentRoute: () => protectedLayoutRoute,
  path: '/generate',
  component: Generate,
})

const productsRoute = createRoute({
  getParentRoute: () => protectedLayoutRoute,
  path: '/products',
  component: Products,
})

const productAddRoute = createRoute({
  getParentRoute: () => protectedLayoutRoute,
  path: '/products/add',
  component: ProductAdd,
})

const productEditRoute = createRoute({
  getParentRoute: () => protectedLayoutRoute,
  path: '/products/$id/edit',
  component: ProductEdit,
})

const historyRoute = createRoute({
  getParentRoute: () => protectedLayoutRoute,
  path: '/history',
  component: History,
})

const settingsRoute = createRoute({
  getParentRoute: () => protectedLayoutRoute,
  path: '/settings',
  component: Settings,
})

const draftsRoute = createRoute({
  getParentRoute: () => protectedLayoutRoute,
  path: '/drafts',
  component: Drafts,
})

const scheduleRoute = createRoute({
  getParentRoute: () => protectedLayoutRoute,
  path: '/schedule',
  component: Schedule,
})

// 4. Gabungkan semua route
const routeTree = rootRoute.addChildren([
  authLayoutRoute.addChildren([loginRoute, registerRoute]),
  protectedLayoutRoute.addChildren([
    indexRoute,
    generateRoute,
    productsRoute,
    productAddRoute,
    productEditRoute,
    historyRoute,
    draftsRoute,
    settingsRoute,
    scheduleRoute,
  ]),
])

// 5. Buat router
const router = createRouter({ routeTree })

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}

// 6. Render
createRoot(document.getElementById('root')!).render(
  <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
    <AuthProvider>
      <RouterProvider router={router} />
      <Toaster position="top-right" richColors closeButton />
    </AuthProvider>
  </ThemeProvider>
) 