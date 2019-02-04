<template>
  <div>
    <h1>{{msg}}</h1>
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
      this.$axios.post('/users/signin').then((response) => {
        location.href = response.data.URL
      }, (error) => {
        console.log(error)
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
