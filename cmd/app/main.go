package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"markitos-it-svc-golden/internal/domain/model"
	"markitos-it-svc-golden/internal/infrastructure/configuration"
	"markitos-it-svc-golden/internal/infrastructure/database"
	"markitos-it-svc-golden/internal/infrastructure/gapi"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var repository model.GoldenRepository
var config configuration.GoldenConfiguration

// #[.'.]:> Main function that orchestrates the startup and controlled shutdown of the application
// #[.'.]:> STEP 1: Show startup banner
// #[.'.]:> These logs help clearly identify the service startup
// #[.'.]:> STEP 2: Load configuration from files or environment variables
// #[.'.]:> This function sets all operational parameters
// #[.'.]:> STEP 3: Initialize database connection and repository
// #[.'.]:> Prepare data access and table structure
// #[.'.]:> STEP 4: Start gRPC servers
// #[.'.]:> Start entry points for gRPC clients
// #[.'.]:> STEP 5: Show shutdown banner when finished
// #[.'.]:> These logs clearly mark the end of the service execution
func main() {
	log.Println("['.']:>")
	log.Println("['.']:>--------------------------------------------")
	log.Println("['.']:>--- <starting markitos-it-svc-golden>  ---")

	loadConfiguration()
	log.Println("['.']:>------- configuration loaded")

	loadDatabase()
	log.Println("['.']:>------- database initialized")

	startServers()

	log.Println("['.']:>--------------------------------------------")
	log.Println("['.']:>--- <markitos-it-svc-golden stopped>  ---")
	log.Println("['.']:>")
}

// #[.'.]:> This function loads the service configuration
// #[.'.]:> STEP 1: Try to load configuration from file or environment variables
// #[.'.]:> Looks for "app.env" in the current directory, or uses environment variables if not found
// #[.'.]:> If there's an error, terminate the application immediately
// #[.'.]:> Can't operate without valid configuration
// #[.'.]:> STEP 2: Store configuration in a global variable
// #[.'.]:> Makes it accessible to the rest of the program functions
func loadConfiguration() {
	loadedConfig, err := configuration.LoadConfiguration(".")
	if err != nil {
		log.Fatal("['.']:>------- unable to load configuration: ", err)
	}

	config = loadedConfig
}

// #[.'.]:> This function initializes the database and repository
// #[.'.]:> STEP 1: Establish connection to PostgreSQL using the connection string
// #[.'.]:> GORM abstracts connection details and database handling
// #[.'.]:> If unable to connect to the database, it's a fatal error
// #[.'.]:> STEP 2: Run automatic migrations to create or update tables
// #[.'.]:> Ensures the database structure matches our models
// #[.'.]:> If migrations fail, can't continue
// #[.'.]:> STEP 3: Create a repository instance with the database connection
// #[.'.]:> The repository encapsulates all data access logic
func loadDatabase() {
	db, err := gorm.Open(postgres.Open(config.DatabaseDsn), &gorm.Config{})
	if err != nil {
		log.Fatal("['.']:> error unable to connect to database:", err)
	}

	err = db.AutoMigrate(&model.Golden{}, &model.GoldenEvent{})
	if err != nil {
		log.Fatal("['.']:> error unable to migrate database:", err)
	}

	repo := database.NewGoldenPostgresRepository(db)
	repository = &repo
}

// #[.'.]:> This function starts the servers and manages their lifecycle
// #[.'.]:> STEP 1: Create a cancelable context to signal shutdown
// #[.'.]:> This context will be propagated to the servers to manage their lifecycle
// #[.'.]:> STEP 2: Set up a channel to capture OS signals
// #[.'.]:> Allows detection of Ctrl+C or system shutdown signals
// #[.'.]:> STEP 3: Create a wait group to coordinate server shutdown
// #[.'.]:> The WaitGroup lets us wait for both servers to fully stop
// #[.'.]:> STEP 4: Start the gRPC server in a separate goroutine
// #[.'.]:> STEP 5: Block until a termination signal is received
// #[.'.]:> The application will wait here until Ctrl+C or SIGTERM is received
// #[.'.]:> STEP 6: Cancel the context to start controlled shutdown
// #[.'.]:> This sends the termination signal to both servers
// #[.'.]:> STEP 7: Wait for both servers to fully stop
// #[.'.]:> Won't exit until both servers have completed their shutdown
func startServers() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := runGRPCServer(ctx); err != nil {
			log.Printf("['.']:> error running gRPC server: %v", err)
		}
	}()
	<-stop
	log.Println("['.']:>------- shutting down servers gracefully...")
	cancel()
	wg.Wait()
}

// #[.'.]:> This function starts and manages the gRPC server lifecycle
// #[.'.]:> STEP 1: Create a network listener
// #[.'.]:> This listener will listen for TCP requests at the configured address and port
// #[.'.]:> STEP 2: Create a new gRPC server instance
// #[.'.]:> This object is the heart of the server and will handle all requests
// #[.'.]:> STEP 3: Create the implementation of our service
// #[.'.]:> This part contains the actual business logic
// #[.'.]:> STEP 4: Register our service with the gRPC server
// #[.'.]:> Connects our implementations with the gRPC system
// #[.'.]:> STEP 5: Enable reflection to facilitate testing
// #[.'.]:> Reflection allows tools like grpcurl to discover our services
// #[.'.]:> STEP 6: Set up controlled (graceful) shutdown
// #[.'.]:> This goroutine runs in the background and waits for the shutdown signal
// #[.'.]:> Blocks until the context is canceled (shutdown signal)
// #[.'.]:> Logs a message indicating the server is shutting down
// #[.'.]:> Performs a graceful shutdown:
// #[.'.]:> - Stops accepting new connections
// #[.'.]:> - Waits for ongoing calls to finish
// #[.'.]:> - Closes all connections cleanly
// #[.'.]:> STEP 7: Log that the server is running
// #[.'.]:> STEP 8: Start the server (this method blocks until an error occurs)
// #[.'.]:> The server now actively listens for incoming requests
func runGRPCServer(ctx context.Context) error {
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(gapi.UnaryLoggingInterceptor()))
	server := gapi.NewServer(config.GRPCServerAddress, repository, config)
	gapi.RegisterGoldenserviceServer(grpcServer, server)
	reflection.Register(grpcServer)

	go func() {
		<-ctx.Done()
		log.Println("['.']:> shutting down gRPC server...")
		grpcServer.GracefulStop()
	}()
	log.Printf("['.']:> gRPC server running at %s", config.GRPCServerAddress)

	return grpcServer.Serve(listener)
}
