import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { configBackupApi, Backup } from '@/api/endpoints/backup';
import { RefreshCw, Download, Upload, Trash2, RotateCcw, Database, Server, HardDrive } from 'lucide-react';

export default function ConfigBackup() {
  const queryClient = useQueryClient();
  const [selectedBackup, setSelectedBackup] = useState<Backup | null>(null);

  // Get all backups
  const { data, isLoading, refetch } = useQuery({
    queryKey: ['backups'],
    queryFn: () => configBackupApi.listBackups(),
  });

  // Backup OLT mutation
  const backupOLTMutation = useMutation({
    mutationFn: configBackupApi.backupOLT,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['backups'] });
      alert('OLT backup created successfully!');
    },
    onError: (error: any) => {
      alert(`Failed to create backup: ${error.message}`);
    },
  });

  // Delete backup mutation
  const deleteMutation = useMutation({
    mutationFn: (backupId: string) => configBackupApi.deleteBackup(backupId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['backups'] });
      setSelectedBackup(null);
      alert('Backup deleted successfully!');
    },
    onError: (error: any) => {
      alert(`Failed to delete backup: ${error.message}`);
    },
  });

  // Restore mutation
  const restoreMutation = useMutation({
    mutationFn: (backupId: string) => configBackupApi.restoreFromBackup(backupId),
    onSuccess: () => {
      alert('Configuration restored successfully!');
    },
    onError: (error: any) => {
      alert(`Failed to restore: ${error.message}`);
    },
  });

  const backups = (data as any)?.data || [];

  const formatDate = (dateStr: string) => {
    try {
      return new Date(dateStr).toLocaleString();
    } catch {
      return dateStr;
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Config Backup</h1>
          <p className="text-muted-foreground">
            Backup and restore OLT configurations
          </p>
        </div>
        <div className="flex gap-2">
          <button
            onClick={() => {
              if (confirm('Create a full OLT configuration backup?')) {
                backupOLTMutation.mutate();
              }
            }}
            disabled={backupOLTMutation.isPending}
            className="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-md bg-primary text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
          >
            <Database className="h-4 w-4" />
            {backupOLTMutation.isPending ? 'Creating...' : 'Backup OLT'}
          </button>
          <button
            onClick={() => refetch()}
            disabled={isLoading}
            className="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-md border border-input bg-background hover:bg-accent hover:text-accent-foreground disabled:opacity-50"
          >
            <RefreshCw className={`h-4 w-4 ${isLoading ? 'animate-spin' : ''}`} />
          </button>
        </div>
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        {/* Backup List */}
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <HardDrive className="h-5 w-5" />
              Saved Backups
            </CardTitle>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="text-center py-8 text-muted-foreground">Loading backups...</div>
            ) : backups.length === 0 ? (
              <div className="text-center py-8">
                <HardDrive className="h-12 w-12 mx-auto text-muted-foreground mb-2" />
                <p className="text-muted-foreground">No backups found</p>
                <p className="text-xs text-muted-foreground mt-1">
                  Create a backup to protect your configuration
                </p>
              </div>
            ) : (
              <div className="space-y-2">
                {backups.map((backup: Backup) => (
                  <div
                    key={backup.id}
                    onClick={() => setSelectedBackup(backup)}
                    className={`p-4 rounded-lg border cursor-pointer transition-colors ${
                      selectedBackup?.id === backup.id
                        ? 'border-primary bg-primary/10'
                        : 'border-border hover:bg-muted'
                    }`}
                  >
                    <div className="flex justify-between items-start">
                      <div>
                        <div className="font-medium flex items-center gap-2">
                          {backup.type === 'olt' ? (
                            <Server className="h-4 w-4 text-blue-500" />
                          ) : (
                            <Database className="h-4 w-4 text-green-500" />
                          )}
                          {backup.type === 'olt' ? 'Full OLT Backup' : 'ONU Backup'}
                        </div>
                        <div className="text-xs text-muted-foreground mt-1">
                          {formatDate(backup.created_at)}
                        </div>
                        {backup.description && (
                          <div className="text-sm mt-1">{backup.description}</div>
                        )}
                      </div>
                      <div className="text-xs text-muted-foreground">
                        {backup.onu_count && `${backup.onu_count} ONUs`}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Backup Actions */}
        <Card>
          <CardHeader>
            <CardTitle>Backup Actions</CardTitle>
          </CardHeader>
          <CardContent>
            {selectedBackup ? (
              <div className="space-y-4">
                <div className="p-3 bg-muted rounded-lg">
                  <div className="text-xs text-muted-foreground">Selected Backup</div>
                  <div className="font-medium">{selectedBackup.type === 'olt' ? 'Full OLT' : 'ONU'} Backup</div>
                  <div className="text-xs text-muted-foreground mt-1">
                    {formatDate(selectedBackup.created_at)}
                  </div>
                </div>

                <div className="space-y-2">
                  <button
                    onClick={() => {
                      if (confirm(`Restore configuration from this backup?\n\nThis will overwrite current settings!`)) {
                        restoreMutation.mutate(selectedBackup.id);
                      }
                    }}
                    disabled={restoreMutation.isPending}
                    className="w-full inline-flex items-center justify-center gap-2 px-4 py-2 text-sm font-medium rounded-md bg-blue-500 text-white hover:bg-blue-600 disabled:opacity-50"
                  >
                    <RotateCcw className="h-4 w-4" />
                    {restoreMutation.isPending ? 'Restoring...' : 'Restore'}
                  </button>

                  <button
                    onClick={async () => {
                      try {
                        const blob = await configBackupApi.exportBackup(selectedBackup.id);
                        const url = window.URL.createObjectURL(blob);
                        const a = document.createElement('a');
                        a.href = url;
                        a.download = `backup-${selectedBackup.id}.json`;
                        a.click();
                        window.URL.revokeObjectURL(url);
                      } catch (error: any) {
                        alert(`Failed to export: ${error.message}`);
                      }
                    }}
                    className="w-full inline-flex items-center justify-center gap-2 px-4 py-2 text-sm font-medium rounded-md border border-input bg-background hover:bg-accent"
                  >
                    <Download className="h-4 w-4" />
                    Export
                  </button>

                  <button
                    onClick={() => {
                      if (confirm(`Delete this backup?\n\nThis action cannot be undone!`)) {
                        deleteMutation.mutate(selectedBackup.id);
                      }
                    }}
                    disabled={deleteMutation.isPending}
                    className="w-full inline-flex items-center justify-center gap-2 px-4 py-2 text-sm font-medium rounded-md bg-red-500 text-white hover:bg-red-600 disabled:opacity-50"
                  >
                    <Trash2 className="h-4 w-4" />
                    {deleteMutation.isPending ? 'Deleting...' : 'Delete'}
                  </button>
                </div>
              </div>
            ) : (
              <div className="text-center py-8 text-muted-foreground">
                <Database className="h-12 w-12 mx-auto mb-2 opacity-50" />
                <p>Select a backup to manage</p>
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Quick Actions */}
      <Card>
        <CardHeader>
          <CardTitle>Quick Backup</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-2">
            <div className="p-4 border rounded-lg">
              <div className="flex items-center gap-2 mb-2">
                <Server className="h-5 w-5 text-blue-500" />
                <span className="font-medium">Full OLT Backup</span>
              </div>
              <p className="text-sm text-muted-foreground mb-3">
                Backup all ONU configurations, VLAN settings, and profiles
              </p>
              <button
                onClick={() => {
                  if (confirm('Create full OLT backup?')) {
                    backupOLTMutation.mutate();
                  }
                }}
                disabled={backupOLTMutation.isPending}
                className="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-md bg-blue-500 text-white hover:bg-blue-600 disabled:opacity-50"
              >
                <Database className="h-4 w-4" />
                Create Backup
              </button>
            </div>

            <div className="p-4 border rounded-lg">
              <div className="flex items-center gap-2 mb-2">
                <Upload className="h-5 w-5 text-green-500" />
                <span className="font-medium">Import Backup</span>
              </div>
              <p className="text-sm text-muted-foreground mb-3">
                Import a previously exported backup file
              </p>
              <label className="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-md bg-green-500 text-white hover:bg-green-600 cursor-pointer">
                <Upload className="h-4 w-4" />
                Select File
                <input
                  type="file"
                  accept=".json"
                  className="hidden"
                  onChange={async (e) => {
                    const file = e.target.files?.[0];
                    if (file) {
                      try {
                        await configBackupApi.importBackup(file);
                        queryClient.invalidateQueries({ queryKey: ['backups'] });
                        alert('Backup imported successfully!');
                      } catch (error: any) {
                        alert(`Failed to import: ${error.message}`);
                      }
                    }
                  }}
                />
              </label>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
