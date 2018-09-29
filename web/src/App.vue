<template>
  <div id="app">
    <!-- <img alt="Vue logo" src="./assets/logo.png"> -->
    <HelloWorld msg="Welcome to Your Vue.js App"/>
  </div>
</template>

<script>
import HelloWorld from './components/HelloWorld.vue'

export default {
  name: 'app',
  components: {
    HelloWorld
  }
}
</script>

<style>
@import url('https://fonts.googleapis.com/css?family=PT+Sans');
@import url('https://fonts.googleapis.com/css?family=Noto+Serif+KR');

body {
  /* background-image: linear-gradient(to bottom, white, #9e8a8a); */
  height: 1000px;
}

#app {
  font-family: 'Noto Serif KR', sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
  max-width: 720px;
  margin: auto;
}
</style>
