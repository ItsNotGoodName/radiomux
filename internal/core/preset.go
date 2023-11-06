package core

import (
	"net/url"
	"strconv"
	"strings"
)

type Preset struct {
	ID   int64
	Name string
	Slug *url.URL
}

func PresetSlugParse(slug *url.URL) (any, error) {
	switch slug.Scheme {
	case "radiomux+file":
		return asPresetFile(slug)
	case "radiomux+subsonic":
		return asPresetSubsonic(slug)
	default:
		return asPresetURL(slug), nil
	}
}

func PresetSlugIsRadiomux(slug *url.URL) bool {
	return strings.HasPrefix(slug.Scheme, "radiomux+")
}

type PresetURL string

func asPresetURL(slug *url.URL) PresetURL {
	return PresetURL(slug.String())
}

type PresetFile struct {
	SourceID int64
	Path     string
}

// asPresetFile has a schema of "radiomux+file:///path/to/file?source-id=999".
func asPresetFile(slug *url.URL) (PresetFile, error) {
	sourceID, err := strconv.ParseInt(slug.Query().Get("source-id"), 10, 16)
	if err != nil {
		return PresetFile{}, err
	}

	path := slug.Path

	return PresetFile{
		SourceID: sourceID,
		Path:     path,
	}, nil
}

type PresetSubsonic struct {
	SourceID int64
	MediaID  string
}

// asPresetSubsonic has a schema of "radiomux+subsonic://?source-id=999&media-id=999".
func asPresetSubsonic(slug *url.URL) (PresetSubsonic, error) {
	sourceID, err := strconv.ParseInt(slug.Query().Get("source-id"), 10, 16)
	if err != nil {
		return PresetSubsonic{}, err
	}

	mediaID := slug.Query().Get("media-id")

	return PresetSubsonic{
		SourceID: sourceID,
		MediaID:  mediaID,
	}, nil
}
