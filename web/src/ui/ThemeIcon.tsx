import {
  RiDesignContrastLine,
  RiWeatherMoonLine,
  RiWeatherSunLine,
} from "solid-icons/ri";
import { Component, Match, Switch } from "solid-js";
import { IconProps } from "solid-icons";
import {
  DARK_MODE,
  LIGHT_MODE,
  themeMode,
} from "./theme-mode";

export const ThemeIcon: Component<IconProps> = (props) => {
  return (
    <Switch fallback={<RiDesignContrastLine {...props} />}>
      <Match when={themeMode() == DARK_MODE}>
        <RiWeatherMoonLine {...props} />
      </Match>
      <Match when={themeMode() == LIGHT_MODE}>
        <RiWeatherSunLine {...props} />
      </Match>
    </Switch>
  );
};
