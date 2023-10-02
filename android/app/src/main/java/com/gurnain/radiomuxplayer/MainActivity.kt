package com.gurnain.radiomuxplayer

import android.content.ComponentName
import android.content.Intent
import android.os.Bundle
import android.util.Log
import android.view.WindowManager
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.MoreVert
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.material3.TopAppBar
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.media3.common.C
import androidx.media3.common.DeviceInfo
import androidx.media3.common.MediaItem
import androidx.media3.common.MediaMetadata
import androidx.media3.common.PlaybackException
import androidx.media3.common.PlaybackParameters
import androidx.media3.common.Player
import androidx.media3.session.MediaController
import androidx.media3.session.SessionToken
import androidx.preference.PreferenceManager
import com.google.common.util.concurrent.MoreExecutors
import com.gurnain.radiomuxplayer.protos.Message
import com.gurnain.radiomuxplayer.protos.MessageConverter
import com.gurnain.radiomuxplayer.ui.theme.RadioMuxPlayerTheme

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            RadioMuxPlayerTheme {
                // A surface container using the 'background' color from the theme
                Surface(
                    modifier = Modifier.fillMaxSize(), color = MaterialTheme.colorScheme.background
                ) {
                    App()
                }
            }
        }
    }

    override fun onStart() {
        super.onStart()

        val sharedPreferences =
            PreferenceManager.getDefaultSharedPreferences(this /* Activity context */)
        val url = sharedPreferences.getString("url", "")

        // TODO: remove this when background play is setup
        window.addFlags(WindowManager.LayoutParams.FLAG_KEEP_SCREEN_ON)

        val sessionToken = SessionToken(this, ComponentName(this, PlaybackService::class.java))
        val mediaControllerFuture = MediaController.Builder(this, sessionToken).buildAsync()
        mediaControllerFuture.addListener({
            mediaController = mediaControllerFuture.get()
            url?.let { connection = Connection(it, connectionListener) }
        }, MoreExecutors.directExecutor())
    }

    override fun onStop() {
        connection?.disconnect()

        mediaController?.let {
            it.removeListener(playerListener)
            it.stop()
            it.release()
        }

        super.onStop()
    }

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

    private val connectionListener = object : Connection.Listener {
        private val TAG = "connectionListener"

        override fun rpc(payload: Message.Rpc) {
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
                        mediaController?.let { m ->
                            playerListener.onMediaMetadataChanged(m.mediaMetadata)
                            playerListener.onPlaylistMetadataChanged(m.playlistMetadata)
                            playerListener.onIsLoadingChanged(m.isLoading)
                            playerListener.onPlaybackStateChanged(m.playbackState)
                            playerListener.onIsPlayingChanged(m.isPlaying)
                            m.playerError?.let { playerListener.onPlayerError(it) }
                            playerListener.onPlaybackParametersChanged(m.playbackParameters)
                            playerListener.onVolumeChanged(m.volume)
                            playerListener.onDeviceInfoChanged(m.deviceInfo)
                            playerListener.onDeviceVolumeChanged(m.deviceVolume, m.isDeviceMuted)
                            playerListener.onCurrentUriChanged(currentMediaUri)
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
            MessageConverter.rpcReply(payload.id).let { connection?.send(it) }
        }
    }

    private val playerListener = object : Player.Listener {
        private val TAG = "playerListener"

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

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun App() {
    val context = LocalContext.current
    Scaffold(topBar = {
        TopAppBar(title = {
            Text(
                stringResource(id = R.string.app_name),
                maxLines = 1,
                overflow = TextOverflow.Ellipsis
            )
        }, actions = {
            IconButton(onClick = {
                context.startActivity(
                    Intent(
                        context, SettingsActivity::class.java
                    )
                )
            }) {
                Icon(
                    imageVector = Icons.Filled.MoreVert, contentDescription = "More Options"
                )
            }
        })
    }, content = { innerPadding ->
        LazyColumn(
            contentPadding = innerPadding, verticalArrangement = Arrangement.spacedBy(8.dp)
        ) { }
    })
}

@Preview(showBackground = true)
@Composable
fun GreetingPreview() {
    RadioMuxPlayerTheme {
        App()
    }
}