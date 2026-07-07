package api

import (
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	dateFormat = "20060102"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("Empty repetition rule field")
	}

	dateStart, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", errors.New("Invalid source time format" + err.Error())
	}

	splitRepeat := strings.Split(repeat, " ")
	letter := splitRepeat[0]

	switch letter {
	case "d":
		return daily(now, dateStart, repeat)
	case "y":
		return yearly(now, dateStart, repeat)
	case "w":
		return weekly(now, dateStart, repeat)
	case "m":
		return monthly(now, dateStart, repeat)
	default:
		return "", errors.New("Unsupported repeat rule format")
	}
}

func daily(now time.Time, dateStart time.Time, repeat string) (string, error) {
	splitRepeat := strings.Split(repeat, " ")
	if len(splitRepeat) != 2 {
		return "", errors.New("Invalid rule format d: expected 'd <number>'")
	}

	days, err := strconv.Atoi(splitRepeat[1])
	if err != nil {
		return "", errors.New("Invalid value in rule d: expected number")
	}

	if days < 1 {
		return "", errors.New("The number of days must be a positive")
	}

	if days > 400 {
		return "", errors.New("The number of days must be less than 400")
	}

	nextDate := dateStart

	for {
		nextDate = nextDate.AddDate(0, 0, days)
		if nextDate.After(now) {
			break
		}
	}

	return nextDate.Format(dateFormat), nil
}

func yearly(now time.Time, dateStart time.Time, repeat string) (string, error) {
	nextDate := dateStart

	for {
		nextDate = nextDate.AddDate(1, 0, 0)
		if nextDate.After(now) {
			break
		}
	}

	return nextDate.Format(dateFormat), nil
}

func weekly(now time.Time, dateStart time.Time, repeat string) (string, error) {
	splitRepeat := strings.Split(repeat, " ")
	if len(splitRepeat) != 2 {
		return "", errors.New("Invalid rule w: expected 'w <comma-separated numbers 1 through 7>'")
	}

	splitDays := strings.Split(splitRepeat[1], ",")
	weekdays := make([]time.Weekday, 0)

	for _, v := range splitDays {
		day, err := strconv.Atoi(v)
		if err != nil {
			return "", errors.New("Invalid value in rule w: expected number" + err.Error())
		}

		if day < 1 || day > 7 {
			return "", errors.New("The day of the week must be between 1 and 7")
		}

		if day == 7 {
			weekdays = append(weekdays, time.Sunday)
		} else {
			weekdays = append(weekdays, time.Weekday(day))
		}
	}

	nextDate := dateStart

	for {
		nextWeekday := nextDate.Weekday()
		for _, w := range weekdays {
			if nextWeekday == w && nextDate.After(now) {
				return nextDate.Format(dateFormat), nil
			}
		}
		nextDate = nextDate.AddDate(0, 0, 1)
	}
}

func monthly(now time.Time, dateStart time.Time, repeat string) (string, error) {
	splitRepeat := strings.Split(repeat, " ")
	if len(splitRepeat) != 2 && len(splitRepeat) != 3 {
		return "", errors.New("Invalid rule format m: expected 'm <comma-separated 1 to 31, -1, -2> [comma-separated 1 to 12]'")
	}

	splitDays := strings.Split(splitRepeat[1], ",")
	days := make([]int, 0)

	for _, sd := range splitDays {
		day, err := strconv.Atoi(sd)
		if err != nil {
			return "", errors.New("Invalid value in rule m: expected number" + err.Error())
		}

		if day >= 1 && day <= 31 || day == -1 || day == -2 {
			days = append(days, day)
		} else {
			return "", errors.New("Invalid day of month value: must be between 1 and 31, -1 or -2")
		}
	}

	months := make([]time.Month, 0)
	if len(splitRepeat) == 3 {
		splitMonths := strings.Split(splitRepeat[2], ",")

		for _, sm := range splitMonths {
			month, err := strconv.Atoi(sm)
			if err != nil {
				return "", errors.New("Invalid month format: expected number" + err.Error())
			}
			if month < 1 || month > 12 {
				return "", errors.New("The month must be between 1 and 12")
			}
			months = append(months, time.Month(month))
		}
	}

	sort.Ints(days)

	nextDate := dateStart

	for {
		year := nextDate.Year()
		month := nextDate.Month()

		if len(months) > 0 {
			needMonth := false
			for _, m := range months {
				if m == month {
					needMonth = true
					break
				}
			}
			if !needMonth {
				nextDate = nextDate.AddDate(0, 1, 0)
				continue
			}
		}

		lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()

		for _, d := range days {
			if d < 0 {
				continue
			}

			day := d

			if day > lastDay {
				continue
			}

			nextDay := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

			if nextDay.After(now) {
				return nextDay.Format(dateFormat), nil
			}

		}

		for _, d := range days {
			if d >= 0 {
				continue
			}

			day := lastDay + d + 1
			if day < 1 {
				continue
			}

			nextDay := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

			if nextDay.After(now) {
				return nextDay.Format(dateFormat), nil
			}
		}

		nextDate = nextDate.AddDate(0, 1, 0)
	}
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowReq := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	if date == "" {
		http.Error(w, "The 'date' field must be filled in", http.StatusBadRequest)
		return
	}

	if repeat == "" {
		http.Error(w, "The 'repeat' field must be filled in", http.StatusBadRequest)
		return
	}

	var now time.Time
	if nowReq == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(dateFormat, nowReq)
		if err != nil {
			http.Error(w, "Invalid 'now' format, expected YYYYMMDD", http.StatusBadRequest)
			return
		}
	}

	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}
