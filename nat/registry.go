package nat

var (
	DefaultRegistry = AttributeRegistry{}
)

// AttributeRegistry is used to map parse and print functions to a type.
type AttributeRegistry map[AttributeType]Registration

type Registration struct {
	Name  string
	Parse AttributeParserFunc
	Print AttributePrinterFunc
}

func (a AttributeRegistry) Register(t AttributeType, name string, parse AttributeParserFunc, print AttributePrinterFunc) {
	a[t] = Registration{name, parse, print}
}
