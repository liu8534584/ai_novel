import { createRouter, createWebHistory } from 'vue-router'
import MainLayout from '../layouts/MainLayout.vue'
import Home from '../views/Home.vue'
import Settings from '../views/Settings.vue'
import PromptTemplates from '../views/PromptTemplates.vue'
import Plans from '../views/Plans.vue'
import Writing from '../views/Writing.vue'
import MasterOutline from '../views/MasterOutline.vue'

const router = createRouter({
  history: createWebHistory('/ui/'),
  routes: [
    {
      path: '/',
      component: MainLayout,
      children: [
        { path: '', component: Home },
        { path: 'plans/:bookId?', component: Plans },
        { path: 'outline/:bookId?', component: MasterOutline },
        { path: 'writing/:bookId', component: Writing },
        { path: 'settings', component: Settings },
        { path: 'prompts', component: PromptTemplates }
      ]
    }
  ]
})

export default router
