import { styled } from "@macaron-css/solid";
import { For } from "solid-js";
import { Button } from "~/ui/Button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "~/ui/Card";
import { mixin, theme } from "~/ui/theme";
import { Dropdown, DropdownCard, DropdownCardContent } from "~/ui/Dropdown";
import { Input } from "~/ui/Input";
import { Table, TableBody, TableCaption, TableData, TableHead, TableHeader, TableRow } from "~/ui/Table";
import { Badge } from "~/ui/Badge";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "~/ui/Dialog";
import { Textarea } from "~/ui/Textarea";

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
const badgeVariants = ["default", "secondary", "destructive", "outline"]

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
      <Dropdown button={(ref) => <Button ref={ref}>Dropdown Button</Button>}>
        {ref => (
          <DropdownCard ref={ref}>
            <DropdownCardContent>
              Dropdown Card Content
            </DropdownCardContent>
          </DropdownCard>
        )}
      </Dropdown>
      <Input placeholder="Placeholder" />
      <Table>
        <TableCaption>Table Caption</TableCaption>
        <TableHeader>
          <TableRow>
            <TableHead>Test 1</TableHead>
            <TableHead>Test 2</TableHead>
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
      <Row>
        <For each={badgeVariants}>
          {variant =>
            <Badge variant={variant as any}>Default</Badge>
          }
        </For>
      </Row>
      <Dialog button={ref => <Button ref={ref}>Dialog</Button>}>
        {setOpen =>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Header Title</DialogTitle>
              <DialogDescription>
                Header Description
              </DialogDescription>
            </DialogHeader>
            <DialogFooter>
              <Button onClick={[setOpen, false]}>Footer Button</Button>
            </DialogFooter>
          </DialogContent>
        }
      </Dialog>
      <Textarea placeholder="Placeholder"></Textarea>
    </Root >
  )
}

