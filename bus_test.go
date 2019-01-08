package cromberbus

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type ACommandHandler struct{
	NumberOfHandleCalls int
}
func (h *ACommandHandler) Handle(command Command) error {
	h.NumberOfHandleCalls++
	return nil
}

func TestMapHandlerResolver_Resolve(t *testing.T) {
	isRequire := require.New(t)
	command := struct {}{}
	t.Run("Given a map handler resolver", func(t *testing.T) {
		sut := NewMapHandlerResolver()
		t.Run("When command handler is not found", func(t *testing.T) {
			handler, err := sut.Resolve(command)
			t.Run("Then an error is returned", func(t *testing.T) {
				isRequire.Nil(handler)
				isRequire.Error(err)
			})
		})
		t.Run("When a command with its handler is added", func(t *testing.T) {
			handler := &ACommandHandler{}
			sut.AddHandler(command, handler)
			t.Run("Then the handler is resolved", func(t *testing.T) {
				resolvedHandler, err := sut.Resolve(command)
				isRequire.Equal(handler, resolvedHandler)
				isRequire.NoError(err)
			})
		})
	})
}

func TestCromberBus_Dispatch(t *testing.T) {
	isRequire := require.New(t)
	t.Run("Given a cromberbus command bus", func(t *testing.T) {
		handler := ACommandHandler{}
		handlerResolver := HandlerResolverMock{
			ResolveFunc: func(command Command) (CommandHandler, error) {
				return &handler, nil
			},
		}
		sut := NewCromberBus(&handlerResolver)
		t.Run("When a command is dispatched", func(t *testing.T) {
			command := struct {}{}
			err := sut.Dispatch(command)
			t.Run("Then the resolved command handler handles the command", func(t *testing.T) {
				isRequire.Nil(err)
				resolverHasBeenCalled := len(handlerResolver.ResolveCalls()) > 0
				isRequire.True(resolverHasBeenCalled)
				handlerHasBeenCalled := handler.NumberOfHandleCalls > 0
				isRequire.True(handlerHasBeenCalled)
			})
		})
	})
}