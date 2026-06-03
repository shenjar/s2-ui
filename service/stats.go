package service

import (
	"time"

	"github.com/shenaba/s2-ui/database"
	"github.com/shenaba/s2-ui/database/model"

	"gorm.io/gorm"
)

type onlines struct {
	Inbound  []string `json:"inbound,omitempty"`
	User     []string `json:"user,omitempty"`
	Outbound []string `json:"outbound,omitempty"`
}

var onlineResources = &onlines{}

type StatsService struct {
}

func (s *StatsService) SaveStats(enableTraffic bool) error {
	if corePtr == nil || !corePtr.IsRunning() {
		return nil
	}
	box := corePtr.GetInstance()
	if box == nil {
		return nil
	}
	st := box.StatsTracker()
	if st == nil {
		return nil
	}
	stats := st.GetStats()

	// Reset onlines
	onlineResources.Inbound = nil
	onlineResources.Outbound = nil
	onlineResources.User = nil

	if len(*stats) == 0 {
		return nil
	}

	var err error
	db := database.GetDB()
	tx := db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	for _, stat := range *stats {
		if stat.Resource == "user" {
			if stat.Direction {
				err = tx.Model(model.Client{}).Where("name = ?", stat.Tag).
					UpdateColumn("up", gorm.Expr("up + ?", stat.Traffic)).Error
			} else {
				err = tx.Model(model.Client{}).Where("name = ?", stat.Tag).
					UpdateColumn("down", gorm.Expr("down + ?", stat.Traffic)).Error
			}
			if err != nil {
				return err
			}
		}
		if stat.Direction {
			switch stat.Resource {
			case "inbound":
				onlineResources.Inbound = append(onlineResources.Inbound, stat.Tag)
			case "outbound":
				onlineResources.Outbound = append(onlineResources.Outbound, stat.Tag)
			case "user":
				onlineResources.User = append(onlineResources.User, stat.Tag)
			}
		}
	}

	if !enableTraffic {
		return nil
	}
	err = tx.Create(&stats).Error
	return err
}

func (s *StatsService) GetStats(resource string, tag string, period string) ([]model.Stats, error) {
	now := time.Now().Unix()
	var bucketSec int64
	var startTime int64
	switch period {
	case "day":
		bucketSec = 3600
		startTime = now - 86400
	case "month":
		bucketSec = 86400
		startTime = now - 86400*30
	default: // "hour"
		bucketSec = 60
		startTime = now - 3600
	}

	db := database.GetDB()
	resources := []string{resource}
	if resource == "endpoint" {
		resources = []string{"inbound", "outbound"}
	}

	type bucketRow struct {
		Bucket    int64
		Direction bool
		Traffic   int64
	}
	var rows []bucketRow
	err := db.Raw(
		`SELECT (date_time / ?) * ? AS bucket, direction, SUM(traffic) AS traffic
		 FROM stats
		 WHERE resource IN ? AND tag = ? AND date_time > ? AND date_time <= ?
		 GROUP BY bucket, direction
		 ORDER BY bucket`,
		bucketSec, bucketSec, resources, tag, startTime, now,
	).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	// Build lookup map
	type key struct {
		bucket    int64
		direction bool
	}
	lookup := make(map[key]int64, len(rows))
	for _, r := range rows {
		lookup[key{r.Bucket, r.Direction}] = r.Traffic
	}

	// Fill all buckets including empty ones so x-axis is evenly distributed
	firstBucket := (startTime / bucketSec) * bucketSec
	var result []model.Stats
	for b := firstBucket; b <= now; b += bucketSec {
		for _, dir := range []bool{false, true} {
			result = append(result, model.Stats{
				DateTime:  b,
				Resource:  resource,
				Tag:       tag,
				Direction: dir,
				Traffic:   lookup[key{b, dir}],
			})
		}
	}
	return result, nil
}

func (s *StatsService) GetOnlines() (onlines, error) {
	return *onlineResources, nil
}
func (s *StatsService) DelOldStats(days int) error {
	oldTime := time.Now().AddDate(0, 0, -(days)).Unix()
	db := database.GetDB()
	return db.Where("date_time < ?", oldTime).Delete(model.Stats{}).Error
}
