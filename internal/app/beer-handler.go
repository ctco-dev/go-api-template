package app

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ctco-dev/go-api-template/internal/beer"

	"github.com/ctco-dev/go-api-template/internal/log"
)

// JokeHandler handles /joke route
type beerHandler struct {
	repository beer.Repository
}

func (b *beerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, _ = ShiftPath(r.URL.Path)

	if head != "" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		b.handleGet(w, r)
	case "PUT":
		b.handlePut(w, r)
	case "DELETE":
		b.handleDelete(w, r)
	default:
		http.Error(w, "Only GET, PUT and DELETE are allowed", http.StatusMethodNotAllowed)
	}
}

func (b *beerHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ids, ok := r.URL.Query()["id"]

	if !ok || len(ids[0]) < 1 {
		getAllBeers(ctx, w, b.repository)
		return
	}

	getOneBeer(ctx, w, ids[0], b.repository)
}

func (b *beerHandler) handlePut(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	names, ok := r.URL.Query()["name"]

	if !ok || len(names[0]) < 1 {
		http.Error(w, "Beer name is missing.", http.StatusInternalServerError)
		return
	}

	writeBeer(ctx, w, names[0], b.repository)
}

func (b *beerHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ids, ok := r.URL.Query()["id"]

	if !ok || len(ids[0]) < 1 {
		http.Error(w, "Beer id is missing.", http.StatusInternalServerError)
		return
	}

	removeBeer(ctx, w, ids[0], b.repository)
}

func getOneBeer(ctx context.Context, w http.ResponseWriter, id string, reader beer.Reader) {
	beer, err := reader.Read(ctx, id)
	if err != nil {
		http.Error(w, "Can't get a beer.", http.StatusInternalServerError)
		return
	}

	writeResult(ctx, w, beer)
}

func getAllBeers(ctx context.Context, w http.ResponseWriter, reader beer.AllReader) {
	beers, err := reader.ReadAll(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeResult(ctx, w, beers)
}

func writeBeer(ctx context.Context, w http.ResponseWriter, name string, writer beer.Writer) {
	id, err := writer.Write(ctx, beer.Beer{Name: name})
	if err != nil {
		http.Error(w, "Can't write a beer.", http.StatusInternalServerError)
		return
	}

	writeResult(ctx, w, id)
}

func removeBeer(ctx context.Context, w http.ResponseWriter, id string, remover beer.Remover) {
	err := remover.Remove(ctx, id)
	if err != nil {
		http.Error(w, "Can't delete a beer.", http.StatusInternalServerError)
		return
	}

	writeResult(ctx, w, "success")
}

func writeResult(ctx context.Context, w http.ResponseWriter, data interface{}) {
	bytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Can't encode response.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
	log.WithCtx(ctx).Info("I'm done")
}
