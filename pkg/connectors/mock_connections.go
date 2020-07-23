package connectors

import (
	"fmt"

	"github.com/microlib/simple"
)

// Mock all connections
type MockConnectors struct {
	Logger *simple.Logger
}

func (c *MockConnectors) Error(msg string, val ...interface{}) {
	c.Logger.Error(fmt.Sprintf(msg, val...))
}

func (c *MockConnectors) Info(msg string, val ...interface{}) {
	c.Logger.Info(fmt.Sprintf(msg, val...))
}

func (c *MockConnectors) Debug(msg string, val ...interface{}) {
	c.Logger.Debug(fmt.Sprintf(msg, val...))
}

func (c *MockConnectors) Trace(msg string, val ...interface{}) {
	c.Logger.Trace(fmt.Sprintf(msg, val...))
}

func NewTestConnectors(logger *simple.Logger) Clients {
	conns := &MockConnectors{Logger: logger}
	return conns
}
