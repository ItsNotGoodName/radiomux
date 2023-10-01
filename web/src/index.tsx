/* @refresh reload */
import { render } from 'solid-js/web'
import { QueryClient, QueryClientProvider } from '@tanstack/solid-query'
import { Router } from "@solidjs/router";

import "modern-normalize/modern-normalize.css";

import App from './App'

const root = document.getElementById('root')

const queryClient = new QueryClient()

render(() => (
  <QueryClientProvider client={queryClient}>
    <Router>
      <App />
    </Router>
  </QueryClientProvider>
), root!)
