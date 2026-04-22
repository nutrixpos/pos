import { ref,  getCurrentInstance } from 'vue'
import { useRouter } from 'vue-router'


const TOKEN_KEY = 'nutrix_token'
const USER_KEY = 'nutrix_user'

export interface User {
  id: string
  username: string
  email: string
  roles: string[]
}

export interface LoginResponse {
  token: string
  user: User
}

const currentUser = ref<User | null>(null)
const accessToken = ref<string | null>(null)
const isAuthenticated = ref(false)

function loadFromStorage() {
  const token = localStorage.getItem(TOKEN_KEY)
  const userStr = localStorage.getItem(USER_KEY)
  
  if (token && userStr) {
    accessToken.value = token
    currentUser.value = JSON.parse(userStr)
    isAuthenticated.value = true
  }
}

loadFromStorage()

export const auth = {
  accessToken,
  currentUser,
  isAuthenticated,

  async login(username: string, password: string): Promise<boolean> {
    try {
      const response = await fetch(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
      })

      if (!response.ok) {
        return false
      }

      const data: LoginResponse = await response.json()
      
      accessToken.value = data.token
      currentUser.value = data.user
      isAuthenticated.value = true

      localStorage.setItem(TOKEN_KEY, data.token)
      localStorage.setItem(USER_KEY, JSON.stringify(data.user))

      return true
    } catch (error) {
      console.error('Login failed:', error)
      return false
    }
  },

  async register(username: string, email: string, password: string): Promise<boolean> {
    try {
      const response = await fetch('/api/auth/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, email, password }),
      })

      if (!response.ok) {
        return false
      }

      const data: LoginResponse = await response.json()
      
      accessToken.value = data.token
      currentUser.value = data.user
      isAuthenticated.value = true

      localStorage.setItem(TOKEN_KEY, data.token)
      localStorage.setItem(USER_KEY, JSON.stringify(data.user))

      return true
    } catch (error) {
      console.error('Registration failed:', error)
      return false
    }
  },

  async getCurrentUser(): Promise<User | null> {
    if (!accessToken.value) {
      return null
    }

    try {
      const response = await fetch('/api/auth/me', {
        headers: {
          'Authorization': `Bearer ${accessToken.value}`,
        },
      })

      if (!response.ok) {
        this.signOut()
        return null
      }

      const user: User = await response.json()
      currentUser.value = user
      localStorage.setItem(USER_KEY, JSON.stringify(user))

      return user
    } catch (error) {
      console.error('Get current user failed:', error)
      return null
    }
  },

  signOut() {
    accessToken.value = null
    currentUser.value = null
    isAuthenticated.value = false

    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USER_KEY)
  },

  hasRole(role: string): boolean {
    if (!currentUser.value) {
      return false
    }
    return currentUser.value.roles.includes(role)
  },
}

export default auth