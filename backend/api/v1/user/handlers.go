package user

import (
	"backend/models"
	"backend/utils"
	"errors"
	"fmt"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *UserHandler) SignUp(c echo.Context) error {
	req := &userUpdateRequest{}
	user, err := req.bind(c)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	user, err = h.userStore.Create(user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	token, err := generateJWT(user.Id, user.Username, h.jwtSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusCreated, newUserTokenResponse(user, token))
}

func (h *UserHandler) Login(c echo.Context) error {
	req := &userLoginRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	user, err := h.userStore.GetByUsername(req.User.Username)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewError(err))
	}

	if user == nil {
		return c.JSON(http.StatusNotFound, utils.NewError(errors.New("Not found")))
	}

	match, err := argon2id.ComparePasswordAndHash(req.User.Password, user.Password)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	if !match {
		return c.JSON(http.StatusForbidden, utils.NewError(errors.New("Unauthorized")))
	}

	token, err := generateJWT(user.Id, user.Username, h.jwtSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusCreated, newUserTokenResponse(user, token))
}

func (h *UserHandler) GetMe(c echo.Context) error {
	id, err := userIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewError(err))
	}

	user, err := h.userStore.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewError(err))
	}

	if user == nil {
		return c.JSON(http.StatusNotFound, utils.NewError(errors.New("Not found")))
	}

	token, err := generateJWT(user.Id, user.Username, h.jwtSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, newUserTokenResponse(user, token))
}

func (h *UserHandler) GetUser(c echo.Context) error {
	fmt.Printf("a %v", c.Param("id"))
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	user, err := h.userStore.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, newUserResponse(user))
}

func (h *UserHandler) UpdateMe(c echo.Context) error {
	id, err := userIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewError(err))
	}

	user, err := h.userStore.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewError(err))
	}

	user, err = h.innerUpdateUser(c, user)
	if err != nil {
		return err
	}
	fmt.Printf("XD %v \n", user)
	token, err := generateJWT(user.Id, user.Username, h.jwtSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, newUserTokenResponse(user, token))
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	user, err := h.userStore.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NewError(err))
	}

	user, err = h.innerUpdateUser(c, user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, newUserResponse(user))
}

func (h *UserHandler) innerUpdateUser(c echo.Context, user *models.User) (*models.User, error) {
	req := &userUpdateRequest{}
	new_user, err := req.bind(c)
	if err != nil {
		fmt.Printf("\naa %v\n ", err)
		return nil, err
	}

	new_user.Id = user.Id

	user, err = h.userStore.Update(new_user)
	if err != nil {
		fmt.Printf("\nbb %v\n ", err)
		return nil, err
	}

	fmt.Printf("\nar %v\n ", user)
	return user, nil
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	if err = h.userStore.Delete(id); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, nil)
}
