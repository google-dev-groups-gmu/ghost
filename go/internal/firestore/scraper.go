package firestore

import (
	"context"

	"github.com/google-dev-groups-gmu/ghost/go/internal/types"
)

// saves the aggregated schedule for a room
func SaveRoom(ctx context.Context, room types.Room) error {
	if Client == nil {
		return nil
	}
	_, err := Client.Collection("rooms").Doc(room.ID).Set(ctx, room)
	return err
}
