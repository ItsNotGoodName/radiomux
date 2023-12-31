package com.gurnain.radiomuxplayer.protos

import androidx.media3.common.C
import androidx.media3.common.DeviceInfo
import androidx.media3.common.MediaMetadata
import androidx.media3.common.PlaybackException
import androidx.media3.common.PlaybackParameters
import androidx.media3.common.Timeline
import java.text.SimpleDateFormat
import java.util.Date
import java.util.Locale

object MessageConverter {

    fun rpcReply(id: Int): Message.RpcReply.Builder {
        return Message.RpcReply.newBuilder().setId(id)
    }

    fun eventRpcReply(rpcReply: Message.RpcReply.Builder): Message.Event {
        return Message.Event.newBuilder().setRpcReply(rpcReply).build()
    }

    fun mediaMetadata(mediaMetadata: MediaMetadata): Message.MediaMetadata.Builder {
        var res = Message.MediaMetadata.newBuilder()
        mediaMetadata.title?.let { res = res.setTitle(it.toString()) }
        mediaMetadata.artist?.let { res = res.setArtist(it.toString()) }
        mediaMetadata.albumTitle?.let { res = res.setAlbumTitle(it.toString()) }
        mediaMetadata.albumArtist?.let { res = res.setAlbumArtist(it.toString()) }
        mediaMetadata.displayTitle?.let { res = res.setDisplayTitle(it.toString()) }
        mediaMetadata.subtitle?.let { res = res.setSubtitle(it.toString()) }
        mediaMetadata.description?.let { res = res.setDescription(it.toString()) }
        mediaMetadata.userRating?.let {
            res = res.setUserRating(Message.Rating.newBuilder().setIsRated(it.isRated))
        }
        mediaMetadata.overallRating?.let {
            res = res.setOverallRating(Message.Rating.newBuilder().setIsRated(it.isRated))
        }
        // mediaMetadata.artworkData?.let { res = res.setArtworkData(ByteString.copyFrom(it)) }
        mediaMetadata.artworkDataType?.let { res = res.setArtworkDataType(it) }
        mediaMetadata.artworkUri?.let { res = res.setArtworkUri(it.toString()) }
        mediaMetadata.trackNumber?.let { res = res.setTrackNumber(it) }
        mediaMetadata.totalTrackCount?.let { res = res.setTotalTrackCount(it) }
        mediaMetadata.isBrowsable?.let { res = res.setIsBrowsable(it) }
        mediaMetadata.isPlayable?.let { res = res.setIsPlayable(it) }
        mediaMetadata.recordingYear?.let { res = res.setRecordingYear(it) }
        mediaMetadata.recordingMonth?.let { res = res.setRecordingMonth(it) }
        mediaMetadata.recordingDay?.let { res = res.setRecordingDay(it) }
        mediaMetadata.releaseYear?.let { res = res.setReleaseYear(it) }
        mediaMetadata.releaseMonth?.let { res = res.setReleaseMonth(it) }
        mediaMetadata.releaseDay?.let { res = res.setReleaseDay(it) }
        mediaMetadata.writer?.let { res = res.setWriter(it.toString()) }
        mediaMetadata.composer?.let { res = res.setComposer(it.toString()) }
        mediaMetadata.conductor?.let { res = res.setConductor(it.toString()) }
        mediaMetadata.discNumber?.let { res = res.setDiscNumber(it) }
        mediaMetadata.totalDiscCount?.let { res = res.setTotalDiscCount(it) }
        mediaMetadata.genre?.let { res = res.setGenre(it.toString()) }
        mediaMetadata.compilation?.let { res = res.setCompilation(it.toString()) }
        mediaMetadata.station?.let { res = res.setStation(it.toString()) }
        mediaMetadata.mediaType?.let { res = res.setMediaType(it) }
        return res
    }

    fun eventOnMediaMetadataChanged(mediaMetadata: Message.MediaMetadata.Builder): Message.Event {
        return Message.OnMediaMetadataChanged.newBuilder().setMediaMetadata(mediaMetadata)
            .let { Message.Event.newBuilder().setOnMediaMetadataChanged(it).build() }
    }

    fun eventOnPlaylistMetadataChanged(mediaMetadata: Message.MediaMetadata.Builder): Message.Event {
        return Message.OnPlaylistMetadataChanged.newBuilder().setMediaMetadata(mediaMetadata)
            .let { Message.Event.newBuilder().setOnPlaylistMetadataChanged(it).build() }
    }

