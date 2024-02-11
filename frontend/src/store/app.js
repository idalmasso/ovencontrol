// Utilities
import { defineStore } from "pinia";

export const useAppStore = defineStore("app", {
  state: () => ({
    apiError: "",
  }),
  getters: {
    dialog: (state) => state.apiError != "",
  },
  actions: {
    setAPIError(errorStr) {
      this.apiError = errorStr;
    },
    resetAPIError() {
      this.apiError = "";
    },
  },
});
