import { styled } from "@macaron-css/solid";
import { theme, tw } from "./theme";

export const Input = styled("input", {
  base: {
    ...tw.textSm,
    display: "flex",
    height: theme.space[10],
    width: "100%",
    borderRadius: theme.borderRadius.md,
    border: `1px solid ${theme.color.input}`,
    background: theme.color.background,
    padding: `${theme.space[2]} ${theme.space[3]}`,
    color: theme.color.foreground,
    // file:border-0 file:bg-transparent file:text-sm file:font-medium
    "::placeholder": {
      color: theme.color.mutedForeground
    },
    ":disabled": {
      cursor: "not-allowed",
      opacity: "50%"
    }
  },
})

