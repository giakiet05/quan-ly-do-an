import { writable } from "svelte/store";
import type { User } from "../models/user";
import { getAccessToken, getUser, getRefreshToken } from "../services/storage-service";
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

  const isAuthenticated =
    !!accessToken &&
    !isTokenExpired(accessToken) &&
    !!user;

  return {
    user,
    accessToken,
    refreshToken,
    isAuthenticated,
  };
}


export const authStore = writable<AuthState>(getInitialAuthState());


export function setAuth(
  user: User,
  accessToken: string,
  refreshToken: string
) {
  authStore.set({
    user,
    accessToken,
    refreshToken,
    isAuthenticated: true,
  });
}

export function clearAuth() {
  authStore.set({
    user: null,
    accessToken: null,
    refreshToken: null,
    isAuthenticated: false,
  });
}

