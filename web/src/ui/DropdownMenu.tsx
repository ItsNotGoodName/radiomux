import { DropdownMenu } from "@kobalte/core";
import { styled } from "@macaron-css/solid";
import { animation, theme, tw } from "./theme";

export const DropdownMenuRoot = DropdownMenu.Root
export const DropdownMenuTrigger = DropdownMenu.Trigger
export const DropdownMenuIcon = DropdownMenu.Icon
export const DropdownMenuPortal = DropdownMenu.Portal

export const DropdownMenuContent = styled(DropdownMenu.Content, {
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
  variants: {
    variant: {
      default: {
        display: "flex",
        flexDirection: "column",
        padding: theme.space[2],
        width: theme.space[48]
      }
    }
  },
  defaultVariants: {
    variant: "default"
  }
})

export const DropdownMenuArrow = DropdownMenu.Arrow
export const DropdownMenuSeparator = DropdownMenu.Separator
export const DropdownMenuGroup = DropdownMenu.Group
export const DropdownMenuGroupLabel = DropdownMenu.GroupLabel
export const DropdownMenuSub = DropdownMenu.Sub
export const DropdownMenuSubTrigger = DropdownMenu.SubTrigger
export const DropdownMenuSubContent = DropdownMenu.SubContent
export const DropdownMenuItem = DropdownMenu.Item
