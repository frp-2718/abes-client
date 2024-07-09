package models

import "fmt"

type Response struct {
	Data []byte
}

func (r Response) String() string {
	return fmt.Sprintf("Data = %s", string(r.Data))
}
