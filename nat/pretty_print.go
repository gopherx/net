package nat

import (
	"fmt"
	"io"
	"math"
	"strconv"
)

type rowprinter interface {
	print(w io.Writer) (int, error)
}

type msgPrinter struct {
	msg []byte
	row int
}

var (
	hexLookup = createHexLookupTable()
)

func createHexLookupTable() []string {
	table := []string{}
	for i := 0; i < math.MaxUint8; i++ {
		txt := strconv.FormatUint(uint64(i), 16)
		if len(txt) == 1 {
			txt = "0" + txt
		}

		table = append(table, txt)

	}
	return table
}

func uint8ToHex(v uint8) string {
	return hexLookup[v]
}

func uint16ToHex(v uint16) (string, string) {
	return uint8ToHex(uint8(v >> 8)), uint8ToHex(uint8(v))
}

func (r *msgPrinter) print(w io.Writer) (int, error) {
	b0, b1, b2, b3 := "", "", "", ""
	group := " "
	desc := ""

	format := "0x%s 0x%s 0x%s 0x%s //%s   %s"

	switch {
	case r.row == 0:
		desc = "Request type and message length"
		break

	case r.row == 1:
		desc = "Magic Cookie"
		break

	case r.row == 2:
		group = "}"
		break

	case r.row == 3:
		desc = "TransactionID"
		group = "}"
		break

	case r.row == 4:
		group = "}"
		break

	default:
		return 0, nil
	}

	r.row++
	return fmt.Fprintf(w, format, b0, b1, b2, b3, group, desc)
}

// PrettyPrint prints all the messages side by side.
func PrettyPrint(w io.Writer, msg Message, opts *PrintOptions) error {
	pp := &PrettyPrinter{DefaultPrinter}
	return pp.Print(w, msg, opts)
}

type PrettyPrinter struct {
	Printer *MessagePrinter
}

func (p *PrettyPrinter) Print(w io.Writer, msg Message, opts *PrintOptions) error {
	bytes, err := DefaultPrinter.Print(msg, opts)
	if err != nil {
		return err
	}

	printer := &msgPrinter{bytes, 0}
	written := math.MaxInt32
	for written > 0 {
		n, err := printer.print(w)
		if err != nil {
			return err
		}

		_, err = fmt.Fprint(w, "\n")
		if err != nil {
			return err
		}
		written = n
	}

	return nil
}
