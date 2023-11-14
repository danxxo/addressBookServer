Цель: Создать веб-сервис для управления записями адресной книги с использованием Go, базы данных PostgreSQL и стандартной библиотеки net/http.

Структура проекта:

/: Корневая директория проекта.
    main.go: Точка входа в приложение.
/controllers: Директория для контроллеров.
    /stdhttp: Директория для http контроллера.
        stdhttp.go: Контроллер для работы с записями адресной книги с использованием стандартной библиотеки.
/gate: Директория для работы с базой данных.
    /psg: Директория для gate к базе данных PostgreSQL.
        psg.go: Интерфейс и реализация gate к базе данных.
/pkg: Вспомогательные пакеты.
    phone.go: Функции для работы с номерами телефонов.
/models: Модели данных.
    /dto: Директория для Data Transfer Object.
        record.go: Определение структуры данных Record.
Основной функционал:

Добавление записи: Принимает запись и сохраняет её в базе данных.
Получение записей: Поиск записей по различным полям и их комбинациям. Также возможно получить все записи.
Обновление записи: Принимает запись с Phone и обновляет соответствующую запись в базе данных.
Удаление записи: Принимает номер телефона и удаляет соответствующую запись из базы данных.

в phone.go нужно реализовать функцию для обработки номера телефона, приводящую его к единому формату. Например, +7 (999) 123-45-67 -> 79991234567 или 8-999-123-45-67 -> 79991234567.

`func PhoneNormalize(phone string) (normalizedPhone string, err error) {
}`

`
// Record представляет запись в адресной книге.
type Record struct {
    ID        int64  `json:"id"`
	Name      string `json:"name"`
	LastName  string `json:"last_name"`
	MiddleName string `json:"middle_name,omitempty"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
}
`

Файл controllers/stdhttp/stdhttp.go:
`
package stdhttp

import (
	"net/http"
	"your_project_path/psg"
)

// Controller обрабатывает HTTP запросы для адресной книги.
type Controller struct {
	DB  *psg.Psg
	Srv *http.Server
}

// NewController создает новый Controller.
func NewController(addr string, db *psg.Psg) *Controller {
	return
}

// RecordAdd обрабатывает HTTP запрос для добавления новой записи.
func (c *Controller) RecordAdd(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализовать
}

// RecordsGet обрабатывает HTTP запрос для получения записей на основе предоставленных полей Record.
func (c *Controller) RecordsGet(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализовать
}

// RecordUpdate обрабатывает HTTP запрос для обновления записи.
func (c *Controller) RecordUpdate(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализовать
}

// RecordDeleteByPhone обрабатывает HTTP запрос для удаления записи по номеру телефона.
func (c *Controller) RecordDeleteByPhone(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализовать
}
`

Файл psg/psg.go:
`
package psg

import (
	"github.com/jackc/pgx/v5"
)

// Psg представляет гейт к базе данных PostgreSQL.
type Psg struct {
	Conn *pgxpool.Pool
}

// NewPsg создает новый экземпляр Psg.
func NewPsg(psgAddr string, login, password string) *Psg {
	return
}

// RecordAdd добавляет новую запись в базу данных.
func (p *Psg) RecordAdd(record Record) (int64, error) {
	// TODO: Реализовать добавление записи
	return 0, nil
}

// RecordsGet возвращает записи из базы данных на основе предоставленных полей Record.
func (p *Psg) RecordsGet(record Record) ([]Record, error) {
	// TODO: Реализовать поиск записей
	return nil, nil
}

// RecordUpdate обновляет существующую запись в базе данных по номеру телефона.
func (p *Psg) RecordUpdate(record Record) error {
	// TODO: Реализовать обновление записи
	return nil
}

// RecordDeleteByPhone удаляет запись из базы данных по номеру телефона.
func (p *Psg) RecordDeleteByPhone(phone string) error {
	// TODO: Реализовать удаление записи по номеру телефона
	return nil
}
`
