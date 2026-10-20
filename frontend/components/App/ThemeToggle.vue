<script setup lang="ts">
  import { darkThemes } from "~~/lib/data/themes";
  import { Button } from "@/components/ui/button";
  import MdiWhiteBalanceSunny from "~icons/mdi/white-balance-sunny";
  import MdiWeatherNight from "~icons/mdi/weather-night";
  import { useTheme } from "~/composables/use-theme";

  const preferences = useViewPreferences();
  const { setTheme } = useTheme();

  const isDark = computed(() => darkThemes.includes(preferences.value.theme));

  // Remember the user's last light/dark choices so toggling always returns
  // to the theme they left, instead of a hardcoded pair.
  watch(
    () => preferences.value.theme,
    theme => {
      if (darkThemes.includes(theme)) {
        preferences.value.themeDark = theme;
      } else {
        preferences.value.themeLight = theme;
      }
    }
  );

  const toggle = () => {
    setTheme(isDark.value ? preferences.value.themeLight : preferences.value.themeDark);
  };
</script>

<template>
  <Button
    size="icon"
    variant="outline"
    :aria-label="$t('global.toggle_theme')"
    :title="$t('global.toggle_theme')"
    @click="toggle"
  >
    <MdiWhiteBalanceSunny v-if="isDark" class="size-5" />
    <MdiWeatherNight v-else class="size-5" />
  </Button>
</template>
