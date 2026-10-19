<script setup lang="ts">
  import { useI18n } from "vue-i18n";
  import { toast } from "@/components/ui/sonner";
  import { Button } from "@/components/ui/button";
  import { DialogRoot } from "reka-ui";
  import { DialogContent, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
  import { Input } from "@/components/ui/input";
  import { Label } from "@/components/ui/label";
  import { Switch } from "@/components/ui/switch";
  import BaseContainer from "@/components/Base/Container.vue";
  import BaseSectionHeader from "@/components/Base/SectionHeader.vue";
  import BaseCard from "@/components/Base/Card.vue";

  definePageMeta({
    middleware: ["auth"],
  });

  const { t } = useI18n();

  useHead({
    title: "HomeBox | " + t("menu.setups"),
  });

  const api = useUserApi();

  const { data: setups, refresh } = useAsyncData("setups", async () => {
    const { data, error } = await api.setups.getAll();
    if (error) {
      toast.error(t("setups.toast.failed_load"));
      return [];
    }
    return data;
  });

  const setupList = computed(() => setups.value ?? []);

  const showCreate = ref(false);
  const creating = ref(false);

  const createForm = reactive({
    name: "",
    description: "",
    sport: "",
    isActive: true,
  });

  async function createSetup() {
    if (!createForm.name) {
      return;
    }

    creating.value = true;
    const { error, data } = await api.setups.create({ ...createForm });
    creating.value = false;

    if (error) {
      toast.error(t("setups.toast.failed_create"));
      return;
    }

    toast.success(t("setups.toast.created"));
    showCreate.value = false;
    createForm.name = "";
    createForm.description = "";
    createForm.sport = "";
    createForm.isActive = true;
    await refresh();
    await navigateTo(`/setup/${data.id}`);
  }
</script>

<template>
  <div>
    <BaseContainer class="flex flex-col gap-4">
      <div class="flex items-center justify-between gap-4">
        <BaseSectionHeader> {{ $t("menu.setups") }} </BaseSectionHeader>
        <Button @click="showCreate = true">
          {{ $t("setups.create") }}
        </Button>
      </div>

      <p v-if="setupList.length === 0" class="text-foreground/70">
        {{ $t("setups.empty") }}
      </p>

      <BaseCard v-for="setupItem in setupList" :key="setupItem.id">
        <NuxtLink :to="`/setup/${setupItem.id}`" class="flex items-center justify-between gap-4 p-4">
          <div class="flex flex-col gap-1">
            <span class="font-medium">{{ setupItem.name }}</span>
            <span v-if="setupItem.description" class="text-sm text-foreground/70">{{ setupItem.description }}</span>
            <span v-if="setupItem.sport" class="text-sm text-foreground/70">{{ setupItem.sport }}</span>
          </div>
          <span
            class="rounded-full px-2 py-1 text-xs"
            :class="setupItem.isActive ? 'bg-green-500/15 text-green-600' : 'bg-muted text-muted-foreground'"
          >
            {{ setupItem.isActive ? $t("setups.active") : $t("setups.inactive") }}
          </span>
        </NuxtLink>
      </BaseCard>

      <DialogRoot v-model:open="showCreate">
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{{ $t("setups.create") }}</DialogTitle>
          </DialogHeader>
          <form class="flex flex-col gap-4" @submit.prevent="createSetup">
            <div class="flex flex-col gap-2">
              <Label for="setup-name">{{ $t("global.name") }}</Label>
              <Input id="setup-name" v-model="createForm.name" required />
            </div>
            <div class="flex flex-col gap-2">
              <Label for="setup-description">{{ $t("global.description") }}</Label>
              <Input id="setup-description" v-model="createForm.description" />
            </div>
            <div class="flex flex-col gap-2">
              <Label for="setup-sport">{{ $t("setups.sport") }}</Label>
              <Input id="setup-sport" v-model="createForm.sport" />
            </div>
            <Label class="flex cursor-pointer items-center gap-2">
              <Switch v-model="createForm.isActive" />
              {{ $t("setups.active") }}
            </Label>
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
