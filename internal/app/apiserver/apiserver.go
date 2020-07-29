package apiserver

import (
	"context"
	"encoding/json"
	"encoding/base64"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/http-rest-api/internal/app/model"
	"github.com/http-rest-api/internal/app/store"
	"github.com/sirupsen/logrus"
)

// APIServer ...
type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store  *store.Store
}

// New ...
func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

// Start ...
func (s *APIServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureRouter()

	if err := s.configureStore(); err != nil {
		return err
	}

	s.logger.Info("starting api server")

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
	s.router.HandleFunc("/create-user", s.handleUsersCreate())
	s.router.HandleFunc("/session-create", s.handleSessionCreate())
	s.router.HandleFunc("/session-refresh", s.handleSessionRefresh())
	s.router.HandleFunc("/session-delete", s.handleDeleteSessionRefresh())
	s.router.HandleFunc("/delete-all-sessions", s.handleDeleteAllSessionsRefresh())
}

func (s *APIServer) configureStore() error {
	st := store.New(s.config.Store)
	if err := st.Open(); err != nil {
		return err
	}

	s.store = st

	return nil
}

func (s *APIServer) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}
		if id, err := s.store.User().Create(context.TODO(), u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		
		u.ID, err = primitive.ObjectIDFromHex(id)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *APIServer) handleSessionCreate() http.HandlerFunc {
	type request struct {
		UserID      string `json:"user_id"`
		Fingerprint string `json:"fingerprint"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		_, err := s.store.User().FindByID(context.TODO(), req.UserID)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		t := &model.Token{}
		tokens, err := t.CreatePairTokens(req.UserID, req.Fingerprint)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		err = s.store.RefreshSession().CreateNewSession(context.TODO(), req.UserID, tokens["refreshToken"], req.Fingerprint)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		encodedRefreshToken := base64.StdEncoding.EncodeToString([]byte(tokens["refreshToken"]))
		tokens["refreshToken"] = encodedRefreshToken

		s.respond(w, r, http.StatusOK, tokens)
	}
}

func (s *APIServer) handleSessionRefresh() http.HandlerFunc {
	type request struct {
		RefreshToken string `json:"refresh_token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		token := &model.Token{}

		rt, err := token.DecodeFromBase64(req.RefreshToken)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		err = token.ParseToken(rt)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		sessionID, err := s.store.RefreshSession().CheckRefreshSession(context.TODO(), token.UserID, rt, token.Fingerprint)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		s.store.RefreshSession().DeleteRefreshSession(context.TODO(), sessionID)

		t := &model.Token{}
		tokens, err := t.CreatePairTokens(token.UserID, token.Fingerprint)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		err = s.store.RefreshSession().CreateNewSession(context.TODO(), token.UserID, tokens["refreshToken"], token.Fingerprint)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		encodedRefreshToken := base64.StdEncoding.EncodeToString([]byte(tokens["refreshToken"]))
		tokens["refreshToken"] = encodedRefreshToken

		s.respond(w, r, http.StatusOK, tokens)
	}
}

func (s *APIServer) handleDeleteSessionRefresh() http.HandlerFunc {
	type request struct {
		RefreshToken string `json:"refresh_token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		token := &model.Token{}

		rt, err := token.DecodeFromBase64(req.RefreshToken)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		err = token.ParseToken(rt)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		sessionID, err := s.store.RefreshSession().CheckRefreshSession(context.TODO(), token.UserID, rt, token.Fingerprint)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		s.store.RefreshSession().DeleteRefreshSession(context.TODO(), sessionID)

		data := map[string]string{
			"message": "refresh session was deleted!",
		}

		s.respond(w, r, http.StatusOK, data)
	}
}

func (s *APIServer) handleDeleteAllSessionsRefresh() http.HandlerFunc {
	type request struct {
		RefreshToken string `json:"refresh_token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		token := &model.Token{}

		rt, err := token.DecodeFromBase64(req.RefreshToken)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		err = token.ParseToken(rt)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		err = s.store.RefreshSession().DeleteAllRefreshSessions(context.TODO(), token.UserID)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		data := map[string]string{
			"message": "All sessions was deleted",
		}

		s.respond(w, r, http.StatusOK, data)
	}
}

func (s *APIServer) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *APIServer) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
