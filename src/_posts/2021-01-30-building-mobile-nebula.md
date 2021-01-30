---
title: >-
    Building Mobile Nebula
description: >-
    Getting my hands dirty with Android development.
---

This post is going to be cheating a bit. I want to start working on adding DNS
resolver configuration to the [mobile nebula][mobile_nebula] app (if you don't
know nebula, [check it out][nebula], it's well worth knowing about), but I also
need to write a blog post for this week, so I'm combining the two exercises.
This post will essentially be my notes from my progress on today's task.

(Protip: listen to [this][heilung] while following along to achieve the proper
open-source programming aesthetic.)

The current mobile nebula app works very well, but it is lacking one major
feature: the ability to specify custom DNS resolvers. This is important because
I want to be able to access resources on my nebula network by their hostname,
not their IP. Android does everything in its power to make DNS configuration
impossible, and essentially the only way to actually accomplish this is by
specifying the DNS resolvers within the app. I go into more details about why
Android is broken [here][dns_issue].

## Setup

Before I can make changes to the app I need to make sure I can correctly build
it in the first place, so that's the major task for today. The first step to
doing so is to install the project's dependencies. As described in the
[mobile_nebula][mobile_nebula] README, the dependencies are:

- [`flutter`](https://flutter.dev/docs/get-started/install)
- [`gomobile`](https://godoc.org/golang.org/x/mobile/cmd/gomobile)
- [`android-studio`](https://developer.android.com/studio)
- [Enable NDK](https://developer.android.com/studio/projects/install-ndk)

It should be noted that as of writing I haven't used any of these tools ever,
and have only done a small amount of android programming, probably 7 or 8 years
ago, so I'm going to have to walk the line between figuring out problems on the
fly and not having to completely learning these entire ecosystems; there's only
so many hours in a weekend, after all.

I'm running [Archlinux][arch] so I install android-studio and flutter by
doing:

```bash
yay -Sy android-studio flutter
```

And I install `gomobile`, according to its [documentation][gomobile] via:

```bash
go get golang.org/x/mobile/cmd/gomobile
gomobile init
```

Now I startup android-studio and go through the setup wizard for it. I choose
standard setup because customized setup doesn't actually offer any interesting
options. Next android-studio spends approximately two lifetimes downloading
dependencies while my eyesight goes blurry because I'm drinking my coffee too
fast.

It's annoying that I need to install these dependencies, especially
android-studio, in order to build this project. A future goal of mine is to nix
this whole thing up, and make a build pipeline where you can provide a full
nebula configuration file and it outputs a custom APK file for that specific
config; zero configuration required at runtime. This will be useful for
lazy/non-technical users who want to be part of the nebula network.

Once android-studio starts up I'm not quite done yet: there's still the NDK
which must be enabled. The instructions given by the link in
[mobile_nebula][mobile_nebula]'s README explain doing this pretty well, but it's
important to install the specific version indicated in the mobile_nebula repo
(`21.0.6113669` at time of writing). Only another 1GB of dependency downloading
to go....

While waiting for the NDK to download I run `flutter doctor` to make sure
flutter is working, and it gives me some permissions errors. [This blog
post][flutter_blog] gives some tips on setting up, and after running the
following...

```bash
sudo groupadd flutterusers
sudo gpasswd -a $USER flutterusers
sudo chown -R :flutterusers /opt/flutter
sudo chmod -R g+w /opt/flutter/
newgrp flutterusers
```

... I'm able to run `flutter doctor`. It gives the following output:

```
[✓] Flutter (Channel stable, 1.22.6, on Linux, locale en_US.UTF-8)
 
[!] Android toolchain - develop for Android devices (Android SDK version 30.0.3)
    ✗ Android licenses not accepted.  To resolve this, run: flutter doctor --android-licenses
[!] Android Studio
    ✗ Flutter plugin not installed; this adds Flutter specific functionality.
    ✗ Dart plugin not installed; this adds Dart specific functionality.
[!] Connected device
    ! No devices available

! Doctor found issues in 3 categories.
```

The first issue is easily solved as per the instructions given. The second is
solved by finding the plugin manager in android-studio and installing the
flutter plugin (which installs the dart plugin as a dependency, we call that a
twofer).

After installing the plugin the doctor command still complains about not finding
the plugins, but the above mentioned blog post indicates to me that this is
expected. It's comforting to know that the problems indicated by the doctor may
or may not be real problems.

The [blog post][flutter_blog] also indicates that I need `openjdk-8` installed,
so I do:

```bash
yay -S jdk8-openjdk
```

And use the `archlinux-java` command to confirm that that is indeed the default
version for my shell. The [mobile_nebula][mobile_nebula] helpfully expects an
`env.sh` file to exist in the root, so if openjdk-8 wasn't already the default I
could make it so within that file.

## Build

At this point I think I'm ready to try actually building an APK. Thoughts and
prayers required. I run the following in a terminal, since for some reason the
`Build > Flutter > Build APK` dropdown button in android-studio did nothing.

```
flutter build apk
```

It takes quite a while to run, but in the end it errors with:

```
make: 'mobileNebula.aar' is up to date.
cp: cannot create regular file '../android/app/src/main/libs/mobileNebula.aar': No such file or directory

FAILURE: Build failed with an exception.

* Where:
Build file '/tmp/src/mobile_nebula/android/app/build.gradle' line: 95

* What went wrong:
A problem occurred evaluating project ':app'.
> Process 'command './gen-artifacts.sh'' finished with non-zero exit value 1

* Try:
Run with --stacktrace option to get the stack trace. Run with --info or --debug option to get more log output. Run with --scan to get full insights.

* Get more help at https://help.gradle.org

BUILD FAILED in 1s
Running Gradle task 'bundleRelease'...
Running Gradle task 'bundleRelease'... Done                         1.7s
Gradle task bundleRelease failed with exit code 1
```

I narrow down the problem to the `./gen-artifacts.sh` script in the repo's root,
which takes in either `android` or `ios` as an argument. Running it directly
as `./gen-artifacts.sh android` results in the same error:

```bash
make: 'mobileNebula.aar' is up to date.
cp: cannot create regular file '../android/app/src/main/libs/mobileNebula.aar': No such file or directory
```

So now I gotta figure out wtf that `mobileNebula.aar` file is. The first thing I
note is that not only is that file not there, but the `libs` directory it's
supposed to be present in is also not there. So I suspect that there's a missing
build step somewhere.

I search for the string `mobileNebula.aar` within the project using
[ag][silver_searcher] and find that it's built by `nebula/Makefile` as follows:

```make
mobileNebula.aar: *.go
	gomobile bind -trimpath -v --target=android
```

So that file is made by `gomobile`, good to know! Additionally the file is
actually there in the `nebula` directory, so I suspect there's just a missing
build step to move it into `android/app/src/main/libs`. Via some more `ag`-ing I
find that the code which is supposed to move the `mobileNebula.aar` file is in
the `gen-artifacts.sh` script, but that script doesn't create the `libs` folder
as it ought to. I apply the following diff:

```bash
diff --git a/gen-artifacts.sh b/gen-artifacts.sh
index 601ed7b..4f73b4c 100755
--- a/gen-artifacts.sh
+++ b/gen-artifacts.sh
@@ -16,7 +16,7 @@ if [ "$1" = "ios" ]; then
 elif [ "$1" = "android" ]; then
   # Build nebula for android
   make mobileNebula.aar
-  rm -rf ../android/app/src/main/libs/mobileNebula.aar
+  mkdir -p ../android/app/src/main/libs
   cp mobileNebula.aar ../android/app/src/main/libs/mobileNebula.aar

 else
```

(The `rm -rf` isn't necessary, since a) that file is about to be overwritten by
the subsequent `cp` whether or not it's there, and b) it's just deleting a
single file so the `-rf` is an unnecessary risk).

At this point I re-run `flutter build apk` and receive a new error. Progress!

```
A problem occurred evaluating root project 'android'.
> A problem occurred configuring project ':app'.
   > Removing unused resources requires unused code shrinking to be turned on. See http://d.android.com/r/tools/shrink-resources.html for more information.
```

I recall that in the original [mobile_nebula][mobile_nebula] README it mentions
to run the `flutter build` command with the `--no-shrink` option, so I try:

```bash
flutter build apk --no-shrink
```

Finally we really get somewhere. The command takes a very long time to run as it
downloads yet more dependencies (mostly android SDK stuff from the looks of it),
but unfortunately still errors out:

```
Execution failed for task ':app:processReleaseResources'.
> Could not resolve all files for configuration ':app:releaseRuntimeClasspath'.
   > Failed to transform mobileNebula-.aar (:mobileNebula:) to match attributes {artifactType=android-compiled-dependencies-resources, org.gradle.status=integration}.
      > Execution failed for AarResourcesCompilerTransform: /home/mediocregopher/.gradle/caches/transforms-2/files-2.1/735fc805916d942f5311063c106e7363/jetified-mobileNebula.
         > /home/mediocregopher/.gradle/caches/transforms-2/files-2.1/735fc805916d942f5311063c106e7363/jetified-mobileNebula/AndroidManifest.xml
```

Time for more `ag`-ing. I find the file `android/app/build.gradle`, which has
the following block:

```
    implementation (name:'mobileNebula', ext:'aar') {
        exec {
            workingDir '../../'
            environment("ANDROID_NDK_HOME", android.ndkDirectory)
            environment("ANDROID_HOME", android.sdkDirectory)
            commandLine './gen-artifacts.sh', 'android'
        }
    }
```

I never set up the `ANDROID_HOME` or `ANDROID_NDK_HOME` environment variables,
and I suppose that if I'm running the flutter command outside of android-studio
there wouldn't be a way for flutter to know those values, so I try setting them
within my `env.sh`:

```bash
export ANDROID_HOME=~/Android/Sdk
export ANDROID_NDK_HOME=~/Android/Sdk/ndk/21.0.6113669
```

Re-running the build command still results in the same error. But it occurs to
me that I probably had built the `mobileNebula.aar` without those set
previously, so maybe it was built with the wrong NDK version or something. I
tried deleting `nebula/mobileNebula.aar` and try building again. This time...
new errors! Lots of them! Big ones and small ones!

At this point I'm a bit fed up, and want to try a completely fresh build. I back
up my modified `env.sh` and `gen-artifacts.sh` files, delete the `mobile_nebula`
repo, re-clone it, reinstall those files, and try building again. This time just
a single error:

```
Execution failed for task ':app:lintVitalRelease'.
> Could not resolve all artifacts for configuration ':app:debugRuntimeClasspath'.
   > Failed to transform libs.jar to match attributes {artifactType=processed-jar, org.gradle.libraryelements=jar, org.gradle.usage=java-runtime}.
      > Execution failed for JetifyTransform: /tmp/src/mobile_nebula/build/app/intermediates/flutter/debug/libs.jar.
         > Failed to transform '/tmp/src/mobile_nebula/build/app/intermediates/flutter/debug/libs.jar' using Jetifier. Reason: FileNotFoundException, message: /tmp/src/mobile_nebula/build/app/intermediates/flutter/debug/libs.jar (No such file or directory). (Run with --stacktrace for more details.)
           Please file a bug at http://issuetracker.google.com/issues/new?component=460323.
```

So that's cool, apparently there's a bug with flutter and I should file a
support ticket? Well, probably not. It seems that while
`build/app/intermediates/flutter/debug/libs.jar` indeed doesn't exist in the
repo, `build/app/intermediates/flutter/release/libs.jar` _does_, so this appears
to possibly be an issue in declaring which build environment is being used.

After some googling I found [this flutter issue][flutter_issue] related to the
error. Tldr: gradle's not playing nicely with flutter. Downgrading could help,
but apparently building with the `--debug` flag also works. I don't want to
build a release version anyway, so this sits fine with me. I run...

```bash
flutter build apk --no-shrink --debug
```

And would you look at that, I got a result!

```
✓ Built build/app/outputs/flutter-apk/app-debug.apk.
```

## Install

Building was probably the hard part, but I'm not totally out of the woods yet.
Theoretically I could email this apk to my phone or something, but I'd like
something with a faster turnover time; I need `adb`.

I install `adb` via the `android-tools` package:

```bash
yay -S android-tools
```

Before `adb` will work, however, I need to turn on USB debugging on my phone,
which I do by following [this article][usb_debugging]. Once connected I confirm
that `adb` can talk to my phone by doing:

```bash
adb devices
```

And then, finally, I can install the apk:

```
adb install build/app/outputs/flutter-apk/app-debug.apk
```

NOT SO FAST! MORE ERRORS!

```
adb: failed to install build/app/outputs/flutter-apk/app-debug.apk: Failure [INSTALL_FAILED_UPDATE_INCOMPATIBLE: Package net.defined.mobile_nebula signatures do not match previously installed version; ignoring!]
```

I'm guessing this is because I already have the real nebula app installed. I
uninstall it and try again.

AND IT WORKS!!! FUCK YEAH!

```
Performing Streamed Install
Success
```

I can open the nebula app on my phone and it works... fine. There's some
pre-existing networks already installed, which isn't the case for the Play Store
version as far as I can remember, so I suspect those are only there in the
debugging build. Unfortunately the presence of these test networks causes the
app the throw a bunch of errors because it can't contact those networks. Oh well.

The presence of those test networks, in a way, is actually a good thing, as it
means there's probably already a starting point for what I want to do: building
a per-device nebula app with a config preloaded into it.

## Further Steps

Beyond continuing on towards my actual goal of adding DNS resolvers to this app,
there's a couple of other paths I could potentially go down at this point.

* As mentioned, nixify the whole thing. I'm 99% sure the android-studio GUI
  isn't actually needed at all, and I only used it for installing the CMake and
  NDK plugins because I didn't bother to look up how to do it on the CLI.

* Figuring out how to do a proper release build would be great, just for my own
  education. Based on the [flutter issue][flutter_issue] it's possible that all
  that's needed is to downgrade gradle, but maybe that's not so easy.

* Get an android emulator working so that I don't have to install to my phone
  everytime I want to test the app out. I'm not sure if that will also work for
  the VPN aspect of the app, but it will at least help me iterate on UI changes
  faster.

But at this point I'm done for the day, I'll continue on this project some other
time.

[mobile_nebula]: https://github.com/DefinedNet/mobile_nebula
[nebula]: https://slack.engineering/introducing-nebula-the-open-source-global-overlay-network-from-slack/
[dns_issue]: https://github.com/DefinedNet/mobile_nebula/issues/9
[arch]: https://archlinux.org/
[android_wiki]: https://wiki.archlinux.org/index.php/Android#Making_/opt/android-sdk_group-writeable
[heilung]: https://youtu.be/SMJ7pxqk5d4?t=220
[flutter_blog]: https://www.rockyourcode.com/how-to-get-flutter-and-android-working-on-arch-linux/
[gomobile]: https://pkg.go.dev/golang.org/x/mobile/cmd/gomobile
[silver_searcher]: https://github.com/ggreer/the_silver_searcher
[flutter_issue]: https://github.com/flutter/flutter/issues/58247
[usb_debugging]: https://www.droidviews.com/how-to-enable-developer-optionsusb-debugging-mode-on-devices-with-android-4-2-jelly-bean/
