import { As } from "@kobalte/core";
import { style } from "@macaron-css/core";
import { styled } from "@macaron-css/solid";
import { createForm } from "@modular-forms/solid";
import { For, Show, createSignal, } from "solid-js"
import { playerQrUrl } from "~/api";
import { CreatePlayer, UpdatePlayer } from "~/api/client.gen";
import { usePlayerCreateMutation, usePlayerDeleteMutation, usePlayerGetQuery, usePlayerListQuery, usePlayerTokenRegenerateMutation, usePlayerUpdateMutation, usePlayerWsURLQuery } from "~/hooks/api";
import { useSelection } from "~/hooks/useSelection";
import { Button } from "~/ui/Button";
import { PopoverArrow, PopoverCloseButton, PopoverContent, PopoverPortal, PopoverRoot, PopoverTrigger } from "~/ui/Popover";
import { DialogCloseButton, DialogContent, DialogHeader, DialogHeaderTitle, DialogOverlay, DialogPortal, DialogRoot, DialogTrigger } from "~/ui/Dialog";
import { Input } from "~/ui/Input";
import { Label } from "~/ui/Label";
import { Table, TableBody, TableCaption, TableData, TableHead, TableHeader, TableRow } from "~/ui/Table"
import { mixin, theme } from "~/ui/theme";
import { CheckboxInput, CheckboxControl, CheckboxIndicator, CheckboxRoot, CheckboxIcon } from "~/ui/Checkbox";

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

function CreateDialog(props: { disabled: boolean }) {
  const [open, setOpen] = createSignal(false);

  const playerCreateMutation = usePlayerCreateMutation()
  const [form, { Form, Field }] = createForm<CreatePlayer>();

  return (
    <DialogRoot open={open()} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <As component={Button} disabled={props.disabled} size="sm">
          Create
        </As>
      </DialogTrigger>
      <DialogPortal>
        <DialogOverlay />
        <DialogContent>
          <DialogCloseButton />
          <DialogHeader>
            <DialogHeaderTitle>Create Player</DialogHeaderTitle>
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
      </DialogPortal>
    </DialogRoot>
  )
}

function ViewDialog(props: { disabled: boolean, id?: number }) {
  return (
    <DialogRoot>
      <DialogTrigger asChild>
        <As component={Button} disabled={props.disabled} size="sm" variant="secondary">
          View
        </As>
      </DialogTrigger>
      <DialogPortal>
        <DialogOverlay />
        <DialogContent>
          <DialogCloseButton />
          <Show when={props.id}>
            {id => {
              // Queries
              const player = usePlayerGetQuery(id)
              const playerWsURLQuery = usePlayerWsURLQuery(id)
              const imgUrl = () => playerQrUrl(id())

              // Mutations
              const playerTokenRegenerateMutation = usePlayerTokenRegenerateMutation()

              return (
                <>
                  <DialogHeader>
                    <DialogHeaderTitle>{player.data?.name}</DialogHeaderTitle>
                  </DialogHeader>
                  <div class={style({ wordBreak: "break-all" })}>
                    {playerWsURLQuery.data}
                  </div>
                  <img src={imgUrl()}></img>
                  <Button disabled={playerTokenRegenerateMutation.isLoading} size="sm" onClick={() => playerTokenRegenerateMutation.mutate(id())}>
                    Regenereate Token
                  </Button>
                </>
              )
            }}
          </Show>
        </DialogContent>
      </DialogPortal>
    </DialogRoot >
  )
}

