import { Separator } from "@kobalte/core";
import { theme } from "./theme";
import { styled } from "@macaron-css/solid";

export const SeparatorRoot = styled(Separator.Root, {
  base: {
    border: "none",
    background: theme.color.border,
    selectors: {
      '&[data-orientation="horizontal"]': {
        height: "1px",
        width: "100%"
      }
    }
  }
})

