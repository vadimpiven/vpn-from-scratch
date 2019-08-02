package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

// SearchServer implements fake external service.
func SearchServer(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	accessToken := r.Header.Get("AccessToken")
	switch accessToken {
	case "Timeout":
		time.Sleep(5 * time.Second)
		return
	case "InternalServerError":
		w.WriteHeader(http.StatusInternalServerError)
		return
	case "StatusBadRequest1":
		w.WriteHeader(http.StatusBadRequest)
		res, _ := json.Marshal(&SearchErrorResponse{"..."})
		w.Write(res)
		return
	case "StatusBadRequest2":
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{,,",`))
		return
	case "StatusBadRequest3":
		w.WriteHeader(http.StatusBadRequest)
		res, _ := json.Marshal(&SearchErrorResponse{"ErrorBadOrderField"})
		w.Write(res)
		return
	case "InvalidUser":
		w.Write([]byte(`{,,",`))
		return
	case "LimitReached":
		limit, _ := strconv.Atoi(r.FormValue("limit"))
		res, _ := json.Marshal(make([]User, limit))
		w.Write(res)
		return
	case "LimitUnreached":
		res, _ := json.Marshal(make([]User, 10))
		w.Write(res)
		return
	case "ValidToken":
		break
	default:
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	//body, err := ioutil.ReadAll(r.Body)
}

// TestFindUsers runs the tests for FindUsers function of SearchClient struct.
func TestFindUsers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	cases := []struct {
		Name        string
		AccessToken string
		URL         string
		Request     SearchRequest
		Response    *SearchResponse
		IsError     bool
	}{
		{
			Name: "limit must be > 0",
			URL:  ts.URL,
			Request: SearchRequest{
				Limit: -1,
			},
			Response: nil,
			IsError:  true,
		},
		{
			Name: "offset must be > 0",
			URL:  ts.URL,
			Request: SearchRequest{
				Offset: -1,
			},
			Response: nil,
			IsError:  true,
		},
		{
			Name:     "empty",
			Response: nil,
			IsError:  true,
		},
		{
			Name:        "unauthorized",
			AccessToken: "_____",
			URL:         ts.URL,
			Response:    nil,
			IsError:     true,
		},
		{
			Name:        "internal server error",
			AccessToken: "InternalServerError",
			URL:         ts.URL,
			Response:    nil,
			IsError:     true,
		},
		{
			Name:        "status bad request 1 - unknown",
			AccessToken: "StatusBadRequest1",
			URL:         ts.URL,
			Response:    nil,
			IsError:     true,
		},
		{
			Name:        "status bad request 2 - broken SearchErrorResponse",
			AccessToken: "StatusBadRequest2",
			URL:         ts.URL,
			Response:    nil,
			IsError:     true,
		},
		{
			Name:        "status bad request 3 - ErrorBadOrderField",
			AccessToken: "StatusBadRequest3",
			URL:         ts.URL,
			Response:    nil,
			IsError:     true,
		},
		{
			Name:        "invalid user",
			AccessToken: "InvalidUser",
			URL:         ts.URL,
			Response:    nil,
			IsError:     true,
		},
		{
			Name:        "limit reached",
			AccessToken: "LimitReached",
			URL:         ts.URL,
			Request: SearchRequest{
				Limit: 26,
			},
			Response: &SearchResponse{
				NextPage: true,
				Users:    make([]User, 25),
			},
			IsError: false,
		},
		{
			Name:        "limit unreached",
			AccessToken: "LimitUnreached",
			URL:         ts.URL,
			Request: SearchRequest{
				Limit: 15,
			},
			Response: &SearchResponse{
				NextPage: false,
				Users:    make([]User, 10),
			},
			IsError: false,
		},
		{
			Name:        "timeout",
			AccessToken: "Timeout",
			URL:         ts.URL,
			Response:    nil,
			IsError:     true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			c := &SearchClient{
				URL:         tc.URL,
				AccessToken: tc.AccessToken,
			}
			res, err := c.FindUsers(tc.Request)
			if err != nil && !tc.IsError {
				t.Errorf("Error occured but not expected\n")
			}
			if err == nil && tc.IsError {
				t.Errorf("Error expected but not occured\n")
			}
			tcJson, _ := json.Marshal(tc.Response)
			resJson, _ := json.Marshal(res)
			tcStr, resStr := string(tcJson), string(resJson)
			if tcStr != resStr {
				t.Errorf("Expected:\n%s\nGot:\n%s\n", tcStr, resStr)
			}
		})
	}
	ts.Close()
}
