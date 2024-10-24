package api

import (
	"errors"
	"log/slog"
	"net/http"

	iErrors "github.com/xV0lk/zibutatu/internal/errors"
	eView "github.com/xV0lk/zibutatu/views/error"
)

// apiHandler is a function type that handles API requests.
type apiHandler func(w http.ResponseWriter, r *http.Request) error

// MakeHandler is a higher-order function that wraps an apiHandler function
// and returns an http.HandlerFunc. Its purpose is to centralice error handling
// and being able to return errors from out handler functions.
//
// It handles any errors returned by the apiHandler and logs them using slog.
// If the error is an unhandled ApiError, it will log the error and render an error page
func MakeHandler(mh apiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ae *iErrors.ApiError
		if err := mh(w, r); err != nil {
			println("Inside error")
			if errors.As(err, &ae) {
				println("Inside api error")
				slog.Error(ae.Title,
					slog.Int("status", ae.Status),
					slog.String("errorMsg", ae.Msg),
					slog.String("error", ae.Error()),
					slog.String("location", ae.Trace),
					slog.Any("body", ae.Body),
				)
				if ae.Handled == iErrors.HttpUnhandledError {
					eView.ErrorPage(ae).Render(r.Context(), w)
				}
			}
		}
	}
}
