package android

import "time"

type TimelineWindow struct {
	// Whether it's possible to seek within this window.
	IsSeekable bool
	// Whether this window may change when the timeline is updated.
	IsDynamic bool
	// Whether this is a live stream.
	IsLive bool
	// Whether this window contains placeholder information because the real information has yet to
	// be loaded.
	IsPlaceholder bool
	// Default position relative to the start of the window at which to begin playback, in microseconds.
	// May be C.TIME_UNSET if and only if the window was populated with a non-zero default position projection, and if the specified projection cannot be performed whilst remaining within the bounds of the window.
	DefaultPosition time.Duration
	// Duration of the window in milliseconds, or C.TIME_UNSET if unknown.
	Duration time.Duration
}

type PositionInfo struct {
	Position time.Duration
}

type DeviceInfo struct {
	MinVolume int
	MaxVolume int
}

type PlaybackParameters struct {
	Speed int
	Pitch int
}

type PlaybackState int

const (
	PLAYBACK_STATE_IDLE      PlaybackState = 1
	PLAYBACK_STATE_BUFFERING PlaybackState = 2
	PLAYBACK_STATE_READY     PlaybackState = 3
	PLAYBACK_STATE_ENDED     PlaybackState = 4
)

type MediaMetadata struct {
	Title           string
	Artist          string
	AlbumTitle      string
	AlbumArtist     string
	DisplayTitle    string
	Subtitle        string
	Description     string
	UserRating      Rating
	OverallRating   Rating
	ArtworkData     []byte
	ArtworkDataType PictureType
	ArtworkUri      string
	TrackNumber     int
	TotalTrackCount int
	IsBrowsable     bool
	IsPlayable      bool
	RecordingYear   int
	RecordingMonth  int
	RecordingDay    int
	ReleaseYear     int
	ReleaseMonth    int
	ReleaseDay      int
	Writer          string
	Composer        string
	Conductor       string
	DiscNumber      int
	TotalDiscCount  int
	Genre           string
	Compilation     string
	Station         string
	MediaType       MediaType
}

type MediaType int

type Rating struct {
	IsRated    bool
	RatingType RatingType
}

type RatingType int

const (
	RATING_TYPE_UNSET RatingType = iota - 1
	RATING_TYPE_HEART
	RATING_TYPE_PERCENTAGE
	RATING_TYPE_STAR
	RATING_TYPE_THUMB
)

type PictureType int

const (
	PICTURE_TYPE_OTHER                      PictureType = 0x00
	PICTURE_TYPE_FILE_ICON                  PictureType = 0x01
	PICTURE_TYPE_FILE_ICON_OTHER            PictureType = 0x02
	PICTURE_TYPE_FRONT_COVER                PictureType = 0x03
	PICTURE_TYPE_BACK_COVER                 PictureType = 0x04
	PICTURE_TYPE_LEAFLET_PAGE               PictureType = 0x05
	PICTURE_TYPE_MEDIA                      PictureType = 0x06
	PICTURE_TYPE_LEAD_ARTIST_PERFORMER      PictureType = 0x07
	PICTURE_TYPE_ARTIST_PERFORMER           PictureType = 0x08
	PICTURE_TYPE_CONDUCTOR                  PictureType = 0x09
	PICTURE_TYPE_BAND_ORCHESTRA             PictureType = 0x0A
	PICTURE_TYPE_COMPOSER                   PictureType = 0x0B
	PICTURE_TYPE_LYRICIST                   PictureType = 0x0C
	PICTURE_TYPE_RECORDING_LOCATION         PictureType = 0x0D
	PICTURE_TYPE_DURING_RECORDING           PictureType = 0x0E
	PICTURE_TYPE_DURING_PERFORMANCE         PictureType = 0x0F
	PICTURE_TYPE_MOVIE_VIDEO_SCREEN_CAPTURE PictureType = 0x10
	PICTURE_TYPE_A_BRIGHT_COLORED_FISH      PictureType = 0x11
	PICTURE_TYPE_ILLUSTRATION               PictureType = 0x12
	PICTURE_TYPE_BAND_ARTIST_LOGO           PictureType = 0x13
	PICTURE_TYPE_PUBLISHER_STUDIO_LOGO      PictureType = 0x14
)

// PlaybackError is adapted from androidx.media3.common.PlaybackException
type PlaybackError struct {
	Null        bool
	Code        PlaybackCode
	TimestampMs int64
}

func NewPlaybackError() PlaybackError {
	return PlaybackError{Null: true}
}

type PlaybackCode int

