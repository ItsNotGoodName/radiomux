import { createEffect, createSignal, untrack } from "solid-js";
import { darkClass, lightClass, theme } from "./theme";
import { style } from "@macaron-css/core";

const THEME_KEY = "theme";

export const AUTO_MODE = "auto";
export const LIGHT_MODE = "light";
export const DARK_MODE = "dark";

const query = window.matchMedia("(prefers-color-scheme: dark)");

const [modeAuto, setModeAuto] = createSignal(
  query.matches ? DARK_MODE : LIGHT_MODE
);

query.addEventListener("change", (e: MediaQueryListEvent) => {
  setModeAuto(e.matches ? DARK_MODE : LIGHT_MODE);
});

function get(): string {
  return localStorage.getItem(THEME_KEY) ?? AUTO_MODE;
}

function set(theme: string) {
  setMode(theme);
  theme === AUTO_MODE
    ? localStorage.removeItem(THEME_KEY)
    : localStorage.setItem(THEME_KEY, theme);
}

const [mode, setMode] = createSignal(get());
export const themeMode = mode;

export const toggleThemeMode = () => {
  const theme = untrack(mode);
  if (theme == LIGHT_MODE) {
    set(DARK_MODE);
  } else if (theme == DARK_MODE) {
    set(AUTO_MODE);
  } else {
    set(LIGHT_MODE);
  }
};

export const themeModeClass = () => {
  if (themeMode() == AUTO_MODE) {
    return modeAuto() == DARK_MODE ? darkClass : lightClass;
  }
  return themeMode() == DARK_MODE ? darkClass : lightClass;
};

export const useTheme = () => {
  return createEffect(() => {
    document.getElementsByTagName("body")![0].className = themeModeClass() + " " + style({
      background: theme.color.background,
      color: theme.color.foreground,
    })
  })
}
