import { styled } from "@macaron-css/solid";
import { animation, minScreen, mixin, theme, tw } from "./theme";
import { Dialog } from "@kobalte/core";
import { ComponentProps } from "solid-js";
import { RiSystemCloseLine } from "solid-icons/ri";

export const DialogRoot = Dialog.Root
export const DialogTrigger = Dialog.Trigger
export const DialogPortal = Dialog.Portal

export const DialogOverlay = styled(Dialog.Overlay, {
  base: {
    position: "fixed",
    inset: 0,
    zIndex: 50,
    background: theme.color.backgroundOverlay,
    animation: `${animation.overlayHide} 150ms ease 50ms forwards`,
    selectors: {
      '&[data-expanded]': {
        animation: `${animation.overlayShow} 150ms ease`,
      },
    },
  }
})

export const DialogContent = styled(Dialog.Content, {
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
    animation: `${animation.contentHide} 150ms ease-in forwards`,
    selectors: {
      '&[data-expanded]': {
        animation: `${animation.contentShow} 150ms ease-out`,
      },
    },
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
    ...mixin.stack("1.5"),
    textAlign: "center",
    "@media": {
      [minScreen.sm]: {
        textAlign: "left"
      },
    }
  }
})

const TheDialogCloseButton = styled(Dialog.CloseButton, {
  base: {
    ...tw.transitionOpacity,
    display: "flex",
    background: "transparent",
    color: theme.color.mutedForeground,
    padding: 0,
    border: "none",
    cursor: "pointer",
    position: "absolute",
    right: theme.space[4],
    top: theme.space[4],
    borderRadius: theme.borderRadius.sm,
    opacity: "70%",
    ":hover": {
      color: theme.color.foreground,
      opacity: "100%"
    },
    ":disabled": {
      pointerEvents: "none",
    },
  }
})

export function DialogCloseButton(props: Omit<ComponentProps<typeof Dialog.CloseButton>, "title">) {
  return (
    <TheDialogCloseButton title="Close" {...props}>
      <RiSystemCloseLine />
    </TheDialogCloseButton>
  )
}

export const DialogHeaderTitle = styled(Dialog.Title, {
  base: {
    ...tw.textLg,
    fontWeight: "600",
    lineHeight: "1",
    letterSpacing: "-0.025em",
    margin: "0px",
  }
})

export const DialogHeaderDescription = styled(Dialog.Description, {
  base: {
    ...tw.textSm,
    color: theme.color.mutedForeground,
    margin: "0px",
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