const (
	// Caused by an error whose cause could not be identified.
	PLAYBACK_ERROR_CODE_UNSPECIFIED PlaybackCode = 1000

	// Caused by an unidentified error in a remote Player, which is a Player that runs on a different
	// host or process.
	PLAYBACK_ERROR_CODE_REMOTE_ERROR PlaybackCode = 1001
	// Caused by the loading position falling behind the sliding window of available live content.
	PLAYBACK_ERROR_CODE_BEHIND_LIVE_WINDOW PlaybackCode = 1002
	// Caused by a generic timeout.
	PLAYBACK_ERROR_CODE_TIMEOUT PlaybackCode = 1003
	// Caused by a failed runtime check.
	//
	// This can happen when the application fails to comply with the player's API requirements (for
	// example, by passing invalid arguments), or when the player reaches an invalid state.
	PLAYBACK_ERROR_CODE_FAILED_RUNTIME_CHECK PlaybackCode = 1004

	// Input/Output errors (2xxx).

	/** Caused by an Input/Output error which could not be identified. */
	PLAYBACK_ERROR_CODE_IO_UNSPECIFIED PlaybackCode = 2000
	// Caused by a network connection failure.
	//
	// The following is a non-exhaustive list of possible reasons:
	//
	//
	//   There is no network connectivity (you can check this by querying
	//   ConnectivityManager#getActiveNetwork).
	//   The URL's domain is misspelled or does not exist.
	//   The target host is unreachable.
	//   The server unexpectedly closes the connection.
	PLAYBACK_ERROR_CODE_IO_NETWORK_CONNECTION_FAILED PlaybackCode = 2001
	// Caused by a network timeout, meaning the server is taking too long to fulfill a request.
	PLAYBACK_ERROR_CODE_IO_NETWORK_CONNECTION_TIMEOUT PlaybackCode = 2002
	// Caused by a server returning a resource with an invalid "Content-Type" HTTP header value.
	//
	// For example, this can happen when the player is expecting a piece of media, but the server
	// returns a paywall HTML page, with content type "text/html".
	PLAYBACK_ERROR_CODE_IO_INVALID_HTTP_CONTENT_TYPE PlaybackCode = 2003
	// Caused by an HTTP server returning an unexpected HTTP response status code.
	PLAYBACK_ERROR_CODE_IO_BAD_HTTP_STATUS PlaybackCode = 2004
	// Caused by a non-existent file.
	PLAYBACK_ERROR_CODE_IO_FILE_NOT_FOUND PlaybackCode = 2005
	// Caused by lack of permission to perform an IO operation. For example, lack of permission to
	// access internet or external storage.
	PLAYBACK_ERROR_CODE_IO_NO_PERMISSION PlaybackCode = 2006
	// Caused by the player trying to access cleartext HTTP traffic (meaning http:// rather than
	// https://) when the app's Network Security Configuration does not permit it.
	//
	// See https://developer.android.com/guide/topics/media/issues/cleartext-not-permitted this
	// corresponding troubleshooting topic
	PLAYBACK_ERROR_CODE_IO_CLEARTEXT_NOT_PERMITTED PlaybackCode = 2007
	// Caused by reading data out of the data bound.
	PLAYBACK_ERROR_CODE_IO_READ_POSITION_OUT_OF_RANGE PlaybackCode = 2008

	// Content parsing errors (3xxx).

	// Caused by a parsing error associated with a media container format bitstream.
	PLAYBACK_ERROR_CODE_PARSING_CONTAINER_MALFORMED PlaybackCode = 3001
	// Caused by a parsing error associated with a media manifest. Examples of a media manifest are a
	// DASH or a SmoothStreaming manifest, or an HLS playlist.
	PLAYBACK_ERROR_CODE_PARSING_MANIFEST_MALFORMED PlaybackCode = 3002
	// Caused by attempting to extract a file with an unsupported media container format, or an
	// unsupported media container feature.
	PLAYBACK_ERROR_CODE_PARSING_CONTAINER_UNSUPPORTED PlaybackCode = 3003
	// Caused by an unsupported feature in a media manifest. Examples of a media manifest are a DASH
	// or a SmoothStreaming manifest, or an HLS playlist.
	PLAYBACK_ERROR_CODE_PARSING_MANIFEST_UNSUPPORTED PlaybackCode = 3004

	// Decoding errors (4xxx).

	// Caused by a decoder initialization failure.
	PLAYBACK_ERROR_CODE_DECODER_INIT_FAILED PlaybackCode = 4001
	// Caused by a decoder query failure.
	PLAYBACK_ERROR_CODE_DECODER_QUERY_FAILED PlaybackCode = 4002
	// Caused by a failure while trying to decode media samples.
	PLAYBACK_ERROR_CODE_DECODING_FAILED PlaybackCode = 4003
	// Caused by trying to decode content whose format exceeds the capabilities of the device.
	PLAYBACK_ERROR_CODE_DECODING_FORMAT_EXCEEDS_CAPABILITIES PlaybackCode = 4004
	// Caused by trying to decode content whose format is not supported.
	PLAYBACK_ERROR_CODE_DECODING_FORMAT_UNSUPPORTED PlaybackCode = 4005

	// AudioTrack errors (5xxx).

	// Caused by an AudioTrack initialization failure.
	PLAYBACK_ERROR_CODE_AUDIO_TRACK_INIT_FAILED PlaybackCode = 5001
	// Caused by an AudioTrack write operation failure.
	PLAYBACK_ERROR_CODE_AUDIO_TRACK_WRITE_FAILED PlaybackCode = 5002

	// DRM errors (6xxx).

	// Caused by an unspecified error related to DRM protection.
	PLAYBACK_ERROR_CODE_DRM_UNSPECIFIED PlaybackCode = 6000
	// Caused by a chosen DRM protection scheme not being supported by the device. Examples of DRM
	// protection schemes are ClearKey and Widevine.
	PLAYBACK_ERROR_CODE_DRM_SCHEME_UNSUPPORTED PlaybackCode = 6001
	// Caused by a failure while provisioning the device.
	PLAYBACK_ERROR_CODE_DRM_PROVISIONING_FAILED PlaybackCode = 6002
	// Caused by attempting to play incompatible DRM-protected content.
	//
	// For example, this can happen when attempting to play a DRM protected stream using a scheme
	// (like Widevine) for which there is no corresponding license acquisition data (like a pssh box).
	PLAYBACK_ERROR_CODE_DRM_CONTENT_ERROR PlaybackCode = 6003
	// Caused by a failure while trying to obtain a license.
	PLAYBACK_ERROR_CODE_DRM_LICENSE_ACQUISITION_FAILED PlaybackCode = 6004
	// Caused by an operation being disallowed by a license policy.
	PLAYBACK_ERROR_CODE_DRM_DISALLOWED_OPERATION PlaybackCode = 6005
	// Caused by an error in the DRM system.
	PLAYBACK_ERROR_CODE_DRM_SYSTEM_ERROR PlaybackCode = 6006
	// Caused by the device having revoked DRM privileges.
	PLAYBACK_ERROR_CODE_DRM_DEVICE_REVOKED PlaybackCode = 6007
	// Caused by an expired DRM license being loaded into an open DRM session.
	PLAYBACK_ERROR_CODE_DRM_LICENSE_EXPIRED PlaybackCode = 6008

	// Frame processing errors (7xxx).

	// Caused by a failure when initializing a VideoFrameProcessor. @UnstableApi
	PLAYBACK_ERROR_CODE_VIDEO_FRAME_PROCESSOR_INIT_FAILED PlaybackCode = 7000
	// Caused by a failure when processing a video frame. @UnstableApi
	PLAYBACK_ERROR_CODE_VIDEO_FRAME_PROCESSING_FAILED PlaybackCode = 7001

	// Player implementations that want to surface custom errors can use error codes greater than this
	// value, so as to avoid collision with other error codes defined in this class.
	PLAYBACK_CUSTOM_ERROR_CODE_BASE PlaybackCode = 1000000
)

