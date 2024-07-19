// modern re-implementation of rcp (remote copy protocol) daemon, originally part berkeley r-commands
// written by claude, 2024 | public domain
package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	TIMEOUT = 30 * time.Second
)

var (
	root string
	addr string
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.StringVar(&root, "root_dir", "", "Root directory for file operations")
	flag.StringVar(&addr, "bind_addr", ":514", "bind addr, :514 by default")
	flag.Parse()

	if root == "" {
		log.Fatal("Root directory must be specified using -root_dir flag")
	}

	var err error
	root, err = filepath.Abs(root)
	if err != nil {
		log.Fatalf("Failed to get absolute path for root directory: %v", err)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	log.Printf("RCP server listening on %s with root directory: %s", addr, root)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	remoteAddr := conn.RemoteAddr().String()
	log.Printf("New connection from %s", remoteAddr)

	reader := bufio.NewReader(conn)

	var fullCommand []byte
	commandComplete := false
	startTime := time.Now()

	for !commandComplete && time.Since(startTime) < TIMEOUT {
		conn.SetReadDeadline(time.Now().Add(time.Second))
		b, err := reader.ReadByte()
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			log.Printf("Failed to read byte from %s: %v", remoteAddr, err)
			return
		}
		fullCommand = append(fullCommand, b)

		if bytes.Contains(fullCommand, []byte("rcp ")) && bytes.HasSuffix(fullCommand, []byte{0}) {
			commandComplete = true
		}
	}

	if !commandComplete {
		log.Printf("Timeout waiting for full command from %s", remoteAddr)
		return
	}

	log.Printf("Received full command from %s: %s", remoteAddr, hex.Dump(fullCommand))

	parts := bytes.Split(fullCommand, []byte{0})
	log.Printf("Split command parts: %#v", parts)
	if len(parts) < 4 {
		log.Printf("Invalid command format from %s: not enough parts", remoteAddr)
		return
	}
	command := parts[3]
	log.Printf("Command: %q", command)

	if bytes.HasPrefix(command, []byte("rcp -t ")) {
		handleReceiveFile(conn, reader, remoteAddr, command)
	} else if bytes.HasPrefix(command, []byte("rcp -f ")) {
		handleSendFile(conn, reader, remoteAddr, command)
	} else {
		log.Printf("Invalid command format from %s: %q", remoteAddr, command)
	}
}

func handleReceiveFile(conn net.Conn, reader *bufio.Reader, remoteAddr string, command []byte) {
	targetPath := bytes.TrimPrefix(command, []byte("rcp -t "))
	log.Printf("Target path: %q", targetPath)

	// If targetPath is empty or ".", use the root directory
	if len(targetPath) == 0 || bytes.Equal(targetPath, []byte(".")) {
		targetPath = []byte("/")
	}

	fullTargetPath := filepath.Join(root, string(targetPath))
	fullTargetPath = filepath.Clean(fullTargetPath)

	if !bytes.HasPrefix([]byte(fullTargetPath), []byte(root)) {
		log.Printf("Attempted access outside of root directory: %s", fullTargetPath)
		return
	}

	log.Printf("Resolved target path: %q", fullTargetPath)

	conn.SetDeadline(time.Now().Add(TIMEOUT))
	conn.Write([]byte{0})
	time.Sleep(100 * time.Millisecond)
	conn.SetDeadline(time.Now().Add(TIMEOUT))
	conn.Write([]byte{0})
	log.Printf("Sent acknowledgements to %s", remoteAddr)

	for {
		conn.SetDeadline(time.Now().Add(TIMEOUT))
		command, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				log.Printf("Received EOF from %s, ending transfer", remoteAddr)
				return
			}
			log.Printf("Error reading command from %s: %v", remoteAddr, err)
			return
		}

		command = bytes.TrimSpace(command)
		log.Printf("Received command from %s: %q", remoteAddr, command)

		if bytes.HasPrefix(command, []byte("C")) {
			err = handleBinaryFileTransfer(conn, reader, remoteAddr, fullTargetPath, command)
			if err != nil {
				log.Printf("Error handling file transfer: %v", err)
				return
			}
		} else if bytes.HasPrefix(command, []byte("E")) {
			log.Printf("Received end of transfer from %s", remoteAddr)
			conn.SetDeadline(time.Now().Add(TIMEOUT))
			conn.Write([]byte{0}) // Send final acknowledgement
			return
		} else {
			log.Printf("Unexpected command received from %s: %q", remoteAddr, command)
		}
	}
}

