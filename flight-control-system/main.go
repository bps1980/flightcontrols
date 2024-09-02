package main

import (
	"bufio"
	"encoding/json"
	"flight-control-system/actuators"
	"flight-control-system/pid"
	"flight-control-system/sensors"
	"flight-control-system/shared"
	"flight-control-system/utils"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/tarm/serial"
)

var (
	upgrader     = websocket.Upgrader{}
	clients      = make(map[*websocket.Conn]bool)
	broadcast    = make(chan SensorData)
	airports     = make(map[string]Airport)
	transponders = make(map[string]Transponder)
)

type SensorData struct {
	RPM      int `json:"rpm"`
	Leanness int `json:"leanness"`
	Throttle int `json:"throttle"`
}

type Airport struct {
	Code           string  `json:"code"`
	Name           string  `json:"name"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	RunwayLighting bool    `json:"runway_lighting"`
	GroundFreq     string  `json:"ground_freq"`
	TowerFreq      string  `json:"tower_freq"`
}

type Transponder struct {
	ID        string  `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
	Squawk    string  `json:"squawk"`
}

func main() {
	// Initialize sensors and actuators
	if err := sensors.InitSensors(); err != nil {
		fmt.Printf("Failed to initialize sensors: %v\n", err)
		return
	}
	if err := actuators.InitActuators(); err != nil {
		fmt.Printf("Failed to initialize actuators: %v\n", err)
		return
	}

	// Create a PID controller for pitch control
	pidPitch := pid.NewPIDController(1.0, 0.1, 0.05)

	// Set up HTTP server and WebSocket handler
	router := mux.NewRouter()
	fs := http.FileServer(http.Dir("./public"))
	router.PathPrefix("/").Handler(fs)
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleConnections(w, r, pidPitch)
	})
	router.HandleFunc("/api/airports", getAirports).Methods("GET")
	router.HandleFunc("/api/airports/{code}/toggle_lighting", toggleRunwayLighting).Methods("POST")
	router.HandleFunc("/api/transponders", getTransponders).Methods("GET")
	router.HandleFunc("/api/transponders", updateTransponder).Methods("POST")
	router.HandleFunc("/api/communication", handleCommunication).Methods("POST")
	router.HandleFunc("/api/transponders/{id}/squawk", updateSquawk).Methods("POST")

	go handleSerialData()
	go handleMessages()

	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func handleConnections(w http.ResponseWriter, r *http.Request, pidPitch *pid.PIDController) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to upgrade to WebSocket: %v\n", err)
		return
	}
	defer conn.Close()

	clients[conn] = true

	currentMode := utils.SensorMode
	magnetoSwitches := []bool{false, false, false, false}
	ignitionSwitches := []bool{false, false, false, false}
	const rateLimit = 1.0
	var previousControlSignals shared.ControlSignals

	for {
		// Simulate changing flight modes (you can replace this with actual input handling)
		switch time.Now().Second() % 30 {
		case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9:
			currentMode = utils.SensorMode
		case 10, 11, 12, 13, 14, 15, 16, 17, 18, 19:
			currentMode = utils.PerformanceMode
		default:
			currentMode = utils.SuperPerformanceMode
		}

		// Simulate toggling magneto and ignition switches (replace with actual input handling)
		for i := range magnetoSwitches {
			magnetoSwitches[i] = !magnetoSwitches[i]
			ignitionSwitches[i] = !ignitionSwitches[i]
		}

		// Get the desired pitch setpoint based on the current mode
		setpointPitch := utils.GetPitchSetpoint(currentMode)
		leannessSetting := utils.GetLeannessSetting(currentMode)

		// Read sensor data
		sensorData := sensors.ReadSensorData()

		// Compute control output for pitch
		controlPitch := pidPitch.Update(setpointPitch, sensorData.Pitch)

		// Convert control output to throttle, propeller pitch, and leanness commands
		controlSignals := shared.ControlSignals{
			Throttle:       1000 + controlPitch,
			PropPitch:      10 + controlPitch,
			Leanness:       leannessSetting,
			MagnetoSwitch:  magnetoSwitches,
			IgnitionSwitch: ignitionSwitches,
		}

		// Apply rate limiting
		controlSignals.Throttle = rateLimitControl(previousControlSignals.Throttle, controlSignals.Throttle, rateLimit)
		controlSignals.PropPitch = rateLimitControl(previousControlSignals.PropPitch, controlSignals.PropPitch, rateLimit)
		controlSignals.Leanness = rateLimitControl(previousControlSignals.Leanness, controlSignals.Leanness, rateLimit)

		previousControlSignals = controlSignals

		// Update actuators
		actuators.UpdateActuators(controlSignals)

		// Prepare data for WebSocket transmission
		data := struct {
			Mode           utils.FlightMode      `json:"mode"`
			SensorData     sensors.SensorData    `json:"sensor_data"`
			ControlSignals shared.ControlSignals `json:"control_signals"`
		}{
			Mode:           currentMode,
			SensorData:     sensorData,
			ControlSignals: controlSignals,
		}

		// Send data to the client
		err := conn.WriteJSON(data)
		if err != nil {
			fmt.Printf("Failed to write JSON to WebSocket: %v\n", err)
			delete(clients, conn)
			break
		}

		// Wait for a short duration before the next loop iteration
		time.Sleep(100 * time.Millisecond)
	}
}

