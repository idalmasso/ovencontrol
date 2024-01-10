// Composables
import { createRouter, createWebHistory } from "vue-router";

const routes = [
  {
    path: "/",
    component: () => import("@/layouts/default/DefaultLayout.vue"),
    children: [
      {
        path: "",
        name: "Home",
        // route level code-splitting
        // this generates a separate chunk (Home-[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        component: () => import("@/views/HomeView.vue"),
        meta: { title: "CONTROLLO FORNO" },
      },
    ],
  },
  {
    path: "/list-runs",
    component: () => import("@/layouts/default/DefaultLayout.vue"),
    children: [
      {
        path: "",
        name: "ListProgramsRun",
        // route level code-splitting
        // this generates a separate chunk (Home-[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        component: () => import("@/views/ListProgramRunsView.vue"),
        meta: { title: "Programmi di cottura" },
      },
    ],
  },
  {
    path: "/oven-run/:programName",
    component: () => import("@/layouts/default/DefaultLayout.vue"),
    children: [
      {
        path: "",
        name: "OvenRun",
        // route level code-splitting
        // this generates a separate chunk (Home-[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        component: () => import("@/views/OvenRun.vue"),
        meta: { title: "Cottura" },
        props: true
      },
    ],
  },
  {
    path: "/tests",
    component: () => import("@/layouts/default/DefaultLayout.vue"),
    children: [
      {
        path: "",
        name: "ListTests",
        // route level code-splitting
        // this generates a separate chunk (Home-[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        component: () => import("@/views/ListTestsView.vue"),
        meta: { title: "TESTS" },
      },
    ],
  },
  {
    path: "/configurations",
    component: () => import("@/layouts/default/DefaultLayout.vue"),
    children: [
      {
        path: "",
        name: "ListConfigurations",
        // route level code-splitting
        // this generates a separate chunk (Home-[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        component: () => import("@/views/ListConfigurationsView.vue"),
        meta: { title: "CONFIGURAZIONI" },
      },
    ],
  },
  {
    path: "/configurations/programs",
    component: () => import("@/layouts/default/DefaultLayout.vue"),
    children: [
      {
        path: "",
        name: "ListProgramConfigurations",
        // route level code-splitting
        // this generates a separate chunk (Home-[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        component: () => import("@/views/ListProgramConfigurationsView.vue"),
        meta: { title: "PROGRAMMI FORNO" },
      },
    ],
  },
  {
    path: "/configurations/edit-program/:programName",
    component: () => import("@/layouts/default/DefaultLayout.vue"),
    children: [
      {
        path: "",
        name: "EditOvenProgram",
        // route level code-splitting
        // this generates a separate chunk (Home-[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        component: () => import("@/views/EditProgramView.vue"),
        meta: { title: "MODIFICA PROGRAMMA" },
        props: true
      },
    ],
  },
  {
    path: "/configurations/edit-program",
    component: () => import("@/layouts/default/DefaultLayout.vue"),
    children: [
      {
        path: "",
        name: "NewOvenProgram",
        // route level code-splitting
        // this generates a separate chunk (Home-[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        component: () => import("@/views/EditProgramView.vue"),
        meta: { title: "NUOVO PROGRAMMA" }
      },
    ],
  },
  {
    path: "/configurations/programs",
    component: () => import("@/layouts/default/DefaultLayout.vue"),
    children: [
      {
        path: "",
        name: "ListProgramConfigurations",
        // route level code-splitting
        // this generates a separate chunk (Home-[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        component: () => import("@/views/ListProgramConfigurationsView.vue"),
        meta: { title: "PROGRAMMI FORNO" },
      },
    ],
  },
  {
    path: "/tests/temperature-check",
    component: () => import("@/layouts/default/DefaultLayout.vue"),
    children: [
      {
        path: "",
        name: "TemperatureCheck",
        // route level code-splitting
        // this generates a separate chunk (Home-[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        component: () => import("@/views/TemperatureCheck.vue"),
        meta: { title: "TEST TEMPERATURA" },
      },
    ],
  },
  {
    path: "/tests/test-ramp",
    component: () => import("@/layouts/default/DefaultLayout.vue"),
    children: [
      {
        path: "",
        name: "TestRamp",
        // route level code-splitting
        // this generates a separate chunk (Home-[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        component: () => import("@/views/TestRampView.vue"),
        meta: { title: "TEST RAMPA DI SALITA" },
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

export default router;
