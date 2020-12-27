<template>
  <div id="app">
    <router-view></router-view>
  </div>
</template>

<script>

export default {
  name: 'App',
  data() {
    return {
    }
  }
}
</script>

<style type="text/scss">
  body {
    background: black;
    color: white;
  }
  html, body {
    padding: 0;
    margin: 0;
  }

  .clap-animation-enter-active, .clap {
    position: absolute;
    bottom: -67px;
    left: 150px;
    width: 66px;
    height: 66px;
    animation-name: clap-on;
    animation-duration: 2s;
    animation-fill-mode: forwards;
  }
  @keyframes clap-on {
    to {
      bottom: 250px;
      width: 200px;
      height: 200px;
      opacity: 0;
    }
  }



</style>
