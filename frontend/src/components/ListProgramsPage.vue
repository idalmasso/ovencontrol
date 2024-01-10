<template>
  <ListPageWithButtons
    :listItems="listItems"
    @buttonClicked="buttonClickedHandler"
  ></ListPageWithButtons>
</template>

<script setup>
import { ref, defineEmits } from "vue";
import ListPageWithButtons from "./parts/ListPageWithButtons.vue";
const emit = defineEmits({ programSelected: String });

const getPrograms = () => {
  fetch("http://localhost:3333/api/configuration/programs").then((response) => {
    if (response.ok) {
      response.json().then((data) => {
        listItems.value = Object.keys(data).map((name) => {
          return {
            title: name,
            icon: "mdi-wrench-cog",
            name: name,
            color: data[name]['icon-color']
          };
        });
      });
    }
  });
};
const listItems = ref([]);
getPrograms();
function buttonClickedHandler(name) {
  emit("programSelected", name);
}
</script>
