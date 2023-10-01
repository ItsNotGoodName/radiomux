import createClient from 'openapi-fetch';
import { paths, components } from "./types"

export const { POST, GET } = createClient<paths>({ baseUrl: "/api" });

export type ApiError = components["schemas"]["Error"]

export type ApiPlayer = components["schemas"]["Player"]

export type ApiPlayerState = Required<components["schemas"]["PlayerStatePartial"]>

export type ApiPlayerStatePartial = components["schemas"]["PlayerStatePartial"]


export type ApiEvent = ApiEventPlayerState | ApiEventPlayerStatePartial

type EventBase<T extends components["schemas"]["EventType"], D> = {
  type: T,
  data: D,
}

export type ApiEventPlayerState = EventBase<"player_state", Array<ApiPlayerState>>

export type ApiEventPlayerStatePartial = EventBase<"player_state_partial", Array<ApiPlayerStatePartial>>
