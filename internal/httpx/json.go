package httpx

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type errorBody struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	if v == nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(status)
		return
	}

	buf, err := json.Marshal(v)
	if err != nil {
		slog.Error("marshal json", "err", err)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":{"code":"internal","message":"internal error"}}`))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(buf)
	_, _ = w.Write([]byte("\n"))
}

func WriteError(w http.ResponseWriter, status int, code, msg string) {
	WriteJSON(w, status, errorBody{
		Error: errorDetail{Code: code, Message: msg},
	})
}
