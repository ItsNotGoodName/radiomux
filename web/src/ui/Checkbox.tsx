import { styled } from "@macaron-css/solid";
import { mixin, theme, } from "./theme";
import { JSX, Show, createEffect } from "solid-js";
import { RiSystemCheckLine } from "solid-icons/ri";

const TheCheckbox = styled("button", {
  base: {
    ...mixin.size("4"),
    background: "transparent",
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    borderRadius: theme.borderRadius.sm,
    border: `1px solid ${theme.color.primary}`,
    flexShrink: "0",
    cursor: "pointer",
    // ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2
    ":disabled": {
      cursor: "not-allowed",
      opacity: "50%"
    },
    selectors: {
      '&[data-state="checked"]': {
        background: theme.color.primary,
        color: theme.color.primaryForeground,
      },
    }
  }
})

type Props = {
  checked: boolean
  setChecked: (checked: boolean) => void
} & Omit<JSX.HTMLAttributes<HTMLButtonElement>, "role" | "onClick">

export function Checkbox(props: Props) {
  const onClick = () => {
    props.setChecked(!props.checked)
  }

  let ref: HTMLButtonElement
  createEffect(() => {
    if (props.checked) {
      ref.dataset["state"] = "checked"
    } else {
      delete ref.dataset["state"]
    }
  })

  return <TheCheckbox
    ref={ref!}
    role="checkbox"
    onClick={onClick}
    {...props}
  >
    <Show when={props.checked}>
      <RiSystemCheckLine />
    </Show>
  </TheCheckbox>
}
