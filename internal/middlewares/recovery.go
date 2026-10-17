package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
)

func Recovery(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			defer func() {
				if err := recover(); err != nil {
					stack := debug.Stack()

					logger.ErrorContext(r.Context(), "panic recuperado",
						slog.Any("error", err),
						slog.String("stack", string(stack)),
						slog.String("path", r.URL.Path),
					)

					http.Error(w, "Error interno", http.StatusInternalServerError)
				}

			}()
			fmt.Println("recovery")
			next.ServeHTTP(w, r)

		})
	}

}
