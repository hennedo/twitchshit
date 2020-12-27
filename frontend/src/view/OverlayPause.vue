<template>
  <div>
    <transition-group enter-active-class="clap-animation-enter-active">
      <img src="/clap.png" v-for="clapp of claps" :style="clapp" :key="clapp.id" @animationend="endClap(clapp.id)" class="clap">
    </transition-group>
    <div class="last-follow" v-if="follows.last">
      <span class="username">{{ follows.last }}</span>
      <span class="description">Zuletzt gefolgt</span>
      <img src="/followline.png">
    </div>
    <div class="follow-count">
      Follow Goal:
      <div class="follow-count-bg">
        <div class="follow-count-progress" :style="{width: `${follows.count / 50 * 100}%`}">
          <span>{{ follows.count }} / 50 </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>

  import websockets from '../assets/sockets'

  export default {
    name: 'OverlayPause',
    data() {
      return {
        claps: [],
        lastID: 0,
        follows: {
          last: "",
          count: 0
        }
      }
    } ,
    components: {

    },
    methods: {
      clap: function() {
        let newClap = {
          id: this.lastID++,
          left: (100 + Math.floor(Math.random() * 500)) + "px",
        }
        this.claps.push(newClap)
      },
      endClap: function(id) {
        for(let i = 0; i < this.claps.length; i++) {
          if(this.claps[i].id === id) {
            this.claps.splice(i, 1)
          }
        }
      },
      follow: function(msg) {
        this.follows.last = msg.from_name;
      },
      followers: function(msg) {
        this.follows.count = msg.count;
      }
    },
    mounted() {
      websockets.on("clap", this.clap)
      websockets.on("follow", this.follow)
      websockets.on("followers", this.followers)
      websockets.emit("connect")
    },
    beforeUnmount() {
      websockets.off(this.clap)
      websockets.off(this.follow)
      websockets.off(this.follows)
    }
  }
</script>

<style type="text/scss" scoped>


  .last-follow {
    font-family: "Londrina Solid";
    color: white;
    position: relative;
    left: 1577px;
    top: 995px;
  }
  .last-follow .username {
    font-size: 26px;
    left: 5px;
    top: -26px;
    text-align: center;
    width: 280px;
    overflow: hidden;
  }
  .last-follow .description {
    left: 95px;
    top: 5px;
  }
  .last-follow span {
    display: inline-block;
    position: absolute;
    left: 0;
    top: 0;
  }
  .last-follow img {
    position: absolute;
    left: 0;
    top: 0;
  }

  .follow-count {
    position: absolute;
    left: 50px;
    bottom: 50px;
    font-family: "Londrina Solid";
    font-size: 25px;
    color: white;
  }
  .follow-count-bg {
    position: relative;
    width: 400px;
    height: 30px;
    background: white;
    padding: 2px;
    border-radius: 3px;
  }
  .follow-count-progress span {
    text-align: center;
  }
  .follow-count-progress {
    height: 30px;
    border-radius: 3px;
    position: absolute;
    left: 2px;
    top: 2px;
    background: green;
    text-align: center;
  }
</style>
