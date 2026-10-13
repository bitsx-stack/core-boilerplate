package main

import (
	"fmt"
	"log"
	"net/http"

	"bitpos.com/m/config"
	"bitpos.com/m/database"
	"bitpos.com/m/handlers"
	"bitpos.com/m/middlewares"
	"bitpos.com/m/repos"
)

func main() {

	cfg := config.Load()

	db, err := database.UnganishaDatabase()
	if err != nil {
		log.Fatalf("Database! Disconnected: %v", err)
	}
	defer db.Close()

	log.Println("Database! Connected")

	router := http.NewServeMux()

	userRepo := repos.NewUserRepo(db)
	categoryRepo := repos.NewCategoryRepo(db)
	categoryTypeRepo := repos.NewCategoryTypeRepo(db)
	settingRepo := repos.NewSettingRepo(db)
	peopleRepo := repos.NewPeopleRepo(db)
	regionRepo := repos.NewRegionRepo(db)
	itemRepo := repos.NewItemRepo(db)
	transRepo := repos.NewTransactionRepo(db)
	salesRepo := repos.NewSalesRepo(db)
	bulkRepo := repos.NewBulkRepo(db)
	storeRepo := repos.NewStoreRepo(db)
	billRepo := repos.NewBillRepo(db)
	stockRepo := repos.NewStockRepo(db)
	locationRepo := repos.NewLocationRepo(db)
	roleRepo := repos.NewRoleRepo(db)
	userAccessRepo := repos.NewUserAccessRepo(db)
	paymentRepo := repos.NewPaymentRepo(db)
	dashboardRepo := repos.NewDashboardRepo(db)
	reportRepo := repos.NewReportRepo(db)


	authHandler := handlers.NewAuthHandler(userRepo, cfg)
	userHandler := handlers.NewUserHandler(userRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)
	categoryTypeHandler := handlers.NewCategoryTypeHandler(categoryTypeRepo)
	settingHandler := handlers.NewSettingHandler(settingRepo)
	peopleHandler := handlers.NewPeopleHandler(peopleRepo)
	regionHandler := handlers.NewRegionHandler(regionRepo)
	itemHandler := handlers.NewItemHandler(itemRepo)
	salesHandler := handlers.NewSalesHandler(salesRepo)
	bulkHandler := handlers.NewBulkHandler(bulkRepo)
	storeHandler := handlers.NewStoreHandler(storeRepo)
	transHandler := handlers.NewTransactionHandler(transRepo, billRepo, stockRepo)
	locationHandler := handlers.NewLocationHandler(locationRepo)
	roleHandler := handlers.NewRoleHandler(roleRepo)
	userAccessHandler := handlers.NewUserAccessHandler(userAccessRepo)
	billHandler := handlers.NewBillHandler(billRepo, paymentRepo, stockRepo)
	dashboardHandler := handlers.NewDashboardHandler(dashboardRepo)
	reportHandler := handlers.NewReportHandler(reportRepo)

    router.HandleFunc("GET /api/reports/stock", middlewares.AuthMiddleware(cfg)(reportHandler.StockReport))

	router.HandleFunc("GET /api/dashboard", middlewares.AuthMiddleware(cfg)(dashboardHandler.GetDashboard))

	router.HandleFunc("GET /api/bills", billHandler.GetAllBills)
	router.HandleFunc("POST /api/bills/create", middlewares.AuthMiddleware(cfg)(billHandler.CreateBill))
	router.HandleFunc("GET /api/bills/single/", middlewares.AuthMiddleware(cfg)(billHandler.GetSingleBill))
	router.HandleFunc("PUT /api/bills/update/", middlewares.AuthMiddleware(cfg)(billHandler.UpdateBillItems))
	router.HandleFunc("POST /api/bills/pay/", middlewares.AuthMiddleware(cfg)(billHandler.MakePayment))
	router.HandleFunc("GET /api/bills/payments/", middlewares.AuthMiddleware(cfg)(billHandler.GetPayments))


	router.HandleFunc(
		"GET /api/roles",
		middlewares.AuthMiddleware(cfg)(
			middlewares.RequirePermission(userAccessRepo, "ROLES_READ")(
				roleHandler.GetAllRoles,
			),
		),
	)
	router.HandleFunc("PUT /api/roles/update/{id}", middlewares.AuthMiddleware(cfg)(roleHandler.UpdateRole))
	router.HandleFunc("GET /api/roles/single/{id}", middlewares.AuthMiddleware(cfg)(roleHandler.GetSingleRole))
	router.HandleFunc("POST /api/roles/create/", middlewares.AuthMiddleware(cfg)(roleHandler.CreateRole))



	router.HandleFunc("GET /api/users", userHandler.GetAllUsers)
	router.HandleFunc("GET /api/users/single/", userHandler.GetUser)
	router.HandleFunc("POST /api/users/create", userHandler.CreateUser)
	router.HandleFunc("PUT /api/users/update/", userHandler.UpdateUser)
	router.HandleFunc("DELETE /api/users/delete/", userHandler.DeleteUser)


	router.HandleFunc("PUT /api/users/assign-role/{id}", userAccessHandler.AssignRole)
	router.HandleFunc("PUT /api/users/permissions/{id}", userAccessHandler.UpdatePermissions)


	router.HandleFunc("GET /api/locations", locationHandler.GetAllLocations)

	router.HandleFunc("POST /api/transactions/create", middlewares.AuthMiddleware(cfg)(transHandler.CreateTransaction))
	router.HandleFunc("PUT /api/transactions/complete/{id}", middlewares.AuthMiddleware(cfg)(transHandler.CompleteTransaction))
	router.HandleFunc("GET /api/transactions/", middlewares.AuthMiddleware(cfg)(transHandler.GetAllTransactions))


	router.HandleFunc("GET /api/stores/", middlewares.AuthMiddleware(cfg)(storeHandler.GetAllStores))
	router.HandleFunc("GET /api/stores/single/{id}", middlewares.AuthMiddleware(cfg)(storeHandler.GetSingleStore))
	router.HandleFunc("POST /api/stores/create", middlewares.AuthMiddleware(cfg)(storeHandler.CreateStore))
	router.HandleFunc("PUT /api/stores/update/{id}", middlewares.AuthMiddleware(cfg)(storeHandler.UpdateStore))
	router.HandleFunc("PUT /api/stores/delete/{id}", middlewares.AuthMiddleware(cfg)(storeHandler.DeleteStore))

	router.HandleFunc("POST /api/bulk/items", middlewares.AuthMiddleware(cfg)(bulkHandler.UploadItems))
	router.HandleFunc("POST /api/bulk/stock", middlewares.AuthMiddleware(cfg)(bulkHandler.UploadStock))



	router.HandleFunc("POST /api/sales", middlewares.AuthMiddleware(cfg)(salesHandler.CreateSale))

	router.HandleFunc("POST /api/items/create", middlewares.AuthMiddleware(cfg)(itemHandler.CreateItem))
	router.HandleFunc("GET /api/items", middlewares.AuthMiddleware(cfg)(itemHandler.GetAllItems))
	router.HandleFunc("GET /api/items/single/{id}", middlewares.AuthMiddleware(cfg)(itemHandler.GetSingleItem))
	router.HandleFunc("PUT /api/items/update/{id}", middlewares.AuthMiddleware(cfg)(itemHandler.UpdateItem))
	router.HandleFunc("PUT /api/items/delete/{id}", middlewares.AuthMiddleware(cfg)(itemHandler.DeleteItem))

	router.HandleFunc("GET /api/settings", settingHandler.GetAllSettings)
	router.HandleFunc("GET /api/settings/single/{id}", settingHandler.GetSetting)
	router.HandleFunc("POST /api/settings/create", settingHandler.CreateSetting)
	router.HandleFunc("PUT /api/settings/update/{id}", settingHandler.UpdateSetting)
	router.HandleFunc("PUT /api/settings/delete/{id}", settingHandler.DeleteSetting)

	router.HandleFunc("GET /api/regions", regionHandler.GetAllRegions)

	router.HandleFunc("POST /api/login", authHandler.Login)

	router.HandleFunc("GET /api/users/me", middlewares.AuthMiddleware(cfg)(userHandler.GetUserInfo))
	router.HandleFunc("POST /api/users", userHandler.CreateUser)

	router.HandleFunc("GET /api/category-types", middlewares.AuthMiddleware(cfg)(categoryTypeHandler.GetCategoryTypes))
	router.HandleFunc("GET /api/category-types/single/{id}/", middlewares.AuthMiddleware(cfg)(categoryTypeHandler.GetCategoryType))
	router.HandleFunc("POST /api/category-types/create", middlewares.AuthMiddleware(cfg)(categoryTypeHandler.CreateCategoryType))
	router.HandleFunc("PUT /api/category-types/update/{id}", middlewares.AuthMiddleware(cfg)(categoryTypeHandler.UpdateCategoryType))
	router.HandleFunc("PUT /api/category-types/delete/{id}", middlewares.AuthMiddleware(cfg)(categoryTypeHandler.DeleteCategoryType))

	router.HandleFunc("GET /api/categories", middlewares.AuthMiddleware(cfg)(categoryHandler.GetCategories))
	router.HandleFunc("GET /api/categories/single/{id}", categoryHandler.GetCategory)
	router.HandleFunc("POST /api/categories/create", middlewares.AuthMiddleware(cfg)(categoryHandler.CreateCategory))
	router.HandleFunc("PUT /api/categories/update/{id}", middlewares.AuthMiddleware(cfg)(categoryHandler.UpdateCategory))
	router.HandleFunc("PUT /api/categories/delete/{id}", middlewares.AuthMiddleware(cfg)(categoryHandler.DeleteCategory))

	router.HandleFunc("POST /api/people/create", middlewares.AuthMiddleware(cfg)(peopleHandler.CreatePeople))
	router.HandleFunc("GET /api/people", middlewares.AuthMiddleware(cfg)(peopleHandler.GetAllPeople))
	router.HandleFunc("GET /api/people/single/{id}", middlewares.AuthMiddleware(cfg)(peopleHandler.GetSinglePeople))
	router.HandleFunc("PUT /api/people/update/{id}", middlewares.AuthMiddleware(cfg)(peopleHandler.UpdatePeople))
	router.HandleFunc("PUT /api/people/delete/{id}", middlewares.AuthMiddleware(cfg)(peopleHandler.DeletePeople))




	router.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "Hello, API!") })
	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			http.Error(w, "Database not healthy", http.StatusServiceUnavailable)
			return
		}
		fmt.Fprintf(w, "OK")
	})

	// static files
    router.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	handlers := middlewares.CORS(middlewares.Logging(middlewares.Recovery(router)))

	port := ":8081"
	fmt.Printf("Server starts on port %s\n", port)

	if err := http.ListenAndServe(port, handlers); err != nil {
		log.Fatal(err)
	}
}
