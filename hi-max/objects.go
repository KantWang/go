package objectsapi

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"

	objects "objects/gen/objects"
)

// objects service example implementation.
// The example methods log the requests and return zero values.
type objectssrvc struct {
	logger *log.Logger
	db     *sql.DB
}

// NewObjects returns the objects service implementation.
func NewObjects(logger *log.Logger, db *sql.DB) objects.Service {
	return &objectssrvc{logger, db}
}

// sql querys
const (
	// selectByObservationTime = `SELECT * FROM near_earth_objects WHERE observation_time > $1 order by observation_time offset $2 rows fetch next 5 rows only`
	selectByObservationTime = `SELECT * FROM near_earth_objects WHERE observation_time > $1 order by observation_time offset $2 rows limit 5`
	selectByID              = `SELECT * FROM near_earth_objects WHERE id = $1`
	selectAll               = `SELECT * FROM near_earth_objects`
	addObject               = `INSERT INTO near_earth_objects (id, name, observation_time) VALUES ($1, $2, $3)`
	updateObject            = `UPDATE near_earth_objects SET name=$2, observation_time=$3 WHERE id=$1`
	updateName              = `UPDATE near_earth_objects SET name=$2 WHERE id=$1`
	updateTime              = `UPDATE near_earth_objects SET observation_time=$2 WHERE id=$1`
	deleteObject            = `DELETE FROM near_earth_objects WHERE id=$1 RETURNING id`
)

// formating, parsing forms
const (
	observationTimeFormating = "2006-01-02 15:04:05 -0700 MST"
	observationDateParsing   = "2006-01-02"
)

// List implements list.
func (s *objectssrvc) List(ctx context.Context, p *objects.ListPayload) (res objects.NearEarthObjectManagementCollection, err error) {
	// ObservationDate check
	observationTime, err := time.Parse(observationDateParsing, p.ObservationDate)
	if err != nil {
		s.logger.Printf("ObservationDate format is wrong.")
		return nil, objects.MakeWrongDateFormat(fmt.Errorf("ObservationDate format is wrong: %s", p.ObservationDate))
	}

	// Offset check
	if p.Offset < 0 {
		s.logger.Printf("Negative offsets are not allowed.")
		return nil, objects.MakeWrongOffsetFormat(fmt.Errorf("Negative offsets are not allowed: %d", p.Offset))
	}

	result, err := s.db.Query(
		selectByObservationTime,
		observationTime.Unix()*1000,
		p.Offset*5, // 5 objects per 1 page
	)
	if err != nil {
		s.logger.Print(err)
		return nil, objects.MakeInternalServer(fmt.Errorf("DB error"))
	}

	for result.Next() {
		tmp := new(objects.NearEarthObjectManagement)
		err = result.Scan(&tmp.ID, &tmp.Name, &tmp.ObservationTime)
		if err != nil {
			s.logger.Print(err)
			return nil, objects.MakeInternalServer(fmt.Errorf("DB error"))
		}

		tmp.ObservationDate = time.UnixMilli(tmp.ObservationTime).Format(observationTimeFormating)
		res = append(res, tmp)
	}

	return
}

// Find implements find.
func (s *objectssrvc) Find(ctx context.Context, p *objects.FindPayload) (res *objects.NearEarthObjectManagement, err error) {
	// id check
	if p.ID <= 0 {
		s.logger.Print("Negative ID are not allowed")
		return nil, objects.MakeWrongID(fmt.Errorf("Negative ID are not allowed: %d", p.ID))
	}

	res = new(objects.NearEarthObjectManagement)

	result := s.db.QueryRow(
		selectByID,
		p.ID,
	)

	err = result.Scan(&res.ID, &res.Name, &res.ObservationTime)
	if err != nil {
		s.logger.Print(err) // sql: not found error
		switch {
		case err == sql.ErrNoRows:
			return nil, objects.MakeNotFound(fmt.Errorf("The received ID does not exist: %d", p.ID))
		default:
			return nil, objects.MakeInternalServer(fmt.Errorf("DB error"))
		}
	}
	res.ObservationDate = time.UnixMilli(res.ObservationTime).Format(observationTimeFormating)

	return
}

