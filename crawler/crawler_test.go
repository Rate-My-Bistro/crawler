package crawler

import (
	"os"
	"testing"
	"time"
)

func TestBistroWebCrawling(t *testing.T) {
	//setup
	bistroPageReader, _ := os.Open("./bistro.html")

	got := Crawl(bistroPageReader)

	t.Run("expect the correct size", func(t *testing.T) {
		if len(got) != 21 {
			t.Fatalf("expected 21 days but got %q", len(got))
		}
	})

	t.Run("expect the correct date formant", func(t *testing.T) {
		if !isDate(got[0].Date, t) {
			t.Fatalf("expected the first date '2020-07-13' but got %q", got[0].Date)
		}
	})

	t.Run("expect the correct count of mandatory supplements", func(t *testing.T) {
		if len(got[0].MandatorySupplements) != 1 {
			t.Fatalf("expected 1 days but got %q", len(got[0].MandatorySupplements))
		}
	})

	t.Run("expect the correct count of optional supplements", func(t *testing.T) {
		if len(got[1].OptionalSupplements) != 2 {
			t.Fatalf("expected 2 days but got %q", len(got[1].OptionalSupplements))
		}
	})

	t.Run("expect the correct meal naming", func(t *testing.T) {
		if got[0].Name != "Käsetortellini" {
			t.Fatalf("expected the first meal of the week 'Käsetortellini' but got %q", got[0].Name)
		}
	})

	t.Run("expect correct low kcal parsing", func(t *testing.T) {
		if got[0].LowKcal != false {
			t.Fatalf("expected the first meal of the week to be NOT low kcal")
		}
		if got[3].LowKcal != true {
			t.Fatalf("expected the fourth meal of the week to be low kcal")
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
