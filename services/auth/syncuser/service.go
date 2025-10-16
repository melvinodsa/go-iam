package syncuser

import (
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils"
)

type Service interface {
	HandleEvent(event utils.Event[sdk.User])
}
