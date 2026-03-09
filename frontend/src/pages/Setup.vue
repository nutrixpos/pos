<template>
  <div class="setup-wrapper">
    <div class="setup-card">
      <!-- Logo / Icon -->
      <div class="setup-icon">
        <i class="pi pi-database"></i>
      </div>

      <h1 class="setup-title">Welcome to NutrixPOS</h1>
      <p class="setup-subtitle">
        No database connection has been configured yet.<br />
        Please provide your MongoDB connection details below.
      </p>

      <form class="setup-form" @submit.prevent="submit">
        <!-- Host -->
        <div class="field">
          <label for="host">Host</label>
          <InputText
            id="host"
            v-model="form.host"
            placeholder="e.g. localhost"
            :class="{ 'p-invalid': errors.host }"
            class="w-full"
          />
          <small v-if="errors.host" class="p-error">{{ errors.host }}</small>
        </div>

        <!-- Port -->
        <div class="field">
          <label for="port">Port</label>
          <InputNumber
            id="port"
            v-model="form.port"
            :useGrouping="false"
            placeholder="e.g. 27017"
            :class="{ 'p-invalid': errors.port }"
            class="w-full"
          />
          <small v-if="errors.port" class="p-error">{{ errors.port }}</small>
        </div>

        <!-- Database -->
        <div class="field">
          <label for="database">Database name</label>
          <InputText
            id="database"
            v-model="form.database"
            placeholder="e.g. nutrix"
            :class="{ 'p-invalid': errors.database }"
            class="w-full"
          />
          <small v-if="errors.database" class="p-error">{{ errors.database }}</small>
        </div>

        <!-- Database -->
        <div class="field">
          <label for="username">Username</label>
          <InputText
            id="username"
            v-model="form.username"
            placeholder="e.g. webadmin"
            :class="{ 'p-invalid': errors.username }"
            class="w-full"
          />
          <small v-if="errors.username" class="p-error">{{ errors.username }}</small>
        </div>

        <!-- Database -->
        <div class="field">
          <label for="password">Password</label>
          <InputText
            id="password"
            v-model="form.password"
            type="password"
            placeholder="Your database user password"
            :class="{ 'p-invalid': errors.password }"
            class="w-full"
          />
          <small v-if="errors.password" class="p-error">{{ errors.password }}</small>
        </div>

        <!-- Database -->
        <div class="field">
          <label for="password">Confirm Password</label>
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

        <!-- Success banner -->
        <div v-if="saved" class="success-banner mb-2 mt-4">
          <i class="pi pi-check-circle"></i>
          Configuration saved! The server is setup — please restart the application to apply the new settings.
        </div>

        <!-- Error banner -->
        <div v-if="serverError" class="error-banner mt-3 mb-2">
          <i class="pi pi-times-circle"></i>
          {{ serverError }}
        </div>

        <Button
          v-if="!saved"
          type="submit"
          label="Save & Continue"
          icon="pi pi-arrow-right"
          iconPos="right"
          class="submit-btn"
          :loading="loading"
          :disabled="loading || saved"
        />

        <Button v-if="saved" class="submit-btn mt-2" label="Let's go 🚀" @click="router.push({ path: '/admin' })" />
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Button from 'primevue/button'
import axios from 'axios'
import { useRouter } from 'vue-router'
const router = useRouter()

const backendBase = `http://${import.meta.env.VITE_APP_BACKEND_HOST || 'localhost:8000'}`

const form = reactive({
  host: '',
  port: null as number | null,
  database: '',
  username: '',
  password: '',
  confirm_password: '',
})

const errors = reactive({
  host: '',
  port: '',
  database: '',
  username: '',
  password: '',
  confirm_password: '',
})

const loading = ref(false)
const saved   = ref(false)
const serverError = ref('')

function validate(): boolean {
  errors.host = form.host.trim() ? '' : 'Host is required'
  errors.port = form.port && form.port > 0 ? '' : 'A valid port number is required'
  errors.database = form.database.trim() ? '' : 'Database name is required'
  errors.confirm_password = form.password == form.confirm_password ? '' : 'Passwords do not match'

  return !errors.host && !errors.port && !errors.database && !errors.username && !errors.password && !errors.confirm_password
}

async function submit() {
  if (!validate()) return

  loading.value = true
  serverError.value = ''

  try {
    await axios.post(`${backendBase}/api/setup/config`, {
      host:     form.host.trim(),
      port:     form.port,
      database: form.database.trim(),
    })
    saved.value = true
  } catch (err: any) {
    serverError.value =
      err?.response?.data || 'Failed to save configuration. Check the backend logs.'
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

/* Override PrimeVue input for dark glassmorphism card */
:deep(.p-inputtext),
:deep(.p-inputnumber-input) {
  background: rgba(255, 255, 255, 0.08) !important;
  border: 1px solid rgba(255, 255, 255, 0.18) !important;
  color: #1a2a3a;
  border-radius: 0.6rem !important;
  padding: 0.65rem 0.9rem !important;
  transition: border-color 0.2s, box-shadow 0.2s;
}

:deep(.p-inputtext:focus),
:deep(.p-inputnumber-input:focus) {
  border-color: #1a2a3a !important;
  box-shadow: 0 0 0 3px rgba(255, 220, 0, 0.2) !important;
  outline: none !important;
}

:deep(.p-inputtext.p-invalid),
:deep(.p-inputnumber-input.p-invalid) {
  border-color: #ff6b6b !important;
}

:deep(.p-inputnumber) {
  width: 100%;
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
