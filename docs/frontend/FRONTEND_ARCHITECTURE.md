# Frontend Architecture Documentation
## ZTE C320 OLT Management Dashboard

**Project:** Web-based Dashboard for ZTE C320 OLT  
**Architecture Version:** 1.0  
**Last Updated:** January 12, 2026

---

## 📋 Table of Contents

1. [Technology Stack](#technology-stack)
2. [Architecture Overview](#architecture-overview)
3. [Folder Structure](#folder-structure)
4. [Component Architecture](#component-architecture)
5. [State Management](#state-management)
6. [API Integration](#api-integration)
7. [Routing](#routing)
8. [Styling & Theming](#styling--theming)
9. [Performance Optimization](#performance-optimization)
10. [Security Considerations](#security-considerations)
11. [Deployment Architecture](#deployment-architecture)

---

## 🛠️ Technology Stack

### Core Framework
- **React 18.2+** - UI library dengan modern hooks & concurrent features
- **TypeScript 5.3+** - Type safety & better developer experience
- **Vite 5.0+** - Lightning-fast build tool & dev server

### UI & Styling
- **TailwindCSS 3.4+** - Utility-first CSS framework
- **Shadcn/ui** - High-quality, accessible component library
- **Lucide React** - Beautiful icon library
- **Recharts 2.10+** - Composable charting library

### State Management
- **Zustand 4.4+** - Lightweight state management
- **React Query (TanStack Query) 5.17+** - Server state management
- **React Hook Form 7.49+** - Performant form management

### Data Fetching
- **Axios 1.6+** - HTTP client
- **React Query** - Caching, synchronization, and updates

### Routing
- **React Router v6.20+** - Client-side routing

### Form Validation
- **Zod 3.22+** - Schema validation
- **@hookform/resolvers** - Integration with React Hook Form

### Utilities
- **date-fns 3.0+** - Date manipulation
- **clsx** - Conditional className utility
- **tailwind-merge** - Merge Tailwind classes

### Development Tools
- **ESLint** - Code linting
- **Prettier** - Code formatting
- **TypeScript** - Type checking

### Testing (Optional Phase 7)
- **Jest** - Unit testing framework
- **React Testing Library** - Component testing
- **Playwright** - E2E testing

---

## 🏗️ Architecture Overview

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────┐
│                     Browser (Client)                     │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  ┌────────────────────────────────────────────────┐    │
│  │         React Application (SPA)                 │    │
│  │                                                  │    │
│  │  ┌──────────────┐  ┌──────────────────────┐   │    │
│  │  │   UI Layer   │  │   Component Library   │   │    │
│  │  │  (Pages)     │  │   (Shadcn/ui)        │   │    │
│  │  └──────────────┘  └──────────────────────┘   │    │
│  │                                                  │    │
│  │  ┌──────────────────────────────────────────┐  │    │
│  │  │      State Management Layer              │  │    │
│  │  │  - Zustand (Client State)                │  │    │
│  │  │  - React Query (Server State)            │  │    │
│  │  └──────────────────────────────────────────┘  │    │
│  │                                                  │    │
│  │  ┌──────────────────────────────────────────┐  │    │
│  │  │      API Client Layer (Axios)            │  │    │
│  │  │  - Request/Response Interceptors         │  │    │
│  │  │  - Error Handling                        │  │    │
│  │  │  - Retry Logic                           │  │    │
│  │  └──────────────────────────────────────────┘  │    │
│  │                                                  │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
└─────────────────────────────────────────────────────────┘
                          │
                          │ HTTPS
                          ▼
┌─────────────────────────────────────────────────────────┐
│                    VPS Server (Nginx)                    │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  ┌────────────────────────────────────────────────┐    │
│  │  Nginx Reverse Proxy                           │    │
│  │  - Serve static files (Frontend)               │    │
│  │  - Proxy /api/* to Backend                     │    │
│  │  - SSL/TLS termination                         │    │
│  │  - Gzip compression                            │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
└─────────────────────────────────────────────────────────┘
                          │
                          │ HTTP (internal)
                          ▼
┌─────────────────────────────────────────────────────────┐
│              Go Backend API (Port 8081)                  │
│            (ZTE C320 OLT Management API)                 │
└─────────────────────────────────────────────────────────┘
```

### Application Flow

```
User Interaction
      │
      ▼
React Component
      │
      ├─→ Local State (useState, useReducer)
      │
      ├─→ Global State (Zustand)
      │
      └─→ Server State (React Query)
              │
              ▼
          API Client (Axios)
              │
              ├─→ Request Interceptor
              │   - Add auth token
              │   - Add headers
              │
              ▼
          HTTP Request
              │
              ▼
          Backend API
              │
              ▼
          HTTP Response
              │
              ├─→ Response Interceptor
              │   - Handle errors
              │   - Transform data
              │
              ▼
          React Query Cache
              │
              ▼
          Component Re-render
              │
              ▼
          Updated UI
```

---

## 📁 Folder Structure

```
frontend/
├── public/                         # Static assets
│   ├── favicon.ico
│   ├── logo.svg
│   └── robots.txt
│
├── src/
│   ├── api/                        # API layer
│   │   ├── client.ts              # Axios instance & config
│   │   ├── queryClient.ts         # React Query config
│   │   ├── endpoints/             # API endpoint functions
│   │   │   ├── auth.ts           # Authentication endpoints
│   │   │   ├── onu.ts            # ONU monitoring endpoints
│   │   │   ├── provision.ts      # Provisioning endpoints
│   │   │   ├── vlan.ts           # VLAN management endpoints
│   │   │   ├── traffic.ts        # Traffic profile endpoints
│   │   │   ├── management.ts     # ONU management endpoints
│   │   │   ├── backup.ts         # Config backup endpoints
│   │   │   └── monitoring.ts     # Real-time monitoring
│   │   └── types/                # API type definitions
│   │       ├── api-response.ts   # Generic API response types
│   │       ├── onu.ts            # ONU related types
│   │       ├── vlan.ts           # VLAN related types
│   │       └── index.ts          # Export all types
│   │
│   ├── components/                # React components
│   │   ├── ui/                   # Shadcn/ui base components
│   │   │   ├── button.tsx
│   │   │   ├── card.tsx
│   │   │   ├── dialog.tsx
│   │   │   ├── dropdown-menu.tsx
│   │   │   ├── input.tsx
│   │   │   ├── select.tsx
│   │   │   ├── table.tsx
│   │   │   ├── toast.tsx
│   │   │   └── ... (other Shadcn components)
│   │   │
│   │   ├── layout/               # Layout components
│   │   │   ├── AppLayout.tsx     # Main app layout
│   │   │   ├── Sidebar.tsx       # Sidebar navigation
│   │   │   ├── Header.tsx        # Top header
│   │   │   ├── Footer.tsx        # Footer
│   │   │   └── MobileMenu.tsx    # Mobile navigation
│   │   │
│   │   ├── common/               # Common reusable components
│   │   │   ├── DataTable.tsx     # Advanced table with features
│   │   │   ├── PageHeader.tsx    # Page header with breadcrumb
│   │   │   ├── StatusBadge.tsx   # Status indicator badge
│   │   │   ├── LoadingSpinner.tsx
│   │   │   ├── EmptyState.tsx
│   │   │   ├── ErrorBoundary.tsx
│   │   │   └── ThemeToggle.tsx   # Theme switcher
│   │   │
│   │   ├── charts/               # Chart components
│   │   │   ├── SignalQualityChart.tsx
│   │   │   ├── StatusDistributionChart.tsx
│   │   │   ├── PonUtilizationChart.tsx
│   │   │   └── TrendChart.tsx
│   │   │
│   │   ├── forms/                # Form components
│   │   │   ├── OnuProvisionForm.tsx
│   │   │   ├── VlanConfigForm.tsx
│   │   │   ├── TrafficProfileForm.tsx
│   │   │   └── FormWizard.tsx
│   │   │
│   │   └── features/             # Feature-specific components
│   │       ├── onu/
│   │       │   ├── OnuCard.tsx
│   │       │   ├── OnuDetailModal.tsx
│   │       │   └── OnuActionMenu.tsx
│   │       ├── monitoring/
│   │       │   ├── RealTimeMonitor.tsx
│   │       │   └── OpticalPowerDisplay.tsx
│   │       └── provisioning/
│   │           ├── AutoDiscoveryList.tsx
│   │           └── ProvisionWizard.tsx
│   │
│   ├── pages/                    # Page components
│   │   ├── Dashboard/
│   │   │   ├── index.tsx         # Dashboard main page
│   │   │   ├── OverviewCards.tsx
│   │   │   ├── RecentActivity.tsx
│   │   │   └── QuickActions.tsx
│   │   │
│   │   ├── Monitoring/
│   │   │   ├── index.tsx         # Monitoring page
│   │   │   ├── OnuList.tsx
│   │   │   ├── FilterPanel.tsx
│   │   │   └── OnuDetailView.tsx
│   │   │
│   │   ├── Provisioning/
│   │   │   ├── AutoDiscovery/
│   │   │   │   └── index.tsx
│   │   │   ├── ManualProvision/
│   │   │   │   └── index.tsx
│   │   │   └── BatchOperations/
│   │   │       └── index.tsx
│   │   │
│   │   ├── VlanManagement/
│   │   │   ├── index.tsx
│   │   │   ├── VlanList.tsx
│   │   │   └── ServicePortConfig.tsx
│   │   │
│   │   ├── TrafficManagement/
│   │   │   ├── index.tsx
│   │   │   ├── DbaProfiles.tsx
│   │   │   └── TcontConfig.tsx
│   │   │
│   │   ├── OnuManagement/
│   │   │   ├── index.tsx
│   │   │   └── ConfigBackup.tsx
│   │   │
│   │   ├── Settings/
│   │   │   ├── index.tsx
│   │   │   ├── SystemSettings.tsx
│   │   │   ├── UserPreferences.tsx
│   │   │   └── UserManagement.tsx
│   │   │
│   │   ├── Auth/
│   │   │   ├── Login.tsx
│   │   │   └── ForgotPassword.tsx
│   │   │
│   │   └── ErrorPages/
│   │       ├── NotFound.tsx
│   │       └── ServerError.tsx
│   │
│   ├── hooks/                    # Custom React hooks
│   │   ├── useApi.ts            # Generic API hook
│   │   ├── useAuth.ts           # Authentication hook
│   │   ├── useTheme.ts          # Theme management hook
│   │   ├── useOnuData.ts        # ONU data fetching hook
│   │   ├── useDebounce.ts       # Debounce hook
│   │   ├── useLocalStorage.ts   # Local storage hook
│   │   └── useWebSocket.ts      # WebSocket hook (future)
│   │
│   ├── store/                    # Zustand stores
│   │   ├── authStore.ts         # Auth state
│   │   ├── themeStore.ts        # Theme state
│   │   ├── sidebarStore.ts      # Sidebar state
│   │   ├── onuStore.ts          # ONU filter state
│   │   └── index.ts             # Export all stores
│   │
│   ├── utils/                    # Utility functions
│   │   ├── formatters.ts        # Data formatters
│   │   │   ├── formatDate()
│   │   │   ├── formatPower()
│   │   │   └── formatBytes()
│   │   ├── validators.ts        # Validation functions
│   │   ├── constants.ts         # App constants
│   │   ├── helpers.ts           # Helper functions
│   │   └── cn.ts                # className utility
│   │
│   ├── types/                    # TypeScript type definitions
│   │   ├── index.ts             # Global types
│   │   ├── components.ts        # Component prop types
│   │   └── models.ts            # Data model types
│   │
│   ├── styles/                   # Global styles
│   │   └── globals.css          # Global CSS + Tailwind
│   │
│   ├── config/                   # Configuration files
│   │   ├── env.ts               # Environment variables
│   │   └── routes.ts            # Route constants
│   │
│   ├── lib/                      # Third-party library configs
│   │   └── utils.ts             # Utility exports
│   │
│   ├── App.tsx                   # Root App component
│   ├── main.tsx                  # Entry point
│   └── router.tsx                # Router configuration
│
├── .env.example                  # Environment template
├── .env.development              # Development env
├── .env.production               # Production env
├── .eslintrc.json               # ESLint config
├── .prettierrc                  # Prettier config
├── tsconfig.json                # TypeScript config
├── vite.config.ts               # Vite config
├── tailwind.config.js           # Tailwind config
├── postcss.config.js            # PostCSS config
├── package.json                 # Dependencies
└── README.md                    # Project readme
```

---

## 🧩 Component Architecture

### Component Hierarchy

```
App
├── Router
    ├── AppLayout
    │   ├── Sidebar
    │   │   ├── Logo
    │   │   ├── Navigation Menu
    │   │   │   ├── MenuItem (Dashboard)
    │   │   │   ├── MenuItem (Monitoring)
    │   │   │   ├── MenuItem (Provisioning) [nested]
    │   │   │   │   ├── SubMenuItem (Auto Discovery)
    │   │   │   │   ├── SubMenuItem (Manual)
    │   │   │   │   └── SubMenuItem (Batch)
    │   │   │   ├── MenuItem (VLAN)
    │   │   │   ├── MenuItem (Traffic)
    │   │   │   ├── MenuItem (Management)
    │   │   │   └── MenuItem (Settings)
    │   │   └── ThemeToggle
    │   │
    │   ├── Header
    │   │   ├── Breadcrumb
    │   │   ├── SearchBar
    │   │   ├── NotificationBell
    │   │   └── UserMenu
    │   │
    │   ├── Main Content Area
    │   │   └── [Dynamic Page Component]
    │   │       ├── PageHeader
    │   │       └── Page Content
    │   │
    │   └── Footer
    │
    └── Toast Container
```

### Component Categories

#### 1. **Presentational Components**
Pure components yang hanya menerima props dan render UI.

```typescript
// Example: StatusBadge.tsx
interface StatusBadgeProps {
  status: 'online' | 'offline';
  label?: string;
}

export const StatusBadge: React.FC<StatusBadgeProps> = ({ status, label }) => {
  return (
    <Badge variant={status === 'online' ? 'success' : 'destructive'}>
      {label || status}
    </Badge>
  );
};
```

#### 2. **Container Components**
Components yang mengelola state dan logic.

```typescript
// Example: OnuList.tsx
export const OnuList: React.FC = () => {
  const { data, isLoading } = useOnuData();
  const [filters, setFilters] = useState({});
  
  // Logic here...
  
  return (
    <DataTable
      data={data}
      columns={columns}
      isLoading={isLoading}
    />
  );
};
```

#### 3. **Layout Components**
Components untuk struktur halaman.

```typescript
// Example: AppLayout.tsx
export const AppLayout: React.FC<{ children: ReactNode }> = ({ children }) => {
  return (
    <div className="flex h-screen">
      <Sidebar />
      <div className="flex flex-col flex-1">
        <Header />
        <main className="flex-1 overflow-auto p-6">
          {children}
        </main>
        <Footer />
      </div>
    </div>
  );
};
```

---

## 🗄️ State Management

### State Management Strategy

**3-Layer State Management:**

#### 1. Local Component State (useState, useReducer)
Untuk state yang hanya digunakan dalam satu component.

```typescript
const [isOpen, setIsOpen] = useState(false);
const [formData, setFormData] = useState({ name: '', description: '' });
```

**Use Cases:**
- Toggle states (modal, dropdown)
- Form input values
- UI-only states

#### 2. Global Client State (Zustand)
Untuk state yang dibagikan antar components.

```typescript
// store/authStore.ts
import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface AuthState {
  user: User | null;
  token: string | null;
  login: (username: string, password: string) => Promise<void>;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      login: async (username, password) => {
        // Login logic
        const response = await api.login(username, password);
        set({ user: response.user, token: response.token });
      },
      logout: () => {
        set({ user: null, token: null });
      },
    }),
    {
      name: 'auth-storage',
    }
  )
);
```

**Use Cases:**
- Authentication state
- Theme preferences
- Sidebar collapsed/expanded
- Selected filters
- User preferences

#### 3. Server State (React Query)
Untuk data dari API.

```typescript
// hooks/useOnuData.ts
import { useQuery } from '@tanstack/react-query';
import { onuApi } from '@/api/endpoints/onu';

export const useOnuData = (board: number, pon: number) => {
  return useQuery({
    queryKey: ['onu-list', board, pon],
    queryFn: () => onuApi.getOnuList(board, pon),
    refetchInterval: 30000, // Auto-refresh every 30s
    staleTime: 10000,
  });
};
```

**Use Cases:**
- API data fetching
- Caching
- Background refetching
- Optimistic updates
- Mutations

### Zustand Stores

```typescript
// store/themeStore.ts
export const useThemeStore = create<ThemeState>((set) => ({
  theme: 'light',
  setTheme: (theme) => set({ theme }),
  toggleTheme: () => set((state) => ({
    theme: state.theme === 'light' ? 'dark' : 'light'
  })),
}));

// store/sidebarStore.ts
export const useSidebarStore = create<SidebarState>((set) => ({
  isCollapsed: false,
  toggle: () => set((state) => ({ isCollapsed: !state.isCollapsed })),
}));

// store/onuStore.ts
export const useOnuStore = create<OnuState>((set) => ({
  filters: {
    board: null,
    pon: null,
    status: null,
  },
  setFilters: (filters) => set({ filters }),
  resetFilters: () => set({
    filters: { board: null, pon: null, status: null }
  }),
}));
```

---

## 🔌 API Integration

### API Client Setup

```typescript
// api/client.ts
import axios from 'axios';
import { useAuthStore } from '@/store/authStore';

export const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor
apiClient.interceptors.request.use(
  (config) => {
    const token = useAuthStore.getState().token;
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor
apiClient.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      useAuthStore.getState().logout();
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);
```

### React Query Configuration

```typescript
// api/queryClient.ts
import { QueryClient } from '@tanstack/react-query';

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 3,
      retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30000),
      staleTime: 5 * 60 * 1000, // 5 minutes
      cacheTime: 10 * 60 * 1000, // 10 minutes
      refetchOnWindowFocus: false,
    },
    mutations: {
      retry: 1,
    },
  },
});
```

### API Endpoint Example

```typescript
// api/endpoints/onu.ts
import { apiClient } from '../client';
import type { OnuListResponse, OnuDetailResponse } from '../types/onu';

export const onuApi = {
  // Get ONU list
  getOnuList: async (board: number, pon: number): Promise<OnuListResponse> => {
    return apiClient.get(`/board/${board}/pon/${pon}/`);
  },

  // Get ONU details
  getOnuDetail: async (
    board: number,
    pon: number,
    onuId: number
  ): Promise<OnuDetailResponse> => {
    return apiClient.get(`/board/${board}/pon/${pon}/onu/${onuId}`);
  },

  // Reboot ONU
  rebootOnu: async (payload: RebootPayload) => {
    return apiClient.post('/management/reboot', payload);
  },
};
```

### Custom Hook Pattern

```typescript
// hooks/useOnuDetail.ts
export const useOnuDetail = (board: number, pon: number, onuId: number) => {
  return useQuery({
    queryKey: ['onu-detail', board, pon, onuId],
    queryFn: () => onuApi.getOnuDetail(board, pon, onuId),
    enabled: !!board && !!pon && !!onuId,
  });
};

// Usage in component
const OnuDetailModal = ({ board, pon, onuId }) => {
  const { data, isLoading, error } = useOnuDetail(board, pon, onuId);
  
  if (isLoading) return <Spinner />;
  if (error) return <ErrorMessage />;
  
  return <div>{/* Render data */}</div>;
};
```

---

## 🛣️ Routing

### Route Structure

```typescript
// router.tsx
import { createBrowserRouter } from 'react-router-dom';

export const router = createBrowserRouter([
  {
    path: '/',
    element: <AppLayout />,
    children: [
      {
        index: true,
        element: <Dashboard />,
      },
      {
        path: 'monitoring',
        element: <Monitoring />,
      },
      {
        path: 'provision',
        children: [
          {
            path: 'auto',
            element: <AutoDiscovery />,
          },
          {
            path: 'manual',
            element: <ManualProvision />,
          },
          {
            path: 'batch',
            element: <BatchOperations />,
          },
        ],
      },
      {
        path: 'vlan',
        element: <VlanManagement />,
      },
      {
        path: 'traffic',
        element: <TrafficManagement />,
      },
      {
        path: 'onu',
        element: <OnuManagement />,
      },
      {
        path: 'backup',
        element: <ConfigBackup />,
      },
      {
        path: 'settings',
        element: <Settings />,
      },
    ],
  },
  {
    path: '/login',
    element: <Login />,
  },
  {
    path: '*',
    element: <NotFound />,
  },
]);
```

### Protected Routes

```typescript
// components/common/ProtectedRoute.tsx
export const ProtectedRoute = ({ children }: { children: ReactNode }) => {
  const { token } = useAuthStore();
  
  if (!token) {
    return <Navigate to="/login" replace />;
  }
  
  return <>{children}</>;
};
```

---

## 🎨 Styling & Theming

### TailwindCSS Configuration

```javascript
// tailwind.config.js
module.exports = {
  darkMode: 'class',
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        border: 'hsl(var(--border))',
        input: 'hsl(var(--input))',
        ring: 'hsl(var(--ring))',
        background: 'hsl(var(--background))',
        foreground: 'hsl(var(--foreground))',
        primary: {
          DEFAULT: 'hsl(var(--primary))',
          foreground: 'hsl(var(--primary-foreground))',
        },
        // ... more colors
      },
    },
  },
  plugins: [require('tailwindcss-animate')],
};
```

### CSS Variables (Theme)

```css
/* styles/globals.css */
@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  :root {
    --background: 0 0% 100%;
    --foreground: 222.2 84% 4.9%;
    --primary: 221.2 83.2% 53.3%;
    --primary-foreground: 210 40% 98%;
    /* ... more variables */
  }

  .dark {
    --background: 222.2 84% 4.9%;
    --foreground: 210 40% 98%;
    --primary: 217.2 91.2% 59.8%;
    --primary-foreground: 222.2 47.4% 11.2%;
    /* ... more variables */
  }
}
```

### Theme Provider

```typescript
// components/common/ThemeProvider.tsx
export const ThemeProvider = ({ children }: { children: ReactNode }) => {
  const { theme } = useThemeStore();

  useEffect(() => {
    const root = window.document.documentElement;
    root.classList.remove('light', 'dark');
    root.classList.add(theme);
  }, [theme]);

  return <>{children}</>;
};
```

---

## ⚡ Performance Optimization

### Code Splitting

```typescript
// Lazy load pages
const Dashboard = lazy(() => import('@/pages/Dashboard'));
const Monitoring = lazy(() => import('@/pages/Monitoring'));

// Wrap with Suspense
<Suspense fallback={<LoadingSpinner />}>
  <Dashboard />
</Suspense>
```

### Memoization

```typescript
// Memoize expensive computations
const filteredOnus = useMemo(() => {
  return onus.filter((onu) => onu.status === selectedStatus);
}, [onus, selectedStatus]);

// Memoize callbacks
const handleClick = useCallback(() => {
  // handler logic
}, [dependencies]);
```

### Virtual Scrolling (Large Lists)

```typescript
// For large ONU lists
import { useVirtualizer } from '@tanstack/react-virtual';
```

---

## 🔒 Security Considerations

### 1. Authentication & Authorization
- JWT token stored in Zustand (with persist)
- Token in Authorization header
- Auto logout on 401 response

### 2. XSS Protection
- React escapes by default
- Use `dangerouslySetInnerHTML` with caution
- Sanitize user inputs

### 3. HTTPS
- Force HTTPS in production
- Secure cookies

### 4. Environment Variables
- Never commit `.env` files
- Use `.env.example` as template

---

## 🚀 Deployment Architecture

### Build Process

```bash
# Build production bundle
npm run build

# Output: dist/
dist/
├── index.html
├── assets/
│   ├── index-[hash].js
│   ├── index-[hash].css
│   └── vendor-[hash].js
└── ...
```

### VPS Deployment Structure

```
/var/www/olt-dashboard/
├── frontend/              # Frontend static files
│   ├── index.html
│   └── assets/
└── nginx.conf            # Nginx configuration
```

### Nginx Configuration

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    # Redirect to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;
    
    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;
    
    # Frontend static files
    root /var/www/olt-dashboard/frontend;
    index index.html;
    
    # Gzip compression
    gzip on;
    gzip_types text/plain text/css application/json application/javascript;
    
    # Frontend routes (SPA)
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    # Proxy API requests to backend
    location /api/ {
        proxy_pass http://localhost:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    # Cache static assets
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
```

---

## 📊 Performance Metrics

### Target Metrics
- First Contentful Paint (FCP): < 1.5s
- Time to Interactive (TTI): < 3.0s
- Largest Contentful Paint (LCP): < 2.5s
- Cumulative Layout Shift (CLS): < 0.1
- Bundle Size: < 500KB (gzipped)
- Lighthouse Score: > 90

---

## 📝 Summary

This architecture provides:
- ✅ **Scalable** - Modular structure, easy to extend
- ✅ **Maintainable** - Clear separation of concerns
- ✅ **Performant** - Optimized bundle, code splitting
- ✅ **Type-safe** - Full TypeScript coverage
- ✅ **Modern** - Latest React patterns & best practices
- ✅ **Responsive** - Mobile-first design
- ✅ **Themeable** - Light/Dark mode support
- ✅ **Accessible** - WCAG 2.1 compliant components

---

**Document Version:** 1.0  
**Last Updated:** January 12, 2026  
**Next Review:** After Phase 1 completion

