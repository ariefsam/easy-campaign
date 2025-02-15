package idgenerator_test

import (
	"campaign/idgenerator"
	"campaign/logger"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_idgenerator_Generate(t *testing.T) {
	idgen := idgenerator.New()
	ctx := context.TODO()
	id := idgen.Generate(ctx)
	assert.NotEmpty(t, id)

	logger.PrintJSON(id)
}

func Test_idgenerator_GenerateUUID(t *testing.T) {
	idgen := idgenerator.New()
	ctx := context.TODO()
	id := idgen.GenerateUUID(ctx)
	assert.NotEmpty(t, id)

	logger.PrintJSON(id)
}
