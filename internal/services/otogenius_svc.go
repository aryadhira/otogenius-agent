package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aryadhira/otogenius-agent/internal/agent"
	"github.com/aryadhira/otogenius-agent/internal/repository"
	"github.com/aryadhira/otogenius-agent/utils"
)

type OtogeniusSvc struct {
	db    repository.CarRepo
	agent agent.Agent
}

func NewOtogeniusSvc(db repository.CarRepo, agent agent.Agent) *OtogeniusSvc {
	return &OtogeniusSvc{
		db:    db,
		agent: agent,
	}
}

func (o *OtogeniusSvc) GetRecommendation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSON(w, http.StatusBadRequest, "Bad Request", nil)
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, err.Error(), nil)
	}

	defer r.Body.Close()

	var payload map[string]interface{}

	err = json.Unmarshal(bodyBytes, &payload)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
	}

	input := utils.InterfaceToString(payload["input"])

	if input == "" {
		utils.WriteJSON(w, http.StatusBadRequest, "Input can't be empty", nil)
	}

	// Run the agent to get parameter
	res, err := o.agent.Run(input)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
	}

	var param map[string]any
	err = json.Unmarshal([]byte(res.(string)), &param)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("Error can't parse result from agent: %v", err), nil)
	}

	// Get car info
	datas, err := o.db.GetCarData(param)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
	}

	utils.WriteJSON(w, http.StatusOK, "", datas)
}
