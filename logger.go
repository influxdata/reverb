package reverb

import (
	"fmt"
	"log"
	"net/http"
	"os"
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

type Logger struct {
	*log.Logger
	Durations durations
	Extras    extras
}

func (l *Logger) AddDurations(name string, ts ...time.Duration) {
	if l.Durations[name] == nil {
		l.Durations[name] = []time.Duration{}
	}
	l.Durations[name] = append(l.Durations[name], ts...)
}

func (l *Logger) AddExtras(ex ...string) {
	l.Extras = append(l.Extras, ex...)
}

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
