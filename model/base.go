package model

// ModelBase provides common fields and functionalities for all model
type ModelBase struct {
	ID int
}

// NewModelBase creates a new ModelBase instance and all entities should have an ID, right? :)
func NewModelBase(id int) *ModelBase {
	return &ModelBase{ID: id}
}
