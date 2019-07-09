BIN_NAME=tfc-policy-sets-migration

.PHONY: bin
bin:
	GOOS=linux GOARCH=amd64 go build -o build/x/linux/$(BIN_NAME)
	GOOS=darwin GOARCH=amd64 go build -o build/x/darwin/$(BIN_NAME)
	GOOS=windows GOARCH=amd64 go build -o build/x/windows/$(BIN_NAME)

	zip -j build/$(BIN_NAME)-linux-amd64.zip build/x/linux/$(BIN_NAME)
	zip -j build/$(BIN_NAME)-darwin-amd64.zip build/x/darwin/$(BIN_NAME)
	zip -j build/$(BIN_NAME)-windows-amd64.zip build/x/windows/$(BIN_NAME)

	rm -rf build/x
