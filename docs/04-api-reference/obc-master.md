# OBC Master API Reference

## Overview

OBC Master API merupakan endpoint untuk mengelola data master OBC (Order Batch Confirmation) yang bersumber dari SAP. API ini mencakup fitur import Excel, list, detail, dan auto-generate Production Order.

## Authentication & Authorization

**Required Roles:**
- `ADMIN` - Full access
- `PPIC` - Full access
- `MANAGER` - Read-only access (list & detail)

**Authentication Method:** Bearer Token (JWT)

---

## Endpoints

### 1. Import OBC dari Excel

Import data OBC Master dari file Excel (.xlsx) dengan mapping otomatis ke 39 kolom database.

**Endpoint:** `POST /api/obc/import`

**Content-Type:** `multipart/form-data`

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| file | file | Yes | Excel file (.xlsx) dengan format SAP OBC |
| auto_generate_po | boolean | No | Auto-generate Production Orders (default: false) |

**Request Example:**

```bash
curl -X POST http://localhost:8080/api/obc/import \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@obc_data.xlsx" \
  -F "auto_generate_po=true"
```

**Response Success (200):**

```json
{
  "success": true,
  "message": "Import berhasil",
  "data": {
    "total_rows": 150,
    "success_count": 148,
    "failed_count": 2,
    "failed_rows": [
      {
        "row_number": 45,
        "obc_number": "OBC-2024-001",
        "error": "No OBC tidak boleh kosong"
      },
      {
        "row_number": 89,
        "obc_number": "OBC-2024-055",
        "error": "invalid date format: 32/13/2024"
      }
    ],
    "pos_generated": 320,
    "duration_ms": 2450,
    "file_name": "obc_data.xlsx",
    "file_size": 245760
  }
}
```

**Response Partial Success (206):**

```json
{
  "success": false,
  "message": "Import selesai dengan beberapa error",
  "data": {
    "total_rows": 100,
    "success_count": 95,
    "failed_count": 5,
    "failed_rows": [...],
    "pos_generated": 0,
    "duration_ms": 1850
  }
}
```

**Error Responses:**

- `400 Bad Request` - File tidak valid atau missing
- `401 Unauthorized` - Token tidak valid
- `403 Forbidden` - Role tidak memiliki akses
- `500 Internal Server Error` - Database atau parsing error

**Excel Column Mapping:**

| Excel Column | Database Field | Type | Required | Notes |
|--------------|----------------|------|----------|-------|
| No OBC | obc_number | VARCHAR(20) | Yes | Unique identifier |
| Tgl OBC | obc_date | DATE | No | Format: DD/MM/YYYY atau Excel date |
| Material | material | VARCHAR(50) | No | Indexed untuk search |
| SERI | seri | VARCHAR(50) | No | Indexed untuk search |
| WARNA | warna | VARCHAR(50) | No | Indexed untuk search |
| KODE_PABRIK | factory_code | VARCHAR(50) | No | Indexed untuk search |
| QTY PESAN | quantity_ordered | INT | No | Untuk kalkulasi PO |
| JHT | jht | VARCHAR(100) | No | - |
| RPB | rpb | DECIMAL(15,2) | No | - |
| HJE | hje | DECIMAL(15,2) | No | - |
| BPB | bpb | INT | No | - |
| RENCET | rencet | INT | No | - |
| Tgl JTempo | due_date | DATE | No | Format: DD/MM/YYYY atau Excel date |
| Perso / non Perso | personalization | VARCHAR(20) | No | "Perso" or "Non Perso" |
| Perekat | adhesive_type | VARCHAR(50) | No | - |
| GR | gr | VARCHAR(50) | No | - |
| No Pelat | plat_number | VARCHAR(50) | No | - |
| Type | type | VARCHAR(50) | No | - |
| Created On | created_on | DATE | No | Format: DD/MM/YYYY atau Excel date |
| Sales Doc. | sales_document | VARCHAR(50) | No | - |
| Item | item_code | VARCHAR(50) | No | - |
| Material description | material_description | VARCHAR(255) | No | - |
| BUn | base_unit | VARCHAR(20) | No | - |
| Gol. PCA | pca_category | VARCHAR(50) | No | - |
| Kadar Alkohol PCA | alcohol_percentage | DECIMAL(5,2) | No | - |
| Kadar HPTL | hptl_content | DECIMAL(5,2) | No | - |
| Kode Wilayah | region_code | VARCHAR(20) | No | - |
| OBC Awal | obc_initial | VARCHAR(50) | No | - |
| Peruntukan | allocation | VARCHAR(255) | No | - |
| PESANAN | total_order_obc | INT | No | - |
| Plnt | plant_code | VARCHAR(10) | No | - |
| SATUAN | unit | VARCHAR(20) | No | - |
| Tahun | production_year | INT | No | - |
| Tarif Per Liter | excise_rate_per_liter | DECIMAL(15,2) | No | - |
| Volume PCA | pca_volume | DECIMAL(15,2) | No | - |
| Warna MMEA | mmea_color_code | VARCHAR(50) | No | - |

