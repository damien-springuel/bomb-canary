package messagelogger

import "github.com/damien-springuel/bomb-canary/server/messagebus"

type printer interface {
	PrintCommand(m messagebus.Message)
	PrintEvent(m messagebus.Message)
}

type logger struct {
	printer printer
}

func New(printer printer) logger {
	return logger{
		printer: printer,
	}
}

func (l logger) Consume(m messagebus.Message) {
	switch m.Type() {
	case messagebus.CommandMessage:
		l.printer.PrintCommand(m)
	case messagebus.EventMessage:
		l.printer.PrintEvent(m)
	}
}
