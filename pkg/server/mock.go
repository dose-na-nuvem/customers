package server

import "github.com/dose-na-nuvem/customers/pkg/model"

type mockStore struct {
	createCustomerFunc func(name string) (*model.Customer, error)
}

func (m *mockStore) CreateCustomer(name string) (*model.Customer, error) {
	if m.createCustomerFunc != nil {
		return m.createCustomerFunc(name)
	}

	return nil, nil
}
