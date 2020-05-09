package api

import (
	. "github.com/vi-sense/vi-sense/app/model"
	"net/http"
)

// QueryRoomModels godoc
// @Summary Query models
// @Description Query all available room models.
// @Tags models
// @Produce  json
// @Success 200 {array} model.RoomModel
// @Failure 400 {string} string "bad request"
// @Failure 500 {string} string "internal server error"
// @Router /models [get]
func QueryRoomModels() (int, string) {
	var q []RoomModel
	DB.Find(&q)
	return http.StatusOK, asJSON(&q)
}

// QueryRoomModel godoc
// @Summary Query room model
// @Description Query a single room model by id with containing sensors.
// @Tags models
// @Produce json
// @Param id path int true "RoomModelID"
// @Success 200 {object} model.RoomModel
// @Failure 400 {string} string "bad request"
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal server error"
// @Router /models/{id} [get]
func QueryRoomModel(id string) (int, string)  {
	var q RoomModel
	DB.Preload("Sensors").First(&q, id)
	if q.ID == 0 {
		return http.StatusNotFound, ""
	}
	return http.StatusOK, asJSON(&q)
}
