.PHONY: test
test:
	@(go test ./... -v -cover)

.PHONY: copytest
copytest:
	@(go run main.go cp -p -f sessionA:tools/rundeck -t sessionB:teste14/tools)

.PHONY: clean
clean:
	@(rm -rf /tmp/gitlabctl;go run main.go rm group -f sessionB:teste14/tools)
