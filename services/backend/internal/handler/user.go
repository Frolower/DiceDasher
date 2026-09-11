package handler

import (
	"backend/internal/auth"
	"context"
	"diceDasher/pkg/httputil"
	"diceDasher/pkg/logger"
	"errors"
	"net/http"
)

// Registrar — минимальная возможность, которая нужна HTTP-слою.
// Обработчик ничего не знает о PostgreSQL и bcrypt: в main ему передают
// готовый сервис, а в тестах — простую заглушку с тем же методом.
type Registrar interface {
	Register(context.Context, auth.CreateInput) (auth.Created, error)
}

// Handler хранит зависимости, общие для запросов. Данные самого запроса
// остаются локальными переменными методов и не смешиваются между клиентами.
type Handler struct{ registrar Registrar }

func New(registrar Registrar) *Handler { return &Handler{registrar: registrar} }

func (h *Handler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	// Ограничиваем тело до декодирования. Для трёх полей 16 КиБ достаточно;
	// слишком большой запрос не должен занимать неограниченную память.
	r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
	var req createUserRequest
	// Декодируем сразу в структуру: UnpackJSON отклоняет неизвестные поля,
	// некорректный JSON и несколько JSON-значений в одном теле.
	if err := httputil.UnpackJSON(r, &req); err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
		} else {
			// Не возвращаем текст декодера: он может содержать присланные значения.
			http.Error(w, "invalid request body", http.StatusBadRequest)
		}
		return
	}

	// HTTP-модель превращается во входные данные сценария регистрации.
	// Проверка полей, хеширование и запись выполняются внутри auth.Service.
	_, err := h.registrar.Register(r.Context(), auth.CreateInput{
		Username: req.Username, Email: req.Email, Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidInput):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, auth.ErrAlreadyExists):
			http.Error(w, "username or email already exists", http.StatusConflict)
		default:
			// Клиент не получает детали БД. Репозиторий также не переносит
			// PostgreSQL Detail в ошибки: там может оказаться email пользователя.
			logger.Logf(r.Context(), "ERROR registering user: %s", err)
			http.Error(w, "failed to create user", http.StatusInternalServerError)
		}
		return
	}

	// 201 означает, что пользователь уже сохранён. Сессию, cookie и токены
	// регистрация пока не создаёт; внутренний UUID наружу не требуется.
	if err := httputil.PackJSON(w, http.StatusCreated, createUserResponse{Status: "created"}); err != nil {
		// Заголовки уже отправлены: второй HTTP-ответ здесь писать нельзя.
		logger.Logf(r.Context(), "ERROR writing registration response: %s", err)
	}
}
