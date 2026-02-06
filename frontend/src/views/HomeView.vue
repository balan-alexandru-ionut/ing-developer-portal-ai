<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { ArrowRightIcon, Sparkles } from 'lucide-vue-next'
import {useRouter} from "vue-router";
import {useGenerationStore} from "@/stores/generation.ts";
import type {File} from "@/types/files.ts";

const router = useRouter()
const generationStore = useGenerationStore()

const language = ref('Go')
const api = ref('Showcase API')

const supportedLanguages = ['Go', 'Java', 'Python', 'JavaScript', 'TypeScript', 'Rust']
const supportedApis = [
  'Showcase API',
  'Account Information API',
  'Confirmation Availability of Funds API',
  'Payment Initiation API',
  'Real-time Account Reporting API',
]

const languageChosen = ref(false)
const apiChosen = ref(false)

let languageIndex = 0
let apiIndex = 0

let languageInterval: ReturnType<typeof setInterval> | null = null
let apiInterval: ReturnType<typeof setInterval> | null = null

function startLanguageCycling() {
  languageInterval = setInterval(() => {
    languageIndex = (languageIndex + 1) % supportedLanguages.length
    language.value = supportedLanguages[languageIndex]!!
  }, 2000)
}

function startApiCycling() {
  apiInterval = setInterval(() => {
    apiIndex = (apiIndex + 1) % supportedApis.length
    api.value = supportedApis[apiIndex]!!
  }, 3500)
}

function chooseLanguage(chosenLanguage: string) {
  language.value = chosenLanguage
  languageChosen.value = true
  if (languageInterval) {
    clearInterval(languageInterval)
    languageInterval = null
  }
}

function chooseApi(chosenApi: string) {
  api.value = chosenApi
  apiChosen.value = true
  if (apiInterval) {
    clearInterval(apiInterval)
    apiInterval = null
  }
}

function toggleLanguageCycling() {
  if (languageChosen.value) return

  if (languageInterval) {
    clearInterval(languageInterval)
    languageInterval = null
    return
  }

  startLanguageCycling()
}

function toggleApiCycling() {
  if (apiChosen.value) return

  if (apiInterval) {
    clearInterval(apiInterval)
    apiInterval = null
    return
  }

  startApiCycling()
}

function getGeneratedCode() {
  const prompt =
    'Build me a fully working ' + language.value + ' application that calls ' + api.value
  fetch('api/generate', { method: 'GET' })
    .then((res) => res.json())
    .then((data: {files: File[]}) => {
      generationStore.language = language.value
      generationStore.api = api.value
      generationStore.files = data.files
    })
    .then(code => router.push({name: 'code', params: {language: language.value, api: api.value}}))
}

onMounted(() => {
  startLanguageCycling()
  startApiCycling()
})

onBeforeUnmount(() => {
  if (languageInterval) clearInterval(languageInterval)
  if (apiInterval) clearInterval(apiInterval)
})
</script>

<template>
  <div class="flex flex-col justify-center h-screen">
    <div class="flex justify-center w-full mb-16">
      <h1 class="text-5xl font-bold text-accent">Generate your starter application</h1>
    </div>

    <div class="flex justify-center items-center w-full">
      <span class="text-2xl">Build me a fully working</span>
      <button
        ref="langSelectEl"
        class="cycling-btn btn btn-ghost text-2xl text-primary font-bold px-4"
        popovertarget="language-popover"
        style="anchor-name: --language-anchor"
      >
        {{ language }}
      </button>
      <ul
        class="dropdown menu w-36 rounded-box bg-base-100 shadow-sm text-base"
        popover
        id="language-popover"
        style="position-anchor: --language-anchor"
        @toggle="toggleLanguageCycling"
      >
        <li v-for="lang in supportedLanguages" @click="chooseLanguage(lang)">{{ lang }}</li>
      </ul>
      <span class="text-2xl">application that calls</span>
      <button
        ref="apiSelectEl"
        class="cycling-btn btn btn-ghost text-2xl text-primary font-bold px-4"
        popovertarget="api-popover"
        style="anchor-name: --api-anchor"
      >
        {{ api }}
      </button>
      <ul
        class="dropdown menu w-80 rounded-box bg-base-100 shadow-sm text-base"
        popover
        id="api-popover"
        style="position-anchor: --api-anchor"
        @toggle="toggleApiCycling"
      >
        <li v-for="api in supportedApis" @click="chooseApi(api)">{{ api }}</li>
      </ul>
    </div>

    <div class="flex justify-center items-center w-full mt-8">
      <button
        class="btn btn-primary btn-lg mx-auto"
        :class="!languageChosen || !apiChosen ? 'btn-disabled' : ''"
        @click="getGeneratedCode"
      >
        <Sparkles />
        Generate
      </button>
    </div>

    <div class="flex justify-center items-center gap-2 w-full mt-8">
      <ArrowRightIcon class="text-info" />
      <span class="text-lg text-info"
        >Click on the language or API buttons to select your preferences.</span
      >
    </div>
  </div>
</template>

<style scoped>
ul li {
  cursor: pointer;
  padding: 0.5rem 1rem;
  margin: 0.5rem 0;
  border-radius: 0.5rem;

  &:hover {
    background-color: var(--color-accent);
    color: var(--color-accent-content);
    transition: all 0.2s ease-in-out;
  }
}

.cycling-btn {
  transition: width 0.5s ease;
  overflow: hidden;
  white-space: nowrap;
}
</style>
