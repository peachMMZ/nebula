package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type JSONTime time.Time

const TimeFormat = "2006-01-02 15:04:05"

// timeFormatWithFrac SQLite 等可能带小数秒
const timeFormatWithFrac = "2006-01-02 15:04:05.999999999"

func (t JSONTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format(TimeFormat))
}

func (t *JSONTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsed, err := time.Parse(TimeFormat, s)
	if err != nil {
		return err
	}
	*t = JSONTime(parsed)
	return nil
}

// Value 实现 driver.Valuer，供 GORM 写入数据库
func (t JSONTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}

// Scan 实现 sql.Scanner，供 GORM 从数据库读取
func (t *JSONTime) Scan(value interface{}) error {
	if value == nil {
		*t = JSONTime(time.Time{})
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*t = JSONTime(v)
		return nil
	case []byte:
		tm, err := parseTime(string(v))
		if err != nil {
			return err
		}
		*t = JSONTime(tm)
		return nil
	case string:
		tm, err := parseTime(v)
		if err != nil {
			return err
		}
		*t = JSONTime(tm)
		return nil
	default:
		*t = JSONTime(time.Time{})
		return nil
	}
}

func parseTime(s string) (time.Time, error) {
	if t, err := time.Parse(timeFormatWithFrac, s); err == nil {
		return t, nil
	}
	return time.Parse(TimeFormat, s)
}
