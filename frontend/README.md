# ZTE C320 OLT Management Dashboard

Modern web dashboard untuk monitoring dan mengelola ZTE C320 OLT dengan fitur real-time monitoring, provisioning ONU, dan konfigurasi VLAN.

> **⚠️ PENTING:** Frontend ini **TIDAK** di-push ke GitHub. Deployment hanya ke VPS production.

## 🎨 Features

- ✅ **Modern UI/UX** - Clean, professional, dan intuitive interface
- ✅ **Light & Dark Theme** - Automatic theme switching dengan preference persistence
- ✅ **Fully Responsive** - Optimized untuk semua ukuran layar (mobile, tablet, desktop)
- ✅ **Sidebar Navigation** - Collapsible sidebar dengan icon-based menu
- ✅ **Real-time Data** - Auto-refresh monitoring data
- ✅ **Type-Safe** - Full TypeScript coverage
- ✅ **Performance Optimized** - Code splitting, lazy loading, caching

## 🛠️ Tech Stack

### Core
- **React 18.2** - UI Library
- **TypeScript 5.3** - Type safety
- **Vite 5.0** - Build tool & dev server

### UI & Styling
- **TailwindCSS 3.4** - Utility-first CSS framework
- **Shadcn/ui** - High-quality component library
- **Lucide React** - Icon library
- **Recharts** - Charting library

### State Management
- **Zustand** - Lightweight state management
- **React Query** - Server state management

### Data Fetching
- **Axios** - HTTP client
- **React Query** - Caching & synchronization

### Routing
- **React Router v6** - Client-side routing

### Form Handling
- **React Hook Form** - Form management
- **Zod** - Schema validation

## 📋 Prerequisites

- Node.js 18+ dan npm
- Backend API running (default: `http://localhost:8081`)
- VPS dengan Nginx (untuk deployment)

## 🚀 Quick Start

### 1. Install Dependencies

```bash
cd frontend
npm install
```

### 2. Configure Environment

Copy `.env.example` ke `.env.development`:

```bash
cp .env.example .env.development
```

Edit `.env.development` sesuai konfigurasi:

```env
VITE_API_BASE_URL=http://localhost:8081/api/v1
VITE_APP_NAME=ZTE C320 OLT Dashboard (Dev)
VITE_AUTO_REFRESH_INTERVAL=30000
```

### 3. Run Development Server

```bash
npm run dev
```

Dashboard akan berjalan di `http://localhost:3000`

## 🏗️ Build untuk Production

```bash
npm run build
```

Build output akan berada di folder `dist/`.

## 🚀 Deployment ke VPS

### Menggunakan Deployment Script

#### Linux/Mac:

```bash
chmod +x deploy-frontend.sh
./deploy-frontend.sh production
```

#### Windows (PowerShell):

```powershell
.\deploy-frontend.ps1 -Environment production
```

### Prerequisites untuk Deployment

1. **SSH Key-based Authentication**
   - Setup SSH key untuk akses ke VPS tanpa password
   - Copy public key ke VPS: `ssh-copy-id user@vps-host`

2. **Environment Variables** (optional)
   ```bash
   export VPS_HOST=192.168.54.230
   export VPS_USER=root
   export VPS_PORT=22
   ```

3. **VPS Requirements**
   - Nginx installed dan running
   - Directory `/var/www/olt-dashboard` dengan write permissions

### Manual Deployment

Jika tidak menggunakan script:

```bash
# 1. Build
npm run build

# 2. Compress
tar -czf frontend-dist.tar.gz -C dist .

# 3. Transfer ke VPS
scp frontend-dist.tar.gz user@vps:/tmp/

# 4. SSH ke VPS dan extract
ssh user@vps
cd /var/www/olt-dashboard/frontend
tar -xzf /tmp/frontend-dist.tar.gz
chown -R www-data:www-data .
systemctl reload nginx
```

## 📁 Project Structure

