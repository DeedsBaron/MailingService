build:
	@echo "\033[0;32mBuilding binary...\033[m"
	@$(MAKE) -s -C spam

all: build
	@$(MAKE) run -s -C spam

containers:
	docker-compose build
	docker-compose up -d --force-recreate

logs:
	docker-compose logs

clean:
	-docker-compose down
	-docker volume rm $$(docker volume ls -q)
	-rm -rf spam/spam

exec:
	docker exec -it spam bash

status:
	docker ps -a

test:
	@$(MAKE) test -s -C spam

.PHONY: all lib clean fclean re

.DEFAULT_GOAL := all
