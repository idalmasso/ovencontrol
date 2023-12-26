package hwinterface

import "fmt"

type ProcessingError struct {
	Operation string
}

func (e ProcessingError) Error() string {
	return fmt.Sprintf("Cannot execute the operation '%s' while processing", e.Operation)
}