```
frontend/
├── public/                 # Static assets
├── src/
│   ├── api/               # API client & endpoints
│   │   ├── client.ts
│   │   ├── queryClient.ts
│   │   ├── endpoints/
│   │   └── types/
│   ├── components/        # React components
│   │   ├── ui/           # Shadcn/ui components
│   │   ├── layout/       # Layout components
│   │   ├── common/       # Common components
│   │   └── features/     # Feature-specific components
│   ├── pages/            # Page components
│   ├── hooks/            # Custom React hooks
│   ├── store/            # Zustand stores
│   ├── utils/            # Utility functions
│   ├── styles/           # Global styles
│   ├── lib/              # Third-party configs
│   ├── types/            # TypeScript types
│   ├── main.tsx          # Entry point
│   └── router.tsx        # Routes configuration
├── deploy-frontend.sh    # Deployment script (Linux/Mac)
├── deploy-frontend.ps1   # Deployment script (Windows)
├── package.json
├── vite.config.ts
├── tsconfig.json
├── tailwind.config.js
└── .env.example
```

## 🎨 Theme

### Light Mode
- Primary: Blue (#3b82f6)
- Background: White (#ffffff)
- Surface: Gray (#f8fafc)

### Dark Mode
- Primary: Blue (#60a5fa)
- Background: Dark (#0f172a)
- Surface: Dark Gray (#1e293b)

Toggle theme dengan button di sidebar.

## 🔌 API Integration

Frontend berkomunikasi dengan backend API:

- **Base URL:** Configurable via `VITE_API_BASE_URL`
- **Authentication:** Bearer token (jika diimplementasikan)
- **Auto-retry:** Failed requests di-retry otomatis
- **Caching:** React Query cache strategy

### Example API Call

```typescript
import { useQuery } from '@tanstack/react-query';
import { onuApi } from '@/api/endpoints/onu';

const { data, isLoading } = useQuery({
  queryKey: ['onu-list', board, pon],
  queryFn: () => onuApi.getOnuList(board, pon),
  refetchInterval: 30000, // Auto-refresh every 30s
});
```

## 📱 Responsive Design

Dashboard fully responsive:

- **Mobile** (< 640px): Sidebar collapse to hamburger menu
- **Tablet** (640px - 1024px): Compact sidebar
- **Desktop** (> 1024px): Full sidebar dengan labels

## 🔒 Security

- XSS protection (React escapes by default)
- HTTPS enforcement di production
- Secure API communication
- Input validation dengan Zod schemas

## 📚 Documentation

Lengkap documentation tersedia di `docs/frontend/`:

- [FRONTEND_ROADMAP.md](../docs/frontend/FRONTEND_ROADMAP.md) - Development roadmap
- [FRONTEND_WORKLOAD.md](../docs/frontend/FRONTEND_WORKLOAD.md) - Task breakdown
- [FRONTEND_ARCHITECTURE.md](../docs/frontend/FRONTEND_ARCHITECTURE.md) - Technical architecture

## 🐛 Troubleshooting

### Build Errors

```bash
# Clear node_modules dan reinstall
rm -rf node_modules package-lock.json
npm install
npm run build
```

### API Connection Issues

1. Verify backend is running: `http://localhost:8081/api/v1/`
2. Check `.env.development` for correct API URL
3. Check browser console for CORS errors

### Deployment Issues

1. Verify SSH access: `ssh user@vps`
2. Check Nginx status: `systemctl status nginx`
3. Check Nginx logs: `tail -f /var/log/nginx/error.log`
4. Verify file permissions: `ls -la /var/www/olt-dashboard/frontend`

## 📝 Development Guidelines

### Code Style

- ESLint + Prettier configured
- Run linting: `npm run lint`
- Format code: `npm run format`

### Adding New Components

```typescript
// components/features/MyFeature.tsx
import { Card } from '@/components/ui/card';

export const MyFeature = () => {
  return (
    <Card>
      {/* Component content */}
    </Card>
  );
};
```

### Adding New API Endpoints

```typescript
// api/endpoints/myendpoint.ts
import { apiClient } from '../client';

export const myApi = {
  getData: async (id: number) => {
    return apiClient.get(`/my-endpoint/${id}`);
  },
};
```

## 🤝 Contributing

Karena ini private deployment (tidak di GitHub):

1. Develop di local branch
2. Test thoroughly
3. Deploy ke VPS dengan deployment script
4. Verify deployment success

## 📄 License

Private project - Not for public distribution

## 👤 Contact

Untuk pertanyaan atau issue, hubungi development team.

---

**Version:** 1.0.0  
**Last Updated:** January 12, 2026
