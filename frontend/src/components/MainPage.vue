<template>
  <ListPageWithButtons
    :listItems="listItems"
    @buttonClicked="buttonClickedHandler"
  ></ListPageWithButtons>
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import ListPageWithButtons from "./parts/ListPageWithButtons.vue";
const router = useRouter();
const listItems = ref([
  {
    title: "Cottura",
    icon: "mdi-fire-circle",
    name: "ListProgramsRun",
  },
  {
    title: "Tests",
    icon: "mdi-chart-timeline-variant-shimmer",
    name: "ListTests",
  },
  {
    title: "Configurazione",
    icon: "mdi-cog",
    name: "ListConfigurations",
  },
]);
function buttonClickedHandler(name) {
  router.push({ name: name });
}
</script>
