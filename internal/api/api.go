package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"weekly-radio-programme/internal/show"
)

type Api struct {
	Router *chi.Mux
	srv    *show.Service
	port   int
}

func New(srv *show.Service, port int) *Api {
	api := Api{}
	api.Router = chi.NewRouter()

	api.Router.Use(middleware.RequestID)
	api.Router.Use(middleware.Logger)

	api.Router.Get("/health", api.healthCheck)
	api.Router.Post("/show", api.create)
	api.Router.Get("/show/{id}", api.get)
	api.Router.Put("/show/{id}", api.update)
	api.Router.Delete("/show/{id}", api.delete)
	api.Router.Get("/show", api.getAll)
	api.srv = srv
	api.port = port
	return &api
}

func (o *Api) Run() error {
	log.Println("serving to port", o.port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", o.port), o.Router); err != nil {
		return err
	}
	return nil
}

func (o *Api) healthCheck(w http.ResponseWriter, r *http.Request) {
	renderJson(w, r, http.StatusOK, "hello!")
}

func (o *Api) getAll(w http.ResponseWriter, r *http.Request) {
	shows, err := o.srv.GetAll(r.Context())
	if err != nil {
		log.Println(err)
		renderJson(w, r, http.StatusInternalServerError, Response{})
		return
	}
	renderJson(w, r, http.StatusOK, shows)
}

func (o *Api) get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		renderJson(w,r,http.StatusBadRequest,Response{Message: "bad id given"})
	}
	s, err := o.srv.Get(r.Context(), id)
	if err != nil {
		log.Println(err)
		renderJson(w, r, http.StatusInternalServerError, Response{})
		return
	}
	renderJson(w, r, http.StatusOK, s)
}

type createReq struct {
	Title       string `json:"title"`
	Weekday     string `json:"weekday"`
	Timeslot    string `json:"timeslot"`
	Description string `json:"description"`
}

func (o *Api) create(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		renderJson(w, r, http.StatusBadRequest, Response{Message: "cannot read body"})
		return
	}
	defer r.Body.Close()
	var s show.Show
	if err := json.Unmarshal(bodyBytes, &s); err != nil {
		renderJson(w, r, http.StatusBadRequest, Response{Message: "cannot read body"})
		return
	}
	if err := s.Validate(); err != nil {
		renderJson(w, r, http.StatusBadRequest, Response{Message: err.Error()})
		return
	}
	if err := o.srv.CheckForConflicts(r.Context(), s); err != nil {
		switch err.(type) {
		case show.ErrTimeslotConflict:
			renderJson(w, r, http.StatusConflict, Response{Message: err.Error()})
		default:
			log.Println(err)
			renderJson(w, r, http.StatusInternalServerError, Response{})
		}
		return
	}
	id, err := o.srv.Add(r.Context(), s)
	if err != nil {
		log.Println(err)
		renderJson(w, r, http.StatusInternalServerError, Response{})
		return
	}
	renderJson(w, r, http.StatusCreated,
		Response{Message: fmt.Sprintf("show with id %d was created", id)})
}

func (o *Api) update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		renderJson(w,r,http.StatusBadRequest,Response{Message: "bad id given"})
	}
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		renderJson(w, r, http.StatusBadRequest, Response{Message: "cannot read body"})
		return
	}
	defer r.Body.Close()
	var s show.Show
	if err := json.Unmarshal(bodyBytes, &s); err != nil {
		renderJson(w, r, http.StatusBadRequest, Response{Message: "cannot read body"})
		return
	}
	if err := s.Validate(); err != nil {
		renderJson(w, r, http.StatusBadRequest, Response{Message: err.Error()})
		return
	}
	s.Id = id
	if err := o.srv.CheckForConflicts(r.Context(), s); err != nil {
		switch err.(type) {
		case show.ErrTimeslotConflict:
			renderJson(w, r, http.StatusConflict, Response{Message: err.Error()})
		default:
			log.Println(err)
			renderJson(w, r, http.StatusInternalServerError, Response{})
		}
		return
	}
	if err := o.srv.Update(r.Context(), s); err != nil {
		log.Println(err)
		renderJson(w, r, http.StatusInternalServerError, Response{})
		return
	}
	renderJson(w, r, http.StatusCreated,
		Response{Message: fmt.Sprintf("show with id %d was updated", s.Id)})
}

func (o *Api) delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		renderJson(w,r,http.StatusBadRequest,Response{Message: "bad id given"})
	}
	if err := o.srv.Delete(r.Context(), id); err != nil {
		log.Println(err)
		renderJson(w, r, http.StatusInternalServerError, Response{})
		return
	}
	renderJson(w, r, http.StatusOK,
		Response{Message: fmt.Sprintf("show with id %d was deleted", id)})
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

type Response struct {
	Message string `json:"message"`
}

func renderJson(w http.ResponseWriter, r *http.Request, statusCode int, res interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := bufferPool.Get().(*bytes.Buffer)
	defer func() {
		response.Reset()
		bufferPool.Put(response)
	}()
	var err error
	if res != nil {
		err = json.NewEncoder(response).Encode(res)
		if err != nil {
			apiError := Response{Message: "Cannot json encode response"}
			_ = json.NewEncoder(response).Encode(apiError)
			statusCode = 500
		}
	}
	w.WriteHeader(statusCode)
	io.Copy(w, response)
}
