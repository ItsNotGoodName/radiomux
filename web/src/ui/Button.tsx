import { styled } from "@macaron-css/solid";
import { theme } from "./theme";

export const Button = styled("button", {
  base: {
    background: theme.color.background,
    color: theme.color.foreground,
    border: "none",
    cursor: "pointer",
    display: "inline-flex",
    alignItems: "center",
    justifyContent: "center",
    borderRadius: theme.borderRadius.ok,
    touchAction: "manipulation",
    userSelect: "none",
    ":disabled": {
      pointerEvents: "none",
      opacity: "50%"
    }
  },
  variants: {
    variant: {
      default: {
        background: theme.color.primary,
        color: theme.color.primaryForeground,
        selectors: {
          ["&:hover:enabled"]: {
            background: theme.color.primaryHover,
          },
        },
      },
      destructive: {
        background: theme.color.destructive,
        color: theme.color.destructiveForeground,
        selectors: {
          ["&:hover:enabled"]: {
            background: theme.color.destructiveHover,
          },
        },
      },
      outline: {
        background: theme.color.background,
        border: `1px solid ${theme.color.input}`,
        selectors: {
          ["&:hover:enabled"]: {
            background: theme.color.accent,
            color: theme.color.accentForeground,
          },
        },
      },
      secondary: {
        background: theme.color.secondary,
        color: theme.color.secondaryForeground,
        selectors: {
          ["&:hover:enabled"]: {
            background: theme.color.secondaryHover,
          },
        },
      },
      ghost: {
        selectors: {
          ["&:hover:enabled"]: {
            background: theme.color.accent,
            color: theme.color.accentForeground,
          },
        },
      },
      link: {
        color: theme.color.primary,
        textUnderlineOffset: theme.space[1],
        selectors: {
          ["&:hover:enabled"]: {
            textDecoration: "underline"
          },
        },
      }
    },
    size: {
      default: {
        height: theme.space["10"],
        padding: `${theme.space['2']} ${theme.space['4']}`,
      },
      sm: {
        height: theme.space[9],
        paddingLeft: theme.space[3],
        paddingRight: theme.space[3]

      },
      lg: {
        height: theme.space[11],
        paddingLeft: theme.space[8],
        paddingRight: theme.space[8],

      },
      icon: {
        height: theme.space[10],
        width: theme.space[10],
      }
    }
  },
  defaultVariants: {
    variant: "default",
    size: "default"
  }
})

