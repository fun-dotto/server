package handler

import (
	"net/http"

	api "github.com/fun-dotto/server/gen/admin"
	"github.com/gin-gonic/gin"
)

func (h *Handler) BusTripsV1List(c *gin.Context, params api.BusTripsV1ListParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (h *Handler) BusTimetableStopsV1List(c *gin.Context, tripId string) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
