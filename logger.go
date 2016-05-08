package reverb

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/markbates/going/randx"
)

type extras []string

func (a extras) String() string {
	if len(a) == 0 {
		return ""
	}
	return fmt.Sprintf(" (%s)", strings.Join(a, " | "))
}

type durations map[string][]time.Duration

func (d durations) String() string {
	if len(d) == 0 {
		return ""
	}

	as := []string{}
	keys := []string{}
	for k := range d {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		ts := d[k]
		total := int64(0)
		for _, t := range ts {
			total += t.Nanoseconds()
		}
		as = append(as, fmt.Sprintf("%s: %s", k, time.Duration(total)))
	}

	return strings.Join(as, " | ")
}

// Logger gives you all the benefits of a `log.Logger`, but it also
// let's you log durations as well as extra logging data that will
// all get printed out nicely at the end of a request.
type Logger struct {
	*log.Logger
	Durations durations
	Extras    extras
}

// AddDurations let's you add `n` `time.Duration` values to a particular
// key. For example you could log multiple database requests to a single key,
// or log multiple page rendering times to a key. These durations will all be
// rolled up to one duration for that key when printed out in the log.
func (l *Logger) AddDurations(name string, ts ...time.Duration) {
	if l.Durations[name] == nil {
		l.Durations[name] = []time.Duration{}
	}
	l.Durations[name] = append(l.Durations[name], ts...)
}

// AddExtras is where you can add extra data you want to print out with the
// final request log message.
func (l *Logger) AddExtras(ex ...string) {
	l.Extras = append(l.Extras, ex...)
}

func (l *Logger) Error(err error) {
	if err != nil {
		l.Println(l.FmtError(err, 2))
	}
}

func (l *Logger) FmtError(err error, skip int) error {
	if err != nil {
		// notice that we're using 1, so it will actually log the where
		// the error happened, 0 = this function, we don't want that.
		pc, fn, line, _ := runtime.Caller(skip)

		err = fmt.Errorf("%s: %s[%s:%d] %v", err.Error(), runtime.FuncForPC(pc).Name(), fn, line, err)
	}
	return err
}

// NewLogger returns a `Logger` value and sets up default values such as log
// format, an "ID" for the request, etc...
func NewLogger(ctx *echo.Context) *Logger {
	id := getID(ctx)
	l := &Logger{
		Logger:    log.New(os.Stdout, fmt.Sprintf("[%s] ", id), log.LstdFlags),
		Durations: durations{},
		Extras:    extras{},
	}
	return l
}

func getID(ctx *echo.Context) string {
	c, err := ctx.Request().Cookie("_session_id")
	if err != nil {
		c = &http.Cookie{
			Name:    "_session_id",
			Value:   randx.String(10),
			Expires: time.Now().Add(10 * 365 * 24 * time.Hour), // 10 years
		}
		res := ctx.Response()
		res.Header().Add("Set-Cookie", c.String())
	}
	return c.Value
}
