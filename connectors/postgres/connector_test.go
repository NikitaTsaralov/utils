package postgres

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

type TraceSuite struct {
	suite.Suite
}

func (s *TraceSuite) TestStartStop() {
	db, err := Start(Config{
		Password: "dev",
	})
	s.Require().Nil(err)

	defer func(db *sqlx.DB) {
		err := db.Close()
		s.Require().Nil(err)
	}(db)

	s.Require().NotNil(db)
}

func TestTraceSuite(t *testing.T) {
	suite.Run(t, new(TraceSuite))
}
