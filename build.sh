#!/bin/sh

rm -f dist/gblog_linux_amd64  dist/gblog_windows_amd64.exe dist/*.xz

packr2 
gox -tags="embed" -osarch="linux/amd64" -output="dist/{{.Dir}}_{{.OS}}_{{.Arch}}" && upx -9 dist/gblog_linux_amd64
gox -tags="embed" -osarch="windows/amd64" -output="dist/{{.Dir}}_{{.OS}}_{{.Arch}}" && upx -9 dist/gblog_windows_amd64.exe

cd dist

tar Jcf gblog_linux_amd64.tar.xz  gblog_linux_amd64 articles static config.toml
tar Jcf gblog_windows_amd64.tar.xz  gblog_windows_amd64.exe articles static config.toml