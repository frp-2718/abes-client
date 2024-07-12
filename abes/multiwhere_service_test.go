package abes

import (
	"encoding/xml"
	"fmt"
	"os"
	"testing"
)

func TestUnmarshalling(t *testing.T) {
	data, err := os.ReadFile("data.xml")
	if err != nil {
		t.Fatal("unable to read the data file")
	}

	var res serviceResult
	err = xml.Unmarshal(data, &res)
	if err != nil {
		fmt.Println(err)
		t.Fatal("unable to unmarshal data")
	}

	fmt.Println(res)

	for _, q := range res.Queries {
		fmt.Println(q.PPN)
		for _, l := range q.Result.Libraries {
			fmt.Println(l)
		}
		fmt.Println("=====================================")
	}

}
