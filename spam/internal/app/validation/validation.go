package validation

import (
	"bufio"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"regexp"
	"spam/pkg/utils"
	"strings"
	"time"
)

type User struct {
	PhoneNum       *string `json:"Phone_Num"`
	MobileCode     *string `json:"Mobile_Code"`
	Tag            *string `json:"Tag"`
	TimezoneAbbrev *string `json:"Timezone_abbrev"`
}

type MailingList struct {
	ID             int
	LaunchDate     *string `json:"Launch_date"`
	Message        *string `json:"Message"`
	Filter         *string `json:"Filter"`
	FinishDate     *string `json:"Finish_date"`
	LaunchDateTime *time.Time
	FinishDateTime *time.Time
}

type APIResponse struct {
	Code    *int    `json:"code"`
	Message *string `json:"message"`
}

func ValidateExternAPIResponse(resp *http.Response) bool {
	//inverse logic for loop
	if resp != nil {
		if resp.StatusCode == 200 {
			{
				t := APIResponse{}
				d := json.NewDecoder(resp.Body)
				err := d.Decode(&t)
				resp.Body.Close()
				if err != nil {
					//fmt.Println(err.Error())
					return true
				}
				if d.More() {
					//fmt.Println("More content")
					return true
				}
				return false
			}
		}
	}
	return true
}

func (t *User) validatePhoneNum(str string) bool {
	re := regexp.MustCompile(`^7[0-9]{10}$`)
	buf := string(re.Find([]byte(str)))
	if buf == "" {
		return false
	} else {
		return true
	}
}

func (t *MailingList) validateFilter(str string) bool {
	if !strings.Contains(str, ";") {
		return false
	}
	return true
}

func (t *MailingList) validateDate(str string) *time.Time {
	layout := "02 Jan 06 15:04 MST"
	tim, err := time.Parse(layout, str)
	if err != nil {
		return nil
	}
	return &tim
}

func (t *MailingList) ValidateMailingListJSON(w http.ResponseWriter, logger *logrus.Logger, d *json.Decoder) bool {
	if d.More() {
		utils.HttpErrorWithoutBackSlashN(w, "Extraneous data after JSON object", http.StatusBadRequest)
		logger.Error("Extraneous data after JSON object")
		return false
	}
	if t.Filter != nil && *t.Filter != "" {
		if t.validateFilter(*t.Filter) == false {
			utils.HttpErrorWithoutBackSlashN(w, "Filter field must be '<mobile_code>;<tag>'", http.StatusBadRequest)
			logger.Error("Filter field must be '<mobile_code>;<tag>'")
			return false
		}
	}
	if t.LaunchDate != nil {
		t.LaunchDateTime = t.validateDate(*t.LaunchDate)
		if t.LaunchDateTime == nil {
			utils.HttpErrorWithoutBackSlashN(w, "Wrong Finish_date format! Date must be in RFC822 format : '02 Jan 06 15:04 MST'", http.StatusBadRequest)
			logger.Error("Wrong Finish_date format! Date must be in RFC822 format : '02 Jan 06 15:04 MST'")
			return false
		}
	}
	if t.FinishDate != nil {
		t.FinishDateTime = t.validateDate(*t.FinishDate)
		if t.FinishDateTime == nil {
			utils.HttpErrorWithoutBackSlashN(w, "Wrong Finish_date format! Date must be in RFC822 format : '02 Jan 06 15:04 MST'", http.StatusBadRequest)
			logger.Error("Wrong Finish_date format! Date must be in RFC822 format : '02 Jan 06 15:04 MST'")
			return false
		}
	}
	return true
}

func (t *User) validateTimeZone(str string, logger *logrus.Logger) bool {
	path, err := os.Getwd()
	if err != nil {
		logger.Fatal(err.Error())
	}
	f, err := os.Open(path + "/internal/app/validation/TimeZonesAbbr")
	if err != nil {
		logger.Fatal(err.Error())
	}
	s := bufio.NewScanner(f)
	for s.Scan() {
		if str == s.Text() {
			return true
		}
	}
	err = s.Err()
	if err != nil {
		logger.Fatal(err.Error())
	}
	return false
}

func (t *User) ValidateUpdateUserJSON(w http.ResponseWriter, logger *logrus.Logger, d *json.Decoder) bool {
	if d.More() {
		utils.HttpErrorWithoutBackSlashN(w, "Extraneous data after JSON object", http.StatusBadRequest)
		logger.Error("Extraneous data after JSON object")
		return false
	}
	if t.PhoneNum != nil {
		if t.validatePhoneNum(*t.PhoneNum) == false {
			utils.HttpErrorWithoutBackSlashN(w, "Not valid Phone_num field. It must be 7XXXXXXXXXX!", http.StatusBadRequest)
			logger.Error("Not valid Phone_num field. It must be 7XXXXXXXXXX!")
			return false
		}
	}
	if t.TimezoneAbbrev != nil {
		if t.validateTimeZone(*t.TimezoneAbbrev, logger) == false {
			utils.HttpErrorWithoutBackSlashN(w, "Not valid Timezone_abbrev field. It must be from this list https://en.wikipedia.org/wiki/List_of_time_zone_abbreviations!", http.StatusBadRequest)
			logger.Error("Not valid Timezone_abbrev field. It must be from this list https://en.wikipedia.org/wiki/List_of_time_zone_abbreviations!")
			return false
		}
	}
	return true
}

func (t *User) ValidateAddUserJSON(w http.ResponseWriter, logger *logrus.Logger, d *json.Decoder) bool {
	if t.PhoneNum == nil {
		utils.HttpErrorWithoutBackSlashN(w, "Missing field 'PhoneNum' from JSON object", http.StatusBadRequest)
		logger.Error("Missing field 'PhoneNum' from JSON object")
		return false
	}
	if t.MobileCode == nil {
		utils.HttpErrorWithoutBackSlashN(w, "Missing field 'MobileCode' from JSON object", http.StatusBadRequest)
		logger.Error("Missing field 'MobileCode' from JSON object")
		return false
	}
	if t.Tag == nil {
		utils.HttpErrorWithoutBackSlashN(w, "Missing field 'Tag' from JSON object", http.StatusBadRequest)
		logger.Error("Missing field 'Tag' from JSON object")
		return false
	}
	if t.TimezoneAbbrev == nil {
		utils.HttpErrorWithoutBackSlashN(w, "Missing field 'TimezoneAbbrev' from JSON object", http.StatusBadRequest)
		logger.Error("Missing field 'TimezoneAbbrev' from JSON object")
		return false
	}
	if d.More() {
		utils.HttpErrorWithoutBackSlashN(w, "Extraneous data after JSON object", http.StatusBadRequest)
		logger.Error("Extraneous data after JSON object")
		return false
	}
	if t.validatePhoneNum(*t.PhoneNum) == false {
		utils.HttpErrorWithoutBackSlashN(w, "Not valid Phone_num field. It must be 7XXXXXXXXXX!", http.StatusBadRequest)
		logger.Error("Not valid Phone_num field. It must be 7XXXXXXXXXX!")
		return false
	}
	if t.validateTimeZone(*t.TimezoneAbbrev, logger) == false {
		utils.HttpErrorWithoutBackSlashN(w, "Not valid Timezone_abbrev field. It must be from this list https://en.wikipedia.org/wiki/List_of_time_zone_abbreviations!", http.StatusBadRequest)
		logger.Error("Not valid Timezone_abbrev field. It must be from this list https://en.wikipedia.org/wiki/List_of_time_zone_abbreviations!")
		return false
	}
	return true
}
