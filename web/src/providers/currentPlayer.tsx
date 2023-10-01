import { Accessor, Setter, createEffect, createMemo, createSignal, on, } from 'solid-js'
import { ApiPlayerState } from '~/api';
import {
  createContext,
  ParentComponent,
  useContext
} from "solid-js";
import { usePlayerStates } from './playerStates';

type CurrentPlayerContextType = {
  currentPlayerId: Accessor<number>
  setCurrentPlayerId: Setter<number>
  currentPlayerState: Accessor<ApiPlayerState | undefined>,
};

const CurrentPlayerContext = createContext<CurrentPlayerContextType>();

type CurrentPlayerContextProps = {};

export const CurrentPlayerProvider: ParentComponent<CurrentPlayerContextProps> = (props) => {
  // Storage
  const [currentPlayerId, setCurrentPlayerId] = createSignal<number>(parseInt(localStorage.getItem("currentPlayerId") || "") || 0);
  createEffect(
    on(
      currentPlayerId,
      () => {
        localStorage.setItem("currentPlayerId", currentPlayerId().toString());
      },
      { defer: true }
    )
  );

  // Queries
  const { playerStates } = usePlayerStates()
  const currentPlayerState = createMemo(() => playerStates.find((p) => p.id == currentPlayerId()))

  const store: CurrentPlayerContextType = {
    currentPlayerId,
    setCurrentPlayerId,
    currentPlayerState,
  };

  return (
    <CurrentPlayerContext.Provider value={store}>
      {props.children}
    </CurrentPlayerContext.Provider>)
};

export function useCurrentPlayer(): CurrentPlayerContextType {
  const result = useContext(CurrentPlayerContext);
  if (!result) throw new Error("useCurrentPlayer must be used within a CurrentPlayerProvider");
  return result;
}

