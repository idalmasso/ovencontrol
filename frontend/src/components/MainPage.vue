<template>
  <v-container>
    <v-row justify="center" class="fill-height" align-content="center">
      <v-col>
        <ListPageButton
          title="Cottura"
          @click="console.log('NOT IMPLEMENTED')"
        ></ListPageButton> </v-col
      ><v-col>
        <ListPageButton
          title="Test"
          @click="$router.push({ name: 'Tests' })"
        ></ListPageButton> </v-col
      ><v-col> <ListPageButton title="Configurazione"></ListPageButton> </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import ListPageButton from "./ListPageButton.vue";
</script>