**Import Behavior:**

- **Upsert Logic**: Jika OBC number sudah ada, data akan di-update. Jika baru, data akan di-create.
- **Transaction-based**: Semua operasi dalam satu transaction untuk data consistency.
- **Error Handling**: Row yang error akan di-skip, tidak mempengaruhi row lain.
- **Date Parsing**: Support Excel serial date dan text format (DD/MM/YYYY, DD-MM-YYYY, dll).
- **Number Parsing**: Handle thousands separator (comma/dot) dengan benar.

---

### 2. List OBC Masters

Mengambil list OBC Masters dengan pagination dan filtering.

**Endpoint:** `GET /api/obc`

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | int | 1 | Page number |
| page_size | int | 20 | Items per page (max: 100) |
| material | string | - | Filter by material (LIKE search) |
| seri | string | - | Filter by seri (LIKE search) |
| warna | string | - | Filter by warna (LIKE search) |
| factory_code | string | - | Filter by factory code (LIKE search) |
| obc_number | string | - | Filter by OBC number (LIKE search) |

**Request Example:**

```bash
curl -X GET "http://localhost:8080/api/obc?page=1&page_size=20&material=PITA&seri=2024" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response Success (200):**

```json
{
  "success": true,
  "message": "List OBC Master berhasil diambil",
  "data": {
    "items": [
      {
        "id": 1,
        "obc_number": "OBC-2024-001",
        "obc_date": "2024-01-15",
        "material": "PITA CUKAI 2024",
        "seri": "SER-2024-A",
        "warna": "MERAH",
        "factory_code": "F001",
        "quantity_ordered": 50000,
        "material_description": "Pita Cukai Hasil Tembakau Tahun 2024",
        "due_date": "2024-03-15",
        "personalization": "Non Perso",
        "created_on": "2024-01-10",
        "created_at": "2024-01-15 10:30:45"
      }
    ],
    "total": 148,
    "page": 1,
    "page_size": 20,
    "total_pages": 8
  }
}
```

---

### 3. Detail OBC Master

Mengambil detail OBC Master beserta Production Orders terkait.

**Endpoint:** `GET /api/obc/:id`

**Path Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | uint64 | Yes | OBC Master ID |

**Request Example:**

```bash
curl -X GET http://localhost:8080/api/obc/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response Success (200):**

