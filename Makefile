IMAGE_NAME=image_server

PORT=8080

.PHONY: all build run test clean

build:
	docker build -t $(IMAGE_NAME) .

run:
	docker run -d -p $(PORT):$(PORT) --name $(IMAGE_NAME) $(IMAGE_NAME) /app/http_server


test: run
	@echo
	@sleep 5
	@pytest -v tests.py > result.log; \
	RESULT=$$?; \
	cat result.log; \
	if [ $$RESULT -eq 0 ]; then \
		echo "Тесты прошли успешно!"; \
	else \
		echo "Тесты не пройдены!"; \
	fi
	$(MAKE) stop

stop:
	docker stop $(IMAGE_NAME)
	docker rm $(IMAGE_NAME)


clean:
	docker rmi $(IMAGE_NAME)
