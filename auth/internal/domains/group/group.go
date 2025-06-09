package domains

import "time"

// Group представляет доменную сущность «группа»
// Группа объединяет пользователей по какому-либо признаку или правам доступа
type Group struct {
	// id — уникальный идентификатор группы
	id GroupID
	// name — название группы
	name string
	// description — описание или цель группы
	description string
	// createdAt — момент создания группы
	createdAt time.Time
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
