export const USER_CODE_KEY = 'beruangan_user_code'

export const getSavedCode = () => localStorage.getItem(USER_CODE_KEY) || ''

export const saveCode = (code) => localStorage.setItem(USER_CODE_KEY, code)

export const clearCode = () => localStorage.removeItem(USER_CODE_KEY)

export const authHeaders = (code) => ({
  'Content-Type': 'application/json',
  'X-User-ID': code,
})
