package crawler

import (
	"github.com/ansgarS/rate-my-bistro-crawler"
	"os"
	"testing"
	"time"
)

func TestParseBistroWebsite(t *testing.T) {
	//setup
	bistroPageReader, _ := os.Open("./bistro.html")

	got := Start(bistroPageReader)

	t.Run("expect the correct size", func(t *testing.T) {
		if len(got) != 25 {
			t.Fatalf("expected 5 days but got %q", len(got))
		}
	})

	t.Run("expect the correct date formant", func(t *testing.T) {
		if !isDate(got[0].date, t) {
			t.Fatalf("expected the first date '2020-07-13' but got %q", got[0].date)
		}
	})

	t.Run("expect the correct meal naming", func(t *testing.T) {
		if got[0].Name != "Käsespätzle" {
			t.Fatalf("expected the first meal of the week 'Käsespätzle' but got %q", got[0].Name)
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

func getKeys(mealMap map[string][]main.Meal) (keys []string) {
	for k := range mealMap {
		keys = append(keys, k)
	}
	return keys
}
