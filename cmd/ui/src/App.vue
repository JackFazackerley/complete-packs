<template>
  <div id="app">
    <Sizes/>
  </div>
</template>

<script>
import Sizes from './components/Sizes.vue'
export default {
  name: 'app',
  components: {
    Sizes
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>