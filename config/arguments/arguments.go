package arguments

import (
	"errors"
	"github.com/docopt/docopt.go"
	"pandora-pay/config"
	"pandora-pay/config/globals"
)

func InitArguments(argv []string) (err error) {

	if globals.Arguments, err = docopt.Parse(commands, argv, false, config.VERSION_STRING, false, false); err != nil {
		return errors.New("Error processing arguments" + err.Error())
	}

	return
}
