package save

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "go-url-shortner/internal/lib/api/response"
	"go-url-shortner/internal/lib/logger/sl"
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
	}
}
