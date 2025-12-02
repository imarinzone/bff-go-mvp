# EV Charging API Contract Documentation

## 1) API: Search Charging Connectors

**Method:**
POST

**Endpoints**
`/v1/search`

**Description**
Search for available EV charging connectors and offers near a location or by specific filters. User can search either by EVSE ID or by geographic coordinates with distance radius.

**Headers:**

| **Header Name** | **Data Type** | **Mandatory** | **Description** |
|-----------------|---------------|---------------|-----------------|
| Authorization | String | Yes | Bearer token for authentication |
| X-Transaction-Id | String | No | Transaction identifier for request tracking (auto-generated if not provided) |

**Query Parameters:**

| **Name** | **Data Type** | **Mandatory** | **Description** |
|----------|---------------|---------------|-----------------|
| page | Integer | No | Page number for pagination (default: 1, min: 1) |
| per_page | Integer | No | Page size for pagination (default: 20, min: 1, max: 100) |

**Request Payload:**

| **Name** | **Data Type** | **Mandatory** | **Description** |
|----------|---------------|---------------|-----------------|
| evse_id | String | Conditional* | Specific EVSE identifier to look up |
| geo_coordinates | Array[Number] | Conditional* | Geographic coordinates as [latitude, longitude] |
| distance_meters | Number | Conditional* | Radius of search in meters (0-50000) |
| time_window | Object | No | Desired time window for availability |
| time_window.start | String (DateTime) | No | Start date and time |
| time_window.end | String (DateTime) | No | End date and time |
| filters | Object | No | Additional search filters |
| filters.cpo | String | No | Charge Point Operator identifier |
| filters.station | String | No | Charging station identifier |
| filters.connector_type | String | No | Type of connector (TYPE_1, TYPE_2, CCS2) |
| filters.vehicle | Object | No | Vehicle information |
| filters.vehicle.make | String | No | Vehicle manufacturer |
| filters.vehicle.model | String | No | Vehicle model |
| filters.vehicle.type | String | No | Vehicle type (2-wheeler, 3-wheeler, 4-wheeler) |

*Note: Either evse_id OR (geo_coordinates + distance_meters) is required

**Example Request JSON:**

```json
{
  "geo_coordinates": [12.9716, 77.5946],
  "distance_meters": 5000,
  "time_window": {
    "start": "2025-12-02T10:00:00Z",
    "end": "2025-12-02T18:00:00Z"
  },
  "filters": {
    "connector_type": "CCS2",
    "vehicle": {
      "type": "4-wheeler",
      "make": "Tesla",
      "model": "Model 3"
    }
  }
}
```

**Response Payload:**

| **Name** | **Data Type** | **Mandatory** | **Description** |
|----------|---------------|---------------|-----------------|
| total | Integer | Yes | Total number of catalogs matching search criteria |
| page | Integer | Yes | Current page number |
| per_page | Integer | Yes | Number of items per page |
| catalogs | Array[Object] | Yes | List of catalogs containing connectors and offers |
| catalogs[].id | String | Yes | Unique catalog identifier |
| catalogs[].connectors | Array[Object] | Yes | List of charging connectors |
| catalogs[].connectors[].id | String | Yes | Unique connector identifier |
| catalogs[].connectors[].geo_coordinates | Array[Number] | Yes | [latitude, longitude] |
| catalogs[].connectors[].availabilityWindow | Array[Object] | Yes | Time windows when connector is available |
| catalogs[].connectors[].rating | Object | Yes | User rating information |
| catalogs[].connectors[].rating.value | Number | Yes | Average rating value |
| catalogs[].connectors[].rating.count | Integer | Yes | Number of ratings |
| catalogs[].connectors[].isActive | Boolean | Yes | Whether connector is currently active |
| catalogs[].connectors[].provider | Object | Yes | Provider information |
| catalogs[].connectors[].connectorAttributes | Object | Yes | Technical attributes of the connector |
| catalogs[].connectors[].connectorAttributes.connectorType | String | Yes | Type (CCS2, Type2, Type1) |
| catalogs[].connectors[].connectorAttributes.maxPowerKW | Number | Yes | Maximum power in kilowatts |
| catalogs[].connectors[].connectorAttributes.status | String | Yes | Available, Occupied, Reserved, OutOfOrder |
| catalogs[].connectors[].connectorAttributes.chargingSpeed | String | Yes | SLOW, NORMAL, FAST, ULTRAFAST |
| catalogs[].offers | Array[Object] | Yes | List of pricing offers |