func (e PlaybackError) Error() string {
	return e.String()
}

func (e PlaybackError) String() string {
	switch e.Code {
	case PLAYBACK_ERROR_CODE_UNSPECIFIED:
		return "ERROR_CODE_UNSPECIFIED"
	case PLAYBACK_ERROR_CODE_REMOTE_ERROR:
		return "ERROR_CODE_REMOTE_ERROR"
	case PLAYBACK_ERROR_CODE_BEHIND_LIVE_WINDOW:
		return "ERROR_CODE_BEHIND_LIVE_WINDOW"
	case PLAYBACK_ERROR_CODE_TIMEOUT:
		return "ERROR_CODE_TIMEOUT"
	case PLAYBACK_ERROR_CODE_FAILED_RUNTIME_CHECK:
		return "ERROR_CODE_FAILED_RUNTIME_CHECK"
	case PLAYBACK_ERROR_CODE_IO_UNSPECIFIED:
		return "ERROR_CODE_IO_UNSPECIFIED"
	case PLAYBACK_ERROR_CODE_IO_NETWORK_CONNECTION_FAILED:
		return "ERROR_CODE_IO_NETWORK_CONNECTION_FAILED"
	case PLAYBACK_ERROR_CODE_IO_NETWORK_CONNECTION_TIMEOUT:
		return "ERROR_CODE_IO_NETWORK_CONNECTION_TIMEOUT"
	case PLAYBACK_ERROR_CODE_IO_INVALID_HTTP_CONTENT_TYPE:
		return "ERROR_CODE_IO_INVALID_HTTP_CONTENT_TYPE"
	case PLAYBACK_ERROR_CODE_IO_BAD_HTTP_STATUS:
		return "ERROR_CODE_IO_BAD_HTTP_STATUS"
	case PLAYBACK_ERROR_CODE_IO_FILE_NOT_FOUND:
		return "ERROR_CODE_IO_FILE_NOT_FOUND"
	case PLAYBACK_ERROR_CODE_IO_NO_PERMISSION:
		return "ERROR_CODE_IO_NO_PERMISSION"
	case PLAYBACK_ERROR_CODE_IO_CLEARTEXT_NOT_PERMITTED:
		return "ERROR_CODE_IO_CLEARTEXT_NOT_PERMITTED"
	case PLAYBACK_ERROR_CODE_IO_READ_POSITION_OUT_OF_RANGE:
		return "ERROR_CODE_IO_READ_POSITION_OUT_OF_RANGE"
	case PLAYBACK_ERROR_CODE_PARSING_CONTAINER_MALFORMED:
		return "ERROR_CODE_PARSING_CONTAINER_MALFORMED"
	case PLAYBACK_ERROR_CODE_PARSING_MANIFEST_MALFORMED:
		return "ERROR_CODE_PARSING_MANIFEST_MALFORMED"
	case PLAYBACK_ERROR_CODE_PARSING_CONTAINER_UNSUPPORTED:
		return "ERROR_CODE_PARSING_CONTAINER_UNSUPPORTED"
	case PLAYBACK_ERROR_CODE_PARSING_MANIFEST_UNSUPPORTED:
		return "ERROR_CODE_PARSING_MANIFEST_UNSUPPORTED"
	case PLAYBACK_ERROR_CODE_DECODER_INIT_FAILED:
		return "ERROR_CODE_DECODER_INIT_FAILED"
	case PLAYBACK_ERROR_CODE_DECODER_QUERY_FAILED:
		return "ERROR_CODE_DECODER_QUERY_FAILED"
	case PLAYBACK_ERROR_CODE_DECODING_FAILED:
		return "ERROR_CODE_DECODING_FAILED"
	case PLAYBACK_ERROR_CODE_DECODING_FORMAT_EXCEEDS_CAPABILITIES:
		return "ERROR_CODE_DECODING_FORMAT_EXCEEDS_CAPABILITIES"
	case PLAYBACK_ERROR_CODE_DECODING_FORMAT_UNSUPPORTED:
		return "ERROR_CODE_DECODING_FORMAT_UNSUPPORTED"
	case PLAYBACK_ERROR_CODE_AUDIO_TRACK_INIT_FAILED:
		return "ERROR_CODE_AUDIO_TRACK_INIT_FAILED"
	case PLAYBACK_ERROR_CODE_AUDIO_TRACK_WRITE_FAILED:
		return "ERROR_CODE_AUDIO_TRACK_WRITE_FAILED"
	case PLAYBACK_ERROR_CODE_DRM_UNSPECIFIED:
		return "ERROR_CODE_DRM_UNSPECIFIED"
	case PLAYBACK_ERROR_CODE_DRM_SCHEME_UNSUPPORTED:
		return "ERROR_CODE_DRM_SCHEME_UNSUPPORTED"
	case PLAYBACK_ERROR_CODE_DRM_PROVISIONING_FAILED:
		return "ERROR_CODE_DRM_PROVISIONING_FAILED"
	case PLAYBACK_ERROR_CODE_DRM_CONTENT_ERROR:
		return "ERROR_CODE_DRM_CONTENT_ERROR"
	case PLAYBACK_ERROR_CODE_DRM_LICENSE_ACQUISITION_FAILED:
		return "ERROR_CODE_DRM_LICENSE_ACQUISITION_FAILED"
	case PLAYBACK_ERROR_CODE_DRM_DISALLOWED_OPERATION:
		return "ERROR_CODE_DRM_DISALLOWED_OPERATION"
	case PLAYBACK_ERROR_CODE_DRM_SYSTEM_ERROR:
		return "ERROR_CODE_DRM_SYSTEM_ERROR"
	case PLAYBACK_ERROR_CODE_DRM_DEVICE_REVOKED:
		return "ERROR_CODE_DRM_DEVICE_REVOKED"
	case PLAYBACK_ERROR_CODE_DRM_LICENSE_EXPIRED:
		return "ERROR_CODE_DRM_LICENSE_EXPIRED"
	case PLAYBACK_ERROR_CODE_VIDEO_FRAME_PROCESSOR_INIT_FAILED:
		return "ERROR_CODE_VIDEO_FRAME_PROCESSOR_INIT_FAILED"
	case PLAYBACK_ERROR_CODE_VIDEO_FRAME_PROCESSING_FAILED:
		return "ERROR_CODE_VIDEO_FRAME_PROCESSING_FAILED"
	default:
		if e.Code >= PLAYBACK_CUSTOM_ERROR_CODE_BASE {
			return "custom error code"
		} else {
			return "invalid error code"
		}
	}
}
