package services

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type ServiceHandler struct {
	otogenius *OtogeniusSvc
}

func NewServiceHandler(otogenius *OtogeniusSvc) *ServiceHandler {
	return &ServiceHandler{
		otogenius: otogenius,
	}
}

func (h *ServiceHandler) registerHandler() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/recommendation", h.otogenius.GetRecommendation).Methods(http.MethodPost)
	return router
}

func (h *ServiceHandler) Start() error {
	apiHost := os.Getenv("API_HOST")
	apiPort := os.Getenv("API_PORT")
	listenAddr := fmt.Sprintf("%s:%s", apiHost, apiPort)

	router := h.registerHandler()

	server := new(http.Server)
	server.Handler = router
	server.Addr = listenAddr

	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			return err
		}

		methods, _ := route.GetMethods()
		fmt.Printf("Path: %s, Methods: %v\n", pathTemplate, methods)
		return nil
	})

	if err != nil {
		return err
	}

	fmt.Println("services running on: ", listenAddr)
	err = server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
