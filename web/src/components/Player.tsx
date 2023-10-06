import { style } from "@macaron-css/core";
import { styled } from "@macaron-css/solid";
import { RiDeviceDeviceFill, RiMediaPauseFill, RiMediaPlayFill, RiMediaVolumeDownFill, RiMediaVolumeMuteFill, RiMediaVolumeUpFill, RiSystemRefreshFill } from "solid-icons/ri";
import { For, Show, } from "solid-js";
import { Button } from "~/ui/Button";
import { minScreen, mixin, theme, tw } from "~/ui/theme";
import { Dropdown, DropdownCard, DropdownPositioner } from "~/ui/Dropdown";

const Root = styled("div", {
  base: {
    background: theme.color.nav,
    color: theme.color.navForeground,
    height: theme.space[28],
    borderTop: `1px solid ${theme.color.border}`,
    "@media": {
      [minScreen.md]: {
        height: theme.space[14],
      },
    },
  }
})

const Content = styled("div", {
  base: {
    display: "flex",
    "height": "100%",
    paddingLeft: theme.space[2],
    paddingRight: theme.space[2],
    flexDirection: "column",
    "@media": {
      [minScreen.md]: {
        flexDirection: "row",
        gap: theme.space[2]
      },
    },
  }
})

const ContentChild = styled("div", {
  base: {
    ...mixin.row("2"),
    alignItems: "center",
    flex: 1,
    overflow: "hidden",
    ":last-child": {
      borderTop: `1px solid ${theme.color.border}`
    },
    "@media": {
      [minScreen.md]: {
        ":last-child": {
          borderTop: `unset`
        },
      },
    },
  }
})

const MainControl = styled("div", {
  base: {
  }
})

const SubControl = styled("div", {
  base: {
    ...mixin.row("2"),
    flex: "1",
    alignItems: "center",
    overflowX: "hidden",
  }
})

const Text = styled("p", {
  base: {
    ...mixin.textLine()
  }
})

const Thumbnail = styled("img", {
  base: {
    ...mixin.size("10"),
    background: theme.color.muted
  }
})

export const ConnectionIndicator = styled("div", {
  base: {
    ...mixin.size("4"),
    borderRadius: theme.borderRadius.full,
  },
  variants: {
    connected: {
      true: {
        background: "lime"
      },
    },
    disconnected: {
      true: {
        background: "red"
      }
    },
    connecting: {
      true: {
        background: "orange"
      }
    }
  }
})

const PlayerGroup = styled("div", {
  base: {
    ...mixin.row("2"),
    overflowX: "hidden",
    alignItems: "center",
    flex: "1"
  }
})

const VolumeGroup = styled("div", {
  base: {
    ...mixin.row("2"),
    alignItems: "center"
  }
})

const IconRiDeviceDeviceFill = styled(RiDeviceDeviceFill, {
  base: {
    ...mixin.size("6")
  },
  variants: {
    selected: {
      true: {
        color: theme.color.primary
      }
    }
  }
})

const IconRiSystemRefreshFill = styled(RiSystemRefreshFill, {
  base: {
    ...mixin.size("6")
  },
  variants: {
    loading: {
      true: {
        ...tw.animateSpin,
      }
    }
  }
})

const PlayerTableHead = styled("th", {
  base: {
    textAlign: "left"
  }
})

type Props = {
  player?: {
    id: number
    name: string
    ready: boolean
    connected: boolean
    playing: boolean
    loading: boolean,
    volume: number
    muted: boolean
    title: string
    genre: string
    station: string
    uri: string
  }
  players: Array<{
    id: number
    name: string
    connected: boolean
  }>
  onPlayClick?: () => void
  playDisabled?: boolean,
  onVolumeUpClick?: () => void
  onVolumeDownClick?: () => void
  onVolumeClick?: () => void
  onPlayerClick: (id: number) => void
  onSeekClick?: () => void
  seekDisabled?: boolean
}


