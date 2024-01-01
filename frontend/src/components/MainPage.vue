<template>
  <ListPageWithButtons :listItems="listItems"></ListPageWithButtons>
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import ListPageWithButtons from "./parts/ListPageWithButtons.vue";
const router = useRouter();
const listItems = ref([
  {
    button: { title: "Cottura", icon: "mdi-fire-circle" },
    action: () => console.log("NOT IMPLEMENTED"),
  },
  {
    button: { title: "Tests", icon: "mdi-chart-timeline-variant-shimmer" },
    action: () => router.push({ name: "ListTests" }),
  },
  {
    button: { title: "Configurazione", icon: "mdi-cog" },
    action: () => router.push({ name: "ListConfigurations" }),
  },
]);
</script>
