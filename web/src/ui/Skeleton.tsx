import { styled } from "@macaron-css/solid";
import { Skeleton } from "@kobalte/core";
import { keyframes } from "@macaron-css/core";
import { theme } from "./theme";

const skeletonFade = keyframes({
  "0%": {
    opacity: "40%",
  },
  "100%": {
    opacity: "40%",
  },
  "50%": {
    opacity: "100%"
  }
})

export const SkeletonRoot = styled(Skeleton.Root, {
  base: {
    height: "auto",
    width: "100%",
    position: "relative",
    transform: "translateZ(0)",
    selectors: {
      ['&[data-animate="true"]']: {
        animation: `${skeletonFade} 1500ms linear infinite`,
      },
      ['&[data-visible="true"]']: {
        overflow: "hidden"
      },
      ['&[data-visible="true"]::before']: {
        position: "absolute",
        content: "",
        inset: 0,
        zIndex: 10,
        background: theme.color.mutedForeground,
      },
      ['&[data-visible="true"]::after']: {
        position: "absolute",
        content: "",
        inset: 0,
        zIndex: 11,
        background: theme.color.muted,
      },
    }
  }
})
