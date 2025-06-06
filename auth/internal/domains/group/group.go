package domains

import "time"

type Group struct {
	id          GroupID
	name        string
	description string
	createdAt   time.Time
}

func NewGroup(
	id GroupID,
	name string,
	description string,
	createdAt time.Time,
) Group {
	return Group{
		id:          id,
		name:        name,
		description: description,
		createdAt:   createdAt,
	}
}

func (g Group) ID() GroupID {
	return g.id
}

func (g Group) Name() string {
	return g.name
}

func (g Group) Description() string {
	return g.description
}

func (g Group) CreatedAt() time.Time {
	return g.createdAt
}
