package intmail

import (
	"context"
	"fmt"
	"reflect"
	"sort"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/crypto"
	"github.com/stlimtat/remiges-smtp/pkg/pmail"
)

type DefaultMailProcessorFactory struct {
	Cfgs          []config.MailProcessorConfig
	CryptoFactory *crypto.CryptoFactory
	Processors    []IMailProcessor
	Registry      map[string]reflect.Type
}

func NewDefaultMailProcessorFactory(
	_ context.Context,
	cfgs []config.MailProcessorConfig,
	cryptoFactory *crypto.CryptoFactory,
) (*DefaultMailProcessorFactory, error) {
	result := &DefaultMailProcessorFactory{
		Cfgs:          cfgs,
		CryptoFactory: cryptoFactory,
	}
	result.Registry = make(map[string]reflect.Type)
	result.Registry[BodyHeadersProcessorType] = reflect.TypeOf(BodyHeadersProcessor{})
	result.Registry[BodyProcessorType] = reflect.TypeOf(BodyProcessor{})
	result.Registry[DKIMProcessorType] = reflect.TypeOf(DKIMProcessor{})
	result.Registry[MergeBodyProcessorType] = reflect.TypeOf(MergeBodyProcessor{})
	result.Registry[UnixDosProcessorType] = reflect.TypeOf(UnixDosProcessor{})
	return result, nil
}

func (f *DefaultMailProcessorFactory) Init(
	ctx context.Context,
	_ config.MailProcessorConfig,
) error {
	var err error
	// Ignore the config, we will use the cfgs from the NewFactory
	// This is to map to the IMailProcessor interface
	f.Processors, err = f.NewMailProcessors(ctx, f.Cfgs)
	if err != nil {
		return err
	}
	return nil
}

func (f *DefaultMailProcessorFactory) NewMailProcessors(
	ctx context.Context,
	cfgs []config.MailProcessorConfig,
) ([]IMailProcessor, error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().
		Interface("cfgs", cfgs).
		Msg("DefaultMailProcessorFactory")

	if len(cfgs) < 1 {
		return nil, fmt.Errorf("no processors found")
	}
	f.Cfgs = cfgs

	rawProcessors := make(map[int]IMailProcessor)
	processorIndices := make([]int, 0)
	for _, cfg := range cfgs {
		processor, err := f.NewMailProcessor(ctx, cfg)
		if err != nil {
			return nil, err
		}
		rawProcessors[processor.Index()] = processor
		processorIndices = append(processorIndices, processor.Index())
	}
	// sort the processors by index
	if len(rawProcessors) < 1 {
		return nil, fmt.Errorf("no processors found")
	}
	result := make([]IMailProcessor, 0)
	// get the list of indices, then sort them
	sort.Ints(processorIndices)
	for _, processorIdx := range processorIndices {
		result = append(result, rawProcessors[processorIdx])
	}
	f.Processors = result

	return result, nil
}

func (f *DefaultMailProcessorFactory) NewMailProcessor(
	ctx context.Context,
	cfg config.MailProcessorConfig,
) (IMailProcessor, error) {
	// create a single processor based on the processor config
	logger := zerolog.Ctx(ctx).With().
		Str("type", cfg.Type).
		Int("index", cfg.Index).
		Interface("args", cfg.Args).
		Logger()
	logger.Debug().Msg("Creating mail processor")
	var err error

	// use reflection to create the processor
	processorType, ok := f.Registry[cfg.Type]
	if !ok {
		return nil, fmt.Errorf("processor type cannot be found")
	}
	result := reflect.New(processorType).Interface().(IMailProcessor)
	// initialize the processor properly
	err = result.Init(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// Specific processor initialization for dkim
	if cfg.Type == DKIMProcessorType {
		dkimProcessor, ok := result.(*DKIMProcessor)
		if !ok {
			return nil, fmt.Errorf("processor is not a DKIMProcessor")
		}
		err = dkimProcessor.InitDKIMCrypto(ctx, f.CryptoFactory)
		if err != nil {
			return nil, err
		}
		result = dkimProcessor
	}

	return result, nil
}

func (_ *DefaultMailProcessorFactory) Index() int {
	return -1
}

func (f *DefaultMailProcessorFactory) Process(
	ctx context.Context,
	inMail *pmail.Mail,
) (*pmail.Mail, error) {
	logger := zerolog.Ctx(ctx)
	var err error
	// Builder function: process the mail through the processors
	for _, processor := range f.Processors {
		logger.Debug().
			Int("idx", processor.Index()).
			Str("processor", reflect.TypeOf(processor).String()).
			Msg("Running processor")
		inMail, err = processor.Process(ctx, inMail)
		if err != nil {
			return nil, err
		}
	}
	return inMail, nil
}
