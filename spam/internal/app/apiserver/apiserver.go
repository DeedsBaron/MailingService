package apiserver

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"spam/internal/app/config"
	"spam/internal/app/store"
	"spam/internal/app/validation"
	"spam/pkg/utils"
	"strconv"
	"time"
)

type APIServer struct {
	config  *config.Config
	logger  *logrus.Logger
	router  *mux.Router
	storage store.Storage
}

func New(config *config.Config) *APIServer {
	return &APIServer{
		config:  config,
		logger:  logrus.New(),
		router:  mux.NewRouter(),
		storage: store.InitStorage(config),
	}
}

func (s *APIServer) Start() error {
	time.LoadLocation("Local")
	if err := s.configureLogger(); err != nil {
		return err
	}
	s.configureRouter()
	s.logger.Info("Starting API server")
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/AddUser", s.AddUser()).Methods("POST")
	s.router.HandleFunc("/UpdateUser/{id}", s.UpdateUser()).Methods("PATCH")
	s.router.HandleFunc("/DeleteUser/{id}", s.DeleteUser()).Methods("DELETE")
	s.router.HandleFunc("/AddMailingList", s.AddMailingList()).Methods("POST")
	s.router.HandleFunc("/GeneralStats", s.GeneralStats()).Methods("GET")
	s.router.HandleFunc("/docs/", s.DocsRedirect()).Methods("GET")
	s.router.HandleFunc("/DetailedStats/{id}", s.GetDetailedStats()).Methods("GET")
	s.router.HandleFunc("/UpdateMailingList/{id}", s.UpdateMailingList()).Methods("PATCH")
	s.router.HandleFunc("/DeleteMailingList/{id}", s.DeleteMailingList()).Methods("DELETE")
}

func (s *APIServer) DeleteMailingList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(mux.Vars(r)["id"])
		err := s.storage.DeleteMailingList(context.Background(), id)
		if err != nil {
			utils.HttpErrorWithoutBackSlashN(w, err.Error(), http.StatusBadRequest)
			s.logger.Error(err.Error())
			return
		}
		io.WriteString(w, "Successfully deleted MailingList from database!")
		s.logger.Debug("DELETE DeleteMailingList method SUCCESS")
	}
}

func (s *APIServer) UpdateMailingList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := validation.MailingList{}
		d := json.NewDecoder(r.Body)
		err := d.Decode(&t)
		if err != nil {
			utils.HttpErrorWithoutBackSlashN(w, "Bad JSON!", http.StatusBadRequest)
			s.logger.Error(err.Error())
			return
		}
		ok := t.ValidateMailingListJSON(w, s.logger, d)
		if ok == false {
			return
		}
		id, _ := strconv.Atoi(mux.Vars(r)["id"])
		err = s.storage.UpdateMailingList(context.Background(), &t, id)
		if err != nil {
			utils.HttpErrorWithoutBackSlashN(w, err.Error(), http.StatusBadRequest)
			s.logger.Error(err.Error())
			return
		}
		io.WriteString(w, "Successfully updated Client in database!")
		s.logger.Debug("PATCH UpdateMailingList method SUCCESS")
	}
}

func (s *APIServer) GetDetailedStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := s.storage.GetDetailedStats(context.Background(), w, mux.Vars(r)["id"])
		if err != nil {
			utils.HttpErrorWithoutBackSlashN(w, err.Error(), http.StatusBadRequest)
			s.logger.Error(err.Error())
			return
		}
		s.logger.Debug("GET GetDetailedStats method SUCCESS")
	}
}

func (s *APIServer) AddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := validation.User{}
		d := json.NewDecoder(r.Body)
		err := d.Decode(&t)
		if err != nil {
			utils.HttpErrorWithoutBackSlashN(w, "Bad JSON!", http.StatusBadRequest)
			s.logger.Error(err.Error())
			return
		}
		if t.ValidateAddUserJSON(w, s.logger, d) == false {
			return
		}
		err = s.storage.AddUser(context.Background(), &t)
		if err != nil {
			utils.HttpErrorWithoutBackSlashN(w, err.Error(), http.StatusBadRequest)
			s.logger.Error(err.Error())
			return
		}
		io.WriteString(w, "Successfully created Client in database!")
		s.logger.Debug("POST CreateUser method SUCCESS")
	}
}

func (s *APIServer) UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := validation.User{}
		d := json.NewDecoder(r.Body)
		err := d.Decode(&t)
		if err != nil {
			utils.HttpErrorWithoutBackSlashN(w, "Bad JSON!", http.StatusBadRequest)
			s.logger.Error(err.Error())
			return
		}
		ok := t.ValidateUpdateUserJSON(w, s.logger, d)
		if ok == false {
			return
		}
		id, _ := strconv.Atoi(mux.Vars(r)["id"])
		err = s.storage.UpdateUser(context.Background(), &t, id)
		if err != nil {
			utils.HttpErrorWithoutBackSlashN(w, err.Error(), http.StatusBadRequest)
			s.logger.Error(err.Error())
			return
		}
		io.WriteString(w, "Successfully updated Client in database!")
		s.logger.Debug("POST UpdateUser method SUCCESS")
	}
}

func (s *APIServer) DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			utils.HttpErrorWithoutBackSlashN(w, "Wrong ID format!", http.StatusBadRequest)
			s.logger.Error("Wrong ID format!")
			return
		}
		err = s.storage.DeleteUser(context.Background(), id)
		if err != nil {
			utils.HttpErrorWithoutBackSlashN(w, err.Error(), http.StatusBadRequest)
			s.logger.Error(err.Error())
			return
		}
		io.WriteString(w, "Successfully deleted Client from database!")
		s.logger.Debug("DELETE DeleteUser method SUCCESS")
	}
}

func (s *APIServer) AddMailingList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := validation.MailingList{}
		d := json.NewDecoder(r.Body)
		err := d.Decode(&t)
		if err != nil {
			utils.HttpErrorWithoutBackSlashN(w, "Bad JSON!", http.StatusBadRequest)
			s.logger.Error(err.Error())
			return
		}
		if t.ValidateMailingListJSON(w, s.logger, d) == false {
			return
		}
		err = s.storage.AddMailingList(context.Background(), &t)
		if err != nil {
			utils.HttpErrorWithoutBackSlashN(w, err.Error(), http.StatusBadRequest)
			s.logger.Error(err.Error())
			return
		}

		go s.storage.StartSpam(context.Background(), w, &t, s.config, s.logger)

		io.WriteString(w, "Successfully created MailingList in database and tried to send requests to extern server!")
		s.logger.Debug("POST AddMailingList method SUCCESS")
	}
}

func (s *APIServer) GeneralStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := s.storage.GeneralStats(context.Background(), w)
		if err != nil {
			utils.HttpErrorWithoutBackSlashN(w, err.Error(), http.StatusBadRequest)
			s.logger.Error(err.Error())
			return
		}
		s.logger.Debug("GET GeneralStats method SUCCESS")
	}
}

func (s *APIServer) DocsRedirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://app.swaggerhub.com/apis-docs/DeedsBaron/spam/1.0.0", http.StatusMovedPermanently)
		s.logger.Debug("Redirect SUCCESS")
	}
}