// Add implements add.
func (s *objectssrvc) Add(ctx context.Context, p *objects.AddPayload) (res bool, err error) {
	// id, observationTime check
	switch {
	case p.ID <= 0:
		s.logger.Print("Negative ID are not allowed")
		return false, objects.MakeWrongID(fmt.Errorf("Negative ID are not allowed: %d", p.ID))
	case p.ObservationTime <= 0:
		s.logger.Print("Negative ObservationTime are not allowed")
		return false, objects.MakeWrongTime(fmt.Errorf("Negative ObservationTime are not allowed: %d", p.ObservationTime))
	}

	// insert
	_, err = s.db.Exec(
		addObject,
		p.ID,
		p.Name,
		p.ObservationTime,
	)

	if err != nil {
		switch {
		case err.(*pq.Error).Code == "23505": // unique_violation error - 23505
			s.logger.Printf("id: %d is already exists", p.ID)
			return false, objects.MakeAlreadyExists(fmt.Errorf("already exists"))
		default:
			s.logger.Print(err)
			return false, objects.MakeInternalServer(fmt.Errorf("DB error"))
		}
	}

	return true, nil
}

// Update implements update.
func (s *objectssrvc) Update(ctx context.Context, p *objects.UpdatePayload) (res bool, err error) {
	switch {
	case p.ID <= 0:
		s.logger.Print("Negative ID are not allowed")
		return false, objects.MakeWrongID(fmt.Errorf("Negative ID are not allowed: %d", p.ID))
	case p.ObservationTime != nil && *p.ObservationTime <= 0:
		s.logger.Print("Negative ObservationTime are not allowed")
		return false, objects.MakeWrongTime(fmt.Errorf("Negative ObservationTime are not allowed: %d", p.ObservationTime))
	case p.Name == nil && p.ObservationTime == nil:
		s.logger.Print("You must enter at least one")
		return false, objects.MakeNothingToUpdate(fmt.Errorf("You must enter at least one"))
	}

	// check exist or not
	var id, observationTime int
	var name string

	result := s.db.QueryRow(
		selectByID,
		p.ID,
	)

	err = result.Scan(&id, &name, &observationTime)
	if err != nil {
		s.logger.Print(err)
		if err == sql.ErrNoRows {
			return false, objects.MakeNotFound(fmt.Errorf("The received ID does not exist: %d", p.ID))
		} else {
			return false, objects.MakeInternalServer(fmt.Errorf("DB error"))
		}
	}

	// update
	switch {
	case p.Name != nil && p.ObservationTime != nil:
		s.logger.Printf("Update completed! id: %d, name: %s -> %s, time: %d -> %d", p.ID, name, *p.Name, observationTime, *p.ObservationTime)
		_, err = s.db.Exec(updateObject, p.ID, p.Name, p.ObservationTime)
	case p.Name == nil && p.ObservationTime != nil:
		s.logger.Printf("Update completed! id: %d, time: %d -> %d", p.ID, observationTime, *p.ObservationTime)
		_, err = s.db.Exec(updateTime, p.ID, p.ObservationTime)
	case p.Name != nil && p.ObservationTime == nil:
		s.logger.Printf("Update completed! id: %d, name: %s -> %s", p.ID, name, *p.Name)
		_, err = s.db.Exec(updateName, p.ID, p.Name)
	}

	if err != nil {
		s.logger.Print(err)
		return false, objects.MakeInternalServer(fmt.Errorf("DB error"))
	}

	return true, nil
}

// Delete implements delete.
func (s *objectssrvc) Delete(ctx context.Context, p *objects.DeletePayload) (res bool, err error) {
	// id check
	if p.ID <= 0 {
		s.logger.Print("Negative ID are not allowed")
		return false, objects.MakeWrongID(fmt.Errorf("Negative ID are not allowed: %d", p.ID))
	}

	// check exist or not
	var id, observationTime int
	var name string

	result := s.db.QueryRow(
		selectByID,
		p.ID,
	)

	err = result.Scan(&id, &name, &observationTime)
	if err != nil {
		s.logger.Print(err)
		switch {
		case err == sql.ErrNoRows:
			return false, objects.MakeNotFound(fmt.Errorf("The received ID does not exist: %d", p.ID))
		default:
			return false, objects.MakeInternalServer(fmt.Errorf("DB error"))
		}
	}

	// delete
	_, err = s.db.Exec(deleteObject, p.ID)
	if err != nil {
		s.logger.Print(err)
		return false, objects.MakeInternalServer(fmt.Errorf("DB error"))
	}

	return true, nil
}
