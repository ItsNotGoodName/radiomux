import { components } from "./types"
import { PlayerService, PresetService, StateService } from './client.gen';

export const playerService = new PlayerService("", fetch)
export const presetService = new PresetService("", fetch)
export const stateService = new StateService("", fetch)

export function playerQrUrl(id: number): string {
  return `/api/players/${id}/qr?t=` + + new Date().getTime();
}

export type ApiError = components["schemas"]["Error"]

export type ApiPlayerState = Required<components["schemas"]["PlayerStatePartial"]>

export type ApiPlayerStatePartial = components["schemas"]["PlayerStatePartial"]

export type ApiNotification = components["schemas"]["Notification"]

export type ApiEvent = ApiEventPlayerState | ApiEventPlayerStatePartial | ApiEventNotification

type EventBase<T extends components["schemas"]["EventType"], D> = {
  type: T,
  data: D,
}

export type ApiEventPlayerState = EventBase<"player_state", Array<ApiPlayerState>>

export type ApiEventPlayerStatePartial = EventBase<"player_state_partial", Array<ApiPlayerStatePartial>>

export type ApiEventNotification = EventBase<"notification", ApiNotification>
