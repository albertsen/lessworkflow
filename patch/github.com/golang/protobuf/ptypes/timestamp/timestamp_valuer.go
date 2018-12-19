package timestamp

import (
	"database/sql/driver"
	"time"
)

func (ts *Timestamp) Value() (driver.Value, error) {
	return time.Unix(ts.Seconds, int64(ts.Nanos)).UTC(), nil
}
