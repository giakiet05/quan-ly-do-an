// src/services/auth-service.ts
import {
  setAccessToken,
  getAccessToken,
  clearAccessToken,
  setRefreshToken,
  getRefreshToken,
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
  RefreshTokenRequestDto,
  RefreshTokenResponseDto,
  LogoutRequestDto,
  ForgotPasswordRequestDto,
  VerifyResetOtpRequestDto,
  VerifyResetOtpResponseDto,
  ResetPasswordRequestDto,
  CompleteGoogleSetupRequestDto,
} from "../dtos/auth-dto";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "";

/* ---------------- Helpers ---------------- */
function safeJson(res: Response) {
  return res.json().catch(() => null);
}

function throwIfError(res: Response, json: any) {
  if (!res.ok) {
    // prefer { message, error_code } or json itself
    const err = json ?? { message: `HTTP ${res.status}` };
    throw err;
  }
}

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

/** Login - returns LoginResponseDto and persists tokens/user */
export async function login(payload: LoginRequestDto): Promise<LoginResponseDto> {
  const res = await fetch(`${API_BASE_URL}/api/auth/local/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

  const json = await safeJson(res);
  throwIfError(res, json);

  const data = (json?.data ?? json) as LoginResponseDto;

  setAccessToken(data.access_token);
  setRefreshToken(data.refresh_token);
  setUser(data.user);

  return data;
}

/** Send email verification (start register flow) */
export async function sendEmailVerification(
  payload: SendVerificationRequestDto
): Promise<void> {
  const res = await fetch(`${API_BASE_URL}/api/auth/local/send-verification`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

  const json = await safeJson(res);
  throwIfError(res, json);
  return;
}

/** Verify OTP for email - returns verificationToken */
export async function verifyEmail(
  payload: VerifyEmailRequestDto
): Promise<VerifyEmailResponseDto> {
  const res = await fetch(`${API_BASE_URL}/api/auth/local/verify-email`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

  const json = await safeJson(res);
  throwIfError(res, json);

  return (json?.data ?? json) as VerifyEmailResponseDto;
}

/** Complete registration using verificationToken -> returns RegisterResponseDto and persists tokens/user */
export async function completeRegistration(
  payload: CompleteRegistrationRequestDto
): Promise<RegisterResponseDto> {
  const res = await fetch(`${API_BASE_URL}/api/auth/local/complete-registration`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
  console.log("Response from completeRegistration:", payload, res);
  const json = await safeJson(res);
  throwIfError(res, json);

  const data = (json?.data ?? json) as RegisterResponseDto;

  setAccessToken(data.access_token);
  setRefreshToken(data.refresh_token);
  setUser(data.user);

  return data;
}

/** Refresh tokens using refreshToken from storage (returns new tokens) */
export async function refreshToken(): Promise<RefreshTokenResponseDto | null> {
  const refreshToken = getRefreshToken();
  if (!refreshToken) return null;

  const payload: RefreshTokenRequestDto = { refresh_token: refreshToken };

  const res = await fetch(`${API_BASE_URL}/api/auth/refresh`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

  const json = await safeJson(res);

  if (!res.ok) {
    // refresh failed -> force logout on client
    logout();
    return null;
  }

  const data = (json?.data ?? json) as RefreshTokenResponseDto;
  setAccessToken(data.access_token);
  setRefreshToken(data.refresh_token);
  return data;
}

/** Get valid access token: if expired attempt refresh; returns accessToken or null */
export async function getValidAccessToken(): Promise<string | null> {
  const token = getAccessToken();
  if (!token) return null;
  if (!isTokenExpired(token)) return token;

  const refreshed = await refreshToken();
  return refreshed?.access_token ?? null;
}

/** Logout: notify backend if possible, then clear client storage */
export async function logout(payload?: LogoutRequestDto): Promise<void> {
  try {
    // Try notify backend (if route exists). If backend requires tokens, payload should be provided.
    await fetch(`${API_BASE_URL}/api/auth/logout`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload ?? {}),
    }).catch(() => null);
  } catch {
    // ignore
  } finally {
    clearAccessToken();
    clearRefreshToken();
    clearUser();
    // navigate to login - leave routing to app if needed
    try {
      window.location.href = "/#/login";
    } catch { }
  }
}

/* ---------------- Password reset flow ---------------- */
export async function forgotPassword(payload: ForgotPasswordRequestDto): Promise<void> {
  const res = await fetch(`${API_BASE_URL}/api/auth/local/forgot-password`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

  const json = await safeJson(res);
  throwIfError(res, json);
  return;
}

export async function verifyResetOtp(
  payload: VerifyResetOtpRequestDto
): Promise<VerifyResetOtpResponseDto> {
  const res = await fetch(`${API_BASE_URL}/api/auth/local/verify-reset-otp`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

  const json = await safeJson(res);
  throwIfError(res, json);
  return (json?.data ?? json) as VerifyResetOtpResponseDto;
}

export async function resetPassword(payload: ResetPasswordRequestDto): Promise<void> {
  const res = await fetch(`${API_BASE_URL}/api/auth/local/reset-password`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

  const json = await safeJson(res);
  throwIfError(res, json);
  return;
}

/* ---------------- Google setup ---------------- */
export function loginWithGoogle() {
  window.location.href = `${API_BASE_URL}/api/auth/google/login`;
}

export function handleLoginCallback(): { success: boolean; user?: any; accessToken?: string } {
  const params = new URLSearchParams(window.location.search);
  const accessToken = params.get('accessToken');
  const refreshToken = params.get('refreshToken');
  const userStr = params.get('user');

  if (!accessToken) {
    return { success: false };
  }

  setAccessToken(accessToken);

  if (refreshToken) {
    setRefreshToken(refreshToken);
  }

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

export async function completeGoogleSetup(
  payload: CompleteGoogleSetupRequestDto
): Promise<RegisterResponseDto> {
  const res = await fetch(`${API_BASE_URL}/api/auth/google/complete-setup`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

  const json = await safeJson(res);
  throwIfError(res, json);

  const data = (json?.data ?? json) as RegisterResponseDto;
  setAccessToken(data.access_token);
  setRefreshToken(data.refresh_token);
  setUser(data.user);
  return data;
}

/* ---------------- Exports & compatibility ---------------- */
export const authService = {
  // helpers
  isTokenExpired,
  getValidAccessToken,
  // core
  login,
  sendEmailVerification,
  verifyEmail,
  completeRegistration,
  refreshToken,
  logout,
  // password
  forgotPassword,
  verifyResetOtp,
  resetPassword,
  // google
  completeGoogleSetup,
  loginWithGoogle,
  handleLoginCallback,
};

export default authService;