```json
{
  "success": true,
  "message": "Detail OBC Master berhasil diambil",
  "data": {
    "id": 1,
    "obc_number": "OBC-2024-001",
    "obc_date": "2024-01-15",
    "material": "PITA CUKAI 2024",
    "seri": "SER-2024-A",
    "warna": "MERAH",
    "factory_code": "F001",
    "quantity_ordered": 50000,
    "jht": "JHT-001",
    "rpb": 15000.00,
    "hje": 20000.00,
    "bpb": 1000,
    "rencet": 50,
    "due_date": "2024-03-15",
    "personalization": "Non Perso",
    "adhesive_type": "Type A",
    "gr": "GR-001",
    "plat_number": "PLT-2024-001",
    "type": "Standard",
    "created_on": "2024-01-10",
    "sales_document": "SD-2024-001",
    "item_code": "ITEM-001",
    "material_description": "Pita Cukai Hasil Tembakau Tahun 2024",
    "base_unit": "PCS",
    "pca_category": "Category A",
    "alcohol_percentage": 0.00,
    "hptl_content": 0.00,
    "region_code": "REG-01",
    "obc_initial": "OBC-INIT-001",
    "allocation": "Domestic Market",
    "total_order_obc": 50000,
    "plant_code": "P001",
    "unit": "PCS",
    "production_year": 2024,
    "excise_rate_per_liter": 15000.00,
    "pca_volume": 0.00,
    "mmea_color_code": "RED-001",
    "created_at": "2024-01-15 10:30:45",
    "updated_at": "2024-01-15 10:30:45",
    "production_orders": [
      {
        "id": 101,
        "po_number": 1705305045001,
        "quantity_ordered": 40000,
        "priority": "NORMAL",
        "current_stage": "KHAZWAL_MATERIAL_PREP",
        "current_status": "WAITING_MATERIAL_PREP",
        "due_date": "2024-03-15",
        "created_at": "2024-01-15 10:30:45"
      },
      {
        "id": 102,
        "po_number": 1705305045002,
        "quantity_ordered": 13000,
        "priority": "NORMAL",
        "current_stage": "KHAZWAL_MATERIAL_PREP",
        "current_status": "WAITING_MATERIAL_PREP",
        "due_date": "2024-03-15",
        "created_at": "2024-01-15 10:30:45"
      }
    ],
    "total_pos": 2,
    "total_po_quantity": 53000,
    "total_with_buffer": 53000,
    "is_personalized": false
  }
}
```

**Error Responses:**

- `400 Bad Request` - ID tidak valid
- `404 Not Found` - OBC Master tidak ditemukan
- `401 Unauthorized` - Token tidak valid
- `403 Forbidden` - Role tidak memiliki akses

---

### 4. Generate Production Orders dari OBC

Generate Production Orders secara manual dari OBC Master dengan formula otomatis.

**Endpoint:** `POST /api/obc/:id/generate-po`

**Path Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | uint64 | Yes | OBC Master ID |

**Request Example:**

```bash
curl -X POST http://localhost:8080/api/obc/1/generate-po \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response Success (201):**

```json
{
  "success": true,
  "message": "Production Orders berhasil di-generate",
  "data": {
    "pos_generated": 2,
    "production_orders": [
      {
        "id": 103,
        "po_number": 1705305145001,
        "quantity_ordered": 40000,
        "priority": "NORMAL",
        "current_stage": "KHAZWAL_MATERIAL_PREP",
        "current_status": "WAITING_MATERIAL_PREP",
        "due_date": "2024-03-15",
        "created_at": "2024-01-15 10:35:45"
      },
      {
        "id": 104,
        "po_number": 1705305145002,
        "quantity_ordered": 13000,
        "priority": "NORMAL",
        "current_stage": "KHAZWAL_MATERIAL_PREP",
        "current_status": "WAITING_MATERIAL_PREP",
        "due_date": "2024-03-15",
        "created_at": "2024-01-15 10:35:45"
      }
    ]
  }
}
```

**PO Generation Logic:**

1. **Total Calculation**: `Total = QTY + (QTY * 6%)` - Buffer 6% untuk antisipasi reject
2. **PO Splitting**: `PO Count = CEIL(Total / 40000)` - Max 40,000 per PO
3. **PO Number**: Timestamp-based untuk uniqueness
4. **Due Date**: Menggunakan due_date dari OBC, atau 30 hari dari sekarang
5. **Priority**: Default NORMAL, bisa di-adjust berdasarkan due date
6. **Product Specs**: Auto-map dari fields OBC (material, seri, warna, dll)

**Error Responses:**

- `400 Bad Request` - ID tidak valid atau QuantityOrdered = 0
- `404 Not Found` - OBC Master tidak ditemukan
- `500 Internal Server Error` - Database error

---

## Business Logic & Validations

### Import Validations

1. **File Format**: Hanya menerima `.xlsx` (Excel 2007+)
2. **Required Fields**: Minimal `No OBC` harus terisi
3. **Unique Constraint**: OBC number harus unique (upsert jika duplikat)
4. **Date Format**: Support Excel serial date dan text formats
5. **Number Format**: Handle comma/dot separator

### PO Generation Rules

1. **Minimum Quantity**: QuantityOrdered harus > 0
2. **Buffer Calculation**: Selalu tambah 6% untuk buffer reject
3. **Max per PO**: 40,000 pcs per Production Order
4. **Denormalization**: OBC data di-copy ke PO untuk performance (OBCNumber, ProductName, etc)
5. **Relationship**: PO memiliki foreign key ke OBCMaster untuk full data access

---

## Data Relationship

```
OBCMaster (1) ----< (N) ProductionOrder
                         |
                         +----< (1) KhazwalMaterialPreparation
                         +----< (N) POStageTracking
