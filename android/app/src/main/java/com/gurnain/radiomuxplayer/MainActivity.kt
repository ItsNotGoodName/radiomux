package com.gurnain.radiomuxplayer

import android.Manifest
import android.content.Intent
import android.os.Build
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.result.ActivityResultLauncher
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.CheckCircle
import androidx.compose.material.icons.filled.Close
import androidx.compose.material.icons.filled.MoreVert
import androidx.compose.material.icons.filled.Refresh
import androidx.compose.material3.DropdownMenu
import androidx.compose.material3.DropdownMenuItem
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.material3.TopAppBar
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import androidx.core.app.ActivityCompat
import androidx.preference.PreferenceManager
import com.gurnain.radiomuxplayer.ui.theme.RadioMuxPlayerTheme
import com.journeyapps.barcodescanner.ScanContract
import com.journeyapps.barcodescanner.ScanOptions

class MainActivity : ComponentActivity() {

    @OptIn(ExperimentalMaterial3Api::class)
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            RadioMuxPlayerTheme {
                // A surface container using the 'background' color from the theme
                Surface(
                    modifier = Modifier.fillMaxSize(), color = MaterialTheme.colorScheme.background
                ) {
                    val context = LocalContext.current
                    Scaffold(
                        topBar = {
                            TopAppBar(title = {
                                Text(
                                    stringResource(id = R.string.app_name),
                                    maxLines = 1,
                                    overflow = TextOverflow.Ellipsis
                                )
                            }, actions = {
                                var expanded by remember { mutableStateOf(false) }
                                val websocketState = ""
                                when (websocketState) {
                                    "connected" -> {
                                        Icon(
                                            imageVector = Icons.Filled.CheckCircle,
                                            contentDescription = "Connected",
                                        )
                                    }

                                    "connecting" -> {
                                        Icon(
                                            imageVector = Icons.Filled.Refresh,
                                            contentDescription = "Connecting",
                                        )
                                    }

                                    else -> {
                                        Icon(
                                            imageVector = Icons.Filled.Close,
                                            contentDescription = "Disconnected",
                                        )
                                    }
                                }
                                Box(contentAlignment = Alignment.Center) {
                                    IconButton(
                                        onClick = { expanded = true }
                                    ) {
                                        Icon(
                                            imageVector = Icons.Filled.MoreVert,
                                            contentDescription = "More Options",
                                        )
                                    }
                                    DropdownMenu(
                                        modifier = Modifier.width(196.dp),
                                        expanded = expanded,
                                        onDismissRequest = { expanded = false }) {
                                        DropdownMenuItem(text = { Text("Settings") }, onClick = {
                                            context.startActivity(
                                                Intent(
                                                    context,
                                                    SettingsActivity::class.java
                                                )
                                            )
                                        })
                                        DropdownMenuItem(text = { Text("Scan") }, onClick = {
                                            onScan(0)
                                        })
                                        DropdownMenuItem(text = { Text("Scan Front") }, onClick = {
                                            onScan(1)
                                        })
                                        DropdownMenuItem(text = { Text("Start") }, onClick = {
                                            Intent(
                                                applicationContext,
                                                ManagerService::class.java
                                            ).let {
                                                it.action = ManagerService.Actions.START.toString()
                                                startService(it)
                                            }
                                        })
                                        DropdownMenuItem(text = { Text("Stop") }, onClick = {
                                            Intent(
                                                applicationContext,
                                                ManagerService::class.java
                                            ).let {
                                                it.action = ManagerService.Actions.STOP.toString()
                                                startService(it)
                                            }
                                        })
                                    }
                                }
                            })
                        },
                        content = { innerPadding ->
                            LazyColumn(
                                contentPadding = innerPadding,
                                verticalArrangement = Arrangement.spacedBy(8.dp)
                            ) {}
                        })

                }
            }
        }
    }

    override fun onStart() {
        super.onStart()

        // I don't know if I need this for foreground services to not be killed...
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.TIRAMISU) {
            ActivityCompat.requestPermissions(
                this,
                arrayOf(Manifest.permission.POST_NOTIFICATIONS),
                0
            )
        }

        // Start manager service unconditionally because I don't have a way to remember last service state
        Intent(
            applicationContext,
            ManagerService::class.java
        ).also {
            it.action = ManagerService.Actions.START.toString()
            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
                startForegroundService(it)
            } else {
                startService(it)
            }
        }

        // Allow scanning to QR code from web ui to configure the app
        barcodeLauncher = registerForActivityResult(ScanContract()) { result ->
            result.contents?.let {
                PreferenceManager.getDefaultSharedPreferences(this /* Activity context */)
                    .edit().putString("url", it).apply()
            }
        }
    }

    private var barcodeLauncher: ActivityResultLauncher<ScanOptions>? = null

    private fun onScan(cameraId: Int) {
        barcodeLauncher?.let {
            val options = ScanOptions()
            options.setOrientationLocked(false)
            options.setCameraId(cameraId)
            it.launch(options)
        }
    }
}