export function Player(props: Props) {
  const playerDisabled = () => props.player == undefined || !props.player.connected || !props.player.ready

  const seekDisabled = () => playerDisabled() || props.seekDisabled

  const playDisabled = () => playerDisabled() || props.playDisabled || props.seekDisabled
  const playStatus = () => {
    if (props.player == undefined) {
      return "Unknown"
    }
    return props.player?.playing ? "Playing" : "Paused"
  }

  return (
    <Root>
      <Content>
        <ContentChild>
          <div class={style({ display: "flex", alignItems: "center" })}>
            <Dropdown options={{ placement: "top" }} button={
              ref =>
                <Button disabled={playerDisabled()} title="Media" ref={ref} size="icon" variant="ghost" class={style({
                  color: theme.color.cardForeground,
                  borderRadius: theme.borderRadius.full,
                  overflow: "hidden"
                })}>
                  <Thumbnail src="/favicon.svg" alt="Media Thumbnail" />
                </Button>
            }>
              {ref =>
                <DropdownPositioner ref={ref} class={style({ maxWidth: theme.size.md, width: "100%" })}>
                  <DropdownCard class={style({ padding: theme.space[2], overflowX: "auto" })}>
                    <Show when={props.player}>
                      {player =>
                        <table>
                          <tbody>
                            <tr>
                              <PlayerTableHead>Title</PlayerTableHead>
                              <td>{player().title}</td>
                            </tr>
                            <tr>
                              <PlayerTableHead>Genre</PlayerTableHead>
                              <td>{player().genre}</td>
                            </tr>
                            <tr>
                              <PlayerTableHead>Station</PlayerTableHead>
                              <td>{player().station}</td>
                            </tr>
                            <tr>
                              <PlayerTableHead>URI</PlayerTableHead>
                              <td>{player().uri}</td>
                            </tr>
                          </tbody>
                        </table>
                      }
                    </Show>
                  </DropdownCard>
                </DropdownPositioner>}
            </Dropdown>
          </div>
          <Text title={props.player?.title}>{props.player?.title}</Text>
          <Button size="icon" variant="ghost" title="Seek" disabled={seekDisabled()} onClick={props.onSeekClick} class={style({ marginLeft: "auto" })}>
            <IconRiSystemRefreshFill loading={props.player?.loading || props.seekDisabled} />
          </Button>
        </ContentChild>
        <ContentChild>
          <MainControl>
            <Button
              size="icon"
              variant="ghost"
              disabled={playDisabled()}
              onClick={props.onPlayClick}
              title={playStatus()}
            >
              <Show when={props.player?.playing} fallback={
                <RiMediaPlayFill class={style({ ...mixin.size("10") })} />
              }>
                <RiMediaPauseFill class={style({ ...mixin.size("10") })} />
              </Show>
            </Button>
          </MainControl>
          <SubControl>
            <Show when={props.player} fallback={<PlayerGroup />}>
              {player => (
                <>
                  <PlayerGroup>
                    <div>
                      <ConnectionIndicator connected={player().connected} disconnected={!player().connected && !player().ready} connecting={!player().ready && player().connected} />
                    </div>
                    <Text title={player().name}>{player().name}</Text>
                  </PlayerGroup>
                  <VolumeGroup>
                    <Button disabled={playerDisabled()} size="icon" variant="ghost" onClick={props.onVolumeDownClick} title="Volume Down">
                      <RiMediaVolumeDownFill class={style({ ...mixin.size("6") })} />
                    </Button>
                    <Button disabled={playerDisabled()} size="icon" variant="ghost" onClick={props.onVolumeClick} title={player().muted ? "Volume Muted" : undefined}>
                      <Show when={!player().muted} fallback={
                        <RiMediaVolumeMuteFill class={style({ ...mixin.size("6"), color: "red" })} />
                      }>
                        {player().volume}
                      </Show>
                    </Button>
                    <Button disabled={playerDisabled()} size="icon" variant="ghost" onClick={props.onVolumeUpClick} title="Volume Up">
                      <RiMediaVolumeUpFill class={style({ ...mixin.size("6") })} />
                    </Button>
                  </VolumeGroup>
                </>
              )}
            </Show>
            <Dropdown button={ref =>
              <Button ref={ref} size="icon" variant="ghost" title="Players" >
                <IconRiDeviceDeviceFill selected={!!props.player?.id} />
              </Button>
            }>
              {ref =>
                <DropdownPositioner ref={ref} class={style({ maxWidth: theme.space[60], width: "100%" })}>
                  <DropdownCard class={style({ padding: theme.space[2] })}>
                    <For each={props.players}>
                      {player => (
                        <Button
                          onClick={[props.onPlayerClick, player.id]}
                          size="sm"
                          variant={player.id == props.player?.id ? "default" : "ghost"}
                          class={style({ ...mixin.row("2",), width: "100%", justifyContent: "start" })}
                        >
                          <div>
                            <ConnectionIndicator connected={player.connected} disconnected={!player.connected} />
                          </div>
                          <Text title={player.name}>{player.name}</Text>
                        </Button>
                      )}
                    </For>
                  </DropdownCard>
                </DropdownPositioner>
              }
            </Dropdown>
          </SubControl>
        </ContentChild>
      </Content>
    </Root>
  )
}