**Example Response JSON (Success):**

```json
{
  "total": 1,
  "page": 1,
  "per_page": 100,
  "catalogs": [
    {
      "id": "catalog-ev-charging-001",
      "connectors": [
        {
          "id": "ev-charger-ccs2-001",
          "geo_coordinates": [12.9716, 77.5946],
          "availabilityWindow": [
            {
              "startTime": "06:00:00",
              "endTime": "22:00:00"
            }
          ],
          "rating": {
            "value": 4.5,
            "count": 128
          },
          "isActive": true,
          "provider": {
            "id": "ecopower-charging",
            "descriptor": {
              "name": "EcoPower Charging Pvt Ltd"
            }
          },
          "connectorAttributes": {
            "connectorType": "CCS2",
            "maxPowerKW": 60,
            "minPowerKW": 5,
            "socketCount": 2,
            "reservationSupported": true,
            "status": "Available",
            "chargingSpeed": "FAST",
            "powerType": "DC",
            "connectorFormat": "CABLE"
          }
        }
      ],
      "offers": [
        {
          "id": "offer-ccs2-60kw-kwh",
          "descriptor": {
            "name": "Per-kWh Tariff - CCS2 60kW"
          },
          "items": ["ev-charger-ccs2-001"],
          "price": {
            "currency": "INR",
            "value": 18,
            "applicableQuantity": {
              "unitText": "Kilowatt Hour",
              "unitCode": "KWH",
              "unitQuantity": 1
            }
          },
          "acceptedPaymentMethod": ["UPI", "Card", "Wallet"]
        }
      ]
    }
  ]
}
```

---

## 2) API: Get Charging Estimate

**Method:**
POST

**Endpoints**
`/v1/estimate`

**Description**
Get cost and time estimates for a prospective charging session based on vehicle, connector, and charging preferences.

**Headers:**

| **Header Name** | **Data Type** | **Mandatory** | **Description** |
|-----------------|---------------|---------------|-----------------|
| Authorization | String | Yes | Bearer token for authentication |
| X-Transaction-Id | String | Yes | Transaction identifier for request tracking |
| X-Bpp-Id | String | Yes | Backend provider platform identifier |

**Request Payload:**

| **Name** | **Data Type** | **Mandatory** | **Description** |
|----------|---------------|---------------|-----------------|
| evse_id | String | Yes | EVSE identifier of the selected connector/station |
| connector_id | String | Yes | Specific connector identifier |
| vehicle | Object | Yes | Vehicle information |
| vehicle.make | String | Yes | Vehicle manufacturer |
| vehicle.model | String | Yes | Vehicle model |
| vehicle.type | String | Yes | Vehicle type (2-wheeler, 3-wheeler, 4-wheeler) |
| time_window | Object | No | Desired charging time window |
| time_window.start | String (DateTime) | No | Start date and time |
| time_window.end | String (DateTime) | No | End date and time |
| energy | Object | No | Target energy for estimation |
| energy.value | Number | No | Energy value |
| energy.unit | String | No | Unit of energy (e.g., kWh) |
| amount | Object | No | Amount for estimation |
| amount.value | Number | No | Amount value (minimum: 0) |
| amount.currency | String | No | Currency code (ISO 4217) |
| offer_id | String | No | Specific offer to apply |

**Example Request JSON:**

```json
{
  "evse_id": "IN*ECO*BTM*01*CCS2*A",
  "connector_id": "ev-charger-ccs2-001",
  "vehicle": {
    "make": "Tesla",
    "model": "Model 3",
    "type": "4-wheeler"
  },
  "time_window": {
    "start": "2025-12-02T14:00:00Z",
    "end": "2025-12-02T15:30:00Z"
  },
  "energy": {
    "value": 30,
    "unit": "kWh"
  },
  "amount": {
    "value": 200,
    "currency": "INR"
  },
  "offer_id": "offer-ccs2-60kw-kwh"
}
```

**Response Payload:**

