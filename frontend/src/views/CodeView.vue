<script setup lang="ts">

import {useGenerationStore} from "@/stores/generation.ts";
import {codeToHtml} from "shiki";
import {ref, watchEffect} from "vue";
import type {File} from "@/types/files.ts";
import FileTree from "@/components/FileTree.vue";
import {Download, Play} from "lucide-vue-next";

const generationStore = useGenerationStore()
const renderedFiles = ref<Map<string, string>>(new Map())
const visibleFile = ref('')

async function renderFile(file: File): Promise<string> {
  return codeToHtml(file.code, {
    lang: generationStore.language.toLowerCase(),
    theme: 'catppuccin-mocha',
  })
}

function downloadArchive() {
  fetch('/api/generate/archive', {method: 'GET', headers: {'Accept': 'application/zip'}})
    .then(response => response.blob())
    .then(blob => {
      const url = globalThis.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = 'archive.zip'
      document.body.appendChild(a)
      a.click()
      a.remove()
      globalThis.URL.revokeObjectURL(url)
    })
}

watchEffect(async () => {
  const entries = new Map<string, string>()
  console.log(generationStore.files)
  for (const file of generationStore.files) {
    const html = await  renderFile(file)
    entries.set(file.filePath, html)
  }
  renderedFiles.value = entries
  visibleFile.value = renderedFiles.value.get(generationStore.files[0]?.filePath as string) ?? ''
})

function fileOpen(filePath: string) {
  visibleFile.value = renderedFiles.value.get(filePath) ?? ''
}

</script>

<template>
  <div class="flex gap-8 p-12 bg-base-300 h-screen w-full">
    <FileTree @fileOpen="path => fileOpen(path)" :filePaths="generationStore.files.map(f => f.filePath)" />
    <div class="mockup-window bg-base-100 border border-base-300 size-full overflow-scroll">
      <div class="flex justify-end gap-4 w-full p-4">
        <button class="btn btn-primary" @click="downloadArchive">
          <Download />
          Download
        </button>
        <button class="btn btn-disabled">
          Run
          <Play />
        </button>
      </div>
      <div class="p-4 rounded-lg" v-html="visibleFile"></div>
    </div>
  </div>
</template>

<style scoped>

</style>
