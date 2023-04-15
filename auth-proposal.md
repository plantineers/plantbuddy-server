# Authentication and Authorization (might be done later)

## For Sensors/Micro-Controllers

When a sensor is created, a certificate will be generated and passed to the sensor along with its ID. The certificate
will be used to authenticate the sensor.

## For Users

Users will be authenticated using OAuth2 along with Google Firebase. We implement a role model for different use cases.

Alternative: Basic Auth

**Roles**:

- **Chief Buddy**: Can create, update, delete sensors and plants and plant groups.
- **Buddy**: Can view all data.

*Thought: use user/pw for both sensors and users because it is easier to have just one auth method.*
