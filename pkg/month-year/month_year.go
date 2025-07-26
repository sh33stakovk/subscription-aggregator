package monthyear

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type MonthYear struct {
	time.Time
}

func (my *MonthYear) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("01-2006", s)
	if err != nil {
		return err
	}
	my.Time = t
	return nil
}

func (my MonthYear) MarshalJSON() ([]byte, error) {
	return []byte(`"` + my.Format("01-2006") + `"`), nil
}

func (my MonthYear) Value() (driver.Value, error) {
	return my.Format("2006-01-02"), nil
}

func (my *MonthYear) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		my.Time = time.Date(v.Year(), v.Month(), 1, 0, 0, 0, 0, time.UTC)
		return nil
	case []byte:
		t, err := time.Parse("2006-01-02", string(v))
		if err != nil {
			return err
		}
		my.Time = t
		return nil
	case string:
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
		my.Time = t
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into MonthYear", value)
	}
}
