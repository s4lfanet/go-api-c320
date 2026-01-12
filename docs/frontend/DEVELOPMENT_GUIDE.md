# Frontend Development Guide

Quick reference untuk development frontend ZTE C320 OLT Dashboard.

## 🚀 Quick Commands

```bash
# Development
npm run dev              # Start dev server (http://localhost:3000)

# Build
npm run build           # Production build
npm run preview         # Preview production build

# Code Quality
npm run lint            # Run ESLint
npm run format          # Format with Prettier

# Deployment
./deploy-frontend.sh    # Deploy to VPS (Linux/Mac)
.\deploy-frontend.ps1   # Deploy to VPS (Windows)
```

## 📂 Where to Add New Code

### New Page
```
src/pages/MyNewPage/
├── index.tsx           # Main page component
└── components/         # Page-specific components
```

### New API Endpoint
```
src/api/endpoints/
└── myendpoint.ts       # API functions
```

### New Component
```
src/components/features/
└── MyFeature.tsx       # Feature component
```

### New Store
```
src/store/
└── myStore.ts          # Zustand store
```

## 🎨 Using Components

### Button
```tsx
import { Button } from '@/components/ui/button';

<Button variant="default">Click me</Button>
<Button variant="outline">Outline</Button>
<Button variant="ghost">Ghost</Button>
```

### Card
```tsx
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';

<Card>
  <CardHeader>
    <CardTitle>Title</CardTitle>
  </CardHeader>
  <CardContent>
    Content here
  </CardContent>
</Card>
```

### Input
```tsx
import { Input } from '@/components/ui/input';

<Input type="text" placeholder="Enter text..." />
```

## 🔌 API Integration

### Fetching Data
```tsx
import { useQuery } from '@tanstack/react-query';
import { onuApi } from '@/api/endpoints/onu';

const MyComponent = () => {
  const { data, isLoading, error } = useQuery({
    queryKey: ['onu-list', 1, 1],
    queryFn: () => onuApi.getOnuList(1, 1),
  });

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error!</div>;

  return <div>{/* Use data */}</div>;
};
```

### Mutations (POST/PUT/DELETE)
```tsx
import { useMutation } from '@tanstack/react-query';

const mutation = useMutation({
  mutationFn: (data) => apiClient.post('/endpoint', data),
  onSuccess: () => {
    // Handle success
  },
});

mutation.mutate({ key: 'value' });
```

## 🗂️ State Management

### Zustand Store
```tsx
// store/myStore.ts
import { create } from 'zustand';

interface MyState {
  count: number;
  increment: () => void;
}

export const useMyStore = create<MyState>((set) => ({
  count: 0,
  increment: () => set((state) => ({ count: state.count + 1 })),
}));

// Usage in component
const count = useMyStore((state) => state.count);
const increment = useMyStore((state) => state.increment);
```

## 🎨 Theming

### Using Theme
```tsx
import { useThemeStore } from '@/store/themeStore';

const { theme, toggleTheme } = useThemeStore();
```

### Custom Colors
Edit `tailwind.config.js` untuk custom colors.

## 📱 Responsive Design

```tsx
// Mobile-first approach
<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
  {/* Content */}
</div>

// Breakpoints:
// sm: 640px
// md: 768px
// lg: 1024px
// xl: 1280px
// 2xl: 1536px
```

## 🔗 Routing

### Adding New Route
```tsx
// router.tsx
{
  path: 'my-route',
  element: <MyPage />,
}
```

### Navigation
```tsx
import { Link, useNavigate } from 'react-router-dom';

// Link
<Link to="/my-route">Go to page</Link>

// Programmatic
const navigate = useNavigate();
navigate('/my-route');
```

## 🛠️ Utilities

### Formatters
```tsx
import { formatDate, formatPower } from '@/utils/formatters';

formatDate('2024-01-12');           // "2024-01-12 10:30:00"
formatPower(-22.5);                  // "-22.50 dBm"
```

### Class Names
```tsx
import { cn } from '@/lib/utils';

<div className={cn('base-class', isActive && 'active-class')} />
```

## 🚨 Error Handling

### Error Boundary
```tsx
import { ErrorBoundary } from '@/components/common/ErrorBoundary';

<ErrorBoundary>
  <MyComponent />
</ErrorBoundary>
```

## 📦 Adding Dependencies

```bash
npm install package-name
```

Then import and use:
```tsx
import something from 'package-name';
```

## 🐛 Debugging

### React Query DevTools
Already configured. Open browser devtools.

### Console Logging
```tsx
console.log('Debug:', data);
```

### TypeScript Errors
Run type checking:
```bash
npm run build
```

## 📝 Best Practices

1. **Component naming:** PascalCase (e.g., `MyComponent`)
2. **File naming:** Same as component (e.g., `MyComponent.tsx`)
3. **Export:** Named exports preferred
4. **Types:** Always add TypeScript types
5. **Hooks:** Prefix with `use` (e.g., `useMyHook`)

## 🔍 Troubleshooting

### Port already in use
```bash
# Kill process on port 3000
# Windows:
netstat -ano | findstr :3000
taskkill /PID <PID> /F

# Linux/Mac:
lsof -ti:3000 | xargs kill -9
```

### Module not found
```bash
rm -rf node_modules package-lock.json
npm install
```

## 📚 Learn More

- [React Documentation](https://react.dev/)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [TailwindCSS](https://tailwindcss.com/docs)
- [Shadcn/ui](https://ui.shadcn.com/)
- [React Query](https://tanstack.com/query/latest)

---

Happy coding! 🚀
