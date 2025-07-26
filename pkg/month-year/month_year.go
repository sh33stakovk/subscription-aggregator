package monthyear

import (
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
