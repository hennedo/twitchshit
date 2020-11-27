<template>
  <div id="app">
    <button @click="sendClap()">Send Clap</button>
    <img src="/clap.png" class="clap" v-for="clapp of claps" :style="clapp" :key="clapp.id">
  </div>
</template>

<script>

  import websockets from './assets/sockets'

export default {
  name: 'App',
  data() {
    return {
      claps: [],
      lastID: 0
    }
  } ,
  components: {

  },
  methods: {
    clap: function() {
      let newClap = {
        id: this.lastID++,
        left: (100 + Math.floor(Math.random() * 500)) + "px"
      }
      this.claps.push(newClap)
      setTimeout(() => {
        for(let i = 0; i<this.claps.length; i++) {
          if(this.claps[i].id === newClap.id) {
            this.claps.splice(i, 1)
          }
        }
      }, 3000)
    },
    sendClap: function() {
      websockets.emit("clap", {})
    }
  },
  mounted() {
    websockets.on("clap", this.clap)
  },
  beforeUnmount() {
    websockets.off(this.clap)
  }
}
</script>

<style>
  body {
    background: black;
  }

  .clap {
    position: absolute;
    bottom: -67px;
    left: 150px;
    animation-name: clap-on;
    animation-duration: 2s;
  }
  @keyframes clap-on {
    to {
      bottom: 250px;
      opacity: 0;
    }
  }



</style>
