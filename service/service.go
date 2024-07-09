package service

import "github.com/frp-2718/sudoc-client/models"

type Service interface {
	GetResponse(requestParams map[string]string) (*models.Response, error)
}
