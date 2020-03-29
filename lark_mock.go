package main

import "github.com/stretchr/testify/mock"

//MockLark implements a mock version of Lark
type MockLark struct {
	mock.Mock
}

//GetAccessToken provides a mocked version of `GetAccessToken`
func (m *MockLark) GetAccessToken(appDetails AppDetails) (LarkToken, error) {
	rets := m.Called(appDetails)
	return rets.Get(0).(LarkToken), rets.Error(1)
}
//GetBotGroups provides a mocked version of `GetBotGroups`
func (m *MockLark) GetBotGroups(token string) (LarkData, error) {
	rets := m.Called(token)
	return rets.Get(0).(LarkData), rets.Error(1)
}

//SendMessage provides a mocked version of `SendMessage`
func (m *MockLark) SendMessage(token string, larkMessageRequest LarkMessageRequest) int {
	rets := m.Called(token, larkMessageRequest)
	return rets.Get(0).(int)
}

//InitMockLark initializes the lark variable, but it assigns
//a new MockLark instance to it, instead of an actual lark
func InitMockLark() *MockLark {
	l := new(MockLark)
	lark = l
	return l
}
