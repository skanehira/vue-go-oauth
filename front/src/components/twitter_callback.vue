<template>
  <div>
    <h1>{{msg}}</h1>
  </div>
</template>

<script>
export default {
  data () {
    return {
      msg: 'Authenticating...'
    }
  },
  methods: {
    auth () {
      let url = window.location.href
      // TODO remove this
      console.log(url)

      this.$axios.get('/auth/twiiter', {url: url}).then((response) => {
        this.$router.push('mypage')
      }, (error) => {
        alert(error)
        this.$router.push('/')
      })
    }
  },
  mounted () {
    this.auth()
  }
}
</script>
