docker build -t mifare .

docker stop mifare
docker rm mifare 

docker run -d --name mifare -p 8888:8888 -v "$(pwd)/mifare.db:/app/mifare.db" mifare-api
