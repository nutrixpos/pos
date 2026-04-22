<template>
  <div class="setup-wrapper">
    <div class="setup-card">
      <div class="flex justify-content-end">
        <a href="https://nutrixpos.com/userguide/installation.html" target="_blank">
            <Button icon="pi pi-info-circle" class="p-button-text" />
        </a>
      </div>
      <div class="setup-icon">
        <i class="pi pi-user-plus"></i>
      </div>

      <h1 class="setup-title">Create Admin User</h1>
      <p class="setup-subtitle">
        No admin user exists yet.<br />
        Create your administrator account to get started.
      </p>

      <form class="setup-form" @submit.prevent="submit">
        <div class="field">
          <label for="username">Username</label>
          <InputText
            id="username"
            v-model="form.username"
            placeholder="e.g. admin"
            :class="{ 'p-invalid': errors.username }"
            class="w-full"
          />
          <small v-if="errors.username" class="p-error">{{ errors.username }}</small>
        </div>

        <div class="field">
          <label for="email">Email</label>
          <InputText
            id="email"
            v-model="form.email"
            type="email"
            placeholder="e.g. admin@example.com"
            :class="{ 'p-invalid': errors.email }"
            class="w-full"
          />
          <small v-if="errors.email" class="p-error">{{ errors.email }}</small>
        </div>

        <div class="field">
          <label for="password">Password</label>
          <InputText
            id="password"
            v-model="form.password"
            type="password"
            placeholder="Your password"
            :class="{ 'p-invalid': errors.password }"
            class="w-full"
          />
          <small v-if="errors.password" class="p-error">{{ errors.password }}</small>
        </div>

        <div class="field">
          <label for="confirm_password">Confirm Password</label>
          <InputText
            id="confirm_password"
            v-model="form.confirm_password"
            type="password"
            placeholder="Enter your password again"
            :class="{ 'p-invalid': errors.confirm_password }"
            class="w-full"
          />
          <small v-if="errors.confirm_password" class="p-error">{{ errors.confirm_password }}</small>
        </div>

        <div v-if="serverError" class="error-banner mt-3 mb-2">
          <i class="pi pi-times-circle"></i>
          {{ serverError }}
        </div>

        <div v-if="success" class="success-banner mb-2 mt-4">
          <i class="pi pi-check-circle"></i>
          Admin created successfully! Redirecting to home...
        </div>

        <div v-if="!success" class="flex flex-column gap-1 mt-1">
          <Button
            type="submit"
            label="Create Admin"
            icon="pi pi-user-plus"
            iconPos="right"
            class="submit-btn"
            :loading="loading"
            :disabled="loading"
          />
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import InputText from 'primevue/inputtext'
import Button from 'primevue/button'
import axios from 'axios'
import { useRouter } from 'vue-router'
import auth from '@/services/auth'

const router = useRouter()

const apiUrl = `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}`

const form = reactive({
  username: '',
  email: '',
  password: '',
  confirm_password: '',
})

const errors = reactive({
  username: '',
  email: '',
  password: '',
  confirm_password: '',
})

const loading = ref(false)
const serverError = ref('')
const success = ref(false)

function validate(): boolean {
  errors.username = form.username.trim() ? '' : 'Username is required'
  errors.email = form.email.trim() ? '' : 'Email is required'
  errors.email = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email) ? '' : 'Invalid email format'
  errors.password = form.password.length >= 6 ? '' : 'Password must be at least 6 characters'
  errors.confirm_password = form.password === form.confirm_password ? '' : 'Passwords do not match'

  return !errors.username && !errors.email && !errors.password && !errors.confirm_password
}

