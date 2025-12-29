# Book of Shadows API Documentation

This document describes the REST API endpoints for the Book of Shadows application.

## Base URL

All API endpoints are relative to the application root: `/api/`

## Content Types

- Requests with a body should use `Content-Type: application/json`
- Responses may be HTML (for HTMX integration) or JSON depending on the endpoint

## Error Response Format

All errors follow this JSON structure:

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable message"
  }
}
```

### Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `NOT_FOUND` | 404 | Resource not found |
| `BAD_REQUEST` | 400 | Invalid request format |
| `INVALID_JSON` | 400 | Malformed JSON in request body |
| `MISSING_FIELD` | 400 | Required field is missing |
| `VALIDATION_ERROR` | 400 | Data validation failed |
| `CONFLICT` | 409 | Resource already exists |
| `PAYLOAD_TOO_LARGE` | 413 | Request body exceeds maximum size |
| `INTERNAL_ERROR` | 500 | Server error |

---

## Endpoints

### Investigators

#### List Investigators
```
GET /api/investigator
```

Returns HTML with all investigators for the current session.

**Response:** HTML (investigator list template)

---

#### Get Investigator
```
GET /api/investigator/{id}
```

Returns HTML for a single investigator's character sheet.

**Path Parameters:**
| Name | Type | Description |
|------|------|-------------|
| id | string | Investigator ID |

**Response:** HTML (character sheet template)

**Errors:**
- `404 NOT_FOUND` - Investigator not found

---

#### Create Investigator
```
POST /api/investigator/
```

Creates a new investigator with the provided base information.

**Request Body:**
```json
{
  "name": "John Smith",
  "age": "30",
  "residence": "Boston",
  "birthplace": "New York",
  "archetype": "Adventurer",
  "occupation": "Antiquarian"
}
```

**Response:**
```json
{
  "Key": "investigator-uuid"
}
```

**Errors:**
- `400 BAD_REQUEST` - Invalid request body

---

#### Update Investigator
```
PUT /api/investigator/{id}
```

Updates a field on an existing investigator.

**Path Parameters:**
| Name | Type | Description |
|------|------|-------------|
| id | string | Investigator ID |

**Request Body:**
```json
{
  "section": "personalInfo|attributes|skills|stats|status|weapons|gear|backstory",
  "field": "FieldName",
  "value": "new value"
}
```

**Section Types:**
- `personalInfo` - Name, age, residence, birthplace
- `attributes` - STR, DEX, INT, CON, APP, POW, SIZ, EDU, LUCK
- `skills` - Any skill name
- `stats` - Current HP, Sanity, Magic Points
- `status` - Insane, temporary insane, major wound, unconscious
- `weapons` - Weapon entries
- `gear` - Equipment and possessions
- `backstory` - Personal details and history

**Response:** `200 OK`

**Errors:**
- `404 NOT_FOUND` - Investigator not found
- `400 BAD_REQUEST` - Invalid field or value

---

#### Delete Investigator
```
DELETE /api/investigator/{id}
```

Deletes an investigator.

**Path Parameters:**
| Name | Type | Description |
|------|------|-------------|
| id | string | Investigator ID |

**Response:** `200 OK`

**Headers:**
- `HX-Trigger: deleted` - For HTMX integration

**Errors:**
- `404 NOT_FOUND` - Investigator not found

---

### Export/Import

#### Export Investigators
```
GET /api/investigator/list/export
```

Generates an export code containing all investigators.

**Response:**
```json
"export-code-uuid"
```

---

#### Import Investigators
```
POST /api/investigator/list/import/
```

Imports investigators from an export code.

**Request Body:**
```json
{
  "ImportCode": "export-code-uuid"
}
```

**Response:** `201 Created`

**Headers:**
- `HX-Trigger: import` - For HTMX integration

**Errors:**
- `400 BAD_REQUEST` - Missing or invalid import code

---

#### Export as PDF
```
POST /api/investigator/PDF/{id}
```

Generates a PDF export of the investigator's character sheet.

**Path Parameters:**
| Name | Type | Description |
|------|------|-------------|
| id | string | Investigator ID |

**Response:** `application/pdf` binary

**Errors:**
- `404 NOT_FOUND` - Investigator not found

---

### Random Generation

#### Generate Random Investigator
```
GET /api/generate/
```

Creates a random investigator.

**Query Parameters:**
| Name | Type | Default | Description |
|------|------|---------|-------------|
| mode | string | pulp | Game mode: `pulp` or `classic` |

**Response:** HTML (character sheet template)

---

### Archetypes

#### Get Archetype Occupations
```
GET /api/archetype/{name}/occupations/
```

Returns available occupations for a given archetype.

**Path Parameters:**
| Name | Type | Description |
|------|------|-------------|
| name | string | Archetype name |

**Response:**
```json
{
  "suggested": ["Occupation1", "Occupation2"],
  "others": ["Occupation3", "Occupation4"]
}
```

---

### Reporting

#### Report Issue
```
POST /api/report-issue
```

Submits a bug report or issue.

**Request Body:**
```json
{
  "issueType": "bug|feature|other",
  "description": "Issue description",
  "email": "optional@email.com",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Issue reported successfully"
}
```

---

## Wizard Endpoints

The wizard endpoints return HTML templates for the character creation flow.

### Base Step
```
GET /wizard/base/{key}
```

Returns the base information form (name, age, archetype, occupation).

### Attributes Step
```
GET /wizard/attributes/{key}
```

Returns the attribute allocation form.

### Skills Step
```
GET /wizard/skills/{key}
```

Returns the skill allocation form.

---

## HTMX Integration

Many endpoints are designed to work with HTMX for dynamic HTML updates:

- Responses include HTML fragments for target replacement
- Custom HX-Trigger headers signal events:
  - `deleted` - Investigator was deleted
  - `import` - Investigators were imported
- Use `hx-target` and `hx-swap` attributes for proper integration

## Request Limits

- Maximum request body size: 1MB
- POST/PUT requests require `Content-Type: application/json`
