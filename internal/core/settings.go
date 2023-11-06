package core

import (
	"fmt"
	"net/url"
	"slices"
)

var Settings settings

type settings struct {
	WSPrefix   string
	HTTPPrefix string
}

func Init(httpURL *url.URL) {
	if httpURL == nil {
		return
	}

	wsSchema := ""
	httpSchema := ""
	if slices.Contains([]string{"https", "wss"}, httpURL.Scheme) {
		wsSchema = "wss://"
		httpSchema = "https://"
	} else {
		wsSchema = "ws://"
		httpSchema = "http://"
	}

	host := httpURL.Host

	Settings = settings{
		WSPrefix:   wsSchema + host,
		HTTPPrefix: httpSchema + host,
	}
}

func (s settings) FileURL(sourceID int64, path string) string {
	return fmt.Sprintf("%s/api/files/%d%s", s.HTTPPrefix, sourceID, path)
}

func (s settings) PlayerWSURL(p Player) string {
	return fmt.Sprintf("%s/api/android/ws?id=%d&token=%s", s.WSPrefix, p.ID, p.Token)
}
