import { Accessor, batch, } from 'solid-js'
import { createReconnectingWS, createWSState } from '@solid-primitives/websocket';
import { createStore } from 'solid-js/store';
import { ApiPlayerState, ApiEvent } from '~/api';
import {
  createContext,
  ParentComponent,
  useContext
} from "solid-js";
import { relativeWsUrl } from '~/common';

type PlayerStatesContextType = {
  connection: Accessor<string>
  playerStates: Array<ApiPlayerState>,
};

const PlayerStatesContext = createContext<PlayerStatesContextType>();

type PlayerStatesContextProps = {};

export const PlayerStatesProvider: ParentComponent<PlayerStatesContextProps> = (props) => {
  const ws = createReconnectingWS(relativeWsUrl("/api/ws"));

  const wsStates = ["Connecting", "Connected", "Disconnecting", "Disconnected"];
  const wsState = createWSState(ws);
  const connection = () => wsStates[wsState()]

  const [playerStates, setPlayerStates] = createStore<Array<ApiPlayerState>>([])
  ws.addEventListener("message", (msg) => {
    const event = JSON.parse(msg.data) as ApiEvent

    switch (event.type) {
      case "player_state":
        setPlayerStates(event.data)
        break
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
    }
  })

  const store: PlayerStatesContextType = {
    playerStates,
    connection
  };

  return (
    <PlayerStatesContext.Provider value={store}>
      {props.children}
    </PlayerStatesContext.Provider>)
};

export function usePlayerStates(): PlayerStatesContextType {
  const result = useContext(PlayerStatesContext);
  if (!result) throw new Error("usePlayerStates must be used within a PlayerStatesProvider");
  return result;
}

