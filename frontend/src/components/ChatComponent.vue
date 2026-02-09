<script lang="ts" setup>

import {onMounted, ref} from "vue";
import {useChatStore} from "@/stores/chat.ts";
import {marked} from "marked";

const chatStore = useChatStore();

const prompt = ref('');
const processing = ref(false);

async function startChat(): Promise<Response> {
  const res = await fetch('/api/chat/start', {method: 'POST'});
  if (res.ok) {
    console.log('Chat session started');
  }
  return res;
}

function sendMessage() {
  chatStore.userMessages.push(prompt.value);
  processing.value = true;
  fetch('/api/chat/message', {method: 'POST', body: JSON.stringify({prompt: prompt.value})})
    .then(res => res.json())
    .then(data => data.message)
    .then(msg => marked.parse(msg, {async: false}))
    .then(html => chatStore.modelMessages.push(html))
    .finally(() => processing.value = false)
    .finally(() => prompt.value = '');
}

onMounted(async () => {
  await startChat()
})
</script>

<template>
  <div class="flex flex-col gap-4 p-4 bg-base-100 rounded-lg">
    <div class="flex flex-col gap-2 overflow-y-auto">
      <template v-for="(_, i) in chatStore.userMessages" :key="i">
        <div class="chat chat-end">
          <div class="chat-bubble chat-bubble-primary">{{ chatStore.userMessages[i] }}</div>
        </div>
        <div v-if="chatStore.modelMessages[i]" class="chat chat-start">
          <div class="chat-bubble chat-bubble-secondary">
            <div v-html="chatStore.modelMessages[i]"/>
          </div>
        </div>
        <div v-else-if="processing" class="chat chat-start">
          <span class="loading loading-infinity loading-md"></span>
          <div class="chat-bubble chat-bubble-secondary">Processing...</div>
        </div>
      </template>
    </div>
    <div class="flex gap-4">
      <input v-model="prompt" class="input rounded-xl w-full" placeholder="Type your message here..."
             type="text" @keyup.enter="sendMessage"/>
      <button class="btn btn-primary" @click="sendMessage">Send</button>
    </div>
  </div>
</template>

<style scoped>

</style>
