# Data we need

## Possible Endpoints

- `/v1/sensor-gateway`: Populate data (time series DB)
- `/v1/sensors`: Get all sensors
- `/v1/sensor/{id}`: Get a single sensor

- `/v1/sensor-types`: Get all sensor types (we should hard code the types)
- `/v1/sensor-types/{id}`: Get a single sensor type

- `/v1/plants`: Get all plants
- `/v1/plant/{id}`: Get a single plant

- `/v1/plant-groups`: Get all plant groups
- `/v1/plant-group/{id}`: Get a single plant group

- `/v1/plant-ql`: Execute a query

## Sensor Gateway

```json
{
    "sensors": [
        {
            "id": 1,
            "timestamp": "2000-01-01T00:00:00.000000",
            "value": 50.0
        }
    ]
}
```

PlantQL: `Sensor.type="humidity" & Data.value > 50 & Plant.type="Cactaceae"`

Optional extension: send custom labels/tags that can also be filtered.

## Sensor Manifest

```json
{
  "id": 1,
  "plant": 1,
  "interval": "60",
  "type": "humidity"
}
```

## Plant Manifest

### Specific Plant

```json
{
  "id": 1,
  "group": 1,
  "additionalCareTips": [
    "Massage every morning"
  ]
}
```

### Plant Group

```json
{
  "id": 1,
  "type": "Cactaceae",
  "description": "The big one behind the door",
  "careTips": [
    "Play Mozart music every evening.",
    "Place in the sun"
  ],
  "ranges": [
    {
      "sensorType": "humidity",
      "min": 40,
      "max": 60
    },
    {
      "sensorType": "temperature",
      "min": 10,
      "max": 50
    }
  ]
}
```

## Authentication and Authorization (might be done later)

### For Sensors/Micro-Controllers

When a sensor is created, a certificate will be generated and passed to the sensor along with its ID. The certificate
will be used to authenticate the sensor.

### For Users

Users will be authenticated using OAuth2 along with Google Firebase. We implement a role model for different use cases.

Alternative: Basic Auth

**Roles**:

- **Chief Buddy**: Can create, update, delete sensors and plants and plant groups.
- **Buddy**: Can view all data.

*Thought: use user/pw for both sensors and users because it is easier to have just one auth method.*
