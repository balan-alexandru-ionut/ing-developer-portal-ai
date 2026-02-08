<script setup lang="ts">
import { computed } from "vue";
import FileTreeNode from "./FileTreeNode.vue";

const props = defineProps({
  filePaths: { type: Array<string>, required: true },
});

export type TreeNode = {
  name: string;
  path: string;
  children: TreeNode[];
  isFile: boolean;
};

function buildTree(paths: string[]): TreeNode[] {
  const root: TreeNode = { name: "", path: "", children: [], isFile: false };

  for (const filePath of paths) {
    const parts = filePath.split("/");
    let current = root;

    for (let i = 0; i < parts.length; i++) {
      const part = parts[i]!;
      const isFile = i === parts.length - 1;
      let child = current.children.find((c) => c.name === part);

      if (!child) {
        child = {
          name: part,
          path: parts.slice(0, i + 1).join("/"),
          children: [],
          isFile,
        };
        current.children.push(child);
      }

      current = child;
    }
  }

  function sortTree(nodes: TreeNode[]) {
    nodes.sort((a, b) => {
      if (a.isFile !== b.isFile) return a.isFile ? 1 : -1;
      return a.name.localeCompare(b.name);
    });
    for (const node of nodes) {
      sortTree(node.children);
    }
  }

  sortTree(root.children);
  return root.children;
}

const tree = computed(() => buildTree(props.filePaths));

const emit = defineEmits<{
  fileOpen: [filePath: string];
}>();
</script>

<template>
  <div class="flex flex-col p-4 bg-base-100 rounded-lg text-sm">
    <FileTreeNode
      v-for="node in tree"
      :key="node.path"
      :node="node"
      :depth="0"
      @file-open="(path: string) => emit('fileOpen', path)"
    />
  </div>
</template>
