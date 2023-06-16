package md

import (
	"testing"
	"time"
)

// create some tests to check parallel execution of tests
// go test -v -run Test_One
func Test_1(t *testing.T) {

	t.Run("1 - 1", func(t *testing.T) {
		t.Log("1 - 1 - start")
		time.Sleep(time.Second)
		t.Log("1 - 1 - end")
	})

	t.Run("1 - 2", func(t *testing.T) {
		t.Log("1 - 2 - start")
		time.Sleep(time.Second)
		t.Log("1 - 2 - end")
	})

	t.Run("1 -3", func(t *testing.T) {
		t.Log("1 - 3 - start")
		time.Sleep(time.Second)
		t.Log("1 - 3 - end")
	})

}

func Test_2(t *testing.T) {

	t.Run("2 - 1", func(t *testing.T) {
		t.Log("2 - 1 - start")
		time.Sleep(time.Second)
		t.Log("2 - 1 - end")
	})

	t.Run("2 - 2", func(t *testing.T) {
		t.Log("2 - 2 - start")
		time.Sleep(time.Second)
		t.Log("2 - 2 - end")
	})

	t.Run("2 -3", func(t *testing.T) {
		t.Log("2 - 3 - start")
		time.Sleep(time.Second)
		t.Log("2 - 3 - end")
	})

}

func Test_3(t *testing.T) {

	t.Run("3 - 1", func(t *testing.T) {
		t.Log("3 - 1 - start")
		time.Sleep(time.Second)
		t.Log("3 - 1 - end")
	})

	t.Run("3 - 2", func(t *testing.T) {
		t.Log("3 - 2 - start")
		time.Sleep(time.Second)
		t.Log("3 - 2 - end")
	})

	t.Run("3 -3", func(t *testing.T) {
		t.Log("3 - 3 - start")
		time.Sleep(time.Second)
		t.Log("3 - 3 - end")
	})

}
