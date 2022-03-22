package apiserver

import (
	"bytes"
	"fmt"
	"github.com/DeedsBaron/colors"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"runtime"
	"spam/internal/app/config"
	"testing"
)

type ResponseURL struct {
	shortURL string
}

func init() {
	var _ = func() bool {
		testing.Init()
		return true
	}()
	//change test dir for correct default config parse
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../../")
	err := os.Chdir(dir)
	fmt.Println(os.Getwd())
	if err != nil {
		panic(err)
	}
	logrus.SetLevel(logrus.FatalLevel)
	Cfg, err = config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	s = New(Cfg)
	s.logger.SetLevel(logrus.FatalLevel)
}

var (
	Cfg *config.Config
	s   *APIServer
)

func TestAddUser(t *testing.T) {
	fmt.Println(colors.Yellow + "Testing POST method \"AddUser\"" + colors.Res)
	jsons := []struct {
		user []byte
	}{
		{
			user: []byte(`{"Phone_num": "79164887712","Mobile_code": "7","Tag": "cool","Timezone_abbrev" : "MSK"}`),
		},
		{
			user: []byte(`asdasdasdasd`),
		},
		{
			user: []byte(`{"Mobile_code": "7","Tag": "cool","Timezone_abbrev" : "MSK"}`),
		},
		{
			user: []byte(`{"Phone_num": "79164887712","Tag": "cool","Timezone_abbrev" : "MSK"}`),
		},
		{
			user: []byte(`{"Phone_num": "79164887712","Mobile_code": "7","Timezone_abbrev" : "MSK"}`),
		},
		{
			user: []byte(`{"Phone_num": "79164887712","Mobile_code": "7","Tag": "cool"}`),
		},
		{
			user: []byte(`{"Phone_num": "79164887712213123","Mobile_code": "7","Tag": "cool", "Timezone_abbrev" : "MSK"}`),
		},
		{
			user: []byte(`{"Phone_num": "78923456711","Mobile_code": "7","Tag": "cool", "Timezone_abbrev" : "ASDASDAS"}`),
		},
	}
	tests := []struct {
		testDescription string
		req             *http.Request
		user            []byte
	}{
		{
			testDescription: colors.Purple + "TEST1: valid JSON" + colors.Res,
			req: func() *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/AddUser", bytes.NewBuffer(jsons[0].user))
				if err != nil {
					log.Fatal(err.Error())
				}
				return req
			}(),
		},
		{
			testDescription: colors.Purple + "TEST2:not valid JSON" + colors.Res,
			req: func() *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/AddUser", bytes.NewBuffer(jsons[1].user))
				if err != nil {
					log.Fatal(err.Error())
				}
				return req
			}(),
		},
		{
			testDescription: colors.Purple + "TEST3:missing field phone_num" + colors.Res,
			req: func() *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/AddUser", bytes.NewBuffer(jsons[2].user))
				if err != nil {
					log.Fatal(err.Error())
				}
				return req
			}(),
		},
		{
			testDescription: colors.Purple + "TEST4:missing field mobile_code" + colors.Res,
			req: func() *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/AddUser", bytes.NewBuffer(jsons[3].user))
				if err != nil {
					log.Fatal(err.Error())
				}
				return req
			}(),
		},
		{
			testDescription: colors.Purple + "TEST5:missing field tag" + colors.Res,
			req: func() *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/AddUser", bytes.NewBuffer(jsons[4].user))
				if err != nil {
					log.Fatal(err.Error())
				}
				return req
			}(),
		},
		{
			testDescription: colors.Purple + "TEST5:missing field Timezone_abbrev" + colors.Res,
			req: func() *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/AddUser", bytes.NewBuffer(jsons[5].user))
				if err != nil {
					log.Fatal(err.Error())
				}
				return req
			}(),
		},
		{
			testDescription: colors.Purple + "TEST6:not valid Phone_num" + colors.Res,
			req: func() *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/AddUser", bytes.NewBuffer(jsons[6].user))
				if err != nil {
					log.Fatal(err.Error())
				}
				return req
			}(),
		},
		{
			testDescription: colors.Purple + "TEST6:not valid Timezone_abbr" + colors.Res,
			req: func() *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/AddUser", bytes.NewBuffer(jsons[7].user))
				if err != nil {
					log.Fatal(err.Error())
				}
				return req
			}(),
		},
	}
	for i, val := range tests {
		res := httptest.NewRecorder()
		fmt.Println(val.testDescription)
		fmt.Println(string(jsons[i].user))

		s.configureRouter()
		s.router.ServeHTTP(res, val.req)

		fmt.Println(colors.Cyan+"Response body: "+colors.Res, res.Body.String())
		fmt.Println(colors.Cyan+"Response code: "+colors.Res, res.Code)
	}
}

