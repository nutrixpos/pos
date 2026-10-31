<template>
  <div class="w-full">
    <div class="grid mx-2">
      <div class="col-12">
        <h3>{{ t('backups') }}</h3>
        <p class="mt-0 text-color-secondary">{{ t('backups_description') }}</p>
      </div>

      <div class="col-12 lg:col-4">
        <Card>
          <template #title>{{ t('export_backups') }}</template>
          <template #content>
            <p class="mt-0">{{ t('backups_export_description') }}</p>
          </template>
          <template #footer>
            <Button
              :label="t('export_backups')"
              icon="pi pi-file-export"
              :loading="exporting"
              @click="exportBackup"
            />
          </template>
        </Card>
      </div>

      <div class="col-12 lg:col-8">
        <Card>
          <template #title>{{ t('import_backups') }}</template>
          <template #content>
            <p class="mt-0">{{ t('backups_import_description') }}</p>

            <div class="flex flex-wrap gap-2 align-items-center">
              <FileUpload
                mode="basic"
                name="file"
                accept=".zip,application/zip"
                :auto="false"
                :multiple="false"
                :maxFileSize="maxUploadBytes"
                :chooseLabel="t('select_zip_file')"
                @select="onFileSelect"
              />
              <Button
                v-if="selectedFile"
                :label="t('backups_preview')"
                icon="pi pi-eye"
                severity="secondary"
                :loading="previewLoading"
                :disabled="previewLoading"
                @click="loadPreview"
              />
            </div>

            <p v-if="selectedFile" class="text-sm mt-2 mb-0">
              <i class="pi pi-file mr-2"></i>{{ selectedFile.name }}
            </p>

            <div v-if="previewLoading && !preview" class="flex align-items-center gap-2 mt-3 text-color-secondary">
              <ProgressSpinner style="width: 1.1rem; height: 1.1rem" strokeWidth="6" />
              <span>{{ t('backups_preview_loading') }}</span>
            </div>

            <template v-if="preview">
              <Message severity="warn" :closable="false" class="mt-4">
                {{ t('backups_import_warning') }}
              </Message>
              <Message v-if="sections.sales" severity="error" :closable="false" class="mt-2">
                {{ t('backups_sales_wipe_warning') }}
              </Message>
              <Message severity="info" :closable="false" class="mt-2">
                {{ t('backups_notice') }}
              </Message>

              <div class="flex align-items-center gap-2 mt-4">
                <h4 class="m-0">{{ t('backups_sections') }}</h4>
                <ProgressSpinner v-if="previewLoading" style="width: 1.1rem; height: 1.1rem" strokeWidth="6" />
              </div>
              <div :class="{ 'opacity-50 pointer-events-none': previewLoading }">
                <div class="grid">
                  <div class="col-12 md:col-6 flex align-items-center gap-2">
                    <Checkbox v-model="sections.materials" binary inputId="section-materials" />
                    <label for="section-materials">{{ t('backups_section_materials') }}</label>
                  </div>
                  <div class="col-12 md:col-6 flex align-items-center gap-2 md:pl-4">
                    <Checkbox v-model="includeStock" binary inputId="section-stock" />
                    <label for="section-stock">{{ t('backups_section_stock') }}</label>
                  </div>
                  <div class="col-12 md:col-6 flex align-items-center gap-2">
                    <Checkbox v-model="sections.products" binary inputId="section-products" />
                    <label for="section-products">{{ t('backups_section_products') }}</label>
                  </div>
                  <div class="col-12 md:col-6 flex align-items-center gap-2 md:pl-4">
                    <Checkbox v-model="includeImages" binary inputId="section-images" />
                    <label for="section-images">{{ t('backups_section_images') }}</label>
                  </div>
                  <div class="col-12 md:col-6 flex align-items-center gap-2">
                    <Checkbox v-model="sections.categories" binary inputId="section-categories" />
                    <label for="section-categories">{{ t('backups_section_categories') }}</label>
                  </div>
                  <div class="col-12 md:col-6 flex align-items-center gap-2 md:pl-4">
                    <Checkbox v-model="sections.sales" binary inputId="section-sales" />
                    <label for="section-sales">{{ t('sales_orders') }}</label>
                  </div>
                </div>
                <p class="text-sm text-color-secondary">{{ t('backups_dependency_note') }}</p>

                <DataTable :value="countRows" size="small" stripedRows class="mt-3">
                  <Column field="label" :header="t('backups_section')" />
                  <Column field="incoming" :header="t('backups_incoming')" />
                  <Column field="current" :header="t('backups_current')" />
                </DataTable>
              </div>

              <div class="flex justify-content-end mt-4">
                <Button
                  :label="t('confirm_backups_import')"
                  icon="pi pi-exclamation-triangle"
                  severity="danger"
                  :disabled="!hasSelection"
                  @click="openConfirm"
                />
              </div>
            </template>

            <Message v-if="importError" severity="error" :closable="false" class="mt-4">
              {{ importError }}
            </Message>
            <Button
              v-if="failedBackup"
              :label="t('download_backup')"
              icon="pi pi-download"
              severity="secondary"
              class="mt-2"
              @click="downloadBackup(failedBackup)"
            />

            <template v-if="lastResult">
              <Message severity="success" :closable="false" class="mt-4">
                {{ t('backups_imported_success') }}
              </Message>
              <h4 class="mt-3">{{ t('backups_imported_counts') }}</h4>
              <DataTable :value="importedRows" size="small" stripedRows>
                <Column field="label" :header="t('backups_section')" />
                <Column field="count" :header="t('backups_incoming')" />
              </DataTable>
              <Button
                v-if="lastResult.backup_file"
                :label="t('download_backup')"
                icon="pi pi-download"
                class="mt-3"
                @click="downloadBackup(lastResult.backup_file)"
              />
            </template>
          </template>
        </Card>
      </div>
    </div>

    <Dialog
      v-model:visible="confirmVisible"
      modal
      :header="t('confirm_backups_import')"
      :style="{ width: '42rem' }"
      :breakpoints="{ '1199px': '60vw', '575px': '90vw' }"
    >
      <Message severity="error" :closable="false">
        {{ t('backups_replace_warning') }}
      </Message>
      <ul class="mt-3">
        <li v-for="label in selectionLabels" :key="label">{{ label }}</li>
      </ul>
      <p class="text-sm text-color-secondary">{{ t('backups_notice') }}</p>
      <div class="flex align-items-center gap-2 mt-4">
        <Checkbox v-model="acknowledged" binary inputId="backups-ack" />
        <label for="backups-ack">{{ t('backups_confirm_ack') }}</label>
      </div>
      <template #footer>
        <Button :label="t('cancel')" severity="secondary" @click="confirmVisible = false" />
        <Button
          :label="t('confirm_backups_import')"
          severity="danger"
          :loading="importing"
          :disabled="!acknowledged"
          @click="importBackup"
        />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import axios, { type AxiosError } from 'axios'
