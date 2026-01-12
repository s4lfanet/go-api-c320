import { apiClient } from '../client';

export interface UnconfiguredONU {
  serial_number: string;
  pon_port: string;
  onu_type: string;
  last_seen?: string;
}

export interface RegisterONUPayload {
  pon_port: number;
  onu_id: number;
  serial_number: string;
  onu_type: string;
  name: string;
  description?: string;
}

export const provisioningApi = {
  // Get all unconfigured ONUs (from all PON ports with status "Logging")
  getUnconfiguredONUs: async (): Promise<any> => {
    // Fallback: scan all PON ports for ONUs with status "Logging"
    const unconfigured: any[] = [];
    
    // Scan PON ports 1-16
    for (let pon = 1; pon <= 16; pon++) {
      try {
        const response = await apiClient.get(`/board/1/pon/${pon}`);
        const onus = response.data || response || [];
        
        // Filter ONUs with status "Logging" (unconfigured)
        const logging = onus.filter((onu: any) => 
          onu.status === 'Logging' || onu.status === 'logging'
        );
        
        unconfigured.push(...logging.map((onu: any) => ({
          serial_number: onu.serial_number,
          pon_port: onu.pon.toString(),
          onu_type: onu.onu_type || onu.type,
          onu_id: onu.onu_id,
          last_seen: new Date().toISOString()
        })));
      } catch (err) {
        // Skip PON ports with errors
        console.debug(`PON ${pon} scan skipped:`, err);
      }
    }
    
    return { data: unconfigured };
  },

  // Get unconfigured ONUs by PON port
  getUnconfiguredONUsByPON: async (pon: number): Promise<any> => {
    return apiClient.get(`/onu/unconfigured/${pon}`);
  },

  // Register new ONU
  registerONU: async (payload: RegisterONUPayload): Promise<any> => {
    return apiClient.post('/onu/register', payload);
  },

  // Delete ONU
  deleteONU: async (pon: number, onuId: number): Promise<any> => {
    return apiClient.delete(`/onu/${pon}/${onuId}`);
  },

  // Get empty ONU IDs for a PON port
  getEmptyOnuIds: async (board: number, pon: number): Promise<any> => {
    return apiClient.get(`/board/${board}/pon/${pon}/onu_id/empty`);
  },
};
