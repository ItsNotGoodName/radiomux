import { For, Show } from 'solid-js'
import { theme, mixin } from '~/ui/theme'
import { styled, } from '@macaron-css/solid'
import { Button } from '~/ui/Button'
import { style } from '@macaron-css/core'
import { usePlayerPresetMutation, usePlayerMediaMutation, usePlayersQuery, usePresetsQuery } from '~/hooks/api'
import { Card, CardContent, CardHeader } from '~/ui/Card'
import { Input } from '~/ui/Input'
import { useCurrentPlayer } from '~/providers/currentPlayer'

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
  return (
    <Root>
      <Content>
        <PlayerControl />
        <StatusCard />
      </Content>
    </Root >
  )
}

function StatusCard() {
  const presetsQuery = usePresetsQuery()
  const playersQuery = usePlayersQuery()

  return (
    <Card>
      <CardHeader>Test</CardHeader>
      <CardContent class={style({ overflowX: "auto" })}>
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Name</th>
              <th>Token</th>
            </tr>
          </thead>
          <tbody>
            <For each={playersQuery.data}>
              {p => (
                <tr>
                  <td>{p.id}</td>
                  <td>{p.name}</td>
                  <td>{p.token}</td>
                </tr>)
              }
            </For>
          </tbody>
        </table>
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Name</th>
              <th>URL</th>
            </tr>
          </thead>
          <tbody>
            <For each={presetsQuery.data}>
              {p => (
                <tr>
                  <td>{p.id}</td>
                  <td>{p.name}</td>
                  <td>{p.url}</td>
                </tr>)
              }
            </For>
          </tbody>
        </table>
      </CardContent>
    </Card>
  )
}

function PlayerControl() {
  const { currentPlayerState, currentPlayerId } = useCurrentPlayer()
  const presetsQuery = usePresetsQuery()

  // Play media by preset
  const playerPlayPresetMutation = usePlayerPresetMutation()

  // Play media by URI
  const playerPlayUriMutation = usePlayerMediaMutation()
  let playerPlayUri: HTMLInputElement
  const playerPlayUriSubmit = () =>
    playerPlayUriMutation.mutate({ id: currentPlayerId(), uri: playerPlayUri.value })

  return (
    <Show when={currentPlayerState()}>
      {currentPlayerState =>
        <div class={style({ ...mixin.stack("2") })}>
          <div class={style({ ...mixin.row("2"), alignItems: "center" })}>
            <Input class={style({ flex: "1" })} disabled={playerPlayUriMutation.isLoading} ref={playerPlayUri!} type="text" placeholder="URL"></Input>
            <Button disabled={playerPlayUriMutation.isLoading} onClick={playerPlayUriSubmit}>Play</Button>
          </div>
          <For each={presetsQuery.data}>
            {preset =>
              <Button disabled={playerPlayPresetMutation.isLoading} variant={preset.url == currentPlayerState().uri ? "default" : "outline"} onClick={() => playerPlayPresetMutation.mutate({ preset: preset.id, id: currentPlayerState().id })}>
                <div class={style({ ...mixin.textLine() })}>
                  {preset.name}
                </div>
              </Button>}
          </For>
        </div>
      }
    </Show>
  )
}
