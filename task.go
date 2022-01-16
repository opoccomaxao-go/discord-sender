package discordsender

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//nolint:maligned // ObjectID must be inline
type Task struct {
	ID         primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	HookURL    string             `bson:"hook_url"`
	PostData   string             `bson:"post_data"`
	Expiration time.Time          `bson:"expiration"`
	Executed   bool               `bson:"executed"`
}
