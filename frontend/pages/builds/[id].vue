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
  import type {
    EntitySummary,
    BuildOut,
    BuildPartOut,
    BuildStatus,
    UsageSessionOut,
  } from "~~/lib/api/types/data-contracts";

  definePageMeta({
    middleware: ["auth"],
  });

  const { t } = useI18n();
  const route = useRoute();
  const api = useUserApi();

  const { data: sportsData } = useAsyncData("build-sports", async () => {
    const { data } = await api.builds.getSports();
    return data ?? [];
  });
  const sports = computed(() => (Array.isArray(sportsData.value) ? sportsData.value : []));

  const buildId = computed(() => route.params.id as string);

  const build = ref<BuildOut | null>(null);
  const parts = ref<BuildPartOut[]>([]);
  const sessions = ref<UsageSessionOut[]>([]);
  const loading = ref(true);

  async function loadAll() {
    loading.value = true;

    const [buildRes, partsRes, sessionsRes] = await Promise.all([
      api.builds.get(buildId.value),
      api.builds.getParts(buildId.value),
      api.builds.getSessions(buildId.value),
    ]);

    if (buildRes.error) {
      toast.error(t("builds.toast.failed_load"));
      loading.value = false;
      return;
    }
    build.value = buildRes.data;
    parts.value = partsRes.error ? [] : partsRes.data;
    sessions.value = sessionsRes.error ? [] : sessionsRes.data;
    loading.value = false;
  }

  onMounted(loadAll);

  const currentParts = computed(() => parts.value.filter(p => !p.removedDate));

  // Positions are strictly scoped to the build's sport. A build without a
  // sport is generic, so it may use any sport's positions (deduplicated,
  // first-seen order). A part's own sport never restricts this list.
  const currentSportPositions = computed(() => {
    if (!build.value?.sport) {
      const seen = new Set<string>();
      const all: string[] = [];
      for (const s of sports.value) {
        for (const pos of s.positions) {
          if (!seen.has(pos)) {
            seen.add(pos);
            all.push(pos);
          }
        }
      }
      return all;
    }
    return sports.value.find(s => s.slug === build.value?.sport)?.positions ?? [];
  });

  // The part picker is a flat list: the build's sport already constrains the
  // position dropdown, and cross-sport parts are warned about on mount, so
  // grouping by sport adds no value here.

  // Prettified label for a raw position slug ("front_wheel" -> "Front Wheel").
  // Values written to the API stay the raw slug.
  function humanizePosition(pos: string): string {
    return pos
      .split("_")
      .map(word => (word ? word.charAt(0).toUpperCase() + word.slice(1) : word))
      .join(" ");
  }

  // --- edit / delete build ---
  const showEdit = ref(false);
  const saving = ref(false);
  const editForm = reactive({ name: "", description: "", sport: "", status: "assembled" as BuildStatus });

  function openEdit() {
    if (!build.value) {
      return;
    }
    editForm.name = build.value.name;
    editForm.description = build.value.description;
    editForm.sport = build.value.sport;
    editForm.status = build.value.status;
    showEdit.value = true;
  }

  async function saveEdit() {
    if (!build.value || !editForm.name) {
      return;
    }
    saving.value = true;
    const { error } = await api.builds.update(build.value.id, { ...editForm });
    saving.value = false;
    if (error) {
      toast.error(t("builds.toast.failed_update"));
      return;
    }
    toast.success(t("builds.toast.updated"));
    showEdit.value = false;
    await loadAll();
  }

  const showDelete = ref(false);
  const deleting = ref(false);

  async function deleteBuild() {
    if (!build.value) {
      return;
    }
    deleting.value = true;
    const { error } = await api.builds.delete(build.value.id);
    deleting.value = false;
    if (error) {
      toast.error(t("builds.toast.failed_delete"));
      return;
    }
    toast.success(t("builds.toast.deleted"));
    await navigateTo("/builds");
  }

  // --- parts ---
  const showMount = ref(false);
  const mounting = ref(false);
  const allEntities = ref<EntitySummary[]>([]);
  const mountForm = reactive({ entityId: "", position: "" });

  // Pre-select the most likely position for the chosen part, based on its
  // name (parts are bought/used by category: a "Deck" goes on the deck,
  // wheels and bearings mount to a truck). Cross-sport parts picked via
  // "Show All" get no guess — the user decides deliberately.
  function guessPosition(entity: EntitySummary | undefined): string {
    if (!entity) {
      return "";
    }
    if (build.value?.sport && entity.sport && build.value.sport !== entity.sport) {
      return "";
    }
    const name = entity.name.toLowerCase();
    if (name.includes("deck") || name.includes("grip")) {
      return "deck";
    }
    if (name.includes("truck") || name.includes("wheel") || name.includes("bearing")) {
      return name.includes("rear") ? "rear_truck" : "front_truck";
    }
    return "";
  }

  watch(
    () => mountForm.entityId,
    id => {
      if (!id) {
        mountForm.position = "";
        return;
      }
      mountForm.position = guessPosition(allEntities.value.find(e => e.id === id));
    }
  );

  const showAllParts = ref(false);

  // Strict mode: a sported build only lists its own sport's parts plus
  // generic (sportless) parts. "Show All" lifts the filter for rare
  // cross-sport mounts; the server-side mismatch warning still applies.
  const selectableEntities = computed(() => {
    if (showAllParts.value || !build.value?.sport) {
      return allEntities.value;
    }
    return allEntities.value.filter(e => !e.sport || e.sport === build.value?.sport);
  });

  async function openMount() {
    showMount.value = true;
    showAllParts.value = false;
    if (allEntities.value.length === 0) {
      const { data, error } = await api.items.getAll();
      if (error) {
        toast.error(t("builds.toast.failed_load"));
        return;
      }
      allEntities.value = data.items;
    }
  }

  async function mountPart() {
    if (!mountForm.entityId || !mountForm.position) {
      return;
    }
    mounting.value = true;
    const { data, error } = await api.builds.addPart(buildId.value, {
      entityId: mountForm.entityId,
      position: mountForm.position,
    });
    mounting.value = false;
    if (error) {
      toast.error(t("builds.toast.part_failed_add"));
      return;
    }
    if (data && data.warnings && data.warnings.length > 0) {
      // Success first so the warning stacks on top, then let the warning
      // persist until dismissed instead of auto-dismissing behind it.
      toast.success(t("builds.toast.part_added"));
      toast.warning(data.warnings.map(code => t(`builds.warning.${code}`)).join(" "), {
        duration: Infinity,
      });
    } else {
      toast.success(t("builds.toast.part_added"));
    }
    showMount.value = false;
    mountForm.entityId = "";
    mountForm.position = "";
    await loadAll();
  }

  async function unmountPart(part: BuildPartOut) {
    const { error } = await api.builds.removePart(buildId.value, part.id);
    if (error) {
      toast.error(t("builds.toast.part_failed_remove"));
      return;
    }
    toast.success(t("builds.toast.part_removed"));
    await loadAll();
  }

  async function purgePart(part: BuildPartOut) {
    const { error } = await api.builds.purgePart(buildId.value, part.id);
    if (error) {
      toast.error(t("builds.toast.part_failed_delete"));
      return;
    }
    toast.success(t("builds.toast.part_deleted"));
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
    const { error } = await api.builds.createSession(buildId.value, {
      ...sessionForm,
      date: new Date(sessionForm.date).toISOString(),
    });
    loggingSession.value = false;
    if (error) {
      toast.error(t("builds.toast.session_failed_log"));
      return;
    }
    toast.success(t("builds.toast.session_logged"));
    showSession.value = false;
    sessionForm.durationMinutes = 0;
    sessionForm.location = "";
    sessionForm.conditions = "";
    sessionForm.notes = "";
    sessionForm.rating = 0;
    await loadAll();
  }

  async function deleteSession(session: UsageSessionOut) {
    const { error } = await api.builds.deleteSession(buildId.value, session.id);
    if (error) {
      toast.error(t("builds.toast.session_failed_delete"));
      return;
    }
    toast.success(t("builds.toast.session_deleted"));
    await loadAll();
  }
