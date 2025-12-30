# OBC Master Import Testing Guide

## Overview

Panduan testing untuk fitur OBC Master Excel Import yang mencakup test scenarios, sample data, dan expected results.

---

## Prerequisites

### 1. Setup Environment

```bash
# Start backend server
cd backend
go run main.go

# Database harus sudah running
# Pastikan GORM AutoMigrate sudah jalan untuk create tables
```

### 2. Create Test User (PPIC Role)

```sql
-- Via MySQL client atau DBeaver
INSERT INTO users (nip, full_name, email, phone, password_hash, role, department, must_change_password, status) 
VALUES (
    '12345', 
    'Test PPIC User', 
    'ppic@test.local', 
    '081234567890',
    '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIRQ8DrKIW', -- Password: Admin@123
    'PPIC',
    'PPIC',
    FALSE,
    'ACTIVE'
);
```

### 3. Get JWT Token

```bash
# Login to get JWT token
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "nip": "12345",
    "password": "Admin@123"
  }'

# Save the token dari response
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

## Test Scenarios

### Test 1: Basic Import - Valid Data

**Objective:** Import Excel dengan data yang valid untuk verify happy path.

**Sample Data:** `test_data/obc_basic.xlsx`

**Test Steps:**

```bash
# 1. Import tanpa auto-generate PO
curl -X POST http://localhost:8080/api/obc/import \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@test_data/obc_basic.xlsx"

# Expected Response:
# - success: true
# - total_rows: 5
# - success_count: 5
# - failed_count: 0
# - pos_generated: 0
```

**Verification:**

```sql
-- Check database
SELECT COUNT(*) FROM obc_masters; -- Should be 5
SELECT obc_number, material, quantity_ordered FROM obc_masters;
```

---

### Test 2: Import dengan Auto-Generate PO

**Objective:** Verify PO generation logic dengan buffer 6% dan splitting.

**Sample Data:** `test_data/obc_with_po.xlsx`

**Test Data:**
- OBC-TEST-001: QTY = 50,000 → Total with buffer = 53,000 → 2 POs (40k + 13k)
- OBC-TEST-002: QTY = 30,000 → Total with buffer = 31,800 → 1 PO (31,800)
- OBC-TEST-003: QTY = 100,000 → Total with buffer = 106,000 → 3 POs (40k + 40k + 26k)

**Test Steps:**

```bash
# Import dengan auto_generate_po=true
curl -X POST http://localhost:8080/api/obc/import \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@test_data/obc_with_po.xlsx" \
  -F "auto_generate_po=true"

# Expected Response:
# - success: true
# - total_rows: 3
# - success_count: 3
# - failed_count: 0
# - pos_generated: 6 (2 + 1 + 3)
```

**Verification:**

```sql
-- Check OBC Masters
SELECT id, obc_number, quantity_ordered FROM obc_masters 
WHERE obc_number LIKE 'OBC-TEST-%';

-- Check Production Orders generated
SELECT po_number, obc_master_id, obc_number, quantity_ordered 
FROM production_orders 
WHERE obc_number LIKE 'OBC-TEST-%'
ORDER BY obc_number, po_number;

-- Verify quantities
-- OBC-TEST-001: Should have 2 POs with total ~53,000
-- OBC-TEST-002: Should have 1 PO with ~31,800
-- OBC-TEST-003: Should have 3 POs with total ~106,000
```

---

### Test 3: Upsert Logic - Update Existing OBC

**Objective:** Verify bahwa import ulang dengan OBC number sama akan update data.

**Test Steps:**

```bash
# 1. Import pertama kali
curl -X POST http://localhost:8080/api/obc/import \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@test_data/obc_update_v1.xlsx"

# 2. Check data awal
curl -X GET "http://localhost:8080/api/obc?obc_number=OBC-UPDATE-001" \
  -H "Authorization: Bearer $TOKEN"
# Note: quantity_ordered = 10,000

# 3. Import file kedua dengan OBC number sama tapi quantity berbeda
curl -X POST http://localhost:8080/api/obc/import \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@test_data/obc_update_v2.xlsx"

# 4. Check data setelah update
curl -X GET "http://localhost:8080/api/obc?obc_number=OBC-UPDATE-001" \
  -H "Authorization: Bearer $TOKEN"
