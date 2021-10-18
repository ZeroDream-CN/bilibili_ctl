#!/bin/sh

mkdir build
mkdir release
../rsrc_linux_amd64 -manifest nac.manifest -o nac.syso

echo "正在构建 Windows x32"
export GOOS=windows
export GOARCH=386
cp nac.syso build
go build -o build/Bilibili_Ctl_windows_i686.exe

echo "正在构建 Windows x64"
export GOOS=windows
export GOARCH=amd64
rm -rf nac.syso
rm -rf build/nac.syso
../rsrc_linux_amd64 -manifest nac64.manifest -o nac.syso
cp nac.syso build
go build -o build/Bilibili_Ctl_windows_amd64.exe

echo "正在构建 Linux x32"
export GOOS=linux
export GOARCH=386
rm -rf nac.syso
rm -rf build/nac.syso
go build -o build/Bilibili_Ctl_linux_i686

echo "正在构建 Linux x64"
export GOOS=linux
export GOARCH=amd64
go build -o build/Bilibili_Ctl_linux_amd64

echo "构建完成，正在打包 Windows 版本..."
zip release/Bilibili_Ctl_windows.zip build/Bilibili_Ctl_windows_amd64.exe build/Bilibili_Ctl_windows_i686.exe

echo "构建完成，正在打包 Linux 版本..."
zip release/Bilibili_Ctl_linux.zip build/Bilibili_Ctl_linux_amd64 build/Bilibili_Ctl_linux_i686

echo "打包完成，正在清理目录..."

rm -rf build/

echo "全部构建已完成，构建的文件已保存至 release/ 目录下"