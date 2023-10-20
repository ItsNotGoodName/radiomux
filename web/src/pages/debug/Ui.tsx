import { styled } from "@macaron-css/solid";
import { For, } from "solid-js";
import { Button } from "~/ui/Button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "~/ui/Card";
import { mixin, theme } from "~/ui/theme";
import { Input } from "~/ui/Input";
import { Table, TableBody, TableCaption, TableData, TableHead, TableHeader, TableRow } from "~/ui/Table";
import { Badge } from "~/ui/Badge";
import { DialogContent, DialogHeaderDescription, DialogFooter, DialogHeader, DialogHeaderTitle, DialogRoot, DialogPortal, DialogOverlay, DialogTrigger, DialogCloseButton } from "~/ui/Dialog";
import { Textarea } from "~/ui/Textarea";
import { Label } from "~/ui/Label";
import { As } from "@kobalte/core";
import { style } from "@macaron-css/core";
import { DropdownMenuArrow, DropdownMenuContent, DropdownMenuItem, DropdownMenuPortal, DropdownMenuRoot, DropdownMenuTrigger } from "~/ui/DropdownMenu";
import { ToastList, ToastRegion, ToastTitle, ToastCloseButton, ToastContent, ToastProgressTrack, ToastProgressFill, ToastDescription, toast } from "~/ui/Toast";
import { Portal } from "solid-js/web";

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
  const showToast = () => {
    toast.custom(() =>
      <ToastContent>
        <ToastCloseButton />
        <ToastTitle>Title</ToastTitle>
        <ToastDescription>Description</ToastDescription>
        <ToastProgressTrack>
          <ToastProgressFill />
        </ToastProgressTrack>
      </ToastContent>
    )
    toast.show("Hello World")
  }

  return (
    <Root>
      <Portal>
        <ToastRegion>
          <ToastList />
        </ToastRegion>
      </Portal>
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
      <Button onClick={showToast} >Toast</Button>
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
      <DropdownMenuRoot>
        <DropdownMenuTrigger asChild>
          <As component={Button}>
            Dropdown
          </As>
        </DropdownMenuTrigger>
        <DropdownMenuPortal>
          <DropdownMenuContent>
            <DropdownMenuArrow />
            <DropdownMenuItem asChild>
              <As component={Button} class={style({ width: "100%" })}>
                Button
              </As>
            </DropdownMenuItem>
            <DropdownMenuItem asChild>
              <As component={Button} class={style({ width: "100%" })}>
                Button
              </As>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenuPortal>
      </DropdownMenuRoot>
      <Stack>
        <Label>Input Label</Label>
        <Input placeholder="Placeholder" />
      </Stack>
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
      <DialogRoot>
        <DialogTrigger asChild>
          <As component={Button}>Dialog</As>
        </DialogTrigger>
        <DialogPortal>
          <DialogOverlay />
          <DialogContent>
            <DialogHeader>
              <DialogCloseButton />
              <DialogHeaderTitle>Header Title</DialogHeaderTitle>
              <DialogHeaderDescription>
                Header Description
              </DialogHeaderDescription>
            </DialogHeader>
            <DialogFooter>
              <Button>Footer Button</Button>
            </DialogFooter>
          </DialogContent>
        </DialogPortal>
      </DialogRoot>
      <Stack>
        <Label>Textarea Label</Label>
        <Textarea placeholder="Placeholder"></Textarea>
      </Stack>
    </Root>
  )
}