func handleBinaryFileTransfer(conn net.Conn, reader *bufio.Reader, remoteAddr, targetPath string, fileInfo []byte) error {
	parts := bytes.SplitN(bytes.TrimSpace(fileInfo[1:]), []byte(" "), 3)
	if len(parts) != 3 {
		return fmt.Errorf("invalid file info format: %q", fileInfo)
	}

	fileMode, err := strconv.ParseInt(string(parts[0]), 8, 32)
	if err != nil {
		return fmt.Errorf("invalid file mode: %v", err)
	}

	fileSize, err := strconv.ParseInt(string(parts[1]), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid file size: %v", err)
	}

	fileName := string(parts[2])
	filePath := filepath.Join(targetPath, fileName)
	filePath = filepath.Clean(filePath)

	if !bytes.HasPrefix([]byte(filePath), []byte(root)) {
		return fmt.Errorf("attempted access outside of root directory: %s", filePath)
	}

	log.Printf("Preparing to receive file from %s: %q (mode: %o, size: %d)", remoteAddr, filePath, fileMode, fileSize)

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(fileMode))
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	conn.SetDeadline(time.Now().Add(TIMEOUT))
	if _, err := conn.Write([]byte{0}); err != nil {
		return fmt.Errorf("failed to send acknowledgement: %v", err)
	}

	buffer := make([]byte, 32*1024)
	var totalReceived int64
	for totalReceived < fileSize {
		conn.SetDeadline(time.Now().Add(TIMEOUT))
		n, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading file content: %v", err)
		}
		if n > 0 {
			if _, err := file.Write(buffer[:n]); err != nil {
				return fmt.Errorf("error writing to file: %v", err)
			}
			totalReceived += int64(n)
		}
	}

	if totalReceived != fileSize {
		log.Printf("Warning: received %d bytes, expected %d bytes", totalReceived, fileSize)
	}

	log.Printf("File transfer completed: %s (received %d bytes)", filePath, totalReceived)

	// Read the null byte that indicates end of file
	conn.SetDeadline(time.Now().Add(TIMEOUT))
	endByte, err := reader.ReadByte()
	if err != nil {
		if err != io.EOF {
			log.Printf("Warning: error reading end-of-file marker: %v", err)
		}
	} else if endByte != 0 {
		log.Printf("Warning: unexpected end-of-file marker: %d", endByte)
	}

	// Send final acknowledgement
	conn.SetDeadline(time.Now().Add(TIMEOUT))
	if _, err := conn.Write([]byte{0}); err != nil {
		return fmt.Errorf("failed to send final acknowledgement: %v", err)
	}

	return nil
}

func handleSendFile(conn net.Conn, reader *bufio.Reader, remoteAddr string, command []byte) {
	sourcePath := bytes.TrimPrefix(command, []byte("rcp -f "))
	log.Printf("Source path: %q", sourcePath)

	fullSourcePath := filepath.Join(root, string(sourcePath))
	fullSourcePath = filepath.Clean(fullSourcePath)

	if !bytes.HasPrefix([]byte(fullSourcePath), []byte(root)) {
		log.Printf("Attempted access outside of root directory: %s", fullSourcePath)
		return
	}

	log.Printf("Resolved source path: %q", fullSourcePath)

	file, err := os.Open(fullSourcePath)
	if err != nil {
		log.Printf("Failed to open file %s: %v", fullSourcePath, err)
		conn.SetDeadline(time.Now().Add(TIMEOUT))
		conn.Write([]byte{1})
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("Failed to get file info for %s: %v", fullSourcePath, err)
		conn.SetDeadline(time.Now().Add(TIMEOUT))
		conn.Write([]byte{1})
		return
	}

	conn.SetDeadline(time.Now().Add(TIMEOUT))
	if _, err := conn.Write([]byte{0}); err != nil {
		log.Printf("Failed to send initial byte to client %s: %v", remoteAddr, err)
		return
	}

	fileInfoStr := fmt.Sprintf("C%04o %d %s\n", fileInfo.Mode().Perm(), fileInfo.Size(), filepath.Base(fullSourcePath))
	conn.SetDeadline(time.Now().Add(TIMEOUT))
	if _, err := conn.Write([]byte(fileInfoStr)); err != nil {
		log.Printf("Failed to send file info for %s: %v", fullSourcePath, err)
		return
	}
	log.Printf("Sent file info to client %s: %s", remoteAddr, fileInfoStr)

	conn.SetDeadline(time.Now().Add(TIMEOUT))
	ack, err := reader.ReadByte()
	if err != nil || ack != 0 {
		log.Printf("Failed to receive acknowledgement for file info from client %s: %v", remoteAddr, err)
		return
	}

	buffer := make([]byte, 32*1024)
	var totalSent int64
	for totalSent < fileInfo.Size() {
		conn.SetDeadline(time.Now().Add(TIMEOUT))
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			log.Printf("Error reading file %s: %v", fullSourcePath, err)
			return
		}
		if n > 0 {
			conn.SetDeadline(time.Now().Add(TIMEOUT))
			if _, err := conn.Write(buffer[:n]); err != nil {
				log.Printf("Error sending file content for %s: %v", fullSourcePath, err)
				return
			}
			totalSent += int64(n)
		}
		if err == io.EOF {
			break
		}
	}

	log.Printf("File transfer completed: %s (sent %d bytes)", fullSourcePath, totalSent)

	conn.SetDeadline(time.Now().Add(TIMEOUT))
	if _, err := conn.Write([]byte{0}); err != nil {
		log.Printf("Failed to send end-of-file marker to client %s: %v", remoteAddr, err)
		return
	}

	conn.SetDeadline(time.Now().Add(TIMEOUT))
	ack, err = reader.ReadByte()
	if err != nil || ack != 0 {
		log.Printf("Failed to receive final acknowledgement from client %s: %v", remoteAddr, err)
		return
	}

	log.Printf("File %s sent successfully to client %s", fullSourcePath, remoteAddr)
}
