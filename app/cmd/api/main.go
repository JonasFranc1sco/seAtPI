package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"app/app/cmd/cli"
	"app/app/internal/database"
	"app/app/internal/handlers"
	"app/app/internal/repositories"
	"app/app/internal/services"
)

func main() {
	db, err := database.Connect("data/inventory.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatal(err)
	}

	equipmentRepository := repositories.NewEquipmentRepository(db)
	equipmentService := services.NewEquipmentService(equipmentRepository)
	equipmentHandler := handlers.NewEquipmentHandler(equipmentService)

	antivirusRepository := repositories.NewAntivirusRepository(db)
	antivirusService := services.NewAntivirusService(antivirusRepository)
	antivirusHandler := handlers.NewAntivirusHandler(antivirusService)

	inventoryRepository := repositories.NewInventoryRepository(db)
	inventoryService := services.NewInventoryService(inventoryRepository)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)

	reportRepository := repositories.NewReportRepository(db)
	reportService := services.NewReportService(reportRepository)
	reportHandler := handlers.NewReportHandler(reportService)

	router := mux.NewRouter()

	// Registro de rotas
	router.HandleFunc("/health", handlers.HealthHandler).Methods(http.MethodGet)
	router.HandleFunc("/version", handlers.HealthHandler).Methods(http.MethodGet)
	router.HandleFunc("/info", handlers.InfoHandler).Methods(http.MethodGet)

	// Rotas de equipamentos
	router.HandleFunc("/equipamentos", equipmentHandler.ListEquipmentsHandler).Methods(http.MethodGet)
	router.HandleFunc("/equipamentos", equipmentHandler.CreateEquipmentHandler).Methods(http.MethodPost)
	router.HandleFunc("/equipamentos/{tombo}", equipmentHandler.UpdateEquipmentHandler).Methods(http.MethodPut)
	router.HandleFunc("/equipamentos/{tombo}", equipmentHandler.DeleteEquipmentHandler).Methods(http.MethodDelete)
	router.HandleFunc("/equipamentos/{tombo}", equipmentHandler.GetEquipmentByAssetTagHandler).Methods(http.MethodGet)
	router.HandleFunc("/equipamentos/{serial}", equipmentHandler.GetEquipmentBySerialNumberHandler).Methods(http.MethodGet)

	// Rotas para importar CSVs
	router.HandleFunc("/imports/patrimonio", equipmentHandler.ImportPatrimonyCSVHandler).Methods(http.MethodPost)
	router.HandleFunc("/imports/antivirus", antivirusHandler.ImportTrendCSVHandler).Methods(http.MethodPost)
	router.HandleFunc("/imports/inventario", inventoryHandler.ImportInventoryCSVHandler).Methods(http.MethodPost)

	router.HandleFunc("/reports/sem-antivirus", reportHandler.EquipmentsWithoutAntivirusHandler).Methods(http.MethodGet)
	router.HandleFunc("/reports/antivirus-com-problema", reportHandler.EquipmentWithAntivirusProblemsHandler).Methods(http.MethodGet)
	router.HandleFunc("/reports/juncao", reportHandler.InventoryEquipmentAntivirusUnifiedHandler).Methods(http.MethodGet)

	// Iniciando o servidor HTTP na porta 8080.
	go func() {
		err = http.ListenAndServe(":8080", router)
		if err != nil {
			log.Fatal(err)
		}
	}()

	cli.Interface()
}
