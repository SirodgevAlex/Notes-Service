Привет!
Это инструкция

1. Я прикрутил сюда докер, поэтому запускать так

```bash
docker-compose up --build
```

2. Все, мы все запустили, можно делать сами запросы
3. Запрос для создания пользователя

```bash
curl -i -X POST http://localhost:8080/register \
-H 'Content-Type: application/json' \
-d '{"Email": "sirodgev@yandex.ru", "Password": "Sneiieir1_"}'
```

получим в ответ такой результат

```bash
Date: Mon, 28 Apr 2024 08:13:01 GMT
Content-Length: 64
Content-Type: text/plain; charset=utf-8

{"ID":1,"Email":"sirodgev@yandex.ru","Password":"Sneiieir1_"}
```

4. Запрос для аутентификации

```bash
curl -i -X POST http://localhost:8080/authorize \
-H 'Content-Type: application/json' \
-d '{"Email": "sirodgev@yandex.ru", "Password": "Sneiieir1_"}'
```

получим в ответ такой результат. P S токены разные будут, можно скопировать из терминала токен из результат, потом вставить его в след запрос. Тогда все будет хорошо, объявление создастся

```bash
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 28 Apr 2024 08:29:43 GMT
Content-Length: 146

{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjEsImV4cCI6MTcxNDM3OTY4Mywic3ViIjoiMSJ9.6ZawYuH5jBYtM6nGMEgh2REVr8cCKLSyPJAx5DuXRZo"}
```

5. curl запрос для создания заметки

```bash
curl -i -X POST http://localhost:8080/notes \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjEsImV4cCI6MTcxNDM3OTY4Mywic3ViIjoiMSJ9.6ZawYuH5jBYtM6nGMEgh2REVr8cCKLSyPJAx5DuXRZo" \
-H "Content-Type: application/json" \
-d '{
  "title": "This is a title",
  "text": "This asdsdis a newf note"
}'
```

получим в ответ такой результат

```bash
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 28 Apr 2024 08:35:29 GMT
Content-Length: 86

{"note_id":3,"text":"This asdsdis a newf note","title":"This is a title","user_id":1}
```

6. curl запрос для получения заметки

```bash
curl -i -X GET http://localhost:8080/notes/3 \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjEsImV4cCI6MTcxNDQ3MjM2Nywic3ViIjoiMSJ9.litAzj2OprQnlYIxYboTMPkLKfSI84Pipwaw-wXzL6o"
```

получим в ответ такой результат

```bash
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 29 Apr 2024 11:01:13 GMT
Content-Length: 118

{"id":3,"created_at":"2024-04-29T11:35:29+03:00","author_id":1,"title":"Updated title","text":"Updated note content"}
```

7. Запрос для редактирования заметки

```bash
curl -i -X PATCH http://localhost:8080/notes/3 \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjEsImV4cCI6MTcxNDM3OTY4Mywic3ViIjoiMSJ9.6ZawYuH5jBYtM6nGMEgh2REVr8cCKLSyPJAx5DuXRZo" \
-H "Content-Type: application/json" \
-d '{
  "title": "Updated title",
  "text": "Updated note content"
}'
```

получим в ответ такой результат

```bash
HTTP/1.1 200 OK
Date: Mon, 28 Apr 2024 08:39:26 GMT
Content-Length: 0
```

8. curl запрос для удаления заметки

```bash
curl -i -X DELETE http://localhost:8080/notes/2 \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjEsImV4cCI6MTcxNDQ3MjM2Nywic3ViIjoiMSJ9.litAzj2OprQnlYIxYboTMPkLKfSI84Pipwaw-wXzL6o"
```

получим в ответ такой результат

```bash
HTTP/1.1 200 OK
Date: Mon, 28 Apr 2024 10:39:26 GMT
Content-Length: 0
```

если что-то не заработает, то мой тг @sirodgevalex
