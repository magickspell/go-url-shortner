package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	resp "go-url-shortner/internal/lib/api/response"
	"go-url-shortner/internal/lib/logger/sl"
	"go-url-shortner/internal/lib/random"
	"go-url-shortner/internal/storage"
	"log/slog"
	"net/http"
)

// дублирует метод из storage (описываем интерфейсы там где они используются)
type URLSaver interface {
	Save(urlToSave string, alias string) (int64, error)
}

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alisa string `json:"alisa,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

// todo (move to config or db)
const aliasLength = 8

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request body"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		// validate
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, resp.Error("invalid request"))
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		alias := req.Alisa
		if alias == "" {
			// todo check unique alias
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.Save(req.URL, alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLExists) {
				log.Info("url already exists", slog.String("url", req.URL))
				render.JSON(w, r, resp.Error("url already exists"))
				return
			}

			log.Info("failed to save url", sl.Err(err))
		}

		log.Info("url saved", slog.Int64("id", id))

		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
