import { styled } from "@macaron-css/solid";
import { tw } from "./theme";

export const Label = styled("label", {
  base: {
    ...tw.textSm,
    fontWeight: 500,
    lineHeight: 1,
    ":disabled": {
      cursor: "not-allowed",
      opacity: "50%"
    }
  }
})
