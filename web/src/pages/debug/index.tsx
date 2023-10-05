import { Route } from '@solidjs/router'
import { Ui } from './Ui';
import { Home } from './Home';

export function Debug() {
  return (
    <Route path="/debug">
      <Route path="/*" component={Home} />
      <Route path="/ui" component={Ui} />
    </Route>
  )
}

