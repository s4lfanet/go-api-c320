import { Activity, Users, Signal, AlertTriangle, RefreshCw } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { useQuery } from '@tanstack/react-query';
import { onuApi } from '@/api/endpoints/onu';

export default function Dashboard() {
  const { data: stats, isLoading, error, refetch } = useQuery({
    queryKey: ['dashboard-stats'],
    queryFn: () => onuApi.getDashboardStats(),
    refetchInterval: 30000, // Refresh every 30 seconds
  });

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Dashboard</h1>
          <p className="text-muted-foreground">
            Overview of your ZTE C320 OLT network status
          </p>
        </div>
        <button
          onClick={() => refetch()}
          disabled={isLoading}
          className="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-md border border-input bg-background hover:bg-accent hover:text-accent-foreground disabled:opacity-50"
        >
          <RefreshCw className={`h-4 w-4 ${isLoading ? 'animate-spin' : ''}`} />
          Refresh
        </button>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total ONUs</CardTitle>
            <Users className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="text-2xl font-bold animate-pulse">...</div>
            ) : error ? (
              <div className="text-2xl font-bold text-red-500">Error</div>
            ) : (
              <>
                <div className="text-2xl font-bold">{stats?.totalOnus || 0}</div>
                <p className="text-xs text-muted-foreground">
                  Across all boards and PON ports
                </p>
              </>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Online</CardTitle>
            <Activity className="h-4 w-4 text-green-500" />
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="text-2xl font-bold animate-pulse">...</div>
            ) : error ? (
              <div className="text-2xl font-bold text-red-500">Error</div>
            ) : (
              <>
                <div className="text-2xl font-bold text-green-500">{stats?.onlineOnus || 0}</div>
                <p className="text-xs text-muted-foreground">{stats?.uptime || 0}% uptime</p>
              </>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Offline</CardTitle>
            <Signal className="h-4 w-4 text-red-500" />
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="text-2xl font-bold animate-pulse">...</div>
            ) : error ? (
              <div className="text-2xl font-bold text-red-500">Error</div>
            ) : (
              <>
                <div className="text-2xl font-bold text-red-500">{stats?.offlineOnus || 0}</div>
                <p className="text-xs text-muted-foreground">
                  {stats?.totalOnus ? ((stats.offlineOnus / stats.totalOnus) * 100).toFixed(1) : 0}% offline
                </p>
              </>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Alerts</CardTitle>
            <AlertTriangle className="h-4 w-4 text-yellow-500" />
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="text-2xl font-bold animate-pulse">...</div>
            ) : error ? (
              <div className="text-2xl font-bold text-red-500">Error</div>
            ) : (
              <>
                <div className="text-2xl font-bold text-yellow-500">{stats?.alerts || 0}</div>
                <p className="text-xs text-muted-foreground">Requires attention</p>
              </>
            )}
          </CardContent>
        </Card>
      </div>

      {/* PON Ports Summary */}
      <Card>
        <CardHeader>
          <CardTitle>Active PON Ports</CardTitle>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="text-sm text-muted-foreground">Loading PON ports...</div>
          ) : error ? (
            <div className="text-sm text-red-500">Failed to load PON port data</div>
          ) : stats?.ponPorts && stats.ponPorts.length > 0 ? (
            <div className="space-y-2">
              <div className="grid grid-cols-7 gap-2 text-xs font-medium text-muted-foreground border-b pb-2">
                <div>Board</div>
                <div>PON</div>
                <div>Total</div>
                <div>Online</div>
                <div>Offline</div>
                <div className="col-span-2">Status</div>
              </div>
              <div className="max-h-96 overflow-y-auto space-y-1">
                {stats.ponPorts.map(({ board, pon, total, online, offline }) => (
                  <div key={`${board}-${pon}`} className="grid grid-cols-7 gap-2 text-sm py-2 hover:bg-muted/50 rounded">
                    <div className="font-medium">{board}</div>
                    <div className="font-medium">{pon}</div>
                    <div>{total}</div>
                    <div className="text-green-600">{online}</div>
                    <div className="text-red-600">{offline}</div>
                    <div className="col-span-2">
                      <div className="flex items-center gap-1">
                        <div className="flex-1 bg-muted rounded-full h-2 overflow-hidden">
                          <div
                            className="bg-green-500 h-full rounded-full"
                            style={{ width: `${(online / total) * 100}%` }}
                          />
                        </div>
                        <span className="text-xs text-muted-foreground">
                          {((online / total) * 100).toFixed(0)}%
                        </span>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          ) : (
            <div className="text-sm text-muted-foreground">
              No active PON ports found. Please check your OLT connection.
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
