package controllers

import (
	"testing"

	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/migrate/model"
)

func Test_rondom_create_user(t *testing.T) {
	model.Random_create_user()
}

func Test_rondom_create_tags(t *testing.T) {
	model.Random_create_tags()
}
