import { styled } from "@macaron-css/solid"
import { theme, tw } from "./theme"
import { JSX, ParentProps } from "solid-js"

const Root = styled("div", {
  base: {
    position: "relative",
    width: "100%",
    overflow: "auto"
  }
})

const TheTable = styled("table", {
  base: {
    ...tw.textSm,
    width: "100%",
    captionSide: "bottom",
    borderCollapse: "collapse"
  }
})

export function Table(props: ParentProps<JSX.HTMLAttributes<HTMLTableElement>>) {
  return (
    <Root>
      <TheTable {...props}>
        {props.children}
      </TheTable>
    </Root >
  )
}

export const TableHeader = styled("thead", {
  base: {}
})

export const TableBody = styled("tbody", {
  base: {}
})

export const TableFooter = styled("tfoot", {
  base: {
    background: theme.color.primary,
    fontWeight: "500",
    color: theme.color.primaryForeground
  }
})

export const TableRow = styled("tr", {
  base: {
    ...tw.transitionColors,
    borderBottom: `1px solid ${theme.color.border}`,
    ":hover": {
      background: theme.color.mutedHover,
    },
    selectors: {
      [`${TableHeader} &`]: {
        borderBottom: `1px solid ${theme.color.border}`,
        background: "unset"
      },
      [`${TableBody} &:last-child`]: {
        borderBottom: "none"
      },
      '&[data-state-selected]': {
        background: theme.color.muted
      },
    }
  }
})

export const TableHead = styled("th", {
  base: {
    height: theme.space[12],
    paddingLeft: theme.space[4],
    paddingRight: theme.space[4],
    textAlign: "left",
    verticalAlign: "middle",
    fontWeight: "500",
    color: theme.color.mutedForeground
    // [&:has([role=checkbox])]:pr-0"
  }
})

export const TableData = styled("td", {
  base: {
    padding: theme.space[4],
    verticalAlign: "middle",
    // [&:has([role=checkbox])]:pr-0", className)}
  }
})

export const TableCaption = styled("caption", {
  base: {
    ...tw.textSm,
    marginTop: theme.space[4],
    color: theme.color.mutedForeground,
  }
})

