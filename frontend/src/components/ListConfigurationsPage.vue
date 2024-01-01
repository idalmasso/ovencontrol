<template>
  <ListPageWithButtons :listItems="listItems"></ListPageWithButtons>
</template>

<script setup>
import { useRouter } from "vue-router";
import { ref } from "vue";
import ListPageWithButtons from "./parts/ListPageWithButtons.vue";
const router = useRouter();
const listItems = ref([
  {
    button: { title: "Configura programmi", icon: "mdi-tune-vertical-variant" },
    action: () => router.push({ name: "ListProgramConfigurations" }),
  },
  {
    button: { title: "Configurazione parametri", icon: "mdi-wrench-cog" },
    action: () => console.log("NOT IMPLEMENTED"),
  },
]);
</script>
