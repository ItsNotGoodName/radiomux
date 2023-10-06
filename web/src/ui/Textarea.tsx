import { styled } from "@macaron-css/solid";
import { theme, tw } from "./theme";

export const Textarea = styled("textarea", {
  base: {
    ...tw.textSm,
    display: "flex",
    minHeight: "80px",
    width: "100%",
    borderRadius: theme.borderRadius.md,
    border: `1px solid ${theme.color.input}`,
    background: theme.color.background,
    color: theme.color.foreground,
    padding: `${theme.space[2]} ${theme.space[3]}`,
    // focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2
    "::placeholder": {
      color: theme.color.mutedForeground
    },
    ":disabled": {
      cursor: "not-allowed",
      opacity: "50%"
    }
  }
})
