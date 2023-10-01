export function relativeWsUrl(uri: string): string {
  return `${window.location.protocol === "https:" ? "wss:" : "ws:"}//${window.location.host}${uri}`
}

export type ApiResponse<T> = {
  error?: {
    message: string
  }
  data?: T
  response: Response
}

export function extractApiData<T>(res: ApiResponse<T>): T {
  if (res.error) {
    throw new Error(res.error.message)
  }
  return res.data as T
}
