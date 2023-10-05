import { styled } from "@macaron-css/solid";
import { theme, tw } from "./theme";

export const Badge = styled("div", {
  base: {
    ...tw.textXs,
    ...tw.transitionColors,
    display: "inline-flex",
    alignItems: "center",
    borderRadius: theme.borderRadius.full,
    padding: `${theme.space["0.5"]} ${theme.space["2.5"]}`,
    fontWeight: "600",

  },
  variants: {
    variant: {
      default: {
        borderColor: "transparent",
        background: theme.color.primary,
        color: theme.color.primaryForeground,
        ":hover": {
          background: theme.color.primaryOverlay,
        }
      },
      secondary: {
        borderColor: "transparent",
        background: theme.color.secondary,
        color: theme.color.secondaryForeground,
        ":hover": {
          background: theme.color.secondaryHover,
        }
      },
      destructive: {
        borderColor: "transparent",
        background: theme.color.destructive,
        color: theme.color.destructiveForeground,
        ":hover": {
          background: theme.color.destructiveOverlay,
        }
      },
      outline: {
        color: theme.color.foreground,
      }
    }
  },
  defaultVariants: {
    variant: "default"
  }
})

