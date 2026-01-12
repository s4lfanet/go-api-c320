import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { trafficApi, DBAProfile } from '@/api/endpoints/traffic';
import { RefreshCw, Gauge, Plus, Trash2 } from 'lucide-react';

export default function Traffic() {
  const queryClient = useQueryClient();
  const [showProfileForm, setShowProfileForm] = useState(false);
  const [editingProfile, setEditingProfile] = useState<DBAProfile | null>(null);
  const [profileForm, setProfileForm] = useState<DBAProfile>({
    name: '',
    type: 1,
    fixed_bandwidth: 0,
    assured_bandwidth: 10240,
    maximum_bandwidth: 102400,
  });

  // Get all DBA profiles
  const { data, isLoading, refetch } = useQuery({
    queryKey: ['dba-profiles'],
    queryFn: () => trafficApi.getAllDBAProfiles(),
  });

  // Create DBA profile mutation
  const createMutation = useMutation({
    mutationFn: trafficApi.createDBAProfile,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['dba-profiles'] });
      setShowProfileForm(false);
      resetForm();
      alert('DBA Profile created successfully!');
    },
    onError: (error: any) => {
      alert(`Failed to create profile: ${error.message}`);
    },
  });

  // Delete DBA profile mutation
  const deleteMutation = useMutation({
    mutationFn: trafficApi.deleteDBAProfile,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['dba-profiles'] });
      alert('DBA Profile deleted successfully!');
    },
    onError: (error: any) => {
      alert(`Failed to delete profile: ${error.message}`);
    },
  });

  const profiles = (data as any)?.data || [];

  const resetForm = () => {
    setProfileForm({
      name: '',
      type: 1,
      fixed_bandwidth: 0,
      assured_bandwidth: 10240,
      maximum_bandwidth: 102400,
    });
    setEditingProfile(null);
  };

  const handleSubmit = () => {
    if (!profileForm.name) {
      alert('Please enter a profile name');
      return;
    }
    createMutation.mutate(profileForm);
  };

  const formatBandwidth = (bw: number) => {
    if (bw >= 1024) {
      return `${(bw / 1024).toFixed(1)} Mbps`;
    }
    return `${bw} Kbps`;
  };

  const getTypeLabel = (type: number) => {
    switch (type) {
      case 1: return 'Type 1 (Fixed)';
      case 2: return 'Type 2 (Assured)';
      case 3: return 'Type 3 (Assured + Max)';
      case 4: return 'Type 4 (Max)';
      case 5: return 'Type 5 (Fixed + Assured + Max)';
      default: return `Type ${type}`;
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Traffic Control</h1>
          <p className="text-muted-foreground">
            Manage DBA profiles and bandwidth allocation
          </p>
        </div>
        <div className="flex gap-2">
          <button
            onClick={() => {
              resetForm();
              setShowProfileForm(true);
            }}
            className="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-md bg-primary text-primary-foreground hover:bg-primary/90"
          >
            <Plus className="h-4 w-4" />
            New Profile
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

      <div className="grid gap-6 lg:grid-cols-2">
        {/* DBA Profiles List */}
        <Card className={showProfileForm ? '' : 'lg:col-span-2'}>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Gauge className="h-5 w-5" />
              DBA Profiles
            </CardTitle>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="text-center py-8 text-muted-foreground">Loading profiles...</div>
            ) : profiles.length === 0 ? (
              <div className="text-center py-8">
                <Gauge className="h-12 w-12 mx-auto text-muted-foreground mb-2" />
                <p className="text-muted-foreground">No DBA profiles configured</p>
                <button
                  onClick={() => setShowProfileForm(true)}
                  className="mt-4 inline-flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-md bg-primary text-primary-foreground hover:bg-primary/90"
                >
                  <Plus className="h-4 w-4" />
                  Create Profile
                </button>
              </div>
            ) : (
              <div className="space-y-3">
                {profiles.map((profile: DBAProfile, idx: number) => (
                  <div
                    key={idx}
                    className="p-4 rounded-lg border border-border hover:bg-muted/50 transition-colors"
                  >
                    <div className="flex justify-between items-start">
                      <div>
                        <div className="font-medium">{profile.name}</div>
                        <div className="text-xs text-muted-foreground mt-1">
                          {getTypeLabel(profile.type)}
                        </div>
                      </div>
                      <button
                        onClick={() => {
                          if (confirm(`Delete profile "${profile.name}"?`)) {
                            deleteMutation.mutate(profile.name);
                          }
                        }}
                        disabled={deleteMutation.isPending}
                        className="text-red-500 hover:text-red-700 p-1"
                      >
                        <Trash2 className="h-4 w-4" />
                      </button>
                    </div>
                    <div className="grid grid-cols-3 gap-4 mt-3 text-sm">
                      <div>
                        <div className="text-xs text-muted-foreground">Fixed</div>
                        <div>{formatBandwidth(profile.fixed_bandwidth || 0)}</div>
                      </div>
                      <div>
                        <div className="text-xs text-muted-foreground">Assured</div>
                        <div>{formatBandwidth(profile.assured_bandwidth || 0)}</div>
                      </div>
                      <div>
                        <div className="text-xs text-muted-foreground">Maximum</div>
                        <div>{formatBandwidth(profile.maximum_bandwidth || 0)}</div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Create/Edit Profile Form */}
        {showProfileForm && (
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Plus className="h-5 w-5" />
                {editingProfile ? 'Edit' : 'Create'} DBA Profile
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div>
                  <label className="text-sm font-medium mb-1 block">Profile Name *</label>
                  <Input
                    placeholder="e.g., PROFILE-100M"
                    value={profileForm.name}
                    onChange={(e) => setProfileForm({ ...profileForm, name: e.target.value })}
                  />
                </div>
                <div>
                  <label className="text-sm font-medium mb-1 block">DBA Type</label>
                  <select
                    className="w-full h-10 rounded-md border border-input bg-background px-3 py-2 text-sm"
                    value={profileForm.type}
                    onChange={(e) => setProfileForm({ ...profileForm, type: parseInt(e.target.value) })}
                  >
                    <option value={1}>Type 1 - Fixed Bandwidth</option>
                    <option value={2}>Type 2 - Assured Bandwidth</option>
                    <option value={3}>Type 3 - Assured + Maximum</option>
                    <option value={4}>Type 4 - Maximum Only</option>
                    <option value={5}>Type 5 - Fixed + Assured + Maximum</option>
                  </select>
                </div>
                <div>
                  <label className="text-sm font-medium mb-1 block">Fixed Bandwidth (Kbps)</label>
                  <Input
                    type="number"
                    min={0}
                    value={profileForm.fixed_bandwidth}
                    onChange={(e) => setProfileForm({ ...profileForm, fixed_bandwidth: parseInt(e.target.value) })}
                  />
                </div>
                <div>
                  <label className="text-sm font-medium mb-1 block">Assured Bandwidth (Kbps)</label>
                  <Input
                    type="number"
                    min={0}
                    value={profileForm.assured_bandwidth}
                    onChange={(e) => setProfileForm({ ...profileForm, assured_bandwidth: parseInt(e.target.value) })}
                  />
                </div>
                <div>
                  <label className="text-sm font-medium mb-1 block">Maximum Bandwidth (Kbps)</label>
                  <Input
                    type="number"
                    min={0}
                    value={profileForm.maximum_bandwidth}
                    onChange={(e) => setProfileForm({ ...profileForm, maximum_bandwidth: parseInt(e.target.value) })}
                  />
                </div>
                <div className="text-xs text-muted-foreground">
                  <p>Common values: 10Mbps = 10240, 50Mbps = 51200, 100Mbps = 102400</p>
                </div>
                <div className="flex gap-2 pt-2">
                  <button
                    onClick={handleSubmit}
                    disabled={createMutation.isPending}
                    className="flex-1 inline-flex items-center justify-center gap-2 px-4 py-2 text-sm font-medium rounded-md bg-primary text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
                  >
                    {createMutation.isPending ? 'Creating...' : 'Create Profile'}
                  </button>
                  <button
                    onClick={() => {
                      setShowProfileForm(false);
                      resetForm();
                    }}
                    className="px-4 py-2 text-sm font-medium rounded-md border border-input bg-background hover:bg-accent"
                  >
                    Cancel
                  </button>
                </div>
              </div>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}
