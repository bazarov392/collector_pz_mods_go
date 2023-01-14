go build -o ./build/builderMods
GOOS=windows GOARCH=386 go build -o ./build/builderMods.exe
GOOS=windows GOARCH=amd64 go build -o ./build/builderMods64.exe
echo "builderMods - Для Linux"
echo "builderMods.exe - Для Windows x32-bit"
echo "builderMods64.exe - Для Windows x64-bit"