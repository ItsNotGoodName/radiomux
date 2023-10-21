import { Switch } from "@kobalte/core"
import { styled } from "@macaron-css/solid"
import { ComponentProps } from "solid-js"
import { mixin, theme, tw } from "./theme"

export const SwitchRoot = Switch.Root
export const SwitchInput = Switch.Input

const TheSwitchControl = styled(Switch.Control, {
  base: {
    ...tw.transitionColors,
    display: "inline-flex",
    height: "24px",
    width: "44px",
    flexShrink: "0",
    cursor: "pointer",
    alignItems: "center",
    paddingLeft: "2px", // TODO: this should not be needed
    borderRadius: theme.borderRadius.full,
    borderWidth: theme.space[2],
    borderColor: "transparent",
    ":disabled": {
      cursor: "not-allowed",
      opacity: "50%"
    },
    background: theme.color.input,
    selectors: {
      ['&[data-checked]']: {
        background: theme.color.primary
      },
    }
  }
})

const TheSwitchThumb = styled(Switch.Thumb, {
  base: {
    ...mixin.size("5"),
    ...tw.shadowLg,
    ...tw.transitionTransform,
    pointerEvents: "none",
    display: "block",
    borderRadius: theme.borderRadius.full,
    background: theme.color.background,
    transform: "translateX(0)",
    selectors: {
      ['&[data-checked]']: {
        transform: `translateX(${theme.space[5]})`
      },
    }
  }
})

export function SwitchControl(props: ComponentProps<typeof Switch.Control>) {
  return (
    <TheSwitchControl {...props}>
      <TheSwitchThumb />
    </TheSwitchControl >
  )
}

export const SwitchLabel = Switch.Label
export const SwitchDescription = Switch.Description
export const SwitchErrorMessage = Switch.ErrorMessage
