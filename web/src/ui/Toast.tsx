import { Toast, toaster } from "@kobalte/core";
import { styled } from "@macaron-css/solid";
import { animation, minScreen, mixin, theme, tw } from "./theme";
import { keyframes } from "@macaron-css/core";
import { ComponentProps, JSX } from "solid-js";
import { RiSystemCloseLine } from "solid-icons/ri";

export const ToastRegion = Toast.Region

export const ToastList = styled(Toast.List, {
  base: {
    position: "fixed",
    top: 0,
    left: 0,
    right: 0,
    zIndex: 150,
    display: "flex",
    maxHeight: "100vh",
    width: "100%",
    flexDirection: "column",
    padding: theme.space[4],
    gap: theme.space[2],
    margin: 0,
    "@media": {
      [minScreen.sm]: {
        left: "auto",
      },
      [minScreen.md]: {
        maxWidth: "420px"
      },
    },
  }
})

const slideIn = keyframes({
  from: {
    transform: `translateX(100%)`
  },
  to: {
    transform: "translateX(var(--kb-toast-swipe-end-x))"
  }
})

const slideRight = keyframes({
  from: {
    transform: "translateX(var(--kb-toast-swipe-end-x))"
  },
  to: {
    transform: "translateX(100%)",
  }
})

export const ToastRoot = styled(Toast.Root, {
  base: {
    ...tw.shadowLg,
    pointerEvents: "auto",
    position: "relative",
    width: "100%",
    overflow: "hidden",
    borderRadius: theme.borderRadius.md,
    transition: "all",
    selectors: {
      '&[data-opened]': {
        animation: `${slideIn} 150ms cubic-bezier(0.16, 1, 0.3, 1);`
      },
      '&[data-closed]': {
        animation: `${animation.overlayHide} 100ms ease-in`,

      },
      '&[data-swipe="move"]': {
        transform: "translateX(var(--kb-toast-swipe-move-x))"
      },
      '&[data-swipe="cancel"]': {
        transform: "translateX(0)",
        transition: "transform 200ms ease-out",
      },
      '&[data-swipe="end"]': {
        animation: `${slideRight} 100ms ease-out`,
      },
    },
  },
  variants: {
    variant: {
      default: {
        background: theme.color.background,
        color: theme.color.foreground,
        border: `1px solid ${theme.color.border}`,
      },
      destructive: {
        background: theme.color.destructive,
        color: theme.color.destructiveForeground,
      }
    }
  },
  defaultVariants: {
    variant: "default"
  }
})

export const ToastContent = styled("div", {
  base: {
    ...mixin.stack("2"),
    padding: theme.space[4],
  }
})

export const TheToastCloseButton = styled(Toast.CloseButton, {
  base: {
    ...tw.transitionOpacity,
    display: "flex",
    background: "transparent",
    color: theme.color.mutedForeground,
    padding: theme.space[1],
    border: "none",
    cursor: "pointer",
    position: "absolute",
    right: theme.space[2],
    top: theme.space[2],
    borderRadius: theme.borderRadius.sm,
    opacity: "70%",
    ":hover": {
      color: theme.color.foreground,
      opacity: "100%"
    },
    // BUG: this works with `pnpm run build` but not `pnpm run dev` because ToastRoot uses an imported element
    // selectors: {
    //   [`${ToastRoot}:hover &`]: {
    //     opacity: "100%",
    //   },
    // TODO: add destructive variant but is blocked by the bug above
    // },
  }
})

export function ToastCloseButton(props: ComponentProps<typeof Toast.CloseButton>) {
  return (
    <TheToastCloseButton title="Close" {...props}>
      <RiSystemCloseLine />
    </TheToastCloseButton>
  )
}

export const ToastTitle = styled(Toast.Title, {
  base: {
    ...tw.textSm,
    fontWeight: 700
  }
})

export const ToastDescription = styled(Toast.Description, {
  base: {
    ...tw.textSm,
    opacity: "90%"
  }
})

export const ToastProgressTrack = styled(Toast.ProgressTrack, {
  base: {
    height: theme.space[2],
    width: "100%",
    background: theme.color.primaryForeground,
    borderRadius: theme.borderRadius.ok
  }
})

export const ToastProgressFill = styled(Toast.ProgressFill, {
  base: {
    background: theme.color.primary,
    borderRadius: theme.borderRadius.ok,
    height: "100%",
    width: "var(--kb-toast-progress-fill-width)",
    transition: "width 250ms linear"
  }
})

function show(message: string) {
  return toaster.show(props => (
    <ToastRoot toastId={props.toastId}>
      <ToastContent>
        <ToastCloseButton />
        {message}
      </ToastContent>
    </ToastRoot>
  ));
}

function success(message: string) {
  return toaster.show(props => (
    <ToastRoot toastId={props.toastId}>
      <ToastContent>
        <ToastCloseButton />
        {message}
      </ToastContent>
    </ToastRoot>
  ));
}

function error(title: string, message: string) {
  return toaster.show(props => (
    <ToastRoot toastId={props.toastId} variant="destructive">
      <ToastContent>
        <ToastCloseButton />
        <ToastTitle>{title}</ToastTitle>
        <ToastDescription>
          {message}
        </ToastDescription>
      </ToastContent>
    </ToastRoot>
  ));
}

function custom(ele: () => JSX.Element) {
  return toaster.show(props => <ToastRoot toastId={props.toastId}>{ele as any}</ToastRoot>);
}

function dismiss(id: number) {
  return toaster.dismiss(id);
}

export const toast = {
  show,
  success,
  error,
  custom,
  dismiss,
};
