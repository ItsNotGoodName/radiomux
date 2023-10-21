import { WebrpcError } from "~/api/client.gen"
import { toast } from "~/ui/Toast"

export function toastWebrpcError(err: WebrpcError) {
  if (err.cause != undefined) {
    toast.error("Webrpc Error", `${err.message}: ${err.cause}`)
  } else {
    toast.error("Webrpc Error", err.message)
  }
}
