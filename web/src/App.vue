
<template>
  <n-config-provider :locale="zhCN" :date-locale="dateZhCN">
    <div class="h-screen w-full">
      <router-view />
    </div>
  </n-config-provider>
</template>

<script setup lang="ts">
import {
  zhCN,
  dateZhCN,
  NConfigProvider
} from 'naive-ui'
</script>

<style>
/* Reset */
body {
    margin: 0;
    padding: 0;
    font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Fira Sans', 'Droid Sans', 'Helvetica Neue', sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
}

</style>