# Expected: quantity_ordered = 20,000 (updated)
```

**Verification:**

```sql
-- Should only have 1 record (not 2)
SELECT COUNT(*) FROM obc_masters WHERE obc_number = 'OBC-UPDATE-001';
-- Result: 1

-- Verify updated quantity
SELECT quantity_ordered FROM obc_masters WHERE obc_number = 'OBC-UPDATE-001';
-- Result: 20000
```

---

### Test 4: Error Handling - Missing Required Field

**Objective:** Verify bahwa rows dengan missing OBC number akan di-skip.

**Sample Data:** `test_data/obc_with_errors.xlsx`

**Test Data:**
- Row 1: Valid data
- Row 2: Missing OBC number (should fail)
- Row 3: Valid data
- Row 4: Missing OBC number (should fail)
- Row 5: Valid data

**Test Steps:**

```bash
curl -X POST http://localhost:8080/api/obc/import \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@test_data/obc_with_errors.xlsx"

# Expected Response:
# - success: false
# - total_rows: 5
# - success_count: 3
# - failed_count: 2
# - failed_rows: [
#     { row_number: 2, obc_number: "", error: "No OBC tidak boleh kosong" },
#     { row_number: 4, obc_number: "", error: "No OBC tidak boleh kosong" }
#   ]
```

---

### Test 5: Date Format Parsing

**Objective:** Verify berbagai format date bisa di-parse dengan benar.

**Sample Data:** `test_data/obc_date_formats.xlsx`

**Test Data:**
- Row 1: Date in Excel serial format (e.g., 44927 = 2023-01-15)
- Row 2: Date as DD/MM/YYYY (e.g., 15/01/2023)
- Row 3: Date as DD-MM-YYYY (e.g., 15-01-2023)
- Row 4: Date as YYYY-MM-DD (e.g., 2023-01-15)

**Test Steps:**

```bash
curl -X POST http://localhost:8080/api/obc/import \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@test_data/obc_date_formats.xlsx"

# Expected: All 4 rows should import successfully
```

**Verification:**

```sql
-- All dates should be parsed correctly to 2023-01-15
SELECT obc_number, obc_date FROM obc_masters 
WHERE obc_number LIKE 'OBC-DATE-%';
```

---

### Test 6: Large File Import

**Objective:** Test performance dengan file besar (100+ rows).

**Sample Data:** `test_data/obc_large_100.xlsx` (100 rows)

**Test Steps:**

```bash
# Measure import time
time curl -X POST http://localhost:8080/api/obc/import \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@test_data/obc_large_100.xlsx" \
  -F "auto_generate_po=true"

# Expected:
# - Duration: < 3 seconds untuk 100 rows
# - All rows should import successfully
```

---

### Test 7: List & Filter

**Objective:** Test list endpoint dengan berbagai filter.

**Test Steps:**

```bash
# 1. List all (with pagination)
curl -X GET "http://localhost:8080/api/obc?page=1&page_size=20" \
  -H "Authorization: Bearer $TOKEN"

# 2. Filter by material
curl -X GET "http://localhost:8080/api/obc?material=PITA" \
  -H "Authorization: Bearer $TOKEN"

# 3. Filter by seri
curl -X GET "http://localhost:8080/api/obc?seri=2024" \
  -H "Authorization: Bearer $TOKEN"

# 4. Combine filters
curl -X GET "http://localhost:8080/api/obc?material=PITA&seri=2024&warna=MERAH" \
  -H "Authorization: Bearer $TOKEN"

# 5. Search by OBC number
curl -X GET "http://localhost:8080/api/obc?obc_number=OBC-2024-001" \
  -H "Authorization: Bearer $TOKEN"
```

---

### Test 8: Detail OBC with POs

**Objective:** Test detail endpoint untuk melihat OBC beserta Production Orders.

**Test Steps:**

```bash
# 1. Get OBC detail by ID
curl -X GET "http://localhost:8080/api/obc/1" \
  -H "Authorization: Bearer $TOKEN"

