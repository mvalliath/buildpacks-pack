package pack_test

import (
	"github.com/buildpack/pack"
	"github.com/buildpack/pack/config"
	"github.com/buildpack/pack/mocks"
	"github.com/fatih/color"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"testing"

	h "github.com/buildpack/pack/testhelpers"
)

func TestInspectBuilder(t *testing.T) {
	color.NoColor = true
	spec.Run(t, "inspect-builder", testInspectBuilder, spec.Parallel(), spec.Report(report.Terminal{}))
}

func testInspectBuilder(t *testing.T, when spec.G, it spec.S) {
	var (
		inspector        *pack.BuilderInspector
		mockController   *gomock.Controller
		mockImageFactory *mocks.MockImageFactory
	)

	it.Before(func() {
		mockController = gomock.NewController(t)
		mockImageFactory = mocks.NewMockImageFactory(mockController)

		inspector = &pack.BuilderInspector{
			Config:       &config.Config{},
			ImageFactory: mockImageFactory,
		}
	})

	it.After(func() {
		mockController.Finish()
	})

	when("#Inspect", func() {
		when("builder has valid metatadata label", func() {
			it.Before(func() {
				mockBuilderImage := mocks.NewMockImage(mockController)
				mockImageFactory.EXPECT().NewRemote("some/builder").Return(mockBuilderImage, nil)
				mockBuilderImage.EXPECT().Label("io.buildpacks.pack.metadata").Return(`{"runImages": ["some/default", "gcr.io/some/default"]}`, nil)
			})

			when("builder exists in config", func() {
				it.Before(func() {
					inspector.Config.Builders = []config.Builder{
						{
							Image:     "some/builder",
							RunImages: []string{"some/run"},
						},
					}
				})

				it("returns the builder with the given name", func() {
					builder, err := inspector.Inspect("some/builder")
					h.AssertNil(t, err)
					h.AssertEq(t, builder.Image, "some/builder")
				})

				it("set the local run images", func() {
					builder, err := inspector.Inspect("some/builder")
					h.AssertNil(t, err)
					h.AssertEq(t, builder.LocalRunImages, []string{"some/run"})
				})
				it("set the defaults run images", func() {
					builder, err := inspector.Inspect("some/builder")
					h.AssertNil(t, err)
					h.AssertEq(t, builder.DefaultRunImages, []string{"some/default", "gcr.io/some/default"})
				})
			})

			when("builder does not exist in config", func() {
				it("returns the builder with default run images", func() {
					builder, err := inspector.Inspect("some/builder")
					h.AssertNil(t, err)
					h.AssertEq(t, builder.Image, "some/builder")
					h.AssertNil(t, builder.LocalRunImages)
					h.AssertEq(t, builder.DefaultRunImages, []string{"some/default", "gcr.io/some/default"})
				})
			})
		})

		when("builder has missing metadata label", func() {
			it.Before(func() {
				mockBuilderImage := mocks.NewMockImage(mockController)
				mockImageFactory.EXPECT().NewRemote("some/builder").Return(mockBuilderImage, nil)
				mockBuilderImage.EXPECT().Label("io.buildpacks.pack.metadata").Return("", errors.New("error!"))
			})

			it("returns an error", func() {
				_, err := inspector.Inspect("some/builder")
				h.AssertError(t, err, "failed to find run images for builder 'some/builder': error!")
			})
		})

		when("builder has invalid metadata label", func() {
			it.Before(func() {
				mockBuilderImage := mocks.NewMockImage(mockController)
				mockImageFactory.EXPECT().NewRemote("some/builder").Return(mockBuilderImage, nil)
				mockBuilderImage.EXPECT().Label("io.buildpacks.pack.metadata").Return("junk", nil)
			})

			it("returns an error", func() {
				_, err := inspector.Inspect("some/builder")
				h.AssertNotNil(t, err)
				h.AssertContains(t, err.Error(), "failed to parse run images for builder 'some/builder':")
			})
		})
	})
}