function UpdateDialog(props: { disabled: boolean, id: number }) {
  const [open, setOpen] = createSignal(false);

  return (
    <DialogRoot open={open()} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <As component={Button} disabled={props.disabled} size="sm" variant="outline">
          Update
        </As>
      </DialogTrigger>
      <DialogPortal>
        <DialogOverlay />
        <DialogContent>
          <DialogCloseButton />
          <DialogHeader>
            <DialogHeaderTitle>Update Player</DialogHeaderTitle>
          </DialogHeader>
          <Show when={props.id}>
            {id => {
              // Queries
              const player = usePlayerGetQuery(id, {
                keepPreviousData: false,
              })

              return (
                <Show when={player.data}>
                  {player => {
                    const playerUpdateMutation = usePlayerUpdateMutation()
                    const [form, { Form, Field }] = createForm<UpdatePlayer>({
                      initialValues: {
                        ...player()
                      }
                    });

                    return (
                      <Form
                        onSubmit={(data) => playerUpdateMutation.mutateAsync({ ...data, id: id() }).then(() => setOpen(false))}
                        class={formClass}
                      >
                        <Field name="name">
                          {(fields, props) => (
                            <FormControl>
                              <Label>Name</Label>
                              <Input value={fields.value} placeholder="Name" {...props} />
                            </FormControl>
                          )}
                        </Field>
                        <Button disabled={form.submitting} type="submit">Update</Button>
                      </Form>
                    )
                  }}
                </Show>
              )
            }}
          </Show>
        </DialogContent>
      </DialogPortal>
    </DialogRoot>
  )
}

function DeletePopover(props: { disabled: boolean, ids: Array<number>, onDelete: () => void }) {
  const [open, setOpen] = createSignal(false);

  const playerDeleteMutation = usePlayerDeleteMutation()

  return (
    <PopoverRoot open={open()} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <As component={Button} disabled={props.disabled} variant="destructive" size="sm">
          Delete
        </As>
      </PopoverTrigger>
      <PopoverPortal>
        <PopoverContent class={style({ ...mixin.stack("2"), padding: theme.space[2], width: theme.space[48] })}>
          <PopoverArrow />
          <div>Are you sure you wish to delete {props.ids.length} items?</div>
          <div class={style({ ...mixin.row("2"), flexDirection: "row-reverse" })}>
            <PopoverCloseButton asChild>
              <As component={Button} size="sm">
                No
              </As>
            </PopoverCloseButton>
            <Button
              onClick={() => playerDeleteMutation.mutateAsync(props.ids).then(() => {
                props.onDelete()
                setOpen(false)
              })}
              disabled={props.disabled || playerDeleteMutation.isLoading}
              size="sm"
              variant="destructive"
            >
              Yes
            </Button>
          </div>
        </PopoverContent>
      </PopoverPortal>
    </PopoverRoot>
  )
}

export function Players() {
  const playerListQuery = usePlayerListQuery()
  const playerListSelection = useSelection(() => playerListQuery.data?.players ?? [])

  return (
    <Root>
      <Content>
        <Row>
          <CreateDialog disabled={playerListSelection.ids().length > 0} />
          <ViewDialog disabled={playerListSelection.ids().length != 1} id={playerListSelection.ids()[0]} />
          <UpdateDialog disabled={playerListSelection.ids().length != 1} id={playerListSelection.ids()[0]} />
          <DeletePopover disabled={playerListSelection.ids().length == 0} ids={playerListSelection.ids()} onDelete={playerListSelection.clear} />
        </Row>
        <Table>
          <TableCaption>{playerListQuery.data?.count ?? 0} Players</TableCaption>
          <TableHeader>
            <TableRow>
              <TableHead>
                <CheckboxRoot
                  checked={playerListSelection.isMultiple()}
                  onChange={() => playerListSelection.checkMultiple(!playerListSelection.isMultiple())}
                >
                  <CheckboxInput />
                  <CheckboxControl>
                    <CheckboxIndicator >
                      <CheckboxIcon />
                    </CheckboxIndicator>
                  </CheckboxControl>
                </CheckboxRoot>
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
                    <CheckboxRoot
                      checked={playerListSelection.record()[p.id]}
                      onChange={() => playerListSelection.check(p.id, !playerListSelection.record()[p.id])}
                    >
                      <CheckboxInput />
                      <CheckboxControl>
                        <CheckboxIndicator >
                          <CheckboxIcon />
                        </CheckboxIndicator>
                      </CheckboxControl>
                    </CheckboxRoot>
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
