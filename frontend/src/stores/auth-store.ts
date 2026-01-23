import { writable } from "svelte/store";
import type { User } from "../models/user"; // Đảm bảo đúng đường dẫn model
import {
  getAccessToken,
  getUser,
  getRefreshToken,
  // 👇 Import thêm các hàm này
  clearAccessToken,
  clearRefreshToken,
  clearUser,
  setAccessToken, // (Optional) dùng cho setAuth
  setRefreshToken, // (Optional) dùng cho setAuth
  setUser as setStorageUser // (Optional) đổi tên để tránh trùng
} from "../services/storage-service";
import { isTokenExpired } from "../services/auth-service";

interface AuthState {
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
}

function getInitialAuthState(): AuthState {
  const accessToken = getAccessToken();
  const refreshToken = getRefreshToken();
  const user = getUser();

  // Kiểm tra token còn hạn không
  const isAuthenticated =
    !!accessToken &&
    !isTokenExpired(accessToken) &&
    !!user;

  // Nếu token hết hạn ngay lúc init, nên clear luôn storage cho sạch
  if (accessToken && !isAuthenticated) {
    clearAccessToken();
    clearRefreshToken();
    clearUser();
  }

  return {
    user: isAuthenticated ? user : null,
    accessToken: isAuthenticated ? accessToken : null,
    refreshToken: isAuthenticated ? refreshToken : null,
    isAuthenticated,
  };
}

export const authStore = writable<AuthState>(getInitialAuthState());

export function setAuth(
  user: User,
  accessToken: string,
  refreshToken: string
) {
  // 1. Lưu vào LocalStorage (để F5 không bị mất)
  setAccessToken(accessToken);
  setRefreshToken(refreshToken);
  setStorageUser(user);

  // 2. Cập nhật Store
  authStore.set({
    user,
    accessToken,
    refreshToken,
    isAuthenticated: true,
  });
}

// 👇 Hàm quan trọng nhất cho Logout
export function clearAuth() {
  // 1. Xóa sạch LocalStorage
  clearAccessToken();
  clearRefreshToken();
  clearUser();

  // 2. Reset Store về null
  authStore.set({
    user: null,
    accessToken: null,
    refreshToken: null,
    isAuthenticated: false,
  });
}