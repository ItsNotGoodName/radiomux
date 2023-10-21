import { For, createMemo } from 'solid-js'
import { mixin, theme } from '~/ui/theme'
import { Button } from '~/ui/Button'
import { style } from '@macaron-css/core'
import { Input } from '~/ui/Input'
import { useCurrentPlayer } from '~/providers/currentPlayer'
import { styled } from '@macaron-css/solid'
import { usePresetListQuery, useStateMediaSetMutation } from '~/hooks/api'
import { toastWebrpcError } from '~/common/toast'

const Root = styled("div", {
  base: {
    display: "flex",
    justifyContent: "center",
    padding: theme.space[4]
  },
});

const Content = styled("div", {
  base: {
    ...mixin.stack("2"),
    maxWidth: theme.size.md,
    width: "100%"
  }
})

const Text = styled("div", {
  base: {
    ...mixin.textLine()
  }
})

export function Home() {
  // Queries
  const presetListQuery = usePresetListQuery()

  // Mutations
  const stateMediaSetMutation = useStateMediaSetMutation()

  const { currentPlayerState, currentPlayerId } = useCurrentPlayer()
  // TODO: if this is called more than once in other components, then move it into useCurrentPlayer 
  const disabled = createMemo(() => currentPlayerState() == undefined || !currentPlayerState()!.ready || stateMediaSetMutation.isLoading)

  let uriElement: HTMLInputElement
  const onUriSubmit = () =>
    stateMediaSetMutation.mutateAsync({ id: currentPlayerId(), uri: uriElement.value }).catch(toastWebrpcError)

  const onPresetClick = (presetId: number) =>
    stateMediaSetMutation.mutateAsync({ presetId: presetId, id: currentPlayerState()?.id || 0 }).catch(toastWebrpcError)

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
    </Root>
  )
}
