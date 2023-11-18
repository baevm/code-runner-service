# code-runner-service

### 1. **Запуск кода**
Создает контейнер с нужным ЯП. Запускает, выполнят и после удаляет.

Запрос:
```
curl localhost:8080/run_code -X POST -H "Content-Type: application/json" \
-d '{"language": "python", "code": "from time import sleep\nsleep(2)\nres = 0\nprint(res)"}'
```

Ответ:
```
{"result":"0\r\n"}
```
