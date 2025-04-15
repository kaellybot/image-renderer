# image-renderer

Scripts to generate images like sets or tutorials, written in Go

## Tutorials

- For quality and performances reasons, FFMPEG is used so install the .exe somewhere and add the PATH to the binary in $PATH
- Generate .exe (Windows only)
```
GOOS=windows GOARCH=amd64 go build -o screenrec.exe ./cmd/tutorials
```

Run it, let's go!