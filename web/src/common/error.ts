import { WebrpcError } from "~/api/client.gen"
import { toast } from "~/ui/Toast"

export function webrpcErrorMessage(err: WebrpcError) {
  if (err.cause != undefined) {
    return `${err.message}: ${err.cause}`
  } else {
    return err.message
  }
}

export function toastWebrpcError(err: WebrpcError) {
  toast.error("Webrpc Error", webrpcErrorMessage(err))
}

