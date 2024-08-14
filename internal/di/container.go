package di

import (
	"conduit/internal/repository/postgres"
	"conduit/internal/services"

	"github.com/sarulabs/di"
)

func BuildContainer() (di.Container, error) {
	builder, _ := di.NewBuilder()

	builder.Add(di.Def{
		Name: "db",
		Build: func(ctn di.Container) (interface{}, error) {
			return postgres.NewDB()
		},
	})

	builder.Add(di.Def{
		Name: "strategy_service",
		Build: func(ctn di.Container) (interface{}, error) {
			return services.NewStrategyService(ctn.Get("db").(*postgres.DB)), nil
		},
	})

	return builder.Build()
}
