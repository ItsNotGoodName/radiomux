import { styled } from "@macaron-css/solid";
import { Component, For, } from "solid-js";
import { Button } from "~/ui/Button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "~/ui/Card";
import { mixin, theme } from "~/ui/theme";
import { Dropdown } from "~/ui/Dropdown";

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

export const Ui: Component = () => {
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
        {(ref) => (
          <Card ref={ref}>
            <CardHeader>
              <CardTitle>Hello World</CardTitle>
            </CardHeader>
            <CardContent>
              Hello World
            </CardContent>
          </Card>
        )}
      </Dropdown>
    </Root>
  )
}

