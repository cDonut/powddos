package pow

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Challenge struct for POW.
type Challenge struct {
	Level     int64
	Timestamp int64
	Data      string
}

// Creates new Challenge.
func NewChallenge(level int64, data string) *Challenge {
	return &Challenge{Level: level, Data: data, Timestamp: time.Now().Unix()}
}

// Casts Challenge to string.
func (c *Challenge) String() string {
	return fmt.Sprintf("%d:%d:%s", c.Level, c.Timestamp, c.Data)
}

// Solves POW Challenge.
func (c *Challenge) Solve(attempts uint64) (string, error) {
	var counter uint64
	s := c.String()

	for {
		solution := fmt.Sprintf("%s:%d", s, counter)
		if checkHash(sha1.Sum([]byte(solution)), c.Level) == false {
			counter++

			if counter > attempts {
				return "", errors.New("Attempts limit reached")
			}

			continue
		}

		return solution, nil
	}
}

// Parses Challenge from string. Challenge format - level:timestamp:data
func ParseChallenge(s string) (*Challenge, error) {
	parts := strings.Split(s, ":")
	if len(parts) < 3 {
		return nil, errors.New("Bad challenge format")
	}

	l, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, err
	}

	t, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, err
	}

	return &Challenge{Level: l, Timestamp: t, Data: parts[2]}, nil
}

// Checks POW solution
func CheckSolution(solution string, level int64, data string, expire int64) bool {
	c, err := ParseChallenge(solution)
	if err != nil {
		return false
	}

	if c.Level < level || c.Data != data {
		return false
	}

	t := time.Now()

	if c.Timestamp > t.Unix() || (t.Unix()-c.Timestamp > expire) {
		return false
	}

	return checkHash(sha1.Sum([]byte(solution)), level)
}

func checkHash(hash [20]byte, level int64) bool {
	if level > 20 {
		level = 20
	}

	var i int64
	for i = 0; i < level; i++ {
		if hash[i] != 0 {
			return false
		}
	}

	return true
}
