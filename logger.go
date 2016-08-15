package reverb

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/labstack/echo"
	elog "github.com/labstack/echo/log"
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
	elog.Logger
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

// NewLogger returns a `Logger` value and sets up default values such as log
// format, an "ID" for the request, etc...
func NewLogger(ctx echo.Context) *Logger {
	// id := getID(ctx)
	l := &Logger{
		Logger:    ctx.Logger(),
		Durations: durations{},
		Extras:    extras{},
	}
	return l
}
