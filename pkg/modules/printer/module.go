package printer

import (
	"context"
	"fmt"

	"github.com/dipdup-io/workerpool"
	"github.com/rs/zerolog/log"
)

// predefined constants
const (
	ModuleName = "printer"
	InputName  = "Input"
)

// Module - the structure which is responsible for print received messages
type Module struct {
	Input chan string

	g workerpool.Group
}

// NewModule - constructor of printer structure
func NewModule() *Module {
	return &Module{
		Input: make(chan string, 16),
		g:     workerpool.NewGroup(),
	}
}

// Name -
func (printer *Module) Name() string {
	return ModuleName
}

// Close - gracefully stops module
func (printer *Module) Close() error {
	printer.g.Wait()

	close(printer.Input)
	return nil
}

// Start - starts module
func (printer *Module) Start(ctx context.Context) {
	printer.g.GoCtx(ctx, printer.listen)
}

func (printer *Module) listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-printer.Input:
			if !ok {
				return
			}
			log.Info().Str("obj_type", fmt.Sprintf("%T", msg)).Msgf("%##v", msg)
		}
	}
}
