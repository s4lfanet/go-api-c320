# Frontend Development Workload
## ZTE C320 OLT Management Dashboard

**Project:** Web-based Dashboard for ZTE C320 OLT  
**Document Version:** 1.0  
**Last Updated:** January 12, 2026

---

## 📋 Table of Contents

1. [Project Setup Tasks](#project-setup-tasks)
2. [Core Features Development](#core-features-development)
3. [UI/UX Components](#uiux-components)
4. [API Integration](#api-integration)
5. [Testing & Quality Assurance](#testing--quality-assurance)
6. [Deployment & DevOps](#deployment--devops)
7. [Time Estimation](#time-estimation)

---

## 🔧 Project Setup Tasks

### 1.1 Initial Configuration (8 hours)
- [x] Create project structure
- [ ] Initialize Vite + React + TypeScript project
- [ ] Configure TailwindCSS
- [ ] Setup Shadcn/ui components
- [ ] Configure ESLint + Prettier
- [ ] Setup Git (local only, no GitHub)
- [ ] Create folder structure
- [ ] Setup environment variables (.env files)

**Files to Create:**
```
frontend/
├── .env.example
├── .env.development
├── .env.production
├── .eslintrc.json
├── .prettierrc
├── tsconfig.json
├── vite.config.ts
├── tailwind.config.js
├── postcss.config.js
└── index.html
```

### 1.2 Dependencies Installation (2 hours)
**Core Dependencies:**
```json
{
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-router-dom": "^6.20.0",
    "axios": "^1.6.0",
    "@tanstack/react-query": "^5.17.0",
    "zustand": "^4.4.7",
    "react-hook-form": "^7.49.0",
    "zod": "^3.22.4",
    "@hookform/resolvers": "^3.3.3",
    "lucide-react": "^0.300.0",
    "recharts": "^2.10.0",
    "date-fns": "^3.0.0",
    "clsx": "^2.0.0",
    "tailwind-merge": "^2.2.0"
  },
  "devDependencies": {
    "@types/react": "^18.2.0",
    "@types/react-dom": "^18.2.0",
    "@vitejs/plugin-react": "^4.2.0",
    "typescript": "^5.3.0",
    "vite": "^5.0.0",
    "tailwindcss": "^3.4.0",
    "autoprefixer": "^10.4.16",
    "postcss": "^8.4.32",
    "eslint": "^8.56.0",
    "prettier": "^3.1.0"
  }
}
```

### 1.3 Project Structure Setup (4 hours)
```
frontend/src/
├── api/                    # API client & endpoints
│   ├── client.ts
│   ├── endpoints/
│   │   ├── onu.ts
│   │   ├── vlan.ts
│   │   ├── traffic.ts
│   │   ├── provision.ts
│   │   └── monitoring.ts
│   └── types/
│       └── api-types.ts
├── components/             # Reusable components
│   ├── ui/                # Shadcn/ui components
│   ├── layout/
│   │   ├── Sidebar.tsx
│   │   ├── Header.tsx
│   │   ├── Footer.tsx
│   │   └── Layout.tsx
│   ├── common/
│   │   ├── Card.tsx
│   │   ├── Table.tsx
│   │   ├── Modal.tsx
│   │   └── Toast.tsx
│   └── charts/
│       ├── SignalChart.tsx
│       ├── StatusChart.tsx
│       └── PonUtilization.tsx
├── pages/                 # Page components
│   ├── Dashboard/
│   ├── Monitoring/
│   ├── Provisioning/
│   ├── VlanManagement/
│   ├── OnuManagement/
│   └── Settings/
├── hooks/                 # Custom React hooks
│   ├── useApi.ts
│   ├── useAuth.ts
│   ├── useTheme.ts
│   └── useWebSocket.ts
├── store/                 # Zustand stores
│   ├── authStore.ts
│   ├── themeStore.ts
│   └── onuStore.ts
├── utils/                 # Utility functions
│   ├── formatters.ts
│   ├── validators.ts
│   └── constants.ts
├── styles/               # Global styles
│   └── globals.css
├── types/                # TypeScript types
│   └── index.ts
├── App.tsx
├── main.tsx
└── router.tsx
```

---

## 🎨 UI/UX Components

### 2.1 Layout Components (16 hours)

#### Sidebar Menu (6 hours)
- [ ] Responsive sidebar (collapsible on mobile)
- [ ] Menu items dengan icons
- [ ] Active state highlighting
- [ ] Nested menu support
- [ ] Light/Dark theme toggle
- [ ] Logo & branding area

**Menu Structure:**
```typescript
const menuItems = [
  { icon: LayoutDashboard, label: 'Dashboard', path: '/' },
  { icon: Activity, label: 'Monitoring', path: '/monitoring' },
  { 
    icon: Network, 
    label: 'Provisioning', 
    children: [
      { label: 'Auto Discovery', path: '/provision/auto' },
      { label: 'Manual Provision', path: '/provision/manual' },
      { label: 'Batch Operations', path: '/provision/batch' }
    ]
  },
  { icon: Cable, label: 'VLAN Management', path: '/vlan' },
  { icon: Gauge, label: 'Traffic Profiles', path: '/traffic' },
  { icon: Settings2, label: 'ONU Management', path: '/onu' },
  { icon: Database, label: 'Config Backup', path: '/backup' },
  { icon: Settings, label: 'Settings', path: '/settings' }
];
```

#### Header Component (4 hours)
- [ ] Breadcrumb navigation
- [ ] User profile dropdown
- [ ] Notification bell
- [ ] Search bar
- [ ] Theme switcher
- [ ] Mobile menu toggle

#### Footer Component (2 hours)
- [ ] Copyright info
- [ ] Version number
- [ ] Quick links

#### Main Layout (4 hours)
- [ ] Responsive grid system
- [ ] Content area with padding
- [ ] Mobile-first approach
- [ ] Smooth transitions

### 2.2 Common Components (24 hours)

#### Data Table Component (8 hours)
- [ ] Sortable columns
- [ ] Pagination
- [ ] Row selection (single/multiple)
- [ ] Search/filter
- [ ] Custom cell renderers
- [ ] Loading states
- [ ] Empty states
- [ ] Export to CSV/Excel
- [ ] Responsive (card view on mobile)

#### Modal/Dialog Component (4 hours)
- [ ] Customizable sizes
- [ ] Backdrop click handling
- [ ] ESC key close
- [ ] Confirm dialogs
- [ ] Form modals

#### Toast/Notification Component (3 hours)
- [ ] Success/Error/Warning/Info types
- [ ] Auto-dismiss
- [ ] Action buttons
- [ ] Stack multiple toasts
- [ ] Position variants

#### Form Components (6 hours)
- [ ] Input field dengan validation
- [ ] Select/Dropdown
- [ ] Checkbox/Radio
- [ ] Date picker
- [ ] Multi-select
- [ ] Form wizard component

#### Loading States (3 hours)
- [ ] Skeleton loaders
- [ ] Spinner variants
- [ ] Progress bars
- [ ] Shimmer effects

### 2.3 Chart Components (12 hours)

#### Signal Quality Chart (4 hours)
- [ ] Real-time RX/TX power visualization
- [ ] Line chart dengan multiple series
- [ ] Threshold indicators
- [ ] Zoom & pan support
- [ ] Tooltip dengan detailed info

#### Status Distribution (3 hours)
- [ ] Pie/Donut chart
- [ ] Online/Offline statistics
- [ ] Color-coded segments
- [ ] Legend & labels

#### PON Utilization Chart (3 hours)
- [ ] Bar chart per PON port
- [ ] Color-coded by utilization %
- [ ] Sortable
- [ ] Click to drill-down

#### Historical Trends (2 hours)
- [ ] Area chart
- [ ] Date range selector
- [ ] Export chart as image

---

## 🔌 API Integration

### 3.1 API Client Setup (8 hours)

#### Axios Configuration (3 hours)
- [ ] Base URL configuration
- [ ] Request interceptors (auth token)
- [ ] Response interceptors (error handling)
- [ ] Retry logic
- [ ] Timeout configuration
- [ ] Request cancellation

```typescript
// api/client.ts
const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});
```

#### React Query Setup (3 hours)
- [ ] Query client configuration
- [ ] Default query options
- [ ] Cache configuration
- [ ] Devtools integration
- [ ] Mutation handling

#### Error Handling (2 hours)
- [ ] Custom error class
- [ ] Error boundary component
- [ ] User-friendly error messages
- [ ] Retry strategies

### 3.2 API Endpoints Implementation (40 hours)

#### ONU Monitoring APIs (8 hours)
- [ ] `GET /board/{board_id}/pon/{pon_id}/` - List ONUs
- [ ] `GET /board/{board_id}/pon/{pon_id}/onu/{onu_id}` - ONU details
- [ ] `GET /board/{board_id}/pon/{pon_id}/info` - PON info
- [ ] Real-time monitoring endpoint
- [ ] TypeScript types untuk responses

#### Provisioning APIs (8 hours)
- [ ] `GET /provision/unconfigured` - Auto-discovery
- [ ] `POST /provision/authorize` - Authorize ONU
- [ ] `POST /provision/create` - Manual provision
- [ ] `POST /provision/batch` - Batch provision
- [ ] Form validation schemas

#### VLAN Management APIs (6 hours)
- [ ] `POST /vlan/create` - Create VLAN
- [ ] `DELETE /vlan/delete` - Delete VLAN
- [ ] VLAN assignment endpoints
- [ ] Service-port configuration

#### Traffic Profile APIs (6 hours)
- [ ] DBA profile endpoints
- [ ] T-CONT configuration
- [ ] GEM port settings
- [ ] Profile templates

#### ONU Management APIs (8 hours)
- [ ] `POST /management/reboot` - Reboot
- [ ] `POST /management/block` - Block ONU
- [ ] `POST /management/unblock` - Unblock ONU
- [ ] `DELETE /management/delete` - Delete ONU
- [ ] `PUT /management/description` - Update description

#### Config Backup APIs (4 hours)
- [ ] `POST /config/backup` - Backup config
- [ ] `POST /config/restore` - Restore config
- [ ] Download backup file
- [ ] Upload backup file

---

## 📄 Page Development

### 4.1 Dashboard Page (16 hours)

#### Overview Cards (4 hours)
- [ ] Total ONUs card
- [ ] Online/Offline status
- [ ] Signal quality summary
- [ ] Active alarms count

#### Charts Section (6 hours)
- [ ] Status distribution chart
- [ ] PON utilization chart
- [ ] Signal trends chart

#### Recent Activity Table (4 hours)
- [ ] Recent ONU registrations
- [ ] Recent alarms
- [ ] Auto-refresh (30s interval)

#### Quick Actions (2 hours)
- [ ] Quick provision button
- [ ] Quick search
- [ ] Export reports

### 4.2 Monitoring Page (20 hours)

#### Filters Section (4 hours)
- [ ] Board selector
- [ ] PON selector
- [ ] Status filter
- [ ] Signal quality filter
- [ ] Search by name/SN

#### ONU List Table (10 hours)
- [ ] Sortable columns
- [ ] Real-time status updates
- [ ] Signal quality indicators
- [ ] Action buttons (detail/reboot/delete)
- [ ] Bulk selection
- [ ] Export data
- [ ] Pagination

#### ONU Detail Modal (6 hours)
- [ ] Complete ONU information
- [ ] Real-time optical power
- [ ] Uptime & last seen
- [ ] Configuration summary
- [ ] Quick actions

### 4.3 Provisioning Pages (24 hours)

#### Auto-Discovery Page (8 hours)
- [ ] Scan for unconfigured ONUs
- [ ] ONU list dengan serial numbers
- [ ] Quick authorize button
- [ ] Batch authorization
- [ ] Filter & search

#### Manual Provision Page (10 hours)
- [ ] Multi-step wizard
  - Step 1: Board/PON selection
  - Step 2: ONU info (SN, Type, Name)
  - Step 3: Traffic profile
  - Step 4: VLAN configuration
  - Step 5: Review & submit
- [ ] Form validation
- [ ] Progress indicator
- [ ] Save as template

#### Batch Operations Page (6 hours)
- [ ] Upload CSV file
- [ ] Preview import data
- [ ] Validation errors display
- [ ] Progress tracking
- [ ] Success/failure report

### 4.4 VLAN Management Page (16 hours)

#### VLAN List (6 hours)
- [ ] VLAN table
- [ ] Add/Edit/Delete VLAN
- [ ] Service-port list per VLAN
- [ ] Search & filter

#### Service-Port Configuration (8 hours)
- [ ] Create service-port form
- [ ] ONU selection
- [ ] VLAN assignment
- [ ] User/Inner VLAN
- [ ] Priority & bandwidth

#### Bulk VLAN Operations (2 hours)
- [ ] Multi-select ONUs
- [ ] Assign VLAN to multiple ONUs
- [ ] Progress tracking

### 4.5 Traffic Management Page (12 hours)

#### DBA Profiles (4 hours)
- [ ] Profile list
- [ ] Create/Edit profile
- [ ] Bandwidth settings
- [ ] Apply to ONUs

#### T-CONT Configuration (4 hours)
- [ ] T-CONT list per ONU
- [ ] Configure T-CONT
- [ ] Traffic type selection

#### GEM Port Settings (4 hours)
- [ ] GEM port list
- [ ] Configuration form
- [ ] Mapping to services

### 4.6 ONU Management Page (16 hours)

#### ONU Operations (8 hours)
- [ ] ONU list dengan action buttons
- [ ] Reboot ONU (single/bulk)
- [ ] Block/Unblock
- [ ] Delete confirmation
- [ ] Update description

#### Configuration Backup (8 hours)
- [ ] Backup configuration
- [ ] Restore configuration
- [ ] Backup history
- [ ] Download/Upload files
- [ ] Schedule automatic backups

### 4.7 Settings Page (12 hours)

#### System Settings (4 hours)
- [ ] OLT connection settings
- [ ] API endpoint configuration
- [ ] Timeout settings
- [ ] Auto-refresh intervals

#### User Preferences (4 hours)
- [ ] Theme selection
- [ ] Language selection (ID/EN)
- [ ] Date format
- [ ] Timezone

#### User Management (4 hours)
- [ ] User list
- [ ] Add/Edit/Delete users
- [ ] Role assignment
- [ ] Password change

---

## 🎨 Theme Implementation

### 5.1 Light/Dark Theme (12 hours)

#### Theme Provider (4 hours)
- [ ] Create ThemeContext
- [ ] Zustand theme store
- [ ] Persist theme preference (localStorage)
- [ ] System preference detection

#### Theme Styles (6 hours)
- [ ] CSS variables untuk colors
- [ ] TailwindCSS dark mode config
- [ ] Component variants
- [ ] Smooth transitions

#### Theme Toggle (2 hours)
- [ ] Toggle button component
- [ ] Icon animation
- [ ] Keyboard shortcut (Ctrl+Shift+T)

---

## 🧪 Testing & Quality Assurance

### 6.1 Unit Testing (20 hours)
- [ ] Setup Jest + React Testing Library
- [ ] Component tests (80% coverage)
- [ ] Hook tests
- [ ] Utility function tests
- [ ] Store tests

### 6.2 Integration Testing (16 hours)
- [ ] API integration tests
- [ ] Form submission tests
- [ ] Navigation tests
- [ ] Authentication flow tests

### 6.3 E2E Testing (12 hours)
- [ ] Setup Playwright
- [ ] Critical user flows
- [ ] Provisioning workflow
- [ ] VLAN configuration workflow

### 6.4 Manual Testing (8 hours)
- [ ] Cross-browser testing (Chrome, Firefox, Safari)
- [ ] Responsive testing (mobile, tablet, desktop)
- [ ] Accessibility testing
- [ ] Performance testing

---

## 🚀 Deployment & DevOps

### 7.1 Build Configuration (4 hours)
- [ ] Production build optimization
- [ ] Code splitting
- [ ] Asset optimization
- [ ] Source map configuration

### 7.2 VPS Deployment Script (8 hours)
- [ ] Build script
- [ ] Compression script
- [ ] SCP/SFTP transfer script
- [ ] Deployment automation
- [ ] Rollback mechanism
- [ ] Health check

**Deployment Script (`deploy-frontend.sh`):**
```bash
#!/bin/bash
# 1. Build production
# 2. Compress dist/
# 3. Transfer to VPS
# 4. Extract on server
# 5. Update Nginx config
# 6. Reload Nginx
# 7. Health check
```

### 7.3 Server Configuration (6 hours)
- [ ] Nginx configuration
- [ ] SSL/TLS setup (Let's Encrypt)
- [ ] Reverse proxy to API
- [ ] Gzip compression
- [ ] Caching headers
- [ ] Security headers

### 7.4 Monitoring Setup (4 hours)
- [ ] Error tracking (Sentry optional)
- [ ] Performance monitoring
- [ ] Uptime monitoring
- [ ] Log aggregation

---

## 📚 Documentation

### 8.1 User Documentation (12 hours)
- [ ] User guide
- [ ] Feature documentation
- [ ] Screenshots & GIFs
- [ ] FAQ section
- [ ] Troubleshooting guide

### 8.2 Technical Documentation (8 hours)
- [ ] Architecture overview
- [ ] Component documentation
- [ ] API integration guide
- [ ] Deployment guide
- [ ] Development setup guide

### 8.3 Code Documentation (6 hours)
- [ ] JSDoc comments
- [ ] README files
- [ ] Inline comments
- [ ] Type definitions

---

## ⏱️ Time Estimation Summary

| Category | Tasks | Estimated Hours |
|----------|-------|-----------------|
| **Project Setup** | 3 | 14 |
| **UI/UX Components** | 11 | 52 |
| **API Integration** | 8 | 48 |
| **Page Development** | 7 | 116 |
| **Theme Implementation** | 3 | 12 |
| **Testing & QA** | 4 | 56 |
| **Deployment & DevOps** | 4 | 22 |
| **Documentation** | 3 | 26 |
| **TOTAL** | **43** | **346 hours** |

**Estimated Timeline:**
- **Full-time (8h/day):** ~43 working days (~9 weeks)
- **Part-time (4h/day):** ~86 working days (~18 weeks)
- **Weekend only (8h/week):** ~43 weeks

---

## 📝 Priority Matrix

### 🔴 High Priority (Must Have - Phase 1-3)
- Project setup & configuration
- Authentication
- Dashboard overview
- ONU monitoring
- Basic provisioning
- VLAN management
- Responsive layout
- Light/Dark theme

### 🟡 Medium Priority (Should Have - Phase 4-5)
- Advanced provisioning features
- Traffic management
- Batch operations
- Configuration backup
- Charts & analytics

### 🟢 Low Priority (Nice to Have - Phase 6-7)
- User management
- Advanced analytics
- Custom reports
- Alarm system
- Multi-language support

---

## ✅ Task Completion Checklist

### Phase 1: Foundation
- [ ] Project initialized
- [ ] Dependencies installed
- [ ] Folder structure created
- [ ] Theme system implemented
- [ ] Layout components built
- [ ] API client configured

### Phase 2: Core Features
- [ ] Dashboard completed
- [ ] Monitoring page completed
- [ ] ONU detail view completed
- [ ] Real-time updates working

### Phase 3: Provisioning
- [ ] Auto-discovery implemented
- [ ] Manual provision wizard completed
- [ ] Batch operations working

### Phase 4: Management
- [ ] VLAN management completed
- [ ] Traffic profiles completed
- [ ] ONU operations completed

### Phase 5: Finalization
- [ ] All tests passing
- [ ] Documentation completed
- [ ] Deployment script working
- [ ] Production deployed

---

**Last Updated:** January 12, 2026  
**Next Review:** Weekly during development

