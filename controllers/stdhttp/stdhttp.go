package stdhttp

import (
	psg "addressBookServer/gate/psg"
	"addressBookServer/models/dto"
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
	rec := dto.Record{}
	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		err = errors.Wrap(err, "stdhttp.RecordAdd():json.Decode error")
		fmt.Println(err)
		return
	}
	fmt.Println("Recieved ", rec)
	_, err = c.DB.RecordAdd(rec)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (c *Controller) RecordsGet(w http.ResponseWriter, r *http.Request) {
	rec := dto.Record{}
	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Recieved ", rec)

	records, err := c.DB.RecordsGet(rec)
	if err != nil {
		fmt.Println(err)
	}

	var recordsToJson dto.Records = records

	fmt.Println(recordsToJson)
	jsondData, err := json.Marshal(recordsToJson)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = w.Write(jsondData)
	if err != nil {
		fmt.Println(err)
		return
	}

}

func (c *Controller) RecordUpdate(w http.ResponseWriter, r *http.Request) {
	rec := dto.Record{}
	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Recieved ", rec)

	err = c.DB.RecordUpdate(rec)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Updated")
}

func (c *Controller) RecordDeleteByPhone(w http.ResponseWriter, r *http.Request) {
	rec := dto.Record{}
	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Recieved ", rec)

	recPhone := rec.Phone

	err = c.DB.RecordDeleteByPhone(recPhone)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Deleted")
}
