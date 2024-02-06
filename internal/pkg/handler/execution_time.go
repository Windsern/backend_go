package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

type BuildingStateRequest struct {
	AccessKey     int64 `json:"access_key"`
	BuildingState int   `json:"state"`
}

type Request struct {
	BuildingId     int64 `json:"building_id"`
	VerificationId int64 `json:"verification_id"`
}

func (h *Handler) issueBuildingState(c *gin.Context) {
	var input Request
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println("handler.issueBuildingState:", input)

	c.Status(http.StatusOK)

	go func() {
		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
		sendBuildingStateRequest(input)
	}()
}

func sendBuildingStateRequest(request Request) {

	var state = -1
	if rand.Intn(10)%10 >= 2 {
		state = rand.Intn(100)
	}

	answer := BuildingStateRequest{
		AccessKey:     123,
		BuildingState: state,
	}

	client := &http.Client{}

	jsonAnswer, _ := json.Marshal(answer)
	bodyReader := bytes.NewReader(jsonAnswer)

	requestURL := fmt.Sprintf("http://127.0.0.1:8000/api/verifications/%d/update_building/%d/", request.VerificationId, request.BuildingId)

	req, _ := http.NewRequest(http.MethodPut, requestURL, bodyReader)

	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending PUT request:", err)
		return
	}

	defer response.Body.Close()

	fmt.Println("PUT Request Status:", response.Status)
}
