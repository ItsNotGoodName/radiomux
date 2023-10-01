package com.gurnain.radiomuxplayer

import android.os.Handler
import android.os.Looper
import android.util.Log
import com.gurnain.radiomuxplayer.protos.Message
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.Response
import okhttp3.WebSocket
import okhttp3.WebSocketListener
import okio.ByteString
import okio.ByteString.Companion.toByteString

class Connection constructor(private val url: String, private val listener: Listener) {
    private var webSocket: WebSocket? = null
    private var shouldReconnect = true

    private fun initWebSocket() {
        val client = OkHttpClient()
        val request = Request.Builder().url(url).build()
        webSocket = client.newWebSocket(request, webSocketListener)
        client.dispatcher.executorService.shutdown()
    }

    fun connect() {
        shouldReconnect = true
        initWebSocket()
    }

    fun reconnect() {
        initWebSocket()
    }

    fun disconnect() {
        webSocket?.close(1000, "Disconnecting...")
        webSocket = null
        shouldReconnect = false
    }

    fun send(event: Message.Event) {
        webSocket?.send(event.toByteArray().toByteString())
    }

    private val webSocketListener = object : WebSocketListener() {
        private val TAG = "webSocketListener"

        override fun onOpen(webSocket: WebSocket, response: Response) {
            Log.v(TAG, "onOpen()")
        }

        override fun onMessage(webSocket: WebSocket, text: String) {
            Log.v(TAG, "onMessage(): text=$text")
        }

        override fun onMessage(webSocket: WebSocket, bytes: ByteString) {
            Log.v(TAG, "onMessage(): bytes")
            Message.Rpc.parseFrom(bytes.toByteArray())?.let {
                Handler(Looper.getMainLooper()).post { listener.rpc(it) }
            }
        }

        override fun onFailure(
            webSocket: WebSocket, t: Throwable, response: Response?
        ) {
            Log.e(TAG, "onFailure(): " + response.toString() + ":" + t.toString())
            if (shouldReconnect) reconnect()
        }

        override fun onClosed(webSocket: WebSocket, code: Int, reason: String) {
            Log.e(TAG, "onClosed: $code: $reason")
        }
    }

    interface Listener {
        fun rpc(payload: Message.Rpc)
    }
}