    fun eventOnIsLoadingChanged(isLoading: Boolean): Message.Event {
        return Message.OnIsLoadingChanged.newBuilder().setIsLoading(isLoading)
            .let { Message.Event.newBuilder().setOnIsLoadingChanged(it).build() }
    }

    fun eventOnPlaybackStateChanged(playbackState: Int): Message.Event {
        return Message.OnPlaybackStateChanged.newBuilder().setPlaybackState(playbackState)
            .let { Message.Event.newBuilder().setOnPlaybackStateChanged(it).build() }
    }

    fun eventOnIsPlayingChanged(isPlaying: Boolean): Message.Event {
        return Message.OnIsPlayingChanged.newBuilder().setIsPlaying(isPlaying)
            .let { Message.Event.newBuilder().setOnIsPlayingChanged(it).build() }
    }

    fun eventOnPlayerError(error: PlaybackException): Message.Event {
        return Message.PlaybackException.newBuilder().setErrorCode(error.errorCode)
            .setTimestampMs(error.timestampMs)
            .let { Message.OnPlayerError.newBuilder().setError(it) }
            .let { Message.Event.newBuilder().setOnPlayerError(it).build() }
    }

    fun eventOnPlaybackParametersChanged(playbackParameters: PlaybackParameters): Message.Event {
        return Message.PlaybackParameters.newBuilder().setPitch(playbackParameters.pitch)
            .setSpeed(playbackParameters.speed)
            .let { Message.OnPlaybackParametersChanged.newBuilder().setPlaybackParameters(it) }
            .let { Message.Event.newBuilder().setOnPlaybackParametersChanged(it).build() }
    }

    fun eventOnVolumeChanged(volume: Float): Message.Event {
        return Message.OnVolumeChanged.newBuilder().setVolume(volume)
            .let { Message.Event.newBuilder().setOnVolumeChanged(it).build() }
    }

    fun eventOnDeviceInfoChanged(deviceInfo: DeviceInfo): Message.Event {
        return Message.DeviceInfo.newBuilder().setMinVolume(deviceInfo.minVolume)
            .setMaxVolume(deviceInfo.maxVolume)
            .let { Message.OnDeviceInfoChanged.newBuilder().setDeviceInfo(it) }
            .let { Message.Event.newBuilder().setOnDeviceInfoChanged(it).build() }
    }

    fun eventOnDeviceVolumeChanged(volume: Int, muted: Boolean): Message.Event {
        return Message.OnDeviceVolumeChanged.newBuilder().setVolume(volume).setMuted(muted)
            .let { Message.Event.newBuilder().setOnDeviceVolumeChanged(it).build() }
    }

    fun eventOnCurrentUriChanged(currentUri: String): Message.Event {
        return Message.OnCurrentUriChanged.newBuilder().setUri(currentUri)
            .let { Message.Event.newBuilder().setOnCurrentUriChanged(it).build() }
    }

    fun eventOnTimelineChanged(window: Timeline.Window): Message.Event {
        return Message.TimelineWindow.newBuilder()
            .setPresentationStartTimeMs(window.presentationStartTimeMs)
            .setWindowStartTimeMs(window.windowStartTimeMs)
            .setElapsedRealtimeEpochOffsetMs(window.elapsedRealtimeEpochOffsetMs)
            .setIsSeekable(window.isSeekable)
            .setIsDynamic(window.isDynamic)
            .setIsPlaceholder(window.isPlaceholder)
            .setDefaultPositionMs(window.defaultPositionMs)
            .setDurationMs(window.durationMs)
            .let { Message.OnTimelineChanged.newBuilder().setWindow(it).setTimeUnset(C.TIME_UNSET) }
            .let { Message.Event.newBuilder().setOnTimelineChanged(it).build() }
    }

    fun positionInfo(positionMs: Long): Message.PositionInfo.Builder {
        return Message.PositionInfo.newBuilder()
            .setPositionMs(positionMs)
    }

    private val time = SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss.SSSXXX", Locale.ROOT)

    fun eventOnPositionChanged(
        oldPosition: Message.PositionInfo.Builder,
        newPosition: Message.PositionInfo.Builder,
        date: Date
    ): Message.Event {
        return Message.OnPositionChanged.newBuilder()
            .setOldPosition(oldPosition)
            .setNewPosition(newPosition)
            .setTimestamp(time.format(date))
            .build()
            .let { Message.Event.newBuilder().setOnPositionChanged(it).build() }
    }
}
