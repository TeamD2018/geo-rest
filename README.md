# DC Courier Monitor Service
[![Build Status](https://travis-ci.com/TeamD2018/geo-rest.svg?branch=suite_test)](https://travis-ci.com/TeamD2018/geo-rest)

Микросервис, позволяющий отслеживать геолокацию курьеров в рамках приложения [Delivery Club](https://www.delivery-club.ru).

Разработка сервиса осуществлялась в рамках итогого задания образовательного проекта [Технопарк](https://park.mail.ru/) от компании [Mail.Ru](https://mail.ru).

Микросервис написан на языке Go и использует в качестве зависимостей БД Elasticsearch и Tarantool

## Установка

Для удобства установки всех зависимостей проекта и самого микросервиса используется [Docker](https://www.docker.com)

### Зависимости

* Docker
> Для установки Docker воспользуйтесь официальным [туториалом](https://docs.docker.com/install/linux/docker-ce/centos/)
* [Kubernetes](https://kubernetes.io/docs/setup/) или [docker-compose](https://docs.docker.com/compose/install/#install-compose)

## [Установка через Kubernetes](https://github.com/TeamD2018/geo-rest-stuff/tree/master/deploy)
## Установка через Docker-Compose

* Склонируйте репозиторий, выполнив

```
git clone https://github.com/TeamD2018/geo-rest
```


В качестве конфигурации сервис использует конфиг-файл с расширением .toml

* Создайте файл `geo_rest_config.toml` в директории проекта или переименуйте `geo_rest_config.example.toml`

* Выполните

```
docker-compose -f docker-compose.yml up
```
or
```
docker-compose -f docker-compose.yml up -d
```

для запуска в фоновом режиме

Сервис будет доступен по адресу `server.url` указанному в конфиге

Доступное API расположено по адресу [openapi.track-delivery.club](http://openapi.track-delivery.club/)

## Команда

* [Данила Масленников](https://github.com/Dnnd)
* [Даниил Котельников](https://github.com/zwirec)
* [Олег Уткин](https://github.com/oleggator)
