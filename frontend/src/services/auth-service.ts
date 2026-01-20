// src/services/auth-service.ts
import {
  setAccessToken,
  getAccessToken,
  setRefreshToken,
  getRefreshToken,
  clearAccessToken,
  clearRefreshToken,
  setUser,
  clearUser,
} from "./storage-service";

import type {
  LoginRequestDto,
  LoginResponseDto,
  SendVerificationRequestDto,
  VerifyEmailRequestDto,
  VerifyEmailResponseDto,
  CompleteRegistrationRequestDto,
  RegisterResponseDto,
  RefreshTokenResponseDto,
  LogoutRequestDto,
  ForgotPasswordRequestDto,
  VerifyResetOtpRequestDto,
  VerifyResetOtpResponseDto,
  ResetPasswordRequestDto,
  CompleteGoogleSetupRequestDto,
} from "../dtos/auth-dto";

import { apiFetch } from "../utils/api-fetch";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "";

/* ---------------- Helpers ---------------- */
function parseJwt(token: string): any {
  try {
    const base64 = token.split(".")[1];
    return JSON.parse(atob(base64));
  } catch {
    return null;
  }
}

export function isTokenExpired(token: string): boolean {
  const payload = parseJwt(token);
  if (!payload?.exp) return true;
  return Date.now() >= payload.exp * 1000;
}

/* ---------------- Core Auth Services ---------------- */

/** Đăng nhập */
export async function login(payload: LoginRequestDto): Promise<LoginResponseDto> {
  const data = await apiFetch<LoginResponseDto>("/api/auth/local/login", {
    method: "POST",
    body: JSON.stringify(payload),
    skipAuth: true
  });

  setAccessToken(data.accessToken);
  setRefreshToken(data.refreshToken);
  setUser(data.user);
  return data;
}

/** Gửi mã xác thực email (Đăng ký) */
export async function sendEmailVerification(payload: SendVerificationRequestDto): Promise<void> {
  await apiFetch("/api/auth/local/send-verification", {
    method: "POST",
    body: JSON.stringify(payload),
    skipAuth: true
  });
}

/** Xác thực OTP email */
export async function verifyEmail(payload: VerifyEmailRequestDto): Promise<VerifyEmailResponseDto> {
  return await apiFetch<VerifyEmailResponseDto>("/api/auth/local/verify-email", {
    method: "POST",
    body: JSON.stringify(payload),
    skipAuth: true
  });
}

/** Hoàn tất đăng ký */
export async function completeRegistration(payload: CompleteRegistrationRequestDto): Promise<RegisterResponseDto> {
  const data = await apiFetch<RegisterResponseDto>("/api/auth/local/complete-registration", {
    method: "POST",
    body: JSON.stringify(payload),
    skipAuth: true
  });

  setAccessToken(data.accessToken);
  setRefreshToken(data.refreshToken);
  setUser(data.user);
  return data;
}

/** Làm mới Token */
export async function refreshToken(): Promise<RefreshTokenResponseDto | null> {
  const currentToken = getRefreshToken();
  if (!currentToken) return null;

  try {
    const data = await apiFetch<RefreshTokenResponseDto>("/api/auth/refresh", {
      method: "POST",
      body: JSON.stringify({ refreshToken: currentToken }),
      skipAuth: true
    });

    setAccessToken(data.accessToken);
    setRefreshToken(data.refreshToken);
    return data;
  } catch (err) {
    logout();
    return null;
  }
}

/** Lấy Access Token còn hạn */
export async function getValidAccessToken(): Promise<string | null> {
  const token = getAccessToken();
  if (!token) return null;
  if (!isTokenExpired(token)) return token;

  const refreshed = await refreshToken();
  return refreshed?.accessToken ?? null;
}

/** Đăng xuất */
export async function logout(payload?: LogoutRequestDto): Promise<void> {
  try {
    await apiFetch("/api/auth/logout", {
      method: "POST",
      body: JSON.stringify(payload ?? {}),
    }).catch(() => null);
  } finally {
    clearAccessToken();
    clearRefreshToken();
    clearUser();
    window.location.href = "/#/login";
  }
}

/* ---------------- Password Reset Flow ---------------- */

export async function forgotPassword(payload: ForgotPasswordRequestDto): Promise<void> {
  await apiFetch("/api/auth/local/forgot-password", {
    method: "POST",
    body: JSON.stringify(payload),
    skipAuth: true
  });
}

export async function verifyResetOtp(payload: VerifyResetOtpRequestDto): Promise<VerifyResetOtpResponseDto> {
  return await apiFetch<VerifyResetOtpResponseDto>("/api/auth/local/verify-reset-otp", {
    method: "POST",
    body: JSON.stringify(payload),
    skipAuth: true
  });
}

export async function resetPassword(payload: ResetPasswordRequestDto): Promise<void> {
  await apiFetch("/api/auth/local/reset-password", {
    method: "POST",
    body: JSON.stringify(payload),
    skipAuth: true
  });
}

/* ---------------- Google Auth ---------------- */

export function loginWithGoogle() {
  window.location.href = `${API_BASE_URL}/api/auth/google/login`;
}

export function handleLoginCallback() {
  const hash = window.location.hash;
  const queryString = hash.includes("?") ? hash.split("?")[1] : "";
  const params = new URLSearchParams(queryString);

  const accessToken = params.get('access_token');
  const refreshToken = params.get('refresh_token');
  const userStr = params.get('user');
  const setupToken = params.get('setup_token');

  if (setupToken) return { success: false, setupRequired: true, setupToken };
  if (!accessToken) return { success: false };

  setAccessToken(accessToken);
  if (refreshToken) setRefreshToken(refreshToken);

  let userObj = null;
  if (userStr) {
    try {
      userObj = JSON.parse(decodeURIComponent(userStr));
      setUser(userObj);
    } catch (e) {
      console.error("Lỗi parse user info từ URL:", e);
    }
  }

  return { success: true, user: userObj, accessToken };
}

export async function completeGoogleSetup(payload: CompleteGoogleSetupRequestDto): Promise<RegisterResponseDto> {
  const data = await apiFetch<RegisterResponseDto>("/api/auth/google/complete-setup", {
    method: "POST",
    body: JSON.stringify(payload),
    skipAuth: true
  });

  setAccessToken(data.accessToken);
  setRefreshToken(data.refreshToken);
  setUser(data.user);
  return data;
}

/* ---------------- Exports ---------------- */
export const authService = {
  isTokenExpired,
  getValidAccessToken,
  login,
  sendEmailVerification,
  verifyEmail,
  completeRegistration,
  refreshToken,
  logout,
  forgotPassword,
  verifyResetOtp,
  resetPassword,
  completeGoogleSetup,
  loginWithGoogle,
  handleLoginCallback,
};

export default authService;