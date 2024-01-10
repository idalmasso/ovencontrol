<template>
  <ListPageWithButtons
    :listItems="listItems"
    @buttonClicked="buttonClickedHandler"
  ></ListPageWithButtons>
</template>

<script setup>
import { useRouter } from "vue-router";
import { ref } from "vue";
import ListPageWithButtons from "./parts/ListPageWithButtons.vue";
const router = useRouter();
const listItems = ref([
  {
    title: "Configura programmi",
    icon: "mdi-tune-vertical-variant",
    name: "ListProgramConfigurations",
  },
  {
    title: "Configurazione parametri",
    icon: "mdi-wrench-cog",
    name: "OvenConfig",
  },
]);
function buttonClickedHandler(name) {
  router.push({ name: name });
}
</script>
