package android

import (
	"context"
)

// BusCommand is used to run commands against a player.
type BusCommand interface {
	Handle(ctx context.Context, id int64, cmd Command) error
}

type Command interface {
	isCommand()
}

type CommandStop struct{}

func (CommandStop) isCommand() {}

type CommandPlay struct{}

func (CommandPlay) isCommand() {}

type CommandPause struct{}

func (CommandPause) isCommand() {}

type CommandPrepare struct{}

func (CommandPrepare) isCommand() {}

type CommandSetPlayWhenReady struct {
	PlayWhenReady bool
}

func (CommandSetPlayWhenReady) isCommand() {}

type CommandSetMediaItem struct {
	URI string
}

func (CommandSetMediaItem) isCommand() {}

type CommandSetVolume struct {
	Volume float64
}

func (CommandSetVolume) isCommand() {}

type CommandSetDeviceVolume struct {
	Volume int
}

func (CommandSetDeviceVolume) isCommand() {}

type CommandIncreaseDeviceVolume struct{}

func (CommandIncreaseDeviceVolume) isCommand() {}

type CommandDecreaseDeviceVolume struct{}

func (CommandDecreaseDeviceVolume) isCommand() {}

type CommandSetDeviceMuted struct {
	Muted bool
}

func (CommandSetDeviceMuted) isCommand() {}

type CommandDisconnect struct{}

func (CommandDisconnect) isCommand() {}

type CommandSyncState struct{}

func (CommandSyncState) isCommand() {}

type CommandSeekToDefaultPosition struct{}

func (CommandSeekToDefaultPosition) isCommand() {}
