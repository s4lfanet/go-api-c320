import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { useQuery } from '@tanstack/react-query';
import { onuApi } from '@/api/endpoints/onu';
import { RefreshCw, Activity, Signal, Wifi, Server } from 'lucide-react';

export default function Monitoring() {
  const [selectedPon, setSelectedPon] = useState<number | null>(null);

  // Get OLT summary
  const { data: oltData, isLoading: oltLoading, refetch: refetchOlt } = useQuery({
    queryKey: ['olt-monitoring'],
    queryFn: () => onuApi.getOLTMonitoring(),
    refetchInterval: 15000,
  });

  // Get PON monitoring when selected
  const { data: ponData, isLoading: ponLoading } = useQuery({
    queryKey: ['pon-monitoring', selectedPon],
    queryFn: () => selectedPon ? onuApi.getPONMonitoring(selectedPon) : null,
    enabled: !!selectedPon,
    refetchInterval: 10000,
  });

  const oltStats = (oltData as any)?.data || oltData || {};
  const ponPorts = oltStats.pon_ports || [];
  
  // Remove duplicate PON ports
  const uniquePonPorts = Array.from(
    new Map(ponPorts.map((p: any) => [p.pon_port, p])).values()
  ).sort((a: any, b: any) => parseInt(a.pon_port) - parseInt(b.pon_port));

  const ponMonitoringData = (ponData as any)?.data || ponData || {};

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">ONU Monitoring</h1>
          <p className="text-muted-foreground">
            Real-time monitoring of all ONUs across PON ports
          </p>
        </div>
        <button
          onClick={() => refetchOlt()}
          disabled={oltLoading}
          className="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-md border border-input bg-background hover:bg-accent hover:text-accent-foreground disabled:opacity-50"
        >
          <RefreshCw className={`h-4 w-4 ${oltLoading ? 'animate-spin' : ''}`} />
          Refresh
        </button>
      </div>

      {/* OLT Summary */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total ONUs</CardTitle>
            <Server className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{oltStats.total_onus || 0}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Online</CardTitle>
            <Activity className="h-4 w-4 text-green-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-500">{oltStats.online_onus || 0}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Offline</CardTitle>
            <Signal className="h-4 w-4 text-red-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-500">{oltStats.offline_onus || 0}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">PON Ports</CardTitle>
            <Wifi className="h-4 w-4 text-blue-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-blue-500">{uniquePonPorts.length}</div>
          </CardContent>
        </Card>
      </div>

      {/* PON Ports Grid */}
      <Card>
        <CardHeader>
          <CardTitle>PON Ports Status</CardTitle>
        </CardHeader>
        <CardContent>
          {oltLoading ? (
            <div className="text-center py-4 text-muted-foreground">Loading PON ports...</div>
          ) : (
            <div className="grid grid-cols-4 md:grid-cols-8 lg:grid-cols-16 gap-2">
              {uniquePonPorts.map((pon: any) => (
                <button
                  key={pon.pon_port}
                  onClick={() => setSelectedPon(parseInt(pon.pon_port))}
                  className={`p-3 rounded-lg border text-center transition-colors ${
                    selectedPon === parseInt(pon.pon_port)
                      ? 'border-primary bg-primary/10'
                      : 'border-border hover:bg-muted'
                  }`}
                >
                  <div className="text-xs text-muted-foreground">PON</div>
                  <div className="text-lg font-bold">{pon.pon_port}</div>
                  <div className="text-xs">
                    <span className="text-green-500">{pon.online_count}</span>
                    /
                    <span className="text-muted-foreground">{pon.onu_count}</span>
                  </div>
                </button>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Selected PON Details */}
      {selectedPon && (
        <Card>
          <CardHeader>
            <CardTitle>
              PON {selectedPon} Details
              {ponLoading && <span className="text-sm font-normal ml-2">Loading...</span>}
            </CardTitle>
          </CardHeader>
          <CardContent>
            {ponMonitoringData.onus && ponMonitoringData.onus.length > 0 ? (
              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b">
                      <th className="text-left p-2">ONU ID</th>
                      <th className="text-left p-2">Name</th>
                      <th className="text-left p-2">Serial Number</th>
                      <th className="text-left p-2">Type</th>
                      <th className="text-left p-2">Status</th>
                      <th className="text-left p-2">RX Power</th>
                      <th className="text-left p-2">TX Power</th>
                      <th className="text-left p-2">Distance</th>
                    </tr>
                  </thead>
                  <tbody>
                    {ponMonitoringData.onus.map((onu: any) => (
                      <tr key={onu.onu_id} className="border-b hover:bg-muted/50">
                        <td className="p-2 font-medium">{onu.onu_id}</td>
                        <td className="p-2">{onu.name || '-'}</td>
                        <td className="p-2 font-mono text-xs">{onu.serial_number}</td>
                        <td className="p-2">{onu.onu_type || '-'}</td>
                        <td className="p-2">
                          <span
                            className={`inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium ${
                              onu.status === 'online'
                                ? 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300'
                                : 'bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300'
                            }`}
                          >
                            {onu.status === 'online' ? (
                              <><Activity className="h-3 w-3" /> Online</>
                            ) : (
                              <><Signal className="h-3 w-3" /> Offline</>
                            )}
                          </span>
                        </td>
                        <td className="p-2">
                          <span className={`${
                            parseFloat(onu.rx_power) < -27 ? 'text-red-600' :
                            parseFloat(onu.rx_power) < -25 ? 'text-yellow-600' : 'text-green-600'
                          }`}>
                            {onu.rx_power || '-'} dBm
                          </span>
                        </td>
                        <td className="p-2">{onu.tx_power || '-'} dBm</td>
                        <td className="p-2">{onu.distance || '-'} m</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            ) : (
              <div className="text-center py-8 text-muted-foreground">
                {ponLoading ? 'Loading...' : 'No ONUs found on this PON port'}
              </div>
            )}
          </CardContent>
        </Card>
      )}
    </div>
  );
}
