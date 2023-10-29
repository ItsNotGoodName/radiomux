import { style } from "@macaron-css/core";
import { styled } from "@macaron-css/solid";
import { RiDeviceDeviceFill, RiMediaPauseFill, RiMediaPlayFill, RiMediaSkipBackFill, RiMediaSkipForwardFill, RiMediaStopFill, RiMediaVolumeDownFill, RiMediaVolumeMuteFill, RiMediaVolumeUpFill, RiSystemRefreshFill } from "solid-icons/ri";
import { For, Show, } from "solid-js";
import { Button } from "~/ui/Button";
import { minScreen, mixin, theme, tw } from "~/ui/theme";
import { As, Popover, Progress, } from "@kobalte/core";
import { durationHumanize } from "~/common";
import { DropdownMenuArrow, DropdownMenuContent, DropdownMenuItem, DropdownMenuPortal, DropdownMenuRoot, DropdownMenuTrigger } from "~/ui/DropdownMenu";
import { PopoverContent, PopoverPortal, PopoverRoot, PopoverTrigger } from "~/ui/Popover";

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
    ...mixin.row("2")
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
    flexShrink: 0,
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

const iconClass = style({
  ...mixin.size("6")
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

const IconRiMediaVolumeMuteFill = styled(RiMediaVolumeMuteFill, {
  base: {
    ...mixin.size("6"),
    color: "red"
  }
})

const MediaTable = styled("table", {
  base: {
    borderCollapse: "collapse",
    width: "100%",
  }
})

const MediaTableRow = styled("tr", {
  base: {
    borderBottom: `1px solid ${theme.color.border}`,
    selectors: {
      [`&:last-child`]: {
        borderBottom: "none"
      },
    }
  }
})

const MediaTableHead = styled("th", {
  base: {
    textAlign: "left",
    padding: theme.space[1],
  }
})

const MediaTableData = styled("td", {
  base: {
    padding: theme.space[1],
  }
})

const MediaControlPopoverContent = styled("div", {
  base: {
    ...mixin.stack("2"),
    padding: theme.space[2],
    width: theme.space[48]
  }
})

const MediaControlPopoverControls = styled("div", {
  base: {
    ...mixin.row("2"),
    flexWrap: "wrap",
    justifyContent: "center",
    "@media": {
      [minScreen.lg]: {
        display: "none"
      },
    }
  }
})

const MediaControls = styled("div", {
  base: {
    display: "none",
    "@media": {
      [minScreen.lg]: {
        ...mixin.row("2")
      },
    }
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
  onPlayPauseClick?: () => void
  playPauseDisabled?: boolean,
  onStopClick?: () => void
  stopDisabled?: boolean,
  onVolumeUpClick?: () => void
  onVolumeDownClick?: () => void
  onVolumeClick?: () => void
  onPlayerChange: (id: number) => void
  onSeekBackClick?: () => void
  seekBackDisabled?: boolean
  onSeekForwardClick?: () => void
  seekForwardDisabled?: boolean
}

export function Player(props: Props) {
  const playerDisabled = () => props.player == undefined || !props.player.connected || !props.player.ready

  const playPauseTitle = () => {
    if (props.player == undefined) {
      return "Unknown"
    }
    return props.player?.playing ? "Playing" : "Paused"
  }

  return (
    <Root>
      <Content>
        <ContentChild>
          <div>
            <PopoverRoot placement="top">
              <PopoverTrigger asChild>
                <As component={Button} disabled={playerDisabled()} title="Media" size="icon" variant="ghost" class={style({
                  color: theme.color.cardForeground,
                  borderRadius: theme.borderRadius.full,
                  overflow: "hidden"
                })}>
                  <Thumbnail src="/favicon.svg" alt="Media Thumbnail" />
                </As>
              </PopoverTrigger>
              <PopoverPortal>
                <PopoverContent class={style({ padding: theme.space[2], width: theme.size.sm })}>
                  <Popover.Arrow />
                  <Show when={!playerDisabled() && props.player}>
                    {player =>
                      <div class={style({ overflowX: "auto", })}>
                        <MediaTable>
                          <tbody>
                            <MediaTableRow>
                              <MediaTableHead>Title</MediaTableHead>
                              <MediaTableData>{player().title}</MediaTableData>
                            </MediaTableRow>
                            <MediaTableRow>
                              <MediaTableHead>Genre</MediaTableHead>
                              <MediaTableData>{player().genre}</MediaTableData>
                            </MediaTableRow>
                            <MediaTableRow>
                              <MediaTableHead>Station</MediaTableHead>
                              <MediaTableData>{player().station}</MediaTableData>
                            </MediaTableRow>
                            <MediaTableRow>
                              <MediaTableHead>URI</MediaTableHead>
                              <MediaTableData>{player().uri}</MediaTableData>
                            </MediaTableRow>
                            <Show when={player().playback_error}>
                              <MediaTableRow>
                                <MediaTableHead>Error</MediaTableHead>
                                <MediaTableData>{player().playback_error}</MediaTableData>
                              </MediaTableRow>
                            </Show>
                          </tbody>
                        </MediaTable>
                      </div>
                    }
                  </Show>
                </PopoverContent>
              </PopoverPortal>
            </PopoverRoot>
          </div>
          <div class={style({ flex: 1, display: "flex", flexDirection: "column", overflowX: "hidden" })}>
            <div class={style({ flex: 1, display: "flex" })}>
              <Text title={props.player?.title}>{props.player?.title}</Text>
            </div>
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
          <PopoverRoot placement="top">
            <PopoverTrigger asChild>
              <As component={Button} size="icon" variant="ghost" title="Media Controls" disabled={playerDisabled()}>
                <IconRiSystemRefreshFill loading={props.player?.loading || props.seekBackDisabled} />
              </As>
            </PopoverTrigger>
            <PopoverPortal>
              <PopoverContent>
                <Popover.Arrow />
                <MediaControlPopoverContent>
                  <div>Lorem ipsum dolor sit, amet consectetur adipisicing elit. Fugiat voluptatibus qui aliquam obcaecati distinctio consequuntur eveniet reprehenderit debitis. Distinctio nam delectus laborum exercitationem dicta ipsa rerum dolor maxime cum maiores!</div>
                  <Show when={!playerDisabled()}>
                    <MediaControlPopoverControls>
                      <Button size="icon" variant="ghost" onClick={props.onSeekBackClick} disabled={props.seekBackDisabled} title="Seek Back">
                        <RiMediaSkipBackFill class={iconClass} />
                      </Button>
                      <Button size="icon" variant="ghost" onClick={props.onStopClick} disabled={props.stopDisabled} title="Stop">
                        <RiMediaStopFill class={iconClass} />
                      </Button>
                      <Button size="icon" variant="ghost" onClick={props.onSeekForwardClick} disabled={props.seekForwardDisabled} title="Seek Forward">
                        <RiMediaSkipForwardFill class={iconClass} />
                      </Button>
                    </MediaControlPopoverControls>
                  </Show>
                </MediaControlPopoverContent>
              </PopoverContent>
            </PopoverPortal>
          </PopoverRoot>
        </ContentChild>
        <ContentChild>
          <MainControl>
            <Button
              size="icon"
              disabled={playerDisabled() || props.playPauseDisabled}
              onClick={props.onPlayPauseClick}
              title={playPauseTitle()}
              class={style({
                borderRadius: theme.borderRadius.full
              })}
            >
              <Show when={props.player?.playing} fallback={
                <RiMediaPlayFill class={iconClass} />
              }>
                <RiMediaPauseFill class={iconClass} />
              </Show>
            </Button>
            <Show when={!playerDisabled()}>
              <MediaControls>
                <Button size="icon" variant="ghost" onClick={props.onSeekBackClick} disabled={props.seekBackDisabled} title="Seek Back">
                  <RiMediaSkipBackFill class={iconClass} />
                </Button>
                <Button size="icon" variant="ghost" onClick={props.onStopClick} disabled={props.stopDisabled} title="Stop">
                  <RiMediaStopFill class={iconClass} />
                </Button>
                <Button size="icon" variant="ghost" onClick={props.onSeekForwardClick} disabled={props.seekForwardDisabled} title="Seek Forward">
                  <RiMediaSkipForwardFill class={iconClass} />
                </Button>
              </MediaControls>
            </Show>
          </MainControl>
          <SubControl>
            <PlayerGroup>
              <Show when={props.player}>
                {player => (
                  <>
                    <ConnectionIndicator connected={player().connected} disconnected={!player().connected && !player().ready} connecting={!player().ready && player().connected} />
                    <Text title={player().name}>{player().name}</Text>
                  </>
                )}
              </Show>
            </PlayerGroup>
            <Show when={props.player}>
              {player => (
                <VolumeGroup>
                  <Button disabled={playerDisabled()} size="icon" variant="ghost" onClick={props.onVolumeDownClick} title="Volume Down">
                    <RiMediaVolumeDownFill class={iconClass} />
                  </Button>
                  <Button disabled={playerDisabled()} size="icon" variant="ghost" onClick={props.onVolumeClick} title={player().muted ? "Volume Muted" : "Volume " + player().volume}>
                    <Show when={!player().muted} fallback={
                      <IconRiMediaVolumeMuteFill />
                    }>
                      {player().volume}
                    </Show>
                  </Button>
                  <Button disabled={playerDisabled()} size="icon" variant="ghost" onClick={props.onVolumeUpClick} title="Volume Up">
                    <RiMediaVolumeUpFill class={iconClass} />
                  </Button>
                </VolumeGroup>
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
                          onClick={[props.onPlayerChange, player.id]}
                          size="sm"
                          variant={player.id == props.player?.id ? "default" : "ghost"}
                          class={style({ ...mixin.row("2",), width: "100%", justifyContent: "start" })}
                        >
                          <ConnectionIndicator connected={player.connected} disconnected={!player.connected} />
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
    </Root>
  )
}
