import { styled } from "@macaron-css/solid"
import { theme, tw } from "./theme"
import { Alert } from "@kobalte/core"

export const AlertRoot = styled(Alert.Root, {
  base: {
    width: "100%",
    borderRadius: theme.borderRadius.lg,
    border: `1px solid ${theme.color.border}`,
    padding: theme.space[4],
  },
  variants: {
    variant: {
      default: {
        background: theme.color.background,
        text: theme.color.foreground,
      },
      destructive: {
        background: theme.color.destructive,
        borderColor: theme.color.destructiveBorder,
        text: theme.color.destructive,
      }
    }
  },
  defaultVariants: {
    variant: "default"
  }
})

export const AlertTitle = styled("h5", {
  base: {
    fontWeight: 700,
    lineHeight: 1,
    letterSpacing: "-0.025em",
    margin: `0 0 ${theme.space[1]} 0`,
  }
})

export const AlertDescription = styled("div", {
  base: {
    ...tw.textSm,
    margin: 0,
  }
})