| **Name** | **Data Type** | **Mandatory** | **Description** |
|----------|---------------|---------------|-----------------|
| order | Object | Yes | Order information |
| order.id | String | Yes | Order identifier |
| order.mode | String | Yes | Order mode (reservation, instant) |
| order.status | String | Yes | Order status (quoted_price) |
| amount | Object | Yes | Total estimated amount |
| amount.value | Number | Yes | Amount value |
| amount.currency | String | Yes | Currency code (ISO 4217) |
| durationInMinutes | String | Yes | Estimated duration in minutes |
| percentageOfBatteryCharged | String | Yes | Estimated battery percentage gain |
| energy | Object | Yes | Estimated energy |
| energy.value | Number | Yes | Energy value |
| energy.unit | String | Yes | Unit of energy (e.g., kWh) |
| validity | Object | Yes | Estimate validity period |
| validity.startDate | String (DateTime) | Yes | Start date and time |
| validity.endDate | String (DateTime) | Yes | End date and time |
| priceComponents | Array[Object] | Yes | Breakdown of price components |
| priceComponents[].type | String | Yes | Component type (UNIT, SURCHARGE, DISCOUNT, FEE) |
| priceComponents[].value | Number | Yes | Component value |
| priceComponents[].currency | String | Yes | Currency code |
| priceComponents[].description | String | Yes | Human-readable description |
| cancellation | Object | Yes | Cancellation policy and fees |
| cancellation.fee | Object | Yes | Cancellation fee information |
| cancellation.fee.percentage | String | Yes | Cancellation fee percentage |

**Example Response JSON (Success):**

```json
{
  "order": {
    "id": "1231208-id",
    "mode": "reservation",
    "status": "quoted_price"
  },
  "amount": {
    "value": 128.64,
    "currency": "INR"
  },
  "durationInMinutes": "15",
  "percentageOfBatteryCharged": "80",
  "energy": {
    "value": 30,
    "unit": "kWh"
  },
  "validity": {
    "startDate": "2025-12-02T14:00:00Z",
    "endDate": "2025-12-02T18:00:00Z"
  },
  "priceComponents": [
    {
      "type": "UNIT",
      "value": 100,
      "currency": "INR",
      "description": "Base charging session cost (100 INR)"
    },
    {
      "type": "SURCHARGE",
      "value": 20,
      "currency": "INR",
      "description": "Surge price (20%)"
    },
    {
      "type": "DISCOUNT",
      "value": -15,
      "currency": "INR",
      "description": "Offer discount (15%)"
    },
    {
      "type": "FEE",
      "value": 10,
      "currency": "INR",
      "description": "Service fee"
    }
  ],
  "cancellation": {
    "fee": {
      "percentage": "30"
    },
    "externalRef": {
      "mimetype": "text/html",
      "url": "https://example-company.com/charge/tnc.html"
    }
  }
}
```

---

## 3) API: Initiate Payment

**Method:**
POST

**Endpoints**
`/v1/orders/{order_id}/payment`

**Description**
Initiate or fetch payment instructions for the specified order after getting estimate.

**Headers:**

| **Header Name** | **Data Type** | **Mandatory** | **Description** |
|-----------------|---------------|---------------|-----------------|
| Authorization | String | Yes | Bearer token for authentication |
| X-Transaction-Id | String | Yes | Transaction identifier for request tracking |
| X-Bpp-Id | String | Yes | Backend provider platform identifier |

**Path Parameters:**

| **Name** | **Data Type** | **Mandatory** | **Description** |
|----------|---------------|---------------|-----------------|
| order_id | String | Yes | Unique order identifier from estimate response |

**Request Payload:**

| **Name** | **Data Type** | **Mandatory** | **Description** |
|----------|---------------|---------------|-----------------|
| payment_method | String | No | Preferred payment method (UPI, Card, Wallet, BankTransfer) |

**Response Payload:**

| **Name** | **Data Type** | **Mandatory** | **Description** |
|----------|---------------|---------------|-----------------|
| order | Object | Yes | Order information |
| order.id | String | Yes | Order identifier |
| order.status | String | Yes | Order status (ACTIVE) |
| order.mode | String | Yes | Order mode (RESERVATION) |
| currency | String | Yes | Currency code |
| amount | Number | Yes | Total payment amount |
| benificiary_id | String | No | Beneficiary identifier |
| acceptedPaymentMethod | Array[String] | Yes | List of accepted payment methods |
| payment_url | String | No | URL for payment gateway (if applicable) |
| validity | Object | Yes | Payment validity period |

**Example Response JSON (Success):**

```json
{
  "order": {
    "id": "order-bpp-789012",
    "status": "ACTIVE",
    "mode": "RESERVATION"
  },
  "currency": "INR",
  "amount": 128.64,
  "benificiary_id": "",
  "acceptedPaymentMethod": ["BankTransfer", "UPI", "Wallet"],
  "payment_url": "",
  "validity": {
    "startDate": "2025-12-02T14:00:00Z",
    "endDate": "2025-12-02T18:00:00Z"
  }
}
```