func handleSerialData() {
	var config *serial.Config

	// Adjust the serial port based on your OS and setup
	if isWindows() {
		config = &serial.Config{Name: "COM3", Baud: 9600}
	} else {
		config = &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600} // Adjust for your Linux/WSL setup
	}

	port, err := serial.OpenPort(config)
	if err != nil {
		fmt.Println("Error opening serial port:", err)
		return
	}
	defer port.Close()

	reader := bufio.NewReader(port)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from serial port:", err)
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) == 3 {
			rpm, leanness, throttle := parseSensorData(parts)
			sensorData := SensorData{
				RPM:      rpm,
				Leanness: leanness,
				Throttle: throttle,
			}
			broadcast <- sensorData
		}
	}
}

func isWindows() bool {
	return strings.Contains(runtime.GOOS, "windows")
}

func rateLimitControl(previous, current, rateLimit float64) float64 {
	if current > previous+rateLimit {
		return previous + rateLimit
	} else if current < previous-rateLimit {
		return previous - rateLimit
	}
	return current
}

func handleMessages() {
	for {
		data := <-broadcast
		message, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error marshaling data:", err)
			continue
		}
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				fmt.Println("Error sending message:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func parseSensorData(parts []string) (int, int, int) {
	// Parse and return sensor data from the parts
	// Adjust based on the actual format of your data
	var rpm, leanness, throttle int
	fmt.Sscanf(parts[0], "RPM:%d", &rpm)
	fmt.Sscanf(parts[1], "Leanness:%d", &leanness)
	fmt.Sscanf(parts[2], "Throttle:%d", &throttle)
	return rpm, leanness, throttle
}

func getAirports(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(airports)
}

func toggleRunwayLighting(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]
	airport, exists := airports[code]
	if !exists {
		http.Error(w, "Airport not found", http.StatusNotFound)
		return
	}
	airport.RunwayLighting = !airport.RunwayLighting
	airports[code] = airport
	json.NewEncoder(w).Encode(airport)
}

func getTransponders(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(transponders)
}

func updateTransponder(w http.ResponseWriter, r *http.Request) {
	var t Transponder
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	transponders[t.ID] = t
	json.NewEncoder(w).Encode(t)
}

func updateSquawk(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var req struct{ Squawk string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	transponder, exists := transponders[id]
	if !exists {
		http.Error(w, "Transponder not found", http.StatusNotFound)
		return
	}
	transponder.Squawk = req.Squawk
	transponders[id] = transponder
	json.NewEncoder(w).Encode(transponder)
}

func handleCommunication(w http.ResponseWriter, r *http.Request) {
	var msg struct {
		To      string `json:"to"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("Communication to %s: %s\n", msg.To, msg.Message)
	w.WriteHeader(http.StatusOK)
}
