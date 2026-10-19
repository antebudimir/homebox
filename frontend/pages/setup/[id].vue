<script setup lang="ts">
  import { useI18n } from "vue-i18n";
  import { toast } from "@/components/ui/sonner";
  import { Button, ButtonGroup } from "@/components/ui/button";
  import { DialogRoot } from "reka-ui";
  import { DialogContent, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
  import { Input } from "@/components/ui/input";
  import { Label } from "@/components/ui/label";
  import { Switch } from "@/components/ui/switch";
  import BaseContainer from "@/components/Base/Container.vue";
  import BaseSectionHeader from "@/components/Base/SectionHeader.vue";
  import BaseCard from "@/components/Base/Card.vue";
  import type { EntitySummary, SetupOut, SetupPartOut, UsageSessionOut } from "~~/lib/api/types/data-contracts";

  definePageMeta({
    middleware: ["auth"],
  });

  const { t } = useI18n();
  const route = useRoute();
  const api = useUserApi();

  const setupId = computed(() => route.params.id as string);

  const setup = ref<SetupOut | null>(null);
  const parts = ref<SetupPartOut[]>([]);
  const sessions = ref<UsageSessionOut[]>([]);
  const loading = ref(true);

  async function loadAll() {
    loading.value = true;

    const [setupRes, partsRes, sessionsRes] = await Promise.all([
      api.setups.get(setupId.value),
      api.setups.getParts(setupId.value),
      api.setups.getSessions(setupId.value),
    ]);

    if (setupRes.error) {
      toast.error(t("setups.toast.failed_load"));
      loading.value = false;
      return;
    }
    setup.value = setupRes.data;
    parts.value = partsRes.error ? [] : partsRes.data;
    sessions.value = sessionsRes.error ? [] : sessionsRes.data;
    loading.value = false;
  }

  onMounted(loadAll);

  const currentParts = computed(() => parts.value.filter(p => !p.removedDate));

  // --- edit / delete setup ---
  const showEdit = ref(false);
  const saving = ref(false);
  const editForm = reactive({ name: "", description: "", sport: "", isActive: true });

  function openEdit() {
    if (!setup.value) {
      return;
    }
    editForm.name = setup.value.name;
    editForm.description = setup.value.description;
    editForm.sport = setup.value.sport;
    editForm.isActive = setup.value.isActive;
    showEdit.value = true;
  }

  async function saveEdit() {
    if (!setup.value || !editForm.name) {
      return;
    }
    saving.value = true;
    const { error } = await api.setups.update(setup.value.id, { ...editForm });
    saving.value = false;
    if (error) {
      toast.error(t("setups.toast.failed_update"));
      return;
    }
    toast.success(t("setups.toast.updated"));
    showEdit.value = false;
    await loadAll();
  }

  const showDelete = ref(false);
  const deleting = ref(false);

  async function deleteSetup() {
    if (!setup.value) {
      return;
    }
    deleting.value = true;
    const { error } = await api.setups.delete(setup.value.id);
    deleting.value = false;
    if (error) {
      toast.error(t("setups.toast.failed_delete"));
      return;
    }
    toast.success(t("setups.toast.deleted"));
    await navigateTo("/setups");
  }

  // --- parts ---
  const showMount = ref(false);
  const mounting = ref(false);
  const allEntities = ref<EntitySummary[]>([]);
  const mountForm = reactive({ entityId: "", position: "" });

  async function openMount() {
    showMount.value = true;
    if (allEntities.value.length === 0) {
      const { data, error } = await api.items.getAll();
      if (error) {
        toast.error(t("setups.toast.failed_load"));
        return;
      }
      allEntities.value = data.items;
    }
  }

  async function mountPart() {
    if (!mountForm.entityId) {
      return;
    }
    mounting.value = true;
    const { data, error } = await api.setups.addPart(setupId.value, {
      entityId: mountForm.entityId,
      position: mountForm.position,
    });
    mounting.value = false;
    if (error) {
      toast.error(t("setups.toast.part_failed_add"));
      return;
    }
    if (data && data.warnings && data.warnings.length > 0) {
      toast.warning(data.warnings.map(code => t(`setups.warning.${code}`)).join(" "));
    }
    toast.success(t("setups.toast.part_added"));
    showMount.value = false;
    mountForm.entityId = "";
    mountForm.position = "";
    await loadAll();
  }

  async function unmountPart(part: SetupPartOut) {
    const { error } = await api.setups.removePart(setupId.value, part.id);
    if (error) {
      toast.error(t("setups.toast.part_failed_remove"));
      return;
    }
    toast.success(t("setups.toast.part_removed"));
    await loadAll();
  }

  async function purgePart(part: SetupPartOut) {
    const { error } = await api.setups.purgePart(setupId.value, part.id);
    if (error) {
      toast.error(t("setups.toast.part_failed_delete"));
      return;
    }
    toast.success(t("setups.toast.part_deleted"));
    await loadAll();
  }

  const showHistory = ref(false);

  // --- sessions ---
  const showSession = ref(false);
  const loggingSession = ref(false);
  const sessionForm = reactive({
    date: new Date().toISOString().slice(0, 10),
    durationMinutes: 0,
    location: "",
    conditions: "",
    notes: "",
    rating: 0,
  });

  async function logSession() {
    loggingSession.value = true;
    const { error } = await api.setups.createSession(setupId.value, {
      ...sessionForm,
      date: new Date(sessionForm.date).toISOString(),
    });
    loggingSession.value = false;
    if (error) {
      toast.error(t("setups.toast.session_failed_log"));
      return;
    }
    toast.success(t("setups.toast.session_logged"));
    showSession.value = false;
    sessionForm.durationMinutes = 0;
    sessionForm.location = "";
    sessionForm.conditions = "";
    sessionForm.notes = "";
    sessionForm.rating = 0;
    await loadAll();
  }

  async function deleteSession(session: UsageSessionOut) {
    const { error } = await api.setups.deleteSession(setupId.value, session.id);
    if (error) {
      toast.error(t("setups.toast.session_failed_delete"));
      return;
    }
    toast.success(t("setups.toast.session_deleted"));
    await loadAll();
  }
</script>

<template>
  <div>
    <BaseContainer class="flex flex-col gap-6">
      <div v-if="loading">{{ $t("global.loading") }}</div>

      <template v-else-if="setup">
        <div class="flex items-start justify-between gap-4">
          <div class="flex flex-col gap-2">
            <NuxtLink to="/setups" class="text-sm text-foreground/70 hover:underline">
              ← {{ $t("menu.setups") }}
            </NuxtLink>
            <BaseSectionHeader class="!mb-0">
              {{ setup.name }}
              <span
                class="ml-3 rounded-full px-2 py-1 text-xs"
                :class="setup.isActive ? 'bg-green-500/15 text-green-600' : 'bg-muted text-muted-foreground'"
              >
                {{ setup.isActive ? $t("setups.active") : $t("setups.inactive") }}
              </span>
            </BaseSectionHeader>
            <p v-if="setup.description" class="text-foreground/70">{{ setup.description }}</p>
          </div>
          <ButtonGroup>
            <Button variant="outline" @click="openEdit">{{ $t("global.edit") }}</Button>
            <Button variant="destructive" @click="showDelete = true">{{ $t("global.delete") }}</Button>
          </ButtonGroup>
        </div>

        <!-- Parts -->
        <BaseCard collapsable>
          <template #title> {{ $t("setups.parts.title") }} ({{ currentParts.length }}) </template>
          <template #title-actions>
            <div class="mt-2 flex items-center gap-4">
              <Label class="flex cursor-pointer items-center gap-2">
                <Switch v-model="showHistory" />
                {{ $t("setups.parts.history") }}
              </Label>
              <Button size="sm" @click="openMount">{{ $t("setups.parts.mount") }}</Button>
            </div>
          </template>

          <p v-if="currentParts.length === 0" class="px-6 pb-4 text-foreground/70">
            {{ $t("setups.parts.empty") }}
          </p>

          <div v-else class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead>
                <tr class="border-b text-left text-foreground/70">
                  <th class="p-4">{{ $t("setups.parts.entity") }}</th>
                  <th class="p-4">{{ $t("setups.parts.position") }}</th>
                  <th class="p-4">{{ $t("setups.parts.mounted") }}</th>
                  <th v-if="showHistory" class="p-4">{{ $t("setups.parts.removed") }}</th>
                  <th class="p-4" />
                </tr>
              </thead>
              <tbody>
                <tr v-for="part in showHistory ? parts : currentParts" :key="part.id" class="border-b">
                  <td class="p-4">
                    <NuxtLink :to="`/item/${part.entityId}`" class="hover:underline">
                      {{ part.entityName }}
                    </NuxtLink>
                  </td>
                  <td class="p-4">{{ part.position }}</td>
                  <td class="p-4">{{ new Date(part.addedDate).toLocaleDateString() }}</td>
                  <td v-if="showHistory" class="p-4">
                    {{ part.removedDate ? new Date(part.removedDate).toLocaleDateString() : "—" }}
                  </td>
                  <td class="p-4 text-right">
                    <Button v-if="!part.removedDate" size="sm" variant="outline" @click="unmountPart(part)">
                      {{ $t("setups.parts.unmount") }}
                    </Button>
                    <Button v-else size="sm" variant="destructive" @click="purgePart(part)">
                      {{ $t("setups.parts.purge") }}
                    </Button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </BaseCard>

        <!-- Sessions -->
        <BaseCard collapsable>
          <template #title> {{ $t("setups.sessions.title") }} ({{ sessions.length }}) </template>
          <template #title-actions>
            <div class="mt-2">
              <Button size="sm" @click="showSession = true">{{ $t("setups.sessions.log") }}</Button>
            </div>
          </template>

          <p v-if="sessions.length === 0" class="px-6 pb-4 text-foreground/70">
            {{ $t("setups.sessions.empty") }}
          </p>

          <div v-else class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead>
                <tr class="border-b text-left text-foreground/70">
                  <th class="p-4">{{ $t("setups.sessions.date") }}</th>
                  <th class="p-4">{{ $t("setups.sessions.duration_minutes") }}</th>
                  <th class="p-4">{{ $t("setups.sessions.location") }}</th>
                  <th class="p-4">{{ $t("setups.sessions.conditions") }}</th>
                  <th class="p-4">{{ $t("setups.sessions.notes") }}</th>
                  <th class="p-4">{{ $t("setups.sessions.rating") }}</th>
                  <th class="p-4" />
                </tr>
              </thead>
              <tbody>
                <tr v-for="session in sessions" :key="session.id" class="border-b">
                  <td class="p-4">{{ new Date(session.date).toLocaleDateString() }}</td>
                  <td class="p-4">{{ session.durationMinutes }}</td>
                  <td class="p-4">{{ session.location }}</td>
                  <td class="p-4">{{ session.conditions }}</td>
                  <td class="p-4">{{ session.notes }}</td>
                  <td class="p-4">{{ session.rating }}</td>
                  <td class="p-4 text-right">
                    <Button size="sm" variant="destructive" @click="deleteSession(session)">
                      {{ $t("global.delete") }}
                    </Button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </BaseCard>

        <!-- Edit dialog -->
        <DialogRoot v-model:open="showEdit">
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{{ $t("setups.edit") }}</DialogTitle>
            </DialogHeader>
            <form class="flex flex-col gap-4" @submit.prevent="saveEdit">
              <div class="flex flex-col gap-2">
                <Label for="edit-name">{{ $t("global.name") }}</Label>
                <Input id="edit-name" v-model="editForm.name" required />
              </div>
              <div class="flex flex-col gap-2">
                <Label for="edit-description">{{ $t("global.description") }}</Label>
                <Input id="edit-description" v-model="editForm.description" />
              </div>
              <div class="flex flex-col gap-2">
                <Label for="edit-sport">{{ $t("setups.sport") }}</Label>
                <Input id="edit-sport" v-model="editForm.sport" />
              </div>
              <Label class="flex cursor-pointer items-center gap-2">
                <Switch v-model="editForm.isActive" />
                {{ $t("setups.active") }}
              </Label>
              <DialogFooter>
                <Button type="button" variant="outline" @click="showEdit = false">{{ $t("global.cancel") }}</Button>
                <Button type="submit" :disabled="saving">{{ $t("global.save") }}</Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </DialogRoot>

        <!-- Delete confirmation -->
        <DialogRoot v-model:open="showDelete">
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{{ $t("global.delete") }}</DialogTitle>
            </DialogHeader>
            <p>{{ $t("setups.delete_confirm") }}</p>
            <DialogFooter>
              <Button variant="outline" @click="showDelete = false">{{ $t("global.cancel") }}</Button>
              <Button variant="destructive" :disabled="deleting" @click="deleteSetup">
                {{ $t("global.delete") }}
              </Button>
            </DialogFooter>
          </DialogContent>
        </DialogRoot>

        <!-- Mount part dialog -->
        <DialogRoot v-model:open="showMount">
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{{ $t("setups.parts.mount") }}</DialogTitle>
            </DialogHeader>
            <form class="flex flex-col gap-4" @submit.prevent="mountPart">
              <div class="flex flex-col gap-2">
                <Label for="mount-entity">{{ $t("setups.parts.entity") }}</Label>
                <select
                  id="mount-entity"
                  v-model="mountForm.entityId"
                  class="flex h-9 w-full rounded-md border bg-background px-3 text-sm"
                  required
                >
                  <option value="" disabled>{{ $t("setups.parts.entity_placeholder") }}</option>
                  <option v-for="entity in allEntities" :key="entity.id" :value="entity.id">
                    {{ entity.name }}
                  </option>
                </select>
              </div>
              <div class="flex flex-col gap-2">
                <Label for="mount-position">{{ $t("setups.parts.position") }}</Label>
                <Input id="mount-position" v-model="mountForm.position" />
              </div>
              <DialogFooter>
                <Button type="button" variant="outline" @click="showMount = false">{{ $t("global.cancel") }}</Button>
                <Button type="submit" :disabled="mounting">{{ $t("global.add") }}</Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </DialogRoot>

        <!-- Log session dialog -->
        <DialogRoot v-model:open="showSession">
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{{ $t("setups.sessions.log") }}</DialogTitle>
            </DialogHeader>
            <form class="flex flex-col gap-4" @submit.prevent="logSession">
              <div class="flex flex-col gap-2">
                <Label for="session-date">{{ $t("setups.sessions.date") }}</Label>
                <Input id="session-date" v-model="sessionForm.date" type="date" required />
              </div>
              <div class="flex flex-col gap-2">
                <Label for="session-duration">{{ $t("setups.sessions.duration_minutes") }}</Label>
                <Input id="session-duration" v-model.number="sessionForm.durationMinutes" type="number" min="0" />
              </div>
              <div class="flex flex-col gap-2">
                <Label for="session-location">{{ $t("setups.sessions.location") }}</Label>
                <Input id="session-location" v-model="sessionForm.location" />
              </div>
              <div class="flex flex-col gap-2">
                <Label for="session-conditions">{{ $t("setups.sessions.conditions") }}</Label>
                <Input id="session-conditions" v-model="sessionForm.conditions" />
              </div>
              <div class="flex flex-col gap-2">
                <Label for="session-notes">{{ $t("setups.sessions.notes") }}</Label>
                <Input id="session-notes" v-model="sessionForm.notes" />
              </div>
              <div class="flex flex-col gap-2">
                <Label for="session-rating">{{ $t("setups.sessions.rating") }}</Label>
                <Input id="session-rating" v-model.number="sessionForm.rating" type="number" min="0" max="5" />
              </div>
              <DialogFooter>
                <Button type="button" variant="outline" @click="showSession = false">{{ $t("global.cancel") }}</Button>
                <Button type="submit" :disabled="loggingSession">{{ $t("global.save") }}</Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </DialogRoot>
      </template>
    </BaseContainer>
  </div>
</template>
