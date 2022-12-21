GOOS=darwin GOARCH=amd64 go build -o binaries/go-noughts-and-crosses-darwin-amd64
GOOS=linux GOARCH=amd64 go build -o binaries/go-noughts-and-crosses-linux-amd64
GOOS=windows GOARCH=amd64 go build -o binaries/go-noughts-and-crosses-windows-amd64.exe

read -p "Enter version number (e.g. 2.3.4): " version
gh release create "v$version" binaries/go-noughts-and-crosses-darwin-amd64 binaries/go-noughts-and-crosses-linux-amd64 binaries/go-noughts-and-crosses-windows-amd64.exe --title "v$version" --generate-notes