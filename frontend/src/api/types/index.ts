// API Response Types
export interface ApiResponse<T> {
  code: number;
  status: string;
  data: T;
}

// ONU Types
export interface Onu {
  board: number;
  pon: number;
  onu_id: number;
  name: string;
  description?: string;
  onu_type: string;
  serial_number: string;
  rx_power: string;
  tx_power?: string;
  status: string;
  ip_address?: string;
  last_online?: string;
  last_offline?: string;
  uptime?: string;
  last_down_time_duration?: string;
  offline_reason?: string;
  gpon_optical_distance?: string;
}

export interface OnuListResponse extends ApiResponse<Onu[]> {}
export interface OnuDetailResponse extends ApiResponse<Onu> {}

// PON Types
export interface PonInfo {
  board: number;
  pon: number;
  admin_status: string;
  operational_status: string;
  total_onus: number;
  online_onus: number;
  offline_onus: number;
  rx_power?: string;
  tx_power?: string;
}

export interface PonInfoResponse extends ApiResponse<PonInfo> {}

// Provision Types
export interface UnconfiguredOnu {
  board: number;
  pon: number;
  serial_number: string;
  onu_type?: string;
  last_seen?: string;
}

export interface UnconfiguredOnuResponse extends ApiResponse<UnconfiguredOnu[]> {}

export interface ProvisionPayload {
  board: number;
  pon: number;
  onu_id: number;
  serial_number: string;
  onu_type: string;
  name: string;
  description?: string;
}

// Management Types
export interface ManagementPayload {
  board: number;
  pon: number;
  onu_id: number;
  description?: string;
}

// VLAN Types
export interface VlanPayload {
  board: number;
  pon: number;
  onu_id: number;
  vlan_id: number;
  svlan?: number;
  cvlan?: number;
  priority?: number;
}

// Error Type
export interface ApiError {
  message: string;
  status?: number;
  data?: any;
}

// Dashboard Statistics
export interface DashboardStats {
  totalOnus: number;
  onlineOnus: number;
  offlineOnus: number;
  alerts: number;
  uptime: string;
  ponPorts: Array<{
    board: number;
    pon: number;
    total: number;
    online: number;
    offline: number;
  }>;
}
