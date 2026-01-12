export const APP_NAME = import.meta.env.VITE_APP_NAME || 'ZTE C320 OLT Dashboard';
export const APP_VERSION = import.meta.env.VITE_APP_VERSION || '1.0.0';
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1';
export const AUTO_REFRESH_INTERVAL = parseInt(
  import.meta.env.VITE_AUTO_REFRESH_INTERVAL || '30000'
);

// OLT Configuration
export const MAX_BOARD_ID = 20;
export const MAX_PON_ID = 16;
export const MAX_ONU_ID = 128;

// Status
export const ONU_STATUS = {
  ONLINE: 'Online',
  OFFLINE: 'Offline',
  LOGGING: 'Logging',
  LOS: 'LOS',
  DYINGGASP: 'DyingGasp',
} as const;

// Signal Quality Thresholds (dBm)
export const SIGNAL_THRESHOLDS = {
  EXCELLENT: -20,
  GOOD: -23,
  FAIR: -25,
  POOR: -27,
};

// Colors
export const STATUS_COLORS = {
  online: 'text-green-500',
  offline: 'text-red-500',
  warning: 'text-yellow-500',
} as const;

// Routes
export const ROUTES = {
  DASHBOARD: '/',
  MONITORING: '/monitoring',
  PROVISION_AUTO: '/provision/auto',
  PROVISION_MANUAL: '/provision/manual',
  PROVISION_BATCH: '/provision/batch',
  VLAN: '/vlan',
  TRAFFIC: '/traffic',
  ONU: '/onu',
  BACKUP: '/backup',
  SETTINGS: '/settings',
} as const;
