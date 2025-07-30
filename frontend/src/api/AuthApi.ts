import api, { setAccessToken } from "../api";

export interface LoginResponse {
  access_token: string;
}

export const authApi = {
  async login(username: string, password: string): Promise<LoginResponse> {
    const res = await api.post("/login", { username, password });
    setAccessToken(res.data.data.access_token);
    return res.data.data;
  },

  async refresh(): Promise<LoginResponse> {
    const res = await api.post("/refresh", {}, { withCredentials: true });
    setAccessToken(res.data.data.access_token);
    return res.data.data;
  },

  logout(): void {
    setAccessToken(null);
    window.location.href = "/login";
  },
};
