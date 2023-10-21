import { playerService, stateService, presetService } from "~/api"
import { CreateQueryOptions, createMutation, createQuery, useQueryClient } from '@tanstack/solid-query'
import { CreatePlayer, Player, PlayerListResult, Preset, SetStateAction, SetStateMedia, SetStateVolume, UpdatePlayer, WebrpcError } from "~/api/client.gen"
import { Accessor } from "solid-js"

export const usePlayerListQuery = (options?: CreateQueryOptions<PlayerListResult, WebrpcError>) =>
  createQuery<PlayerListResult, WebrpcError>(() => ["/players"], () =>
    playerService.playerList().
      then(res => res.res), options as any)

export const usePresetListQuery = (options?: CreateQueryOptions<Array<Preset>, WebrpcError>) =>
  createQuery<Array<Preset>, WebrpcError>(() => ["/presets"],
    () => presetService.presetList().
      then(res => res.presets), options as any)

export const usePlayerGetQuery = (id: Accessor<number>, options?: CreateQueryOptions<Player, WebrpcError>) =>
  createQuery<Player, WebrpcError>(() => ["/presets/{id}", id()],
    () => playerService.playerGet({ id: id() }).
      then(res => res.player), options as any)

export const usePlayerWsURLQuery = (id: Accessor<number>, options?: CreateQueryOptions<string, WebrpcError>) =>
  createQuery<string, WebrpcError>(() => ["playerWsURL", id()],
    () => playerService.playerWsURL({ id: id() }).
      then(res => res.url), options as any)

export const useStateActionSetMutation = (options?: Omit<CreateQueryOptions<unknown, WebrpcError>, "onSuccess">) =>
  createMutation<unknown, WebrpcError, SetStateAction>({
    ...options as any,
    mutationFn: (req) =>
      stateService.stateActionSet({ req }, options as any)
  })

export const useStateVolumeSetMutation = (options?: Omit<CreateQueryOptions<unknown, WebrpcError>, "onSuccess">) =>
  createMutation<unknown, WebrpcError, SetStateVolume>({
    ...options as any,
    mutationFn: (req) =>
      stateService.stateVolumeSet({ req }, options as any)
  })

export const useStateMediaSetMutation = (options?: Omit<CreateQueryOptions<unknown, WebrpcError>, "onSuccess">) =>
  createMutation<unknown, WebrpcError, SetStateMedia>({
    ...options as any,
    mutationFn: (req) =>
      stateService.stateMediaSet({ req }, options as any)
  })

export const usePlayerDeleteMutation = (options?: Omit<CreateQueryOptions<unknown, WebrpcError>, "onSuccess">) => {
  const queryClient = useQueryClient()
  return createMutation<unknown, WebrpcError, Array<number>>({
    ...options as any,
    mutationFn: (ids) => playerService.playerDelete({ ids }),
    onSuccess(_, ids) {
      ids.map((id) => queryClient.invalidateQueries(["/players/{id}", id]))
      queryClient.invalidateQueries(["/players"], options as any)
    },
  })
}

export const usePlayerCreateMutation = (options?: Omit<CreateQueryOptions<unknown, WebrpcError>, "onSuccess">) => {
  const queryClient = useQueryClient()
  return createMutation<unknown, WebrpcError, CreatePlayer>({
    ...options as any,
    mutationFn: (req) =>
      playerService.playerCreate({ req }),
    onSuccess() {
      queryClient.invalidateQueries(["/players"])
    },
  })
}

export const usePlayerUpdateMutation = (options?: Omit<CreateQueryOptions<unknown, WebrpcError>, "onSuccess">) => {
  const queryClient = useQueryClient()
  return createMutation<unknown, WebrpcError, UpdatePlayer>({
    ...options as any,
    mutationFn: (req) =>
      playerService.playerUpdate({ req }),
    onSuccess() {
      queryClient.invalidateQueries(["/players"])
    },
  })
}

export const usePlayerTokenRegenerateMutation = (options?: Omit<CreateQueryOptions<unknown, WebrpcError>, "onSuccess">) => {
  const queryClient = useQueryClient()
  return createMutation<unknown, WebrpcError, number>({
    ...options as any,
    mutationFn: (id) =>
      playerService.playerTokenRegenerate({ id }),
    onSuccess(_, id) {
      queryClient.invalidateQueries(["/players/{id}", id])
      queryClient.invalidateQueries(["/players"])
      queryClient.invalidateQueries(["playerWsURL", id])
    },
  })
}
