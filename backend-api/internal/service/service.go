package service

import "shop/backend-api/internal/repository"

func NewAuthService(dao repository.DAO) AuthService {
	return &authService{
		userRepo: dao.NewUserRepo(),
	}
}

func NewProductService(dao repository.DAO) ProductService {
	return &productService{
		productRepo: dao.NewProductRepo(),
	}
}
