<script setup lang="ts">
import ing from '@/assets/ing.svg'
</script>

<template>
  <nav class="navbar bg-base-100 shadow-xl">
    <RouterLink to="/" class="flex items-center gap-4 px-4">
      <img :src="ing" alt="ING logo" class="w-12" />
      <span class="text-xl font-bold">ING Sandbox Playground</span>
    </RouterLink>
  </nav>
</template>

<style scoped></style>
