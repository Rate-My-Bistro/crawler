package jobs

import (
	"testing"
)

func TestAddJobsToQueue(t *testing.T) {
	t.Run("when adding 3 jobs to queue they should be present", func(t *testing.T) {
		EnqueueJob("2020-08-03")
		EnqueueJob("2020-07-03")
		EnqueueJob("2020-06-03")

		if len(JobQueue) != 3 {
			t.Fatalf("The queue size should be 3 but is %q", len(JobQueue))
		}

		job1 := DequeueJob()
		if job1.DateToParse != "2020-08-03" {
			t.Fatalf("The first dequeued job should parse the date 2020-08-03 but got %q", job1.DateToParse)

		}
		if len(JobQueue) != 2 {
			t.Fatalf("The queue size should be 2 but is %q", len(JobQueue))

		}

		job2 := DequeueJob()
		if job2.DateToParse != "2020-07-03" {
			t.Fatalf("The first dequeued job should parse the date 2020-08-03 but got %q", job2.DateToParse)

		}
		if len(JobQueue) != 1 {
			t.Fatalf("The queue size should be 1 but is %q", len(JobQueue))

		}

		job3 := DequeueJob()
		if job3.DateToParse != "2020-06-03" {
			t.Fatalf("The first dequeued job should parse the date 2020-08-03 but got %q", job3.DateToParse)

		}
		if len(JobQueue) != 0 {
			t.Fatalf("The queue size should be 0 but is %q", len(JobQueue))

		}
	})
}
