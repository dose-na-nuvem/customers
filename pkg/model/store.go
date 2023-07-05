package model

import (
	"gorm.io/gorm"
)

type CustomerStore struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *CustomerStore {
	return &CustomerStore{
		db: db,
	}
}

func (s *CustomerStore) CreateCustomer(name string) (*Customer, error) {
	c := &Customer{
		Name: name,
	}
	s.db.Create(c)
	return c, nil
}

func (s *CustomerStore) ListCustomers() ([]Customer, error) {
	var customers []Customer
	result := s.db.Find(&customers)
	
	return customers, result.Error
}
