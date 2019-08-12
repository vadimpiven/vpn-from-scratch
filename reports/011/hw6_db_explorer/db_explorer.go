package main

import (
	"bytes"
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
	Status int         `json:"-"`
	Error  string      `json:"error,omitempty"`
	Result interface{} `json:"response,omitempty"`
}

// write performs buffered write of entire string.
func write(out io.Writer, str []byte) (err error) {
	var n int
	for n, err = out.Write(str); err == nil && n < len(str); {
		str = str[n:]
	}
	return
}

// errorOccurred function checks for an error and handles it if any occurs.
func errorOccurred(err error, w http.ResponseWriter) bool {
	if err != nil {
		writeResponse(Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, w)
		return true
	}
	return false
}

// writeResponse writes JSON response and sets status code.
func writeResponse(r Response, w http.ResponseWriter) {
	if res, err := json.Marshal(r); !errorOccurred(err, w) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(r.Status)
		if err = write(w, res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

type (
	Table struct {
		Fields []reflect.StructField
		Struct reflect.Type
	}
)

// getTableList creates the list of existing tables.
func getTableList(db *sql.DB) []string {
	var tables []string
	rows, err := db.Query("SHOW TABLES;")
	check(err)
	defer Close(rows)
	for rows.Next() {
		var table string
		check(rows.Scan(&table))
		tables = append(tables, table)
	}
	return tables
}

// getTables builds the description and receiver for each table in database.
func getTables(db *sql.DB) (tableList []string, tables map[string]Table) {
	tableList = getTableList(db)
	tables = make(map[string]Table)
	for _, table := range tableList {
		rows, err := db.Query("SELECT * FROM " + table + " LIMIT 0;")
		check(err)
		var tableFields []reflect.StructField
		fields, err := rows.ColumnTypes()
		check(err)
		for _, field := range fields {
			ty := field.ScanType()
			// BUG(vadimpiven): sql.RawBytes could represent not only string type.
			// ScanType returns sql.RawBytes in case of Nullable with this driver,
			// hopefully only strings could be null with this test data.
			// Type *string is used to make holding nil value possible.
			if ty == reflect.TypeOf((*sql.RawBytes)(nil)).Elem() {
				ty = reflect.TypeOf((*string)(nil))
			}
			tableFields = append(tableFields, reflect.StructField{
				Name: strings.Title(field.Name()),
				Type: ty,
				Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s"`, field.Name())),
			})
		}
		Close(rows)
		//rows, err := db.Query("SHOW COLUMNS FROM " + table)
		//check(err)
		//var tableFields []reflect.StructField
		//var fName, fType, fNull, fPri, fDefault, fExtra sql.NullString
		//for rows.Next() {
		//	err = rows.Scan(&fName, &fType, &fNull, &fPri, &fDefault, &fExtra)
		//	fName, fType, fNull := fName.String, fType.String, fNull.String
		//	check(err)
		//	switch {
		//	case strings.HasPrefix(fType, "int"):
		//		tableFields = append(tableFields, reflect.StructField{
		//					Name: strings.Title(fName),
		//					Type: reflect.TypeOf((int64)(0)),
		//					Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s"`, fName)),
		//				})
		//	case strings.HasPrefix(fType, "varchar"), fType == "text":
		//		switch fNull {
		//		case "NO":
		//			tableFields = append(tableFields, reflect.StructField{
		//				Name: strings.Title(fName),
		//				Type: reflect.TypeOf(""),
		//				Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s"`, fName)),
		//			})
		//		default:
		//			tableFields = append(tableFields, reflect.StructField{
		//				Name: strings.Title(fName),
		//				Type: reflect.TypeOf((*string)(nil)),
		//				Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s"`, fName)),
		//			})
		//		}
		//	}
		//}
		//Close(rows)

		tables[table] = Table{
			Fields: tableFields,
			Struct: reflect.StructOf(tableFields),
		}
	}
	return
}

// selectFromDB performs SELECT query from given table.
func selectFromDB(db *sql.DB, table string, resType reflect.Type, limit, offset uint64) (interface{}, error) {
	var res []interface{}
	rows, err := db.Query("SELECT * FROM " + table + " LIMIT ? OFFSET ?;", limit, offset)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		temp := reflect.New(resType).Elem()
		var tempFields []interface{}
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
	return res, nil
}

// selectByID performs SELECT by ID from given table.
func selectByID(db *sql.DB, table, idField string, resType reflect.Type, id uint64) (interface{}, error) {
	rows := db.QueryRow("SELECT * FROM " + table + " WHERE " + idField + " = ?;", id)
	res := reflect.New(resType).Elem()
	var tempFields []interface{}
	for i := 0; i < res.NumField(); i++ {
		tempFields = append(tempFields, res.Field(i).Addr().Interface())
	}
	err := rows.Scan(tempFields...)
	if err != nil {
		return nil, err
	}
	return res.Interface(), nil
}

// insert performs INSERT query to given table.
func insert(db *sql.DB, table string, data reflect.Value) (id int64, err error) {
	var args []interface{}
	q := "INSERT INTO " + table + "("
	// i == 0 is skipped because this program assumes that first field is ID identifier (autoincrement)
	for i := 1; i < data.NumField(); i++ {
		if i > 1 {
			q += ", "
		}
		q += data.Type().Field(i).Name
		args = append(args, data.Field(i).Interface())
	}
	q += ") VALUES ("
	for i := 0; i < len(args); i++ {
		if i > 0 {
			q += ", "
		}
		q += "?"
	}
	q += ");"
	res, Err := db.Exec(q, args...)
	for Err != nil && strings.Contains(Err.Error(), "cannot be null") {
		q = "INSERT INTO " + table + "("
		for i := 1; i < data.NumField(); i++ {
			if i > 1 {
				q += ", "
			}
			field := strings.ToLower(data.Type().Field(i).Name)
			q += field
			if strings.Contains(Err.Error(), field) {
				args[i-1] = ""
			}
		}
		q += ") VALUES ("
		for i := 0; i < len(args); i++ {
			if i > 0 {
				q += ", "
			}
			q += "?"
		}
		q += ");"
		res, Err = db.Exec(q, args...)
	}
	if Err != nil {
		return
	}
	id, err = res.LastInsertId()
	return
}

// updateById performs UPDATE query to given table with certain id.
func updateById(db *sql.DB, table string, data reflect.Value, postData []byte, id uint64) error {
	var args []interface{}
	q := "UPDATE " + table + " SET "
	c := 0
	// i == 0 is errored because this program assumes that first field is ID identifier (autoincrement)
	for i := 0; i < data.NumField(); i++ {
		field := strings.ToLower(data.Type().Field(i).Name)
		if !bytes.Contains(postData, []byte(field)) {
			continue
		}
		if i == 0 {
			return fmt.Errorf("field %s have invalid type", field)
		}
		if c > 0 {
			q += ", "
		}
		q += field + " = ?"
		args = append(args, data.Field(i).Interface())
		c++
	}
	if c == 0 {
		// error not specified by tests
		return fmt.Errorf("no data provided")
	}
	q += " WHERE " + data.Type().Field(0).Name + " = ?;"
	args = append(args, id)
	_, err := db.Exec(q, args...)
	if err != nil {
		for i := 1; i < data.NumField(); i++ {
			field := strings.ToLower(data.Type().Field(i).Name)
			if strings.Contains(err.Error(), field) {
				return fmt.Errorf("field %s have invalid type", field)
			}
		}
	}
	return err
}

// deleteByID deletes from given table row with given ID.
func deleteByID(db *sql.DB, table, idField string, id uint64) (int64, error) {
	res, err := db.Exec("DELETE FROM " + table + " WHERE " + idField + " = ?;", id)
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
	tableList, tables := getTables(db)
	return func(w http.ResponseWriter, r *http.Request) {
		method, url, query := r.Method, strings.Split(strings.Trim(strings.ToLower(r.URL.Path), "/"), "/"), r.URL.Query()
		if url[0] == "" {
			writeResponse(Response{
				Status: http.StatusOK,
				Result: map[string]interface{}{"tables": tableList},
			}, w)
			return
		}
		table, ok := tables[url[0]]
		if !ok {
			writeResponse(Response{
				Status: http.StatusNotFound,
				Error: "unknown table",
			}, w)
			return
		}
		switch len(url) {
		case 1:
			switch method {
			case "PUT":
				buf, err := ioutil.ReadAll(r.Body)
				if errorOccurred(err, w) {
					return
				}
				if err = r.Body.Close(); errorOccurred(err, w) {
					return
				}
				body := reflect.New(table.Struct).Elem()
				if err := json.Unmarshal(buf, body.Addr().Interface()); errorOccurred(err, w) {
					return
				}
				id, err := insert(db, url[0], body)
				if errorOccurred(err, w) {
					return
				}
				writeResponse(Response{
					Status: http.StatusOK,
					Result: map[string]interface{}{strings.ToLower(table.Fields[0].Name): id},
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
				res, err := selectFromDB(db, url[0], table.Struct, limit, offset)
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
			// error unspecified by tests
			if err != nil {
				writeResponse(Response{
					Status: http.StatusBadRequest,
					Error: "id has invalid type",
				}, w)
				return
			}
			switch method {
			case "POST":
				buf, err := ioutil.ReadAll(r.Body)
				if errorOccurred(err, w) {
					return
				}
				if err = r.Body.Close(); errorOccurred(err, w) {
					return
				}
				body := reflect.New(table.Struct).Elem()
				err = json.Unmarshal(buf, body.Addr().Interface())
				if newErr, ok := err.(*json.UnmarshalTypeError); err != nil && ok {
					writeResponse(Response{
						Status: http.StatusBadRequest,
						Error: fmt.Sprint("field ", strings.ToLower(newErr.Field), " have invalid type"),
					}, w)
					return
				} else if errorOccurred(err, w) {
					return
				}
				err = updateById(db, url[0], body, buf, id)
				if err != nil {
					writeResponse(Response{
						Status: http.StatusBadRequest,
						Error: err.Error(),
					}, w)
					return
				}
				writeResponse(Response{
					Status: http.StatusOK,
					Result: map[string]interface{}{"updated": 1},
				}, w)
				return
			case "DELETE":
				num, err := deleteByID(db, url[0], table.Fields[0].Name, id)
				if err != nil {
					writeResponse(Response{
						Status: http.StatusBadRequest,
						Error: err.Error(),
					}, w)
					return
				}
				writeResponse(Response{
					Status: http.StatusOK,
					Result: map[string]interface{}{"deleted": num},
				}, w)
				return
			default: // GET
				res, err := selectByID(db, url[0], table.Fields[0].Name, table.Struct, id)
				if err != nil {
					writeResponse(Response{
						Status: http.StatusNotFound,
						Error: "record not found",
					}, w)
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
			writeResponse(Response{
				Status: http.StatusBadRequest,
				Error: "too many arguments",
			}, w)
			return
		}
	}, nil
}
