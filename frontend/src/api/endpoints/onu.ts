import { apiClient } from '../client';
import type {
  OnuListResponse,
  OnuDetailResponse,
  PonInfoResponse,
  DashboardStats,
} from '../types';

// OLT Monitoring response type
interface OLTMonitoringResponse {
  code: number;
  status: string;
  data: {
    total_onus: number;
    online_onus: number;
    offline_onus: number;
    pon_ports: Array<{
      pon_port: string;
      pon_index: number;
      onu_count: number;
      online_count: number;
      offline_count: number;
      statistics: {
        rx_packets: number;
        rx_bytes: number;
        rx_rate: string;
      };
      last_update: string;
    }>;
    last_update: string;
  };
}

export const onuApi = {
  // Get ONU list on a PON port (tanpa trailing slash)
  getOnuList: async (board: number, pon: number): Promise<OnuListResponse> => {
    return apiClient.get(`/board/${board}/pon/${pon}`);
  },

  // Get specific ONU details
  getOnuDetail: async (
    board: number,
    pon: number,
    onuId: number
  ): Promise<OnuDetailResponse> => {
    return apiClient.get(`/board/${board}/pon/${pon}/onu/${onuId}`);
  },

  // Get PON port information
  getPonInfo: async (board: number, pon: number): Promise<PonInfoResponse> => {
    return apiClient.get(`/board/${board}/pon/${pon}/info`);
  },

  // Get OLT monitoring summary - use actual backend endpoint
  getOLTMonitoring: async (): Promise<OLTMonitoringResponse> => {
    return apiClient.get('/monitoring/olt');
  },

  // Get PON monitoring
  getPONMonitoring: async (pon: number): Promise<any> => {
    return apiClient.get(`/monitoring/pon/${pon}`);
  },

  // Get ONU monitoring
  getONUMonitoring: async (pon: number, onuId: number): Promise<any> => {
    return apiClient.get(`/monitoring/onu/${pon}/${onuId}`);
  },

  // Get dashboard statistics using OLT monitoring endpoint
  getDashboardStats: async (): Promise<DashboardStats> => {
    try {
      const response: any = await apiClient.get('/monitoring/olt');
      const data = response.data || response;
      
      // Aggregate unique PON ports (remove duplicates by pon_port)
      const uniquePonPorts = new Map<string, any>();
      if (data.pon_ports) {
        data.pon_ports.forEach((p: any) => {
          if (!uniquePonPorts.has(p.pon_port) || p.onu_count > 0) {
            uniquePonPorts.set(p.pon_port, p);
          }
        });
      }

      const ponPorts = Array.from(uniquePonPorts.values())
        .filter((p: any) => p.onu_count > 0)
        .map((p: any) => ({
          board: 1, // Default board
          pon: parseInt(p.pon_port),
          total: p.onu_count,
          online: p.online_count,
          offline: p.offline_count,
        }))
        .sort((a, b) => a.pon - b.pon);

      const totalOnus = data.total_onus || 0;
      const onlineOnus = data.online_onus || 0;
      const offlineOnus = data.offline_onus || 0;

      return {
        totalOnus,
        onlineOnus,
        offlineOnus,
        alerts: offlineOnus, // Offline ONUs as alerts
        uptime: totalOnus > 0 ? ((onlineOnus / totalOnus) * 100).toFixed(1) : '0',
        ponPorts,
      };
    } catch (error) {
      console.error('Failed to get dashboard stats:', error);
      return {
        totalOnus: 0,
        onlineOnus: 0,
        offlineOnus: 0,
        alerts: 0,
        uptime: '0',
        ponPorts: [],
      };
    }
  },
};
