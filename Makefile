.PHONY: bin
bin:
	GOOS=linux GOARCH=amd64 go build -o bin/linux/tfc-policy-sets-migration
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin/tfc-policy-sets-migration
	GOOS=windows GOARCH=amd64 go build -o bin/windows/tfc-policy-sets-migration

	zip -j -m bin/linux-amd64.zip bin/linux/tfc-policy-sets-migration
	zip -j -m bin/darwin-amd64.zip bin/darwin/tfc-policy-sets-migration
	zip -j -m bin/windows-amd64.zip bin/windows/tfc-policy-sets-migration
