/**
 * Kills the builds feature at the URL level: when the group has builds
 * disabled, any /builds route lands on home. Data is never deleted — this
 * is display-only.
 */
export default defineNuxtRouteMiddleware(async () => {
  const { selectedCollection, load } = useCollections();
  if (!selectedCollection.value?.id) {
    await load();
  }
  const cid = selectedCollection.value?.id;
  if (!cid) {
    return;
  }
  const api = useUserApi();
  const { data } = await api.group.get(cid);
  if (data && data.showBuilds === false) {
    return navigateTo("/");
  }
});
