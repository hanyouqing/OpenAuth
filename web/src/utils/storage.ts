const TOKEN_KEY = 'access_token'
const REFRESH_TOKEN_KEY = 'refresh_token'
const LANGUAGE_KEY = 'language'

export const storage = {
  getToken(): string | null {
    return localStorage.getItem(TOKEN_KEY)
  },

  setToken(token: string): void {
    localStorage.setItem(TOKEN_KEY, token)
  },

  removeToken(): void {
    localStorage.removeItem(TOKEN_KEY)
  },

  getRefreshToken(): string | null {
    return localStorage.getItem(REFRESH_TOKEN_KEY)
  },

  setRefreshToken(token: string): void {
    localStorage.setItem(REFRESH_TOKEN_KEY, token)
  },

  removeRefreshToken(): void {
    localStorage.removeItem(REFRESH_TOKEN_KEY)
  },

  getLanguage(): string {
    return localStorage.getItem(LANGUAGE_KEY) || 'en'
  },

  setLanguage(lang: string): void {
    localStorage.setItem(LANGUAGE_KEY, lang)
  },

  clear(): void {
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(REFRESH_TOKEN_KEY)
    this.clearUserInfo()
  },

  getUserInfo(): any {
    const userStr = localStorage.getItem('user_info')
    return userStr ? JSON.parse(userStr) : null
  },

  setUserInfo(user: any): void {
    localStorage.setItem('user_info', JSON.stringify(user))
  },

  clearUserInfo(): void {
    localStorage.removeItem('user_info')
  },
}
