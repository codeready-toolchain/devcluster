package auth

import (
	"fmt"
	"hash/fnv"
	"time"

	"github.com/codeready-toolchain/devcluster/pkg/log"

	uuid "github.com/satori/go.uuid"
)

// GenerateShortID generates a short ID.
// If the provided prefix is not an empty string then the generated ID will be in the following format: <prefix>-<date>-<id>.
// Where <date> is the current date in the following format: "Mmm-dd".
// For example for "rhd" prefix: "rhd-Jan02-5294981621"
func GenerateShortIDWithDate(prefix string) string {
	return shortID(prefix, true)
}

// GenerateShortID generates a short ID.
// If the provided prefix is not an empty string then the generated ID will be in the following format: <prefix>-<id>.
// For example for "rhd" prefix: "rhd-5294981621"
func GenerateShortID(prefix string) string {
	return shortID(prefix, false)
}

func shortID(prefix string, date bool) string {
	if prefix != "" {
		prefix = prefix + "-"
	}
	if date {
		prefix = fmt.Sprintf("%s%s-", prefix, time.Now().Format("Jan02"))
	}
	return fmt.Sprintf("%s%d", prefix, hash(uuid.NewV4().String()))
}

func hash(s string) uint32 {
	h := fnv.New32a()
	_, err := h.Write([]byte(s))
	if err != nil {
		log.Error(nil, err, "unable to generate cluster name")
	}
	return h.Sum32()
}
