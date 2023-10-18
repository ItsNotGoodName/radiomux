export function relativeWsUrl(uri: string): string {
  return `${window.location.protocol === "https:" ? "wss:" : "ws:"}//${window.location.host}${uri}`
}

export function durationHumanize(durationMs: number): string {
  const durationS = durationMs / 1000
  const seconds = (Math.round(durationS % 60) + "").padStart(2, "0")
  const minutes = (Math.round(durationS / 60) + "").padStart(2, "0")

  return `${minutes}:${seconds}`
}

// export type ApiResponse<T> = {
//   error?: {
//     message: string
//   }
//   data?: T
//   response: Response
// }

// export function extractApiData<T>(res: ApiResponse<T>): T {
//   if (res.error) {
//     throw new Error(res.error.message)
//   }
//   return res.data as T
// }
