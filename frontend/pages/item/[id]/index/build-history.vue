<script setup lang="ts">
  import { useI18n } from "vue-i18n";
  import { toast } from "@/components/ui/sonner";
  import BaseContainer from "@/components/Base/Container.vue";
  import BaseCard from "@/components/Base/Card.vue";
  import type { BuildOut, BuildPartOut } from "~~/lib/api/types/data-contracts";

  const props = defineProps<{
    item: { id: string; name: string };
  }>();

  const { t } = useI18n();
  const api = useUserApi();

  const history = ref<BuildPartOut[]>([]);
  const buildsById = ref<Record<string, BuildOut>>({});
  const loading = ref(true);

  onMounted(async () => {
    const [historyRes, buildsRes] = await Promise.all([
      api.builds.getEntityHistory(props.item.id),
      api.builds.getAll(),
    ]);

    if (historyRes.error) {
      toast.error(t("builds.toast.history_failed_load"));
    } else {
      history.value = historyRes.data;
    }

    if (!buildsRes.error) {
      const byId: Record<string, BuildOut> = {};
      for (const build of buildsRes.data) {
        byId[build.id] = build;
      }
      buildsById.value = byId;
    }

    loading.value = false;
  });

  const currentlyMounted = computed(() => history.value.filter(p => !p.removedDate));
  const pastMounts = computed(() => history.value.filter(p => p.removedDate));
</script>

<template>
  <BaseContainer class="flex flex-col gap-6">
    <div v-if="loading">{{ $t("global.loading") }}</div>

    <template v-else>
      <p v-if="history.length === 0" class="text-foreground/70">{{ $t("builds.history.empty") }}</p>

      <template v-else>
        <BaseCard v-if="currentlyMounted.length > 0">
          <template #title> {{ $t("builds.history.currently_mounted") }} </template>
          <ul class="flex flex-col gap-2 p-4">
            <li v-for="part in currentlyMounted" :key="part.id" class="flex items-center justify-between gap-4">
              <NuxtLink :to="`/builds/${part.buildId}`" class="font-medium hover:underline">
                {{ buildsById[part.buildId]?.name ?? part.buildId }}
              </NuxtLink>
              <span class="text-sm text-foreground/70">
                {{ $t("builds.parts.mounted") }}: {{ new Date(part.addedDate).toLocaleDateString() }}
              </span>
            </li>
          </ul>
        </BaseCard>

        <BaseCard v-if="pastMounts.length > 0" collapsable>
          <template #title> {{ $t("builds.history.title") }} </template>
          <div class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead>
                <tr class="border-b text-left text-foreground/70">
                  <th class="p-4">{{ $t("builds.history.build") }}</th>
                  <th class="p-4">{{ $t("builds.parts.position") }}</th>
                  <th class="p-4">{{ $t("builds.parts.mounted") }}</th>
                  <th class="p-4">{{ $t("builds.parts.removed") }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="part in pastMounts" :key="part.id" class="border-b">
                  <td class="p-4">
                    <NuxtLink :to="`/builds/${part.buildId}`" class="hover:underline">
                      {{ buildsById[part.buildId]?.name ?? part.buildId }}
                    </NuxtLink>
                  </td>
                  <td class="p-4">{{ part.position }}</td>
                  <td class="p-4">{{ new Date(part.addedDate).toLocaleDateString() }}</td>
                  <td class="p-4">
                    {{ part.removedDate ? new Date(part.removedDate).toLocaleDateString() : "—" }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </BaseCard>
      </template>
    </template>
  </BaseContainer>
</template>
