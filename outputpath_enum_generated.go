// Code generated by ggenums; DO NOT EDIT.
package sentinel

import (
	"encoding/json"
	"fmt"
)

type OutputPathEnum string

const (
	OutputPathStdout OutputPathEnum = "stdout"
	OutputPathStderr OutputPathEnum = "stderr"
)

var AllOutputPaths = []OutputPathEnum{
	OutputPathStdout,
	OutputPathStderr,
}

func (e OutputPathEnum) String() string {
	return string(e)
}

func (e OutputPathEnum) Validate() error {
	switch e {
	case OutputPathStdout, OutputPathStderr:
		return nil
	default:
		return fmt.Errorf("invalid OutputPath: %s", e)
	}
}

func ParseOutputPath(s string) (OutputPathEnum, error) {
	e := OutputPathEnum(s)
	if err := e.Validate(); err != nil {
		return "", err
	}
	return e, nil
}

func (e OutputPathEnum) MarshalJSON() ([]byte, error) {
	if err := e.Validate(); err != nil {
		return []byte("null"), nil
	}
	return json.Marshal(string(e))
}

func (e *OutputPathEnum) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	parsed, err := ParseOutputPath(s)
	if err != nil {
		return err
	}

	*e = parsed
	return nil
}