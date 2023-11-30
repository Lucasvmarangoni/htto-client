package entities

import (
	"errors"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

type File struct {
	ID        string    `json:"document_id" valid:"uuid"`
	Type      string    `json:"type" valid:"notnull"`
	CreatedAt time.Time `json:"created_at" valid:"-"`
	Customer  string    `json:"customer" valid:"notnull"`	
}

func (f *File) validate() error {

	typ := strings.ToLower(f.Type)

	if typ != "contract" && typ != "extract" && typ != "invoice" {
		return errors.New("invalid type")
	}

	_, err := govalidator.ValidateStruct(f)

	if err != nil {
		return err
	}
	return nil
}

func NewFile(typ string, createdAt time.Time, customer string) *File {
	return &File{
		ID:        uuid.NewV4().String(),
		Type:      typ,
		CreatedAt: createdAt,
		Customer:  customer,		
	}
}