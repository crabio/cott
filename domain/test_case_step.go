package domain

import "bytes"

type TestCaseStep struct {
	Name     string
	StepFunc func() error
}

func (s *TestCaseStep) String() string {
	var buf bytes.Buffer

	buf.WriteByte('{')
	buf.WriteString("name: ")
	buf.WriteString(s.Name)
	buf.WriteByte('}')

	return buf.String()
}
