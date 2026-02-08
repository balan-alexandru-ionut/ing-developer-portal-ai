<script setup lang="ts">

import {useGenerationStore} from "@/stores/generation.ts";
import {codeToHtml} from "shiki";
import {ref, watchEffect} from "vue";
import type {File} from "@/types/files.ts";
import FileTree from "@/components/FileTree.vue";

const generationStore = useGenerationStore()
const renderedFiles = ref<Map<string, string>>(new Map())
const visibleFile = ref('')

async function renderFile(file: File): Promise<string> {
  return codeToHtml(file.code, {
    lang: generationStore.language.toLowerCase(),
    theme: 'catppuccin-mocha',
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
      <div class="p-4 rounded-lg" v-html="visibleFile"></div>
    </div>
  </div>
</template>

<style scoped>

</style>
