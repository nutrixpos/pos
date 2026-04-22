<template>
  <div class="login-wrapper">
    <div class="login-card">
      <div class="login-logo">
        <img src="@/assets/logo.png" alt="NutrixPOS" />
      </div>

      <h1 class="login-title">Welcome Back</h1>
      <p class="login-subtitle">Sign in to continue to NutrixPOS</p>

      <form class="login-form" @submit.prevent="login">
        <div class="field">
          <label for="username">Username</label>
          <InputText
            id="username"
            v-model="form.username"
            placeholder="Enter your username"
            :class="{ 'p-invalid': errors.username }"
            class="w-full"
            autocomplete="username"
          />
          <small v-if="errors.username" class="p-error">{{ errors.username }}</small>
        </div>

        <div class="field">
          <label for="password">Password</label>
          <Password
            id="password"
            v-model="form.password"
            placeholder="Enter your password"
            :class="{ 'p-invalid': errors.password }"
            class="w-full"
            :feedback="false"
            toggleMask
            autocomplete="current-password"
            inputClass="w-full"
          />
          <small v-if="errors.password" class="p-error">{{ errors.password }}</small>
        </div>

        <div v-if="serverError" class="error-banner">
          <i class="pi pi-times-circle"></i>
          {{ serverError }}
        </div>

        <Button
          type="submit"
          label="Sign In"
          icon="pi pi-sign-in"
          class="login-btn"
          :loading="loading"
          :disabled="loading"
        />
      </form>

      <Divider>
        <span class="text-muted">or</span>
      </Divider>

      <div class="text-center">
        <small class="text-muted">
          Don't have an account?
          <router-link to="/register" class="text-primary font-semibold">Contact admin</router-link>
        </small>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import Button from 'primevue/button'
import Divider from 'primevue/divider'
import { useRouter } from 'vue-router'
import auth from '@/services/auth'

const router = useRouter()

const form = reactive({
  username: '',
  password: '',
})

const errors = reactive({
  username: '',
  password: '',
})

const loading = ref(false)
const serverError = ref('')

function validate(): boolean {
  errors.username = form.username.trim() ? '' : 'Username is required'
  errors.password = form.password ? '' : 'Password is required'
  return !errors.username && !errors.password
}

async function login() {
  if (!validate()) return

  loading.value = true
  serverError.value = ''

  try {
    const success = await auth.login(form.username.trim(), form.password)

    if (success) {
      router.push({ path: '/' })
    } else {
      serverError.value = 'Invalid username or password'
    }
  } catch (err: any) {
    serverError.value = err?.response?.data?.error || 'Login failed. Please try again.'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-wrapper {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  background: linear-gradient(135deg, #1a2a3a 0%, #2d4a6a 50%, #1a2a3a 100%);
  background-attachment: fixed;
}

.login-card {
  background: white;
  border-radius: 1.5rem;
  padding: 3rem 2.5rem;
  width: 100%;
  max-width: 420px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  animation: slideUp 0.5s ease-out;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.login-logo {
  display: flex;
  justify-content: center;
  margin-bottom: 1.5rem;
}

.login-logo img {
  height: 48px;
  object-fit: contain;
}

.login-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: #1a2a3a;
  text-align: center;
  margin: 0 0 0.5rem;
}

.login-subtitle {
  color: #64748b;
  text-align: center;
  font-size: 0.95rem;
  margin: 0 0 2rem;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.field label {
  color: #374151;
  font-size: 0.875rem;
  font-weight: 500;
}

:deep(.p-inputtext),
:deep(.p-password-input) {
  padding: 0.75rem 1rem;
  border-radius: 0.5rem;
  border: 1px solid #e2e8f0;
  transition: all 0.2s;
}

:deep(.p-inputtext:focus),
:deep(.p-password-input:focus) {
  border-color: #14977B;
  box-shadow: 0 0 0 3px rgba(20, 151, 123, 0.1);
}

:deep(.p-inputtext.p-invalid),
:deep(.p-password-input.p-invalid) {
  border-color: #ef4444;
}

:deep(.p-password) {
  width: 100%;
}

.error-banner {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 0.5rem;
  padding: 0.75rem 1rem;
  color: #dc2626;
  font-size: 0.875rem;
}

.error-banner .pi {
  font-size: 1rem;
  flex-shrink: 0;
}

.login-btn {
  width: 100%;
  justify-content: center;
  padding: 0.875rem !important;
  font-size: 1rem !important;
  font-weight: 600 !important;
  border-radius: 0.5rem !important;
  background: linear-gradient(135deg, #14977B 0%, #1fb590 100%) !important;
  border: none !important;
  transition: all 0.2s;
}

.login-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(20, 151, 123, 0.4);
}

.p-error {
  color: #ef4444;
  font-size: 0.8rem;
}

.text-center {
  text-align: center;
}

.text-muted {
  color: #94a3b8;
  font-size: 0.875rem;
}

.text-primary {
  color: #14977B;
}

.text-primary:hover {
  text-decoration: underline;
}
</style>