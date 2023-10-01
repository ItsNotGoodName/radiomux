import { styled } from "@macaron-css/solid";
import { theme } from "./theme";

export const Input = styled("Input", {
  base: {
    display: "flex",
    height: theme.space[10],
    width: "100%",
    borderRadius: theme.borderRadius.md,
    border: `1px solid ${theme.color.input}`,
    background: theme.color.background,
    padding: `${theme.space[2]} ${theme.space[3]}`,
    // text-sm
    // placeholder:text-muted-foreground
    ":disabled": {
      cursor: "not-allowed",
      opacity: "50%"
    }
  },
})

