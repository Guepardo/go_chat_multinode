package exchanges

import "example.com/gochat/models"

type Exchange interface {
	GetPublishChannel() chan *models.Message
	GetConsumeChannel() chan *models.Message
	Start()
}
