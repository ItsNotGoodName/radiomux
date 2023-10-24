/**
 * This file was auto-generated by openapi-typescript.
 * Do not make direct changes to the file.
 */


export interface paths {
  "/players/{id}/qr": {
    /** @description Get player's QR code for the Android app. */
    get: {
      parameters: {
        path: {
          id: number;
        };
      };
      responses: {
        200: {
          content: {
            "image/png": string;
          };
        };
        /** @description unexpected error */
        default: {
          content: {
            "application/json": components["schemas"]["Error"];
          };
        };
      };
    };
  };
}

export type webhooks = Record<string, never>;

export interface components {
  schemas: {
    Error: {
      message: string;
    };
    Event: unknown;
    EventBase: {
      type: components["schemas"]["EventType"];
    };
    /** @enum {string} */
    EventType: "player_state" | "player_state_partial" | "notification";
    EventDataPlayerStatePartial: unknown;
    EventDataPlayerState: unknown;
    EventDataNotification: unknown;
    Notification: {
      error: boolean;
      title: string;
      description: string;
    };
    /** @enum {string} */
    PlayerPlaybackState: "IDLE" | "BUFFERING" | "READY" | "ENDED";
    PlayerStatePartial: {
      /** Format: int64 */
      id: number;
      name?: string;
      connected?: boolean;
      ready?: boolean;
      min_volume?: number;
      max_volume?: number;
      volume?: number;
      muted?: boolean;
      playback_state?: components["schemas"]["PlayerPlaybackState"];
      playback_error?: string;
      playing?: boolean;
      loading?: boolean;
      title?: string;
      genre?: string;
      station?: string;
      uri?: string;
      timeline_is_seekable?: boolean;
      timeline_is_live?: boolean;
      timeline_is_placeholder?: boolean;
      /** Format: int64 */
      timeline_default_position?: number;
      /** Format: int64 */
      timeline_duration?: number;
      /** Format: int64 */
      position?: number;
      /** Format: date-time */
      position_time?: string;
    };
    PlayerState: components["schemas"]["PlayerStatePartial"] & Record<string, never>;
  };
  responses: never;
  parameters: never;
  requestBodies: never;
  headers: never;
  pathItems: never;
}

export type $defs = Record<string, never>;

export type external = Record<string, never>;

export type operations = Record<string, never>;
