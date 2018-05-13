package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/essem/learngo/addressbook/pb"
	_ "github.com/go-sql-driver/mysql"

	"golang.org/x/net/context"
)

type contextKey int

const (
	idContextKey contextKey = 0
)

type webService struct {
	db         *appDB
	httpServer *http.Server
}

func (s *webService) init(address string) {
	router := http.NewServeMux()
	router.HandleFunc("/people", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.list(w, r)
		case http.MethodPost:
			s.create(w, r)
		default:
			http.NotFound(w, r)
		}
	})
	router.HandleFunc("/people/", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.URL.Path[len("/people/"):], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newContext := r.WithContext(context.WithValue(r.Context(), idContextKey, id))

		switch r.Method {
		case http.MethodGet:
			s.read(w, newContext)
		case http.MethodPatch, http.MethodPut:
			s.update(w, newContext)
		case http.MethodDelete:
			s.delete(w, newContext)
		default:
			http.NotFound(w, r)
		}
	})

	hander := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			next.ServeHTTP(w, r)
		})
	}

	s.httpServer = &http.Server{
		Addr:    address,
		Handler: hander(router),
	}
}

func (s *webService) serve() error {
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *webService) shutdown() {
	s.httpServer.Shutdown(context.Background())
}

func writeJSON(w http.ResponseWriter, p interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// curl localhost:8090/people
func (s *webService) list(w http.ResponseWriter, r *http.Request) {
	people, err := s.db.list()
	if err != nil {
		log.Fatal(err)
	}

	writeJSON(w, people)
}

// curl -X POST -d '{"name":"test","email":"test@test.com"}' localhost:8090/people
func (s *webService) create(w http.ResponseWriter, r *http.Request) {
	person := pb.Person{}
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := s.db.create(&person)
	if err != nil {
		log.Fatal(err)
	}

	writeJSON(w, id)
}

// curl localhost:8090/people/1
func (s *webService) read(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(idContextKey).(int64)

	person, err := s.db.read(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, person)
}

// curl -X PATCH -d '{"name":"test2","email":"test2@test.com"}' localhost:8090/people/1
func (s *webService) update(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(idContextKey).(int64)

	person := pb.Person{}
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	person.Id = id

	err = s.db.update(&person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, true)
}

// curl -X DELETE localhost:8090/people/1
func (s *webService) delete(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(idContextKey).(int64)

	err := s.db.delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, true)
}
