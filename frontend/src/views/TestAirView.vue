<template>
  <v-container>
    <v-row
      ><v-btn color="blue" @click="$router.push({ name: 'Home' })"
        >Indietro</v-btn
      ></v-row
    ></v-container
  >
  <ListPageWithButtons
    :listItems="listItems"
    @buttonClicked="buttonClickedHandler"
  ></ListPageWithButtons>
</template>

<script setup>
import { ref } from "vue";
import { useAppStore } from "@/store/app";
import ListPageWithButtons from "../components/parts/ListPageWithButtons.vue";
const store = useAppStore();
const listItems = ref([
  {
    title: "Apri aria",
    icon: "mdi-valve-open",
    name: "OpenAir",
  },
  {
    title: "Chiudi aria",
    icon: "mdi-valve-closed",
    name: "CloseAir",
  },
]);
function buttonClickedHandler(name) {
  switch (name) {
    case "OpenAir":
      fetch("http://localhost:3333/api/processes/open-air", {
        method: "POST",
      }).then((a) => {
        if (!a.ok) {
          a.json().then((data) => {
            store.setAPIError(data["Error"]);
          });
        }
      });
      break;
    case "CloseAir":
      fetch("http://localhost:3333/api/processes/close-air", {
        method: "POST",
      }).then((a) => {
        if (!a.ok) {
          a.json().then((data) => {
            store.setAPIError(data["Error"]);
          });
        }
      });
      break;
  }
}
</script>
