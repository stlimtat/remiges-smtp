package mail

import (
	"context"
	"fmt"
	"reflect"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
)

type DefaultMailProcessorFactory struct {
	cfgs       []config.MailProcessorConfig
	processors []IMailProcessor
	registry   map[string]reflect.Type
}

func NewDefaultMailProcessorFactory(
	_ context.Context,
) (IMailProcessorFactory, error) {
	result := &DefaultMailProcessorFactory{}
	result.registry = make(map[string]reflect.Type)
	result.registry[UnixDosProcessorType] = reflect.TypeOf(UnixDosProcessor{})
	return result, nil
}

func (f *DefaultMailProcessorFactory) NewMailProcessors(
	ctx context.Context,
	cfgs []config.MailProcessorConfig,
) ([]IMailProcessor, error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("Factory creating mail processors")

	if len(cfgs) < 1 {
		return nil, fmt.Errorf("no processors found")
	}
	f.cfgs = cfgs
	rawProcessors := make([]IMailProcessor, 0)
	for _, cfg := range cfgs {
		processor, err := f.NewMailProcessor(ctx, cfg)
		if err != nil {
			return nil, err
		}
		err = processor.Init(ctx, cfg)
		if err != nil {
			return nil, err
		}
		rawProcessors = append(rawProcessors, processor)
	}
	// sort the processors by index
	if len(rawProcessors) < 1 {
		return nil, fmt.Errorf("no processors found")
	}
	f.processors = make([]IMailProcessor, len(rawProcessors))
	for _, processor := range rawProcessors {
		// TODO: need to handle when the index is not in order
		f.processors[processor.Index()] = processor
	}

	return f.processors, nil
}

func (f *DefaultMailProcessorFactory) NewMailProcessor(
	ctx context.Context,
	cfg config.MailProcessorConfig,
) (IMailProcessor, error) {
	// create a single processor based on the processor config
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("Creating mail processor")
	var err error

	// use reflection to create the processor
	processorType, ok := f.registry[cfg.Type]
	if !ok {
		return nil, fmt.Errorf("processor type cannot be found")
	}
	processor := reflect.New(processorType).Interface().(IMailProcessor)
	// initialize the processor properly
	err = processor.Init(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return processor, nil
}

func (f *DefaultMailProcessorFactory) Process(
	ctx context.Context,
	inMail *Mail,
) (outMail *Mail, err error) {
	logger := zerolog.Ctx(ctx)
	// Builder function: process the mail through the processors
	for _, processor := range f.processors {
		logger.Info().Msgf("Processing mail through processor %d", processor.Index())
		inMail, err = processor.Process(ctx, inMail)
		if err != nil {
			return nil, err
		}
	}
	return inMail, nil
}
