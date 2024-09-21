import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(), // import.meta.env.BASE_URL
  routes: [
    {
      path: "/",
      redirect: "/home",
    },
    {
      path: "/home",
      name: "home",
      component: () => import("@/views/HomeView.vue")
    },
    {
      path: "/main",
      name: "main",
      component: () => import("@/views/MainView.vue")
    },
    {
      path: "/test",
      name: "test",
      component: () => import("@/views/TestView.vue")
    }
  ],
});

export default router
