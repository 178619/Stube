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

## Stube

Stube is a fork of the original repository tube.

The main difference between the original project and this one is the video controls based on JavaScript. Although being JS-Free is great for both security and privacy, it makes controls rely on browser native controller which usually doesn't really provide a good user experience. Also it becomes pointless when you want to use it for a local video streaming server rather than for a public channel.

Still, it will not use any third-party script or CDN, to maintain an isolated enviroment so it can prevent any leak. You also can still play videos without Javscript. Turning off JavaScript will make you able to use browser native controls.

While Stube currently is just a fork repository that adds barely more than some Javascripts, its ultimate goal is to create a complete personal media streaming server while being light-weighted, portable and simple enough. You will be able to access all the media files you own in your Stube.

Supported Video Files: MP4
Supported Audio Files: MP3, FLAC, OGG

All differences from the original work:
- Has a video controller based on Javascript, for both pointers and keyboards. Enabled when Javascript is on.
- Loads files in subfolders.
- Video-only page for embeds (iframes).

Todo:
- Home (/) with a sitemap and a searcher
- Books (/b/) ePub/PDF
- Games (/g/) Browser Based Games
- Images (/i/) png/jpeg/gif/webp
- Comics (/c/) Listed Images
- Series (/s/) Listed Videos/Audios

# installation

## from release

1. Download [release](https://github.com/wybiral/tube/releases) for your platform
2. Extract zip archive
3. Run `tube` executable to start server (this will output the URL for accessing from a browser)
4. Move videos to `videos` directory
5. Open the URL from step 3 and enjoy!

## from source

1. [Install Golang](https://golang.org/doc/install) if you don't already have it
2. `go get github.com/wybiral/tube`
3. `cd $GOPATH/src/github.com/wybiral/tube`
4. `go run main.go` (this will output the URL for accessing from a browser)
5. Move videos to `$GOPATH/src/github.com/wybiral/tube/videos`
6. Open the URL from step 4 and enjoy!
