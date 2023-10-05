import { For, Show } from 'solid-js'
import { theme, mixin } from '~/ui/theme'
import { styled, } from '@macaron-css/solid'
import { Button } from '~/ui/Button'
import { style } from '@macaron-css/core'
import { usePlayerPresetMutation, usePlayerMediaMutation, usePlayersQuery, usePresetsQuery } from '~/hooks/api'
import { Input } from '~/ui/Input'
import { useCurrentPlayer } from '~/providers/currentPlayer'
import { Table, TableBody, TableCaption, TableData, TableHead, TableHeader, TableRow } from '~/ui/Table'

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

function PlayerPresets() {
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

function DebugCard() {
  const presetsQuery = usePresetsQuery()
  const playersQuery = usePlayersQuery()

  return (
    <>
      <Table>
        <TableCaption>Players</TableCaption>
        <TableHeader>
          <TableRow>
            <TableHead>ID</TableHead>
            <TableHead>Name</TableHead>
            <TableHead>Token</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <For each={playersQuery.data}>
            {p => (
              <TableRow>
                <TableData>{p.id}</TableData>
                <TableData>{p.name}</TableData>
                <TableData>{p.token}</TableData>
              </TableRow>)
            }
          </For>
        </TableBody>
      </Table>
      <Table>
        <TableCaption>Presets</TableCaption>
        <TableHeader>
          <TableRow>
            <TableHead>ID</TableHead>
            <TableHead>Name</TableHead>
            <TableHead>URL</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <For each={presetsQuery.data}>
            {p => (
              <TableRow>
                <TableData>{p.id}</TableData>
                <TableData>{p.name}</TableData>
                <TableData>{p.url}</TableData>
              </TableRow>)
            }
          </For>
        </TableBody>
      </Table>
    </>
  )
}

export function Home() {
  return (
    <Root>
      <Content>
        <PlayerPresets />
        <DebugCard />
      </Content>
    </Root >
  )
}
