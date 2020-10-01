package server

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

// func  GetScheduleHandler(w http.ResponseWriter, r *http.Request) {
// 	session, err := getCloudSessions(r)
// 	if err != nil {
// 		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
// 		return
// 	}

// 	vars := mux.Vars(r)

// 	accountID, ok := vars["accountID"]

// 	if !ok {
// 		handleError(w, http.StatusBadRequest, "could not get accountID")
// 		return
// 	}

// 	docs, err := session.GetDocument(accountID)
// 	if err != nil {
// 		handleError(w, http.StatusUnauthorized, "could not get docs", err.Error())
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	e := json.NewEncoder(w)
// 	e.Encode(docs)
// 	return
// }

func SetScheduleHandler(c echo.Context) error {
	session, err := getCloudSessions(c)
	if err != nil {
		return err
	}

	var body map[string]interface{}
	decoder := json.NewDecoder(c.Request().Body)
	err = decoder.Decode(&body)
	if err != nil {
		return err
	}

	accountID := c.Param("accountID")

	if err := session.CreateDocument(accountID, body); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, StatusOK{Message: "success"})
}

// TODO: complete this
func DeleteScheduleHandler(c echo.Context) error {
	_, err := getCloudSessions(c)
	if err != nil {
		return err
	}
	return nil
}

func UpdateScheduleHandler(c echo.Context) error {
	_, err := getCloudSessions(c)
	if err != nil {
		return err
	}
	return nil
}

func GetAllScheduleHandler(c echo.Context) error {
	session, err := getCloudSessions(c)
	if err != nil {
		return err
	}

	accountID := c.Param("accountID")

	docs, err := session.GetAllDocument(accountID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, docs)
}