# Expected Response should include:
# - Full OBC fields (39 columns)
# - production_orders array
# - total_pos
# - total_po_quantity
# - total_with_buffer
# - is_personalized
```

---

### Test 9: Manual PO Generation

**Objective:** Test manual generate PO dari OBC yang belum punya PO.

**Test Steps:**

```bash
# 1. Import OBC tanpa auto-generate
curl -X POST http://localhost:8080/api/obc/import \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@test_data/obc_no_po.xlsx"

# 2. Get OBC ID dari list
OBC_ID=$(curl -s -X GET "http://localhost:8080/api/obc?obc_number=OBC-MANUAL-001" \
  -H "Authorization: Bearer $TOKEN" | jq '.data.items[0].id')

# 3. Generate PO manually
curl -X POST "http://localhost:8080/api/obc/$OBC_ID/generate-po" \
  -H "Authorization: Bearer $TOKEN"

# Expected:
# - POs created successfully
# - pos_generated > 0
```

**Verification:**

```sql
-- Check POs created
SELECT COUNT(*) FROM production_orders WHERE obc_master_id = $OBC_ID;
```

---

### Test 10: Authorization Test

**Objective:** Verify bahwa non-PPIC users tidak bisa import.

**Test Steps:**

```bash
# 1. Login as STAFF_KHAZWAL (non-PPIC user)
STAFF_TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"nip": "staff_nip", "password": "password"}' | jq -r '.data.access_token')

# 2. Try to import (should fail)
curl -X POST http://localhost:8080/api/obc/import \
  -H "Authorization: Bearer $STAFF_TOKEN" \
  -F "file=@test_data/obc_basic.xlsx"

# Expected Response:
# - Status: 403 Forbidden
# - Message: "Akses ditolak"

# 3. Try to generate PO (should fail)
curl -X POST "http://localhost:8080/api/obc/1/generate-po" \
  -H "Authorization: Bearer $STAFF_TOKEN"

# Expected Response:
# - Status: 403 Forbidden

# 4. Try to list (should work for MANAGER)
MANAGER_TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"nip": "manager_nip", "password": "password"}' | jq -r '.data.access_token')

curl -X GET "http://localhost:8080/api/obc" \
  -H "Authorization: Bearer $MANAGER_TOKEN"

# Expected: Should work (read-only access)
```

---

## Sample Data Files

### obc_basic.xlsx Structure

| No OBC | Tgl OBC | Material | SERI | WARNA | KODE_PABRIK | QTY PESAN | ... |
|--------|---------|----------|------|-------|-------------|-----------|-----|
| OBC-2024-001 | 15/01/2024 | PITA CUKAI MMEA 2024 | SER-2024-A | MERAH | F001 | 50000 | ... |
| OBC-2024-002 | 16/01/2024 | PITA CUKAI MMEA 2024 | SER-2024-A | BIRU | F001 | 30000 | ... |
| OBC-2024-003 | 17/01/2024 | PITA CUKAI HPTL 2024 | SER-2024-B | HIJAU | F002 | 40000 | ... |
| OBC-2024-004 | 18/01/2024 | PITA CUKAI MMEA 2024 | SER-2024-A | KUNING | F001 | 25000 | ... |
| OBC-2024-005 | 19/01/2024 | PITA CUKAI HPTL 2024 | SER-2024-B | PUTIH | F002 | 35000 | ... |

### CSV Sample (untuk quick testing)

Save as `obc_sample.csv` and convert to .xlsx:

```csv
No OBC,Tgl OBC,Material,SERI,WARNA,KODE_PABRIK,QTY PESAN,JHT,RPB,HJE,BPB,RENCET,Tgl JTempo,Perso / non Perso,Perekat,GR,No Pelat,Type,Created On,Sales Doc.,Item,Material description,BUn,Gol. PCA,Kadar Alkohol PCA,Kadar HPTL,Kode Wilayah,OBC Awal,Peruntukan,PESANAN,Plnt,SATUAN,Tahun,Tarif Per Liter,Volume PCA,Warna MMEA
OBC-2024-001,15/01/2024,PITA CUKAI MMEA 2024,SER-2024-A,MERAH,F001,50000,JHT-001,15000,20000,1000,50,15/03/2024,Non Perso,Type A,GR-001,PLT-001,Standard,10/01/2024,SD-001,ITEM-001,Pita Cukai MMEA Tahun 2024,PCS,Cat A,0,0,REG-01,OBC-INIT-001,Domestic,50000,P001,PCS,2024,15000,0,RED-001
OBC-2024-002,16/01/2024,PITA CUKAI MMEA 2024,SER-2024-A,BIRU,F001,30000,JHT-002,15000,20000,1000,50,20/03/2024,Non Perso,Type A,GR-002,PLT-002,Standard,11/01/2024,SD-002,ITEM-002,Pita Cukai MMEA Tahun 2024,PCS,Cat A,0,0,REG-01,OBC-INIT-002,Domestic,30000,P001,PCS,2024,15000,0,BLUE-001
```

---

## Automated Testing Script

Create `test_obc_import.sh`:

```bash
#!/bin/bash

