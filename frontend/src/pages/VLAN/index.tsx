import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { vlanApi } from '@/api/endpoints/vlan';
import { RefreshCw, Network, Plus, Trash2, Edit } from 'lucide-react';

export default function VLAN() {
  const queryClient = useQueryClient();
  const [pon, setPon] = useState(1);
  const [onuId, setOnuId] = useState(1);
  const [showConfigForm, setShowConfigForm] = useState(false);
  const [vlanConfig, setVlanConfig] = useState({
    vlan_id: 100,
    svlan: 0,
    cvlan: 0,
    vlan_mode: 'tag',
    cos: 0,
  });

  // Get ONU VLAN
  const { data: vlanData, isLoading, refetch } = useQuery({
    queryKey: ['onu-vlan', pon, onuId],
    queryFn: () => vlanApi.getONUVLAN(pon, onuId),
    enabled: pon > 0 && onuId > 0,
  });

  // Get all service ports
  const { data: servicePortsData } = useQuery({
    queryKey: ['service-ports'],
    queryFn: () => vlanApi.getAllServicePorts(),
  });

  // Configure VLAN mutation
  const configureMutation = useMutation({
    mutationFn: () => vlanApi.configureVLAN({
      pon_port: pon,
      onu_id: onuId,
      ...vlanConfig,
    }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['onu-vlan'] });
      setShowConfigForm(false);
      alert('VLAN configured successfully!');
    },
    onError: (error: any) => {
      alert(`Failed to configure VLAN: ${error.message}`);
    },
  });

  // Delete VLAN mutation
  const deleteMutation = useMutation({
    mutationFn: () => vlanApi.deleteVLAN(pon, onuId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['onu-vlan'] });
      alert('VLAN deleted successfully!');
    },
    onError: (error: any) => {
      alert(`Failed to delete VLAN: ${error.message}`);
    },
  });

  const vlanInfo = (vlanData as any)?.data || null;
  const servicePorts = (servicePortsData as any)?.data || [];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">VLAN Management</h1>
          <p className="text-muted-foreground">
            Configure VLAN settings for ONUs
          </p>
        </div>
      </div>

      {/* ONU Selector */}
      <Card>
        <CardHeader>
          <CardTitle>Select ONU</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex gap-4 items-end">
            <div className="flex-1">
              <label className="text-sm font-medium mb-1 block">PON Port</label>
              <Input
                type="number"
                min={1}
                max={16}
                value={pon}
                onChange={(e) => setPon(parseInt(e.target.value))}
              />
            </div>
            <div className="flex-1">
              <label className="text-sm font-medium mb-1 block">ONU ID</label>
              <Input
                type="number"
                min={1}
                max={128}
                value={onuId}
                onChange={(e) => setOnuId(parseInt(e.target.value))}
              />
            </div>
            <button
              onClick={() => refetch()}
              disabled={isLoading}
              className="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-md border border-input bg-background hover:bg-accent hover:text-accent-foreground disabled:opacity-50"
            >
              <RefreshCw className={`h-4 w-4 ${isLoading ? 'animate-spin' : ''}`} />
              Query
            </button>
          </div>
        </CardContent>
      </Card>

      <div className="grid gap-6 lg:grid-cols-2">
        {/* Current VLAN Config */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Network className="h-5 w-5" />
              Current VLAN Configuration
            </CardTitle>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="text-center py-8 text-muted-foreground">Loading...</div>
            ) : vlanInfo ? (
              <div className="space-y-3">
                <div className="grid grid-cols-2 gap-4 p-4 bg-muted rounded-lg">
                  <div>
                    <div className="text-xs text-muted-foreground">VLAN ID</div>
                    <div className="font-medium">{vlanInfo.vlan_id || '-'}</div>
                  </div>
                  <div>
                    <div className="text-xs text-muted-foreground">Service Port</div>
                    <div className="font-medium">{vlanInfo.service_port_id || '-'}</div>
                  </div>
                  <div>
                    <div className="text-xs text-muted-foreground">SVLAN</div>
                    <div className="font-medium">{vlanInfo.svlan || '-'}</div>
                  </div>
                  <div>
                    <div className="text-xs text-muted-foreground">CVLAN</div>
                    <div className="font-medium">{vlanInfo.cvlan || '-'}</div>
                  </div>
                  <div>
                    <div className="text-xs text-muted-foreground">Mode</div>
                    <div className="font-medium">{vlanInfo.vlan_mode || 'tag'}</div>
                  </div>
                  <div>
                    <div className="text-xs text-muted-foreground">CoS</div>
                    <div className="font-medium">{vlanInfo.cos || 0}</div>
                  </div>
                </div>
                <div className="flex gap-2">
                  <button
                    onClick={() => setShowConfigForm(true)}
                    className="flex-1 inline-flex items-center justify-center gap-2 px-4 py-2 text-sm font-medium rounded-md border border-input bg-background hover:bg-accent"
                  >
                    <Edit className="h-4 w-4" />
                    Modify
                  </button>
                  <button
                    onClick={() => {
                      if (confirm('Are you sure you want to delete this VLAN configuration?')) {
                        deleteMutation.mutate();
                      }
                    }}
                    disabled={deleteMutation.isPending}
                    className="inline-flex items-center justify-center gap-2 px-4 py-2 text-sm font-medium rounded-md bg-red-500 text-white hover:bg-red-600 disabled:opacity-50"
                  >
                    <Trash2 className="h-4 w-4" />
                  </button>
                </div>
              </div>
            ) : (
              <div className="text-center py-8">
                <Network className="h-12 w-12 mx-auto text-muted-foreground mb-2" />
                <p className="text-muted-foreground">No VLAN configured</p>
                <button
                  onClick={() => setShowConfigForm(true)}
                  className="mt-4 inline-flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-md bg-primary text-primary-foreground hover:bg-primary/90"
                >
                  <Plus className="h-4 w-4" />
                  Configure VLAN
                </button>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Configure Form */}
        {showConfigForm && (
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Plus className="h-5 w-5" />
                {vlanInfo ? 'Modify' : 'Configure'} VLAN
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div>
                  <label className="text-sm font-medium mb-1 block">VLAN ID *</label>
                  <Input
                    type="number"
                    min={1}
                    max={4094}
                    value={vlanConfig.vlan_id}
                    onChange={(e) => setVlanConfig({ ...vlanConfig, vlan_id: parseInt(e.target.value) })}
                  />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="text-sm font-medium mb-1 block">SVLAN</label>
                    <Input
                      type="number"
                      min={0}
                      max={4094}
                      value={vlanConfig.svlan}
                      onChange={(e) => setVlanConfig({ ...vlanConfig, svlan: parseInt(e.target.value) })}
                    />
                  </div>
                  <div>
                    <label className="text-sm font-medium mb-1 block">CVLAN</label>
                    <Input
                      type="number"
                      min={0}
                      max={4094}
                      value={vlanConfig.cvlan}
                      onChange={(e) => setVlanConfig({ ...vlanConfig, cvlan: parseInt(e.target.value) })}
                    />
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="text-sm font-medium mb-1 block">Mode</label>
                    <select
                      className="w-full h-10 rounded-md border border-input bg-background px-3 py-2 text-sm"
                      value={vlanConfig.vlan_mode}
                      onChange={(e) => setVlanConfig({ ...vlanConfig, vlan_mode: e.target.value })}
                    >
                      <option value="tag">Tag</option>
                      <option value="untag">Untag</option>
                      <option value="transparent">Transparent</option>
                    </select>
                  </div>
                  <div>
                    <label className="text-sm font-medium mb-1 block">CoS (0-7)</label>
                    <Input
                      type="number"
                      min={0}
                      max={7}
                      value={vlanConfig.cos}
                      onChange={(e) => setVlanConfig({ ...vlanConfig, cos: parseInt(e.target.value) })}
                    />
                  </div>
                </div>
                <div className="flex gap-2 pt-2">
                  <button
                    onClick={() => configureMutation.mutate()}
                    disabled={configureMutation.isPending}
                    className="flex-1 inline-flex items-center justify-center gap-2 px-4 py-2 text-sm font-medium rounded-md bg-primary text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
                  >
                    {configureMutation.isPending ? 'Saving...' : 'Save Configuration'}
                  </button>
                  <button
                    onClick={() => setShowConfigForm(false)}
                    className="px-4 py-2 text-sm font-medium rounded-md border border-input bg-background hover:bg-accent"
                  >
                    Cancel
                  </button>
                </div>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Service Ports List */}
        <Card className={showConfigForm ? 'lg:col-span-2' : ''}>
          <CardHeader>
            <CardTitle>Service Ports</CardTitle>
          </CardHeader>
          <CardContent>
            {servicePorts.length === 0 ? (
              <div className="text-center py-4 text-muted-foreground">
                No service ports configured
              </div>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b">
                      <th className="text-left p-2">ID</th>
                      <th className="text-left p-2">VLAN</th>
                      <th className="text-left p-2">PON</th>
                      <th className="text-left p-2">ONU</th>
                      <th className="text-left p-2">Mode</th>
                    </tr>
                  </thead>
                  <tbody>
                    {servicePorts.slice(0, 10).map((sp: any, idx: number) => (
                      <tr key={idx} className="border-b hover:bg-muted/50">
                        <td className="p-2">{sp.service_port_id}</td>
                        <td className="p-2">{sp.vlan_id}</td>
                        <td className="p-2">{sp.pon_port}</td>
                        <td className="p-2">{sp.onu_id}</td>
                        <td className="p-2">{sp.vlan_mode || 'tag'}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
                {servicePorts.length > 10 && (
                  <div className="text-xs text-muted-foreground mt-2">
                    Showing 10 of {servicePorts.length} service ports
                  </div>
                )}
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
