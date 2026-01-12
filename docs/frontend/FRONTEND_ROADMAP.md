# Frontend Development Roadmap
## ZTE C320 OLT Management Dashboard

**Project:** Web-based Dashboard for ZTE C320 OLT Management  
**Status:** Planning Phase  
**Last Updated:** January 12, 2026

---

## 🎯 Project Overview

Membangun web dashboard modern dan responsif untuk mengelola ZTE C320 OLT dengan fitur monitoring real-time, provisioning ONU, dan konfigurasi VLAN menggunakan REST API yang telah tersedia.

---

## 📅 Development Phases

### Phase 1: Foundation & Core Setup (Week 1-2)
**Status:** 🟡 Planning

#### Objectives
- Setup project dengan modern tech stack
- Implementasi autentikasi dasar
- Desain sistem dan UI/UX framework
- Konfigurasi deployment ke VPS

#### Deliverables
- ✅ Project initialization dengan Vite + React
- ✅ TailwindCSS + Shadcn/ui setup
- ✅ Light/Dark theme implementation
- ✅ Responsive layout dengan sidebar
- ✅ Basic authentication flow
- ✅ API client configuration (Axios)
- ✅ Environment configuration
- ✅ Deployment script ke VPS

#### Tech Stack
- **Framework:** React 18 + TypeScript
- **Build Tool:** Vite
- **Styling:** TailwindCSS + Shadcn/ui
- **State Management:** Zustand
- **API Client:** Axios + React Query
- **Routing:** React Router v6
- **Charts:** Recharts / Chart.js
- **Icons:** Lucide React
- **Forms:** React Hook Form + Zod

---

### Phase 2: Dashboard & Monitoring (Week 3-4)
**Status:** 🔴 Not Started

#### Objectives
- Dashboard overview dengan statistik real-time
- Monitoring ONU per PON port
- Visualisasi optical power dan signal quality

#### Features
1. **Dashboard Overview**
   - Total ONU online/offline cards
   - Signal quality distribution chart
   - Recent activities log
   - System status indicators
   - PON port utilization chart

2. **ONU Monitoring Page**
   - Filter by Board/PON/Status
   - Real-time data table dengan auto-refresh
   - Optical power visualization (RX/TX)
   - Export data (CSV/Excel)
   - Search & pagination

3. **PON Port Management**
   - PON port list dengan status
   - Per-port ONU statistics
   - Signal quality heatmap

#### API Endpoints Integration
- `GET /board/{board_id}/pon/{pon_id}/`
- `GET /board/{board_id}/pon/{pon_id}/onu/{onu_id}`
- `GET /board/{board_id}/pon/{pon_id}/info`
- `GET /monitoring/realtime/board/{board_id}/pon/{pon_id}/onu/{onu_id}`

---

### Phase 3: ONU Provisioning (Week 5-6)
**Status:** 🔴 Not Started

#### Objectives
- Wizard-based ONU provisioning
- Auto-discovery ONU management
- Bulk provisioning support

#### Features
1. **ONU Auto-Discovery**
   - List unprovisioned ONUs
   - One-click authorization
   - Batch authorization

2. **Manual ONU Provisioning**
   - Step-by-step wizard
   - Serial number validation
   - ONU type selection
   - Name & description assignment
   - Traffic profile selection
   - VLAN configuration

3. **Provisioning Templates**
   - Save configuration templates
   - Template management
   - Quick provisioning dari template

#### API Endpoints Integration
- `GET /provision/unconfigured/board/{board_id}/pon/{pon_id}`
- `POST /provision/authorize`
- `POST /provision/create`
- `POST /provision/batch`

---

### Phase 4: VLAN & Traffic Management (Week 7-8)
**Status:** 🔴 Not Started

#### Objectives
- VLAN configuration interface
- Traffic profile management
- Service-port assignment

#### Features
1. **VLAN Management**
   - Create/edit/delete VLAN
   - Service-port configuration
   - VLAN assignment per ONU
   - Bulk VLAN operations

2. **Traffic Profiles**
   - DBA profile management
   - T-CONT configuration
   - GEM port settings
   - Profile templates

3. **Service Configuration**
   - Per-ONU service setup
   - Bandwidth allocation
   - QoS configuration

#### API Endpoints Integration
- `POST /vlan/create`
- `DELETE /vlan/delete`
- `POST /traffic/dba-profile`
- `POST /traffic/tcont`
- `POST /traffic/gem-port`

---

### Phase 5: ONU Management Operations (Week 9-10)
**Status:** 🔴 Not Started

#### Objectives
- ONU lifecycle management
- Batch operations interface
- Configuration backup/restore

#### Features
1. **ONU Operations**
   - Reboot ONU
   - Block/Unblock ONU
   - Delete ONU
   - Update description
   - Reset configuration

2. **Batch Operations**
   - Multi-select ONUs
   - Batch reboot
   - Batch delete
   - Batch provisioning
   - Progress tracking

3. **Configuration Management**
   - Backup ONU configuration
   - Restore configuration
   - Configuration history
   - Export/import configs

