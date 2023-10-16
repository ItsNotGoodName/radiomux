package com.gurnain.radiomuxplayer

import android.annotation.SuppressLint
import androidx.media3.common.C
import androidx.media3.exoplayer.ExoPlayer
import androidx.media3.session.MediaSession
import androidx.media3.session.MediaSessionService

class PlaybackService : MediaSessionService() {

    private var mediaSession: MediaSession? = null

    // If desired, validate the controller before returning the media session
    override fun onGetSession(controllerInfo: MediaSession.ControllerInfo): MediaSession? =
        mediaSession

    // Create your player and media session in the onCreate lifecycle event
    @SuppressLint("UnsafeOptInUsageError")
    override fun onCreate() {
        super.onCreate()
        val player = ExoPlayer.Builder(this).setDeviceVolumeControlEnabled(true)
            .setWakeMode(C.WAKE_MODE_NETWORK).build()
        mediaSession = MediaSession.Builder(this, player).build()
    }

    // Remember to release the player and media session in onDestroy
    override fun onDestroy() {
        mediaSession?.run {
            player.release()
            release()
            mediaSession = null
        }
        super.onDestroy()
    }
}
