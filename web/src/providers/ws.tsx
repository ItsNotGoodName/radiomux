import { Accessor, Show, batch, onCleanup } from 'solid-js'
import { createReconnectingWS, createWSState } from '@solid-primitives/websocket';
import { createStore } from 'solid-js/store';
import { ApiEvent, ApiPlayerState } from '~/api';
import {
  createContext,
  ParentComponent,
  useContext
} from "solid-js";
import { relativeWsUrl } from '~/common';
import { usePlayerListQuery } from '~/hooks/api';
import { ToastTitle, ToastContent, ToastDescription, toast, ToastCloseButton } from '~/ui/Toast';

type WSContextType = {
  webSocketState: Accessor<WebSocketState>
  playerStates: Array<ApiPlayerState>,
};

const WSContext = createContext<WSContextType>();

type WSContextProps = {};

export enum WebSocketState {
  Connecting,
  Connected,
  Disconnecting,
  Disconnected,
}

export const WSProvider: ParentComponent<WSContextProps> = (props) => {
  const ws = createReconnectingWS(relativeWsUrl("/api/ws"));

  const webSocketState = createWSState(ws);

  const [playerStates, setPlayerStates] = createStore<Array<ApiPlayerState>>([])

  usePlayerListQuery({
    onSuccess: (data) => {
      batch(() =>
        data.players.map((player) =>
          setPlayerStates(state => state.id == player.id, (state) =>
            ({ ...state, name: player.name }))))
    }
  })

  const onMessage = (msg: MessageEvent<string>) => {
    const event = JSON.parse(msg.data) as ApiEvent

    switch (event.type) {
      case "player_state_partial":
        // TODO: figure out if reconcile can be used here
        batch(() => {
          for (const s of event.data) {
            for (const [k, v] of Object.entries(s)) {
              setPlayerStates(state => state.id == s.id, k as any, v)
            }
          }
        })
        break
      case "player_state":
        setPlayerStates(event.data)
        break
      case "notification":
        toast.custom(() => (
          <ToastContent>
            <ToastCloseButton />
            <Show when={event.data.title}>
              {title => <ToastTitle>{title()}</ToastTitle>}
            </Show>
            <Show when={event.data.description}>
              {description => <ToastDescription>{description()}</ToastDescription>}
            </Show>
          </ToastContent>
        ), { variant: event.data.error ? "destructive" : "default" })
    }
  }

  let errorToastId: number | undefined
  const onError = (event: Event) => {
    console.log(event)
    errorToastId && toast.dismiss(errorToastId)
    errorToastId = toast.error("WebSocket", "Disconnected from server due to an error.")
  }

  ws.addEventListener("message", onMessage)
  ws.addEventListener("error", onError)
  onCleanup(() => {
    ws.removeEventListener("message", onMessage)
    ws.removeEventListener("error", onError)
  })

  const store: WSContextType = {
    playerStates,
    webSocketState,
  };

  return (
    <WSContext.Provider value={store}>
      {props.children}
    </WSContext.Provider>)
};

export function useWS(): WSContextType {
  const result = useContext(WSContext);
  if (!result) throw new Error("useWS must be used within a WSProvider");
  return result;
}

