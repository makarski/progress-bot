.PHONY: build config run

build:
	@docker build -t progress-bot .

config:
	@./interactive_configs.sh

run: config env
	$(call run_app)

define run_app
	@docker run -it progress-bot
endef
