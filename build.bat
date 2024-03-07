go generate ./...
go build -o build/lanty.exe -ldflags -H=windowsgui .\cmd\
copy .\settings.yaml .\build\settings.yaml
del "cmd\*.syso"
