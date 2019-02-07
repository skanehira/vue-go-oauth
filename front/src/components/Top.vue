<template>
  <div>
    <h1>{{msg}}</h1>
    <div v-if="is_signed_in">
      <button class="button" type="submit" @click="singout">Sign Out</button>
    </div>
    <div v-else>
      <button class="button" type="submit" @click="singin">Sign In</button>
    </div>
  </div>
</template>

<script>
export default {
  data () {
    return {
      msg: 'Gorilla Gorilla Gorilla Gorilla Gorilla Gorilla Gorilla',
      is_signed_in: false
    }
  },
  methods: {
    singin () {
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
        alert(error)
      })
    },
    singout () {
      this.$axios.post('/users/signout').then((response) => {
        this.is_signed_in = response.data.is_signed_in
      }, (error) => {
        alert(error)
      })
    }
  },
  mounted () {
    let session = this.$cookie.get('test_session')
    if (session != null) {
      this.is_signed_in = true
    } else {
      this.is_signed_in = false
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
