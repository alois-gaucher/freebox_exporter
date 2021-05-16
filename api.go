package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	apiErrors = map[string]error{
		"invalid_token":           errors.New("the app token you are trying to use is invalid or has been revoked"),
		"insufficient_rights":     errors.New("your app permissions does not allow accessing this API"),
		"denied_from_external_ip": errors.New("you are trying to get an app_token from a remote IP"),
		"invalid_request":         errors.New("your request is invalid"),
		"ratelimited":             errors.New("too many auth error have been made from your IP"),
		"new_apps_denied":         errors.New("new application token request has been disabled"),
		"apps_denied":             errors.New("API access from apps has been disabled"),
		"internal_error":          errors.New("internal error"),
		"db_error":                errors.New("the database you are trying to access doesn't seem to exist"),
		"nodev":                   errors.New("invalid interface"),
	}
)

type ApiResponse interface {
	Status() (bool, string)
}

func (r *apiResponse) Status() (bool, string) {
	return r.Success, r.ErrorCode
}

func getApiData(authInf *authInfo, pr *postRequest, xSessionToken *string, response ApiResponse, requestBody io.Reader) error {
	freeboxToken, err := setFreeboxToken(authInf, xSessionToken)
	if err != nil {
		return err
	}

	client := http.Client{}
	req, err := http.NewRequest(pr.method, pr.url, requestBody)
	if err != nil {
		return err
	}
	req.Header.Add(pr.header, *xSessionToken)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == 404 {
		return errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, response)
	if err != nil {
		if debug {
			log.Println(string(body))
		}
		return err
	}

	_, errorCode := response.Status()
	if errorCode == "auth_required" {
		var err error
		*xSessionToken, err = getSessToken(freeboxToken, authInf, xSessionToken)
		if err != nil {
			return err
		}
	} else if errorCode != "" {
		if apiErrors[errorCode] == nil {
			return fmt.Errorf("%s: The API returns an unknown error_code: %s", pr.url, errorCode)
		}
		return apiErrors[errorCode]
	}

	return nil
}

func getRrdData(authInf *authInfo, pr *postRequest, xSessionToken *string, db string, fields []string) ([]int64, error) {
	d := &database{
		DB:        db,
		Fields:    fields,
		Precision: 10,
		DateStart: int(time.Now().Unix() - 10),
	}

	body, err := buildBody(d)
	if err != nil {
		return []int64{}, err
	}
	rrdTest := rrd{}
	err = getApiData(authInf, pr, xSessionToken, &rrdTest, body)
	if err != nil {
		return []int64{}, err
	}

	if len(rrdTest.Result.Data) == 0 {
		return []int64{}, nil
	}

	var result []int64
	for _, field := range fields {
		result = append(result, rrdTest.Result.Data[0][field])
	}
	return result, nil
}

func buildBody(d *database) (io.Reader, error) {
	if d == nil {
		return nil, nil
	}
	r, err := json.Marshal(*d)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(r), nil
}
