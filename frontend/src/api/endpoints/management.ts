import { apiClient } from '../client';

export interface ONUAction {
  pon_port: number;
  onu_id: number;
}

export interface ONUDescription {
  pon_port: number;
  onu_id: number;
  description: string;
}

export interface BatchONUAction {
  onus: ONUAction[];
}

export interface BatchDescriptions {
  onus: ONUDescription[];
}

export const onuManagementApi = {
  // Single ONU operations
  rebootONU: async (action: ONUAction): Promise<any> => {
    return apiClient.post('/onu-management/reboot', action);
  },

  blockONU: async (action: ONUAction): Promise<any> => {
    return apiClient.post('/onu-management/block', action);
  },

  unblockONU: async (action: ONUAction): Promise<any> => {
    return apiClient.post('/onu-management/unblock', action);
  },

  updateDescription: async (data: ONUDescription): Promise<any> => {
    return apiClient.put('/onu-management/description', data);
  },

  deleteONU: async (pon: number, onuId: number): Promise<any> => {
    return apiClient.delete(`/onu-management/${pon}/${onuId}`);
  },

  // Batch operations
  batchReboot: async (data: BatchONUAction): Promise<any> => {
    return apiClient.post('/batch/reboot', data);
  },

  batchBlock: async (data: BatchONUAction): Promise<any> => {
    return apiClient.post('/batch/block', data);
  },

  batchUnblock: async (data: BatchONUAction): Promise<any> => {
    return apiClient.post('/batch/unblock', data);
  },

  batchDelete: async (data: BatchONUAction): Promise<any> => {
    return apiClient.post('/batch/delete', data);
  },

  batchUpdateDescriptions: async (data: BatchDescriptions): Promise<any> => {
    return apiClient.put('/batch/descriptions', data);
  },
};
