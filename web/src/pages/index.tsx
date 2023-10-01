import { Component } from 'solid-js'
import { styled, } from '@macaron-css/solid'
import { usePlayerPauseMutation, usePlayerPlayMutation, usePlayerVolumeMutation, } from '~/hooks/api'
import { Player } from '~/components/Player'
import { Outlet, Route } from '@solidjs/router'
import { Home } from './Home'
import { mixin, theme } from '~/ui/theme'
import { AUTO_MODE, DARK_MODE, LIGHT_MODE, themeMode, themeModeClass, toggleThemeMode } from '~/ui/theme-mode'
import { PlayerStatesProvider, usePlayerStates } from '~/providers/playerStates'
import { CurrentPlayerProvider, useCurrentPlayer } from '~/providers/currentPlayer'
import { Button } from '~/ui/Button'
import { ThemeIcon } from '~/ui/ThemeIcon'
import { style } from '@macaron-css/core'
import { RiSystemMenuLine } from 'solid-icons/ri'

export const Pages: Component = () => {
  return (
    <Route path="/" component={App}>
      <Route path="/" component={Home} />
    </Route>
  )
}

const Root = styled("div", {
  base: {
    display: "flex",
    flexDirection: "column",
    inset: "0",
    position: "fixed",
    background: theme.color.background,
    color: theme.color.foreground,
  },
});

const Content = styled("div", {
  base: {
    flex: "1",
    overflowY: "auto"
  },
});

function App() {
  const themeTitle = () => {
    switch (themeMode()) {
      case AUTO_MODE:
        return "Theme Auto"
      case LIGHT_MODE:
        return "Theme Light"
      case DARK_MODE:
        return "Theme Dark"
      default:
        return "Theme Unknown"
    }
  }

  return (
    <PlayerStatesProvider>
      <CurrentPlayerProvider>
        <Root class={themeModeClass()}>
          <Header>
            <HeaderContent>
              <Button size='icon' variant='ghost'>
                <RiSystemMenuLine class={style({ ...mixin.size("8") })} />
              </Button>
              <Button size='icon' variant='ghost' onClick={toggleThemeMode} title={themeTitle()}>
                <ThemeIcon class={style({ ...mixin.size("8") })} />
              </Button>
            </HeaderContent>
          </Header>
          <Content>
            <Outlet />
          </Content>
          <div>
            <ThePlayer />
          </div>
        </Root>
      </CurrentPlayerProvider>
    </PlayerStatesProvider>
  )
}

const Header = styled("div", {
  base: {
    height: theme.space[14],
    width: "100%",
    background: theme.color.card,
    borderBottom: `${theme.space.px} solid ${theme.color.border}`,
  },
});

const HeaderContent = styled("div", {
  base: {
    ...mixin.row("4"),
    justifyContent: "space-between",
    alignItems: "center",
    height: "100%",
    paddingLeft: theme.space[4],
    paddingRight: theme.space[4],
  },
});

function ThePlayer() {
  // Queries
  const { playerStates } = usePlayerStates()
  const { currentPlayerId, setCurrentPlayerId, currentPlayerState } = useCurrentPlayer()

  // Mutations
  const playerVolumeMutation = usePlayerVolumeMutation()
  const playerPlayMutation = usePlayerPlayMutation()
  const playerPauseMutation = usePlayerPauseMutation()

  return (
    <Player
      player={currentPlayerState()}
      players={playerStates}
      onPlayClick={() => currentPlayerState()?.playing ? playerPauseMutation.mutate(currentPlayerId()) : playerPlayMutation.mutate(currentPlayerId())}
      playLoading={playerPauseMutation.isLoading || playerPlayMutation.isLoading}
      onVolumeDownClick={() => playerVolumeMutation.mutate({ id: currentPlayerId(), delta: -1 })}
      onVolumeUpClick={() => playerVolumeMutation.mutate({ id: currentPlayerId(), delta: 1 })}
      onVolumeClick={() => playerVolumeMutation.mutate({ id: currentPlayerId(), delta: 0, mute: !currentPlayerState()?.muted })}
      onPlayerClick={(id) => setCurrentPlayerId((prev) => prev == id ? 0 : id)}
    />
  )
}
