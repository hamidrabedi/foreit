const AUTH_TOKEN_KEY = "admin_token";

export function getAuthToken(): string | null {
  return localStorage.getItem(AUTH_TOKEN_KEY);
}

export function isAuthenticated(): boolean {
  return !!getAuthToken();
}

export function clearAuth() {
  localStorage.removeItem(AUTH_TOKEN_KEY);
}
