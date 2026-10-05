# In progress

- Бэкенд
  - [x] Регистрация — POST /register, без сессии; [описание](services/backend/README.md)
  - [x] Вход — username + password, access JWT и refresh-cookie
  - [x] Создание сессии — при успешном входе, хранится хеш refresh-токена
  - [ ] Смена пароля
  - [x] Отзыв сессии — POST /logout, удаление refresh-cookie
  - [ ] Обновление токенов (refresh)
