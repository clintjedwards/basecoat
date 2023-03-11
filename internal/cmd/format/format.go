package format

import (
	"fmt"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/clintjedwards/basecoat/proto"
	"github.com/fatih/color"

	"github.com/dustin/go-humanize"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// UnixMilli returns a humanized version of time given in unix millisecond. The zeroMsg is the string returned when
// the time is 0 and assumed to be not set.
func UnixMilli(unix int64, zeroMsg string, detail bool) string {
	if unix == 0 {
		return zeroMsg
	}

	if !detail {
		return humanize.Time(time.UnixMilli(unix))
	}

	relativeTime := humanize.Time(time.UnixMilli(unix))
	realTime := time.UnixMilli(unix).Format(time.RFC850)

	return fmt.Sprintf("%s (%s)", realTime, relativeTime)
}

// Takes a string enum and turns them into title case. If the value is unknown we turn it into
// a string of your choosing.
func NormalizeEnumValue[s ~string](value s, unknownString string) string {
	toTitle := cases.Title(language.AmericanEnglish)
	toLower := cases.Lower(language.AmericanEnglish)
	state := toTitle.String(toLower.String(string(value)))

	if strings.Contains(strings.ToLower(state), "unknown") {
		return unknownString
	}

	return state
}

func GenerateGenericTable(data [][]string, sep string, indent int) string {
	tableString := &strings.Builder{}
	table := tabwriter.NewWriter(tableString, 0, 2, 1, ' ', tabwriter.TabIndent)

	for _, item := range data {
		fmttedRow := ""

		for i := 1; i < indent; i++ {
			fmttedRow += " "
		}

		fmttedRow += strings.Join(item, fmt.Sprintf("\t%s ", sep))
		fmt.Fprintln(table, fmttedRow)
	}
	table.Flush()
	return tableString.String()
}

// Duration returns a humanized duration time for two epoch milli second times.
func Duration(start, end int64) string {
	if start == 0 {
		return "0s"
	}

	startTime := time.UnixMilli(start)
	endTime := time.Now()

	if end != 0 {
		endTime = time.UnixMilli(end)
	}

	duration := endTime.Sub(startTime)

	if duration > time.Second {
		truncate := time.Second
		return "~" + duration.Truncate(truncate).String()
	}

	return "~" + duration.String()
}

func ColorizeAccountState(state string) string {
	switch strings.ToUpper(state) {
	case proto.AccountState_ACTIVE.String():
		return color.GreenString(state)
	case proto.AccountState_DISABLED.String():
		return color.YellowString(state)
	case proto.AccountState_UNKNOWN.String():
		return color.RedString(state)
	default:
		return state
	}
}
