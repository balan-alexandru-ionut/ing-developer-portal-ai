<script setup lang="ts">
import {CircleDashed, CircleCheck} from "lucide-vue-next";
import {onMounted, ref} from "vue";
import {useRouter} from "vue-router";

type GenerationStatusResponse = {
  time: string,
  status: 'not_started' | 'generating' | 'generated' | 'formatting' | 'done'
}

const status = ref<GenerationStatusResponse['status']>('not_started')

const POLL_INTERVAL = 1500

const router = useRouter()

async function pollStatus() {
  const res = await fetch('/api/generate/status', {method: 'GET'})
  const data = await res.json() as GenerationStatusResponse
  status.value = data.status
  if (data.status === 'done') {
    setTimeout(() => router.push('/code'), 1500)
  } else {
    setTimeout(pollStatus, POLL_INTERVAL)
  }
}

onMounted(pollStatus)


</script>

<template>
<div class="flex justify-center w-full">
  <ul class="timeline">
    <li>
      <div class="timeline-middle">
        <span v-if="['not_started', 'generating'].includes(status)" class="loading loading-spinner loading-lg text-primary"></span>
        <CircleCheck v-else class="text-success" />
      </div>
      <div class="timeline-end timeline-box text-base">
        <span v-if="['not_started', 'generating'].includes(status)">Generating code</span>
        <span v-else>Code generated</span>
      </div>
      <hr class="w-32" />
    </li>

    <li :class="['done', 'formatting'].includes(status) ? 'text-base-content' : 'text-gray-500'">
      <hr class="w-32" />
      <div class="timeline-middle">
        <span v-if="status === 'formatting'" class="loading loading-spinner loading-lg text-primary"></span>
        <CircleDashed v-if="!['formatting', 'done'].includes(status)" />
        <CircleCheck class="text-success" v-if="status === 'done'" />
      </div>
      <div class="timeline-end timeline-box text-base">
        <span v-if="status !== 'done'">Formatting code</span>
        <span v-if="status === 'done'">Code formatted</span>
      </div>
    </li>
  </ul>
</div>
</template>

<style scoped>

</style>
