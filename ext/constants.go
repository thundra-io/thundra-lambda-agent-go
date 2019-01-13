package ext

const (
	// ThundraTagPrefix is prefix for the tags that are used in thundra internals
	ThundraTagPrefix = "thundra.span"
	// ClassNameKey defines class name 
	ClassNameKey = ThundraTagPrefix + ".className"
	DomainNameKey = ThundraTagPrefix + ".domainName"
	OperationTypeKey = "operation.type"
)