# pdf-microservice
## Cервис, генерящий пдфки

### Запуск:

При первом запуске:

```shell
git clone github.com/Basty64/pdf-microservice
```
>Далее

```shell
cd ./pdf-microservice
```
```shell
cp ./.pdf-microservice-config-dev.sample.toml ./.pdf-microservice-config-dev.toml
```
>Затем

```shell
docker compose up pdf-microservice
```
>При последующих запусках:
```shell
docker compose up
```

Завершение работы:
```shell
docker compose stop
```


Используемый стэк:
 > go-1.23 || fpdf || minio-client || chi-v5 || viper

