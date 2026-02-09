import {defineStore} from "pinia";
import {ref} from "vue";

export const useChatStore = defineStore('chat', () => {
  const userMessages = ref<string[]>([])
  const modelMessages = ref<string[]>([])

  return {userMessages, modelMessages}
})