---

## 4) API: Get Order Information

**Method:**
GET

**Endpoints**
`/v1/orders/{order_id}`

**Description**
Retrieve order details, current status, and real-time telemetry for the specified order.

**Headers:**

| **Header Name** | **Data Type** | **Mandatory** | **Description** |
|-----------------|---------------|---------------|-----------------|
| Authorization | String | Yes | Bearer token for authentication |
| X-Transaction-Id | String | Yes | Transaction identifier for request tracking |
| X-Bpp-Id | String | Yes | Backend provider platform identifier |

**Path Parameters:**

| **Name** | **Data Type** | **Mandatory** | **Description** |
|----------|---------------|---------------|-----------------|
| order_id | String | Yes | Unique order identifier |

**Request Payload:**
None (GET request)

**Response Payload:**

| **Name** | **Data Type** | **Mandatory** | **Description** |
|----------|---------------|---------------|-----------------|
| order | Object | Yes | Order information |
| order.id | String | Yes | Order identifier |
| order.status | String | Yes | Order status (ACTIVE, COMPLETED, CANCELLED) |
| order.mode | String | Yes | Order mode (RESERVATION) |
| payment | Object | Yes | Payment information |
| payment.status | String | Yes | Payment status (PAID, PENDING) |
| charging | Object | Yes | Charging session information |
| charging.status | String | Yes | Charging status (ACTIVE, COMPLETED, IDLE) |
| connectorId | String | No | Connector identifier |
| connectorType | String | No | Type of connector used |
| vehicle | Object | No | Vehicle information |
| trackingUrl | String | No | URL to track charging session |
| chargingTelemetry | Object | No | Real-time charging telemetry data |
| chargingTelemetry.eventTime | String (DateTime) | Yes | Timestamp of telemetry event |
| chargingTelemetry.metrics | Array[Object] | Yes | Array of charging metrics |
| chargingTelemetry.metrics[].name | String | Yes | Metric name (STATE_OF_CHARGE, POWER, ENERGY, VOLTAGE, CURRENT, SESSION_DURATION) |
| chargingTelemetry.metrics[].value | Number | Yes | Metric value |
| chargingTelemetry.metrics[].unitCode | String | Yes | Unit code (PERCENTAGE, KWH, KW, VLT, AMP, min) |

**Example Response JSON (Success):**

```json
{
  "order": {
    "id": "order-bpp-789012",
    "status": "ACTIVE",
    "mode": "RESERVATION"
  },
  "payment": {
    "status": "PAID"
  },
  "charging": {
    "status": "ACTIVE"
  },
  "connectorId": "ev-charger-ccs2-001",
  "connectorType": "CCS2",
  "vehicle": {
    "type": "4-wheeler",
    "model": "Model 3",
    "make": "Tesla"
  },
  "trackingUrl": "https://track.bluechargenet-aggregator.io/session/SESSION-9876543210",
  "chargingTelemetry": {
    "eventTime": "2025-12-02T14:30:00Z",
    "metrics": [
      {
        "name": "STATE_OF_CHARGE",
        "value": 62.5,
        "unitCode": "PERCENTAGE"
      },
      {
        "name": "POWER",
        "value": 18.4,
        "unitCode": "KW"
      },
      {
        "name": "ENERGY",
        "value": 10.2,
        "unitCode": "KWH"
      },
      {
        "name": "VOLTAGE",
        "value": 392,
        "unitCode": "VLT"
      },
      {
        "name": "CURRENT",
        "value": 47,
        "unitCode": "AMP"
      },
      {
        "name": "SESSION_DURATION",
        "value": 10,
        "unitCode": "min"
      }
    ]
  }
}
```

---

## 5) API: Start Charging Session

**Method:**
PUT

**Endpoints**
`/v1/orders/{order_id}/start`

**Description**
Start a charging session for an existing order.

**Headers:**

| **Header Name** | **Data Type** | **Mandatory** | **Description** |
|-----------------|---------------|---------------|-----------------|
| Authorization | String | Yes | Bearer token for authentication |
| X-Transaction-Id | String | Yes | Transaction identifier for request tracking |
| X-Bpp-Id | String | Yes | Backend provider platform identifier |

**Path Parameters:**

| **Name** | **Data Type** | **Mandatory** | **Description** |
|----------|---------------|---------------|-----------------|
| order_id | String | Yes | Unique order identifier |

**Request Payload:**
None

**Response Payload:**

