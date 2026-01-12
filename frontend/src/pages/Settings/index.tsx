import { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Server, Moon, Sun, Save, RotateCcw } from 'lucide-react';
import { useThemeStore } from '@/store/themeStore';

interface AppSettings {
  apiBaseUrl: string;
  refreshInterval: number;
  theme: 'light' | 'dark';
}

const defaultSettings: AppSettings = {
  apiBaseUrl: '/api/v1',
  refreshInterval: 30,
  theme: 'light',
};

export default function Settings() {
  const { theme, setTheme } = useThemeStore();
  const [settings, setSettings] = useState<AppSettings>(() => {
    const saved = localStorage.getItem('app-settings');
    return saved ? { ...defaultSettings, ...JSON.parse(saved) } : defaultSettings;
  });
  const [saved, setSaved] = useState(false);

  useEffect(() => {
    setSettings(prev => ({ ...prev, theme }));
  }, [theme]);

  const handleSave = () => {
    localStorage.setItem('app-settings', JSON.stringify(settings));
    setTheme(settings.theme);
    setSaved(true);
    setTimeout(() => setSaved(false), 2000);
  };

  const handleReset = () => {
    if (confirm('Reset all settings to default?')) {
      setSettings(defaultSettings);
      localStorage.removeItem('app-settings');
      setTheme('light');
      setSaved(true);
      setTimeout(() => setSaved(false), 2000);
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Settings</h1>
        <p className="text-muted-foreground">
          Configure application settings and preferences
        </p>
      </div>

      {/* API Settings */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Server className="h-5 w-5" />
            API Configuration
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div>
            <label className="text-sm font-medium mb-1 block">API Base URL</label>
            <Input
              value={settings.apiBaseUrl}
              onChange={(e) => setSettings({ ...settings, apiBaseUrl: e.target.value })}
              placeholder="/api/v1"
            />
            <p className="text-xs text-muted-foreground mt-1">
              The base URL for API requests. Default: /api/v1 (proxied through Nginx)
            </p>
          </div>

          <div>
            <label className="text-sm font-medium mb-1 block">Auto-refresh Interval (seconds)</label>
            <Input
              type="number"
              min={5}
              max={300}
              value={settings.refreshInterval}
              onChange={(e) => setSettings({ ...settings, refreshInterval: parseInt(e.target.value) })}
            />
            <p className="text-xs text-muted-foreground mt-1">
              How often to refresh data automatically. Set to 0 to disable.
            </p>
          </div>
        </CardContent>
      </Card>

      {/* Theme Settings */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Sun className="h-5 w-5" />
            Appearance
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-3 gap-4">
            <button
              onClick={() => setSettings({ ...settings, theme: 'light' })}
              className={`p-4 rounded-lg border transition-colors ${
                settings.theme === 'light'
                  ? 'border-primary bg-primary/10'
                  : 'border-border hover:bg-muted'
              }`}
            >
              <Sun className="h-6 w-6 mx-auto mb-2" />
              <div className="text-sm font-medium">Light</div>
            </button>

            <button
              onClick={() => setSettings({ ...settings, theme: 'dark' })}
              className={`p-4 rounded-lg border transition-colors ${
                settings.theme === 'dark'
                  ? 'border-primary bg-primary/10'
                  : 'border-border hover:bg-muted'
              }`}
            >
              <Moon className="h-6 w-6 mx-auto mb-2" />
              <div className="text-sm font-medium">Dark</div>
            </button>
          </div>
        </CardContent>
      </Card>

      {/* About */}
      <Card>
        <CardHeader>
          <CardTitle>About</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-2 text-sm">
            <div className="flex justify-between">
              <span className="text-muted-foreground">Application</span>
              <span className="font-medium">ZTE C320 OLT Dashboard</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Version</span>
              <span className="font-medium">1.0.0</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Backend API</span>
              <span className="font-medium">Go SNMP OLT API</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Supported OLT</span>
              <span className="font-medium">ZTE C320 v2.1.0</span>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Save/Reset Buttons */}
      <div className="flex gap-4">
        <button
          onClick={handleSave}
          className="inline-flex items-center gap-2 px-6 py-2 text-sm font-medium rounded-md bg-primary text-primary-foreground hover:bg-primary/90"
        >
          <Save className="h-4 w-4" />
          {saved ? 'Saved!' : 'Save Settings'}
        </button>
        <button
          onClick={handleReset}
          className="inline-flex items-center gap-2 px-6 py-2 text-sm font-medium rounded-md border border-input bg-background hover:bg-accent"
        >
          <RotateCcw className="h-4 w-4" />
          Reset to Default
        </button>
      </div>
    </div>
  );
}
