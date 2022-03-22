package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

func DoWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	i := 1
	for j := attempts; j > 0; {
		logrus.Info("Trying to connect to database attempt ", i, "(", attempts, ")")
		i += 1
		if err = fn(); err != nil {
			time.Sleep(delay)
			j--
			continue
		}
		return nil
	}
	return
}

func HttpErrorWithoutBackSlashN(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	io.WriteString(w, error)
}

func FmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02dh:%02dm", h, m)
}
