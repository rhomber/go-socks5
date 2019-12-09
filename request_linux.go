package socks5

import (
	"github.com/hanwen/go-fuse/v2/splice"
	"io"
	"net"
)

// proxy is used to suffle data from src to destination, and sends errors
// down a dedicated channel
func proxy(dst io.Writer, src io.Reader, errCh chan error) {
	pair, err := splice.Get()
	pair.MaxGrow()

	TcpDst := dst.(*net.TCPConn)
	TcpSrc := src.(*net.TCPConn)

	TcpSrc.SetReadBuffer(256 * 1024)
	TcpSrc.SetWriteBuffer(256 * 1024)

	TcpDst.SetReadBuffer(256 * 1024)
	TcpDst.SetWriteBuffer(256 * 1024)

	FdDst, _ := TcpDst.File()
	FdSrc, _ := TcpSrc.File()

	if err == nil {
		for {
			w, err := splice.SpliceCopy(FdDst, FdSrc, pair)
			if err != nil || w == 0 {
				break
			}
		}
		pair.Close()
		errCh <- err
	}
}
