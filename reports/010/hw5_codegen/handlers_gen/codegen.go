package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

const Template = `// GENERATED, DO NOT MODIFY

package {{.Package}}

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// Response represents the response structure.
type Response struct {
	Status int         ` + "`json:" + `"-"` + "`" + `
	Error  string      ` + "`json:" + `"error"` + "`" + `
	Result interface{} ` + "`json:" + `"response,omitempty"` + "`" + `
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
		switch err := err.(type) {
		case ApiError:
			writeResponse(Response{
				Status: err.HTTPStatus,
				Error:  err.Error(),
			}, w)
		default:
			writeResponse(Response{
				Status: http.StatusInternalServerError,
				Error:  err.Error(),
			}, w)
		}
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

// isPost checks if request is POST and throws error if not.
func isPost(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		writeResponse(Response{
			Status: http.StatusNotAcceptable,
			Error:  "bad method",
		}, w)
		return false
	}
	return true
}

// isAuthorized checks if authorization header is provided and throws error if not.
func isAuthorized(w http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get("X-Auth") != "100500" {
		writeResponse(Response{
			Status: http.StatusForbidden,
			Error:  "unauthorized",
		}, w)
		return false
	}
	return true
}

// getString reads string from request end throws error if value required but not provided.
func getString(name string, data map[string][]string, required bool) (string, error) {
	res := data[name]
	if len(res) == 0 {
		if required {
			return "", ApiError{
				HTTPStatus: http.StatusBadRequest,
				Err:        fmt.Errorf("%s must me not empty", name),
			}
		}
		return "", nil
	}
	return res[0], nil
}

// getInt reads int from request end throws error if value required but not provided.
func getInt(name string, data map[string][]string, required bool) (int, error) {
	tmp := data[name]
	if len(tmp) == 0 {
		if required {
			return 0, ApiError{
				HTTPStatus: http.StatusBadRequest,
				Err:        fmt.Errorf("%s must me not empty", name),
			}
		}
		return 0, nil
	}
	res, err := strconv.Atoi(tmp[0])
	if err != nil {
		return 0, ApiError{
			HTTPStatus: http.StatusBadRequest,
			Err:        fmt.Errorf("%s must be int", name),
		}
	}
	return res, nil
}

// inEnum sets value to default if empty string, checks whether value is in enum and throws error if not.
func inEnum(name string, value string, enum []string, defVal string) (string, error) {
	if value == "" && defVal != "" {
		return defVal, nil
	}
	for _, str := range enum {
		if value == str {
			return value, nil
		}
	}
	return "", ApiError{
		HTTPStatus: http.StatusBadRequest,
		Err:        fmt.Errorf("%s must be one of [%s]", name, strings.Join(enum, ", ")),
	}
}

// checkMinString checks whether string length is greater then minimum required.
func checkMinString(name string, value string, min int) error {
	if len(value) < min {
		return ApiError{
			HTTPStatus: http.StatusBadRequest,
			Err:        fmt.Errorf("%s len must be >= %d", name, min),
		}
	}
	return nil
}

// checkMaxString checks whether string length is smaller then minimum required.
func checkMaxString(name string, value string, max int) error {
	if len(value) > max {
		return ApiError{
			HTTPStatus: http.StatusBadRequest,
			Err:        fmt.Errorf("%s len must be <= %d", name, max),
		}
	}
	return nil
}

// checkMinInt checks whether value is greater then minimum required.
func checkMinInt(name string, value int, min int) error {
	if value < min {
		return ApiError{
			HTTPStatus: http.StatusBadRequest,
			Err:        fmt.Errorf("%s must be >= %d", name, min),
		}
	}
	return nil
}

// checkMaxInt checks whether value is smaller then minimum required.
func checkMaxInt(name string, value int, max int) error {
	if value > max {
		return ApiError{
			HTTPStatus: http.StatusBadRequest,
			Err:        fmt.Errorf("%s must be <= %d", name, max),
		}
	}
	return nil
}{{range $name, $methods := .ApiList}}

// ServeHTTP defines functions serving different URLs for {{$name}} method handlers.
func (srv *{{$name}}) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path { {{- range $methods}}
	case "{{.URL}}":
		srv.handler{{.Name}}(w, r){{end}}
	default:
		writeResponse(Response{
			Status: http.StatusNotFound,
			Error:  "unknown method",
		}, w)
	}
}{{range $methods}}

// handler{{.Name}} is a handler for {{$name}}.{{.Name}} function.
func (srv *{{$name}}) handler{{.Name}}(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); {{if .Post}}!isPost(w, r) || {{end}}{{if .Auth}}!isAuthorized(w, r) || {{end}}errorOccurred(err, w) {
		return
	}
	data := r.{{if .Post}}Post{{end}}Form{{range .Params.Fields}}
	{{.Name}}, err := get{{if .IsString}}String{{else}}Int{{end}}("{{.ParamName}}", data, {{.Required}})
	if errorOccurred(err, w){{if .Min.IsSet}} ||
	   errorOccurred(checkMin{{if .IsString}}String{{else}}Int{{end}}("{{.ParamName}}", {{.Name}}, {{.Min.Edge}}), w){{end}}{{if .Max.IsSet}} ||
	   errorOccurred(checkMax{{if .IsString}}String{{else}}Int{{end}}("{{.ParamName}}", {{.Name}}, {{.Max.Edge}}), w){{end}} {
		return
	}{{if .Enum}}
	{{.Name}}, err = inEnum("{{.ParamName}}", {{.Name}}, {{printf "%#v" .Enum}}, "{{.Default}}")
	if errorOccurred(err, w) {
		return
	}{{end}}{{end}}
	res, err := srv.{{.Name}}(r.Context(), {{.Params.Name}}{ {{- range .Params.Fields}}
		{{.Name}}: {{.Name}},{{end}}
	})
	if !errorOccurred(err, w) {
		writeResponse(Response{
			Status: http.StatusOK,
			Result: res,
		}, w)
	}
}{{end}}{{end}}
`

