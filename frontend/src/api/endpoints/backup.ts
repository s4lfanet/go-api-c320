import { apiClient } from '../client';

export interface Backup {
  id: string;
  type: string;
  created_at: string;
  description?: string;
  onu_count?: number;
}

export const configBackupApi = {
  // Backup operations
  backupONU: async (pon: number, onuId: number): Promise<any> => {
    return apiClient.post(`/config/backup/onu/${pon}/${onuId}`, {});
  },

  backupOLT: async (): Promise<any> => {
    return apiClient.post('/config/backup/olt', {});
  },

  importBackup: async (file: File): Promise<any> => {
    const formData = new FormData();
    formData.append('file', file);
    return apiClient.post('/config/backup/import', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
  },

  // Backup management
  listBackups: async (): Promise<any> => {
    return apiClient.get('/config/backups');
  },

  getBackup: async (backupId: string): Promise<any> => {
    return apiClient.get(`/config/backup/${backupId}`);
  },

  deleteBackup: async (backupId: string): Promise<any> => {
    return apiClient.delete(`/config/backup/${backupId}`);
  },

  exportBackup: async (backupId: string): Promise<any> => {
    return apiClient.get(`/config/backup/${backupId}/export`, {
      responseType: 'blob',
    });
  },

  // Restore operations
  restoreFromBackup: async (backupId: string): Promise<any> => {
    return apiClient.post(`/config/restore/${backupId}`, {});
  },
};
