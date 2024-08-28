package atsemail_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/crewlinker/atsemail"
)

var JobApplicationNotificationExamples = []atsemail.JobApplicationNotification{
	{},
}

func TestRenderJobApplicationNotification(t *testing.T) {
	t.Parallel()

	for i, example := range JobApplicationNotificationExamples {
		t.Run(fmt.Sprintf("example %d", i), func(t *testing.T) {
			t.Parallel()

			render, err := atsemail.New(example)
			if err != nil {
				t.Fatalf("failed to init render: %v", err)
			}

			var txtbuf, htbuf bytes.Buffer
			if err := render.Render(&txtbuf, &htbuf); err != nil {
				t.Errorf("failed to render example %d: %v", i, err)
			}
		})
	}
}
