<template>
  <ListPageWithButtons
    :listItems="listItems"
    @buttonClicked="buttonClickedHandler"
  ></ListPageWithButtons>
</template>

<script setup>
import { useRouter } from "vue-router";
import { ref } from "vue";
import { useAppStore } from "@/store/app";
import ListPageWithButtons from "./parts/ListPageWithButtons.vue";
const router = useRouter();
const store = useAppStore();
const listItems = ref([
  {
    title: "Configura programmi",
    icon: "mdi-tune-vertical-variant",
    name: "ListProgramConfigurations",
  },
  {
    title: "Esporta programmi terminati",
    icon: "mdi-content-save-move",
    name: "Export",
  },
  {
    title: "Configurazione parametri",
    icon: "mdi-wrench-cog",
    name: "OvenConfig",
  },
]);
function buttonClickedHandler(name) {
  console.log("NAME:" + name);
  if (name == "Export") {
    fetch("http://localhost:3333/api/configuration/move-runs-usb", {
      method: "POST",
    }).then((response) => {
      if (response.ok) {
        alert("Success!");
      } else {
        response.json().then((data) => store.setAPIError(data["Error"]));
      }
    });
  } else {
    router.push({ name: name });
  }
}
</script>