| **Name** | **Data Type** | **Mandatory** | **Description** |
|----------|---------------|---------------|-----------------|
| order | Object | Yes | Order information |
| order.id | String | Yes | Order identifier |
| order.status | String | Yes | Order status (ACTIVE) |
| order.mode | String | Yes | Order mode (RESERVATION) |
| payment | Object | Yes | Payment information |
| payment.status | String | Yes | Payment status (PAID) |
| charging | Object | Yes | Charging session information |
| charging.status | String | Yes | Charging status (ACTIVE) |

**Example Response JSON (Success):**

```json
{
  "order": {
    "id": "order-bpp-789012",
    "status": "ACTIVE",
    "mode": "RESERVATION"
  },
  "payment": {
    "status": "PAID"
  },
  "charging": {
    "status": "ACTIVE"
  }
}
```

---

## 6) API: Stop Charging Session

**Method:**
PUT

**Endpoints**
`/v1/orders/{order_id}/stop`

**Description**
Stop an active charging session for the specified order and receive completion summary with final charges.

**Headers:**

| **Header Name** | **Data Type** | **Mandatory** | **Description** |
|-----------------|---------------|---------------|-----------------|
| Authorization | String | Yes | Bearer token for authentication |
| X-Transaction-Id | String | Yes | Transaction identifier for request tracking |
| X-Bpp-Id | String | Yes | Backend provider platform identifier |

**Path Parameters:**

| **Name** | **Data Type** | **Mandatory** | **Description** |
|----------|---------------|---------------|-----------------|
| order_id | String | Yes | Unique order identifier |

**Request Payload:**

| **Name** | **Data Type** | **Mandatory** | **Description** |
|----------|---------------|---------------|-----------------|
| reason_code | String | No | Optional code indicating reason for stopping |
| message | String | No | Optional free-form message |

**Example Request JSON:**

```json
{
  "reason_code": "BATTERY_FULL",
  "message": "Battery reached desired level"
}
```

**Response Payload:**

| **Name** | **Data Type** | **Mandatory** | **Description** |
|----------|---------------|---------------|-----------------|
| order | Object | Yes | Order information |
| order.id | String | Yes | Order identifier |
| order.status | String | Yes | Order status (COMPLETED) |
| order.mode | String | Yes | Order mode (RESERVATION) |
| payment | Object | Yes | Payment information |
| payment.status | String | Yes | Payment status (PAID) |
| charging | Object | Yes | Charging session information |
| charging.status | String | Yes | Charging status (COMPLETED) |
| validity | Object | Yes | Validity period |
| price_components | Array[Object] | Yes | Breakdown of final price components |
| price_components[].type | String | Yes | Component type (BASE, SURCHARGE, DISCOUNT, FEE, REFUND) |
| price_components[].value | Number/String | Yes | Component value (negative for refunds) |
| price_components[].currency | String | Yes | Currency code |
| price_components[].description | String | Yes | Human-readable description |

**Example Response JSON (Success):**

```json
{
  "order": {
    "id": "order-bpp-789012",
    "status": "COMPLETED",
    "mode": "RESERVATION"
  },
  "payment": {
    "status": "PAID"
  },
  "charging": {
    "status": "COMPLETED"
  },
  "validity": {
    "startDate": "2025-12-02T14:00:00Z",
    "endDate": "2025-12-02T18:00:00Z"
  },
  "price_components": [
    {
      "type": "BASE",
      "value": 100,
      "currency": "INR",
      "description": "Base charging session cost (100 INR)"
    },
    {
      "type": "SURCHARGE",
      "value": 20,
      "currency": "INR",
      "description": "Surge price (20%)"
    },
    {
      "type": "DISCOUNT",
      "value": -15,
      "currency": "INR",
      "description": "Offer discount (15%)"
    },
    {
      "type": "FEE",
      "value": 10,
      "currency": "INR",
      "description": "Service fee"
    },
    {
      "type": "FEE",
      "value": 13.64,
      "currency": "INR",
      "description": "Overcharge estimation"
    }
  ]
}
```

---

## 7) API: Cancel Order - Estimate

**Method:**
GET

**Endpoints**
`/v1/orders/{order_id}/cancel`

**Description**
Get an estimate of cancellation charges and results for an existing order without performing the actual cancellation.

**Headers:**

| **Header Name** | **Data Type** | **Mandatory** | **Description** |
|-----------------|---------------|---------------|-----------------|
| Authorization | String | Yes | Bearer token for authentication |
| X-Transaction-Id | String | Yes | Transaction identifier for request tracking |
| X-Bpp-Id | String | Yes | Backend provider platform identifier |

