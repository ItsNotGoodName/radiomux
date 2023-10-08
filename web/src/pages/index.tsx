import { Component, } from 'solid-js'
import { styled, } from '@macaron-css/solid'
import { ConnectionIndicator, Player } from '~/components/Player'
import { Link, Outlet, Route } from '@solidjs/router'
import { Home } from './Home'
import { mixin, theme, tw } from '~/ui/theme'
import { AUTO_MODE, DARK_MODE, LIGHT_MODE, themeMode, toggleThemeMode } from '~/ui/theme-mode'
import { PlayerStatesProvider, WebSocketState, usePlayerStates } from '~/providers/playerStates'
import { CurrentPlayerProvider, useCurrentPlayer } from '~/providers/currentPlayer'
import { Button } from '~/ui/Button'
import { ThemeIcon } from '~/ui/ThemeIcon'
import { style } from '@macaron-css/core'
import { RiSystemMenuLine } from 'solid-icons/ri'
import { Dropdown, DropdownCard, DropdownPositioner } from "~/ui/Dropdown";
import { Players } from './Players'
import { Presets } from './Presets'
import { StateAction } from '~/api/client.gen'
import { useStateActionSetMutation, useStateVolumeSetMutation } from '~/hooks/api'

const Header = styled("div", {
  base: {
    ...tw.textXl,
    height: theme.space[14],
    position: "sticky",
    top: "0",
    width: "100%",
    background: theme.color.nav,
    color: theme.color.navForeground,
    borderBottom: `${theme.space.px} solid ${theme.color.border}`,
    zIndex: "10"
  },
});

const menuLinkInactiveClass = style({
  textDecoration: "none",
  padding: theme.space[2],
  borderRadius: theme.borderRadius.ok,
  ":hover": {
    background: theme.color.accent,
    color: theme.color.accentForeground,
  },

  color: theme.color.foreground,
})

const menuLinkActiveClass = style({
  textDecoration: "none",
  padding: theme.space[2],
  borderRadius: theme.borderRadius.ok,

  background: theme.color.accent,
  color: theme.color.accentForeground,
})

function TheHeader() {
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
    <div class={style({
      ...mixin.row("2"),
      alignItems: "center",
      height: "100%",
      paddingLeft: theme.space[2],
      paddingRight: theme.space[2],
    })}>
      <Dropdown button={ref =>
        <Button ref={ref} size='icon' variant='ghost' title="Menu">
          <RiSystemMenuLine class={style({ ...mixin.size("6") })} />
        </Button>
      }>
        {(ref, setOpen) =>
          <DropdownPositioner ref={ref} class={style({ maxWidth: theme.space[60], width: "100%" })}>
            <DropdownCard class={style({
              display: "flex",
              flexDirection: "column",
              padding: theme.space[2],
            })}>
              <Link onClick={[setOpen, false]} activeClass={menuLinkActiveClass} inactiveClass={menuLinkInactiveClass} href="/" end>Home</Link>
              <Link onClick={[setOpen, false]} activeClass={menuLinkActiveClass} inactiveClass={menuLinkInactiveClass} href="/players">Players</Link>
              <Link onClick={[setOpen, false]} activeClass={menuLinkActiveClass} inactiveClass={menuLinkInactiveClass} href="/presets">Presets</Link>
            </DropdownCard>
          </DropdownPositioner>
        }
      </Dropdown>
      <div class={style({
        ...mixin.textLine(),
        display: "flex",
        flex: "1",
        alignItems: "center"
      })}>
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
          <ThemeIcon class={style({ ...mixin.size("6") })} />
        </Button>
      </div>
    </div>
  )
}

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

function ThePlayer() {
  // Queries
  const { playerStates } = usePlayerStates()
  const { currentPlayerId, setCurrentPlayerId, currentPlayerState } = useCurrentPlayer()

  // Mutations
  const stateVolumeSetMutation = useStateVolumeSetMutation()
  const stateActionSetMutation = useStateActionSetMutation()

  return (
    <Player
      player={currentPlayerState()}
      players={playerStates}
      onPlayClick={() => stateActionSetMutation.mutate({ id: currentPlayerId(), action: currentPlayerState()?.playing ? StateAction.PUASE : StateAction.PLAY })}
      playDisabled={stateActionSetMutation.isLoading}
      onVolumeDownClick={() => stateVolumeSetMutation.mutate({ id: currentPlayerId(), delta: -1 })}
      onVolumeUpClick={() => stateVolumeSetMutation.mutate({ id: currentPlayerId(), delta: 1 })}
      onVolumeClick={() => stateVolumeSetMutation.mutate({ id: currentPlayerId(), mute: !currentPlayerState()?.muted })}
      onPlayerClick={(id) => setCurrentPlayerId((prev) => prev == id ? 0 : id)}
      onSeekClick={() => stateActionSetMutation.mutate({ id: currentPlayerId(), action: StateAction.SEEK })}
      seekDisabled={stateActionSetMutation.isLoading}
    />
  )
}

function App() {
  return (
    <PlayerStatesProvider>
      <CurrentPlayerProvider>
        <Header>
          <TheHeader />
        </Header>
        <Content>
          <Outlet />
        </Content>
        <Footer>
          <ThePlayer />
        </Footer>
      </CurrentPlayerProvider>
    </PlayerStatesProvider>
  )
}

export const Pages: Component = () => {
  return (
    <Route path="/" component={App}>
      <Route path="/" component={Home} />
      <Route path="/players" component={Players} />
      <Route path="/presets" component={Presets} />
    </Route>
  )
}

