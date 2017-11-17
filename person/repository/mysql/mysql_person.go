package mysql

import (
	"database/sql"
	"time"

	"github.com/Sirupsen/logrus"

	errors "github.com/habasement/cenacle-server/error"
	models "github.com/habasement/cenacle-server/person"
	"github.com/habasement/cenacle-server/person/repository"
)

type mysqlPersonRepository struct {
	Conn *sql.DB
}

func (m *mysqlPersonRepository) fetch(query string, args ...interface{}) ([]*models.Person, error) {

	rows, err := m.Conn.Query(query, args...)

	if err != nil {
		return nil, errors.ErrInternalServer
	}

	defer rows.Close()

	result := make([]*models.Person, 0)
	for rows.Next() {
		t := new(models.Person)
		err = rows.Scan(
			&t.ID,
			&t.Firstname,
			&t.Lastname,
			&t.EmailAddress,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			return nil, errors.ErrInternalServer
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlPersonRepository) Fetch(cursor string, num int64) ([]*models.Person, error) {

	query := `SELECT id, firstname, lastname, email_address, updated_at, created_at
  						FROM person WHERE id > ? LIMIT ?`

	return m.fetch(query, cursor, num)

}
func (m *mysqlPersonRepository) GetByID(id int64) (*models.Person, error) {
	query := `SELECT  id, firstname, lastname, email_address, updated_at, created_at
  						FROM person WHERE id = ?`

	list, err := m.fetch(query, id)
	if err != nil {
		return nil, err
	}

	person := &models.Person{}
	if len(list) > 0 {
		person = list[0]
	} else {
		return nil, errors.ErrInternalServer
	}

	return person, nil
}

func (m *mysqlPersonRepository) Store(person *models.Person) (int64, error) {

	query := `INSERT person SET firstname=? , lastname=? , email_address=? , updated_at=? , created_at=?`
	stmt, err := m.Conn.Prepare(query)
	if err != nil {
		return 0, err
	}
	now := time.Now()
	person.CreatedAt = now
	person.UpdatedAt = now
	logrus.Debug("Created At: ", now)
	res, err := stmt.Exec(person.Firstname, person.Lastname, person.EmailAddress, now, now)
	if err != nil {

		return 0, err
	}
	return res.LastInsertId()
}

func (m *mysqlPersonRepository) Delete(id int64) (bool, error) {
	query := "DELETE FROM person WHERE id = ?"

	stmt, err := m.Conn.Prepare(query)
	if err != nil {
		return false, err
	}
	res, err := stmt.Exec(id)
	if err != nil {

		return false, err
	}
	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAfected != 1 {
		logrus.Error("Weird  Behaviour. Total Affected ", rowsAfected)
		return false, errors.ErrInternalServer
	}

	return true, nil
}
func (m *mysqlPersonRepository) Update(person *models.Person) (*models.Person, error) {
	query := `UPDATE person set firstname=?, lastname=?, email_address=?, updated_at=? WHERE ID = ?`

	stmt, err := m.Conn.Prepare(query)
	if err != nil {
		return nil, nil
	}
	now := time.Now()
	res, err := stmt.Exec(person.Firstname, person.Lastname, person.EmailAddress, now, person.ID)
	if err != nil {
		return nil, err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affect != 1 {
		logrus.Error("Weird  Behaviour. Total Affected ", affect)
		return nil, errors.ErrInternalServer
	}

	return person, nil
}

//NewMysqlPersonRepository mysql person repostory
func NewMysqlPersonRepository(Conn *sql.DB) repository.PersonRepository {

	return &mysqlPersonRepository{Conn}
}
