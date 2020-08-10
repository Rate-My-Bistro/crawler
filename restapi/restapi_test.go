package restapi

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"testing"
)

func TestInsertOrUpdate(t *testing.T) {
	//start server
	Serve()

	//setup client
	client := resty.New()
	resp, _ := client.R().
		EnableTrace().
		Get("http://localhost/jobs")
	fmt.Println("Body       :\n", resp)

	t.Run("insert a record and update it", func(t *testing.T) {

	})
}
