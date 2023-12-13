package stdhttp

import (
	psg "addressBookServer/gate/psg"
	"addressBookServer/models/dto"
	errorLogger "addressBookServer/pkg/errorLogger"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type Controller struct {
	DB     *psg.Psg
	Srv    *http.Server
	Logger *errorLogger.ErrorLogger
}

func NewController(addr string, db *psg.Psg, logger *errorLogger.ErrorLogger) *Controller {

	controller := &Controller{
		Srv: &http.Server{
			Addr: addr,
		},
		DB:     db,
		Logger: logger,
	}

	http.HandleFunc("/add", controller.RecordAdd)
	http.HandleFunc("/get", controller.RecordsGet)
	http.HandleFunc("/update", controller.RecordUpdate)
	http.HandleFunc("/delete", controller.RecordDeleteByPhone)

	return controller
}

func (c *Controller) RecordAdd(w http.ResponseWriter, r *http.Request) {

	response := dto.Response{}

	defer responseWriteAndReturn(w, &response)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	rec := dto.Record{}
	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.RecordAdd(): json.NewDecoder(r.Body).Decode(&rec)")
		c.Logger.LogError(err)
		response.Wrap("error", nil, err)
		return
	}

	_, err = c.DB.RecordAdd(rec)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.RecordAdd(): c.DB.RecordAdd(rec)")
		c.Logger.LogError(err)
		response.Wrap("error", nil, err)
		return
	}

	response.Wrap("OK", nil, nil)

}

func (c *Controller) RecordsGet(w http.ResponseWriter, r *http.Request) {
	response := dto.Response{}

	defer responseWriteAndReturn(w, &response)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	rec := dto.Record{}
	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.RecordsGet(): json.NewDecoder(r.Body).Decode(&rec)")
		c.Logger.LogError(err)
		response.Wrap("error", nil, err)
		return
	}

	records, err := c.DB.RecordsGet(rec)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.RecordsGet(): c.DB.RecordsGet(rec)")
		c.Logger.LogError(err)
		response.Wrap("error", nil, err)
		return
	}

	recordsJSON, err := json.Marshal(&records)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.RecordsGet(): json.Marshal(&records)")
		c.Logger.LogError(err)
		response.Wrap("error", nil, err)
		return
	}

	response.Wrap("OK", recordsJSON, nil)

}

func (c *Controller) RecordUpdate(w http.ResponseWriter, r *http.Request) {

	response := dto.Response{}

	defer responseWriteAndReturn(w, &response)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// getting record from client
	rec := dto.Record{}
	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.RecordUpdate(): json.NewDecoder(r.Body).Decode(&rec)")
		c.Logger.LogError(err)
		response.Wrap("error", nil, err)
		return
	}

	// making update on DB
	err = c.DB.RecordUpdate(rec)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.RecordUpdate(): c.DB.RecordUpdate(rec)")
		c.Logger.LogError(err)
		response.Wrap("error", nil, err)
		return
	}

	response.Wrap("OK", nil, nil)

}

func (c *Controller) RecordDeleteByPhone(w http.ResponseWriter, r *http.Request) {

	response := dto.Response{}

	defer responseWriteAndReturn(w, &response)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	rec := dto.Record{}

	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.RecordDeleteByPhone(): json.NewDecoder(r.Body).Decode(&rec)")
		c.Logger.LogError(err)
		response.Wrap("error", nil, err)
		return
	}

	recPhone := rec.Phone

	err = c.DB.RecordDeleteByPhone(recPhone)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.RecordDeleteByPhone(): c.DB.RecordDeleteByPhone(recPhone)")
		c.Logger.LogError(err)
		response.Wrap("error", nil, err)
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
