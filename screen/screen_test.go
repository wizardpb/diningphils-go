package screen

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type TestSuite struct {
	suite.Suite
	buffer *bytes.Buffer
}

func TestScreen(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupTest() {
	Initialize()
	s.buffer = &bytes.Buffer{}
	screen.stdOut = s.buffer
}

func (s *TestSuite) TearDownTest() {
	Close()
}

func (s *TestSuite) TestWrite() {
	testString := "this string"
	Write(testString)
	time.Sleep(time.Millisecond)
	s.Assert().Equal(clrScreen+testString, s.buffer.String())
	s.Assert().Equal(cursorPos{1, 1 + len(testString)}, screen.currentCursor)
}

func (s *TestSuite) TestWriteScreenLine() {
	testString := "this string"
	WriteScreenLine(3, 3, testString)
	time.Sleep(time.Millisecond)

	s.Assert().Equal(fmt.Sprintf(clrScreen+cursorPosition+clrLine, 3, 3)+testString+fmt.Sprintf(cursorPosition, 1, 1), s.buffer.String())
	s.Assert().Equal(cursorPos{1, 1}, screen.currentCursor)
}
