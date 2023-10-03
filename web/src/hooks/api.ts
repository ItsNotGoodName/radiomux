import { GET, POST, ApiPlayer } from "~/api"
import { createMutation, createQuery } from '@tanstack/solid-query'
import { extractApiData } from '~/common'

export const usePlayersQuery = () => createQuery(() => ["/players"], () =>
  GET("/players", {}).
    then(extractApiData))

export const usePresetsQuery = () => createQuery(() => ["/presets"],
  () => GET("/presets", {}).
    then(extractApiData))

export const usePlayerQuery = (id: number) => createQuery<ApiPlayer, Error>(() => ["/presets/{id}", id],
  () => GET("/players/{id}", { params: { path: { id } } }).
    then(extractApiData))

export const usePlayerPlayMutation = () => createMutation<void, Error, number>({
  mutationFn: (id: number) =>
    POST("/players/{id}/play", { params: { path: { id } } })
      .then(extractApiData)
})

export const usePlayerPauseMutation = () => createMutation<void, Error, number>({
  mutationFn: (id: number) =>
    POST("/players/{id}/pause", { params: { path: { id } } })
      .then(extractApiData)
})

export const usePlayerSeekMutation = () => createMutation<void, Error, number>({
  mutationFn: (id: number) =>
    POST("/players/{id}/seek", { params: { path: { id } } })
      .then(extractApiData)
})

export const usePlayerVolumeMutation = () => createMutation<void, Error, { id: number, delta?: number, mute?: boolean }>({
  mutationFn: ({ id, delta, mute }) =>
    POST("/players/{id}/volume", { params: { path: { id }, query: { delta, mute } } })
      .then(extractApiData)
})

export const usePlayerPresetMutation = () => createMutation<void, Error, { preset: number, id: number }>({
  mutationFn: ({ preset, id }) =>
    POST("/players/{id}/preset", { params: { path: { id }, query: { preset } } }).
      then(extractApiData)
})

export const usePlayerMediaMutation = () => createMutation<void, Error, { uri: string, id: number }>({
  mutationFn: ({ uri, id }) => POST("/players/{id}/media", { params: { path: { id }, query: { uri } } }).
    then(extractApiData)
})
