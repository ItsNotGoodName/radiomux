import { style } from "@macaron-css/core";
import { styled } from "@macaron-css/solid";
import { RiDeviceDeviceFill, RiMediaPauseFill, RiMediaPlayFill, RiMediaVolumeDownFill, RiMediaVolumeMuteFill, RiMediaVolumeUpFill, RiSystemRefreshFill } from "solid-icons/ri";
import { For, Show, } from "solid-js";
import { Button } from "~/ui/Button";
import { minScreen, mixin, theme, tw } from "~/ui/theme";
import { As, Popover, Progress, } from "@kobalte/core";
import { durationHumanize } from "~/common";
import { DropdownMenuArrow, DropdownMenuContent, DropdownMenuItem, DropdownMenuPortal, DropdownMenuRoot, DropdownMenuTrigger } from "~/ui/DropdownMenu";
import { PopoverContent } from "~/ui/Popover";

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

const Text = styled("div", {
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
    textAlign: "left",
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
    timeline_duration: number
    timeline_is_seekable: boolean
    timeline_is_placeholder: boolean
    playback_error: string
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
            <Popover.Root>
              <Popover.Trigger asChild>
                <As component={Button} disabled={playerDisabled()} title="Media" size="icon" variant="ghost" class={style({
                  color: theme.color.cardForeground,
                  borderRadius: theme.borderRadius.full,
                  overflow: "hidden"
                })}>
                  <Thumbnail src="/favicon.svg" alt="Media Thumbnail" />
                </As>
              </Popover.Trigger>
              <Popover.Portal>
                <PopoverContent class={style({ padding: theme.space[2], width: theme.size.sm })}>
                  <Popover.Arrow />
                  <Show when={props.player}>
                    {player =>
                      <div class={style({ overflowX: "auto", })}>
                        <table class={style({ borderCollapse: "collapse", })}>
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
                            <Show when={player().playback_error}>
                              <tr>
                                <PlayerTableHead>Error</PlayerTableHead>
                                <td>{player().playback_error}</td>
                              </tr>
                            </Show>
                          </tbody>
                        </table>
                      </div>
                    }
                  </Show>
                </PopoverContent>
              </Popover.Portal>
            </Popover.Root>
          </div>
          <div class={style({ overflowX: "hidden", flex: 1 })}>
            <Text title={props.player?.title}>{props.player?.title}</Text>
            <div class={style({ ...mixin.row("1"), alignItems: 'center' })}>
              <Show when={props.player?.timeline_is_seekable && !props.player.timeline_is_placeholder}>
                <Text>{durationHumanize(0)}</Text>
                <Progress.Root value={75} class={style({ flex: 1 })}>
                  <Progress.Track class={style({ height: theme.space[2] })}>
                    <Progress.Fill class={style({ height: "100%", width: "var(--kb-progress-fill-width)", background: theme.color.navForeground, borderRadius: theme.borderRadius.ok })} />
                  </Progress.Track>
                </Progress.Root>
                <Text>{durationHumanize(props.player?.timeline_duration ?? 0)}</Text>
              </Show>
            </div>
          </div>
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
            <DropdownMenuRoot>
              <DropdownMenuTrigger asChild>
                <As component={Button} size="icon" variant="ghost" title="Players">
                  <IconRiDeviceDeviceFill selected={!!props.player?.id} />
                </As>
              </DropdownMenuTrigger>
              <DropdownMenuPortal>
                <DropdownMenuContent>
                  <DropdownMenuArrow />
                  <For each={props.players}>
                    {player => (
                      <DropdownMenuItem asChild closeOnSelect={false}>
                        <As component={Button}
                          onClick={[props.onPlayerClick, player.id]}
                          size="sm"
                          variant={player.id == props.player?.id ? "default" : "ghost"}
                          class={style({ ...mixin.row("2",), width: "100%", justifyContent: "start" })}
                        >
                          <div>
                            <ConnectionIndicator connected={player.connected} disconnected={!player.connected} />
                          </div>
                          <Text title={player.name}>{player.name}</Text>
                        </As>
                      </DropdownMenuItem>
                    )}
                  </For>
                </DropdownMenuContent>
              </DropdownMenuPortal>
            </DropdownMenuRoot>
          </SubControl>
        </ContentChild>
      </Content>
    </Root >
  )
}
