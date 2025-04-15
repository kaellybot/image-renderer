# image-renderer

Scripts to generate images like sets or tutorials, written in Go

## Tutorials

Generate .exe and moves it on windows
```
GOOS=windows GOARCH=amd64 go build -o screenrec.exe ./cmd/tutorials
```
A video of 10s is started. It takes 1min to be generated/saved
