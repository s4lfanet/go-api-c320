import { createBrowserRouter } from 'react-router-dom';
import { AppLayout } from './components/layout/AppLayout';
import Dashboard from './pages/Dashboard';
import Monitoring from './pages/Monitoring';
import Provisioning from './pages/Provisioning';
import VLAN from './pages/VLAN';
import Traffic from './pages/Traffic';
import ONUManagement from './pages/ONUManagement';
import ConfigBackup from './pages/ConfigBackup';
import Settings from './pages/Settings';
import NotFound from './pages/ErrorPages/NotFound';

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
        path: 'provisioning',
        element: <Provisioning />,
      },
      {
        path: 'vlan',
        element: <VLAN />,
      },
      {
        path: 'traffic',
        element: <Traffic />,
      },
      {
        path: 'onu-management',
        element: <ONUManagement />,
      },
      {
        path: 'config-backup',
        element: <ConfigBackup />,
      },
      {
        path: 'settings',
        element: <Settings />,
      },
    ],
  },
  {
    path: '*',
    element: <NotFound />,
  },
]);
