package main

import (
	"context"
	"log"
	"strconv"
	"time"
)

type accountResolver struct {
	server *Server
}

func (resolver *accountResolver) Orders(ctx context.Context, obj *Account) ([]*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	orderList, err := resolver.server.orderClient.GetOrdersForAccount(ctx, obj.ID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var orders []*Order
	for _, order := range orderList {
		var products []*OrderedProduct
		for _, orderedProduct := range order.Products {
			products = append(products, &OrderedProduct{
				ID:          strconv.Itoa(int(orderedProduct.ID)),
				Name:        orderedProduct.Name,
				Description: orderedProduct.Description,
				Price:       orderedProduct.Price,
				Quantity:    int(orderedProduct.Quantity),
			})
		}
		orders = append(orders, &Order{
			ID:         strconv.Itoa(int(order.ID)),
			CreatedAt:  order.CreatedAt,
			TotalPrice: order.TotalPrice,
			Products:   products,
		})
	}

	return orders, nil
}
