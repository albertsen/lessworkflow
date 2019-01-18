package repo

import (
	"net/http"

	"github.com/albertsen/lessworkflow/pkg/db/conn"
	"github.com/go-pg/pg"
)

var (
	Connect = conn.Connect
	Close   = conn.Close
)

func Select(record interface{}) (int, error) {
	if err := conn.DB().Select(record); err != nil {
		if err == pg.ErrNoRows {
			return http.StatusNotFound, nil
		}
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func Insert(record interface{}) (int, error) {
	if err := conn.DB().Insert(record); err != nil {
		pgError, ok := err.(pg.Error)
		if ok && pgError.IntegrityViolation() {
			return http.StatusConflict, err
		}
		return http.StatusInternalServerError, err
	}
	return http.StatusCreated, nil
}

func Update(record interface{}) (int, error) {
	if err := conn.DB().Update(record); err != nil {
		if err == pg.ErrNoRows {
			return http.StatusNotFound, err
		}
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func Delete(record interface{}) (int, error) {
	if err := conn.DB().Delete(record); err != nil {
		if err == pg.ErrNoRows {
			return http.StatusNotFound, err
		}
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
