package domains

type Group struct {
	id          GroupID
	name        string
	description string
}

func NewGroup(
	id GroupID,
	name string,
	description string,
) Group {
	return Group{
		id:          id,
		name:        name,
		description: description,
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
