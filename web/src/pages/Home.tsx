import { For, Show, createMemo } from 'solid-js'
import { mixin, theme } from '~/ui/theme'
import { Button } from '~/ui/Button'
import { style } from '@macaron-css/core'
import { usePlayerPresetMutation, usePlayerMediaMutation, usePresetsQuery } from '~/hooks/api'
import { Input } from '~/ui/Input'
import { useCurrentPlayer } from '~/providers/currentPlayer'
import { styled } from '@macaron-css/solid'

const Root = styled("div", {
  base: {
    display: "flex",
    justifyContent: "center",
    padding: theme.space[2]
  },
});

const Content = styled("div", {
  base: {
    ...mixin.stack("2"),
    maxWidth: theme.size.md,
    width: "100%"
  }
})

export function Home() {
  const { currentPlayerState, currentPlayerId } = useCurrentPlayer()
  const presetsQuery = usePresetsQuery()

  // Play media by preset
  const playerPlayPresetMutation = usePlayerPresetMutation()

  // Play media by URI
  const playerPlayUriMutation = usePlayerMediaMutation()
  let playerPlayUri: HTMLInputElement
  const playerPlayUriSubmit = () =>
    playerPlayUriMutation.mutate({ id: currentPlayerId(), uri: playerPlayUri.value })

  const disabled = createMemo(() => currentPlayerState() == undefined || !currentPlayerState()!.ready)

  return (
    <Root>
      <Content>
        <Show when={currentPlayerState()}>
          {currentPlayerState =>
            <div class={style({ ...mixin.stack("2") })}>
              <div class={style({ ...mixin.row("2"), alignItems: "center" })}>
                <Input class={style({ flex: "1" })} disabled={disabled() || playerPlayUriMutation.isLoading} ref={playerPlayUri!} type="text" placeholder="URL"></Input>
                <Button disabled={disabled() || playerPlayUriMutation.isLoading} onClick={playerPlayUriSubmit}>Play</Button>
              </div>
              <For each={presetsQuery.data}>
                {preset =>
                  <Button disabled={disabled() || playerPlayPresetMutation.isLoading} variant={preset.url == currentPlayerState().uri ? "default" : "outline"} onClick={() => playerPlayPresetMutation.mutate({ preset: preset.id, id: currentPlayerState().id })}>
                    <div class={style({ ...mixin.textLine() })}>
                      {preset.name}
                    </div>
                  </Button>}
              </For>
            </div>
          }
        </Show>
      </Content >
    </Root>
  )
}
