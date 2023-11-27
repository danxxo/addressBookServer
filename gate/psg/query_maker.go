package postgres

import (
	dto "addressBookServer/models/dto"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	gitErors "github.com/pkg/errors"
)

type SelectCond struct {
	Lop    string
	PgxInd string
	Field  string
	Value  any
}

type UpdateCond struct {
	Comma  string
	PgxInd string
	Field  string
	Value  any
}

func CheckAllFieldsIsFilled(r dto.Record) bool {
	return r.Address != "" && r.LastName != "" && r.MiddleName != "" && r.Name != "" && r.Phone != ""
}

func UpdateRecord(r dto.Record) (statement string, values []any, err error) {
	defer func() { err = gitErors.Wrap(err, "UpdateRecord") }()

	sqlFields, values, err := StructToFieldsValues(r, "sql.field")
	if err != nil {
		return
	}

	var conds []UpdateCond
	var lastIndex int

	for i := range sqlFields {
		lastIndex = i + 1
		if i == 0 {
			conds = append(conds, UpdateCond{
				Comma:  "",
				PgxInd: "$" + strconv.Itoa(i+1),
				Field:  sqlFields[i],
				Value:  values[i],
			})
			continue
		}
		conds = append(conds, UpdateCond{
			Comma:  ",",
			PgxInd: "$" + strconv.Itoa(i+1),
			Field:  sqlFields[i],
			Value:  values[i],
		})
	}

	query := `
	UPDATE 
	address_book
	SET
		{{range .}} {{.Comma}} {{.Field}}={{.PgxInd}}{{end}}
	WHERE phone = `

	query += "$" + strconv.Itoa(lastIndex+1) + "\n;"

	tmpl, err := template.New("").Parse(query)
	if err != nil {
		return
	}

	var sb strings.Builder

	err = tmpl.Execute(&sb, conds)
	if err != nil {
		return
	}

	statement = sb.String()
	return
}

func SelectRecord(r dto.Record) (statement string, values []any, err error) {
	defer func() { err = gitErors.Wrap(err, "SelectRecord()") }()

	sqlFields, values, err := StructToFieldsValues(r, "sql.field")
	if err != nil {
		return
	}

	var conds []SelectCond

	for i := range sqlFields {
		if i == 0 {
			conds = append(conds, SelectCond{
				Lop:    "",
				PgxInd: "$" + strconv.Itoa(i+1),
				Field:  sqlFields[i],
				Value:  values[i],
			})
			continue
		}
		conds = append(conds, SelectCond{
			Lop:    "AND",
			PgxInd: "$" + strconv.Itoa(i+1),
			Field:  sqlFields[i],
			Value:  values[i],
		})
	}

	query := `
	SELECT 
	id, name, lastname, middlename, address, phone
	FROM
	    address_book
	WHERE
		{{range .}} {{.Lop}} {{.Field}}={{.PgxInd}}{{end}}
;
`
	tmpl, err := template.New("").Parse(query)
	if err != nil {
		return
	}

	var sb strings.Builder
	err = tmpl.Execute(&sb, conds)
	if err != nil {
		return
	}
	statement = sb.String()
	return
}

func StructToFieldsValues(s any, tag string) (sqlFields []string, values []any, err error) {
	defer func() { err = gitErors.Wrap(err, "StructToFieldsValues()") }()
	rv := reflect.ValueOf(s)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return nil, nil, errors.New("s must be a struct")
	}

	for i := 0; i < rv.NumField(); i++ {
		field := rv.Type().Field(i)
		tg := strings.TrimSpace(field.Tag.Get(tag))
		if tg == "" || tg == "-" {
			continue
		}
		tgs := strings.Split(tg, ",")
		tg = tgs[0]

		fv := rv.Field(i)
		isZero := false
		switch fv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			isZero = fv.Int() == 0
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			isZero = fv.Uint() == 0
		case reflect.Float32, reflect.Float64:
			isZero = fv.Float() == 0
		case reflect.Complex64, reflect.Complex128:
			isZero = fv.Complex() == complex(0, 0)
		case reflect.Bool:
			isZero = !fv.Bool()
		case reflect.String:
			isZero = fv.String() == ""
		case reflect.Array, reflect.Slice:
			isZero = fv.Len() == 0
		}

		if isZero {
			continue
		}

		sqlFields = append(sqlFields, tg)
		values = append(values, fv.Interface())
	}
	return
}
