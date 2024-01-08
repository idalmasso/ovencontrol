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
    title: "Test temperatura",
    icon: "mdi-thermometer-check",
    name: "TemperatureCheck",
  },
  {
    title: "Test rampa",
    icon: "mdi-finance",
    name: "TestRamp",
  },
]);
function buttonClickedHandler(name) {
  router.push({ name: name });
}
</script>
