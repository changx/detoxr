<template>
    <div id="headerBar" class="flex flex-row h-12 w-full items-stretch">
        <div class="flex-none w-32"></div>
        <div class="flex-grow self-center">
            <router-link to="/home">home</router-link>
            <router-link to="/about">about</router-link>
        </div>
        <div class="flex-none w-32"></div>
    </div>
    <div class="container">
      <router-view></router-view>
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
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
#headerBar {
    background-color:lightseagreen;
    color: white;
    font-weight: bold;
}
</style>
