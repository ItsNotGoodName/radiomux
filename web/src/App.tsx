import { Outlet, Route, Routes } from '@solidjs/router'
import { Debug } from './pages/debug';
import { Pages } from './pages';
import { useTheme } from './ui/theme-mode';
import { theme } from './ui/theme';
import { styled } from '@macaron-css/solid';

const TheRoot = styled("div", {
  base: {
    inset: 0,
    position: "fixed",
    overflowY: "auto",
    background: theme.color.background,
    color: theme.color.foreground,
  },
});

function Root() {
  return (
    <TheRoot>
      <Outlet />
    </TheRoot>
  )
}

function App() {
  useTheme()

  return (
    <Routes>
      <Route path="/" component={Root}>
        <Pages />
        <Debug />
      </Route>
    </Routes>
  )
}

export default App
