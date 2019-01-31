<template>
  <div>
    <button class="button" type="submit" @click="login">Login Twitter</button>
  </div>
</template>

<script>
export default {
  data () {
    return {
      msg: 'Welcome to Your Vue.js App'
    }
  },
  methods: {
    login () {
      // TODO remove this
      console.log('start login twitter')

      this.$axios.get('/login/twitter').then((response) => {
        location.href = response.auth_url
      }, (error) => {
        alert(error)
      })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1, h2 {
  font-weight: normal;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
