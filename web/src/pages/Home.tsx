import { For, createMemo } from 'solid-js'
import { minScreen, mixin, theme, tw } from '~/ui/theme'
import { Button } from '~/ui/Button'
import { SeparatorRoot } from '~/ui/Seperator'
import { style } from '@macaron-css/core'
import { Input } from '~/ui/Input'
import { useCurrentPlayer } from '~/providers/currentPlayer'
import { styled } from '@macaron-css/solid'
import { usePresetListQuery, useStateMediaSetMutation } from '~/hooks/api'
import { toastWebrpcError } from '~/common/error'
import { useWS } from '~/providers/ws'
import { ConnectionIndicator } from '~/components/Player'

const Text = styled("div", {
  base: {
    ...mixin.textLine()
  }
})

const Root = styled("div", {
  base: {
    ...mixin.row("4"),
    justifyContent: "center",
    padding: theme.space[4]
  },
});

const Content = styled("div", {
  base: {
    ...mixin.stack("2"),
    maxWidth: theme.size.md,
    width: "100%",
  }
})

const PlayerListRoot = styled("div", {
  base: {
    ...mixin.textLine(),
    display: "none",
    width: theme.space[48],
    flexShrink: 0,
    "@media": {
      [minScreen.lg]: {
        display: "unset"
      },
    },
  }
})

const PlayerListContent = styled("div", {
  base: {
    border: `1px solid ${theme.color.border}`,
    borderRadius: theme.borderRadius.ok,
    padding: theme.space[2],
  }
})

const PlayerListList = styled("div", {
  base: {
    display: "flex",
    flexDirection: "column"
  }
})

const PlayerListTitle = styled("div", {
  base: {
    ...tw.textLg,
  }
})

export function Home() {
  // Queries
  const presetListQuery = usePresetListQuery()

  // Mutations
  const stateMediaSetMutation = useStateMediaSetMutation()

  const { currentPlayerState, currentPlayerId, setCurrentPlayerId } = useCurrentPlayer()
  // TODO: if this is called more than once in other components, then move it into useCurrentPlayer 
  const disabled = createMemo(() => currentPlayerState() == undefined || !currentPlayerState()!.ready || stateMediaSetMutation.isLoading)

  let uriElement: HTMLInputElement
  const onUriSubmit = () =>
    stateMediaSetMutation.mutateAsync({ id: currentPlayerId(), uri: uriElement.value }).catch(toastWebrpcError)

  const onPresetClick = (presetId: number) =>
    stateMediaSetMutation.mutateAsync({ presetId: presetId, id: currentPlayerState()?.id || 0 }).catch(toastWebrpcError)

  const { playerStates } = useWS()
  const toggleCurrentPlayerId = (id: number) => setCurrentPlayerId((prev) => prev == id ? 0 : id)

  return (
    <Root>
      <Content class={style({ ...mixin.stack("2") })}>
        <div class={style({ ...mixin.row("2"), alignItems: "center" })}>
          <Input
            class={style({ flex: "1" })}
            type="text"
            placeholder="URL"
            disabled={disabled()}
            ref={uriElement!}
          />
          <Button
            disabled={disabled()}
            onClick={onUriSubmit}
          >
            Play
          </Button>
        </div>
        <For each={presetListQuery.data}>
          {preset =>
            <Button
              disabled={disabled()}
              variant={preset.url == currentPlayerState()?.uri ? "default" : "outline"}
              onClick={[onPresetClick, preset.id]}
            >
              <Text>{preset.name}</Text>
            </Button>
          }
        </For>
      </Content>
      <PlayerListRoot>
        <PlayerListContent>
          <PlayerListTitle>Players</PlayerListTitle>
          <SeparatorRoot />
          <PlayerListList>
            <For each={playerStates}>
              {player => (
                <Button
                  onClick={[toggleCurrentPlayerId, player.id]}
                  size="sm"
                  variant={currentPlayerId() == player.id ? "default" : "ghost"}
                  class={style({ ...mixin.row("2",), width: "100%", justifyContent: "start" })}
                >
                  <ConnectionIndicator connected={player.connected} disconnected={!player.connected} />
                  <Text title={player.name}>{player.name}</Text>
                </Button>
              )}
            </For>
          </PlayerListList>
        </PlayerListContent>
      </PlayerListRoot>
    </Root>
  )
}
