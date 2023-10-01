import { Route, A } from '@solidjs/router'
import { Component } from "solid-js";
import { Ui } from './Ui';
import { style } from '@macaron-css/core';
import { theme } from '~/ui/theme';

const Home: Component = () => {
  return (
    <div class={style({ padding: theme.space[4] })}>
      <ul>
        <li><A href='./ui'>Ui</A></li>
        <li><A href='./player'>Player</A></li>
      </ul>
    </div>
  )
}

export const Debug: Component = () => {
  return (
    <>
      <Route path="/*" component={Home} />
      <Route path="/ui" component={Ui} />
    </>
  )
}

