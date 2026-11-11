# Frontend

- Nuxt 4 / Vue 3, pnpm. Gate: `pnpm run typecheck` before commit (`pnpm run lint:fix` for style).
- Theme state lives in `frontend/composables/use-preferences.ts` (`themeLight`/`themeDark`) and the generated `frontend/lib/api/types/data-contracts.ts` — new settings types go through the data contracts.
- The Builds kill-switch (`ShowBuilds` group flag) gates the nav entry (`frontend/layouts/default.vue`), the `/builds` route (`frontend/middleware/builds-enabled.ts`) and the sport field (`frontend/components/Entity/CreateModal.vue`). Any new builds-adjacent UI must respect `useBuildsEnabled()`.
- Markdown/content rendering and i18n follow existing patterns in `frontend/` — check neighbouring components before inventing new ones.
