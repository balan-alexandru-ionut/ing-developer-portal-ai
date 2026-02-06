<script setup lang="ts">

import {useGenerationStore} from "@/stores/generation.ts";
import {codeToHtml} from "shiki";
import {ref, watchEffect} from "vue";

const generationStore = useGenerationStore()

const html = ref('')

watchEffect(async () => {
  const file = generationStore.files[3]
  if (file && file.code) {
    html.value = await codeToHtml(file.code, {
      lang: generationStore.language.toLowerCase(),
      theme: 'catppuccin-mocha'
    })
  }
})
</script>

<template>
  <div v-html="html"></div>
</template>

<style scoped>

</style>
