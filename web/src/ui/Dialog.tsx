import { styled } from "@macaron-css/solid";
import { minScreen, mixin, theme, tw } from "./theme";
import { JSX, ParentProps, Setter, Show, createSignal } from "solid-js";
import { themeModeClass } from "./theme-mode";
import { Portal } from "solid-js/web";
import Dismiss from "solid-dismiss";

const DialogOverlay = styled("div", {
  base: {
    position: "fixed",
    inset: 0,
    zIndex: 50,
    background: theme.color.backgroundOverlay
    // backdrop-blur-sm data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0
  }
})

export const DialogContent = styled("div", {
  base: {
    ...tw.shadowLg,
    position: "fixed",
    left: "50%",
    top: "50%",
    zIndex: 50,
    display: "grid",
    width: "100%",
    maxWidth: theme.size.lg,
    translate: "-50% -50%",
    gap: theme.space[4],
    border: `1px solid ${theme.color.border}`,
    background: theme.color.background,
    color: theme.color.foreground,
    padding: theme.space[6],
    animationDuration: "300ms",
    // data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[state=closed]:slide-out-to-left-1/2 data-[state=closed]:slide-out-to-top-[48%] data-[state=open]:slide-in-from-left-1/2 data-[state=open]:slide-in-from-top-[48%] 
    "@media": {
      [minScreen.sm]: {
        borderRadius: theme.borderRadius.lg,
      },
      [minScreen.md]: {
        width: "100%"
      },
    },
  }
})

export const DialogHeader = styled("div", {
  base: {
    ...mixin.stack("1.5"), // flex flex-col space-y-1.5
    textAlign: "center",
    "@media": {
      [minScreen.sm]: {
        textAlign: "left"
      },
    }
  }
})

export const DialogFooter = styled("div", {
  base: {
    display: "flex",
    flexDirection: "column-reverse",
    "@media": {
      [minScreen.sm]: {
        flexDirection: "row",
        justifyContent: "end",
        gap: theme.space[2]
      },
    }
  }
})

export const DialogTitle = styled("div", {
  base: {
    ...tw.textLg,
    fontWeight: "600",
    lineHeight: "1",
    letterSpacing: "-0.025em"
  }
})

export const DialogDescription = styled("div", {
  base: {
    ...tw.textSm,
    color: theme.color.mutedForeground
  }
})

type Props = {
  button: (ref: Setter<HTMLElement | undefined>) => JSX.Element,
}

export function Dialog(props: ParentProps<Props>) {
  const [open, setOpen] = createSignal(false);
  const [anchor, setAnchor] = createSignal<HTMLElement>();

  return (
    <>
      {props.button(setAnchor)}
      <Portal>
        <div class={themeModeClass()}>
          <Show when={open()}>
            <DialogOverlay />
          </Show>
          <Dismiss menuButton={anchor()} open={open} setOpen={setOpen} >
            {props.children}
          </Dismiss>
        </div>
      </Portal>
    </>
  )
}
