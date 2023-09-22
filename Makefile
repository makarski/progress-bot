.PHONY: build config run

build:
	@docker build -t progress-bot .

config:
	@./interactive_configs.sh

run: config
	$(call run_app)

define run_app
	@docker run progress-bot
endef
