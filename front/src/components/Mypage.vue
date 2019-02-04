<template>
  <div>
    <h1>Mypage</h1>
    <table border=1>
      <tr>
        <th>Name</th>
        <td>{{my.name}}</td>
      </tr>
      <tr>
        <th>Description</th>
        <td>{{my.description}}</td>
      </tr>
      <tr>
        <th>Location</th>
        <td>{{my.location}}</td>
      </tr>
      <tr>
        <th>ScreenName</th>
        <td>{{my.screen_name}}</td>
      </tr>
    </table>
  </div>
</template>

<script>
export default {
  data () {
    return {
      my: null
    }
  },
  methods: {
    getUserInfo () {
      this.$axios.get('/users').then((response) => {
        console.log(response)
        this.my = response.data
      }, (error) => {
        console.log(error)
      })
    }
  },
  mounted () {
    this.getUserInfo()
  }
}
</script>
