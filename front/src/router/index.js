import Vue from 'vue'
import Router from 'vue-router'
import Top from '@/components/Top'
import TwitterCallback from '@/components/TwitterCallback'
import NotFound from '@/components/NotFound'

Vue.use(Router)

export default new Router({
  routes: [
    { path: '/', name: 'Top', component: Top },
    { path: '/twitter/callback', name: 'twitter_callback', component: TwitterCallback },
    { path: '*', name: 'NotFound', component: NotFound }
  ]
})