# Configuration
BASE_URL="http://localhost:8080"
TOKEN=""  # Set your JWT token here

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "========================================="
echo "OBC Master Import Testing Suite"
echo "========================================="

# Test 1: Basic Import
echo -e "\n${YELLOW}Test 1: Basic Import${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/api/obc/import" \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@test_data/obc_basic.xlsx")

SUCCESS=$(echo $RESPONSE | jq -r '.success')
if [ "$SUCCESS" == "true" ]; then
  echo -e "${GREEN}✓ Test 1 PASSED${NC}"
else
  echo -e "${RED}✗ Test 1 FAILED${NC}"
  echo $RESPONSE | jq '.'
fi

# Test 2: Import with Auto-Generate PO
echo -e "\n${YELLOW}Test 2: Import with Auto-Generate PO${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/api/obc/import" \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@test_data/obc_with_po.xlsx" \
  -F "auto_generate_po=true")

SUCCESS=$(echo $RESPONSE | jq -r '.success')
POS_GENERATED=$(echo $RESPONSE | jq -r '.data.pos_generated')
if [ "$SUCCESS" == "true" ] && [ "$POS_GENERATED" -gt 0 ]; then
  echo -e "${GREEN}✓ Test 2 PASSED (Generated $POS_GENERATED POs)${NC}"
else
  echo -e "${RED}✗ Test 2 FAILED${NC}"
  echo $RESPONSE | jq '.'
fi

# Add more tests...

echo -e "\n========================================="
echo "Testing Complete"
echo "========================================="
```

---

## Performance Benchmarks

Expected performance metrics:

| Operation | Rows/Records | Expected Time | Notes |
|-----------|--------------|---------------|-------|
| Import (no PO) | 100 rows | < 2s | Basic insert/update |
| Import (with PO) | 100 rows | < 5s | Includes PO generation |
| List (paginated) | 1000 records | < 100ms | With indexes |
| Detail | 1 record | < 50ms | With preload |
| Generate PO | 1 OBC | < 100ms | Depends on quantity |

---

## Troubleshooting

### Issue: Import fails with "column not found"

**Cause:** Excel column headers tidak sesuai dengan mapping.

**Solution:** Check column headers case-sensitive. Use exact names:
- "No OBC" (not "no obc" or "NO OBC")
- "SERI" (not "Seri" or "seri")

### Issue: Dates imported as NULL

**Cause:** Date format tidak recognized.

**Solution:** Use Excel date format atau DD/MM/YYYY text format.

### Issue: Numbers imported incorrectly

**Cause:** Thousands separator or decimal separator issues.

**Solution:** Use standard number format (no separator or comma/dot properly placed).

---

## CI/CD Integration

### GitHub Actions Example

```yaml
name: OBC Import Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.24
      
      - name: Run tests
        run: |
          cd backend
          go test ./services -run TestOBCImportService
          go test ./handlers -run TestOBCHandler
```

---

## Conclusion

Testing checklist:

- [ ] Basic import works
- [ ] Auto PO generation works dengan correct calculation
- [ ] Upsert logic works (update existing OBC)
- [ ] Error handling works (skip bad rows)
- [ ] Date parsing works for multiple formats
- [ ] Large file import performance acceptable
- [ ] List & filter works correctly
- [ ] Detail shows full data with POs
- [ ] Manual PO generation works
- [ ] Authorization works (PPIC only for write, Manager for read)

Jika semua test PASSED, feature siap untuk production deployment.
