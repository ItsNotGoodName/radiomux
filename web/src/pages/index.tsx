import { Component, ComponentProps, createEffect, createSignal } from 'solid-js'
import { styled, } from '@macaron-css/solid'
import { ConnectionIndicator, Player } from '~/components/Player'
import { Link, Outlet, Route } from '@solidjs/router'
import { Home } from './Home'
import { minScreen, mixin, theme, tw } from '~/ui/theme'
import { AUTO_MODE, DARK_MODE, LIGHT_MODE, themeMode, toggleThemeMode } from '~/ui/theme-mode'
import { PlayerStatesProvider, WebSocketState, usePlayerStates } from '~/providers/playerStates'
import { CurrentPlayerProvider, useCurrentPlayer } from '~/providers/currentPlayer'
import { Button } from '~/ui/Button'
import { ThemeIcon } from '~/ui/ThemeIcon'
import { style } from '@macaron-css/core'
import { RiSystemMenuLine } from 'solid-icons/ri'
import { Players } from './Players'
import { Presets } from './Presets'
import { StateAction } from '~/api/client.gen'
import { useStateActionSetMutation, useStateVolumeSetMutation } from '~/hooks/api'
import { As, DropdownMenu } from '@kobalte/core'
import { DropdownMenuContent } from '~/ui/DropdownMenu'
import { ToastList, ToastRegion } from '~/ui/Toast'
import { Portal } from 'solid-js/web'
import { toastWebrpcError } from '~/common/error'

const Header = styled("div", {
  base: {
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

const MenuLink = styled(Link, {
  base: {
    textDecoration: "none",
    padding: theme.space[2],
    borderRadius: theme.borderRadius.ok,
  }
})

const menuLinkInactiveClass = style({
  color: theme.color.foreground,
  ":hover": {
    background: theme.color.primary,
    color: theme.color.primaryForeground,
  },
})

const menuLinkActiveClass = style({
  background: theme.color.primary,
  color: theme.color.primaryForeground,
})

const MenuIcon = styled(RiSystemMenuLine, {
  base: {
    ...mixin.size("6")
  }
})

function HeaderContent(props: { onMenuClick?: () => void }) {
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
      <DropdownMenu.Root>
        <DropdownMenu.Trigger asChild>
          <As component={Button} size='icon' variant='ghost' title="Menu" class={style({
            "@media": {
              [minScreen.md]: {
                display: "none"
              },
            },
          })}>
            <MenuIcon />
          </As>
        </DropdownMenu.Trigger>
        <DropdownMenu.Portal>
          <DropdownMenuContent>
            <DropdownMenu.Arrow />
            <DropdownMenu.Item asChild>
              <As component={MenuLink} activeClass={menuLinkActiveClass} inactiveClass={menuLinkInactiveClass} href="/" end>Home</As>
            </DropdownMenu.Item>
            <DropdownMenu.Item asChild>
              <As component={MenuLink} activeClass={menuLinkActiveClass} inactiveClass={menuLinkInactiveClass} href="/players">Players</As>
            </DropdownMenu.Item>
            <DropdownMenu.Item asChild>
              <As component={MenuLink} activeClass={menuLinkActiveClass} inactiveClass={menuLinkInactiveClass} href="/presets">Presets</As>
            </DropdownMenu.Item>
          </DropdownMenuContent>
        </DropdownMenu.Portal>
      </DropdownMenu.Root>
      <Button onClick={props.onMenuClick} size='icon' variant='ghost' title="Menu" class={style({
        display: "none",
        "@media": {
          [minScreen.md]: {
            display: "flex"
          },
        }
      })}>
        <MenuIcon />
      </Button>
      <div class={style({
        ...tw.textXl,
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
    </div >
  )
}

const Root = styled("div", {
  base: {
    display: "flex",
    minHeight: "100vh",
  },
});

const TheMenuRoot = styled("aside", {
  base: {
    background: theme.color.nav,
    transition: "border-width 250ms",
    borderRight: `0px solid ${theme.color.border}`,
    "@media": {
      [minScreen.md]: {
        selectors: {
          [`&[data-open]`]: {
            borderRight: `1px solid ${theme.color.border}`,
          }
        }
      }
    }
  }
})

function MenuRoot(props: Omit<ComponentProps<typeof TheMenuRoot>, "ref"> & { menuOpen?: boolean }) {
  let ref: HTMLDivElement
  createEffect(() => {
    if (props.menuOpen) {
      ref.dataset.open = ""
    } else {
      delete ref.dataset.open
    }
  })
  return <TheMenuRoot ref={ref!} {...props} />
}

const MenuContent = styled("div", {
  base: {
    transition: "width 250ms",
    overflowX: "hidden",
    width: theme.space[0],
    "@media": {
      [minScreen.md]: {
        selectors: {
          [`${TheMenuRoot}[data-open] &`]: {
            width: theme.space[48],
          }
        }
      }
    }
  }
})

const MenuLinks = styled("div", {
  base: {
    ...mixin.stack("1"),
    padding: theme.space[2]
  }
})

const Content = styled("div", {
  base: {
    flex: "1",
    overflowX: "hidden",
  }
})

const Footer = styled("div", {
  base: {
    bottom: "0",
    position: "sticky",
  }
})

function FooterPlayer() {
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
      onPlayClick={() => stateActionSetMutation.mutateAsync({ id: currentPlayerId(), action: currentPlayerState()?.playing ? StateAction.PUASE : StateAction.PLAY }).catch(toastWebrpcError)}
      playDisabled={stateActionSetMutation.isLoading}
      onVolumeDownClick={() => stateVolumeSetMutation.mutateAsync({ id: currentPlayerId(), delta: -1 }).catch(toastWebrpcError)}
      onVolumeUpClick={() => stateVolumeSetMutation.mutateAsync({ id: currentPlayerId(), delta: 1 }).catch(toastWebrpcError)}
      onVolumeClick={() => stateVolumeSetMutation.mutateAsync({ id: currentPlayerId(), mute: !currentPlayerState()?.muted }).catch(toastWebrpcError)}
      onPlayerClick={(id) => setCurrentPlayerId((prev) => prev == id ? 0 : id)}
      onSeekClick={() => stateActionSetMutation.mutateAsync({ id: currentPlayerId(), action: StateAction.SEEK }).catch(toastWebrpcError)}
      seekDisabled={stateActionSetMutation.isLoading}
    />
  )
}


function App() {
  const [menuOpen, setMenuOpen] = createSignal(true)

  return (
    <PlayerStatesProvider>
      <CurrentPlayerProvider>
        <Portal>
          <ToastRegion>
            <ToastList class={style({ top: theme.space[14] })} />
          </ToastRegion>
        </Portal>
        <Header>
          <HeaderContent onMenuClick={() => setMenuOpen((prev) => !prev)} />
        </Header>
        <Root>
          <MenuRoot menuOpen={menuOpen()}>
            <MenuContent>
              <MenuLinks>
                <MenuLink activeClass={menuLinkActiveClass} inactiveClass={menuLinkInactiveClass} href='/' end>Home</MenuLink>
                <MenuLink activeClass={menuLinkActiveClass} inactiveClass={menuLinkInactiveClass} href='/players'>Players</MenuLink>
                <MenuLink activeClass={menuLinkActiveClass} inactiveClass={menuLinkInactiveClass} href='/presets'>Presets</MenuLink>
              </MenuLinks>
            </MenuContent>
          </MenuRoot>
          <Content>
            <Outlet />
          </Content>
        </Root>
        <Footer>
          <FooterPlayer />
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

