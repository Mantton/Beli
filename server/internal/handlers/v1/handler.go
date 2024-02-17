package v1

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/mantton/beli/internal/cache"
)

type V1Handler struct {
	cache *cache.Cache
}

func New(c *cache.Cache) *V1Handler {
	return &V1Handler{
		cache: c,
	}
}

type tileBody struct {
	X     int `json:"x"`
	Y     int `json:"y"`
	Color int `json:"color"`
}

type board struct {
	EncodedString string `json:"encoded"`
}

func (b tileBody) Valid() bool {
	if b.X < 0 || b.X > 10 || b.Y < 0 || b.Y > 10 || b.Color < 0 || b.Color > 16 {
		return false
	}

	return true
}

/*
POST /draw
*/
func (h *V1Handler) HandleDrawTile(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	var body tileBody
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if !body.Valid() {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = h.cache.SetTile(r.Context(), int64(body.X), int64(body.Y), body.Color)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		slog.Error(err.Error())
	}
}

func (h *V1Handler) HandleGetTile(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	xStr := q.Get("x")
	yStr := q.Get("y")

	fmt.Println(xStr, yStr)

	if len(xStr) == 0 || len(yStr) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	x, err := strconv.ParseInt(xStr, 10, 64)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	y, err := strconv.ParseInt(yStr, 10, 64)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	color, err := h.cache.GetTile(r.Context(), x, y)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}

	res := tileBody{
		X:     int(x),
		Y:     int(y),
		Color: color,
	}

	err = json.NewEncoder(w).Encode(res)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		slog.Error(err.Error())
	}

}

func (h *V1Handler) HandleGetBoard(w http.ResponseWriter, r *http.Request) {

	res, err := h.cache.GetBoard(r.Context(), "CURRENT_BOARD")

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}

	b := board{
		EncodedString: res,
	}

	err = json.NewEncoder(w).Encode(b)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		slog.Error(err.Error())
	}
}
