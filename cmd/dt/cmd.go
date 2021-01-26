package dt

import (
	"github.com/Newlooc/fundtools/pkg/apis"
	"github.com/Newlooc/fundtools/pkg/dt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"time"
)

var (
	start string
	end   string
	code  string
	frq   int
	money float64
	rate  float64
)

func NewDTCmd() *cobra.Command {
	addFlags(dtCmd)
	startTime, err := time.Parse(apis.DateFormat, start)
	if err != nil {
		log.WithError(err).Fatalf("Date format should be YYYY-MM-DD. %s given", start)
	}
	endTime, err := time.Parse(apis.DateFormat, end)
	if err != nil {
		log.WithError(err).Fatalf("Date format should be YYYY-MM-DD. %s given", start)
	}
	dtCmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {

			dtm, err := dt.NewDTManager(code, startTime, endTime, frq, money, rate)
			if err != nil {
				log.WithError(err).Fatal("Manager init with error")
			}

			if err := dtm.Run(); err != nil {
				log.WithError(err).Fatal("Manager run with error")
			}
		},
	}

	return dtCmd
}

func validateFlags() {
}

func addFlags(dtCmd *cobra.Command) {
	dtCmd.PersistentFlags().StringVar(&start, "start", "", "YYYY-MM-DD")
	dtCmd.PersistentFlags().StringVar(&end, "end", "", "YYYY-MM-DD")
	dtCmd.PersistentFlags().StringVar(&code, "code", "", "code")
	dtCmd.PersistentFlags().Float64Var(&money, "money", 100, "money")
	dtCmd.PersistentFlags().IntVar(&frq, "frq", 1, "frq")
	dtCmd.PersistentFlags().Float64Var(&rate, "rate", 0.15, "rate")
}
