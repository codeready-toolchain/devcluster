package auth

import (
	"fmt"
	"hash/fnv"

	"github.com/codeready-toolchain/devcluster/pkg/log"
	uuid "github.com/satori/go.uuid"
)

// GenerateShortID generates a short ID.
// If the provided prefix is not an empty string then the generated ID will be in the following format: <prefix>-<id>.
// For example for "redhat" prefix: "redhat-5294981621"
func GenerateShortID(prefix string) string {
	if prefix != "" {
		prefix = prefix + "-"
	}
	return fmt.Sprintf("%s-%d", prefix, hash(uuid.NewV4().String()))
}

func hash(s string) uint32 {
	h := fnv.New32a()
	_, err := h.Write([]byte(s))
	if err != nil {
		log.Error(nil, err, "unable to generate cluster name")
	}
	return h.Sum32()
}
