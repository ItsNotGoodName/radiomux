import { Outlet, Route, Routes } from '@solidjs/router'
import { Debug } from './pages/debug';
import { Pages } from './pages';
import { themeModeClass } from './ui/theme-mode';
import { theme } from './ui/theme';
import { styled } from '@macaron-css/solid';

const Root = styled("div", {
  base: {
    minHeight: "100vh",
    background: theme.color.background,
    color: theme.color.foreground,
  },
});

function TheRoot() {
  return (
    <Root class={themeModeClass()}>
      <Outlet />
    </Root>
  )
}

function App() {
  return (
    <Routes>
      <Route path="/" component={TheRoot}>
        <Pages />
        <Debug />
      </Route>
    </Routes>
  )
}

export default App
