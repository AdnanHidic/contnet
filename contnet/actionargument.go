package contnet

type ActionArgumentType uint8

const (
	action_argument_type_unknown ActionArgumentType = iota
	action_argument_type_integer
	action_argument_type_string
	action_argument_type_float
)

var ActionArgumentTypes = struct {
	Unknown ActionArgumentType
	Integer ActionArgumentType
	String  ActionArgumentType
	Float   ActionArgumentType
}{
	Unknown: action_argument_type_unknown,
	Integer: action_argument_type_integer,
	String:  action_argument_type_string,
	Float:   action_argument_type_float,
}

type ActionArgument struct {
	Name  string
	Type  ActionArgumentType
	Value interface{}
}
type ActionArgumentFactory struct{}

func (factory ActionArgumentFactory) New(name string, argType ActionArgumentType, value interface{}) *ActionArgument {
	return &ActionArgument{
		Name:  name,
		Type:  argType,
		Value: value,
	}
}

func (arg *ActionArgument) Clone() *ActionArgument {
	return &ActionArgument{
		Name:  arg.Name,
		Type:  arg.Type,
		Value: arg.Value,
	}
}

type ActionArguments []*ActionArgument

func (args ActionArguments) Clone() ActionArguments {
	out := ActionArguments{}

	for i := 0; i < len(args); i++ {
		out = append(out, args[i].Clone())
	}

	return out
}
