import { apiClient } from '../client';

export interface VLANConfig {
  pon_port: number;
  onu_id: number;
  service_port_id?: number;
  vlan_id: number;
  svlan?: number;
  cvlan?: number;
  vlan_mode?: string;
  cos?: number;
}

export const vlanApi = {
  // Get ONU VLAN configuration
  getONUVLAN: async (pon: number, onuId: number): Promise<any> => {
    return apiClient.get(`/vlan/onu/${pon}/${onuId}`);
  },

  // Get all service ports
  getAllServicePorts: async (): Promise<any> => {
    return apiClient.get('/vlan/service-ports');
  },

  // Configure VLAN
  configureVLAN: async (config: VLANConfig): Promise<any> => {
    return apiClient.post('/vlan/onu', config);
  },

  // Modify VLAN
  modifyVLAN: async (config: VLANConfig): Promise<any> => {
    return apiClient.put('/vlan/onu', config);
  },

  // Delete VLAN
  deleteVLAN: async (pon: number, onuId: number): Promise<any> => {
    return apiClient.delete(`/vlan/onu/${pon}/${onuId}`);
  },
};
