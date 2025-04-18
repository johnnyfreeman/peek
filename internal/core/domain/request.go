package domain

type Request struct {
	Name string
	// Transport    TransportKind     // "http", "grpc", "ws", etc.
	// Definition   RequestDefinition // interface, implemented per transport
	Placeholders []Placeholder
}
