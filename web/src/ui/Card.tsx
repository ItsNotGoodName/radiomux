import { styled } from "@macaron-css/solid";
import { tw, theme } from "./theme";

export const Card = styled("div", {
  base: {
    ...tw.shadowSm,
    background: theme.color.card,
    color: theme.color.cardForeground,
    border: `1px solid ${theme.color.border}`,
    borderRadius: theme.borderRadius.lg,
  }
})

export const CardContent = styled("div", {
  base: {
    padding: `0 ${theme.space[6]} ${theme.space[6]} ${theme.space[6]}`
  }
})

export const CardDescription = styled("p", {
  base: {
    ...tw.textSm,
    color: theme.color.mutedForeground
  }
})

export const CardFooter = styled("div", {
  base: {
    padding: `0 ${theme.space[6]} ${theme.space[6]} ${theme.space[6]}`
  }
})

export const CardHeader = styled("div", {
  base: {
    display: "flex",
    flexDirection: "column",
    padding: theme.space[6],
    gap: theme.space[2] // space-y-1.5
  }
})

export const CardTitle = styled("div", {
  base: {
    ...tw.text2xl,
    fontWeight: "600", // font-semibold
    lineHeight: 1, // leading-none
    letterSpacing: "-0.05em", //  tracking-tighter
  }
})
