import { style } from "@macaron-css/core";
import { styled } from "@macaron-css/solid";
import { RiDeviceDeviceFill, RiMediaPauseFill, RiMediaPlayFill, RiMediaVolumeDownFill, RiMediaVolumeMuteFill, RiMediaVolumeUpFill } from "solid-icons/ri";
import { Component, For, Show, } from "solid-js";
import { Button } from "~/ui/Button";
import { minScreen, mixin, theme, tw } from "~/ui/theme";
import { Dropdown } from "~/ui/Dropdown";

const Root = styled("div", {
  base: {
    background: theme.color.card,
    color: theme.color.cardForeground,
    height: theme.space[32],
    borderTop: `1px solid ${theme.color.border}`,
    "@media": {
      [minScreen.md]: {
        height: theme.space[16],
      },
    },
  }
})

const Content = styled("div", {
  base: {
    display: "flex",
    "height": "100%",
    paddingLeft: theme.space[4],
    paddingRight: theme.space[4],
    flexDirection: "column",
    "@media": {
      [minScreen.md]: {
        flexDirection: "row",
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
    color: theme.color.cardForeground,
    border: `1px solid ${theme.color.border}`,
    borderRadius: theme.borderRadius.full
  }
})

const ConnectionIndicator = styled("div", {
  base: {
    ...mixin.size("4"),
    borderRadius: theme.borderRadius.full,
  },
  variants: {
    connected: {
      true: {
        background: "lime"
      },
      false: {
        background: "red"
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

const Popover = styled("div", {
  base: {
    ...tw.shadowMd,
    background: theme.color.popover,
    color: theme.color.popoverForeground,
    border: `1px solid ${theme.color.border}`,
    borderRadius: theme.borderRadius.lg,
  }
})

type Props = {
  player?: {
    id: number
    name: string
    connected: boolean
    title: string
    playing: boolean
    volume: number
    muted: boolean
  }
  players: Array<{
    id: number
    name: string
    connected: boolean
  }>
  onPlayClick?: () => void
  playLoading?: boolean,
  onVolumeUpClick?: () => void
  onVolumeDownClick?: () => void
  onVolumeClick?: () => void
  onPlayerClick: (id: number) => void
}

export const Player: Component<Props> = (props) => {
  const volumeDisabled = () => props.player == undefined || !props.player.connected

  const playDisabled = () => props.player == undefined || !props.player.connected || props.playLoading
  const playTitle = () => {
    if (props.player == undefined) {
      return ""
    }
    return props.player?.playing ? "Playing" : "Paused"
  }

  return (
    <Root>
      <Content>
        <ContentChild>
          <Thumbnail src="/favicon.svg" />
          <Text title={props.player?.title}>{props.player?.title}</Text>
        </ContentChild>
        <ContentChild>
          <MainControl>
            <Button
              disabled={playDisabled()}
              size="icon"
              variant="ghost"
              onClick={props.onPlayClick}
              title={playTitle()}
            >
              <Show when={props.player?.playing} fallback={
                <RiMediaPlayFill class={style({ ...mixin.size("9") })} />
              }>
                <RiMediaPauseFill class={style({ ...mixin.size("9") })} />
              </Show>
            </Button>
          </MainControl>
          <SubControl>
            <Show when={props.player} fallback={<PlayerGroup />}>
              {player => (
                <>
                  <PlayerGroup>
                    <div>
                      <ConnectionIndicator connected={player().connected} />
                    </div>
                    <Text title={player().name}>{player().name}</Text>
                  </PlayerGroup>
                  <VolumeGroup>
                    <Button disabled={volumeDisabled()} size="icon" variant="ghost" onClick={props.onVolumeDownClick} title="Volume Down">
                      <RiMediaVolumeDownFill class={style({ ...mixin.size("6") })} />
                    </Button>
                    <Button disabled={volumeDisabled()} size="icon" variant="ghost" onClick={props.onVolumeClick} title={player().muted ? "Volume Muted" : undefined}>
                      <Show when={!player().muted} fallback={
                        <RiMediaVolumeMuteFill class={style({ ...mixin.size("6"), color: "red" })} />
                      }>
                        {player().volume}
                      </Show>
                    </Button>
                    <Button disabled={volumeDisabled()} size="icon" variant="ghost" onClick={props.onVolumeUpClick} title="Volume Up">
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
                <div ref={ref} class={style({ padding: theme.space[2], width: theme.space[60] })}>
                  <Popover class={style({ padding: theme.space[2] })}>
                    <For each={props.players}>
                      {player => (
                        <Button
                          onClick={[props.onPlayerClick, player.id]}
                          size="sm"
                          variant={player.id == props.player?.id ? "default" : "ghost"}
                          class={style({ ...mixin.row("2",), width: "100%", justifyContent: "start" })}
                        >
                          <div>
                            <ConnectionIndicator connected={player.connected} />
                          </div>
                          <Text title={player.name}>{player.name}</Text>
                        </Button>
                      )}
                    </For>
                  </Popover>
                </div>
              }
            </Dropdown>
          </SubControl>
        </ContentChild>
      </Content>
    </Root>
  )
}
