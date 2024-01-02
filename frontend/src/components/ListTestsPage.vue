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
    button: { title: "Test temperatura", icon: "mdi-thermometer-check" },
    action: () => router.push({ name: "TemperatureCheck" }),
  },
  {
    button: { title: "Test rampa", icon: "mdi-finance" },
    action: () => router.push({ name: "TestRamp" }),
  },
]);
</script>
