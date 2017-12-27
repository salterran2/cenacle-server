package controller

import (
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"

	errors "github.com/habasement/cenacle-server/error"
	models "github.com/habasement/cenacle-server/person"
	personService "github.com/habasement/cenacle-server/person/service"

	"github.com/labstack/echo"
	validator "gopkg.in/go-playground/validator.v9"
)

//PersonController PersonController struct
type PersonController struct {
	PService personService.PersonService
}

//FetchPerson echo
func (controller *PersonController) FetchPerson(c echo.Context) error {

	numS := c.QueryParam("num")
	num, _ := strconv.Atoi(numS)

	cursor := c.QueryParam("cursor")

	personList, nextCursor, err := controller.PService.Fetch(cursor, int64(num))

	if err != nil {
		return c.JSON(getStatusCode(err), err.Error())
	}
	c.Response().Header().Set(`X-Cursor`, nextCursor)
	return c.JSON(http.StatusOK, personList)
}

//GetByID get By ID
func (controller *PersonController) GetByID(c echo.Context) error {

	idP, err := strconv.Atoi(c.Param("id"))
	id := int64(idP)

	person, err := controller.PService.GetByID(id)

	if err != nil {
		return c.JSON(getStatusCode(err), err.Error())
	}
	return c.JSON(http.StatusOK, person)
}

func isRequestValid(m *models.Person) (bool, error) {

	validate := validator.New()

	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

//Store store
func (controller *PersonController) Store(c echo.Context) error {
	var person models.Person
	err := c.Bind(&person)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ok, err := isRequestValid(&person)
	if !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	pr, err := controller.PService.Store(&person)

	if err != nil {
		return c.JSON(getStatusCode(err), err.Error())
	}
	return c.JSON(http.StatusCreated, pr)
}

// Delete delete
func (controller *PersonController) Delete(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	id := int64(idP)

	_, err = controller.PService.Delete(id)

	if err != nil {

		return c.JSON(getStatusCode(err), err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func getStatusCode(err error) int {
	if err != nil {
		logrus.Error(err)
	}
	switch err {
	case errors.ErrInternalServer:

		return http.StatusInternalServerError
	case errors.ErrNotFound:
		return http.StatusNotFound
	case errors.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

//NewPersonController person controller
func NewPersonController(e *echo.Echo, pService personService.PersonService) {
	controller := &PersonController{
		PService: pService,
	}

	e.GET("/person", controller.FetchPerson)
	e.POST("/person", controller.Store)
	e.GET("/person/:id", controller.GetByID)
	e.DELETE("/person/:id", controller.Delete)

}
