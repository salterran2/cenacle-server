package service

import (
	"strconv"
	"time"

	errors "github.com/habasement/cenacle-server/error"
	"github.com/habasement/cenacle-server/person"
	"github.com/habasement/cenacle-server/person/repository"
)

//PersonService interface
type PersonService interface {
	Fetch(cursor string, num int64) ([]*person.Person, string, error)
	GetByID(id int64) (*person.Person, error)
	Update(ar *person.Person) (*person.Person, error)
	Store(*person.Person) (*person.Person, error)
	Delete(id int64) (bool, error)
}

type personService struct {
	personRepos repository.PersonRepository
}

func (p *personService) Fetch(cursor string, num int64) ([]*person.Person, string, error) {
	if num == 0 {
		num = 10
	}

	listPerson, err := p.personRepos.Fetch(cursor, num)
	if err != nil {
		return nil, "", err
	}
	nextCursor := ""

	if size := len(listPerson); size == int(num) {
		lastID := listPerson[num-1].ID
		nextCursor = strconv.Itoa(int(lastID))
	}

	return listPerson, nextCursor, nil
}

func (p *personService) GetByID(id int64) (*person.Person, error) {

	return p.personRepos.GetByID(id)
}

func (p *personService) Update(pr *person.Person) (*person.Person, error) {
	pr.UpdatedAt = time.Now()
	return p.personRepos.Update(pr)
}

func (p *personService) Store(m *person.Person) (*person.Person, error) {

	existedPerson, _ := p.GetByID(m.ID)
	if existedPerson != nil {
		return nil, errors.ErrConflict
	}

	id, err := p.personRepos.Store(m)
	if err != nil {
		return nil, err
	}

	m.ID = id
	return m, nil
}

func (p *personService) Delete(id int64) (bool, error) {
	existedPerson, _ := p.GetByID(id)

	if existedPerson == nil {
		return false, errors.ErrNotFound
	}

	return p.personRepos.Delete(id)
}

//NewPersonService NewPersonService Object
func NewPersonService(p repository.PersonRepository) PersonService {
	return &personService{p}
}
