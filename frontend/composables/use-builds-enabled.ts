import { computed } from "vue";

/**
 * Shared kill-switch state for the builds feature, backed by the group's
 * show_builds setting. Keyed useAsyncData so every consumer shares one
 * fetch (and SSR payload); call `refreshNuxtData("builds-enabled")` after
 * changing the setting to update all consumers.
 */
export const useBuildsEnabled = () => {
  const api = useUserApi();
  const { selectedCollection } = useCollections();

  const { data } = useAsyncData(
    "builds-enabled",
    async () => {
      const cid = selectedCollection.value?.id;
      if (!cid) {
        return true;
      }
      const { data: group, error } = await api.group.get(cid);
      if (error || !group) {
        return true;
      }
      return group.showBuilds;
    },
    { watch: [selectedCollection], default: () => true }
  );

  const buildsEnabled = computed(() => data.value !== false);

  return { buildsEnabled };
};
