package contributors

import (
	"time"

	"github.com/open-sauced/pizza-cli/v2/api/services"
)

type DbContributor struct {
	AuthorLogin string    `json:"author_login"`
	UserID      int       `json:"user_id"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ContribResponse struct {
	Data []DbContributor   `json:"data"`
	Meta services.MetaData `json:"meta"`
}
