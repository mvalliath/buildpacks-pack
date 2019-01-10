package pack

import (
	"encoding/json"
	"github.com/buildpack/lifecycle/image"
	"github.com/buildpack/pack/config"
	"github.com/buildpack/pack/style"
	"github.com/pkg/errors"
)

type BuilderInspector struct {
	Config       *config.Config
	ImageFactory ImageFactory
}

type Builder struct {
	Image            string
	LocalRunImages   []string
	DefaultRunImages []string
}

func DefaultBuilderInspector() (*BuilderInspector, error) {
	cfg, err := config.NewDefault()
	if err != nil {
		return nil, err
	}

	factory, err := image.DefaultFactory()
	if err != nil {
		return nil, err
	}

	return &BuilderInspector{
		Config:       cfg,
		ImageFactory: factory,
	}, nil
}

func (b *BuilderInspector) Inspect(builderName string) (Builder, error) {
	var err error
	var localRunImages, defaultRunImages []string
	if builderConfig := b.Config.GetBuilder(builderName); builderConfig != nil {
		localRunImages = builderConfig.RunImages
	}
	defaultRunImages, err = b.getDefaultRunImages(builderName)
	if err != nil {
		return Builder{}, err
	}

	builder := Builder{
		Image:            builderName,
		LocalRunImages:   localRunImages,
		DefaultRunImages: defaultRunImages,
	}

	return builder, nil
}

func (b *BuilderInspector) getDefaultRunImages(builderName string) ([]string, error) {
	builderImage, err := b.ImageFactory.NewRemote(builderName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get remote image %s", style.Symbol(builderName))
	}
	var metadata BuilderImageMetadata
	label, err := builderImage.Label(MetadataLabel)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find run images for builder %s", style.Symbol(builderName))
	}
	if err := json.Unmarshal([]byte(label), &metadata); err != nil {
		return nil, errors.Wrapf(err, "failed to parse run images for builder %s", style.Symbol(builderName))
	}
	return metadata.RunImages, nil
}
