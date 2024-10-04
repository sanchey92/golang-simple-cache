package server

import (
	"fmt"
	"github.com/sanchey92/golang-simple-cache/internal/config"
	"github.com/sanchey92/golang-simple-cache/internal/peer"
	"github.com/sanchey92/golang-simple-cache/internal/protocol"
	"github.com/sanchey92/golang-simple-cache/internal/storage"
	"github.com/sanchey92/golang-simple-cache/pkg/logger/sl"
	"github.com/tidwall/resp"
	"log/slog"
	"net"
)

type Server struct {
	log        *slog.Logger        // Logger instance for logging server events and errors.
	listenAddr string              // Address where the server listens for incoming connections.
	peers      map[*peer.Peer]bool // Map to keep track of active peers connected to the server.
	ln         net.Listener        // Network listener to accept new connections.
	addPeerCh  chan *peer.Peer     // Channel for adding new peers to the server.
	delPeerCh  chan *peer.Peer     // Channel for removing peers from the server.
	quitCh     chan struct{}       // Channel used to signal the server to shut down.
	msgCh      chan peer.Message   // Channel for handling incoming messages from peers.
	storage    *storage.Storage    // storage instance for handling data operations.
}

func New(cfg *config.Config, log *slog.Logger) *Server {
	return &Server{
		log:        log,
		listenAddr: cfg.ListenAddr,
		peers:      make(map[*peer.Peer]bool),
		addPeerCh:  make(chan *peer.Peer),
		delPeerCh:  make(chan *peer.Peer),
		quitCh:     make(chan struct{}),
		msgCh:      make(chan peer.Message),
		storage:    storage.New(),
	}
}

// Start initializes the server by setting up the listener and starting the main loop.
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		s.log.Error("Listener error", sl.Err(err))
	}
	s.ln = ln

	go s.loop()

	s.log.Info("server running", "listenAddr", s.listenAddr)

	return s.acceptLoop()
}

// loop is the main event loop that processes messages and manages peers.
func (s *Server) loop() {
	for {
		select {
		case msg := <-s.msgCh:
			if err := s.handleMessage(msg); err != nil {
				s.log.Error("raw message error", sl.Err(err))
			}
		case <-s.quitCh:
			return
		case p := <-s.addPeerCh:
			s.log.Info("peer connected", "remoteAddr", p.Conn.RemoteAddr())
			s.peers[p] = true
		case p := <-s.delPeerCh:
			s.log.Info("peer disconnected", "remoteAddr", p.Conn.RemoteAddr())
			delete(s.peers, p)
		}
	}
}

// acceptLoop continuously accepts new incoming connections and handles them.
func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			s.log.Error("accept error", sl.Err(err))
			continue
		}

		go s.handleConn(conn)
	}
}

// handleConn creates a new peer from the connection and starts reading messages.
func (s *Server) handleConn(conn net.Conn) {
	p := peer.New(conn, s.msgCh, s.delPeerCh, s.log)
	s.addPeerCh <- p
	if err := p.ReadLoop(); err != nil {
		s.log.Error("peer read error", sl.Err(err))
	}

	defer func() {
		s.log.Info("Closing connection", "remote_addr", p.Conn.RemoteAddr())
		if err := p.Conn.Close(); err != nil {
			s.log.Error("Failed to close connection", sl.Err(err))
		}
	}()
}

// handleMessage processes the received message based on its command type.
func (s *Server) handleMessage(msg peer.Message) error {
	switch v := msg.Cmd.(type) {
	case protocol.SetCommand:
		if err := s.storage.Set(v.Key, v.Value); err != nil {
			s.log.Error("Set command error", sl.Err(err))
			return err
		}
		if err := resp.NewWriter(msg.Peer.Conn).WriteString("OK"); err != nil {
			s.log.Error("Resp error", sl.Err(err))
			return err
		}
	case protocol.GetCommand:
		val, ok := s.storage.Get(v.Key)
		if !ok {
			return fmt.Errorf("key not found")
		}
		if err := resp.NewWriter(msg.Peer.Conn).WriteString(string(val)); err != nil {
			return err
		}
	}
	return nil
}
