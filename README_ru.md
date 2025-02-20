<a id="readme-top"></a>

<h1 align="center">DNS Forwarder</h1>

<p align="center">
  Простой и эффективный DNS-форвардер на Go!
  <br />
  <br />
  <a href="https://github.com/kartmos/dns-forwarder/issues/new?labels=bug&template=bug-report.md">Сообщить об ошибке</a>
  &middot;
  <a href="https://github.com/kartmos/dns-forwarder/issues/new?labels=enhancement&template=feature-request.md">Предложить улучшение</a>
</p>

<!-- ОГЛАВЛЕНИЕ -->
<details>
  <summary>Оглавление</summary>
  <ol>
    <li><a href="#о-проекте">О проекте</a></li>
    <li><a href="#возможности">Возможности</a></li>
    <li><a href="#начало-работы">Начало работы</a></li>
    <li><a href="#использование">Использование</a></li>
  </ol>
</details>

<!-- О ПРОЕКТЕ -->
## О проекте

DNS Forwarder — это простой и эффективный DNS-форвардер, написанный на Go, который перенаправляет DNS-запросы на указанный DNS-сервер. Проект демонстрирует основы работы с сетевыми запросами и DNS-протоколом в Go. Основные особенности:

* Поддержка UDP
* Настраиваемый адрес DNS-сервера
* Логирование запросов
* Простая конфигурация

<p align="right">(<a href="#readme-top">наверх</a>)</p>

### Возможности

- Перенаправление DNS-запросов
- Поддержка UDP
- Настраиваемый адрес DNS-сервера
- Логирование запросов
- Простая конфигурация

<!-- НАЧАЛО РАБОТЫ -->
## Начало работы

### Требования

- Go 1.24.0+
- Интернет-соединение

### Установка

1. Клонируйте репозиторий.
2. Перейдите в директорию проекта.
3. Соберите бинарный файл в корневой папке проекта.

<p align="right">(<a href="#readme-top">наверх</a>)</p>

<!-- ИСПОЛЬЗОВАНИЕ -->
## Использование Docker-compose

```zsh
docker-compose -f build/deploy/docker-compose.yml up -d
```

<p align="right">(<a href="#readme-top">наверх</a>)</p>