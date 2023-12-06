go-winres make --arch amd64 --out ./cmd/rsrc
fyne bundle -o cmd/bundled.go resources/icon.png
go build -o build/lanty.exe -ldflags -H=windowsgui .\cmd\
copy .\settings.yaml .\build\settings.yaml
del "cmd\*.syso"
