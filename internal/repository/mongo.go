package repository

import (
	"context"
	"time"

	"booking-system/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingRepository struct {
	Collection *mongo.Collection
}

func (r *BookingRepository) Save(ctx context.Context, b *models.Booking) error {
	b.CreatedAt = time.Now()
	_, err := r.Collection.InsertOne(ctx, b)
	return err
}

func (r *BookingRepository) GetAll(ctx context.Context) ([]models.Booking, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var bookings []models.Booking
	err = cursor.All(ctx, &bookings)
	return bookings, err
}

func (r *BookingRepository) Exists(ctx context.Context, userID, ref string) (bool, error) {
	filter := bson.M{
		"user_id":     userID,
		"booking_ref": ref,
	}

	count, err := r.Collection.CountDocuments(ctx, filter)
	return count > 0, err
}

func (r *BookingRepository) GetUpcoming(ctx context.Context) ([]models.Booking, error) {

	filter := bson.M{
		"start_date": bson.M{
			"$gte": time.Now(),
		},
	}

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var bookings []models.Booking
	err = cursor.All(ctx, &bookings)
	return bookings, err
}