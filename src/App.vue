<template>
  <div class="sticky top-0 w-full flex-none border-b border-slate-200 bg-white z-50">
    <div class="max-w-7xl mx-auto py-4">
      <div class="flex flex-row gap-4 text-xl">
        <div class="flex-none font-bold">DNS Detoxifier</div>
        <div class="flex-1 flex flex-row justify-end space-x-6">
          <router-link to="/data/stats">Dashboard</router-link>
          <router-link to="/about">About</router-link>
        </div>
      </div>
    </div>
  </div>
  <div>
    <div class="max-w-7xl mx-auto">
      <router-view></router-view>
    </div>
  </div>
</template>

<script>

export default {
  name: 'App',
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
}
</style>
