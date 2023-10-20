package com.gurnain.radiomuxplayer

import android.content.ComponentName
import android.content.SharedPreferences
import android.util.Log
import androidx.media3.common.C
import androidx.media3.common.DeviceInfo
import androidx.media3.common.MediaItem
import androidx.media3.common.MediaMetadata
import androidx.media3.common.PlaybackException
import androidx.media3.common.PlaybackParameters
import androidx.media3.common.Player
import androidx.media3.common.Timeline
import androidx.media3.session.MediaController
import androidx.media3.session.SessionToken
import androidx.preference.PreferenceManager
import com.google.common.util.concurrent.MoreExecutors
import com.gurnain.radiomuxplayer.protos.Message
import com.gurnain.radiomuxplayer.protos.MessageConverter
import java.util.Date

class Manager(private val context: android.content.Context) :
    SharedPreferences.OnSharedPreferenceChangeListener {

    val TAG = "Manager"

    private var currentMediaUri: String = ""
        set(value) {
            mediaController?.setMediaItem(MediaItem.fromUri(value))
            playerListener.onCurrentUriChanged(value)
            field = value
        }
    private var mediaController: MediaController? = null
        set(value) {
            value?.addListener(playerListener)
            field = value
        }
    private var connection: Connection? = null
        set(value) {
            value?.connect()
            field = value
        }

    init {
        PreferenceManager.getDefaultSharedPreferences(context)
            .registerOnSharedPreferenceChangeListener(this)

        // Setup mediaController and connection
        val sessionToken =
            SessionToken(context, ComponentName(context, PlaybackService::class.java))
        val mediaControllerFuture = MediaController.Builder(context, sessionToken).buildAsync()
        mediaControllerFuture.addListener({
            mediaController = mediaControllerFuture.get()
            onSharedPreferenceChanged(PreferenceManager.getDefaultSharedPreferences(context), "url")
        }, MoreExecutors.directExecutor())
    }

    fun stop() {
        PreferenceManager.getDefaultSharedPreferences(context)
            .unregisterOnSharedPreferenceChangeListener(this)

        connection?.let {
            it.disconnect()
            connection = null
        }

        mediaController?.let {
            it.removeListener(playerListener)
            it.stop()
            it.release()
            mediaController = null
        }
    }

    override fun onSharedPreferenceChanged(sharedPreferences: SharedPreferences, key: String?) {
        sharedPreferences.getString("url", "")?.let { url ->
            connection?.disconnect()
            connection = Connection(url, connectionListener)
        }
    }

    private val connectionListener = object : Connection.Listener {
        private val TAG = "connectionListener"

        override fun rpc(payload: Message.Rpc): Message.RpcReply.Builder {
            Log.v(TAG, payload.toString())

            payload.payloadCase?.let { payloadCase ->
                when (payloadCase) {
                    Message.Rpc.PayloadCase.STOP -> {
                        mediaController?.stop()
                    }

                    Message.Rpc.PayloadCase.PLAY -> {
                        mediaController?.play()
                    }

                    Message.Rpc.PayloadCase.PAUSE -> {
                        mediaController?.pause()
                    }

                    Message.Rpc.PayloadCase.PREPARE -> {
                        mediaController?.prepare()
                    }

                    Message.Rpc.PayloadCase.SETPLAYWHENREADY -> {
                        mediaController?.playWhenReady = payload.setPlayWhenReady.playWhenReady
                    }

                    Message.Rpc.PayloadCase.SETMEDIAITEM -> {
                        currentMediaUri = payload.setMediaItem.uri
                    }

                    Message.Rpc.PayloadCase.SETVOLUME -> {
                        mediaController?.volume = payload.setVolume.volume
                    }

                    Message.Rpc.PayloadCase.SETDEVICEVOLUME -> {
                        mediaController?.setDeviceVolume(
                            payload.setDeviceVolume.volume, C.VOLUME_FLAG_REMOVE_SOUND_AND_VIBRATE
                        )
                    }

                    Message.Rpc.PayloadCase.INCREASEDEVICEVOLUME -> {
                        mediaController?.increaseDeviceVolume(C.VOLUME_FLAG_REMOVE_SOUND_AND_VIBRATE)
                    }

                    Message.Rpc.PayloadCase.DECREASEDEVICEVOLUME -> {
                        mediaController?.decreaseDeviceVolume(C.VOLUME_FLAG_REMOVE_SOUND_AND_VIBRATE)
                    }

                    Message.Rpc.PayloadCase.SETDEVICEMUTED -> {
                        mediaController?.setDeviceMuted(
                            payload.setDeviceMuted.muted, C.VOLUME_FLAG_REMOVE_SOUND_AND_VIBRATE
                        )
                    }

                    Message.Rpc.PayloadCase.SYNCSTATE -> {
                        mediaController?.let { controller ->
                            playerListener.onMediaMetadataChanged(controller.mediaMetadata)
                            playerListener.onPlaylistMetadataChanged(controller.playlistMetadata)
                            playerListener.onIsLoadingChanged(controller.isLoading)
                            playerListener.onPlaybackStateChanged(controller.playbackState)
                            playerListener.onIsPlayingChanged(controller.isPlaying)
                            controller.playerError?.let { playerListener.onPlayerError(it) }
                            playerListener.onPlaybackParametersChanged(controller.playbackParameters)
                            playerListener.onVolumeChanged(controller.volume)
                            playerListener.onDeviceInfoChanged(controller.deviceInfo)
                            playerListener.onDeviceVolumeChanged(
                                controller.deviceVolume,
                                controller.isDeviceMuted
                            )
                            playerListener.onCurrentUriChanged(currentMediaUri)
                            playerListener.onTimelineChanged(
                                controller.currentTimeline,
                                Player.TIMELINE_CHANGE_REASON_PLAYLIST_CHANGED
                            )
                            playerListener.onPositionChanged(
                                controller.currentPosition,
                                controller.currentPosition
                            )
                        }
                    }

                    Message.Rpc.PayloadCase.SEEKTODEFAULTPOSITION -> {
                        mediaController?.seekToDefaultPosition()
                    }

                    Message.Rpc.PayloadCase.PAYLOAD_NOT_SET -> {
                        Log.v(TAG, "payload not set")
                    }
                }
            }

            return MessageConverter.rpcReply(payload.id)
        }
    }

    private val playerListener = object : Player.Listener {
        private val TAG = "playerListener"

        override fun onTimelineChanged(timeline: Timeline, reason: Int) {
            Log.v(TAG, "onTimelineChanged")
            val index = timeline.getFirstWindowIndex(false)
            if (index == C.INDEX_UNSET) {
                return
            }

            Timeline.Window()
                .also { timeline.getWindow(index, it) }
                .let { MessageConverter.eventOnTimelineChanged(it) }
                .let { connection?.send(it) }

        }

        override fun onPositionDiscontinuity(
            oldPosition: Player.PositionInfo,
            newPosition: Player.PositionInfo,
            reason: Int
        ) {
            onPositionChanged(oldPosition.positionMs, newPosition.positionMs)
        }

        fun onPositionChanged(oldPosition: Long, newPosition: Long) {
            Log.v(TAG, "onPositionChanged")
            MessageConverter.eventOnPositionChanged(
                MessageConverter.positionInfo(oldPosition),
                MessageConverter.positionInfo(newPosition),
                Date()
            ).let { connection?.send(it) }
        }

        override fun onMediaMetadataChanged(mediaMetadata: MediaMetadata) {
            Log.v(TAG, "onMediaMetadataChanged")
            MessageConverter.mediaMetadata(mediaMetadata)
                .let { MessageConverter.eventOnMediaMetadataChanged(it) }
                .let { connection?.send(it) }
        }

        override fun onPlaylistMetadataChanged(mediaMetadata: MediaMetadata) {
            Log.v(TAG, "onPlaylistMetadataChanged")
            MessageConverter.mediaMetadata(mediaMetadata)
                .let { MessageConverter.eventOnPlaylistMetadataChanged(it) }
                .let { connection?.send(it) }
        }

        override fun onIsLoadingChanged(isLoading: Boolean) {
            Log.v(TAG, "onIsLoadingChanged: isLoading=$isLoading")
            MessageConverter.eventOnIsLoadingChanged(isLoading).let { connection?.send(it) }
        }

        override fun onPlaybackStateChanged(playbackState: Int) {
            Log.v(TAG, "onPlaybackStateChanged: playbackState=$playbackState")
            MessageConverter.eventOnPlaybackStateChanged(playbackState).let { connection?.send(it) }
        }

        override fun onIsPlayingChanged(isPlaying: Boolean) {
            Log.v(TAG, "onIsPlayingChanged: isPlaying=$isPlaying")
            MessageConverter.eventOnIsPlayingChanged(isPlaying).let { connection?.send(it) }
        }

        override fun onPlayerError(error: PlaybackException) {
            Log.v(TAG, "onPlayerError: $error")
            MessageConverter.eventOnPlayerError(error).let { connection?.send(it) }
        }

        override fun onPlaybackParametersChanged(playbackParameters: PlaybackParameters) {
            Log.v(
                TAG,
                "onPlaybackParametersChanged: speed=${playbackParameters.speed};pitch=${playbackParameters.pitch}"
            )
            MessageConverter.eventOnPlaybackParametersChanged(playbackParameters)
                .let { connection?.send(it) }
        }

        override fun onVolumeChanged(volume: Float) {
            Log.v(TAG, "onVolumeChanged: volume=$volume")
            MessageConverter.eventOnVolumeChanged(volume).let { connection?.send(it) }
        }

        override fun onDeviceInfoChanged(deviceInfo: DeviceInfo) {
            Log.v(
                TAG,
                "onDeviceInfoChanged: minVolume=${deviceInfo.minVolume};maxVolume=${deviceInfo.maxVolume}"
            )
            MessageConverter.eventOnDeviceInfoChanged(deviceInfo).let { connection?.send(it) }
        }

        override fun onDeviceVolumeChanged(volume: Int, muted: Boolean) {
            Log.v(TAG, "onDeviceVolumeChanged: volume=$volume;muted=$muted")
            MessageConverter.eventOnDeviceVolumeChanged(volume, muted).let { connection?.send(it) }
        }

        fun onCurrentUriChanged(currentUri: String) {
            Log.v(TAG, "onCurrentUriChanged: uri=$currentUri")
            MessageConverter.eventOnCurrentUriChanged(currentUri).let { connection?.send(it) }
        }
    }
}
