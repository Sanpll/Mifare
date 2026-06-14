docker build -t mifare-app .
docker stop mifare 2>$null
docker rm mifare 2>$null
docker run -d --name mifare -p 8888:8888 -v "${pwd}/mifare.db:/app/mifare.db" mifare-app