async function submit() {
  if (!validate()) return

  loading.value = true
  serverError.value = ''

  try {
    const response = await axios.post(`${apiUrl}/api/auth/register`, {
      username: form.username.trim(),
      email: form.email.trim(),
      password: form.password
    })

    const data = response.data

    auth.accessToken.value = data.token
    auth.currentUser.value = data.user
    auth.isAuthenticated.value = true

    localStorage.setItem('nutrix_token', data.token)
    localStorage.setItem('nutrix_user', JSON.stringify(data.user))

    success.value = true

    setTimeout(() => {
      router.push({ path: '/' })
    }, 1500)
  } catch (err: any) {
    serverError.value = err?.response?.data?.error || err?.response?.data?.message || 'Failed to create admin user'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>

.setup-wrapper {
  background: #fdf8f3;
  background-image: 
    radial-gradient(ellipse at top, #ffeaa7 0%, transparent 70%),
    radial-gradient(ellipse at bottom, #fab1a0 0%, transparent 50%),
    repeating-linear-gradient(
      45deg,
      transparent,
      transparent 10px,
      rgba(116, 185, 255, 0.03) 10px,
      rgba(116, 185, 255, 0.03) 20px
    );
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
}

.setup-card {
  background: white !important;
  backdrop-filter: blur(18px);
  -webkit-backdrop-filter: blur(18px);
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 1.5rem;
  padding: 3rem 2.5rem;
  width: 100%;
  max-width: 480px;
  box-shadow:
    0 8px 32px rgba(0, 0, 0, 0.4),
    0 0 0 1px rgba(255, 220, 0, 0.08);
  animation: fadeUp 0.4s ease both;
}

@keyframes fadeUp {
  from { opacity: 0; transform: translateY(24px); }
  to   { opacity: 1; transform: translateY(0); }
}

.setup-icon {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: linear-gradient(135deg, #DCDB7A 0%, #EAD60B 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 1.5rem;
  box-shadow: 0 4px 20px rgba(255, 220, 0, 0.35);
}

.setup-icon .pi {
  font-size: 1.8rem;
  color: #2E4762;
}

.setup-title {
  font-size: 1.7rem;
  font-weight: 700;
  color: #1a2a3a;
  text-align: center;
  margin: 0 0 0.5rem;
  letter-spacing: -0.02em;
}

.setup-subtitle {
  color: #1a2a3a;
  text-align: center;
  font-size: 0.95rem;
  line-height: 1.6;
  margin: 0 0 2rem;
}

.setup-form {
  display: flex;
  flex-direction: column;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.field label {
  color: #1a2a3a;
  font-size: 0.875rem;
  font-weight: 500;
}

:deep(.p-inputtext) {
  background: rgba(255, 255, 255, 0.08) !important;
  border: 1px solid rgba(255, 255, 255, 0.18) !important;
  color: #1a2a3a;
  border-radius: 0.6rem !important;
  padding: 0.65rem 0.9rem !important;
  transition: border-color 0.2s, box-shadow 0.2s;
}

:deep(.p-inputtext:focus) {
  border-color: #1a2a3a !important;
  box-shadow: 0 0 0 3px rgba(255, 220, 0, 0.2) !important;
  outline: none !important;
}

:deep(.p-inputtext.p-invalid) {
  border-color: #ff6b6b !important;
}

.p-error {
  color: #ff6b6b;
  font-size: 0.8rem;
}

.submit-btn {
  width: 100%;
  justify-content: center;
  border: none !important;
  border-radius: 0.6rem !important;
  font-weight: 600 !important;
  padding: 0.75rem !important;
  font-size: 1rem !important;
  transition: opacity 0.2s, transform 0.15s !important;
}

.submit-btn:hover:not(:disabled) {
  opacity: 0.9;
  transform: translateY(-1px);
}

.success-banner {
  display: flex;
  align-items: flex-start;
  gap: 0.6rem;
  background: rgba(0, 200, 100, 0.15);
  border: 1px solid rgba(0, 200, 100, 0.35);
  border-radius: 0.6rem;
  padding: 0.85rem 1rem;
  color: #1a2a3a;
  font-size: 0.9rem;
  line-height: 1.5;
}

.success-banner .pi {
  font-size: 1.1rem;
  margin-top: 0.05rem;
  flex-shrink: 0;
}

.error-banner {
  display: flex;
  align-items: center;
  gap: 0.6rem;
  background: rgba(255, 80, 80, 0.15);
  border: 1px solid rgba(255, 80, 80, 0.35);
  border-radius: 0.6rem;
  padding: 0.75rem 1rem;
  color: red;
  font-size: 0.875rem;
}
</style>