import { styled } from "@macaron-css/solid";
import { For } from "solid-js"
import { usePlayersQuery } from "~/hooks/api"
import { Table, TableBody, TableCaption, TableData, TableHead, TableHeader, TableRow } from "~/ui/Table"
import { mixin, theme } from "~/ui/theme";

const Root = styled("div", {
  base: {
    display: "flex",
    justifyContent: "center",
    padding: theme.space[2]
  },
});

const Content = styled("div", {
  base: {
    ...mixin.stack("2"),
    maxWidth: theme.size.xl,
    width: "100%"
  }
})

export function Players() {
  const playersQuery = usePlayersQuery()

  return (
    <Root>
      <Content>
        <Table>
          <TableCaption>Players</TableCaption>
          <TableHeader>
            <TableRow>
              <TableHead>ID</TableHead>
              <TableHead>Name</TableHead>
              <TableHead>Token</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <For each={playersQuery.data}>
              {p => (
                <TableRow>
                  <TableData>{p.id}</TableData>
                  <TableData>{p.name}</TableData>
                  <TableData>{p.token}</TableData>
                </TableRow>)
              }
            </For>
          </TableBody>
        </Table>
      </Content>
    </Root>
  )
}
