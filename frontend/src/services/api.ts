import { getValidAccessToken, logout } from "./auth-service";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "";

interface ApiResponse<T = any> {
  success: boolean;
  message: string;
  data?: T;
  error_code?: string;
}

export async function apiFetch<T = any>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const accessToken = await getValidAccessToken();
  if (!accessToken) {
    logout();
    throw { message: "Not authenticated", error_code: "UNAUTHORIZED" };
  }

  const isFormData = options.body instanceof FormData;

  options.headers = {
    ...((options.headers as Record<string, string>) || {}),
    Authorization: `Bearer ${accessToken}`,
    ...(isFormData ? {} : { "Content-Type": "application/json" }),
  };

  const url =
    path.startsWith("http://") || path.startsWith("https://")
      ? path
      : API_BASE_URL + path;

  const res = await fetch(url, options);

  let body: ApiResponse;
  try {
    body = await res.json();
  } catch {
    throw { message: `HTTP ${res.status}`, error_code: "HTTP_ERROR" };
  }

  if (!res.ok || body.success === false) {
    if (res.status === 401) {
      logout();
    }

    throw {
      message: body.message || "Unknown error",
      error_code: body.error_code || "UNKNOWN_ERROR",
    };
  }

  return body.data as T;
}
