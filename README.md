## Запуск с Docker
1. Установите Docker: `sudo apt install docker.io`
2. Соберите и запустите сервер:
    - `docker build -t chat-server .`
    - `docker run -p 8080:8080 chat-server`
3. (Опционально) Используйте Docker Compose:
    - `docker-compose up --build`

## Запуск клиента с Docker
1. Соберите клиент: `docker build -t chat-client -f Dockerfile.client .`
2. Запустите: `docker run -it --network host chat-client`