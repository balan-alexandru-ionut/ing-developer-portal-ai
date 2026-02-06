import {defineStore} from "pinia";
import {ref} from "vue";
import type {File} from "@/types/files.ts";

export const useGenerationStore = defineStore('generation', () => {
  const language = ref('')
  const api = ref('')
  const files = ref<File[]>([])

  return { language, api, files }
})
