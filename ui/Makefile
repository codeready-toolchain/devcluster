
## build - installs deps and builds the application into build/.
build:
	@yarn install
	@yarn build

## deploy - copies the static resources into the parent project's static folder.
deploy:
	@cp -r build/* ../pkg/assets/

## clean - cleans the build directory.
clean:
	@rm -rf node_modules

.PHONY: help clean
all: help
help: Makefile
	@echo
	@echo " Choose a command run for "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo