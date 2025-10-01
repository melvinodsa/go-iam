package goiamuniverse

const (
	EventRoleUpdated = Event(Role) + ":" + Updated

	EventResourceCreated = Event(Resource) + ":" + Created
	EventResourceDeleted = Event(Resource) + ":" + Deleted

	EventUserCreated = Event(User) + ":" + Created
	EventUserUpdated = Event(User) + ":" + Updated

	EventClientCreated = Event(Client) + ":" + Created
	EventClientUpdated = Event(Client) + ":" + Updated
)

type Event string

const (
	Updated Event = "updated"
	Created Event = "created"
	Deleted Event = "deleted"
)
