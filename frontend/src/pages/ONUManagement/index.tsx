import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { useMutation } from '@tanstack/react-query';
import { onuManagementApi } from '@/api/endpoints/management';
import { Power, Ban, CheckCircle, Trash2, Edit, AlertTriangle } from 'lucide-react';

export default function ONUManagement() {
  const [pon, setPon] = useState(1);
  const [onuId, setOnuId] = useState(1);
  const [description, setDescription] = useState('');
  const [actionResult, setActionResult] = useState<{ type: 'success' | 'error'; message: string } | null>(null);

  // Reboot mutation
  const rebootMutation = useMutation({
    mutationFn: () => onuManagementApi.rebootONU({ pon_port: pon, onu_id: onuId }),
    onSuccess: () => setActionResult({ type: 'success', message: 'ONU reboot command sent successfully!' }),
    onError: (error: any) => setActionResult({ type: 'error', message: `Failed to reboot: ${error.message}` }),
  });

  // Block mutation
  const blockMutation = useMutation({
    mutationFn: () => onuManagementApi.blockONU({ pon_port: pon, onu_id: onuId }),
    onSuccess: () => setActionResult({ type: 'success', message: 'ONU blocked successfully!' }),
    onError: (error: any) => setActionResult({ type: 'error', message: `Failed to block: ${error.message}` }),
  });

  // Unblock mutation
  const unblockMutation = useMutation({
    mutationFn: () => onuManagementApi.unblockONU({ pon_port: pon, onu_id: onuId }),
    onSuccess: () => setActionResult({ type: 'success', message: 'ONU unblocked successfully!' }),
    onError: (error: any) => setActionResult({ type: 'error', message: `Failed to unblock: ${error.message}` }),
  });

  // Update description mutation
  const updateDescMutation = useMutation({
    mutationFn: () => onuManagementApi.updateDescription({ pon_port: pon, onu_id: onuId, description }),
    onSuccess: () => {
      setActionResult({ type: 'success', message: 'Description updated successfully!' });
      setDescription('');
    },
    onError: (error: any) => setActionResult({ type: 'error', message: `Failed to update: ${error.message}` }),
  });

  // Delete mutation
  const deleteMutation = useMutation({
    mutationFn: () => onuManagementApi.deleteONU(pon, onuId),
    onSuccess: () => setActionResult({ type: 'success', message: 'ONU deleted successfully!' }),
    onError: (error: any) => setActionResult({ type: 'error', message: `Failed to delete: ${error.message}` }),
  });

  const isLoading = rebootMutation.isPending || blockMutation.isPending || 
                    unblockMutation.isPending || updateDescMutation.isPending || deleteMutation.isPending;

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">ONU Management</h1>
        <p className="text-muted-foreground">
          Reboot, block, unblock, and manage ONU configurations
        </p>
      </div>

      {/* ONU Selector */}
      <Card>
        <CardHeader>
          <CardTitle>Select ONU</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex gap-4">
            <div className="flex-1">
              <label className="text-sm font-medium mb-1 block">PON Port</label>
              <Input
                type="number"
                min={1}
                max={16}
                value={pon}
                onChange={(e) => {
                  setPon(parseInt(e.target.value));
                  setActionResult(null);
                }}
              />
            </div>
            <div className="flex-1">
              <label className="text-sm font-medium mb-1 block">ONU ID</label>
              <Input
                type="number"
                min={1}
                max={128}
                value={onuId}
                onChange={(e) => {
                  setOnuId(parseInt(e.target.value));
                  setActionResult(null);
                }}
              />
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Action Result */}
      {actionResult && (
        <div className={`p-4 rounded-lg ${
          actionResult.type === 'success' 
            ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200' 
            : 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
        }`}>
          <div className="flex items-center gap-2">
            {actionResult.type === 'success' ? (
              <CheckCircle className="h-5 w-5" />
            ) : (
              <AlertTriangle className="h-5 w-5" />
            )}
            {actionResult.message}
          </div>
        </div>
      )}

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {/* Reboot */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-lg">
              <Power className="h-5 w-5 text-orange-500" />
              Reboot ONU
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground mb-4">
              Restart the ONU device. This will temporarily disconnect the customer.
            </p>
            <button
              onClick={() => {
                if (confirm(`Reboot ONU ${onuId} on PON ${pon}?`)) {
                  rebootMutation.mutate();
                }
              }}
              disabled={isLoading}
              className="w-full inline-flex items-center justify-center gap-2 px-4 py-2 text-sm font-medium rounded-md bg-orange-500 text-white hover:bg-orange-600 disabled:opacity-50"
            >
              <Power className="h-4 w-4" />
              {rebootMutation.isPending ? 'Rebooting...' : 'Reboot'}
            </button>
          </CardContent>
        </Card>

        {/* Block */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-lg">
              <Ban className="h-5 w-5 text-red-500" />
              Block ONU
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground mb-4">
              Disable the ONU. Customer will be disconnected until unblocked.
            </p>
            <button
              onClick={() => {
                if (confirm(`Block ONU ${onuId} on PON ${pon}? Customer will be disconnected.`)) {
                  blockMutation.mutate();
                }
              }}
              disabled={isLoading}
              className="w-full inline-flex items-center justify-center gap-2 px-4 py-2 text-sm font-medium rounded-md bg-red-500 text-white hover:bg-red-600 disabled:opacity-50"
            >
              <Ban className="h-4 w-4" />
              {blockMutation.isPending ? 'Blocking...' : 'Block'}
            </button>
          </CardContent>
        </Card>

        {/* Unblock */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-lg">
              <CheckCircle className="h-5 w-5 text-green-500" />
              Unblock ONU
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground mb-4">
              Re-enable a blocked ONU. Customer connection will be restored.
            </p>
            <button
              onClick={() => {
                if (confirm(`Unblock ONU ${onuId} on PON ${pon}?`)) {
                  unblockMutation.mutate();
                }
              }}
              disabled={isLoading}
              className="w-full inline-flex items-center justify-center gap-2 px-4 py-2 text-sm font-medium rounded-md bg-green-500 text-white hover:bg-green-600 disabled:opacity-50"
            >
              <CheckCircle className="h-4 w-4" />
              {unblockMutation.isPending ? 'Unblocking...' : 'Unblock'}
            </button>
          </CardContent>
        </Card>

        {/* Update Description */}
        <Card className="md:col-span-2">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-lg">
              <Edit className="h-5 w-5 text-blue-500" />
              Update Description
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex gap-4">
              <Input
                placeholder="e.g., Customer Name - Address"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                className="flex-1"
              />
              <button
                onClick={() => {
                  if (!description) {
                    alert('Please enter a description');
                    return;
                  }
                  updateDescMutation.mutate();
                }}
                disabled={isLoading || !description}
                className="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-md bg-blue-500 text-white hover:bg-blue-600 disabled:opacity-50"
              >
                <Edit className="h-4 w-4" />
                {updateDescMutation.isPending ? 'Updating...' : 'Update'}
              </button>
            </div>
          </CardContent>
        </Card>

        {/* Delete ONU */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-lg">
              <Trash2 className="h-5 w-5 text-red-500" />
              Delete ONU
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground mb-4">
              Permanently remove the ONU configuration. This cannot be undone.
            </p>
            <button
              onClick={() => {
                if (confirm(`⚠️ DELETE ONU ${onuId} on PON ${pon}?\n\nThis action cannot be undone!`)) {
                  deleteMutation.mutate();
                }
              }}
              disabled={isLoading}
              className="w-full inline-flex items-center justify-center gap-2 px-4 py-2 text-sm font-medium rounded-md bg-red-600 text-white hover:bg-red-700 disabled:opacity-50"
            >
              <Trash2 className="h-4 w-4" />
              {deleteMutation.isPending ? 'Deleting...' : 'Delete ONU'}
            </button>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