// check function checks for an error and exits if any occurs.
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

// Flush allows not to ignore potential error while calling `Flush()` inside `defer`.
func Flush(w *bufio.Writer) {
	err := w.Flush()
	check(err)
}

type (
	// MethodJSON is used to unmarshal JSON describing Method requirements.
	MethodJSON struct {
		URL	   string `json:"url"`
		Auth   bool   `json:"auth"`
		Method string `json:"method"`
	}
	// All structures below are used to fill in template.
	Boundary struct {
		IsSet bool // default is false
		Edge  int
	}
	Field struct {
		Name	      string // name of structure field
		IsString  bool // if not string - it's int
		Required  bool
		ParamName string // name used to unmarshal incoming data
		Enum	      []string // used only as Enum = ..., no need for initialisation
		Default   string // if empty string - no default
		Min	      Boundary
		Max	      Boundary
	}
	Params struct {
		Name   string // name of handler second argument structure type
		Fields []Field // 0-length slice is created automatically
	}
	Method struct {
		Name   string // original method name
		URL	   string
		Auth   bool
		Post   bool
		Params Params
	}
	TemplateArgs struct {
		Package string // package name
		ApiList map[string][]Method // 0-length map MUST be created manually, 0-length slice inside map is created automatically
	}
)

func main() {
	var debug = flag.Bool("debug", false, "print debug output")
	if flag.Parse(); *debug {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(ioutil.Discard)
	}

	src, err := parser.ParseFile(token.NewFileSet(), flag.Arg(0), nil, parser.ParseComments)
	check(err)
	args := TemplateArgs{
		Package: src.Name.Name,
		ApiList: make(map[string][]Method),
	}

	f, _ := os.Create(flag.Arg(1))
	defer Close(f)

	w := bufio.NewWriter(f)
	defer Flush(w)

	log.Printf("SEARCH for methods to wrap\n")
	for _, decl := range src.Decls {
		currFunc, ok := decl.(*ast.FuncDecl)
		if !ok {
			log.Printf("SKIP %T is not *ast.FuncDecl\n", decl)
			continue
		}

		if currFunc.Doc == nil {
			log.Printf("SKIP function %#v doesnt have comments\n", currFunc.Name.Name)
			continue
		}

		needCodegen := false
		for _, comment := range currFunc.Doc.List {
			if strings.HasPrefix(comment.Text, "// apigen:api ") {
				needCodegen = true

				js := &MethodJSON{}
				err := json.Unmarshal([]byte(strings.TrimPrefix(comment.Text, "// apigen:api ")), js)
				check(err)

				api := currFunc.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name
				args.ApiList[api] = append(args.ApiList[api], Method{
					Name: currFunc.Name.Name,
					URL:  js.URL,
					Auth: js.Auth,
					Post: js.Method == "POST",
					Params: Params{
						Name: currFunc.Type.Params.List[1].Type.(*ast.Ident).Name,
					},
				})

				log.Printf("FOUND function %#v\n", currFunc.Name.Name)
				break
			}
		}
		if !needCodegen {
			log.Printf("SKIP function %#v doesnt have apigen mark\n", currFunc.Name.Name)
		}
	}

	log.Printf("SEARCH for Params.Fields and parse their requirements\n")
	for api, methods := range args.ApiList {
		log.Printf("SEARCH for %s API methods requirements\n", api)
		for i, method := range methods {
			log.Printf("SEARCH for %s.%s method requirements\n", api, method.Name)
		MethodLoop:
			for _, decl := range src.Decls {
				gen, ok := decl.(*ast.GenDecl)
				if !ok {
					log.Printf("SKIP %T is not *ast.GenDecl\n", f)
					continue
				}
				for _, spec := range gen.Specs {
					currType, ok := spec.(*ast.TypeSpec)
					if !ok {
						log.Printf("SKIP %T is not ast.TypeSpec\n", spec)
						continue
					}

					currStruct, ok := currType.Type.(*ast.StructType)
					if !ok {
						log.Printf("SKIP %T is not ast.StructType\n", currStruct)
						continue
					}

					if currType.Name.Name != method.Params.Name {
						log.Printf("SKIP %s struct found, searching for %s struct\n", currType.Name.Name, method.Params.Name)
						continue
					}

					log.Printf("FOUND structure %s for method %s.%s\n", currType.Name.Name, api, method.Name)
				FieldsLoop:
					for _, field := range currStruct.Fields.List {
						if field.Tag != nil {
							tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
							tagVal := tag.Get("apivalidator")
							if tagVal == "-" {
								log.Printf("SKIP %s field of structure %s\n", field.Names[0].Name, currType.Name.Name)
								continue FieldsLoop
							}

							newField := Field{
								Name:	  field.Names[0].Name,
								IsString:  field.Type.(*ast.Ident).Name == "string", // only int and string types are expected
								ParamName: strings.ToLower(field.Names[0].Name),
							}

							log.Printf("PROCESSING %s field of structure %s\n", field.Names[0].Name, currType.Name.Name)
							for _, val := range strings.Split(tagVal, ",") {
								switch {
								case val == "required":
									log.Printf("FOUND 'required' tag of field %s\n", field.Names[0].Name)
									newField.Required = true
								case strings.HasPrefix(val, "paramname="):
									log.Printf("FOUND 'paramname' tag of field %s\n", field.Names[0].Name)
									newField.ParamName = strings.TrimPrefix(val, "paramname=")
								case strings.HasPrefix(val, "enum="):
									log.Printf("FOUND 'enum' tag of field %s\n", field.Names[0].Name)
									newField.Enum = strings.Split(strings.TrimPrefix(val, "enum="), "|")
								case strings.HasPrefix(val, "default="):
									log.Printf("FOUND 'default' tag of field %s\n", field.Names[0].Name)
									newField.Default = strings.TrimPrefix(val, "default=")
								case strings.HasPrefix(val, "min="):
									log.Printf("FOUND 'min' tag of field %s\n", field.Names[0].Name)
									newField.Min.IsSet = true
									newField.Min.Edge, _ = strconv.Atoi(strings.TrimPrefix(val, "min="))
								case strings.HasPrefix(val, "max="):
									log.Printf("FOUND 'max' tag of field %s\n", field.Names[0].Name)
									newField.Max.IsSet = true
									newField.Max.Edge, _ = strconv.Atoi(strings.TrimPrefix(val, "max="))
								default:
									log.Printf("SKIP %s tag of field %s\n", val, field.Names[0].Name)
									continue
								}
							}
							methods[i].Params.Fields = append(methods[i].Params.Fields, newField)
						}
					}

					break MethodLoop
				}
			}
		}
	}

	tmpl, err := template.New("generated file template").Parse(Template)
	check(err)
	err = tmpl.Execute(w, args)
	check(err)
}
