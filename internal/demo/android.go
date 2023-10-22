package demo

import (
	"github.com/ItsNotGoodName/radiomux/internal/core"
	echo "github.com/labstack/echo/v4"
)

type AndroiWSServer struct{}

func (AndroiWSServer) ServeEcho(c echo.Context) error {
	return echo.ErrNotImplemented
}

func (AndroiWSServer) PlayerWSURL(p core.Player) string {
	return "wss://example.com"
}

func NewAndroidWSServer() AndroiWSServer {
	return AndroiWSServer{}
}
