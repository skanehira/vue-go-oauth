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
      msg: 'Gorilla Gorilla Gorilla Gorilla Gorilla Gorilla Gorilla'
    }
  },
  methods: {
    login () {
      this.$axios.post('/users/signin').then((response) => {
        switch (response.data.status) {
          case 200:
            location.href = response.data.url
            break
          case 302:
            this.$router.push(response.data.url)
            break
          default:
        }
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
