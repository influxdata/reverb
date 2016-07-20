package popverb

import (
	"time"

	"github.com/labstack/echo"
	"github.com/markbates/pop"
	"github.com/markbates/reverb"
)

var Transaction = func(db *pop.Connection) echo.MiddlewareFunc {
	return func(handler echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			return db.Transaction(func(tx *pop.Connection) error {
				var lg *reverb.Logger
				clg := ctx.Get("lg")
				if clg != nil {
					lg = clg.(*reverb.Logger)
					pop.Log = func(s string) {
						if pop.Debug {
							lg.Printf("[tx:%s]: %s", tx.ID, s)
						}
					}
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
