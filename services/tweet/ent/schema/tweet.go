package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"time"
)

// Tweet holds the schema definition for the Tweet entity.
type Tweet struct {
	ent.Schema
}

// Fields of the Tweet.
func (Tweet) Fields() []ent.Field {
	return []ent.Field{
		field.String("author_id").NotEmpty(),
		field.String("content").NotEmpty(),
		field.Time("created_at").Default(time.Now),
	}
}

// Indexes of the Tweet.
func (Tweet) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("author_id"),
		index.Fields("created_at"),
	}
}
