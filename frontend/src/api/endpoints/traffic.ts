import { apiClient } from '../client';

export interface DBAProfile {
  name: string;
  type: number;
  fixed_bandwidth?: number;
  assured_bandwidth?: number;
  maximum_bandwidth?: number;
}

export interface TCONTConfig {
  pon_port: number;
  onu_id: number;
  tcont_id: number;
  dba_profile_name: string;
}

export interface GEMPortConfig {
  pon_port: number;
  onu_id: number;
  gemport_id: number;
  tcont_id: number;
  vlan_id?: number;
}

export const trafficApi = {
  // DBA Profile operations
  getAllDBAProfiles: async (): Promise<any> => {
    return apiClient.get('/traffic/dba-profiles');
  },

  getDBAProfile: async (name: string): Promise<any> => {
    return apiClient.get(`/traffic/dba-profile/${name}`);
  },

  createDBAProfile: async (profile: DBAProfile): Promise<any> => {
    return apiClient.post('/traffic/dba-profile', profile);
  },

  modifyDBAProfile: async (profile: DBAProfile): Promise<any> => {
    return apiClient.put('/traffic/dba-profile', profile);
  },

  deleteDBAProfile: async (name: string): Promise<any> => {
    return apiClient.delete(`/traffic/dba-profile/${name}`);
  },

  // TCONT operations
  getONUTCONT: async (pon: number, onuId: number, tcontId: number): Promise<any> => {
    return apiClient.get(`/traffic/tcont/${pon}/${onuId}/${tcontId}`);
  },

  configureTCONT: async (config: TCONTConfig): Promise<any> => {
    return apiClient.post('/traffic/tcont', config);
  },

  deleteTCONT: async (pon: number, onuId: number, tcontId: number): Promise<any> => {
    return apiClient.delete(`/traffic/tcont/${pon}/${onuId}/${tcontId}`);
  },

  // GEMPort operations
  configureGEMPort: async (config: GEMPortConfig): Promise<any> => {
    return apiClient.post('/traffic/gemport', config);
  },

  deleteGEMPort: async (pon: number, onuId: number, gemportId: number): Promise<any> => {
    return apiClient.delete(`/traffic/gemport/${pon}/${onuId}/${gemportId}`);
  },
};
