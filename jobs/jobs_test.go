package jobs

import (
	"testing"
)

func TestAddJobsToQueue(t *testing.T) {
	t.Run("when adding 3 jobs to queue they should be present", func(t *testing.T) {
		EnqueueJob("2020-08-03")
		EnqueueJob("2020-07-03")
		EnqueueJob("2020-06-03")

		job1 := DequeueJob()
		if job1.DateToParse != "2020-08-03" {
			t.Fatalf("The first dequeued job should parse the date 2020-08-03 but got %q", job1.DateToParse)

		}

		job2 := DequeueJob()
		if job2.DateToParse != "2020-07-03" {
			t.Fatalf("The first dequeued job should parse the date 2020-08-03 but got %q", job2.DateToParse)

		}

		job3 := DequeueJob()
		if job3.DateToParse != "2020-06-03" {
			t.Fatalf("The first dequeued job should parse the date 2020-08-03 but got %q", job3.DateToParse)

		}
	})
}
