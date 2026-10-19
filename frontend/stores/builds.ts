import { defineStore } from "pinia";
import type { BuildCreate, BuildOut, BuildUpdate } from "~~/lib/api/types/data-contracts";

export const useBuildStore = defineStore("builds", {
  state: () => ({
    allBuilds: null as BuildOut[] | null,
    client: useUserApi(),
  }),
  getters: {
    /**
     * builds represents the builds currently in the store, synchronized with
     * the server by the actions below.
     */
    builds(state): BuildOut[] {
      return state.allBuilds ?? [];
    },
  },
  actions: {
    async refresh() {
      const result = await this.client.builds.getAll();
      if (result.error) {
        return result;
      }

      this.allBuilds = result.data;
      return result;
    },
    async create(data: BuildCreate) {
      const result = await this.client.builds.create(data);
      if (!result.error) {
        await this.refresh();
      }
      return result;
    },
    async update(id: string, data: BuildUpdate) {
      const result = await this.client.builds.update(id, data);
      if (!result.error) {
        await this.refresh();
      }
      return result;
    },
    async delete(id: string) {
      const result = await this.client.builds.delete(id);
      if (!result.error) {
        await this.refresh();
      }
      return result;
    },
  },
});
