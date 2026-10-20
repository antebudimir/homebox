<script setup lang="ts">
  import { darkThemes } from "~~/lib/data/themes";
  import { SidebarMenuItem, SidebarMenuButton } from "@/components/ui/sidebar";
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
  <SidebarMenuItem>
    <SidebarMenuButton :tooltip="$t('global.toggle_theme')" @click="toggle">
      <MdiWhiteBalanceSunny v-if="isDark" />
      <MdiWeatherNight v-else />
      <span>{{ $t("global.toggle_theme") }}</span>
    </SidebarMenuButton>
  </SidebarMenuItem>
</template>
