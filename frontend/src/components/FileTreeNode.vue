<script setup lang="ts">
import { ref } from "vue";
import type { TreeNode } from "./FileTree.vue";
import {Folder, FolderOpen, FileCode} from "lucide-vue-next";

defineProps<{
  node: TreeNode;
  depth: number;
}>();

const emit = defineEmits<{
  fileOpen: [filePath: string];
}>();

const expanded = ref(true);

function toggle() {
  expanded.value = !expanded.value;
}
</script>

<template>
  <div>
    <div
      class="flex items-center gap-1 py-0.5 px-1 rounded cursor-pointer hover:bg-accent hover:text-accent-content select-none"
      :style="{ paddingLeft: depth * 16 + 'px' }"
      @click="node.isFile ? emit('fileOpen', node.path) : toggle()"
    >
      <span v-if="!node.isFile" class="w-2 text-center text-sm mr-4">
        <FolderOpen v-if="expanded" />
        <Folder v-else />
      </span>
      <span v-else class="w-2 text-center text-sm mr-4">
        <FileCode />
      </span>
      <span class="text-base">{{ node.name }}</span>
    </div>
    <template v-if="!node.isFile && expanded">
      <!-- Vue SFCs can reference themselves by filename for recursion -->
      <FileTreeNode
        v-for="child in node.children"
        :key="child.path"
        :node="child"
        :depth="depth + 1"
        @file-open="(path: string) => emit('fileOpen', path)"
      />
    </template>
  </div>
</template>
