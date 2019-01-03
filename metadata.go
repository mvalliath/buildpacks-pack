package pack

const (
	MetadataLabel = "io.buildpacks.pack.metadata"
)

type BuilderImageMetadata struct {
	RunImages []string `json:"runImages"`
}

