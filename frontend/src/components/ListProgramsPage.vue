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
        console.log(Object.keys(data).map((a) => a));

        listItems.value = Object.keys(data).map((name) => {
          return {
            title: name,
            icon: "mdi-wrench-cog",
            name: name,
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
function buttonClickedHandler(name) {
  console.log("Emitted" + name);
  emit("programSelected", name);
}
</script>
