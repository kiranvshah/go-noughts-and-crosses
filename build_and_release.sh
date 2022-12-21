GOOS=darwin GOARCH=amd64 go build -o binaries/go-noughts-and-crosses-darwin-amd64
GOOS=linux GOARCH=amd64 go build -o binaries/go-noughts-and-crosses-linux-amd64
GOOS=windows GOARCH=amd64 go build -o binaries/go-noughts-and-crosses-windows-amd64.exe

read -p "Enter version number (e.g. 2.3.4): " version
git tag "v$version"
git push origin "v$version"
gh release create "v$version" ./binaries/* --title "v$version" --generate-notes
