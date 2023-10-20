import { Popover } from "@kobalte/core";
import { styled } from "@macaron-css/solid";
import { animation, theme, tw } from "./theme";

export const PopoverRoot = Popover.Root
export const PopoverTrigger = Popover.Trigger
export const PopoverPortal = Popover.Portal

export const PopoverContent = styled(Popover.Content, {
  base: {
    ...tw.shadowMd,
    zIndex: 100,
    background: theme.color.popover,
    color: theme.color.popoverForeground,
    border: `1px solid ${theme.color.border}`,
    borderRadius: theme.borderRadius.lg,
    maxWidth: "var(--kb-popper-content-available-width)",
    transformOrigin: "var(--kb-menu-content-transform-origin)",
    animation: `${animation.contentHide} 150ms ease-in forwards`,
    selectors: {
      '&[data-expanded]': {
        animation: `${animation.contentShow} 150ms ease-out`,
      },
    },
  },
})

export const PopoverArrow = Popover.Arrow
export const PopoverCloseButton = Popover.CloseButton
export const PopoverTitle = Popover.Title
export const PopoverDescription = Popover.Description
