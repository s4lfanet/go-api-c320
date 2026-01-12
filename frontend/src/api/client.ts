import axios, { AxiosError } from 'axios';
import { API_BASE_URL } from '@/utils/constants';

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor
apiClient.interceptors.request.use(
  (config) => {
    // Add auth token if available
    const token = localStorage.getItem('auth-token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor
apiClient.interceptors.response.use(
  (response) => {
    // Return the data property from the response
    return response.data;
  },
  (error: AxiosError) => {
    // Handle different error scenarios
    if (error.response?.status === 401) {
      // Unauthorized - clear auth and redirect to login
      localStorage.removeItem('auth-token');
      window.location.href = '/login';
    }

    // Return a structured error
    return Promise.reject({
      message: error.message,
      status: error.response?.status,
      data: error.response?.data,
    });
  }
);
