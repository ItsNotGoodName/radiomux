import { Options } from '@popperjs/core';
import usePopper from 'solid-popper';
import { JSX, Setter, createSignal } from "solid-js";
import { themeModeClass } from './theme-mode';
import Dismiss from 'solid-dismiss';
import { Portal } from 'solid-js/web';

type Props = {
  options?: Partial<Options>,
  button: (ref: Setter<HTMLElement | undefined>) => JSX.Element,
  children: (ref: Setter<HTMLElement | undefined>, setOpen: Setter<boolean>) => JSX.Element
}

export function Dropdown(props: Props) {
  const [anchor, setAnchor] = createSignal<HTMLElement>();
  const [popper, setPopper] = createSignal<HTMLElement>();

  usePopper(anchor, popper, props.options);

  const [open, setOpen] = createSignal(false);

  return (
    <>
      {props.button(setAnchor)}
      <Portal>
        <Dismiss menuButton={anchor()} open={open} setOpen={setOpen} class={themeModeClass()}>
          {props.children(setPopper, setOpen)}
        </Dismiss>
      </Portal>
    </>
  )
}
