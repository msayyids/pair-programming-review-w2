package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	DB *mongo.Collection
}

func NewRepository(db *mongo.Collection) Repository {
	return Repository{DB: db}
}

// type Products struct {
// 	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
// 	Name        string             `bson:"name" json:"name"`
// 	Description string             `bson:"description" json:"description"`
// 	Price       float64            `bson:"price" json:"price"`
// 	Stock       int                `bson:"stock" json:"stock"`
// }

func (r Repository) CreateProduct(name, description string, price float32, stock int) (Products, error) {
	newProduct := Products{Name: name, Description: description, Price: price, Stock: stock}

	result, err := r.DB.InsertOne(context.Background(), newProduct)
	if err != nil {
		return Products{}, err
	}

	newProduct.ID = result.InsertedID.(primitive.ObjectID)

	return newProduct, nil
}

func (r Repository) GetAllProduct() ([]Products, error) {
	var products []Products

	cursor, err := r.DB.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var product Products
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r Repository) UpdateProduct(id, name, description string, price float32, stock int) (Products, error) {
	newObjectID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": newObjectID}

	update := bson.M{
		"$set": bson.M{"name": name, "description": description, "price": price, "stock": stock},
	}

	_, err := r.DB.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return Products{}, err
	}

	// Create a Products object with the updated values
	updatedProduct := Products{
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
	}

	return updatedProduct, nil
}

func (r Repository) DeleteProduct(id string) error {
	newObjectID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": newObjectID}

	_, err := r.DB.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return nil
}
