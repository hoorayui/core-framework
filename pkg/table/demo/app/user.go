package main

import (
	"context"

	cap "framework/pkg/table/proto"
)

type TestUserInfoProvider struct{}

func (*TestUserInfoProvider) GetCurrentUser(ctx context.Context) (*cap.UserInfo, error) {
	return &cap.UserInfo{
		Id:          "1",
		UserName:    "demo",
		DisplayName: "演示",
	}, nil
}

func (*TestUserInfoProvider) GetUserInfoByID(id string) (*cap.UserInfo, error) {
	return &cap.UserInfo{
		Id:          "1",
		UserName:    "demo",
		DisplayName: "演示",
	}, nil
}
