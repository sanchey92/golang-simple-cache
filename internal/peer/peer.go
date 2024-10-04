package peer

import (
	"errors"
	"fmt"
	"github.com/sanchey92/golang-simple-cache/internal/protocol"
	"github.com/sanchey92/golang-simple-cache/pkg/logger/sl"
	"github.com/tidwall/resp"
	"io"
	"log/slog"
	"net"
)

type Peer struct {
	Conn  net.Conn
	msgCh chan Message
	delCh chan *Peer
	log   *slog.Logger
	read  *resp.Reader
}

type Message struct {
	Cmd  protocol.Command
	Peer *Peer
}

func New(conn net.Conn, msgCh chan Message, delCh chan *Peer, log *slog.Logger) *Peer {
	return &Peer{
		Conn:  conn,
		msgCh: msgCh,
		delCh: delCh,
		log:   log,
		read:  resp.NewReader(conn),
	}
}

// Send transmits a message to the client.
//
// Parameters:
// - msg ([]byte): The message to send.
//
// Returns:
// - int, error : The number of bytes written and an error if any.
func (p *Peer) Send(msg []byte) (int, error) {
	return p.Conn.Write(msg)
}

// ReadLoop continuously reads and processes commands from the client's connection.
//
// Returns:
// - error: An error if the read loop fails.
func (p *Peer) ReadLoop() error {
	const op = "Peer.ReadLoop"
	p.log = p.log.With(slog.String("op", op))

	for {
		v, _, err := p.read.ReadValue()
		if err != nil {
			if errors.Is(err, io.EOF) {
				p.log.Info("connection closed by client")
				p.delCh <- p
				break
			}
			p.log.Error("Failed to read value %v", sl.Err(err))
		}

		cmd, err := parseCommand(v)
		if err != nil {
			p.log.Warn("Invalid command")
			continue
		}

		p.msgCh <- Message{
			Cmd:  cmd,
			Peer: p,
		}
	}
	return nil
}

// parseCommand converts a `resp.Value` into a corresponding protocol command.
//
// Parameters:
// - val (resp.Value): The RESP value received from the client.
//
// Returns:
// - protocol.Command: The parsed command if valid.
// - error: Error if the command is in an invalid format or not recognized.
func parseCommand(val resp.Value) (protocol.Command, error) {
	if val.Type() != resp.Array {
		return nil, fmt.Errorf("invalid command format")
	}

	arr := val.Array()
	if len(arr) < 2 {
		return nil, fmt.Errorf("invalid length of array")
	}

	rawCMD := arr[0].String()

	switch rawCMD {
	case protocol.CommandGET:
		if len(arr) != 2 {
			return nil, fmt.Errorf("invalid GET command format: expected 2 elements, got %d", len(arr))
		}
		return protocol.GetCommand{
			Key: arr[1].Bytes(),
		}, nil
	case protocol.CommandSET:
		if len(arr) != 3 {
			return nil, fmt.Errorf("invalid SET command format: expected 3 elements, got %d", len(arr))
		}
		return protocol.SetCommand{
			Key:   arr[1].Bytes(),
			Value: arr[2].Bytes(),
		}, nil
	default:
		return nil, fmt.Errorf("unhandled command")
	}
}