func TestAddMailingList(t *testing.T) {
	fmt.Println(colors.Yellow + "Testing POST method \"AddMailingList\"" + colors.Res)
	jsons := []struct {
		user []byte
	}{
		{
			user: []byte(`{"Launch_date": "03 Mar 17 12:00 MSK","Message" : "Beeline","Filter" : "cool;7","Finish_date" : "22 Mar 17 12:00 MSK"}`),
		},
		{
			user: []byte(`{"asdasdasdasdasdasdasdasdasdasd" : }`),
		},
		{
			user: []byte(`{"Launch_date": "03 Mar 17 12:00 MSK","Message" : "Hi sweatie!","Filter" : "THECAKEISALIE","Finish_date" : "22 Mar 17 12:00 MSK"}`),
		},
	}
	tests := []struct {
		testDescription string
		req             *http.Request
		user            []byte
	}{
		{
			testDescription: colors.Purple + "TEST1: valid JSON" + colors.Res,
			req: func() *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/AddMailingList", bytes.NewBuffer(jsons[0].user))
				if err != nil {
					log.Fatal(err.Error())
				}
				return req
			}(),
		},
		{
			testDescription: colors.Purple + "TEST2:not valid JSON" + colors.Res,
			req: func() *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/AddMailingList", bytes.NewBuffer(jsons[1].user))
				if err != nil {
					log.Fatal(err.Error())
				}
				return req
			}(),
		},
		{
			testDescription: colors.Purple + "TEST3:bad filter field" + colors.Res,
			req: func() *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/AddMailingList", bytes.NewBuffer(jsons[2].user))
				if err != nil {
					log.Fatal(err.Error())
				}
				return req
			}(),
		},
	}
	for i, val := range tests {
		res := httptest.NewRecorder()
		fmt.Println(val.testDescription)
		fmt.Println(string(jsons[i].user))

		s.configureRouter()
		s.router.ServeHTTP(res, val.req)

		fmt.Println(colors.Cyan+"Response body: "+colors.Res, res.Body.String())
		fmt.Println(colors.Cyan+"Response code: "+colors.Res, res.Code)
	}
}

func TestUpdateMailingList(t *testing.T) {
	fmt.Println(colors.Yellow + "Testing POST method \"UpdateMailingList\"" + colors.Res)
	jsons := []struct {
		user []byte
	}{
		{
			user: []byte(`{"Launch_date": "2022-10-19","Message" : "Wtf!","Filter" : "sad;666","Finish_date" : "2022-05-20"}`),
		},
		{
			user: []byte(`{"Launch_date": "2022-10-19","Message" : "Wtf!","Filter" :"" ,"Finish_date" : "2022-05-20"}`),
		},
		{
			user: []byte(`{"Launch_date": "2022-03-19","Message" : "Hi sweatie!","Filter" : "THECAKEISALIE","Finish_date" : "2022-03-20"}`),
		},
	}
	tests := []struct {
		testDescription string
		req             *http.Request
		user            []byte
	}{
		{
			testDescription: colors.Purple + "TEST1: valid JSON and ID" + colors.Res,
			req: func() *http.Request {
				req, err := http.NewRequest(http.MethodPatch, "/UpdateMailingList/3", bytes.NewBuffer(jsons[0].user))
				if err != nil {
					log.Fatal(err.Error())
				}
				return req
			}(),
		},
		{
			testDescription: colors.Purple + "TEST2:not valid filter" + colors.Res,
			req: func() *http.Request {
				req, err := http.NewRequest(http.MethodPatch, "/UpdateMailingList/3", bytes.NewBuffer(jsons[1].user))
				if err != nil {
					log.Fatal(err.Error())
				}
				return req
			}(),
		},
		{
			testDescription: colors.Purple + "TEST3:bad filter field" + colors.Res,
			req: func() *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/AddMailingList", bytes.NewBuffer(jsons[2].user))
				if err != nil {
					log.Fatal(err.Error())
				}
				return req
			}(),
		},
	}
	for i, val := range tests {
		res := httptest.NewRecorder()
		fmt.Println(val.testDescription)
		fmt.Println(string(jsons[i].user))

		s.configureRouter()
		s.router.ServeHTTP(res, val.req)

		fmt.Println(colors.Cyan+"Response body: "+colors.Res, res.Body.String())
		fmt.Println(colors.Cyan+"Response code: "+colors.Res, res.Code)
	}
}

func TestDeleteUser(t *testing.T) {
	fmt.Println(colors.Yellow + "Testing DELETE method \"DeleteUser\"" + colors.Res)
	tests := []struct {
		testDescription string
		req             *http.Request
	}{
		{
			testDescription: colors.Purple + "TEST1: valid ID" + colors.Res + "\nDELETE /DeleteUser/1",
			req: func() *http.Request {
				req, err := http.NewRequest(http.MethodDelete, "/DeleteUser/18", nil)
				if err != nil {
					log.Fatal(err.Error())
				}
				return req
			}(),
		},
		{
			testDescription: colors.Purple + "TEST2: Not valid ID" + colors.Res + "\nDELETE /DeleteUser/1337",
			req: func() *http.Request {
				req, err := http.NewRequest(http.MethodDelete, "/DeleteUser/1488", nil)
				if err != nil {
					return req
				}
				return req
			}(),
		},
	}
	for _, val := range tests {
		res := httptest.NewRecorder()
		fmt.Println(val.testDescription)

		s.configureRouter()
		s.router.ServeHTTP(res, val.req)

		fmt.Println(colors.Cyan+"Response body: "+colors.Res, res.Body.String())
		fmt.Println(colors.Cyan+"Response code: "+colors.Res, res.Code)
	}
}
