<template>
  <ListPageWithButtons
    :listItems="listItems"
    @buttonClicked="buttonClickedHandler"
  ></ListPageWithButtons>
</template>

<script setup>
import { ref, defineEmits } from "vue";
import ListPageWithButtons from "./parts/ListPageWithButtons.vue";
import { onMounted } from "vue";
import { useAppStore } from "@/store/app";
const emit = defineEmits({ programSelected: String });
const store = useAppStore();
const getPrograms = () => {
  fetch("http://localhost:3333/api/configuration/programs").then((response) => {
    if (response.ok) {
      response.json().then((data) => {
        listItems.value = Object.keys(data).map((name) => {
          return {
            title: name,
            icon: "mdi-wrench-cog",
            name: name,
            color: data[name]["icon-color"],
          };
        });
      });
    } else {
      response.json().then((data) => store.setAPIError(data["Error"]));
    }
  });
};
const listItems = ref([]);
onMounted(() => {
  getPrograms();
});

function buttonClickedHandler(name) {
  emit("programSelected", name);
}
</script>
