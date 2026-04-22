import auth from './auth'

export interface ApiOptions extends RequestInit {
  skipAuth?: boolean
}

async function api<T>(url: string, options: ApiOptions = {}): Promise<T> {
  const { skipAuth, ...fetchOptions } = options

  const headers: HeadersInit = {
    ...options.headers,
  }

  if (!skipAuth && auth.accessToken.value) {
    (headers as Record<string, string>)['Authorization'] = `Bearer ${auth.accessToken.value}`
  }

  const response = await fetch(url, {
    ...fetchOptions,
    headers,
  })

  if (response.status === 401) {
    auth.signOut()
    window.location.href = '/'
    throw new Error('Unauthorized')
  }

  if (!response.ok) {
    throw new Error(`API Error: ${response.status}`)
  }

  return response.json()
}

export default api