# tube [![GoDoc](https://godoc.org/github.com/wybiral/tube?status.svg)](https://godoc.org/github.com/wybiral/tube)

Tube is a Golang project to build a self hosted "tube"-style video player for watching your own video collection over HTTP or hosting your own channel for others to watch.

Some of the key features include:
- Easy to add videos (just move a file into the folder)
- No database (video info pulled from file metadata)
- Can be used without JavaScript Enabled
- Easy to customize CSS and HTML template
- Automatically generates RSS feed (at `/feed.xml`)
- Builtin Tor onion service support
- Clean, simple, familiar UI

Currently only supports MP4 video files so you may need to re-encode your media to MP4 using something like [ffmpeg](https://ffmpeg.org/).

Since all of the video info comes from metadata it's also useful to have a metadata editor such as [EasyTAG](https://github.com/GNOME/easytag) (which supports attaching images as thumbnails too).

By default the server is configured to run on 127.0.0.1:0 which will assign a random port every time you run it. This is to avoid conflicting with other applications and to ensure privacy. You can configure this to be any specific host:port by editing `config.json` before running the server. You can also change the RSS feed details and library path from `config.json`.

## About this repository

This one is a fork of the original repository tube. I personally like to call it Stube.

This project is rather for myself but you are free to use it (if you ever wanted to). I was looking for a way to host some video files in my local network without having to install annoying stuff on neither side of servers and clients, and I found this simple and good enough. Since I am a beginner in programming I'm also using this as my playground to mess with GitHub and some basic codes.

The main difference between the original project and this one is the video controls based on JavaScript. Although being JS-Free is great for both security and privacy, it makes controls rely on browser native controller which usually doesn't really provide a good user experience. Also it becomes pointless when you want to use it for a local video streaming server rather than for a public channel.

Still, it will not use any third-party script or CDN, to maintain an isolated environment so it can prevent any leak. You also can still play videos without Javscript. Turning off JavaScript will make you able to use browser native controls.

While Stube currently is just a fork repository that adds barely more than some Javascripts, its ultimate goal is to create a complete personal media streaming server while being light-weighted, portable and simple enough. You will be able to access all the media files you own in your Stube (in a million years).

Supported File Formats:
- Video
    - MP4
    - OGG (OGV)
- Audio 
    - MP3
    - FLAC
    - WAV
    - OGG (OGA)
- Playable Files (Metadata are ignored)
    - WEBA
    - WEBM
- Captions / Subtitles
    - VTT

All differences from the original work:
- Has a video controller based on Javascript, for both pointers and keyboards. Enabled when Javascript is on.
- Loads files in subfolders.
- Basic Search Function.
- Video-only page for embeds (iframes).
- Loads WebVTT files.

To load VTT files, they have to be in the same folder with the video, and have the filename in the appropriate format.
```
- videos
    - video1.mp4        - MP4 File
    - video1.webm       - MP4 File
    - video1.en.vtt     - VTT File, works for both files
    - video1.mp4.en.vtt - VTT File, works for video.mp4 file
    - video1.vtt        - VTT File, works for both files, not recommended
```

Todo:
- Home (/) with a sitemap and a better searching
- Books (/b/) ePub/PDF
- Games (/g/) Browser Based Games
- Images (/i/) png/jpeg/gif/webp
- Comics (/c/) Listed Images
- Series (/s/) Listed Videos/Audios

# Installation

## From release

1. Download [release](https://github.com/178619/tube/releases) for your platform
2. Extract zip archive
3. Run `tube` executable to start server (this will output the URL for accessing from a browser)
4. Move videos to `videos` directory
5. Open the URL from step 3 and enjoy!

## From source

1. [Install Golang](https://golang.org/doc/install) if you don't already have it
2. `go get github.com/wybiral/tube`
3. `cd $GOPATH/src/github.com/wybiral/tube`
4. `go run main.go` (this will output the URL for accessing from a browser)
5. Move videos to `$GOPATH/src/github.com/wybiral/tube/videos`
6. Open the URL from step 4 and enjoy!
