package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// check exits the program on error.
func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// Close allows not to ignore potential error while calling `Close()` inside `defer`.
func Close(c io.Closer) {
	err := c.Close()
	check(err)
}

// Response represents the response structure.
type Response struct {
	Status int
	Result map[string]interface{}
}

// write performs buffered write of entire string.
func write(out io.Writer, str []byte) (err error) {
	var n int
	for n, err = out.Write(str); err == nil && n < len(str); {
		str = str[n:]
	}
	return
}

type ApiError struct {
	Status int
	Err    error
}

func (ae ApiError) Error() string {
	return ae.Err.Error()
}

// errorOccurred function checks for an error and handles it if any occurs.
func errorOccurred(err error, w http.ResponseWriter) bool {
	if err != nil {
		switch err := err.(type) {
		case ApiError:
			writeResponse(Response{
				Status: err.Status,
				Result: map[string]interface{}{"error": err.Error()},
			}, w)
		default:
			writeResponse(Response{
				Status: http.StatusInternalServerError,
				Result: map[string]interface{}{"error": err.Error()},
			}, w)
		}
		return true
	}
	return false
}

// writeResponse writes JSON response and sets status code.
func writeResponse(r Response, w http.ResponseWriter) {
	var data interface{}
	if r.Status == http.StatusOK {
		data = map[string]interface{}{"response":r.Result}
	} else {
		data = r.Result
	}
	if res, err := json.Marshal(data); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(r.Status)
		if err = write(w, res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ~~~~~~~~~~~~~~~~~ CODE FROM https://github.com/guregu/null/blob/v3.4.0/string.go ~~~~~~~~~~~~~~~~~
// NullString is a nullable string. It supports SQL and JSON serialization.
// It will marshal to null if null. Blank string input will be considered null.
type NullString struct {
	sql.NullString
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (s NullString) ValueOrZero() string {
	if !s.Valid {
		return ""
	}
	return s.String
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string and null input. Blank string input does not produce a null NullString.
// It also supports unmarshalling a sql.NullString.
func (s *NullString) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case string:
		s.String = x
	case map[string]interface{}:
		err = json.Unmarshal(data, &s.NullString)
	case nil:
		s.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.NullString", reflect.TypeOf(v).Name())
	}
	s.Valid = err == nil
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this NullString is null.
func (s NullString) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.String)
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string when this NullString is null.
func (s NullString) MarshalText() ([]byte, error) {
	if !s.Valid {
		return []byte{}, nil
	}
	return []byte(s.String), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null NullString if the input is a blank string.
func (s *NullString) UnmarshalText(text []byte) error {
	s.String = string(text)
	s.Valid = s.String != ""
	return nil
}

// SetValid changes this NullString's value and also sets it to be non-null.
func (s *NullString) SetValid(v string) {
	s.String = v
	s.Valid = true
}

// Ptr returns a pointer to this NullString's value, or a nil pointer if this NullString is null.
func (s NullString) Ptr() *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}

// IsZero returns true for null strings, for potential future omitempty support.
func (s NullString) IsZero() bool {
	return !s.Valid
}
// ~~~~~~~~~~~~~~~~~ END OF CODE FROM https://github.com/guregu/null/blob/v3.4.0/string.go ~~~~~~~~~~~~~~~~~

type (
	Field struct {
		Name    string
		Type    reflect.Type  // could be normal type or sql.Null...
		Default reflect.Value // always *Field.Type as if default is not set - it's nil
	}
	Table struct {
		Name          string
		Fields        []Field
		IDIndex       uint64
		containerType reflect.Type
	}
)

// NewTable creates new table initialising private fields.
func NewTable(name string, fields []Field, idIndex uint64) Table {
	structFields := make([]reflect.StructField, 0, 2)
	for _, field := range fields {
		structFields = append(structFields, reflect.StructField{
			Name: strings.Title(field.Name), // first letter should be capital to make field public
			Type: field.Type,
			Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s"`, field.Name)),
		})
	}
	return Table{
		Name:          name,
		Fields:        fields,
		IDIndex:       idIndex,
		containerType: reflect.StructOf(structFields),
	}
}

// Container returns a structure which could be used to unmarshal json (also sets up defaults).
func (t Table) Container() (structure reflect.Value) {
	structure = reflect.New(t.containerType).Elem()
	for i, field := range t.Fields {
		if !field.Default.IsNil() {
			structure.Field(i).Set(field.Default.Elem())
		}
	}
	return
}

// getTableList creates the list of existing tables.
func getTableList(db *sql.DB) (tables []string) {
	rows, err := db.Query("SHOW TABLES;")
	check(err)
	defer Close(rows)
	for rows.Next() {
		var table string
		check(rows.Scan(&table))
		tables = append(tables, table)
	}
	return
}

// getTables builds the description and receiver for each table in database.
func getTables(db *sql.DB) (tableList []string, tableInfo map[string]Table) {
	tableList = getTableList(db)
	tableInfo = make(map[string]Table)
	for _, tableName := range tableList {
		rows, err := db.Query("SHOW COLUMNS FROM " + tableName)
		check(err)
		tableFields := make([]Field, 0, 2)
		var idIndex uint64
		var fName, fType, fNull, fKey, fExtra string
		var fDefault NullString
		for i := uint64(0); rows.Next(); i++ {
			err = rows.Scan(&fName, &fType, &fNull, &fKey, &fDefault, &fExtra)
			check(err)
			if fKey == "PRI" {
				idIndex = i // expected only one primary key which is also auto incremented
			}
			var rType reflect.Type
			switch { // type
			case strings.HasPrefix(fType, "int"):
				rType = reflect.TypeOf((*uint64)(nil)).Elem()
			case strings.HasPrefix(fType, "varchar"), fType == "text":
				switch fNull {
				case "NO":
					rType = reflect.TypeOf((*string)(nil)).Elem()
				default:                                              // YES
					rType = reflect.TypeOf((*NullString)(nil)).Elem() // NullString is json-compatible sql.NullString
				}
			}
			rDefault := reflect.New(rType)
			if fDefault.Valid {
				switch { // default
				case strings.HasPrefix(fType, "int"):
					val, err := strconv.ParseUint(fDefault.String, 10, 64)
					check(err)
					rDefault.SetUint(val)
				case strings.HasPrefix(fType, "varchar"), fType == "text":
					switch fNull {
					case "NO":
						rDefault.SetString(fDefault.String)
					default: // YES
						rDefault.Set(reflect.ValueOf(fDefault))
					}
				}
			}
			tableFields = append(tableFields, Field{
				Name:    fName,
				Type:    rType,
				Default: rDefault,
			})
		}
		Close(rows)
		tableInfo[tableName] = NewTable(tableName, tableFields, idIndex)
	}
	return
}

// selectFromDB performs SELECT query from given table.
func selectFromDB(db *sql.DB, t Table, limit, offset uint64) (res []interface{}, err error) {
	rows, err := db.Query("SELECT * FROM " + t.Name + " LIMIT ? OFFSET ?;", limit, offset)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		temp := t.Container()
		tempFields := make([]interface{}, 0, 2)
		for i := 0; i < temp.NumField(); i++ {
			tempFields = append(tempFields, temp.Field(i).Addr().Interface())
		}
		err = rows.Scan(tempFields...)
		if err != nil {
			return nil, err
		}
		res = append(res, temp.Interface())
	}
	if err = rows.Close(); err != nil {
		return nil, err
	}
	if len(res) == 0 {
		// error unspecified by tests
		return nil, ApiError{
			Status: http.StatusNotFound,
			Err:    fmt.Errorf("records not found"),
		}
	}
	return res, nil
}

// selectByID performs SELECT by ID from given table.
func selectByID(db *sql.DB, t Table, id uint64) (res interface{}, err error) {
	rows := db.QueryRow("SELECT * FROM " + t.Name + " WHERE " + t.Fields[t.IDIndex].Name + " = ?;", id)
	temp := t.Container()
	tempFields := make([]interface{}, 0, 2)
	for i := 0; i < temp.NumField(); i++ {
		tempFields = append(tempFields, temp.Field(i).Addr().Interface())
	}
	if err = rows.Scan(tempFields...); err == nil {
		return temp.Interface(), nil
	} else if err == sql.ErrNoRows {
		return nil, ApiError{
			Status: http.StatusNotFound,
			Err:    fmt.Errorf("record not found"),
		}
	}
	return
}

// insert performs INSERT query to given table.
func insert(db *sql.DB, t Table, body []byte) (id int64, err error) {
	temp := t.Container()
	if err = unmarshal(body, t, &temp, false); err != nil {
		return
	}
	tempFields := make([]interface{}, 0, 2)
	q := "INSERT INTO " + t.Name + "("
	for i := 0; i < len(t.Fields); i++ {
		if uint64(i) == t.IDIndex {
			continue
		}
		if len(tempFields) > 0 {
			q += ", "
		}
		q += t.Fields[i].Name
		tempFields = append(tempFields, temp.Field(i).Interface())
	}
	q += ") VALUES ("
	for i := 0; i < len(tempFields); i++ {
		if i > 0 {
			q += ", "
		}
		q += "?"
	}
	q += ");"
	res, err := db.Exec(q, tempFields...)
	if err != nil {
		return
	}
	id, err = res.LastInsertId()
	return
}

func unmarshal(body []byte, t Table, container *reflect.Value, errorOnIdChange bool) error {
	var m map[string]interface{}
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	for i, field := range t.Fields {
		val, ok := m[field.Name]
		if !ok {
			continue
		}
		if (val == nil && field.Type != reflect.TypeOf((*NullString)(nil)).Elem()) ||
			(errorOnIdChange && uint64(i) == t.IDIndex) {
			return ApiError{
				Status: http.StatusBadRequest,
				Err:    fmt.Errorf("field %s have invalid type", field.Name),
			}
		}
		switch val.(type) {
		case float64:
			if field.Type != reflect.TypeOf((*uint64)(nil)).Elem() {
				return ApiError{
					Status: http.StatusBadRequest,
					Err:    fmt.Errorf("field %s have invalid type", field.Name),
				}
			}
		case string:
			if field.Type == reflect.TypeOf((*uint64)(nil)).Elem()  {
				return ApiError{
					Status: http.StatusBadRequest,
					Err:    fmt.Errorf("field %s have invalid type", field.Name),
				}
			}
		}
	}
	return json.Unmarshal(body, container.Addr().Interface())
}

// updateById performs UPDATE query to given table with certain id.
func updateById(db *sql.DB, t Table, id uint64, body []byte) error {
	row := db.QueryRow("SELECT * FROM " + t.Name + " WHERE " + t.Fields[t.IDIndex].Name + " = ?;", id)
	temp := t.Container()
	tempFields := make([]interface{}, 0, 2)
	for i := 0; i < temp.NumField(); i++ {
		tempFields = append(tempFields, temp.Field(i).Addr().Interface())
	}
	if err := row.Scan(tempFields...); err != nil {
		if err == sql.ErrNoRows {
			return ApiError{
				Status: http.StatusNotFound,
				Err:    fmt.Errorf("record not found"),
			}
		}
		return err
	}
	if err := unmarshal(body, t, &temp, true); err != nil {
		return err
	}
	tempFields = tempFields[:0]
	q := "UPDATE " + t.Name + " SET "
	for i := 0; i < len(t.Fields); i++ {
		if uint64(i) == t.IDIndex {
			continue
		}
		if len(tempFields) > 0 {
			q += ", "
		}
		q += t.Fields[i].Name + " = ?"
		tempFields = append(tempFields, temp.Field(i).Interface())
	}
	if len(tempFields) == 0 {
		// error not specified by tests
		return fmt.Errorf("no data provided")
	}
	q += " WHERE " + t.Fields[t.IDIndex].Name + " = ?;"
	tempFields = append(tempFields, id)
	_, err := db.Exec(q, tempFields...)
	return err
}

// deleteByID deletes from given table row with given ID.
func deleteByID(db *sql.DB, t Table, id uint64) (int64, error) {
	res, err := db.Exec("DELETE FROM " + t.Name + " WHERE " + t.Fields[t.IDIndex].Name + " = ?;", id)
	if err != nil {
		return 0, err
	}
	num, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return num, nil
}

// NewDbExplorer scans the database and returns initialised handler function.
func NewDbExplorer(db *sql.DB) (http.HandlerFunc, error) {
	tableList, tableInfo := getTables(db)
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		url := strings.Split(strings.Trim(strings.ToLower(r.URL.Path), "/"), "/")
		query := r.URL.Query()
		var body []byte
		var err error
		if method != "GET" {
			body, err = ioutil.ReadAll(r.Body)
			if errorOccurred(err, w) || errorOccurred(r.Body.Close(), w) {
				return
			}
		}

		if url[0] == "" { // GET /
			writeResponse(Response{
				Status: http.StatusOK,
				Result: map[string]interface{}{"tables": tableList},
			}, w)
			return
		}

		table, ok := tableInfo[url[0]] // URL aka /$table...
		if !ok {
			errorOccurred(ApiError{
				Status: http.StatusNotFound,
				Err:    fmt.Errorf("unknown table"),
			}, w)
			return
		}

		switch len(url) {
		case 1:
			switch method {
			case "PUT":
				id, err := insert(db, table, body)
				if errorOccurred(err, w) {
					return
				}
				writeResponse(Response{
					Status: http.StatusOK,
					Result: map[string]interface{}{table.Fields[table.IDIndex].Name: id},
				}, w)
				return
			default: // GET
				limit, offset := uint64(5), uint64(0)
				if temp, ok := query["limit"]; ok {
					if uTemp, err := strconv.ParseUint(temp[0], 10, 64); err == nil {
						limit = uTemp
					}
				}
				if temp, ok := query["offset"]; ok {
					if uTemp, err := strconv.ParseUint(temp[0], 10, 64); err == nil {
						offset = uTemp
					}
				}
				res, err := selectFromDB(db, table, limit, offset)
				if errorOccurred(err, w) {
					return
				}
				writeResponse(Response{
					Status: http.StatusOK,
					Result: map[string]interface{}{"records": res},
				}, w)
				return
			}
		case 2:
			id, err := strconv.ParseUint(url[1], 10, 64)
			if errorOccurred(err, w) {
				return
			}
			switch method {
			case "POST":
				err = updateById(db, table, id, body)
				if errorOccurred(err, w) {
					return
				}
				writeResponse(Response{
					Status: http.StatusOK,
					Result: map[string]interface{}{"updated": 1},
				}, w)
				return
			case "DELETE":
				num, err := deleteByID(db, table, id)
				if errorOccurred(err, w) {
					return
				}
				writeResponse(Response{
					Status: http.StatusOK,
					Result: map[string]interface{}{"deleted": num},
				}, w)
				return
			default: // GET
				res, err := selectByID(db, table, id)
				if errorOccurred(err, w) {
					return
				}
				writeResponse(Response{
					Status: http.StatusOK,
					Result: map[string]interface{}{"record": res},
				}, w)
				return
			}
		default:
			// error unspecified by tests
			errorOccurred(ApiError{
				Status: http.StatusBadRequest,
				Err:    fmt.Errorf("too many arguments"),
			}, w)
			return
		}
	}, nil
}
