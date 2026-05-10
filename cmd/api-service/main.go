package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
    "github.com/storskegg/ampvoice/internal/dbModels"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

var (
    mysqlUsername = ""
    mysqlPassword = ""
    mysqlDatabase = ""
    mysqlHost     = "mariadb_service"
    mysqlPort     = "3306"
)

func init() {
    var ok bool
    mysqlUsername, ok = os.LookupEnv("MYSQL_USER")
    if !ok {
        log.Fatal("MYSQL_USER environment variable not set")
    }
    mysqlPassword, ok = os.LookupEnv("MYSQL_PASSWORD")
    if !ok {
        log.Fatal("MYSQL_PASSWORD environment variable not set")
    }
    mysqlDatabase, ok = os.LookupEnv("MYSQL_DATABASE")
    if !ok {
        log.Fatal("MYSQL_DATABASE environment variable not set")
    }
}

func main() {
    log.Println("Starting server...")

    log.Println("Connecting to database...")
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlUsername, mysqlPassword, mysqlHost, mysqlPort, mysqlDatabase)
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %s", err.Error())
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
    err = gorm.G[dbModels.Parts](db).Create(ctx, dbModels.SamplePart())
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
