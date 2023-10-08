import { styled } from "@macaron-css/solid";
import { For } from "solid-js"
import { usePresetListQuery } from "~/hooks/api";
import { Table, TableBody, TableCaption, TableData, TableHead, TableHeader, TableRow } from "~/ui/Table"
import { mixin, theme } from "~/ui/theme";

const Root = styled("div", {
  base: {
    display: "flex",
    justifyContent: "center",
    padding: theme.space[4]
  },
});

const Content = styled("div", {
  base: {
    ...mixin.stack("2"),
    maxWidth: theme.size.xl,
    width: "100%"
  }
})

export function Presets() {
  const presetListQuery = usePresetListQuery()

  return (
    <Root>
      <Content>
        <Table>
          <TableCaption>Presets</TableCaption>
          <TableHeader>
            <TableRow>
              <TableHead>ID</TableHead>
              <TableHead>Name</TableHead>
              <TableHead>URL</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <For each={presetListQuery.data}>
              {p => (
                <TableRow>
                  <TableData>{p.id}</TableData>
                  <TableData>{p.name}</TableData>
                  <TableData>{p.url}</TableData>
                </TableRow>)
              }
            </For>
          </TableBody>
        </Table>
      </Content>
    </Root>
  )
}
