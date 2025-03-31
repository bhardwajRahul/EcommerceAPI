package main

import (
	"context"
	"log"
	"strconv"
	"time"
)

type queryResolver struct {
	server *Server
}

func (resolver *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id != nil {
		res, err := resolver.server.accountClient.GetAccount(ctx, *id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return []*Account{{
			ID:   strconv.Itoa(int(res.ID)),
			Name: res.Name,
		}}, nil
	}

	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	accountList, err := resolver.server.accountClient.GetAccounts(ctx, skip, take)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var accounts []*Account
	for _, a := range accountList {
		account := &Account{
			ID:    strconv.Itoa(int(a.ID)),
			Name:  a.Name,
			Email: a.Email,
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (resolver *queryResolver) Product(ctx context.Context, pagination *PaginationInput, query, id *string) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Get single
	if id != nil {
		res, err := resolver.server.productClient.GetProduct(ctx, *id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return []*Product{{
			ID:          res.ID,
			Name:        res.Name,
			Description: res.Description,
			Price:       res.Price,
		}}, nil
	}
	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	q := ""
	if query != nil {
		q = *query
	}
	productList, err := resolver.server.productClient.GetProducts(ctx, skip, take, nil, q)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var products []*Product
	for _, product := range productList {
		products = append(products,
			&Product{
				ID:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
			},
		)
	}

	return products, nil
}

func (p PaginationInput) bounds() (uint64, uint64) {
	skipValue := uint64(0)
	takeValue := uint64(100)
	if p.Skip != 0 {
		skipValue = uint64(p.Skip)
	}
	if p.Take != 100 {
		takeValue = uint64(p.Take)
	}
	return skipValue, takeValue
}
