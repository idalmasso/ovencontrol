<template>
  <v-main>
    <div class="text-center">
      <v-dialog v-model="dialog" width="auto">
        <v-card title="Error">
          <v-card-text>
            {{ apiError }}
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary" block @click="store.resetAPIError"
              >Close Dialog</v-btn
            >
          </v-card-actions>
        </v-card>
      </v-dialog>
      <h1>
        {{ $route.meta.title
        }}<span v-if="$route.params.programName">
          - {{ $route.params.programName }}</span
        >
      </h1>
    </div>
    <router-view />
  </v-main>
</template>

<script setup>
import { useAppStore } from "@/store/app.js";
import { storeToRefs } from "pinia";
const store = useAppStore();

const { dialog, apiError } = storeToRefs(store);
</script>
