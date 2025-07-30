import axios, { AxiosError } from "axios";

let accessToken: string | null = null;
export const setAccessToken = (token: string | null) => { accessToken = token; };
export const getAccessToken = () => accessToken;

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  withCredentials: true,
});

api.interceptors.request.use((config) => {
  if (accessToken) {
    config.headers.Authorization = `Bearer ${accessToken}`;
  }
  return config;
});

api.interceptors.response.use(
  (res) => res,
  async (err: AxiosError) => {
    const originalConfig = err.config!;
    if (err.response?.status === 401 && !originalConfig._retry) {
      originalConfig._retry = true;
      try {
        const refreshRes = await axios.post(`${import.meta.env.VITE_API_URL}/refresh`, {}, { withCredentials: true });
        setAccessToken(refreshRes.data.data.access_token);
        originalConfig.headers.Authorization = `Bearer ${refreshRes.data.data.access_token}`;
        return api(originalConfig);
      } catch (refreshError) {
        window.location.href = "/login";
        return Promise.reject(refreshError);
      }
    }
    return Promise.reject(err);
  }
);

export default api;