</script>

<template>
  <div>
    <BaseContainer class="flex flex-col gap-6">
      <div v-if="loading">{{ $t("global.loading") }}</div>

      <template v-else-if="build">
        <div class="flex items-start justify-between gap-4">
          <div class="flex flex-col gap-2">
            <NuxtLink to="/builds" class="text-sm text-foreground/70 hover:underline">
              ← {{ $t("menu.builds") }}
            </NuxtLink>
            <BaseSectionHeader class="!mb-0">
              {{ build.name }}
              <span
                class="ml-3 rounded-full px-2 py-1 text-xs"
                :class="
                  build.status === 'assembled' ? 'bg-green-500/15 text-green-600' : 'bg-muted text-muted-foreground'
                "
              >
                {{ $t(`builds.status.${build.status}`) }}
              </span>
            </BaseSectionHeader>
            <p v-if="build.description" class="text-foreground/70">{{ build.description }}</p>
          </div>
          <ButtonGroup>
            <Button variant="outline" @click="openEdit">{{ $t("global.edit") }}</Button>
            <Button variant="destructive" @click="showDelete = true">{{ $t("global.delete") }}</Button>
          </ButtonGroup>
        </div>

        <!-- Parts -->
        <BaseCard collapsable>
          <template #title> {{ $t("builds.parts.title") }} ({{ currentParts.length }}) </template>
          <template #title-actions>
            <div class="mt-2 flex items-center gap-4">
              <Label class="flex cursor-pointer items-center gap-2">
                <Switch v-model="showHistory" />
                {{ $t("builds.parts.history") }}
              </Label>
              <Button size="sm" @click="openMount">{{ $t("builds.parts.mount") }}</Button>
            </div>
          </template>

          <p v-if="currentParts.length === 0" class="px-6 pb-4 text-foreground/70">
            {{ $t("builds.parts.empty") }}
          </p>

          <div v-else class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead>
                <tr class="border-b text-left text-foreground/70">
                  <th class="p-4">{{ $t("builds.parts.entity") }}</th>
                  <th class="p-4">{{ $t("builds.parts.position") }}</th>
                  <th class="p-4">{{ $t("builds.parts.mounted") }}</th>
                  <th v-if="showHistory" class="p-4">{{ $t("builds.parts.removed") }}</th>
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
                  <td class="p-4">{{ humanizePosition(part.position) || "—" }}</td>
                  <td class="p-4">{{ new Date(part.addedDate).toLocaleDateString() }}</td>
                  <td v-if="showHistory" class="p-4">
                    {{ part.removedDate ? new Date(part.removedDate).toLocaleDateString() : "—" }}
                  </td>
                  <td class="p-4 text-right">
                    <Button v-if="!part.removedDate" size="sm" variant="outline" @click="unmountPart(part)">
                      {{ $t("builds.parts.unmount") }}
                    </Button>
                    <Button v-else size="sm" variant="destructive" @click="purgePart(part)">
                      {{ $t("builds.parts.purge") }}
                    </Button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </BaseCard>

        <!-- Sessions -->
        <BaseCard collapsable>
          <template #title> {{ $t("builds.sessions.title") }} ({{ sessions.length }}) </template>
          <template #title-actions>
            <div class="mt-2">
              <Button size="sm" @click="showSession = true">{{ $t("builds.sessions.log") }}</Button>
            </div>
          </template>

          <p v-if="sessions.length === 0" class="px-6 pb-4 text-foreground/70">
            {{ $t("builds.sessions.empty") }}
          </p>

          <div v-else class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead>
                <tr class="border-b text-left text-foreground/70">
                  <th class="p-4">{{ $t("builds.sessions.date") }}</th>
                  <th class="p-4">{{ $t("builds.sessions.duration_minutes") }}</th>
                  <th class="p-4">{{ $t("builds.sessions.location") }}</th>
                  <th class="p-4">{{ $t("builds.sessions.conditions") }}</th>
                  <th class="p-4">{{ $t("builds.sessions.notes") }}</th>
                  <th class="p-4">{{ $t("builds.sessions.rating") }}</th>
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
              <DialogTitle>{{ $t("builds.edit") }}</DialogTitle>
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
                <Label for="edit-sport">{{ $t("builds.sport") }}</Label>
                <select
                  id="edit-sport"
                  v-model="editForm.sport"
                  class="rounded-md border bg-background px-3 py-2 text-sm"
                >
                  <option value="">—</option>
                  <option v-for="sport in sports" :key="sport.slug" :value="sport.slug">{{ sport.label }}</option>
                </select>
              </div>
              <div class="flex flex-col gap-2">
                <Label for="edit-status">{{ $t("builds.status.label") }}</Label>
                <select
                  id="edit-status"
                  v-model="editForm.status"
                  class="rounded-md border bg-background px-3 py-2 text-sm"
                >
                  <option value="assembled">{{ $t("builds.status.assembled") }}</option>
                  <option value="disassembled">{{ $t("builds.status.disassembled") }}</option>
                  <option value="retired">{{ $t("builds.status.retired") }}</option>
                </select>
              </div>
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
            <p>{{ $t("builds.delete_confirm") }}</p>
            <DialogFooter>
              <Button variant="outline" @click="showDelete = false">{{ $t("global.cancel") }}</Button>
              <Button variant="destructive" :disabled="deleting" @click="deleteBuild">
                {{ $t("global.delete") }}
              </Button>
            </DialogFooter>
          </DialogContent>
        </DialogRoot>

        <!-- Mount part dialog -->
        <DialogRoot v-model:open="showMount">
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{{ $t("builds.parts.mount") }}</DialogTitle>
            </DialogHeader>
            <form class="flex flex-col gap-4" @submit.prevent="mountPart">
              <div class="flex flex-col gap-2">
                <div class="flex items-center justify-between">
                  <Label for="mount-entity">{{ $t("builds.parts.entity") }}</Label>
                  <Label class="flex cursor-pointer items-center gap-2 text-xs text-foreground/70">
                    <Switch v-model="showAllParts" />
                    {{ $t("builds.parts.show_all") }}
                  </Label>
                </div>
                <select
                  id="mount-entity"
                  v-model="mountForm.entityId"
                  class="flex h-9 w-full rounded-md border bg-background px-3 text-sm"
                  required
                >
                  <option value="" disabled>{{ $t("builds.parts.entity_placeholder") }}</option>
                  <option v-for="entity in selectableEntities" :key="entity.id" :value="entity.id">
                    {{ entity.name }}
                  </option>
                </select>
              </div>
              <div class="flex flex-col gap-2">
                <Label for="mount-position">{{ $t("builds.parts.position") }}</Label>
                <select
                  id="mount-position"
                  v-model="mountForm.position"
                  class="rounded-md border bg-background px-3 py-2 text-sm"
                  required
                >
                  <option value="" disabled>{{ $t("builds.parts.position_placeholder") }}</option>
                  <option v-for="pos in currentSportPositions" :key="pos" :value="pos">
                    {{ humanizePosition(pos) }}
                  </option>
                </select>
              </div>
              <DialogFooter>
                <Button type="button" variant="outline" @click="showMount = false">{{ $t("global.cancel") }}</Button>
                <Button type="submit" :disabled="mounting || !mountForm.position">{{ $t("global.add") }}</Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </DialogRoot>

        <!-- Log session dialog -->
        <DialogRoot v-model:open="showSession">
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{{ $t("builds.sessions.log") }}</DialogTitle>
            </DialogHeader>
            <form class="flex flex-col gap-4" @submit.prevent="logSession">
              <div class="flex flex-col gap-2">
                <Label for="session-date">{{ $t("builds.sessions.date") }}</Label>
                <Input id="session-date" v-model="sessionForm.date" type="date" required />
              </div>
              <div class="flex flex-col gap-2">
                <Label for="session-duration">{{ $t("builds.sessions.duration_minutes") }}</Label>
                <Input id="session-duration" v-model.number="sessionForm.durationMinutes" type="number" min="0" />
              </div>
              <div class="flex flex-col gap-2">
                <Label for="session-location">{{ $t("builds.sessions.location") }}</Label>
                <Input id="session-location" v-model="sessionForm.location" />
              </div>
              <div class="flex flex-col gap-2">
                <Label for="session-conditions">{{ $t("builds.sessions.conditions") }}</Label>
                <Input id="session-conditions" v-model="sessionForm.conditions" />
              </div>
              <div class="flex flex-col gap-2">
                <Label for="session-notes">{{ $t("builds.sessions.notes") }}</Label>
                <Input id="session-notes" v-model="sessionForm.notes" />
              </div>
              <div class="flex flex-col gap-2">
                <Label for="session-rating">{{ $t("builds.sessions.rating") }}</Label>
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
