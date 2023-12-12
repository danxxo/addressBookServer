package stdhttp

import (
	psg "addressBookServer/gate/psg"
	"addressBookServer/models/dto"
	errorLogger "addressBookServer/pkg/errorLogger"
	"encoding/json"
	"fmt"
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

	response := dto.Response{}

	defer responseWriteAndReturn(w, &response)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	rec := dto.Record{}
	err = json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		// FIXME: errors wrap врапит только стороку "stdhttp.RecordAdd()"
		// Строк ошибки в логах должна выглядеть так: "<error_text>: stdhttp.RecordAdd(): json.NewDecoder(r.Body).Decode(&rec)"
		err = errors.Wrap(err, "stdhttp.RecordAdd()")
		logger.LogError(err)
		response.ErrorWrap(err)
		return
	}

	_, err = c.DB.RecordAdd(rec)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.RecordAdd()")
		logger.LogError(err)
		response.ErrorWrap(err)
		return
	}

	response.Wrap("OK", nil, nil)

}

func (c *Controller) RecordsGet(w http.ResponseWriter, r *http.Request) {

	// FIXME: во всех методах контроллера, где используется логгер, логгер должен быть создан в самом контроллере, а не в методе.
	// creating logger
	logger, err := errorLogger.NewErrorLogger(".log")
	if err != nil {
		fmt.Println(err)
		return
	}

	response := dto.Response{}

	defer responseWriteAndReturn(w, &response)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	rec := dto.Record{}
	err = json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.RecordsGet()")
		logger.LogError(err)
		response.ErrorWrap(err)
		return
	}

	records, err := c.DB.RecordsGet(rec)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.RecordsGet()")
		logger.LogError(err)
		response.ErrorWrap(err)
		return
	}

	recordsJSON, err := json.Marshal(&records)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.RecordsGet()")
		logger.LogError(err)
		response.ErrorWrap(err)
		return
	}

	response.Wrap("OK", recordsJSON, nil)

}

func (c *Controller) RecordUpdate(w http.ResponseWriter, r *http.Request) {

	// creating logger
	logger, err := errorLogger.NewErrorLogger(".log")
	if err != nil {
		fmt.Println(err)
		return
	}

	response := dto.Response{}

	defer responseWriteAndReturn(w, &response)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// getting record from client
	rec := dto.Record{}
	err = json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.RecordUpdate()")
		logger.LogError(err)
		response.ErrorWrap(err)
		return
	}

	// making update on DB
	err = c.DB.RecordUpdate(rec)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.RecordUpdate()")
		logger.LogError(err)
		response.ErrorWrap(err)
		return
	}

	response.Wrap("OK", nil, nil)

}

func (c *Controller) RecordDeleteByPhone(w http.ResponseWriter, r *http.Request) {

	// creating logger
	logger, err := errorLogger.NewErrorLogger(".log")
	if err != nil {
		fmt.Println(err)
		return
	}

	response := dto.Response{}

	defer responseWriteAndReturn(w, &response)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	rec := dto.Record{}

	err = json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.RecordDeleteByPhone()")
		logger.LogError(err)
		response.ErrorWrap(err)
		return
	}

	recPhone := rec.Phone

	err = c.DB.RecordDeleteByPhone(recPhone)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.RecordDeleteByPhone()")
		logger.LogError(err)
		response.ErrorWrap(err)
		return
	}

	response.Wrap("OK", nil, nil)
}

func responseWriteAndReturn(w http.ResponseWriter, response *dto.Response) {
	errEncode := json.NewEncoder(w).Encode(response)
	if errEncode != nil {
		w.WriteHeader(http.StatusPaymentRequired)
	}
}
