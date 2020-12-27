<template>
  <div>
    Twitch Oauth nvm this screen
  </div>
</template>

<script>
  import websockets from '../assets/sockets'

  export default {
    name: 'TwitchOauth',
    data() {
      return {
      }
    },
    mounted() {
      let hash = this.$route.hash.substr(1)
      let parts = hash.split("&")
      let params = {}
      parts.forEach(part => {
        let buf = part.split("=")
        params[buf[0]] = buf[1]
      })
      websockets.emit("twitchOauth", params.access_token);
    }
  }
</script>

<style type="text/scss">
</style>
