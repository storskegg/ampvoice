package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/storskegg/ampvoice/internal/dbModels"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func main() {
    log.Println("Starting server...")
    dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=America/Chicago"

    log.Println("Connecting to database...")
    db, err := gorm.Open(mysql.New(mysql.Config{
        DSN: dsn,
        //PreferSimpleProtocol: true,
    }), &gorm.Config{})

    if err != nil {
        log.Fatal(err)
    }
    log.Println("Database connection established")

    ctx := context.Background()

    log.Println("Migrating database...")
    // Migrate the schema
    err = db.AutoMigrate(&dbModels.Parts{})
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Database migrated")

    log.Println("Creating initial part data...")
    err = gorm.G[dbModels.Parts](db).Create(ctx, &dbModels.Parts{
        Category: "Passives",
        Type:     "Capacitor",
        Subtype:  "Electrolytic",
        Brand:    "Jupiter",
        Series:   "Cosmos",
        Value:    "100 uF",
        Rating:   "100 V",
        CostUnit: 9.49,
        CostMult: 1.75,
    })
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Initial part data created")

    log.Println("Starting server...")
    router := mux.NewRouter()
    router.HandleFunc("/api/go/part", getPart(ctx, db)).Methods("GET")

    // wrap the router with CORS and JSON content type middlewares
    enhancedRouter := enableCORS(jsonContentTypeMiddleware(router))
    // start server
    err = http.ListenAndServe(":8000", enhancedRouter)
    if err != nil {
        log.Fatal(err)
    }
}

func enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Set CORS headers
        w.Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        // Check if the request is for CORS preflight
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        // Pass down the request to the next middleware (or final handler)
        next.ServeHTTP(w, r)
    })
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Set JSON Content-Type
        w.Header().Set("Content-Type", "application/json")
        next.ServeHTTP(w, r)
    })
}

func getPart(ctx context.Context, db *gorm.DB) http.HandlerFunc {
    part, err := gorm.G[dbModels.Parts](db).First(ctx) // find product with integer primary key
    if err != nil {
        return func(w http.ResponseWriter, r *http.Request) {
            http.Error(w, "Product not found", http.StatusNotFound)
        }
    }
    return func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(part)
    }
}
