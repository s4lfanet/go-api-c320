import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { provisioningApi } from '@/api/endpoints/provisioning';
import { RefreshCw, Plus, Search, AlertCircle } from 'lucide-react';

export default function Provisioning() {
  const queryClient = useQueryClient();
  const [selectedONU, setSelectedONU] = useState<any>(null);
  const [registerForm, setRegisterForm] = useState({
    pon_port: 1,
    onu_id: 1,
    name: '',
    description: '',
  });

  // Get unconfigured ONUs
  const { data, isLoading, refetch } = useQuery({
    queryKey: ['unconfigured-onus'],
    queryFn: () => provisioningApi.getUnconfiguredONUs(),
    refetchInterval: 10000,
  });

  // Register ONU mutation
  const registerMutation = useMutation({
    mutationFn: provisioningApi.registerONU,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['unconfigured-onus'] });
      setSelectedONU(null);
      setRegisterForm({ pon_port: 1, onu_id: 1, name: '', description: '' });
      alert('ONU registered successfully!');
    },
    onError: (error: any) => {
      alert(`Failed to register ONU: ${error.message}`);
    },
  });

  const unconfiguredONUs = (data as any)?.data || [];

  const handleRegister = () => {
    if (!selectedONU) {
      alert('Please select an unconfigured ONU first');
      return;
    }
    if (!registerForm.name) {
      alert('Please enter a name for the ONU');
      return;
    }

    registerMutation.mutate({
      pon_port: parseInt(selectedONU.pon_port),
      onu_id: registerForm.onu_id,
      serial_number: selectedONU.serial_number,
      onu_type: selectedONU.onu_type || 'ZTE-F660',
      name: registerForm.name,
      description: registerForm.description,
    });
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">ONU Provisioning</h1>
          <p className="text-muted-foreground">
            Discover and register new ONUs
          </p>
        </div>
        <button
          onClick={() => refetch()}
          disabled={isLoading}
          className="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-md border border-input bg-background hover:bg-accent hover:text-accent-foreground disabled:opacity-50"
        >
          <RefreshCw className={`h-4 w-4 ${isLoading ? 'animate-spin' : ''}`} />
          Scan
        </button>
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        {/* Unconfigured ONUs */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Search className="h-5 w-5" />
              Unconfigured ONUs
            </CardTitle>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="text-center py-8 text-muted-foreground">
                Scanning for unconfigured ONUs...
              </div>
            ) : unconfiguredONUs.length === 0 ? (
              <div className="text-center py-8">
                <AlertCircle className="h-12 w-12 mx-auto text-muted-foreground mb-2" />
                <p className="text-muted-foreground">No unconfigured ONUs found</p>
                <p className="text-xs text-muted-foreground mt-1">
                  Click Scan to check for new ONUs
                </p>
              </div>
            ) : (
              <div className="space-y-2">
                {unconfiguredONUs.map((onu: any, idx: number) => (
                  <div
                    key={idx}
                    onClick={() => setSelectedONU(onu)}
                    className={`p-3 rounded-lg border cursor-pointer transition-colors ${
                      selectedONU?.serial_number === onu.serial_number
                        ? 'border-primary bg-primary/10'
                        : 'border-border hover:bg-muted'
                    }`}
                  >
                    <div className="flex justify-between items-start">
                      <div>
                        <div className="font-mono text-sm font-medium">
                          {onu.serial_number}
                        </div>
                        <div className="text-xs text-muted-foreground">
                          PON: {onu.pon_port} | Type: {onu.onu_type || 'Unknown'}
                        </div>
                      </div>
                      <span className="text-xs bg-yellow-100 text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300 px-2 py-1 rounded">
                        Unconfigured
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Register Form */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Plus className="h-5 w-5" />
              Register ONU
            </CardTitle>
          </CardHeader>
          <CardContent>
            {selectedONU ? (
              <div className="space-y-4">
                <div className="p-3 bg-muted rounded-lg">
                  <div className="text-sm font-medium">Selected ONU</div>
                  <div className="font-mono text-lg">{selectedONU.serial_number}</div>
                  <div className="text-xs text-muted-foreground">
                    PON Port: {selectedONU.pon_port}
                  </div>
                </div>

                <div className="space-y-3">
                  <div>
                    <label className="text-sm font-medium mb-1 block">ONU ID</label>
                    <Input
                      type="number"
                      min={1}
                      max={128}
                      value={registerForm.onu_id}
                      onChange={(e) => setRegisterForm({ ...registerForm, onu_id: parseInt(e.target.value) })}
                    />
                  </div>
                  <div>
                    <label className="text-sm font-medium mb-1 block">Name *</label>
                    <Input
                      placeholder="e.g., Customer-001"
                      value={registerForm.name}
                      onChange={(e) => setRegisterForm({ ...registerForm, name: e.target.value })}
                    />
                  </div>
                  <div>
                    <label className="text-sm font-medium mb-1 block">Description</label>
                    <Input
                      placeholder="e.g., Jl. Contoh No. 123"
                      value={registerForm.description}
                      onChange={(e) => setRegisterForm({ ...registerForm, description: e.target.value })}
                    />
                  </div>
                </div>

                <div className="flex gap-2 pt-2">
                  <button
                    onClick={handleRegister}
                    disabled={registerMutation.isPending}
                    className="flex-1 inline-flex items-center justify-center gap-2 px-4 py-2 text-sm font-medium rounded-md bg-primary text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
                  >
                    <Plus className="h-4 w-4" />
                    {registerMutation.isPending ? 'Registering...' : 'Register ONU'}
                  </button>
                  <button
                    onClick={() => setSelectedONU(null)}
                    className="px-4 py-2 text-sm font-medium rounded-md border border-input bg-background hover:bg-accent"
                  >
                    Cancel
                  </button>
                </div>
              </div>
            ) : (
              <div className="text-center py-8 text-muted-foreground">
                <Plus className="h-12 w-12 mx-auto mb-2 opacity-50" />
                <p>Select an unconfigured ONU to register</p>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