#### API Endpoints Integration
- `POST /management/reboot`
- `POST /management/block`
- `POST /management/unblock`
- `DELETE /management/delete`
- `PUT /management/description`
- `POST /batch/provision`
- `POST /batch/delete`
- `POST /batch/reboot`
- `POST /config/backup`
- `POST /config/restore`

---

### Phase 6: Advanced Features (Week 11-12)
**Status:** 🔴 Not Started

#### Objectives
- Alarm & notification system
- Advanced analytics
- User management
- System settings

#### Features
1. **Alarm System**
   - Real-time alarm monitoring
   - Alarm history & log
   - Custom alarm rules
   - Notification settings (email/webhook)

2. **Analytics & Reports**
   - Signal quality trends
   - ONU online/offline patterns
   - Bandwidth utilization reports
   - Performance statistics
   - Custom reports generator

3. **User Management**
   - Role-based access control
   - User CRUD operations
   - Activity logging
   - Session management

4. **System Settings**
   - OLT configuration
   - API endpoint settings
   - Theme preferences
   - Language selection (ID/EN)
   - Backup/restore settings

---

### Phase 7: Optimization & Production (Week 13-14)
**Status:** 🔴 Not Started

#### Objectives
- Performance optimization
- Security hardening
- Production deployment
- Documentation

#### Tasks
1. **Performance**
   - Code splitting & lazy loading
   - Image optimization
   - Caching strategy
   - Bundle size optimization
   - Lighthouse score > 90

2. **Security**
   - XSS protection
   - CSRF tokens
   - Secure authentication
   - API key management
   - HTTPS enforcement

3. **Testing**
   - Unit tests (Jest)
   - Integration tests
   - E2E tests (Playwright)
   - Load testing

4. **Documentation**
   - User guide
   - Admin guide
   - API integration docs
   - Deployment guide
   - Troubleshooting guide

5. **Deployment**
   - VPS deployment automation
   - Nginx configuration
   - SSL/TLS setup
   - Monitoring setup
   - Backup strategy

---

## 🎨 Design Principles

### Visual Design
- **Modern & Clean:** Minimalist interface dengan focus pada data
- **Professional:** Color scheme yang professional dan tidak mengganggu
- **Intuitive:** User flow yang jelas dan mudah dipahami
- **Consistent:** Design system yang konsisten di seluruh aplikasi

### Color Scheme
#### Light Mode
- Primary: Blue (#3b82f6)
- Secondary: Slate (#64748b)
- Success: Green (#22c55e)
- Warning: Yellow (#eab308)
- Error: Red (#ef4444)
- Background: White (#ffffff)
- Surface: Gray (#f8fafc)

#### Dark Mode
- Primary: Blue (#60a5fa)
- Secondary: Slate (#94a3b8)
- Success: Green (#4ade80)
- Warning: Yellow (#facc15)
- Error: Red (#f87171)
- Background: Dark (#0f172a)
- Surface: Dark Gray (#1e293b)

### Responsive Breakpoints
- Mobile: < 640px
- Tablet: 640px - 1024px
- Desktop: > 1024px
- Large Desktop: > 1536px

---

## 📊 Success Metrics

### Performance
- [ ] First Contentful Paint < 1.5s
- [ ] Time to Interactive < 3s
- [ ] Lighthouse Score > 90
- [ ] Bundle size < 500KB (gzipped)

### User Experience
- [ ] Mobile-friendly (responsive design)
- [ ] Accessibility score > 90
- [ ] Support all modern browsers
- [ ] Smooth animations (60fps)

### Functionality
- [ ] All 50+ API endpoints integrated
- [ ] Real-time data updates
- [ ] Error handling & retry logic
- [ ] Offline detection

---

## 🚀 Deployment Strategy

### VPS Deployment (No GitHub)
1. Build production bundle locally
2. Transfer ke VPS via SCP/SFTP
3. Deploy dengan Nginx
4. Setup SSL dengan Let's Encrypt
5. Configure reverse proxy ke backend API

### Automated Deployment Script
```bash
# Deploy script akan dibuat untuk:
- Build production bundle
- Compress files
- Transfer to VPS
- Extract & deploy
- Reload Nginx
- Health check
```

---

## 🔄 Maintenance & Updates

### Regular Updates
- Weekly dependency updates
- Security patches monitoring
- Performance monitoring
- User feedback implementation

### Version Control
- Semantic versioning (v1.0.0)
- Changelog documentation
- Backup sebelum update major

---

## 📝 Notes

### Deployment Philosophy
- Frontend **TIDAK** di-push ke GitHub
- Deployment **HANYA** ke VPS production
- Source code disimpan local/private repository
- Build artifacts di-transfer langsung ke server

### API Integration
- Base URL configurable via environment
- Automatic retry untuk failed requests
- Request/response logging
- Error boundary implementation

---

## ✅ Definition of Done

Setiap phase dianggap selesai jika:
- [ ] Semua fitur berfungsi sesuai requirement
- [ ] Responsive di semua breakpoints
- [ ] Light/Dark theme berfungsi sempurna
- [ ] No console errors
- [ ] Code reviewed
- [ ] Tested di Chrome, Firefox, Safari
- [ ] Deployment script tested
- [ ] Documentation updated

---

**Timeline Total:** 14 Weeks  
**Start Date:** TBD  
**Expected Completion:** TBD

