package logger

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// LogFormatter - logrus formatter, implements logrus.Formatter
type LogFormatter struct {
	// FieldsOrder - default: fields sorted alphabetically
	FieldsOrder []string

	// TimestampFormat - default: models.DEFAULT_TIME_FORMAT"
	TimestampFormat string

	// Add caller fields
	WithCallerField bool
}

/*
 * Customized logrus Formatter:
 * Log format: <time> - LEVEL - Message=entry.message {Params: key=value ...}
 */
func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.StampMilli
	}

	// Output buffer
	b := &bytes.Buffer{}

	// Write TIME
	b.WriteString(fmt.Sprintf("<%s>", entry.Time.Format(timestampFormat)))

	// Write LEVEL
	b.WriteString(fmt.Sprintf(" - %s - ", strings.ToUpper(entry.Level.String())))

	// Write MESSAGE
	b.WriteString(fmt.Sprintf("Msg= %s", entry.Message))

	// Write PARAMS
	f.writeFields(b, entry)

	b.WriteByte('\n')

	return b.Bytes(), nil
}

func (f *LogFormatter) writeCaller(b *bytes.Buffer, entry *logrus.Entry) {
	if entry.HasCaller() {
		fmt.Fprintf(b, " caller=%s:%s:%d - ", entry.Caller.File, entry.Caller.Function, entry.Caller.Line)
	}
}

func (f *LogFormatter) writeFields(b *bytes.Buffer, entry *logrus.Entry) {
	if len(entry.Data) != 0 {
		b.WriteString("   {Params: ")
		if f.WithCallerField {
			f.writeCaller(b, entry)
		}
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}

		for _, field := range fields {
			f.writeField(b, entry, field)
			b.WriteString(" - ")
		}
		b.Truncate(len(b.Bytes()) - 3)
		b.WriteByte('}')
	}
}

func (f *LogFormatter) writeOrderedFields(b *bytes.Buffer, entry *logrus.Entry) {
	length := len(entry.Data)
	foundFieldsMap := map[string]bool{}
	for _, field := range f.FieldsOrder {
		if _, ok := entry.Data[field]; ok {
			foundFieldsMap[field] = true
			length--
			f.writeField(b, entry, field)
		}
	}

	if length > 0 {
		notFoundFields := make([]string, 0, length)
		for field := range entry.Data {
			if foundFieldsMap[field] == false {
				notFoundFields = append(notFoundFields, field)
			}
		}

		sort.Strings(notFoundFields)

		for _, field := range notFoundFields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *LogFormatter) writeField(b *bytes.Buffer, entry *logrus.Entry, field string) {
	fmt.Fprintf(b, "%s=%v", field, entry.Data[field])
}