import { computed, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import Button from 'primevue/button'
import Card from 'primevue/card'
import Checkbox from 'primevue/checkbox'
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import Dialog from 'primevue/dialog'
import FileUpload from 'primevue/fileupload'
import Message from 'primevue/message'
import ProgressSpinner from 'primevue/progressspinner'
import auth from '../services/auth'

const { t } = useI18n()
const toast = useToast()

const backendUrl = `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}`
const maxUploadBytes = 512 * 1024 * 1024
const minimumPreviewDurationMs = 2000

interface BackupCounts {
  materials: number
  material_entries: number
  products: number
  categories: number
  sales: number
  orders: number
  images: number
}

interface BackupPreview {
  sections: BackupCounts
  current: BackupCounts
  has_stock: boolean
  has_images: boolean
}

interface BackupImportResult {
  backup_file?: string
  imported: BackupCounts
}

interface BackupFileSelectEvent {
  files?: File[]
}

interface BackupErrorResponse {
  error?: string
  backup_file?: string
}

const exporting = ref(false)
const previewLoading = ref(false)
const importing = ref(false)
const selectedFile = ref<File>()
const preview = ref<BackupPreview>()
const lastResult = ref<BackupImportResult>()
const importError = ref('')
const failedBackup = ref('')
const confirmVisible = ref(false)
const acknowledged = ref(false)

const sections = reactive({
  materials: true,
  products: true,
  categories: true,
  sales: false,
})
const includeStock = ref(true)
const includeImages = ref(true)

watch(
  () => sections.materials,
  (value) => {
    if (!value) {
      sections.products = false
      sections.categories = false
      includeStock.value = false
      includeImages.value = false
    }
  },
)

watch(
  () => sections.products,
  (value) => {
    if (value) {
      sections.materials = true
    } else {
      sections.categories = false
      includeImages.value = false
    }
  },
)

watch(
  () => sections.categories,
  (value) => {
    if (value) {
      sections.products = true
      sections.materials = true
    }
  },
)

watch(includeStock, (value) => {
  if (value) {
    sections.materials = true
  }
})

watch(includeImages, (value) => {
  if (value) {
    sections.products = true
    sections.materials = true
  }
})

const hasSelection = computed(
  () => sections.materials || sections.products || sections.categories || sections.sales,
)

const selectionLabels = computed(() => {
  const labels: string[] = []
  if (sections.materials) labels.push(t('backups_section_materials'))
  if (sections.materials && includeStock.value) labels.push(t('backups_section_stock'))
  if (sections.products) labels.push(t('backups_section_products'))
  if (sections.products && includeImages.value) labels.push(t('backups_section_images'))
  if (sections.categories) labels.push(t('backups_section_categories'))
  if (sections.sales) labels.push(t('sales_orders'))
  return labels
})

const countRows = computed(() => {
  if (!preview.value) return []
  const incoming = preview.value.sections
  const current = preview.value.current
  return [
    { label: t('backups_section_materials'), incoming: incoming.materials, current: current.materials },
    { label: t('backups_section_stock'), incoming: incoming.material_entries, current: current.material_entries },
    { label: t('backups_section_products'), incoming: incoming.products, current: current.products },
    { label: t('backups_section_images'), incoming: incoming.images, current: current.images },
    { label: t('backups_section_categories'), incoming: incoming.categories, current: current.categories },
    {
      label: t('sales_orders'),
      incoming: `${incoming.sales} / ${incoming.orders}`,
      current: `${current.sales} / ${current.orders}`,
    },
  ]
})

const importedRows = computed(() => {
  if (!lastResult.value) return []
  const imported = lastResult.value.imported
  return [
    { label: t('backups_section_materials'), count: imported.materials },
    { label: t('backups_section_stock'), count: imported.material_entries },
    { label: t('backups_section_products'), count: imported.products },
    { label: t('backups_section_images'), count: imported.images },
    { label: t('backups_section_categories'), count: imported.categories },
    { label: t('sales_orders'), count: `${imported.sales} / ${imported.orders}` },
  ]
})

const authHeaders = () => ({
  Authorization: `Bearer ${auth.accessToken.value}`,
})

const apiError = (error: unknown): string => {
  const axiosError = error as AxiosError<BackupErrorResponse | string>
  const data = axiosError.response?.data
  if (typeof data === 'string' && data.trim() !== '') return data
  if (data && typeof data === 'object' && data.error) return data.error
  return axiosError.message || t('backups_unknown_error')
}

const downloadBlob = (data: Blob, filename: string) => {
  const url = window.URL.createObjectURL(data)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  link.remove()
  window.URL.revokeObjectURL(url)
}

const exportBackup = () => {
  exporting.value = true
  axios
    .get(`${backendUrl}/api/backups/export`, {
      headers: authHeaders(),
      responseType: 'blob',
    })
    .then((response) => {
      const stamp = new Date().toISOString().replace(/[:.]/g, '-').slice(0, 19)
      downloadBlob(response.data, `nutrix-backup-${stamp}.zip`)
      toast.add({ severity: 'success', summary: t('done'), detail: t('backups_export_ready'), life: 3000 })
    })
    .catch(() => {
      toast.add({ severity: 'error', summary: t('error'), detail: t('backups_export_failed'), life: 3000 })
    })
    .finally(() => {
      exporting.value = false
    })
}

const onFileSelect = (event: BackupFileSelectEvent) => {
  const file = event.files?.[0]
  if (!file) return
  if (file.size > maxUploadBytes) {
    toast.add({ severity: 'error', summary: t('error'), detail: t('backups_file_too_large'), life: 4000 })
    return
  }

  selectedFile.value = file
  preview.value = undefined
  lastResult.value = undefined
  importError.value = ''
  failedBackup.value = ''
  acknowledged.value = false
  loadPreview()
}

const waitForMinimumPreviewDuration = async (startedAt: number) => {
  const elapsed = Date.now() - startedAt
  if (elapsed < minimumPreviewDurationMs) {
    await new Promise((resolve) => setTimeout(resolve, minimumPreviewDurationMs - elapsed))
  }
}

const loadPreview = async () => {
  if (!selectedFile.value) return

  const startedAt = Date.now()
  const form = new FormData()
  form.append('file', selectedFile.value)

  previewLoading.value = true
  importError.value = ''
  failedBackup.value = ''
  acknowledged.value = false

  try {
    const response = await axios.post(`${backendUrl}/api/backups/import/preview`, form, {
      headers: authHeaders(),
    })
    await waitForMinimumPreviewDuration(startedAt)
    preview.value = response.data.data
  } catch (error) {
    await waitForMinimumPreviewDuration(startedAt)
    preview.value = undefined
    importError.value = apiError(error)
    toast.add({ severity: 'error', summary: t('error'), detail: t('backups_preview_failed'), life: 4000 })
  } finally {
    previewLoading.value = false
  }
}

const openConfirm = () => {
  if (!hasSelection.value) {
    toast.add({ severity: 'warn', summary: t('error'), detail: t('backups_no_sections'), life: 3000 })
    return
  }
  acknowledged.value = false
  confirmVisible.value = true
}

const importBackup = () => {
  if (!selectedFile.value) return

  const form = new FormData()
  form.append('file', selectedFile.value)

  const selected: string[] = []
  if (sections.materials) selected.push('materials')
  if (sections.products) selected.push('products')
  if (sections.categories) selected.push('categories')
  if (sections.sales) selected.push('sales')

  form.append('sections', selected.join(','))
  form.append('include_stock', includeStock.value ? 'true' : 'false')
  form.append('include_images', includeImages.value ? 'true' : 'false')

  importing.value = true
  importError.value = ''
  failedBackup.value = ''

  axios
    .post(`${backendUrl}/api/backups/import`, form, { headers: authHeaders() })
    .then((response) => {
      lastResult.value = response.data.data
      preview.value = undefined
      selectedFile.value = undefined
      confirmVisible.value = false
      acknowledged.value = false
      toast.add({ severity: 'success', summary: t('done'), detail: t('backups_imported_success'), life: 4000 })
    })
    .catch((error) => {
      importError.value = apiError(error)
      const axiosError = error as AxiosError<BackupErrorResponse>
      const data = axiosError.response?.data
      failedBackup.value = data && typeof data === 'object' ? data.backup_file || '' : ''
      toast.add({ severity: 'error', summary: t('error'), detail: t('backups_import_failed'), life: 4000 })
    })
    .finally(() => {
      importing.value = false
    })
}

const downloadBackup = (filename: string) => {
  axios
    .get(`${backendUrl}/api/backups/files/${encodeURIComponent(filename)}`, {
      headers: authHeaders(),
      responseType: 'blob',
    })
    .then((response) => {
      downloadBlob(response.data, filename)
    })
    .catch(() => {
      toast.add({ severity: 'error', summary: t('error'), detail: t('backups_download_failed'), life: 3000 })
    })
}
</script>
