import { createRouter, createWebHistory } from 'vue-router'

import HomeView from '../views/HomeView.vue'
import PostView from '../views/PostView.vue'
import PlaygroundView from '../views/PlaygroundView.vue'
import ProfileView from '../views/ProfileView.vue'

const routes = [
    { path: '/', name: 'home', component: HomeView },
    { path: '/posts/:postId', name: 'post', component: PostView, props: true },
    { path: '/playground', name: 'playground', component: PlaygroundView },
    { path: '/profile', name: 'profile', component: ProfileView },
]

export const router = createRouter({
    history: createWebHistory(),
    routes,
})
