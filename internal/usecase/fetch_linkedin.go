package usecase

import (
	"fmt"
	"profile/internal/model"
)

func FetchLinkedinData(username string) (*model.PlatformUserInfo, error) {
	//課金が必要なため、実装は省略
	_ = username
	return nil, fmt.Errorf("not implemented")
}
