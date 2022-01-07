package services

import (
	"bootcamp/models"
	"testing"
)

var testCasesWrite = []struct {
	name     string
	fruit    *models.Fruit
	err      string
	response error
}{
	{
		"Should return a list of fruits",
		&models.Fruit{},
		"<nil>",
		nil,
	},
	// {
	// 	"Should return error",
	// 	&models.Fruit{},
	// 	"error",
	// 	nil,
	// },
}

func TestWriteMock(t *testing.T) {
	for _, tc := range testCasesWrite {
		t.Run(tc.name, func(t *testing.T) {
			// MOCK
			mock := mockFruitStore{}
			mock.On("CreateFruit").Return(tc.response)
			// SERVICE
			service := NewFruitService(mock)
			err := service.Write(tc.fruit)
			t.Log("ERROR:", err)
		})
	}
}
