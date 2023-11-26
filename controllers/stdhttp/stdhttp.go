package stdhttp

import (
	psg "addressBookServer/gate/psg"
	"addressBookServer/models/dto"
	errorLogger "addressBookServer/pkg/errorLogger"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

type Controller struct {
	DB  *psg.Psg
	Srv *http.Server
}

func NewController(addr string, db *psg.Psg) *Controller {

	controller := &Controller{
		Srv: &http.Server{
			Addr: addr,
		},
		DB: db,
	}

	http.HandleFunc("/add", controller.RecordAdd)
	http.HandleFunc("/get", controller.RecordsGet)
	http.HandleFunc("/update", controller.RecordUpdate)
	http.HandleFunc("/delete", controller.RecordDeleteByPhone)

	return controller
}

func (c *Controller) RecordAdd(w http.ResponseWriter, r *http.Request) {

	// creating logger
	logger, err := errorLogger.NewErrorLogger(".log")
	if err != nil {
		fmt.Println(err)
		return
	}

	if r.Method == "POST" {
		var err error

		// sending response
		defer func() {
			response := dto.Response{ErrorMessage: ""}
			if err != nil {
				response.ErrorMessage = err.Error()
				logger.LogError(err)
			}
			jsonResponse, err := json.Marshal(response)
			if err != nil {
				logger.LogError(err)
			}
			w.Write(jsonResponse)
		}()

		rec := dto.Record{}
		err = json.NewDecoder(r.Body).Decode(&rec)
		if err != nil {
			err = errors.Wrap(err, "stdhttp.RecordAdd():")
			return
		}

		_, err = c.DB.RecordAdd(rec)
		if err != nil {
			err = errors.Wrap(err, "stdhttp.RecordAdd():")
			return
		}
	}
}

func (c *Controller) RecordsGet(w http.ResponseWriter, r *http.Request) {

	// creating logger
	logger, err := errorLogger.NewErrorLogger(".log")
	if err != nil {
		fmt.Println(err)
		return
	}

	if r.Method == "POST" {
		var err error
		var records []dto.Record

		// sending response
		defer func() {
			response := dto.Response{ErrorMessage: ""}
			response.Records = records
			if err != nil {
				response.ErrorMessage = err.Error()
				logger.LogError(err)
			}
			jsonResponse, err := json.Marshal(response)
			if err != nil {
				logger.LogError(err)
			}
			w.Write(jsonResponse)
		}()

		rec := dto.Record{}
		err = json.NewDecoder(r.Body).Decode(&rec)
		if err != nil {
			err = errors.Wrap(err, "stdhttp.RecordsGet()")
			return
		}

		records, err = c.DB.RecordsGet(rec)
		if err != nil {
			err = errors.Wrap(err, "stdhttp.RecordsGet()")
			return
		}

	}
}

func (c *Controller) RecordUpdate(w http.ResponseWriter, r *http.Request) {

	// creating logger
	logger, err := errorLogger.NewErrorLogger(".log")
	if err != nil {
		fmt.Println(err)
		return
	}

	if r.Method == "POST" {
		var err error

		// sending response
		defer func() {
			response := dto.Response{ErrorMessage: ""}
			if err != nil {
				response.ErrorMessage = err.Error()
				logger.LogError(err)
			}
			jsonResponse, err := json.Marshal(response)
			if err != nil {
				logger.LogError(err)
			}
			w.Write(jsonResponse)
		}()

		// getting record from client
		rec := dto.Record{}
		err = json.NewDecoder(r.Body).Decode(&rec)
		if err != nil {
			err = errors.Wrap(err, "stdhttp.RecordUpdate()")
			log.Println(err)
			return
		}

		// making update on DB
		err = c.DB.RecordUpdate(rec)
		if err != nil {
			err = errors.Wrap(err, "stdhttp.RecordUpdate()")
			log.Println(err)
			return
		}
	}
}

func (c *Controller) RecordDeleteByPhone(w http.ResponseWriter, r *http.Request) {

	// creating logger
	logger, err := errorLogger.NewErrorLogger(".log")
	if err != nil {
		fmt.Println(err)
		return
	}

	if r.Method == "POST" {
		var err error

		// sending response
		defer func() {
			response := dto.Response{ErrorMessage: ""}
			if err != nil {
				response.ErrorMessage = err.Error()
				logger.LogError(err)
			}
			jsonResponse, err := json.Marshal(response)
			if err != nil {
				logger.LogError(err)
			}
			w.Write(jsonResponse)
		}()

		rec := dto.Record{}
		err = json.NewDecoder(r.Body).Decode(&rec)
		if err != nil {
			err = errors.Wrap(err, "stdhttp.RecordDeleteByPhone()")
			return
		}

		recPhone := rec.Phone

		err = c.DB.RecordDeleteByPhone(recPhone)
		if err != nil {
			err = errors.Wrap(err, "stdhttp.RecordDeleteByPhone()")
			return
		}
	}
}
