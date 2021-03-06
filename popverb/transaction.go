package popverb

import (
	"time"

	"github.com/influxdata/reverb"
	"github.com/labstack/echo"
	"github.com/markbates/pop"
)

var Transaction = func(db *pop.Connection) echo.MiddlewareFunc {
	return func(handler echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			return db.Transaction(func(tx *pop.Connection) error {
				var lg *reverb.Logger
				clg := ctx.Get("lg")
				if clg != nil {
					lg = clg.(*reverb.Logger)
				}
				ctx.Set("tx", tx)

				before := tx.Elapsed
				err := handler(ctx)
				after := tx.Elapsed
				if clg != nil {
					logPopTimings(lg, []time.Duration{time.Duration(after - before)})
				}
				return err
			})
		}
	}
}

func logPopTimings(lg *reverb.Logger, ts []time.Duration) {
	if len(ts) > 0 {
		lg.AddDurations("POP", ts...)
	}
}