```

**Benefits:**
- **Single Source of Truth**: OBC data hanya di satu tempat
- **Easy Updates**: Update OBC akan reflected di semua PO terkait (via relationship)
- **Performance**: Denormalized fields di PO untuk quick access tanpa join
- **Audit Trail**: Track history changes di OBC Master
- **Reporting**: Easy aggregation & reporting dari OBC level

---

## Best Practices

### Import Excel

1. **Validate File**: Check format dan size sebelum upload
2. **Incremental Import**: Untuk file besar, split menjadi batch kecil
3. **Review Failed Rows**: Periksa error detail untuk correction
4. **Re-import**: Safe untuk re-import, upsert akan handle duplikat

### PO Generation

1. **Review OBC First**: Pastikan data OBC benar sebelum generate PO
2. **Auto-generate on Import**: Gunakan `auto_generate_po=true` untuk efficiency
3. **Manual Generate**: Untuk kontrol lebih detail atau re-generate
4. **Check Quantity**: Validate quantity sebelum generate untuk avoid waste

### Filtering & Search

1. **Use Indexed Fields**: Search by material, seri, warna, factory_code untuk performance
2. **Pagination**: Selalu gunakan pagination untuk list besar
3. **Specific Filters**: Combine multiple filters untuk precise results

---

## Sample Excel Template

Download sample template: [obc_master_sample.xlsx](#)

**Template Structure:**

| No OBC | Tgl OBC | Material | SERI | WARNA | ... | (39 columns total) |
|--------|---------|----------|------|-------|-----|--------------------|
| OBC-2024-001 | 15/01/2024 | PITA CUKAI | SER-A | MERAH | ... | ... |

**Tips:**
- Keep column headers exactly as specified (case-sensitive)
- Date format: DD/MM/YYYY or Excel date serial
- Numbers: Use standard format (with or without thousand separator)
- Required: Only "No OBC" is mandatory

---

## Error Codes & Troubleshooting

| Error | Possible Cause | Solution |
|-------|----------------|----------|
| "File tidak ditemukan" | Missing file parameter | Ensure 'file' is in multipart form |
| "Format file tidak valid" | Not .xlsx file | Use Excel 2007+ format (.xlsx) |
| "No OBC tidak boleh kosong" | Missing OBC number in row | Fill OBC number for all rows |
| "invalid date format" | Wrong date format | Use DD/MM/YYYY or Excel date |
| "database error" | Constraint violation | Check for data integrity issues |

---

## Performance Considerations

- **Import Time**: ~1-2 seconds per 100 rows
- **Batch Size**: Recommended max 1000 rows per import
- **PO Generation**: ~10ms per PO created
- **List Query**: Optimized with indexes on material, seri, warna, factory_code
- **Detail Query**: Includes preload, avg ~20-50ms

---

## Changelog

### v1.0.0 (2024-12-30)
- Initial release OBC Master API
- Excel import dengan 39 fields
- Auto PO generation dengan buffer 6%
- List, detail, dan manual generate endpoints