**Path Parameters:**

| **Name** | **Data Type** | **Mandatory** | **Description** |
|----------|---------------|---------------|-----------------|
| order_id | String | Yes | Unique order identifier |

**Query Parameters:**

| **Name** | **Data Type** | **Mandatory** | **Description** |
|----------|---------------|---------------|-----------------|
| activity | String | No | Activity type related to cancellation |
| cancel_reason | String | No | Reason for cancellation |
| cancel_code | String | No | Standardized cancellation code |

**Request Payload:**
None (GET request)

**Response Payload:**

| **Name** | **Data Type** | **Mandatory** | **Description** |
|----------|---------------|---------------|-----------------|
| order | Object | Yes | Order information |
| order.id | String | Yes | Order identifier |
| order.status | String | Yes | Current order status (ACTIVE) |
| order.mode | String | Yes | Order mode (RESERVATION) |
| payment | Object | Yes | Payment information |
| payment.status | String | Yes | Payment status (PAID) |
| charging | Object | Yes | Charging session information |
| charging.status | String | Yes | Charging status (ACTIVE) |
| validity | Object | Yes | Validity period |
| price_components | Array[Object] | Yes | Breakdown of cancellation price components |
| price_components[].type | String | Yes | Component type (PAID, FEE) |
| price_components[].value | String | Yes | Component value |
| price_components[].currency | String | Yes | Currency code |
| price_components[].description | String | Yes | Human-readable description |

**Example Response JSON (Success):**

```json
{
  "order": {
    "id": "order-bpp-789012",
    "status": "ACTIVE",
    "mode": "RESERVATION"
  },
  "payment": {
    "status": "PAID"
  },
  "charging": {
    "status": "ACTIVE"
  },
  "validity": {
    "startDate": "2025-12-02T14:00:00Z",
    "endDate": "2025-12-02T18:00:00Z"
  },
  "price_components": [
    {
      "type": "PAID",
      "value": "400.00",
      "currency": "INR",
      "description": "Base price"
    },
    {
      "type": "FEE",
      "value": "30.00",
      "currency": "INR",
      "description": "Cancellation charges"
    }
  ]
}
```

---

## 8) API: Cancel Order - Confirm

**Method:**
POST

**Endpoints**
`/v1/orders/{order_id}/cancel`

**Description**
Cancel an existing order and process refund after deducting cancellation charges.

**Headers:**

| **Header Name** | **Data Type** | **Mandatory** | **Description** |
|-----------------|---------------|---------------|-----------------|
| Authorization | String | Yes | Bearer token for authentication |
| X-Transaction-Id | String | Yes | Transaction identifier for request tracking |
| X-Bpp-Id | String | Yes | Backend provider platform identifier |

**Path Parameters:**

| **Name** | **Data Type** | **Mandatory** | **Description** |
|----------|---------------|---------------|-----------------|
| order_id | String | Yes | Unique order identifier |

**Request Payload:**

| **Name** | **Data Type** | **Mandatory** | **Description** |
|----------|---------------|---------------|-----------------|
| cancel_reason | String | No | Reason for cancellation |
| cancel_code | String | No | Standardized cancellation code |

**Example Request JSON:**

```json
{
  "cancel_reason": "Change of plans",
  "cancel_code": "USER_INITIATED"
}
```

**Response Payload:**

| **Name** | **Data Type** | **Mandatory** | **Description** |
|----------|---------------|---------------|-----------------|
| order | Object | Yes | Order information |
| order.id | String | Yes | Order identifier |
| order.status | String | Yes | Order status (CANCELLED) |
| order.mode | String | Yes | Order mode (RESERVATION) |
| payment | Object | Yes | Payment information |
| payment.status | String | Yes | Payment status (PAID) |
| charging | Object | Yes | Charging session information |
| charging.status | String | Yes | Charging status (ACTIVE) |
| price_components | Array[Object] | Yes | Breakdown of cancellation price components |
| price_components[].type | String | Yes | Component type (FEE, REFUND) |
| price_components[].value | String | Yes | Component value (negative for refunds) |
| price_components[].currency | String | Yes | Currency code |
| price_components[].description | String | Yes | Human-readable description |

**Example Response JSON (Success):**

