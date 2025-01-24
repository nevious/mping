package objects

import (
	"testing"
	"time"
)

func TestCalculateAverage(t *testing.T) {
	record := DataRecord{
		durations: []time.Duration{
			time.Duration(1*time.Second),
			time.Duration(2*time.Second),
			time.Duration(3*time.Second),
			time.Duration(4*time.Second),
			time.Duration(5*time.Second),
			time.Duration(6*time.Second),
			time.Duration(7*time.Second),
			time.Duration(8*time.Second),
			time.Duration(9*time.Second),
		},
	}

	value := record.calc_average()
	if value.Seconds() != float64(5.0) {
		t.Errorf("Wrong average value %v", value.Seconds())
	}
}

func TestCalculateStandardDeviation(t *testing.T) {
	// This set should give us a nice sigma of 2
	// primitive test, yes. Be proud of it.
	record := DataRecord{
		durations: []time.Duration{
			time.Duration(1*time.Second),
			time.Duration(2*time.Second),
			time.Duration(3*time.Second),
			time.Duration(4*time.Second),
			time.Duration(5*time.Second),
			time.Duration(6*time.Second),
			time.Duration(7*time.Second),
		},
		Average: time.Duration(4*time.Second),
	}

	value, _ := record.calc_std()
	if value.Seconds() != float64(2) {
		t.Errorf("Wrong sigma value: %v", value.Seconds())
	}
}
