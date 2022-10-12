// Code generated by goa v3.10.0, DO NOT EDIT.
//
// objects service
//
// Command:
// $ goa gen objects/design

package objects

import (
	"context"
	objectsviews "objects/gen/objects/views"

	goa "goa.design/goa/v3/pkg"
)

// API Server
type Service interface {
	// List implements list.
	List(context.Context, *ListPayload) (res NearEarthObjectManagementCollection, err error)
	// Find implements find.
	Find(context.Context, *FindPayload) (res *NearEarthObjectManagement, err error)
	// Add implements add.
	Add(context.Context, *AddPayload) (res bool, err error)
	// Update implements update.
	Update(context.Context, *UpdatePayload) (res bool, err error)
	// Delete implements delete.
	Delete(context.Context, *DeletePayload) (res bool, err error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "objects"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [5]string{"list", "find", "add", "update", "delete"}

// AddPayload is the payload type of the objects service add method.
type AddPayload struct {
	ID              int
	Name            string
	ObservationTime int64
}

// DeletePayload is the payload type of the objects service delete method.
type DeletePayload struct {
	ID int
}

// FindPayload is the payload type of the objects service find method.
type FindPayload struct {
	ID int
}

// ListPayload is the payload type of the objects service list method.
type ListPayload struct {
	ObservationDate string
	Offset          int
}

// NearEarthObjectManagement is the result type of the objects service find
// method.
type NearEarthObjectManagement struct {
	// object ID
	ID int
	// object Name
	Name string
	// object observation time
	ObservationTime int64
	// object observation date
	ObservationDate string
}

// NearEarthObjectManagementCollection is the result type of the objects
// service list method.
type NearEarthObjectManagementCollection []*NearEarthObjectManagement

// UpdatePayload is the payload type of the objects service update method.
type UpdatePayload struct {
	ID              int
	Name            *string
	ObservationTime *int64
}

// MakeWrongDateFormat builds a goa.ServiceError from an error.
func MakeWrongDateFormat(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "wrong_date_format", false, false, false)
}

// MakeWrongOffsetFormat builds a goa.ServiceError from an error.
func MakeWrongOffsetFormat(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "wrong_offset_format", false, false, false)
}

// MakeInternalServer builds a goa.ServiceError from an error.
func MakeInternalServer(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "internal_server", false, false, false)
}

// MakeNotFound builds a goa.ServiceError from an error.
func MakeNotFound(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "not_found", false, false, false)
}

// MakeWrongID builds a goa.ServiceError from an error.
func MakeWrongID(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "wrong_id", false, false, false)
}

// MakeWrongTime builds a goa.ServiceError from an error.
func MakeWrongTime(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "wrong_time", false, false, false)
}

// MakeAlreadyExists builds a goa.ServiceError from an error.
func MakeAlreadyExists(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "already_exists", false, false, false)
}

// MakeNothingToUpdate builds a goa.ServiceError from an error.
func MakeNothingToUpdate(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "nothing_to_update", false, false, false)
}

// NewNearEarthObjectManagementCollection initializes result type
// NearEarthObjectManagementCollection from viewed result type
// NearEarthObjectManagementCollection.
func NewNearEarthObjectManagementCollection(vres objectsviews.NearEarthObjectManagementCollection) NearEarthObjectManagementCollection {
	return newNearEarthObjectManagementCollection(vres.Projected)
}

// NewViewedNearEarthObjectManagementCollection initializes viewed result type
// NearEarthObjectManagementCollection from result type
// NearEarthObjectManagementCollection using the given view.
func NewViewedNearEarthObjectManagementCollection(res NearEarthObjectManagementCollection, view string) objectsviews.NearEarthObjectManagementCollection {
	p := newNearEarthObjectManagementCollectionView(res)
	return objectsviews.NearEarthObjectManagementCollection{Projected: p, View: "default"}
}

// NewNearEarthObjectManagement initializes result type
// NearEarthObjectManagement from viewed result type NearEarthObjectManagement.
func NewNearEarthObjectManagement(vres *objectsviews.NearEarthObjectManagement) *NearEarthObjectManagement {
	return newNearEarthObjectManagement(vres.Projected)
}

// NewViewedNearEarthObjectManagement initializes viewed result type
// NearEarthObjectManagement from result type NearEarthObjectManagement using
// the given view.
func NewViewedNearEarthObjectManagement(res *NearEarthObjectManagement, view string) *objectsviews.NearEarthObjectManagement {
	p := newNearEarthObjectManagementView(res)
	return &objectsviews.NearEarthObjectManagement{Projected: p, View: "default"}
}

// newNearEarthObjectManagementCollection converts projected type
// NearEarthObjectManagementCollection to service type
// NearEarthObjectManagementCollection.
func newNearEarthObjectManagementCollection(vres objectsviews.NearEarthObjectManagementCollectionView) NearEarthObjectManagementCollection {
	res := make(NearEarthObjectManagementCollection, len(vres))
	for i, n := range vres {
		res[i] = newNearEarthObjectManagement(n)
	}
	return res
}

// newNearEarthObjectManagementCollectionView projects result type
// NearEarthObjectManagementCollection to projected type
// NearEarthObjectManagementCollectionView using the "default" view.
func newNearEarthObjectManagementCollectionView(res NearEarthObjectManagementCollection) objectsviews.NearEarthObjectManagementCollectionView {
	vres := make(objectsviews.NearEarthObjectManagementCollectionView, len(res))
	for i, n := range res {
		vres[i] = newNearEarthObjectManagementView(n)
	}
	return vres
}

// newNearEarthObjectManagement converts projected type
// NearEarthObjectManagement to service type NearEarthObjectManagement.
func newNearEarthObjectManagement(vres *objectsviews.NearEarthObjectManagementView) *NearEarthObjectManagement {
	res := &NearEarthObjectManagement{}
	if vres.ID != nil {
		res.ID = *vres.ID
	}
	if vres.Name != nil {
		res.Name = *vres.Name
	}
	if vres.ObservationTime != nil {
		res.ObservationTime = *vres.ObservationTime
	}
	if vres.ObservationDate != nil {
		res.ObservationDate = *vres.ObservationDate
	}
	return res
}

// newNearEarthObjectManagementView projects result type
// NearEarthObjectManagement to projected type NearEarthObjectManagementView
// using the "default" view.
func newNearEarthObjectManagementView(res *NearEarthObjectManagement) *objectsviews.NearEarthObjectManagementView {
	vres := &objectsviews.NearEarthObjectManagementView{
		ID:              &res.ID,
		Name:            &res.Name,
		ObservationTime: &res.ObservationTime,
		ObservationDate: &res.ObservationDate,
	}
	return vres
}
