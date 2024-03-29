package experimental

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setup() {
	random = rand.New(rand.NewSource(42))
}

func TestConstantBackoff(t *testing.T) {
	backoff := ConstantBackoff(time.Millisecond * 500)

	for i := 0; i < 10; i++ {
		actual := backoff.Next()
		assert.Equal(t, time.Millisecond*500, actual)
	}
}

func TestExponentialBackoff(t *testing.T) {
	setup()

	initialDelay := time.Second * 1
	maxDelay := time.Second * 10
	backoff := ExponentialBackoff(initialDelay, maxDelay)

	assert.Equal(t, initialDelay, backoff.Next())

	expected := []time.Duration{
		time.Duration(2051072304),
		time.Duration(3731111519),
		time.Duration(7826726827),
		time.Duration(10000000000),
		time.Duration(10000000000),
	}
	for i := 0; i < 5; i++ {
		actual := backoff.Next()
		assert.Equal(t, expected[i], actual)
	}
}

func TestRandomBackoff(t *testing.T) {
	setup()

	min := time.Second * 1
	max := time.Second * 10
	backoff := RandomBackoff(min, max)

	for i := 0; i < 10; i++ {
		actual := backoff.Next()
		assert.GreaterOrEqual(t, actual, min)
		assert.LessOrEqual(t, actual, max)
	}
}
