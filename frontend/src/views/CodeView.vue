<script setup lang="ts">

import {useGenerationStore} from "@/stores/generation.ts";
import {codeToHtml} from "shiki";
import {ref, watchEffect} from "vue";
import type {File} from "@/types/files.ts";
import FileTree from "@/components/FileTree.vue";
import {Download, Play, Sparkles} from "lucide-vue-next";
import ChatComponent from "@/components/ChatComponent.vue";

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
  <div class="flex flex-col bg-base-300 h-screen gap-4 pt-4">
    <div class="flex gap-8 px-12 h-screen w-full">
      <FileTree :filePaths="generationStore.files.map(f => f.filePath)"
                @fileOpen="path => fileOpen(path)"/>
      <div class="bg-base-100 border border-base-300 size-full overflow-scroll">
        <div class="flex justify-end gap-4 w-full p-4 shrink">
          <button class="btn btn-primary" @click="downloadArchive">
            <Download/>
            Download
          </button>
          <div class="drawer drawer-end w-fit">
            <input id="chat-drawer" class="drawer-toggle" type="checkbox"/>
            <div class="drawer-content">
              <label class="btn btn-accent drawer-button" for="chat-drawer">
                Chat
                <Sparkles/>
              </label>
            </div>
            <div class="drawer-side">
              <label aria-label="close sidebar" class="drawer-overlay" for="chat-drawer"></label>
              <div class="menu bg-base-100 min-h-full w-11/12 p-4 flex flex-col gap-4">
                <ChatComponent/>
              </div>
            </div>
          </div>
          <button class="btn btn-disabled">
            Run
            <Play/>
          </button>
        </div>
        <div class="p-4 rounded-lg" v-html="visibleFile"></div>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>
