# Tutorials

Recording screens in Go works well but GIFs generation sucks: time to produce an image, quality of the SDK palette make it pretty poor compared to FFMPEG.

- Install the [FFMPEG full build](https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-full.7z) and add the binaries folder in $PATH
- Install [Autohotkey](https://www.autohotkey.com/) for automated commands tutorials and add UX binaries (C:\Program Files\AutoHotkey\UX) in $PATH
- Generate .exe (Windows only)
```
GOOS=windows GOARCH=amd64 go build -o screenrec.exe ./cmd/tutorials
```

Run it, let's go!