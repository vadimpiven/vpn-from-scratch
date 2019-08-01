package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// SearchServer implements fake external service.
func SearchServer(w http.ResponseWriter, r *http.Request) {

}

// TestFindUsers runs the tests for FindUsers function of SearchClient struct.
func TestFindUsers(t *testing.T) {
	cases := []struct{
		Name        string
		AccessToken string
		Request     SearchRequest
		Response    *SearchResponse
		IsError     bool
	}{
		{
			Name: "limit must be > 0",
			Request: SearchRequest{
				Limit: -1,
			},
			Response: nil,
			IsError: true,
		},
		{
			Name: "offset must be > 0",
			Request: SearchRequest{
				Offset: -1,
			},
			Response: nil,
			IsError: true,
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			c := &SearchClient{
				URL: ts.URL,
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
