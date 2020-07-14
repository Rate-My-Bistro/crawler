package crawler

import (
	"os"
	"testing"
	"time"
)

func TestParseBistroWebsite(t *testing.T) {
	//setup
	bistroPageReader, _ := os.Open("./bistro.html")

	got := Start(bistroPageReader)

	t.Run("expect the correct size", func(t *testing.T) {
		if len(got) != 5 {
			t.Errorf("expected 5 days but got %q", len(got))
		}
	})

	t.Run("expect the correct date formant", func(t *testing.T) {
		if !isDate(getKeys(got)[0], t) {
			t.Errorf("expected the first date '2020-07-13' but got %q", getKeys(got)[0])
		}
	})

	t.Run("expect not nil", func(t *testing.T) {
		if got["2020-07-13"] == nil || len(got["2020-07-13"]) < 5 {
			t.Error("expected the first not nil but got nil")
		}
	})

	t.Run("expect the correct meal naming", func(t *testing.T) {
		if got["2020-07-13"][0].name != "Käsespätzle" {
			t.Errorf("expected the first meal of the week 'Käsespätzle' but got %q", got["2020-07-13"][0].name)
		}
	})
}

func isDate(dateString string, t *testing.T) bool {
	_, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		t.Error(err)
	}
	return err == nil
}

func getKeys(mealMap map[string][]Meal) (keys []string) {
	for k := range mealMap {
		keys = append(keys, k)
	}
	return keys
}
