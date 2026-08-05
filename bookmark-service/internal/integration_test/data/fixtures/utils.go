package fixtures

import (
	"time"

	"github.com/HemlockPham7/golang-system-design/internal/model"
)

var (
	TestTime = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
)

func GetTestBase(id string) model.Base {
	return model.Base{
		ID:        id,
		CreatedAt: TestTime,
		UpdatedAt: TestTime,
	}
}
