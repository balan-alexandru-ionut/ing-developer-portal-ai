<script setup>
import { ref, computed, nextTick } from 'vue'
import CodeBlock from './CodeBlock.vue'

const loading = ref(false)
const error = ref('')
const raw = ref('')
const parsed = ref(null)
const activeTab = ref(0)
const wrap = ref(false) // toggle for line wrapping

function reset() {
  error.value = ''
  raw.value = ''
  parsed.value = null
  activeTab.value = 0
}

async function generate() {
  reset()
  loading.value = true
  try {
    const res = await fetch('/api/generate', { method: 'GET' })
    if (!res.ok) {
      const text = await res.text().catch(() => '')
      throw new Error(`Request failed (${res.status}): ${text}`)
    }

    const data = await res.json()
    raw.value = typeof data?.code === 'string' ? data.code : JSON.stringify(data, null, 2)

    // Try to parse the raw string to detect files[]
    try {
      parsed.value = JSON.parse(raw.value)
    } catch {
      parsed.value = null
    }

    await nextTick()
  } catch (e) {
    error.value = e?.message ?? 'Unknown error'
  } finally {
    loading.value = false
  }
}

async function copyAll() {
  if (!raw.value) return
  await navigator.clipboard.writeText(raw.value)
}

function downloadAll() {
  if (!raw.value) return
  const blob = new Blob([raw.value], { type: 'application/json;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'generated.json'
  a.click()
  URL.revokeObjectURL(url)
}

// Detect expected file structure
const files = computed(() => {
  if (!parsed.value) return null

  // Preferred: { files: [{ filePath, code }] }
  if (Array.isArray(parsed.value.files)) {
    return parsed.value.files.map((f, idx) => ({
      filePath: f.filePath ?? f.path ?? f.name ?? `file-${idx}`,
      code:
        typeof f.code === 'string'
          ? f.code
          : typeof f.content === 'string'
            ? f.content
            : f.code
              ? JSON.stringify(f.code, null, 2)
              : '',
    }))
  }

  // Alternative: object keyed by path { "src/A.java": "...", ... }
  if (typeof parsed.value === 'object' && parsed.value !== null) {
    const keys = Object.keys(parsed.value)
    const looksLikeFiles = keys.some((k) => k.includes('/') || k.includes('.'))
    if (looksLikeFiles) {
      return keys.map((k) => {
        const v = parsed.value[k]
        return {
          filePath: k,
          code: typeof v === 'string' ? v : JSON.stringify(v, null, 2),
        }
      })
    }
  }

  return null
})

function languageFor(filePath) {
  if (!filePath) return 'json'
  const p = filePath.toLowerCase()
  if (p.endsWith('.java')) return 'java'
  if (p.endsWith('.json')) return 'json'
  if (p.endsWith('.xml') || p.endsWith('pom.xml')) return 'markup'
  if (p.endsWith('.md') || p.includes('readme')) return 'markup' // can switch to markdown renderer later
  if (p.endsWith('.js')) return 'javascript'
  if (p.endsWith('.ts')) return 'typescript'
  return 'json'
}
</script>

<template>
  <div class="min-h-screen bg-base-200">
    <!-- Sticky header -->
    <div class="navbar bg-base-100 shadow sticky top-0 z-10">
      <div class="flex-1 px-4">
        <span class="text-2xl md:text-3xl font-bold">Generate java code for ShowCase API</span>
      </div>

      <div class="flex-none gap-2 px-4 items-center">
        <!-- Generate: fully disabled while loading, pretty loader, stable width -->
        <button
          class="btn btn-primary min-w-[130px] flex items-center justify-center gap-2"
          :class="{ 'btn-disabled': loading }"
          :disabled="loading"
          @click="generate"
        >
          <span v-if="loading" class="loading loading-ring loading-sm"></span>
          <span>{{ loading ? 'Generating…' : 'Generate' }}</span>
        </button>

        <button class="btn" :disabled="!raw || loading" @click="copyAll">Copy</button>

        <button class="btn" :disabled="!raw || loading" @click="downloadAll">Download</button>
      </div>
    </div>

    <div class="container mx-auto px-4 py-6 max-w-7xl min-w-0">
      <!-- Errors -->
      <div v-if="error" class="alert alert-error mb-4">
        <span>{{ error }}</span>
      </div>

      <!-- Loading state with skeletons -->
      <div v-if="loading" class="space-y-4">
        <div class="skeleton h-10 w-64"></div>
        <div class="skeleton h-80 w-full"></div>
      </div>

      <!-- Multi-file output with horizontally scrollable tabs -->
      <div v-if="!loading && files?.length" class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <div role="tablist" class="tabs tabs-boxed overflow-x-auto">
            <button
              v-for="(f, i) in files"
              :key="f.filePath + i"
              role="tab"
              class="tab whitespace-nowrap"
              :class="{ 'tab-active': activeTab === i }"
              @click="activeTab = i"
              :title="f.filePath"
            >
              {{ f.filePath }}
            </button>
          </div>

          <!-- Warn if current file has no code -->
          <div v-if="!files[activeTab].code" class="alert alert-warning mt-4">
            <span
              >No code content found for <b>{{ files[activeTab].filePath }}</b
              >.</span
            >
          </div>

          <div class="mt-4 min-w-0">
            <CodeBlock
              :code="files[activeTab].code"
              :language="languageFor(files[activeTab].filePath)"
              :wrap="wrap"
            />
          </div>
        </div>
      </div>

      <!-- Fallback: single JSON block -->
      <div v-else-if="!loading && raw" class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <p class="text-sm opacity-70 mb-2">Response</p>
          <CodeBlock :code="raw" language="json" :wrap="wrap" />
        </div>
      </div>

      <!-- Empty state -->
      <div v-else class="text-center opacity-60 mt-12">
        Click <strong>Generate</strong> to fetch code.
      </div>
    </div>
  </div>
</template>
