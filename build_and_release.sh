read -p "Enter version number (e.g. 2.3.4): " version

rm -rf binaries
GOOS=darwin GOARCH=amd64 go build -o "binaries/go-noughts-and-crosses-v$version-darwin-amd64"
GOOS=linux GOARCH=amd64 go build -o "binaries/go-noughts-and-crosses-v$version-linux-amd64"
GOOS=windows GOARCH=amd64 go build -o "binaries/go-noughts-and-crosses-v$version-windows-amd64.exe"

git tag "v$version"
git push origin "v$version"
gh release create "v$version" ./binaries/* --title "v$version" --generate-notes
