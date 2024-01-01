<template>
  <ListPageWithButtons :listItems="listItems"></ListPageWithButtons>
</template>

<script setup>
//import { useRouter } from "vue-router";
import { ref } from "vue";
import ListPageWithButtons from "./parts/ListPageWithButtons.vue";
//const router = useRouter();
const getPrograms = () => {
  fetch("http://localhost:3333/api/configuration/programs").then((response) => {
    if (response.ok) {
      response.json().then((data) => {
        console.log(Object.keys(data).map((a) => a));

        listItems.value = Object.keys(data).map((name) => {
          return {
            button: {
              title: name,
              icon: "mdi-wrench-cog",
            },
            action: () => console.log("NOT IMPLEMENTED"),
          };
        });
        programList.value = Object.values(data);
      });
    }
  });
};
const programList = ref([]);
const listItems = ref([]);
getPrograms();
</script>
