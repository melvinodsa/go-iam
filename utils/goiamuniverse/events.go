package goiamuniverse

const (
	EventRoleUpdated = Event(Role) + ":" + Updated

	EventResourceCreated = Event(Resource) + ":" + Created

	EventUserCreated = Event(User) + ":" + Created

	EventClientCreated = Event(Client) + ":" + Created
	EventClientUpdated = Event(Client) + ":" + Updated
)

type Event string

const (
	Updated Event = "updated"
	Created Event = "created"
)
