import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { authApi } from "../api/AuthApi";
import type { AxiosError } from "axios";
import type { ApiError } from "../types/Errors";

export default function Login() {
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState<string | null>(null);
    const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});
    const navigate = useNavigate();

    const handleLogin = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        setFieldErrors({});
        try {
            await authApi.login(username, password);
            navigate("/dashboard");
        } catch (err) {
            const axiosErr = err as AxiosError<ApiError>;
            setError(axiosErr.response?.data?.error || "Login failed");
            if (axiosErr.response?.data?.fields) {
                setFieldErrors(axiosErr.response.data.fields);
            }
        }
    };

    return (
        <div className="min-h-screen bg-gray-900 flex items-center justify-center">
            <div className="bg-gray-800 p-8 rounded shadow-md w-full max-w-md">
                <h2 className="text-2xl text-white font-bold mb-6 text-center">Login</h2>
                {error && <p className="text-red-400 text-center mb-4">{error}</p>}
                <form onSubmit={handleLogin} className="space-y-4">
                    <div>
                        <label className="block text-gray-300">Username</label>
                        <input
                            type="text"
                            className="w-full bg-gray-700 border border-gray-600 text-white px-3 py-2 rounded focus:outline-none focus:border-blue-500"
                            value={username}
                            onChange={(e) => setUsername(e.target.value)}
                            required
                        />
                        {fieldErrors.username && (
                            <p className="text-red-400 text-sm mt-1">{fieldErrors.username}</p>
                        )}
                    </div>

                    <div>
                        <label className="block text-gray-300">Password</label>
                        <input
                            type="password"
                            className="w-full bg-gray-700 border border-gray-600 text-white px-3 py-2 rounded focus:outline-none focus:border-blue-500"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            required
                        />
                        {fieldErrors.password && (
                            <p className="text-red-400 text-sm mt-1">{fieldErrors.password}</p>
                        )}
                    </div>
                    <button
                        type="submit"
                        className="w-full bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded"
                    >
                        Login
                    </button>
                </form>
            </div>
        </div>
    );
}
