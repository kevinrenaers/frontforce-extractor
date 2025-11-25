package influx

import (
	"context"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/rs/zerolog/log"

	"renaers.be/frontforce/internal/configuration"
)

type Influx struct {
	client influxdb2.Client
	org    string
	bucket string
}

func NewInflux(config configuration.Configuration) Influx {
	return Influx{
		client: influxdb2.NewClient(config.InfluxUrl, config.InfluxToken),
		org:    config.InfluxOrg,
		bucket: config.InfluxBucket,
	}
}

func (i Influx) WriteStats(stats map[string]any) {
	writeAPI := i.client.WriteAPIBlocking(i.org, i.bucket)

	for key, value := range stats {
		mappedStat := mapStatToMeasurement(key)
		if mappedStat == "" {
			continue
		}
		p := influxdb2.NewPointWithMeasurement(mappedStat).
			AddTag("location", "sint-truiden").
			AddField("value", value)
		err := writeAPI.WritePoint(context.Background(), p)
		if err != nil {
			log.Error().Err(err).Msgf("influx - failed writing stat %s to influxdb", key)
		}
	}
	// Flush writes
	err := writeAPI.Flush(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("influx - failed flushing data")
	}
}

func mapStatToMeasurement(stat string) string {
	switch stat {
	case "post_availability_percentage":
		return "post.availability_percentage"
	case "post_availability_count":
		return "post.amount_available"
	case "post_callable_count":
		return "post.amount_callable"
	case "post_unavailability_count":
		return "post.amount_unavailable"
	case "Beschikbaar Werk":
		return "post.amount_work"
	case "Beschikbaar snel":
		return "post.amount_fast"
	case "Beschikbaar traag":
		return "post.amount_slow"
	case "Dienstopdracht":
		return "post.amount_dno"
	case "Reserve":
		return "post.amount_reserve"
	case "Kazerne Operationeel":
		return "post.amount_operational"
	case "Kazerne Administratief":
		return "post.amount_administrative"
	case "Officier van dienst":
		return "post.amount_officer"
	case "Officier beschikbaar":
		return "post.amount_officer_available"
	case "Kazerne beschikbaar":
		return "post.amount_on_post"
	case "Ambulance":
		return "post.amount_ambulance"
	case "Niet Operationeel":
		return "post.amount_not_operational"
	case "Interventie":
		return "post.amount_intervention"
	case "Elders beschikbaar":
		return "post.amount_available_elsewhere"
	default:
		return ""
	}
}
