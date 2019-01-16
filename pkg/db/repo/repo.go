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

func Select(Record interface{}) (int, error) {
	if err := conn.DB().Select(Record); err != nil {
		if err == pg.ErrNoRows {
			return http.StatusNotFound, nil
		}
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func Insert(Record interface{}) (int, error) {
	if err := conn.DB().Insert(Record); err != nil {
		pgError, ok := err.(pg.Error)
		if ok && pgError.IntegrityViolation() {
			return http.StatusConflict, err
		}
		return http.StatusInternalServerError, err
	}
	return http.StatusCreated, nil
}

func Update(Record interface{}) (int, error) {
	if err := conn.DB().Update(Record); err != nil {
		if err == pg.ErrNoRows {
			return http.StatusNotFound, err
		}
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func Delete(Record interface{}) (int, error) {
	if err := conn.DB().Delete(Record); err != nil {
		if err == pg.ErrNoRows {
			return http.StatusNotFound, err
		}
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
