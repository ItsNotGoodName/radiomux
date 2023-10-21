import { Routes } from '@solidjs/router'
import { Debug } from './pages/debug';
import { Pages } from './pages';
import { useTheme } from './ui/theme-mode';

function App() {
  useTheme()

  return (
    <Routes>
      <Pages />
      <Debug />
    </Routes>
  )
}

export default App
