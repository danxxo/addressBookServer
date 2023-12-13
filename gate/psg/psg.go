package postgres

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	dto "addressBookServer/models/dto"
	phoneHelper "addressBookServer/pkg/phone"
)

type Psg struct {
	conn *pgxpool.Pool
}

func NewPsg(dburl string, login, pass string) (psg *Psg, err error) {
	defer func() { err = errors.Wrap(err, "postgres.NewPsg(dburl string, login, pass string)") }()

	psg = &Psg{}
	psg.conn, err = parseConnectionString(dburl, login, pass)
	if err != nil {
		return nil, err
	}

	err = psg.conn.Ping(context.Background())
	if err != nil {
		return
	}

	return psg, nil
}

func (p *Psg) RecordsGet(r dto.Record) (records []dto.Record, err error) {
	defer func() { err = errors.Wrap(err, "psg.RecordsGet(r dto.Record)") }()

	if r.Phone != "" {
		r.Phone, err = phoneHelper.PhoneNormalize(r.Phone)
		if err != nil {
			return
		}
	}

	q, values, err := SelectRecord(r)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(values) == 0 {
		q = `
		SELECT * FROM address_book;
		`
	}

	rows, err := p.conn.Query(context.Background(), q, values...)
	if err != nil {
		return
	}

	for rows.Next() {
		var rec dto.Record

		err = rows.Scan(&rec.ID, &rec.Name, &rec.LastName, &rec.MiddleName, &rec.Address, &rec.Phone)
		if err != nil {
			return
		}

		records = append(records, rec)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return records, nil
}

func (psg *Psg) RecordAdd(r dto.Record) (num int64, err error) {
	defer func() { err = errors.Wrap(err, "psg.RecordAdd(r dto.Record)") }()

	if !CheckAllFieldsIsFilled(r) {
		err = errors.New("Not all record fields are filled")
		return
	}

	query := `
	INSERT INTO address_book 
	(name, lastname, middlename, address, phone) 
	VALUES 
	($1, $2, $3, $4, $5);
	`
	r.Phone, err = phoneHelper.PhoneNormalize(r.Phone)
	if err != nil {
		return
	}

	_, err = psg.conn.Exec(context.Background(), query, r.Name, r.LastName, r.MiddleName, r.Address, r.Phone)
	if err != nil {
		return
	}
	return 0, nil
}

// RecordUpdate обновляет существующую запись в базе данных по номеру телефона.
func (p *Psg) RecordUpdate(r dto.Record) (err error) {
	defer func() { err = errors.Wrap(err, "psg.RecordUpdate(r dto.Record)") }()

	if r.Phone == "" {
		err = errors.New("Phone field is required")
		return
	}

	r.Phone, err = phoneHelper.PhoneNormalize(r.Phone)
	if err != nil {
		return
	}

	q, values, err := UpdateRecord(r)
	if err != nil {
		return
	}

	if len(values) < 2 {
		err = errors.New("Not enough Fields to update record. Need something besides Phone")
		return
	}

	values = append(values, r.Phone)

	_, err = p.conn.Exec(context.Background(), q, values...)
	if err != nil {
		return
	}
	return nil
}

// FIXME: все ресиверы должны быть одного названия. Нельзя писать psg, а потом p.
// RecordDeleteByPhone удаляет запись из базы данных по номеру телефона.
func (psg *Psg) RecordDeleteByPhone(phone string) (err error) {
	defer func() { err = errors.Wrap(err, "psg.RecordDeleteByPhone(phone string)") }()

	phone, err = phoneHelper.PhoneNormalize(phone)
	if err != nil {
		return
	}

	query := `
	DELETE FROM address_book
	WHERE phone = $1
	`

	_, err = psg.conn.Exec(context.Background(), query, phone)
	if err != nil {
		return
	}
	return nil
}

func parseConnectionString(dburl, user, password string) (db *pgxpool.Pool, err error) {
	defer func() { errors.Wrap(err, "psg.parseConnectionString(dburl, user, password string)") }()

	var u *url.URL
	if u, err = url.Parse(dburl); err != nil {
		return nil, errors.Wrap(err, "parseConnectionString(dburl, user, password string): url.Parse(dburl)")
	}
	u.User = url.UserPassword(user, password)
	db, err = pgxpool.New(context.Background(), u.String())
	if err != nil {
		return nil, errors.Wrap(err, "parseConnectionString(dburl, user, password string): pgxpool.New(context.Background(), u.String())")
	}
	return
}