```json
{
  "order": {
    "id": "order-bpp-789012",
    "status": "CANCELLED",
    "mode": "RESERVATION"
  },
  "payment": {
    "status": "PAID"
  },
  "charging": {
    "status": "ACTIVE"
  },
  "price_components": [
    {
      "type": "FEE",
      "value": "30.00",
      "currency": "INR",
      "description": "Cancellation charges"
    },
    {
      "type": "REFUND",
      "value": "-300.00",
      "currency": "INR",
      "description": "Cancellation refund"
    }
  ]
}
```

---

## 9) API: Submit Rating & Feedback

**Method:**
POST

**Endpoints**
`/v1/orders/{order_id}/rating`

**Description**
Submit a rating and optional feedback for the specified order after charging session completion.

**Headers:**

| **Header Name** | **Data Type** | **Mandatory** | **Description** |
|-----------------|---------------|---------------|-----------------|
| Authorization | String | Yes | Bearer token for authentication |
| X-Transaction-Id | String | Yes | Transaction identifier for request tracking |
| X-Bpp-Id | String | Yes | Backend provider platform identifier |

**Path Parameters:**

| **Name** | **Data Type** | **Mandatory** | **Description** |
|----------|---------------|---------------|-----------------|
| order_id | String | Yes | Unique order identifier |

**Request Payload:**

| **Name** | **Data Type** | **Mandatory** | **Description** |
|----------|---------------|---------------|-----------------|
| value | Integer | Yes | Rating score between 1 and 5 |
| feedback | Object | No | Optional feedback information |
| feedback.comments | String | No | User comments about the experience |
| feedback.tags | Array[String] | No | Tags describing the experience (e.g., "fast", "clean", "helpful") |

**Example Request JSON:**

```json
{
  "value": 5,
  "feedback": {
    "comments": "Excellent charging experience. Fast and reliable service.",
    "tags": ["fast", "clean", "good_location", "helpful_staff"]
  }
}
```

**Response Payload:**

| **Name** | **Data Type** | **Mandatory** | **Description** |
|----------|---------------|---------------|-----------------|
| order | Object | Yes | Order information |
| order.id | String | Yes | Order identifier |
| order.status | String | Yes | Order status (COMPLETED) |
| order.mode | String | Yes | Order mode (RESERVATION) |
| feedbackForm | Object | No | Additional feedback form information |
| feedbackForm.url | String | No | URL to detailed feedback form |
| feedbackForm.mime_type | String | No | MIME type of feedback form |
| feedbackForm.submission_id | String | No | Unique submission identifier |

**Example Response JSON (Success):**

```json
{
  "order": {
    "id": "order-bpp-789012",
    "status": "COMPLETED",
    "mode": "RESERVATION"
  },
  "feedbackForm": {
    "url": "https://example-bpp.com/feedback/portal",
    "mime_type": "application/xml",
    "submission_id": "feedback-123e4567-e89b-12d3-a456-426614174000"
  }
}
```

---

## 10) API: Get Order Support

**Method:**
GET

**Endpoints**
`/v1/orders/{order_id}/support`

**Description**
Retrieve support contact information and available support channels for the specified order.

**Headers:**

| **Header Name** | **Data Type** | **Mandatory** | **Description** |
|-----------------|---------------|---------------|-----------------|
| Authorization | String | Yes | Bearer token for authentication |
| X-Transaction-Id | String | Yes | Transaction identifier for request tracking |
| X-Bpp-Id | String | Yes | Backend provider platform identifier |

**Path Parameters:**

| **Name** | **Data Type** | **Mandatory** | **Description** |
|----------|---------------|---------------|-----------------|
| order_id | String | Yes | Unique order identifier |

**Request Payload:**
None (GET request)

**Response Payload:**

| **Name** | **Data Type** | **Mandatory** | **Description** |
|----------|---------------|---------------|-----------------|
| order | Object | Yes | Order information |
| order.id | String | Yes | Order identifier |
| order.status | String | Yes | Order status |
| order.mode | String | Yes | Order mode |
| name | String | Yes | Support team name |
| phone | String | Yes | Support phone number |
| email | String | Yes | Support email address |
| url | String | No | URL to support portal or ticket |
| hours | String | Yes | Support hours of operation |
| channels | Array[String] | Yes | Available support channels (phone, email, web, chat) |

**Example Response JSON (Success):**

```json
{
  "order": {
    "id": "order-bpp-789012",
    "status": "CANCELLED",
    "mode": "RESERVATION"
  },
  "name": "BlueCharge Support Team",
  "phone": "18001080",
  "email": "support@bluechargenet-aggregator.io",
  "url": "https://support.bluechargenet-aggregator.io/ticket/SUP-20251202-001",
  "hours": "Mon–Sun 24/7 IST",
  "channels": ["phone", "email", "web", "chat"]
}
```

