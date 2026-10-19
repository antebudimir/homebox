<script setup lang="ts">
  import { useI18n } from "vue-i18n";
  import { toast } from "@/components/ui/sonner";
  import { Button } from "@/components/ui/button";
  import { DialogRoot } from "reka-ui";
  import { DialogContent, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
  import { Input } from "@/components/ui/input";
  import { Label } from "@/components/ui/label";
    import BaseContainer from "@/components/Base/Container.vue";
  import BaseSectionHeader from "@/components/Base/SectionHeader.vue";
  import BaseCard from "@/components/Base/Card.vue";
  import type { BuildStatus } from "~~/lib/api/types/data-contracts";

  definePageMeta({
    middleware: ["auth"],
  });

  const { t } = useI18n();

  useHead({
    title: "HomeBox | " + t("menu.builds"),
  });

  const api = useUserApi();

  const { data: builds, refresh } = useAsyncData("builds", async () => {
    const { data, error } = await api.builds.getAll();
    if (error) {
      toast.error(t("builds.toast.failed_load"));
      return [];
    }
    return data;
  });

  const buildList = computed(() => builds.value ?? []);

  const showCreate = ref(false);
  const creating = ref(false);

  const createForm = reactive({
    name: "",
    description: "",
    sport: "",
    status: "assembled" as BuildStatus,
  });

  async function createBuild() {
    if (!createForm.name) {
      return;
    }

    creating.value = true;
    const { error, data } = await api.builds.create({ ...createForm });
    creating.value = false;

    if (error) {
      toast.error(t("builds.toast.failed_create"));
      return;
    }

    toast.success(t("builds.toast.created"));
    showCreate.value = false;
    createForm.name = "";
    createForm.description = "";
    createForm.sport = "";
    createForm.status = "assembled";
    await refresh();
    await navigateTo(`/builds/${data.id}`);
  }
</script>

<template>
  <div>
    <BaseContainer class="flex flex-col gap-4">
      <div class="flex items-center justify-between gap-4">
        <BaseSectionHeader> {{ $t("menu.builds") }} </BaseSectionHeader>
        <Button @click="showCreate = true">
          {{ $t("builds.create") }}
        </Button>
      </div>

      <p v-if="buildList.length === 0" class="text-foreground/70">
        {{ $t("builds.empty") }}
      </p>

      <BaseCard v-for="buildItem in buildList" :key="buildItem.id">
        <NuxtLink :to="`/builds/${buildItem.id}`" class="flex items-center justify-between gap-4 p-4">
          <div class="flex flex-col gap-1">
            <span class="font-medium">{{ buildItem.name }}</span>
            <span v-if="buildItem.description" class="text-sm text-foreground/70">{{ buildItem.description }}</span>
            <span v-if="buildItem.sport" class="text-sm text-foreground/70">{{ buildItem.sport }}</span>
          </div>
          <span
            class="rounded-full px-2 py-1 text-xs"
            :class="buildItem.status === 'assembled' ? 'bg-green-500/15 text-green-600' : 'bg-muted text-muted-foreground'"
          >
            {{ $t(`builds.status.${buildItem.status}`) }}
          </span>
        </NuxtLink>
      </BaseCard>

      <DialogRoot v-model:open="showCreate">
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{{ $t("builds.create") }}</DialogTitle>
          </DialogHeader>
          <form class="flex flex-col gap-4" @submit.prevent="createBuild">
            <div class="flex flex-col gap-2">
              <Label for="build-name">{{ $t("global.name") }}</Label>
              <Input id="build-name" v-model="createForm.name" required />
            </div>
            <div class="flex flex-col gap-2">
              <Label for="build-description">{{ $t("global.description") }}</Label>
              <Input id="build-description" v-model="createForm.description" />
            </div>
            <div class="flex flex-col gap-2">
              <Label for="build-sport">{{ $t("builds.sport") }}</Label>
              <Input id="build-sport" v-model="createForm.sport" />
            </div>
            <div class="flex flex-col gap-2">
              <Label for="build-status">{{ $t("builds.status.label") }}</Label>
              <select id="build-status" v-model="createForm.status" class="bg-background border rounded-md px-3 py-2 text-sm">
                <option value="assembled">{{ $t("builds.status.assembled") }}</option>
                <option value="disassembled">{{ $t("builds.status.disassembled") }}</option>
                <option value="retired">{{ $t("builds.status.retired") }}</option>
              </select>
            </div>
            <DialogFooter>
              <Button type="button" variant="outline" @click="showCreate = false">
                {{ $t("global.cancel") }}
              </Button>
              <Button type="submit" :disabled="creating">{{ $t("global.create") }}</Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </DialogRoot>
    </BaseContainer>
  </div>
</template>
