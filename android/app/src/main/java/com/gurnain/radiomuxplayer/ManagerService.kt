package com.gurnain.radiomuxplayer

import android.app.NotificationChannel
import android.app.NotificationManager
import android.app.Service
import android.content.Context
import android.content.Intent
import android.os.Build
import android.os.IBinder
import androidx.core.app.NotificationCompat

class ManagerService : Service() {

    private val TAG = "ManagerService"

    override fun onBind(intent: Intent?): IBinder? {
        return null
    }

    override fun onCreate() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val channel = NotificationChannel(
                "manager_channel",
                "Manager Channel",
                NotificationManager.IMPORTANCE_LOW
            )
            val notificationManager =
                getSystemService(Context.NOTIFICATION_SERVICE) as NotificationManager
            notificationManager.createNotificationChannel(channel)
        }
    }

    private var manager: Manager? = null

    private fun start() {
        if (manager == null) {
            manager = Manager(this)
            val notification = NotificationCompat.Builder(this, "manager_channel")
                .setSmallIcon(R.drawable.ic_launcher_foreground).setContentTitle("Manager Service")
                .setContentText("Started")
                .build()
            startForeground(1, notification)
        }
    }

    private fun stop() {
        manager?.let {
            it.stop()
            manager = null
        }
        stopSelf()
    }

    enum class Actions {
        START, STOP
    }

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        when (intent?.action) {
            Actions.START.toString() -> start()
            Actions.STOP.toString() -> stop()
        }
        return super.onStartCommand(intent, flags, startId)
    }
}