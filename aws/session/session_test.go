package session

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCreateAWSSession(t *testing.T) {
	assert := assert.New(t)

	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)

	// Set AWS_REGION to a made-up value to ensure no overlap with real regions
	os.Setenv("AWS_REGION", "us-test-1")

	session := newSession(os.Getenv("AWS_PROFILE"), "", logger)
	assert.Equalf(*session.Session.Config.Region, "us-test-1", "AWS SDK did not load correct region from envvar, expected 'us-test-1', got %s", *session.Session.Config.Region)

	session = newSession(os.Getenv("AWS_PROFILE"), "us-test-2", logger)
	assert.Equalf(*session.Session.Config.Region, "us-test-2", "AWS SDK did not load correct region from envvar, expected 'us-test-2', got %s", *session.Session.Config.Region)
}

func TestCreateAWSSessionPool(t *testing.T) {
	assert := assert.New(t)

	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)

	t.Run("test session permutations", func(t *testing.T) {
		// Ensure that AWS sessions are created for each permutation of region + profile
		// Duplicate values will create duplicate (but not unique) sessions
		regions := []string{"us-test-1", "us-test-2", "us-test-3"}
		profiles := []string{"profile1", "profile2", "profile3"}
		testSessionPool := NewPool(profiles, regions, logger)

		assert.Lenf(testSessionPool.Sessions, 9, "%d sessions in testSessionPool, should be 9 total sessions", len(testSessionPool.Sessions))

		// Requires a region to provided
		var nilRegion []string
		testSessionPool = NewPool(profiles, nilRegion, logger)

		assert.Lenf(testSessionPool.Sessions, 0, "%d sessions in testSessionPool, should be 0 total sessions", len(testSessionPool.Sessions))
	})

	t.Run("test session deduplication", func(t *testing.T) {
		// Because we're using a map to store all the sessions, duplicate sessions should be filtered out
		// This should create 4 unique sessions
		regions := []string{"us-test-1", "us-test-2", "us-test-2"}
		profiles := []string{"profile1", "profile2", "profile2"}
		testSessionPool := NewPool(profiles, regions, logger)

		assert.Lenf(testSessionPool.Sessions, 4, "%d sessions in testSessionPool, should be 4 unique sessions", len(testSessionPool.Sessions))
	})

}
