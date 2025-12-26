package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"

	db "github.com/google-dev-groups-gmu/ghost/go/internal/firestore"
	"github.com/google-dev-groups-gmu/ghost/go/internal/types"
)

// returns static lat/long data for the map
// GET /api/buildings
func GetBuildings(c *gin.Context) {
	buildings := map[string]map[string]interface{}{
		"HORIZN": {"name": "Horizon Hall", "lat": 38.8296, "lng": -77.3072},
		"EXPL":   {"name": "Exploratory Hall", "lat": 38.8291, "lng": -77.3060},
	}
	c.JSON(http.StatusOK, buildings)
}

// returns the full list of rooms and their schedules
// GET /api/rooms?building=HORIZN
func GetRooms(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := db.Client
	if client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database not initialized"})
		return
	}

	// filter by building via query param ?building=HORIZN
	buildingFilter := c.Query("building")

	var iter *firestore.DocumentIterator
	if buildingFilter != "" {
		iter = client.Collection("rooms").Where("building", "==", buildingFilter).Documents(ctx)
	} else {
		iter = client.Collection("rooms").Documents(ctx)
	}
	defer iter.Stop()

	var rooms []types.Room
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("firestore error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error reading database"})
			return
		}

		var room types.Room
		if err := doc.DataTo(&room); err != nil {
			continue
		}
		rooms = append(rooms, room)
	}

	// cache for 60min lets save costs
	c.Header("Cache-Control", "public, max-age=3600")
	c.JSON(http.StatusOK, rooms)
}
