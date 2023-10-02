import { Component } from 'solid-js'
import { styled, } from '@macaron-css/solid'
import { usePlayerPauseMutation, usePlayerPlayMutation, usePlayerVolumeMutation, } from '~/hooks/api'
import { ConnectionIndicator, Player } from '~/components/Player'
import { Outlet, Route } from '@solidjs/router'
import { Home } from './Home'
import { mixin, theme } from '~/ui/theme'
import { AUTO_MODE, DARK_MODE, LIGHT_MODE, themeMode, themeModeClass, toggleThemeMode } from '~/ui/theme-mode'
import { PlayerStatesProvider, WebSocketState, usePlayerStates } from '~/providers/playerStates'
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
    minHeight: "100vh",
    display: "flex",
    flexDirection: "column",
    background: theme.color.background,
    color: theme.color.foreground,
  },
});

const Header = styled("div", {
  base: {
    height: theme.space[14],
    position: "sticky",
    top: "0",
    width: "100%",
    background: theme.color.card,
    borderBottom: `${theme.space.px} solid ${theme.color.border}`,
    zIndex: "10"
  },
});

const HeaderContent = styled("div", {
  base: {
    ...mixin.row("2"),
    alignItems: "center",
    height: "100%",
    paddingLeft: theme.space[2],
    paddingRight: theme.space[2],
  },
});

const Content = styled("div", {
  base: {
    flex: "1",
    paddingBottom: theme.space[28]
  },
});

const Footer = styled("div", {
  base: {
    bottom: "0",
    left: "0",
    right: "0",
    position: "fixed",
  }
})

function TheHeaderContent() {
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

  const { webSocketState } = usePlayerStates()

  const wsStates = ["Connecting", "Connected", "Disconnecting", "Disconnected"];

  return (
    <>
      <Button size='icon' variant='ghost' title="Menu">
        <RiSystemMenuLine class={style({ ...mixin.size("8") })} />
      </Button>
      <div class={style({ ...mixin.textLine(), display: "flex", flex: "1", alignItems: "center" })}>
        RadioMux
      </div>
      <div class={style({ ...mixin.row("2"), alignItems: "center" })}>
        <ConnectionIndicator
          title={"WebSocket " + wsStates[webSocketState()]}
          connected={webSocketState() == WebSocketState.Connected}
          disconnected={webSocketState() == WebSocketState.Disconnected || webSocketState() == WebSocketState.Disconnecting}
          connecting={webSocketState() == WebSocketState.Connecting}
        />
        <Button size='icon' variant='ghost' onClick={toggleThemeMode} title={themeTitle()}>
          <ThemeIcon class={style({ ...mixin.size("8") })} />
        </Button>
      </div>
    </>
  )

}

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

function App() {
  return (
    <PlayerStatesProvider>
      <CurrentPlayerProvider>
        <Root class={themeModeClass()}>
          <Header>
            <HeaderContent>
              <TheHeaderContent />
            </HeaderContent>
          </Header>
          <Content>
            <Outlet />
          </Content>
          <Footer>
            <ThePlayer />
          </Footer>
        </Root>
      </CurrentPlayerProvider>
    </PlayerStatesProvider>
  )
}

