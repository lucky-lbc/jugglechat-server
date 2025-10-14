package events

import "github.com/lucky-lbc/jugglechat-server/storages/models"

type UserRegisteEvent func(user models.User)

var userRegisteEvents []UserRegisteEvent

func init() {
	userRegisteEvents = []UserRegisteEvent{}
}

func RegisteUserRegisteEvent(event UserRegisteEvent) {
	userRegisteEvents = append(userRegisteEvents, event)
}

func TriggerUserRegiste(user models.User) {
	for _, event := range userRegisteEvents {
		event(user)
	}
}
