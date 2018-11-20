# DC Courier Monitor Service
[![Build Status](https://travis-ci.com/TeamD2018/geo-rest.svg?branch=suite_test)](https://travis-ci.com/TeamD2018/geo-rest)

Микросервис, позволяющий отслеживать геолокацию курьеров в рамках приложения [Delivery Club](https://www.delivery-club.ru).

Разработка сервиса осуществлялась в рамках итогого задания образовательного проекта [Технопарк](https://park.mail.ru/) от компании [Mail.Ru](https://mail.ru).

## Getting Started

Микросервис написан на языке Go и использует в качестве зависимостей БД Elasticsearch и Tarantool

## Installation

Для удобства установки всех зависимостей проекта и самого микросервиса используется [Docker](https://www.docker.com)

### Зависимости

* Docker
> Для установки Docker воспользуйтесь официальным [туториалом](https://docs.docker.com/install/linux/docker-ce/centos/)


* Склонируйте репозиторий, выполнив

```
git clone https://github.com/TeamD2018/geo-rest
```


В качестве конфигурации сервис использует конфиг-файл с расширением .toml

* Создайте файл `geo_rest_config.toml` в директории проекта или переименуйте `geo_rest_config.toml.example`

* Выполните

```
docker-compose up -f docker-compose.yml
```

Сервис будет доступен по адресу `server.url` указанному в конфиге

Доступное API расположено по адресу [openapi.track-delivery.club](openapi.track-delivery.club)

## Team

* [Данила Масленников](https://github.com/Dnnd)
* [Даниил Котельников](https://github.com/zwirec)
* [Олег Уткин](https://github.com/oleggator)