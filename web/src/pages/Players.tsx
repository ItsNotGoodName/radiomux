import { style } from "@macaron-css/core";
import { styled } from "@macaron-css/solid";
import { createForm } from "@modular-forms/solid";
import { For, createMemo, } from "solid-js"
import { playerQrUrl } from "~/api";
import { CreatePlayer, Player, UpdatePlayer } from "~/api/client.gen";
import { usePlayerCreateMutation, usePlayerDeleteMutation, usePlayerListQuery, usePlayerTokenRegenerateMutation, usePlayerUpdateMutation } from "~/hooks/api";
import { useSelection } from "~/hooks/useSelection";
import { Button } from "~/ui/Button";
import { Checkbox } from "~/ui/Checkbox";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "~/ui/Dialog";
import { Input } from "~/ui/Input";
import { Label } from "~/ui/Label";
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

const Row = styled("div", {
  base: {
    ...mixin.row("2"),
    flexWrap: "wrap"
  }
})

const formClass = style({
  ...mixin.stack("4")
})

const FormControl = styled("div", {
  base: {
    ...mixin.stack("2")
  }
})

export function Players() {
  const playerListQuery = usePlayerListQuery()
  const playerListSelection = useSelection(() => playerListQuery.data?.players ?? [])
  const playerDeleteMutation = usePlayerDeleteMutation()

  return (
    <Root>
      <Content>
        <Row>
          <Dialog button={(ref) => <Button disabled={playerListSelection.ids().length > 0} ref={ref} size="sm">Create</Button>}>
            {(setOpen) => {
              const playerCreateMutation = usePlayerCreateMutation()
              const [form, { Form, Field }] = createForm<CreatePlayer>();

              return (
                <DialogContent>
                  <DialogHeader>
                    <DialogTitle>Create Player</DialogTitle>
                  </DialogHeader>
                  <Form
                    onSubmit={(data) => playerCreateMutation.mutateAsync(data).then(() => setOpen(false))}
                    class={formClass}
                  >
                    <Field name="name">
                      {(fields, props) => (
                        <FormControl>
                          <Label>Name</Label>
                          <Input value={fields.value} {...props} placeholder="Name" />
                        </FormControl>
                      )}
                    </Field>
                    <Button disabled={form.submitting} type="submit">Create</Button>
                  </Form>
                </DialogContent>
              )
            }}
          </Dialog>
          <Dialog button={(ref) => <Button disabled={playerListSelection.ids().length != 1} ref={ref} size="sm" variant="secondary">View</Button>}>
            {(setOpen) => {
              const player = createMemo((): Player => {
                if (playerListSelection.ids().length == 0) {
                  setOpen(false)
                  return { id: 0, name: "", token: "" }
                }
                const id = playerListSelection.ids()[0]

                const player = playerListQuery.data?.players.find((v) => v.id == id)
                if (player == undefined) {
                  setOpen(false)
                  return { id: 0, name: "", token: "" }
                }

                return player
              })

              const playerTokenRegenerateMutation = usePlayerTokenRegenerateMutation()

              const imgUrl = () => playerQrUrl(player().id)

              return (
                <DialogContent>
                  <DialogHeader>
                    <DialogTitle>{player().name}</DialogTitle>
                  </DialogHeader>
                  <img src={imgUrl()}></img>
                  <Button disabled={playerTokenRegenerateMutation.isLoading} size="sm" onClick={() => playerTokenRegenerateMutation.mutate(player().id)}>
                    Regenereate Token
                  </Button>
                </DialogContent>
              )
            }}
          </Dialog>
          <Dialog button={(ref) => <Button disabled={playerListSelection.ids().length != 1} ref={ref} size="sm" variant="outline">Update</Button>}>
            {(setOpen) => {
              const id = playerListSelection.ids()[0]
              const player = playerListQuery.data?.players.find((v) => v.id == id) as Player

              const playerUpdateMutation = usePlayerUpdateMutation()
              const [form, { Form, Field }] = createForm<UpdatePlayer>({
                initialValues: {
                  name: player.name
                },
              });

              return (
                <DialogContent>
                  <DialogHeader>
                    <DialogTitle>Update Player</DialogTitle>
                  </DialogHeader>
                  <Form
                    onSubmit={(data) => playerUpdateMutation.mutateAsync({ ...data, id }).then(() => setOpen(false))}
                    class={formClass}
                  >
                    <Field name="name">
                      {(fields, props) => (
                        <FormControl>
                          <Label>Name</Label>
                          <Input value={fields.value} {...props} placeholder="Name" />
                        </FormControl>
                      )}
                    </Field>
                    <Button disabled={form.submitting} type="submit">Update</Button>
                  </Form>
                </DialogContent>
              )
            }}
          </Dialog>
          <Button onClick={() => playerDeleteMutation.mutateAsync(playerListSelection.ids()).then(() => playerListSelection.clear())} disabled={playerListSelection.ids().length == 0 || playerDeleteMutation.isLoading} size="sm" variant="destructive">Delete</Button>
        </Row>
        <Table>
          <TableCaption>{playerListQuery.data?.count ?? 0} Players</TableCaption>
          <TableHeader>
            <TableRow>
              <TableHead>
                <Checkbox
                  checked={playerListSelection.isMultiple()}
                  setChecked={() => playerListSelection.checkMultiple(!playerListSelection.isMultiple())}
                />
              </TableHead>
              <TableHead>ID</TableHead>
              <TableHead>Name</TableHead>
              <TableHead>Token</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <For each={playerListQuery.data?.players}>
              {p => (
                <TableRow data-state={playerListSelection.record()[p.id] ? "selected" : ""}>
                  <TableData>
                    <Checkbox
                      checked={playerListSelection.record()[p.id]}
                      setChecked={(checked) => playerListSelection.check(p.id, checked)}
                    />
                  </TableData>
                  <TableData>{p.id}</TableData>
                  <TableData>{p.name}</TableData>
                  <TableData>{p.token}</TableData>
                </TableRow>)
              }
            </For>
          </TableBody>
        </Table>
      </Content>
    </Root >
  )
}
