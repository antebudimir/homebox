import { defineStore } from "pinia";
import type { SetupCreate, SetupOut, SetupUpdate } from "~~/lib/api/types/data-contracts";

export const useSetupStore = defineStore("setups", {
  state: () => ({
    allSetups: null as SetupOut[] | null,
    client: useUserApi(),
  }),
  getters: {
    /**
     * setups represents the setups currently in the store, synchronized with
     * the server by the actions below.
     */
    setups(state): SetupOut[] {
      return state.allSetups ?? [];
    },
  },
  actions: {
    async refresh() {
      const result = await this.client.setups.getAll();
      if (result.error) {
        return result;
      }

      this.allSetups = result.data;
      return result;
    },
    async create(data: SetupCreate) {
      const result = await this.client.setups.create(data);
      if (!result.error) {
        await this.refresh();
      }
      return result;
    },
    async update(id: string, data: SetupUpdate) {
      const result = await this.client.setups.update(id, data);
      if (!result.error) {
        await this.refresh();
      }
      return result;
    },
    async delete(id: string) {
      const result = await this.client.setups.delete(id);
      if (!result.error) {
        await this.refresh();
      }
      return result;
    },
  },
});
