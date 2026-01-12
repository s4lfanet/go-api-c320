import { format, formatDistanceToNow } from 'date-fns';

export const formatDate = (date: string | Date): string => {
  try {
    return format(new Date(date), 'yyyy-MM-dd HH:mm:ss');
  } catch {
    return 'Invalid date';
  }
};

export const formatRelativeTime = (date: string | Date): string => {
  try {
    return formatDistanceToNow(new Date(date), { addSuffix: true });
  } catch {
    return 'Unknown';
  }
};

export const formatPower = (power: string | number): string => {
  const powerNum = typeof power === 'string' ? parseFloat(power) : power;
  if (isNaN(powerNum)) return 'N/A';
  return `${powerNum.toFixed(2)} dBm`;
};

export const formatBytes = (bytes: number): string => {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
};

export const getSignalQuality = (rxPower: string | number): {
  label: string;
  color: string;
} => {
  const power = typeof rxPower === 'string' ? parseFloat(rxPower) : rxPower;
  if (isNaN(power)) return { label: 'Unknown', color: 'gray' };

  if (power >= -20) return { label: 'Excellent', color: 'green' };
  if (power >= -23) return { label: 'Good', color: 'blue' };
  if (power >= -25) return { label: 'Fair', color: 'yellow' };
  if (power >= -27) return { label: 'Poor', color: 'orange' };
  return { label: 'Critical', color: 'red' };
};
