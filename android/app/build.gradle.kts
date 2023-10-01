plugins {
    id("com.android.application")
    id("org.jetbrains.kotlin.android")
}

android {
    namespace = "com.gurnain.radiomuxplayer"
    compileSdk = 33

    defaultConfig {
        applicationId = "com.gurnain.radiomuxplayer"
        minSdk = 24
        targetSdk = 33
        versionCode = 1
        versionName = "1.0"

        testInstrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"
        vectorDrawables {
            useSupportLibrary = true
        }
    }

    buildTypes {
        release {
            isMinifyEnabled = false
            proguardFiles(
                getDefaultProguardFile("proguard-android-optimize.txt"),
                "proguard-rules.pro"
            )
        }
    }
    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_1_8
        targetCompatibility = JavaVersion.VERSION_1_8
    }
    kotlinOptions {
        jvmTarget = "1.8"
    }
    buildFeatures {
        compose = true
    }
    composeOptions {
        kotlinCompilerExtensionVersion = "1.4.3"
    }
    packaging {
        resources {
            excludes += "/META-INF/{AL2.0,LGPL2.1}"
        }
    }
}

dependencies {

    implementation("androidx.core:core-ktx:1.9.0")
    implementation("androidx.lifecycle:lifecycle-runtime-ktx:2.6.2")
    implementation("androidx.activity:activity-compose:1.7.2")
    implementation(platform("androidx.compose:compose-bom:2023.03.00"))
    implementation("androidx.compose.ui:ui")
    implementation("androidx.compose.ui:ui-graphics")
    implementation("androidx.compose.ui:ui-tooling-preview")
    implementation("androidx.compose.material3:material3")
    testImplementation("junit:junit:4.13.2")
    androidTestImplementation("androidx.test.ext:junit:1.1.5")
    androidTestImplementation("androidx.test.espresso:espresso-core:3.5.1")
    androidTestImplementation(platform("androidx.compose:compose-bom:2023.03.00"))
    androidTestImplementation("androidx.compose.ui:ui-test-junit4")
    debugImplementation("androidx.compose.ui:ui-tooling")
    debugImplementation("androidx.compose.ui:ui-test-manifest")

    val media3_version = "1.1.1"

    // For media playback using ExoPlayer
    implementation("androidx.media3:media3-exoplayer:$media3_version")

    // For DASH playback support with ExoPlayer
    implementation("androidx.media3:media3-exoplayer-dash:$media3_version")
    // For HLS playback support with ExoPlayer
    implementation("androidx.media3:media3-exoplayer-hls:$media3_version")
    // For RTSP playback support with ExoPlayer
    implementation("androidx.media3:media3-exoplayer-rtsp:$media3_version")
    // // For ad insertion using the Interactive Media Ads SDK with ExoPlayer
    // implementation("androidx.media3:media3-exoplayer-ima:$media3_version")

    // // For loading data using the Cronet network stack
    // implementation("androidx.media3:media3-datasource-cronet:$media3_version")
    // // For loading data using the OkHttp network stack
    // implementation("androidx.media3:media3-datasource-okhttp:$media3_version")
    // // For loading data using librtmp
    // implementation("androidx.media3:media3-datasource-rtmp:$media3_version")

    // For building media playback UIs
    implementation("androidx.media3:media3-ui:$media3_version")
    // // For building media playback UIs for Android TV using the Jetpack Leanback library
    // implementation("androidx.media3:media3-ui-leanback:$media3_version")

    // For exposing and controlling media sessions
    implementation("androidx.media3:media3-session:$media3_version")

    // // For extracting data from media containers
    // implementation("androidx.media3:media3-extractor:$media3_version")

    // // For integrating with Cast
    // implementation("androidx.media3:media3-cast:$media3_version")

    // // For scheduling background operations using Jetpack Work's WorkManager with ExoPlayer
    // implementation("androidx.media3:media3-exoplayer-workmanager:$media3_version")

    // // For transforming media files
    // implementation("androidx.media3:media3-transformer:$media3_version")

    // // Utilities for testing media components (including ExoPlayer components)
    // implementation("androidx.media3:media3-test-utils:$media3_version")
    // // Utilities for testing media components (including ExoPlayer components) via Robolectric
    // implementation("androidx.media3:media3-test-utils-robolectric:$media3_version")

    // //// Common functionality for media database components
    // implementation("androidx.media3:media3-database:$media3_version")
    // // Common functionality for media decoders
    // implementation("androidx.media3:media3-decoder:$media3_version")
    // // Common functionality for loading data
    // implementation("androidx.media3:media3-datasource:$media3_version")
    // // Common functionality used across multiple media libraries
    // implementation("androidx.media3:media3-common:$media3_version")

    // HTTP client that provides a WebSocket client
    implementation("com.squareup.okhttp3:okhttp:4.10.0")

    // Protobuf runtime
    implementation("com.google.protobuf:protobuf-javalite:3.24.3")
}