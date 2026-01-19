import type { ApiResponse } from "../dtos/api-response-dto";
import { convertKeysToCamel, convertKeysToSnake } from "./case-converter";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";

interface FetchOptions extends RequestInit {
  skipAuth?: boolean;
}

/**
 * Get access token from localStorage
 */
function getAccessToken(): string | null {
  return localStorage.getItem("access_token");
}

/**
 * Get refresh token from localStorage
 */
function getRefreshToken(): string | null {
  return localStorage.getItem("refresh_token");
}

/**
 * Save tokens to localStorage
 */
function saveTokens(accessToken: string, refreshToken?: string): void {
  localStorage.setItem("access_token", accessToken);
  if (refreshToken) {
    localStorage.setItem("refresh_token", refreshToken);
  }
}

/**
 * Clear tokens from localStorage
 */
function clearTokens(): void {
  localStorage.removeItem("access_token");
  localStorage.removeItem("refresh_token");
}

/**
 * Refresh access token using refresh token
 */
async function refreshAccessToken(): Promise<boolean> {
  const refreshToken = getRefreshToken();
  
  if (!refreshToken) {
    return false;
  }

  try {
    const response = await fetch(`${API_BASE_URL}/api/auth/refresh`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });

    if (!response.ok) {
      clearTokens();
      return false;
    }

    const data: ApiResponse<{ access_token: string; refresh_token: string }> =
      await response.json();

    if (data.success && data.data) {
      saveTokens(data.data.access_token, data.data.refresh_token);
      return true;
    }

    clearTokens();
    return false;
  } catch (error) {
    console.error("Failed to refresh token:", error);
    clearTokens();
    return false;
  }
}

/**
 * Enhanced fetch function with automatic authentication and token refresh
 * 
 * @param url - API endpoint (relative to base URL or absolute)
 * @param options - Fetch options with optional skipAuth flag
 * @returns Response data
 * @throws Error if request fails
 */
export async function apiFetch<T = any>(
  url: string,
  options: FetchOptions = {}
): Promise<T> {
  const { skipAuth = false, ...fetchOptions } = options;

  // Build full URL
  const fullUrl = url.startsWith("http") ? url : `${API_BASE_URL}${url}`;

  // Prepare headers
  const headers: Record<string, string> = {
    ...((fetchOptions.headers as Record<string, string>) || {}),
  };

  // Only set Content-Type to JSON if not FormData (FormData needs browser to set multipart/form-data with boundary)
  if (!(fetchOptions.body instanceof FormData)) {
    headers["Content-Type"] = "application/json";
    
    // Convert request body from camelCase to snake_case for backend
    if (fetchOptions.body && typeof fetchOptions.body === 'string') {
      try {
        const parsed = JSON.parse(fetchOptions.body);
        fetchOptions.body = JSON.stringify(convertKeysToSnake(parsed));
      } catch {
        // Keep original if not valid JSON
      }
    }
  }

  // Add authorization header if not skipping auth
  if (!skipAuth) {
    const accessToken = getAccessToken();
    if (accessToken) {
      headers["Authorization"] = `Bearer ${accessToken}`;
    }
  }

  // First attempt
  let response = await fetch(fullUrl, {
    ...fetchOptions,
    headers,
  });

  // If unauthorized and we have a refresh token, try to refresh
  if (response.status === 401 && !skipAuth) {
    const refreshed = await refreshAccessToken();

    if (refreshed) {
      // Retry with new token
      const newAccessToken = getAccessToken();
      if (newAccessToken) {
        headers["Authorization"] = `Bearer ${newAccessToken}`;
      }

      response = await fetch(fullUrl, {
        ...fetchOptions,
        headers,
      });
    } else {
      // Refresh failed, redirect to login
      clearTokens();
      window.location.href = "/#/auth/login";
      throw new Error("Session expired. Please login again.");
    }
  }

  // Handle non-OK responses
  if (!response.ok) {
    let errorMessage = `HTTP ${response.status}: ${response.statusText}`;

    try {
      const errorData: ApiResponse<any> = await response.json();
      if (errorData.message) {
        errorMessage = errorData.message;
      }
      if (errorData.error_code) {
        errorMessage += ` (${errorData.error_code})`;
      }
    } catch {
      // Could not parse error as JSON
    }

    throw new Error(errorMessage);
  }

  // Parse and return response
  const contentType = response.headers.get("content-type");
  
  if (contentType?.includes("application/json")) {
    const jsonResponse: ApiResponse<T> = await response.json();
    
    // Convert response from snake_case to camelCase
    const convertedResponse = convertKeysToCamel<ApiResponse<T>>(jsonResponse);
    
    // Backend wraps all responses in ApiResponse { success, message, data }
    // Extract and return only the data field
    if (convertedResponse.data !== undefined) {
      return convertedResponse.data as T;
    }
    
    // If no data field, return the whole response (e.g., for error responses)
    return convertedResponse as any;
  }

  // For non-JSON responses (e.g., file downloads)
  return response as any;
}

/**
 * Upload file with progress tracking (optional)
 */
export async function uploadFile(
  url: string,
  file: File,
  fieldName: string = "file",
  additionalData?: Record<string, any>,
  onProgress?: (progress: number) => void
): Promise<any> {
  const fullUrl = url.startsWith("http") ? url : `${API_BASE_URL}${url}`;
  const accessToken = getAccessToken();

  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest();
    const formData = new FormData();

    formData.append(fieldName, file);

    // Add additional data
    if (additionalData) {
      Object.entries(additionalData).forEach(([key, value]) => {
        formData.append(key, value);
      });
    }

    // Track upload progress
    if (onProgress) {
      xhr.upload.addEventListener("progress", (e) => {
        if (e.lengthComputable) {
          const progress = (e.loaded / e.total) * 100;
          onProgress(progress);
        }
      });
    }

    // Handle completion
    xhr.addEventListener("load", () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        try {
          const response = JSON.parse(xhr.responseText);
          resolve(response);
        } catch {
          resolve(xhr.responseText);
        }
      } else {
        reject(new Error(`Upload failed: ${xhr.status} ${xhr.statusText}`));
      }
    });

    // Handle errors
    xhr.addEventListener("error", () => {
      reject(new Error("Upload failed due to network error"));
    });

    // Send request
    xhr.open("POST", fullUrl);
    if (accessToken) {
      xhr.setRequestHeader("Authorization", `Bearer ${accessToken}`);
    }
    xhr.send(formData);
  });
}

/**
 * Helper to check if user is authenticated
 */
export function isAuthenticated(): boolean {
  return !!getAccessToken();
}

/**
 * Export token management functions
 */
export { getAccessToken, getRefreshToken, saveTokens, clearTokens };
