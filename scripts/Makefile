.PHONY: build
build: ## Builds an image by name dockerfiles in the repository.
	# Use '-' to ignore errors during script execution
	-@$(CURDIR)/clean_container.sh "$(DIR)"
	-@$(CURDIR)/clean_image.sh "$(DIR)"
	@$(CURDIR)/build.sh "$(DIR)"

.PHONY: push
push: ## Push dockerfile image to repository.
	@$(CURDIR)/build.sh "$(DIR)" "--push"

.PHONY: start
start: ## Run a Dockerfile from the list of available docker image directories (ex. DIR=alpine).
	-@$(CURDIR)/clean_container.sh "$(DIR)"
	@$(CURDIR)/start.sh "$(DIR)" --tail

.PHONY: stop
stop: ## Stop containers.
	-@$(CURDIR)/clean_container.sh "$(DIR)"

.PHONY: clean
clean: ## Cleaning unknown and previous container images.
	-@$(CURDIR)/clean_container.sh "$(DIR)"
	-@$(CURDIR)/clean_image.sh "$(DIR)"

#.PHONY: test
#test: ## Runs the tests on the repository.
#	@$(CURDIR)/run.sh "--test" "$(DIR)"

.PHONY: list
list: ## List all available images.
	@$(CURDIR)/list.sh

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'x