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
      this.$axios.get('/auth/twiiter', {url: window.location.href}).then((response) => {
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