---

## Error Responses

All APIs follow a consistent error response format:

**Error Response Structure:**

| **Name** | **Data Type** | **Description** |
|----------|---------------|-----------------|
| error | Object | Error container |
| error.code | String | Error code (e.g., BAD_REQUEST, UNAUTHORIZED, NOT_FOUND) |
| error.message | String | Human-readable error message |
| error.details | Object | Additional error details (optional) |

**Common HTTP Status Codes:**

| **Status Code** | **Description** | **When It Occurs** |
|-----------------|-----------------|-------------------|
| 400 | Bad Request | Invalid request parameters or malformed request body |
| 401 | Unauthorized | Invalid or missing authentication token |
| 404 | Not Found | Order, EVSE, or connector not found |
| 409 | Conflict | Order in invalid state for requested operation |
| 422 | Unprocessable Entity | Validation error in request data |
| 500 | Internal Server Error | Server error occurred while processing request |

**Example Error Response:**

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Order not found.",
    "details": {
      "order_id": "order-bpp-789012",
      "timestamp": "2025-12-02T14:30:00Z"
    }
  }
}
```

---

## Authentication

All API endpoints require Bearer token authentication.

**Authorization Header Format:**
```
Authorization: Bearer <your_access_token>
```

**Header Requirements:**
- **Authorization**: Required for all endpoints
- **X-Transaction-Id**: Required for tracking requests across systems
- **X-Bpp-Id**: Required for multi-provider environments to identify backend provider

---

## API Flow Diagram

```
1. Search Connectors (POST /v1/search)
   ↓
2. Get Estimate (POST /v1/estimate)
   ↓
3. Initiate Payment (POST /v1/orders/{order_id}/payment)
   ↓
4. Get Order Status (GET /v1/orders/{order_id})
   ↓
5. Start Charging (PUT /v1/orders/{order_id}/start)
   ↓
6. Monitor Status (GET /v1/orders/{order_id}) [Polling/Real-time]
   ↓
7. Stop Charging (PUT /v1/orders/{order_id}/stop)
   ↓
8. Submit Rating (POST /v1/orders/{order_id}/rating)

Alternative Flows:
- Cancel Before Charging: After Step 3 → Cancel Order (POST /v1/orders/{order_id}/cancel)
- Get Support: Any time → Get Support (GET /v1/orders/{order_id}/support)
```

---

## Appendix: Enum Values

### Connector Types
- `TYPE_1`: Type 1 AC connector
- `TYPE_2`: Type 2 AC connector (European standard)
- `CCS2`: Combined Charging System 2 (DC fast charging)

### Vehicle Types
- `2-wheeler`: Two-wheeled electric vehicles
- `3-wheeler`: Three-wheeled electric vehicles
- `4-wheeler`: Four-wheeled electric vehicles (cars)

### Charging Speeds
- `SLOW`: 3-7 kW (several hours)
- `NORMAL`: 7-22 kW (2-4 hours)
- `FAST`: 50-60 kW (30-60 minutes)
- `ULTRAFAST`: 120+ kW (15-30 minutes)

### Connector Status
- `Available`: Connector is available for use
- `Occupied`: Connector is currently in use
- `Reserved`: Connector is reserved
- `OutOfOrder`: Connector is not functioning
- `Unknown`: Status cannot be determined

### Order Status
- `quoted_price`: Price estimate generated
- `ACTIVE`: Order is active and can proceed
- `COMPLETED`: Charging session completed
- `CANCELLED`: Order has been cancelled

### Payment Status
- `PAID`: Payment completed successfully
- `PENDING`: Payment is pending
- `FAILED`: Payment failed
- `REFUNDED`: Payment refunded

### Charging Status
- `ACTIVE`: Charging is in progress
- `COMPLETED`: Charging session completed
- `IDLE`: Vehicle connected but not charging
- `STOPPED`: Charging stopped by user

### Power Types
- `DC`: Direct current (fast charging)
- `AC`: Alternating current (standard charging)
- `AC_3_PHASE`: Three-phase alternating current

### Payment Methods
- `UPI`: Unified Payments Interface
- `Card`: Credit/Debit cards
- `Wallet`: Digital wallets
- `BankTransfer`: Direct bank transfer
- `NetBanking`: Online banking

---

**Document End**