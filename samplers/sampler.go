package samplers

// Sampler interface enables sampling of reported data
type Sampler interface {
	IsSampled(interface{}) bool
}
