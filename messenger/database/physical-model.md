## USER

| Название        | Описание         | Тип данных     | Ограничение        |
| --------------- | ---------------- | -------------- | ------------------ |
| `id`            | ID пользователя  | `SERIAL`       | `PRIMARY KEY`      |
| `nickname`      | Ник пользователя | `VARCHAR(64)`  | `UNIQUE, NOT NULL` |
| `password_hash` | Хэш пароля       | `VARCHAR(128)` | `NOT NULL`         |

## FRIENDSHIP

| Название   | Описание                        | Тип данных | Ограничение                                       |
| ---------- | ------------------------------- | ---------- | ------------------------------------------------- |
| `user_id1` | ID первого пользователя         | `INT`      | `NOT NULL, REFERENCES user(id) ON DELETE CASCADE` |
| `user_id2` | ID второго пользователя         | `INT`      | `NOT NULL, REFERENCES user(id) ON DELETE CASCADE` |

## FRIENDSHIP_INVITATION

| Название       | Описание                   | Тип данных | Ограничение                                       |
| -------------- | -------------------------- | ---------- | ------------------------------------------------- |
| `from_user_id` | ID отправителя приглашения | `INT`      | `NOT NULL, REFERENCES user(id) ON DELETE CASCADE` |
| `to_user_id`   | ID получателя приглашения  | `INT`      | `NOT NULL, REFERENCES user(id) ON DELETE CASCADE` |

## CHAT

| Название | Описание      | Тип данных    | Ограничение   |
| -------- | ------------- | ------------- | ------------- |
| `id`     | ID чата       | `SERIAL`      | `PRIMARY KEY` |
| `name`   | Название чата | `VARCHAR(64)` | `NOT NULL`    |

## CHAT_MEMBERSHIP

| Название  | Описание        | Тип данных | Ограничение                                       |
| --------- | --------------- | ---------- | ------------------------------------------------- |
| `chat_id` | ID чата         | `INT`      | `NOT NULL, REFERENCES chat(id) ON DELETE CASCADE` |
| `user_id` | ID пользователя | `INT`      | `NOT NULL, REFERENCES user(id) ON DELETE CASCADE` |

## CHAT_INVITATION

| Название  | Описание        | Тип данных | Ограничение                                       |
| --------- | --------------- | ---------- | ------------------------------------------------- |
| `chat_id` | ID чата         | `INT`      | `NOT NULL, REFERENCES chat(id) ON DELETE CASCADE` |
| `user_id` | ID пользователя | `INT`      | `NOT NULL, REFERENCES user(id) ON DELETE CASCADE` |

## MESSAGE

| Название   | Описание        | Тип данных | Ограничение                                       |
| ---------- | --------------- | ---------- | ------------------------------------------------- |
| `id`       | ID сообщения    | `SERIAL`   | `PRIMARY KEY`                                     |
| `user_id`  | ID отправителя  | `INT`      | `NOT NULL, REFERENCES user(id) ON DELETE CASCADE` |
| `chat_id`  | ID чата         | `INT`      | `NOT NULL, REFERENCES chat(id) ON DELETE CASCADE` |
| `msg_text` | Текст сообщения | `TEXT`     | `NOT NULL`                                        |
