package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"easycpp/backend/internal/dto"
)

const (
	defaultCPPRunTimeoutMs = 2000
	maxCPPRunTimeoutMs     = 10000
	maxTerminalOutputBytes = 64 * 1024
)

type CPPIdeService interface {
	RunCPP(ctx context.Context, req dto.RunCPPRequest) (dto.RunCPPResponse, error)
}

type cppIdeService struct{}

func NewCPPIdeService() CPPIdeService {
	return &cppIdeService{}
}

func (s *cppIdeService) RunCPP(ctx context.Context, req dto.RunCPPRequest) (dto.RunCPPResponse, error) {
	start := time.Now()
	resp := dto.RunCPPResponse{Language: "cpp"}

	timeLimit := clampTimeout(req.TimeLimitMs)

	workDir, err := os.MkdirTemp("", "easycpp-cpp-*")
	if err != nil {
		return resp, err
	}
	defer os.RemoveAll(workDir)

	sourcePath := filepath.Join(workDir, "main.cpp")
	binaryPath := filepath.Join(workDir, executableName("main"))

	if err := os.WriteFile(sourcePath, []byte(req.Code), 0o600); err != nil {
		return resp, err
	}

	compileStderr, compileExitCode, compileErr := runCommand(ctx, workDir, "g++", []string{"-std=c++17", "-O2", "-pipe", sourcePath, "-o", binaryPath}, "")
	resp.Compile = dto.CPPCompileResult{
		Succeeded: compileErr == nil,
		ExitCode:  compileExitCode,
		Stderr:    compileStderr,
	}

	if compileErr != nil {
		if errors.Is(compileErr, exec.ErrNotFound) {
			resp.Compile.Stderr = "compiler not found: g++ is not installed or not in PATH"
		}
		resp.DurationMs = time.Since(start).Milliseconds()
		return resp, nil
	}

	runCtx, cancel := context.WithTimeout(ctx, time.Duration(timeLimit)*time.Millisecond)
	defer cancel()

	stdout, stderr, runExitCode, timedOut, _ := runBinary(runCtx, binaryPath, req.Stdin)
	resp.Run = &dto.CPPRunResult{
		Succeeded: runExitCode == 0 && !timedOut,
		ExitCode:  runExitCode,
		Stdout:    stdout,
		Stderr:    stderr,
		TimedOut:  timedOut,
	}
	resp.DurationMs = time.Since(start).Milliseconds()

	return resp, nil
}

func clampTimeout(requested int) int {
	if requested <= 0 {
		return defaultCPPRunTimeoutMs
	}
	if requested > maxCPPRunTimeoutMs {
		return maxCPPRunTimeoutMs
	}
	return requested
}

func runCommand(ctx context.Context, workDir, name string, args []string, stdin string) (stderr string, exitCode int, runErr error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = workDir
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}

	var stderrBuffer bytes.Buffer
	cmd.Stderr = &stderrBuffer

	err := cmd.Run()
	return limitOutput(stderrBuffer.String()), extractExitCode(err), err
}

func runBinary(ctx context.Context, binaryPath, stdin string) (stdout string, stderr string, exitCode int, timedOut bool, runErr error) {
	cmd := exec.CommandContext(ctx, binaryPath)
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}

	var stdoutBuffer bytes.Buffer
	var stderrBuffer bytes.Buffer
	cmd.Stdout = &stdoutBuffer
	cmd.Stderr = &stderrBuffer

	err := cmd.Run()
	timeoutHit := errors.Is(ctx.Err(), context.DeadlineExceeded)

	if timeoutHit {
		stderrBuffer.WriteString("\nExecution timed out.")
	}

	return limitOutput(stdoutBuffer.String()), limitOutput(stderrBuffer.String()), extractExitCode(err), timeoutHit, err
}

func extractExitCode(err error) int {
	if err == nil {
		return 0
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
		return 1
	}

	return 1
}

func limitOutput(text string) string {
	if len(text) <= maxTerminalOutputBytes {
		return text
	}
	truncated := text[:maxTerminalOutputBytes]
	return fmt.Sprintf("%s\n\n[output truncated]", truncated)
}

func executableName(base string) string {
	if runtime.GOOS == "windows" {
		return base + ".exe"
	}
	return base
}
