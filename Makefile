# Имя вашего Docker-образа
IMAGE_NAME=image_server

# Порт, на котором будет запускаться сервис внутри контейнера
PORT=8080

# Команды для сборки, запуска и тестирования
.PHONY: all build run test clean

# Собирает Docker-образ
build:
	docker build -t $(IMAGE_NAME) .

# Запускает сервис в контейнере
run:
	docker run -d -p $(PORT):$(PORT) --name $(IMAGE_NAME) $(IMAGE_NAME) /app/http_server

# Запускает тесты
test: run
	@echo "Ждем 5 секунд для запуска сервиса..."
	@sleep 5
	@python tests.py; \
	RESULT=$$?; \
	if [ $$RESULT -eq 0 ]; then \
		echo "Тесты прошли успешно!"; \
	else \
		echo "Тесты не пройдены!"; \
	fi
	$(MAKE) stop

# Останавливает и удаляет контейнер
stop:
	docker stop $(IMAGE_NAME)
	docker rm $(IMAGE_NAME)

# Удаляет Docker-образ
clean:
	docker rmi $(IMAGE_NAME)
