// main.tsx

import { ToastProviderWrapper } from "@/hooks/use-toast"
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { createRoot } from 'react-dom/client'
import {
  createRouter,
  RouterProvider,
  createRootRoute,
  createRoute,
} from '@tanstack/react-router'
import { AuthProvider } from './contexts/AuthContext'
import './index.css'

// Import Layouts
import { Layout } from "@/pages/layout/Layout"
import { AuthLayout } from "@/pages/layout/AuthLayout"

// Import Pages
import { Dashboard } from './routes';
import Generate from './pages/generate/generate';
import Products from './pages/home/Products';
import ProductEdit from './routes/products.$id.edit';
import ProductAdd from './routes/products.add';
import History from './pages/history/History';
import Settings from './pages/settings/Settings';
import Drafts from './pages/draft/Drafts'
import Schedule from './pages/schedule/Schedule';
import Login from './pages/Login'
import Register from './pages/Register'
import InvitePage from "./pages/Invite";
import Landing from './pages/Landing'
import AIManagement from "./pages/ai-management/AIManagement";

// Import AI Management - Provider Pages
import ProviderListPage from "./pages/ai-management/Provider/index";
import CreateProviderPage from "./pages/ai-management/Provider/create";
import ProviderDetailPage from "./pages/ai-management/Provider/[id]/index";
import EditProviderPage from "./pages/ai-management/Provider/[id]/edit";

// Import AI Management - Model Pages
import CreateModelPage from "./pages/ai-management/Model/create";
import EditModelPage from "./pages/ai-management/Model/[id]/edit";

import { Outlet } from '@tanstack/react-router';
import SitemapPage from "./pages/sitemap/Sitemap";

// 1. Buat root route (tanpa layout)
const rootRoute = createRootRoute({
  component: () => <Outlet />,
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

const inviteRoute = createRoute({
  getParentRoute: () => authLayoutRoute,
  path: '/invite',
  component: InvitePage,
})

// 3. Buat protected routes (pakai Layout dengan sidebar)
const protectedLayoutRoute = createRoute({
  getParentRoute: () => rootRoute,
  id: 'protected',
  component: Layout,
})

const indexRoute = createRoute({
  getParentRoute: () => protectedLayoutRoute,
  path: '/dashboard',
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

// AI Management Route
const aiManagementRoute = createRoute({
  getParentRoute: () => protectedLayoutRoute,
  path: '/ai-management',
  component: AIManagement,
})

// Provider Routes (nested under protected layout)
const providerLayoutRoute = createRoute({
  getParentRoute: () => protectedLayoutRoute,
  path: '/provider',
})


// Provider Routes (nested under protected layout)
const SitemapLayoutRoute = createRoute({
  getParentRoute: () => protectedLayoutRoute,
  path: '/sitemap',
  component: SitemapPage,
})

const providerListRoute = createRoute({
  getParentRoute: () => providerLayoutRoute,
  path: '/',
  component: ProviderListPage,
})

const providerAddRoute = createRoute({
  getParentRoute: () => providerLayoutRoute,
  path: '/add',
  component: CreateProviderPage,
})

const providerDetailRoute = createRoute({
  getParentRoute: () => providerLayoutRoute,
  path: '/$id',
  component: ProviderDetailPage,
})

const providerEditRoute = createRoute({
  getParentRoute: () => providerLayoutRoute,
  path: '/$id/edit',
  component: EditProviderPage,
})

// Model Routes (nested under protected layout)
const modelLayoutRoute = createRoute({
  getParentRoute: () => protectedLayoutRoute,
  path: '/model',
})


const modelAddRoute = createRoute({
  getParentRoute: () => modelLayoutRoute,
  path: '/add',
  component: CreateModelPage,
})



const modelEditRoute = createRoute({
  getParentRoute: () => modelLayoutRoute,
  path: '/$id/edit',
  component: EditModelPage,
})

// Landing route (public)
const landingRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  component: Landing,
})

// 4. Gabungkan semua route
const routeTree = rootRoute.addChildren([
  landingRoute,
  authLayoutRoute.addChildren([loginRoute, registerRoute, inviteRoute]),
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
    aiManagementRoute,
    SitemapLayoutRoute,
    // Provider routes
    providerLayoutRoute.addChildren([
      providerListRoute,
      providerAddRoute,
      providerDetailRoute,
      providerEditRoute,
    ]),
    // Model routes
    modelLayoutRoute.addChildren([
    
      modelAddRoute,
    
      modelEditRoute,
    ]),
  ]),
])

// 5. Buat router
const router = createRouter({ routeTree })

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}

const queryClient = new QueryClient();

// 6. Render
createRoot(document.getElementById('root')!).render(
  <ToastProviderWrapper>
    <AuthProvider>
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    </AuthProvider>
  </ToastProviderWrapper>
)