import { Route, Routes } from '@solidjs/router'
import { Debug } from './pages/debug';
import { Pages } from './pages';

function App() {
  return (
    <Routes>
      <Route path="/" >
        <Pages />
      </Route>
      <Route path="/debug"  >
        <Debug />
      </Route>
    </Routes>
  )
}

export default App
