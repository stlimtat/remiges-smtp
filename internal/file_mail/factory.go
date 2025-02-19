package file_mail

import (
	"context"
	"fmt"
	reflect "reflect"
	"sort"

	"github.com/rs/zerolog"
	"github.com/stlimtat/remiges-smtp/internal/config"
	"github.com/stlimtat/remiges-smtp/internal/file"
	"github.com/stlimtat/remiges-smtp/internal/mail"
)

type MailTransformerFactory struct {
	Cfgs         []config.FileMailConfig
	registry     map[string]reflect.Type
	transformers []IMailTransformer
}

func NewMailTransformerFactory(
	_ context.Context,
	cfgs []config.FileMailConfig,
) *MailTransformerFactory {
	result := &MailTransformerFactory{
		Cfgs: cfgs,
	}
	result.registry = make(map[string]reflect.Type)
	result.registry[BodyTransformerType] = reflect.TypeOf(BodyTransformer{})
	result.registry[HeadersTransformerType] = reflect.TypeOf(HeadersTransformer{})
	result.registry[HeaderContentTypeTransformerType] = reflect.TypeOf(HeaderContentTypeTransformer{})
	result.registry[HeaderFromTransformerType] = reflect.TypeOf(HeaderFromTransformer{})
	result.registry[HeaderMsgIDTransformerType] = reflect.TypeOf(HeaderMsgIDTransformer{})
	result.registry[HeaderSubjectTransformerType] = reflect.TypeOf(HeaderSubjectTransformer{})
	result.registry[HeaderToTransformerType] = reflect.TypeOf(HeaderToTransformer{})
	return result
}

func (f *MailTransformerFactory) Init(
	ctx context.Context,
	_ config.FileMailConfig,
) error {
	var err error
	f.transformers, err = f.NewMailTransformers(ctx, f.Cfgs)
	if err != nil {
		return err
	}
	return nil
}

func (f *MailTransformerFactory) NewMailTransformers(
	ctx context.Context,
	cfgs []config.FileMailConfig,
) ([]IMailTransformer, error) {
	logger := zerolog.Ctx(ctx)
	logger.Info().
		Interface("cfgs", cfgs).
		Msg("MailTransformerFactory")

	if len(cfgs) == 0 {
		return nil, fmt.Errorf("no cfgs")
	}
	f.Cfgs = cfgs

	rawTransformers := make(map[int]IMailTransformer)
	transformerIndices := make([]int, 0)
	for _, cfg := range f.Cfgs {
		transformer, err := f.NewMailTransformer(ctx, cfg)
		if err != nil {
			return nil, err
		}
		err = transformer.Init(ctx, cfg)
		if err != nil {
			return nil, err
		}
		rawTransformers[transformer.Index()] = transformer
		transformerIndices = append(transformerIndices, transformer.Index())
	}

	// sort the transformers by index
	if len(rawTransformers) < 1 {
		return nil, fmt.Errorf("no transformers found")
	}
	sort.Ints(transformerIndices)
	result := make([]IMailTransformer, 0)
	for _, transformerIdx := range transformerIndices {
		// TODO: need to handle when the index is not in order
		result = append(result, rawTransformers[transformerIdx])
	}
	return result, nil
}

func (f *MailTransformerFactory) NewMailTransformer(
	ctx context.Context,
	cfg config.FileMailConfig,
) (IMailTransformer, error) {
	// create a single transformer based on the transformer config
	logger := zerolog.Ctx(ctx).With().
		Str("type", cfg.Type).
		Int("index", cfg.Index).
		Interface("args", cfg.Args).
		Logger()
	logger.Debug().Msg("Creating mail transformer")
	var err error

	// use reflection to create the transformer
	transformerType, ok := f.registry[cfg.Type]
	if !ok {
		return nil, fmt.Errorf("transformer type cannot be found")
	}
	result := reflect.New(transformerType).Interface().(IMailTransformer)
	// initialize the transformer properly
	err = result.Init(ctx, cfg)
	if err != nil {
		logger.Error().Err(err).Msg("transformer.Init")
		return nil, err
	}
	return result, nil
}

func (_ *MailTransformerFactory) Index() int {
	return -1
}

func (f *MailTransformerFactory) Transform(
	ctx context.Context,
	fileInfo *file.FileInfo,
	inMail *mail.Mail,
) (outMail *mail.Mail, err error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("Transforming mail")

	for _, transformer := range f.transformers {
		inMail, err = transformer.Transform(ctx, fileInfo, inMail)
		if err != nil {
			return nil, err
		}
	}
	return inMail, nil
}
