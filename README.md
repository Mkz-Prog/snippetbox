# Snippetbox

Веб-сервис для хранения и обмена фрагментами кода (сниппетами). Проект реализован в рамках изучения backend-разработки на языке **Go**.

## 🚀 О проекте

Snippetbox позволяет пользователям создавать и просматривать сниппеты. Основная цель проекта — освоение принципов построения RESTful-приложений, работы с базой данных и организации архитектуры веб-сервисов на Go.

## 🛠 Технологический стек

* **Язык:** Go (Golang)
* **База данных:** MySQL
* **Маршрутизация:** `net/http` (стандартная библиотека)
* **Шаблонизатор:** `html/template`

## 📋 Основной функционал

* Создание новых сниппетов через веб-интерфейс.
* Просмотр списка последних фрагментов кода.
* Чтение конкретного сниппета по уникальному ID.

## ⚙️ Установка и запуск

1. **Клонируйте репозиторий:**

    ```bash
    git clone https://github.com/Mkz-Prog/snippetbox.git
    cd snippetbox
    ```

2. **Настройте базу данных:**

    **Автоматическая настройка**

    🚧 *В РАЗРАБОТКЕ* 🚧

    ---

    **Ручная настройка:**

    * Установите базу данных mariadb-server или MySQL-server:

        ```bash
        apt install mariadb-server
        ```

    * Войдите в термиинал БД:

        ```bash
        sudo mariadb
        ```

    * Создайте БД snippetbox:

        ```sql
        -- Create a new UTF-8 `snippetbox` database.
        CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

        -- Switch to using the `snippetbox` database.
        USE snippetbox;
        ```

    * Создайте таблицу snippets:

        ```sql
        -- Create a `snippets` table.
        CREATE TABLE snippets (
            id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
            title VARCHAR(100) NOT NULL,
            content TEXT NOT NULL,
            created DATETIME NOT NULL,
            expires DATETIME NOT NULL
        );

        -- Add an index on the created column.
        CREATE INDEX idx_snippets_created ON snippets(created);
        ```

    * Внесите записи в таблицу. Пример:

        ```sql
        INSERT INTO snippets (title, content, created, expires) VALUES (
            'An old silent pond',
            'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
            UTC_TIMESTAMP(),
            DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
        );
        ```

    * Создайте пользователя и задайте ему пароль. Замените `pass` на свой пароль:

        ```sql
        CREATE USER 'web'@'localhost';
        GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO 'web'@'localhost';
        ALTER USER 'web'@'localhost' IDENTIFIED BY 'pass';
        ```

    * Запишите свой пароль в переменную окружения `DB_PASSWORD` в `.env` файл, например с помощью команды:

        ```bash
        echo DB_PASSWORD=pass > .env
        ```

3. **Запуск приложения:**

    ```bash
    go run ./cmd/web
    ```

## 📈 Планы по развитию

* Реализация системы авторизации и аутентификации пользователей.
* Контейнеризация приложения с помощью Docker для упрощения деплоя
