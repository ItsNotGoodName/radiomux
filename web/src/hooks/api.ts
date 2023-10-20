import { playerService, stateService, presetService } from "~/api"
import { CreateQueryOptions, createMutation, createQuery, useQueryClient } from '@tanstack/solid-query'
import { CreatePlayer, Player, SetStateAction, SetStateMedia, SetStateVolume, UpdatePlayer, WebrpcError } from "~/api/client.gen"
import { Accessor } from "solid-js"

export const usePlayerListQuery = () => createQuery(() => ["/players"], () =>
  playerService.playerList().
    then(res => res.res))

export const usePresetListQuery = () => createQuery(() => ["/presets"],
  () => presetService.presetList().
    then(res => res.presets))

export const usePlayerGetQuery = (id: Accessor<number>, options?: CreateQueryOptions<Player>) => createQuery<Player, WebrpcError>(() => ["/presets/{id}", id()],
  () => playerService.playerGet({ id: id() }).
    then(res => res.player), options as any)

export const usePlayerWsURLQuery = (id: Accessor<number>, options?: CreateQueryOptions<string>) => createQuery<string, WebrpcError>(() => ["playerWsURL", id()],
  () => playerService.playerWsURL({ id: id() }).
    then(res => res.url), options as any)

export const useStateActionSetMutation = () => createMutation<unknown, WebrpcError, SetStateAction>({
  mutationFn: (req) =>
    stateService.stateActionSet({ req })
})

export const useStateVolumeSetMutation = () => createMutation<unknown, WebrpcError, SetStateVolume>({
  mutationFn: (req) =>
    stateService.stateVolumeSet({ req })
})

export const useStateMediaSetMutation = () => createMutation<unknown, WebrpcError, SetStateMedia>({
  mutationFn: (req) =>
    stateService.stateMediaSet({ req })
})

export const usePlayerDeleteMutation = () => {
  const queryClient = useQueryClient()
  return createMutation<unknown, WebrpcError, Array<number>>({
    mutationFn: (ids) => playerService.playerDelete({ ids }),
    onSuccess(_, ids) {
      ids.map((id) => queryClient.invalidateQueries(["/players/{id}", id]))
      queryClient.invalidateQueries(["/players"])
    },
  })
}

export const usePlayerCreateMutation = () => {
  const queryClient = useQueryClient()
  return createMutation<unknown, WebrpcError, CreatePlayer>({
    mutationFn: (req) =>
      playerService.playerCreate({ req }),
    onSuccess() {
      queryClient.invalidateQueries(["/players"])
    },
  })
}

export const usePlayerUpdateMutation = () => {
  const queryClient = useQueryClient()
  return createMutation<unknown, WebrpcError, UpdatePlayer>({
    mutationFn: (req) =>
      playerService.playerUpdate({ req }),
    onSuccess() {
      queryClient.invalidateQueries(["/players"])
    },
  })
}

export const usePlayerTokenRegenerateMutation = () => {
  const queryClient = useQueryClient()
  return createMutation<unknown, WebrpcError, number>({
    mutationFn: (id) =>
      playerService.playerTokenRegenerate({ id }),
    onSuccess(_, id) {
      queryClient.invalidateQueries(["/players/{id}", id])
      queryClient.invalidateQueries(["/players"])
      queryClient.invalidateQueries(["playerWsURL", id])
    },
  })
}
