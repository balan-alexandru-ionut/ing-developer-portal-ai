<!-- src/components/CodeBlock.vue -->
<script setup>
import { onMounted, onUpdated, ref, watch } from 'vue'
import Prism from 'prismjs'

// Load the languages you need:
import 'prismjs/components/prism-json'
import 'prismjs/components/prism-java'
import 'prismjs/components/prism-markup'
import 'prismjs/components/prism-javascript'
import 'prismjs/components/prism-typescript'

const props = defineProps({
  code: { type: String, required: true },
  language: { type: String, default: 'json' },
  wrap: { type: Boolean, default: false },
})

const codeRef = ref(null)

const highlight = () => {
  // If code string is huge, Prism may be slow; still okay for now
  if (codeRef.value) Prism.highlightElement(codeRef.value)
}

onMounted(highlight)
onUpdated(highlight)

// Re-highlight when content changes
watch(
  () => props.code,
  () => {
    // Avoid flashing; re-run highlight next tick
    requestAnimationFrame(highlight)
  },
)
</script>

<template>
  <div class="w-full min-w-0">
    <div class="rounded-lg border border-base-300 bg-base-100">
      <div class="px-3 py-2 text-xs opacity-60 flex items-center gap-2 border-b border-base-300">
        <span>Preview</span>
        <span class="badge badge-ghost badge-sm">{{ language }}</span>
        <span class="ml-auto">
          <span v-if="wrap" class="badge badge-outline badge-sm">Wrap: ON</span>
          <span v-else class="badge badge-outline badge-sm">Wrap: OFF</span>
        </span>
      </div>

      <div class="max-h-[70vh] overflow-x-auto overflow-y-auto">
        <!-- IMPORTANT: use pre + code; control wrapping via classes -->
        <pre class="m-0 p-4" :class="wrap ? 'whitespace-pre-wrap break-words' : 'whitespace-pre'">
<code ref="codeRef" :class="`language-${language}`">{{ code }}</code>
        </pre>
      </div>
    </div>
  </div>
</template>
