#include <Wire.h>

const int rpmPin = 2;      // Pin for RPM sensor
const int leannessPin = A0; // Pin for leanness sensor (analog)
const int throttlePin = A1; // Pin for throttle sensor (analog)

void setup() {
  Serial.begin(9600);
  pinMode(rpmPin, INPUT);
}

void loop() {
  // Read RPM (example calculation, adjust as per sensor)
  int rpm = pulseIn(rpmPin, HIGH) * 60;

  // Read leanness and throttle (analog sensors)
  int leanness = analogRead(leannessPin);
  int throttle = analogRead(throttlePin);

  // Send data over Serial
  Serial.print("RPM:");
  Serial.print(rpm);
  Serial.print(",Leanness:");
  Serial.print(leanness);
  Serial.print(",Throttle:");
  Serial.println(throttle);

  delay(1000); // Delay 1 second between readings
}
