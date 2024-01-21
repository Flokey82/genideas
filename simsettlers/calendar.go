package simsettlers

import "log"

const daysInYear = 365

type TimeOfDay int

const (
	TimeOfDayMorning TimeOfDay = iota
	TimeOfDayAfternoon
	TimeOfDayEvening
	TimeOfDayNight
)

func (tod TimeOfDay) String() string {
	switch tod {
	case TimeOfDayMorning:
		return "morning"
	case TimeOfDayAfternoon:
		return "afternoon"
	case TimeOfDayEvening:
		return "evening"
	case TimeOfDayNight:
		return "night"
	default:
		return "unknown"
	}
}

type Calendar struct {
	TimeOfDay float64
	Day       uint16
	Year      int
}

func (c *Calendar) GetTimeOfDay() TimeOfDay {
	hour := int(c.TimeOfDay)
	if hour < 6 {
		return TimeOfDayNight
	} else if hour < 12 {
		return TimeOfDayMorning
	} else if hour < 18 {
		return TimeOfDayAfternoon
	} else {
		return TimeOfDayEvening
	}
}

func (c *Calendar) TickCalendar(elapsed float64) {
	prevTOD := c.GetTimeOfDay()

	// One second is one hour.
	c.TimeOfDay += elapsed * 3600.0
	if c.TimeOfDay > 24.0 {
		c.TimeOfDay -= 24.0
		c.Day++
		if c.Day > daysInYear {
			c.Day -= daysInYear
			c.Year++
		}
	}

	if newTOD := c.GetTimeOfDay(); newTOD != prevTOD {
		log.Printf("==========================Time of day changed from %s to %s!!!!!!!!!!!!!", prevTOD, newTOD)
	}
}
