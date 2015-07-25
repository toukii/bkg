package main

import (
	"testing"
)

func TestTargetPath(t *testing.T) {
	paths := targetPath(".")
	t.Log(paths)
}

// 3257793874
