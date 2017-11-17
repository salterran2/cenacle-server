package repository

import models "github.com/habasement/cenacle-server/person"

//PersonRepository interface
type PersonRepository interface {
	Fetch(cursor string, num int64) ([]*models.Person, error)
	GetByID(id int64) (*models.Person, error)
	Update(article *models.Person) (*models.Person, error)
	Store(a *models.Person) (int64, error)
	Delete(id int64) (bool, error)
}
