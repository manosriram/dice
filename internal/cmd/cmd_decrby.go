package cmd

import (
	"strconv"

	"github.com/dicedb/dice/internal/object"
	dstore "github.com/dicedb/dice/internal/store"
	"github.com/dicedb/dice/wire"
)

var cDECRBY = &DiceDBCommand{
	Name:      "DECRBY",
	HelpShort: "evalDECRBY decrements the value of the specified key in args by the specified decrement",
	Eval:      evalDECRBY,
}

// evalDECRBY decrements the value of the specified key in args by the specified decrement
// The function expects exactly two arguments:
//  - The Key to decrement
//  - The Decrement

// If the key does not exist, new key is created with value 0,
// The value of the new key is then decremented by specified decrement.
// If the key exists but does not contain an integer, an error is returned.
//
// Parameters:
//   - c *Cmd: The command context containing the arguments
//   - s *dstore.Store: The data store instance
//
// Returns:
//   - *CmdRes: Response containing the new integer value after decrement
//   - error: Error if wrong number of arguments or wrong value type
func evalDECRBY(c *Cmd, s *dstore.Store) (*CmdRes, error) {
	if len(c.C.Args) != 2 {
		return cmdResNil, errWrongArgumentCount("DECRBY")
	}

	decrAmount, err := strconv.ParseInt(c.C.Args[1], 10, 64)
	if err != nil {
		return cmdResNil, errIntegerOutOfRange
	}

	delta := int64(-decrAmount)

	key := c.C.Args[0]
	obj := s.Get(key)
	if obj == nil {
		obj = s.NewObj(delta, INFINITE_EXPIRATION, object.ObjTypeInt)
		s.Put(key, obj)
		return &CmdRes{R: &wire.Response{
			Value: &wire.Response_VInt{VInt: delta},
		}}, nil
	}

	switch obj.Type {
	case object.ObjTypeInt:
		break
	default:
		return cmdResNil, errWrongTypeOperation("DECRBY")
	}

	val, _ := obj.Value.(int64)
	val += delta

	obj.Value = val
	return &CmdRes{R: &wire.Response{
		Value: &wire.Response_VInt{VInt: val},
	}}, nil
}
