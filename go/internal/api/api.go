package api

import (
	"context"
	"log"
	"net/http"
	"strconv"
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
	c.JSON(http.StatusOK, types.Buildings)
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
	dayFilterStr := c.Query("day")
	timeFilterStr := c.Query("time")

	var iter *firestore.DocumentIterator
	if buildingFilter != "" {
		iter = client.Collection("rooms").Where("building", "==", buildingFilter).Documents(ctx)
	} else {
		iter = client.Collection("rooms").Documents(ctx)
	}
	defer iter.Stop()

	var rooms []types.Room

	var filterDay int = -1
	if dayFilterStr != "" {
		if d, err := strconv.Atoi(dayFilterStr); err == nil {
			filterDay = d
		}
	}

	var filterTime int = -1
	if timeFilterStr != "" {
		if t, err := strconv.Atoi(timeFilterStr); err == nil {
			filterTime = t
		}
	}

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

		if filterDay != -1 {
			var todaysSchedule []types.Meeting

			for _, item := range room.Schedule {
				if item.Day == filterDay {
					// filter only the classes that are ongoing at the specified time
					if filterTime != -1 {
						if item.StartTime <= filterTime && item.EndTime >= filterTime {
							todaysSchedule = append(todaysSchedule, item)
						}
					} else {
						todaysSchedule = append(todaysSchedule, item)
					}
				}
			}
			room.Schedule = todaysSchedule
		}

		rooms = append(rooms, room)
	}

	// cache for 60min lets save costs
	c.Header("Cache-Control", "public, max-age=3600")
	c.JSON(http.StatusOK, rooms)
}

func GetSpecificRoom(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := db.Client
	if client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database not initialized"})
		return
	}

	// filter by room via query param ?room=HORIZN_2014
	roomFilter := c.Query("room")
	doc, err := client.Collection("rooms").Doc(roomFilter).Get(ctx)
	if err != nil {
		log.Printf("firestore error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error reading database"})
		return
	}

	var room types.Room
	if err := doc.DataTo(&room); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error parsing data"})
		return
	}
	c.JSON(http.StatusOK, room)
}
