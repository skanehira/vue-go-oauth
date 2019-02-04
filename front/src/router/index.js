import Vue from 'vue'
import Router from 'vue-router'
import Top from '@/components/Top'
import NotFound from '@/components/NotFound'
import Mypage from '@/components/Mypage'

Vue.use(Router)

export default new Router({
  routes: [
    { path: '/', name: 'Top', component: Top },
    { path: '/mypage', name: 'Mypage', component: Mypage },
    { path: '*', name: 'NotFound', component: NotFound }
  ]
})
