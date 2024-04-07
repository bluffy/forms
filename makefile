help:  ## show help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[$$()% 0-9a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

lang-export:  ## generate files for translation
	goi18n extract --outdir i18n -format yaml
	goi18n merge  --outdir i18n  -format yaml i18n/active.en.yaml i18n/translate.*.yaml	

lang-merge:  ## import tanslations
	goi18n merge  --outdir i18n  -format yaml i18n/active.*.yaml i18n/translate.*.yaml 

dev-prepare:    ## copy Dev Files from Example
	./scripts/bash/prepare-dev.sh

dev-up: 	## run Dev Daimlertuck
	docker compose --env-file ./docker/dev/.env -f docker/dev/docker-compose.yaml -p goapp-dev build
	docker compose --env-file ./docker/dev/.env -f docker/dev/docker-compose.yaml -p goapp-dev up

dev-down: 	## run Dev Daimlertuck
	docker compose --env-file ./docker/dev/.env -f docker/dev/docker-compose.yaml -p goapp-dev down


