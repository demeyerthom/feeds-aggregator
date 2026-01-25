package internal

import (
	"time"

	"github.com/mmcdole/gofeed"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FeedList []Feed

type Feed struct {
	Title  string `json:"title"`
	XMLURL string `json:"xmlUrl"`
}

type Post struct {
	gofeed.Item
	FileName string `json:"fileName"`
}

type FeedItem struct {
	Link  string `json:"link"`
	Title string `json:"title"`
}

// FeedItemDocument is the MongoDB document model for storing feed items
type FeedItemDocument struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Link      string             `bson:"link"`
	Title     string             `bson:"title"`
	Summary   string             `bson:"summary,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
}
