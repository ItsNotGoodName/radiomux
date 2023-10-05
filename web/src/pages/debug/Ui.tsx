import { styled } from "@macaron-css/solid";
import { For, } from "solid-js";
import { Button } from "~/ui/Button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "~/ui/Card";
import { mixin, theme } from "~/ui/theme";
import { Dropdown, DropdownCard, DropdownCardContent } from "~/ui/Dropdown";
import { Input } from "~/ui/Input";
import { Table, TableBody, TableData, TableHead, TableHeader, TableRow } from "~/ui/Table";

const Root = styled("div", {
  base: {
    ...mixin.stack("4"),
    padding: theme.space[4]
  }
})

const Row = styled("div", {
  base: {
    ...mixin.row("4")
  }
})

const Stack = styled("div", {
  base: {
    ...mixin.stack("4")
  }
})

const Text = styled("div", {
  base: {
    ...mixin.textLine()
  }
})

const buttonVariants = ["default", "destructive", "outline", "secondary", "ghost", "link"]
const buttonSizes = ["icon", "sm", "default", "lg"]

export function Ui() {
  return (
    <Root>
      <Card>
        <CardHeader>
          <CardTitle>Title</CardTitle>
          <CardDescription>Description</CardDescription>
        </CardHeader>
        <CardContent>
          Content
        </CardContent>
        <CardFooter>Footer</CardFooter>
      </Card>
      <Stack>
        <For each={buttonSizes}>
          {size =>
            <Row>
              <For each={buttonVariants}>
                {variant =>
                  <Button variant={variant as any} size={size as any}>
                    <Text>
                      Default
                    </Text>
                  </Button>
                }
              </For>
            </Row>
          }
        </For>
      </Stack>
      <Dropdown button={(ref) => <Button ref={ref}>Hi</Button>}>
        {ref => (
          <DropdownCard ref={ref}>
            <DropdownCardContent>
              Hello World
            </DropdownCardContent>
          </DropdownCard>
        )}
      </Dropdown>
      <Input placeholder="Placeholder" />
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Test 1</TableHead>
            <TableHead>Test 1</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow>
            <TableData>Data 1</TableData>
            <TableData>Data 2</TableData>
          </TableRow>
          <TableRow>
            <TableData>Data 1</TableData>
            <TableData>Data 2</TableData>
          </TableRow>
        </TableBody>
      </Table>
    </Root>
  )
}

