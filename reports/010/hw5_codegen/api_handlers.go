// GENERATED, DO NOT MODIFY

package main

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
    Status int         `json:"-"`
    Error  string      `json:"error"`
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
}

// ServeHTTP defines functions serving different URLs for MyApi method handlers.
func (srv *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch r.URL.Path {
    case "/user/profile":
        srv.handlerProfile(w, r)
    case "/user/create":
        srv.handlerCreate(w, r)
    default:
        writeResponse(Response{
            Status: http.StatusNotFound,
            Error:  "unknown method",
        }, w)
    }
}

// handlerProfile is a handler for MyApi.Profile function.
func (srv *MyApi) handlerProfile(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); errorOccurred(err, w) {
        return
    }
    data := r.Form
    Login, err := getString("login", data, true)
    if errorOccurred(err, w) {
        return
    }
    res, err := srv.Profile(r.Context(), ProfileParams{
        Login: Login,
    })
    if !errorOccurred(err, w) {
        writeResponse(Response{
            Status: http.StatusOK,
            Result: res,
        }, w)
    }
}

// handlerCreate is a handler for MyApi.Create function.
func (srv *MyApi) handlerCreate(w http.ResponseWriter, r *http.Request) {
    if !isPost(w, r) {
        return
    }
    if !isAuthorized(w, r) {
        return
    }
    if err := r.ParseForm(); errorOccurred(err, w) {
        return
    }
    data := r.PostForm
    Login, err := getString("login", data, true)
    if errorOccurred(err, w) {
        return
    }
    if err = checkMinString("login", Login, 10); errorOccurred(err, w) {
        return
    }
    Name, err := getString("full_name", data, false)
    if errorOccurred(err, w) {
        return
    }
    Status, err := getString("status", data, false)
    if errorOccurred(err, w) {
        return
    }
    Status, err = inEnum("status", Status, []string{"user", "moderator", "admin"}, "user")
    if errorOccurred(err, w) {
        return
    }
    Age, err := getInt("age", data, false)
    if errorOccurred(err, w) {
        return
    }
    if err = checkMinInt("age", Age, 0); errorOccurred(err, w) {
        return
    }
    if err = checkMaxInt("age", Age, 128); errorOccurred(err, w) {
        return
    }
    res, err := srv.Create(r.Context(), CreateParams{
        Login: Login,
        Name: Name,
        Status: Status,
        Age: Age,
    })
    if !errorOccurred(err, w) {
        writeResponse(Response{
            Status: http.StatusOK,
            Result: res,
        }, w)
    }
}
// ServeHTTP defines functions serving different URLs for OtherApi method handlers.
func (srv *OtherApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch r.URL.Path {
    case "/user/create":
        srv.handlerCreate(w, r)
    default:
        writeResponse(Response{
            Status: http.StatusNotFound,
            Error:  "unknown method",
        }, w)
    }
}

// handlerCreate is a handler for OtherApi.Create function.
func (srv *OtherApi) handlerCreate(w http.ResponseWriter, r *http.Request) {
    if !isPost(w, r) {
        return
    }
    if !isAuthorized(w, r) {
        return
    }
    if err := r.ParseForm(); errorOccurred(err, w) {
        return
    }
    data := r.PostForm
    Username, err := getString("username", data, true)
    if errorOccurred(err, w) {
        return
    }
    if err = checkMinString("username", Username, 3); errorOccurred(err, w) {
        return
    }
    Name, err := getString("account_name", data, false)
    if errorOccurred(err, w) {
        return
    }
    Class, err := getString("class", data, false)
    if errorOccurred(err, w) {
        return
    }
    Class, err = inEnum("class", Class, []string{"warrior", "sorcerer", "rouge"}, "warrior")
    if errorOccurred(err, w) {
        return
    }
    Level, err := getInt("level", data, false)
    if errorOccurred(err, w) {
        return
    }
    if err = checkMinInt("level", Level, 1); errorOccurred(err, w) {
        return
    }
    if err = checkMaxInt("level", Level, 50); errorOccurred(err, w) {
        return
    }
    res, err := srv.Create(r.Context(), OtherCreateParams{
        Username: Username,
        Name: Name,
        Class: Class,
        Level: Level,
    })
    if !errorOccurred(err, w) {
        writeResponse(Response{
            Status: http.StatusOK,
            Result: res,
        }, w)
    }
}
