# chat - simplest web chat

> Simple web chat. Studying purpose project.
> GeekBrains Go #211. 3rd quarter

## architecture

aka clean architecture

Структура папок резко отличается от [golang-standards/project-layout](https://github.com/golang-standards/project-layout) их структура описана ниже.

### domain

Содержит описание базовых сущностей и бизнес правила их использования  
Для чата таковыми являются пользователь, который пишет и читает сообщения - не важно где, как и на чем он это делает. Например на бумажке пишет ручкой и приклеивает на доску, а другой читает пишет в ответ и клеит снизу или же это gRPC клиент подключился и в стриме отправляет сообщения и получает.

### usecases

Содержит модель вариантов использования базовых сущностей  
Есть какой-то чат, в котором производится отправки и чтение сообщений в широковещательном формате

### infra

Слой инфрастуктуры уже описывает каким способом будет осуществляться ввод/вывод информации для вариантов использования  
Еще этот слой часто называют delivery

### repos

В этом слое описываются различные реалиации интерфесов, описанных в domain & usecases
// TODO: @art-frela , описать подробнее после возникновенич самой директории

## build/start

рекомендовано использовать makefile

```makefile
h help:
	@echo "h help 	- this help"
	@echo "build 	- build and the app"
	@echo "run 	- run the app"
	@echo "clean 	- clean app trash"
	@echo "test 	- run all tests"
```
## config

используется пакет [viper](https://github.com/spf13/viper)

```yaml
app: # метаданные приложения
    name: 'chat' # наименование приложения
httpd:
    port: "8000" # http порт, который будет пытаться открыть приложения и принимать на него http запросы
    host: ""  # ip адрес хоста который будет занимать приложение, можно оставить пустым
env: production # тип окружения в котором запускается сервис, production - логи в сокращенном формате
log:
    level: debug # уровень логирования сервиса: debug, info, warn, error
    file: "" # имя файла лога, если пусто или stdout - будет выводить в stdout, если указано имя фацйла, будет писать в него
```

## logging

вывод: stdOUT

используется [zap](https://github.com/uber-go/zap)
