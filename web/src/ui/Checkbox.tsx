import { styled } from "@macaron-css/solid";
import { mixin, theme, } from "./theme";
import { RiSystemCheckLine } from "solid-icons/ri";
import { Checkbox } from "@kobalte/core";

export const CheckboxRoot = Checkbox.Root
export const CheckboxInput = Checkbox.Input

export const CheckboxControl = styled(Checkbox.Control, {
  base: {
    ...mixin.size("4"),
    background: "transparent",
    borderRadius: theme.borderRadius.sm,
    border: `1px solid ${theme.color.primary}`,
    cursor: "pointer",
    // ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2
    ":disabled": {
      cursor: "not-allowed",
      opacity: "50%"
    },
  }
})

export const CheckboxIndicator = styled(Checkbox.Indicator, {
  base: {
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    height: "100%",
    selectors: {
      [`&[data-checked]`]: {
        background: theme.color.primary,
        color: theme.color.primaryForeground,
      },
    }
  }
})

export const CheckboxIcon = RiSystemCheckLine
export const CheckboxLabel = Checkbox.Label
export const CheckboxDescription = Checkbox.Description
export const CheckboxErrorMessage = Checkbox.ErrorMessage

