package down

import (
	"embed"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/spf13/cobra"
	"io/fs"
	"log"
)

var (
	steps *int
	dsn   *string

	Cmd = &cobra.Command{
		Use:   "down",
		Short: "Откатить миграцию",
		Long: "Выполнение отката миграции в mysql на указанное количество шагов назад.\n" +
			"Если шагов не указано, то выполняется откат на один шаг назад.",
		Run: func(cmd *cobra.Command, args []string) {
			// nolint
			rootFS := cmd.Context().Value("fs::resources").(embed.FS)
			resourcesFS, err := fs.Sub(rootFS, "resources")

			if err != nil {
				log.Fatalf("ошибка при работе с resources: %s", err.Error())
			}

			d, err := iofs.New(resourcesFS, "database/migrations")

			if err != nil {
				log.Fatalf("ошибка при чтении с resources: %s", err.Error())
			}

			if *dsn == "" {
				log.Fatal("не указан DSN для подключения к БД")
			}

			m, err := migrate.NewWithSourceInstance("iofs", d, "mysql://"+*dsn)

			if err != nil {
				log.Fatalf("ошибка при подключении к БД: %s", err.Error())
			}

			if *steps == 0 {
				err = m.Down()
			} else if *steps < 0 {
				err = m.Steps(*steps)
			} else {
				err = m.Steps(-*steps)
			}

			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("нет миграций для отката")
			} else if err != nil {
				log.Fatalf("ошибка при откате миграции: %s", err.Error())
			} else {
				version, _, _ := m.Version()

				fmt.Println("откат миграции прошел успешно")
				fmt.Printf("текущая версия БД: %d", version)
			}
		},
	}
)

func init() {
	steps = Cmd.Flags().IntP(
		"steps", "s",
		1,
		"количество шагов миграции",
	)
	dsn = Cmd.Flags().StringP(
		"dsn", "d",
		"",
		"DSN для подключения к БД",
	)
}
