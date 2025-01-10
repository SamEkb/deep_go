package _interface

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v interface_test.go

type UserService struct {
	// not need to implement
	NotEmptyStruct bool
}
type MessageService struct {
	// not need to implement
	NotEmptyStruct bool
}

type Container struct {
	constructors map[string]interface{}
}

func NewContainer() *Container {
	return &Container{
		constructors: map[string]interface{}{},
	}
}

func (c *Container) RegisterType(name string, constructor interface{}) {
	c.constructors[name] = constructor
}

func (c *Container) Resolve(name string) (interface{}, error) {
	value, ok := c.constructors[name]
	if !ok {
		return nil, fmt.Errorf("constructor with name %s not found", name)
	}

	if constructor, ok := value.(func() interface{}); ok {
		return constructor(), nil
	}
	return nil, fmt.Errorf("%s is not a constructor func", name)
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() interface{} {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() interface{} {
		return &MessageService{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)
}
