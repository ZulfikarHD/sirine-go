# OBC Master Test Data

## Overview

Directory ini berisi sample data untuk testing OBC Master Import feature.

## File Structure

```
test_data/
├── README.md                    # This file
├── obc_sample_basic.csv        # Basic import test (5 OBCs)
├── obc_sample_with_po.csv      # PO generation test (3 OBCs with different quantities)
└── (add more test files as needed)
```

## File Descriptions

### obc_sample_basic.csv

**Purpose:** Basic import testing dengan 5 OBC records yang valid.

**Test Scenarios:**
- Valid data import
- Multiple materials (MMEA dan HPTL)
- Different colors dan factory codes
- Mix of personalization types (Perso dan Non Perso)

**Expected Results:**
- 5 OBC Masters created
- All fields populated correctly
- No errors

### obc_sample_with_po.csv

**Purpose:** Test PO generation dengan berbagai quantities.

**Test Data:**

| OBC Number | Quantity | Total (6%) | Expected POs | Split |
|------------|----------|------------|--------------|-------|
| OBC-TEST-001 | 50,000 | 53,000 | 2 | 40,000 + 13,000 |
| OBC-TEST-002 | 30,000 | 31,800 | 1 | 31,800 |
| OBC-TEST-003 | 100,000 | 106,000 | 3 | 40,000 + 40,000 + 26,000 |

**Expected Results:**
- 3 OBC Masters created
- 6 Production Orders generated (2 + 1 + 3)
- Quantities split correctly per max 40,000 per PO

## How to Use

### Convert CSV to Excel (.xlsx)

**Option 1: Using LibreOffice (Command Line)**

```bash
# Install LibreOffice (if not installed)
sudo apt-get install libreoffice

# Convert CSV to XLSX
libreoffice --headless --convert-to xlsx obc_sample_basic.csv
```

**Option 2: Using Excel/LibreOffice GUI**

1. Open CSV file in Excel atau LibreOffice Calc
2. Save As → Choose format: "Excel Workbook (.xlsx)"
3. Save file

**Option 3: Using Python (xlsxwriter)**

```bash
# Install xlsxwriter
pip install xlsxwriter

# Run conversion script (create convert.py first)
python convert_csv_to_xlsx.py
```

convert_csv_to_xlsx.py:
```python
import csv
import xlsxwriter

def convert_csv_to_xlsx(csv_file, xlsx_file):
    workbook = xlsxwriter.Workbook(xlsx_file)
    worksheet = workbook.add_worksheet()
    
    with open(csv_file, 'r', encoding='utf-8') as f:
        reader = csv.reader(f)
        for row_idx, row in enumerate(reader):
            for col_idx, cell in enumerate(row):
                worksheet.write(row_idx, col_idx, cell)
    
    workbook.close()

# Convert files
convert_csv_to_xlsx('obc_sample_basic.csv', 'obc_basic.xlsx')
convert_csv_to_xlsx('obc_sample_with_po.csv', 'obc_with_po.xlsx')
```

### Run Tests

**Manual Test via curl:**

```bash
# Set your JWT token
TOKEN="your_jwt_token_here"

# Test basic import
curl -X POST http://localhost:8080/api/obc/import \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@obc_basic.xlsx"

# Test with PO generation
curl -X POST http://localhost:8080/api/obc/import \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@obc_with_po.xlsx" \
  -F "auto_generate_po=true"
```

**Automated Test Script:**

```bash
# Run test suite
cd /home/sirinedev/WebApp/Developement/sirine-go
./test_obc_import.sh
```

## Creating Additional Test Files

### Template Structure

All test files must have these 36 column headers (in exact order):

```
No OBC, Tgl OBC, Material, SERI, WARNA, KODE_PABRIK, QTY PESAN, JHT, 
RPB, HJE, BPB, RENCET, Tgl JTempo, Perso / non Perso, Perekat, GR, 
No Pelat, Type, Created On, Sales Doc., Item, Material description, 
BUn, Gol. PCA, Kadar Alkohol PCA, Kadar HPTL, Kode Wilayah, OBC Awal, 
Peruntukan, PESANAN, Plnt, SATUAN, Tahun, Tarif Per Liter, Volume PCA, 
Warna MMEA
```

### Test File Naming Convention

- `obc_basic.xlsx` - Happy path testing
- `obc_with_po.xlsx` - PO generation testing
- `obc_with_errors.xlsx` - Error handling testing
- `obc_date_formats.xlsx` - Date parsing testing
- `obc_large_*.xlsx` - Performance testing (e.g., obc_large_100.xlsx)
- `obc_update_v*.xlsx` - Upsert logic testing

## Sample Data Guidelines

### Realistic Values

**OBC Number Format:** `OBC-YYYY-NNN` atau `OBC-TEST-NNN`
- Production: OBC-2024-001
- Testing: OBC-TEST-001

**Material Types:**
- PITA CUKAI MMEA (Mesin Makanan, Elektronik, Alat)
- PITA CUKAI HPTL (High Price Tobacco Leaf)
- PITA CUKAI EA (Etil Alkohol)

**SERI Format:** `SER-YYYY-X` (e.g., SER-2024-A)

**WARNA:**
- MERAH, BIRU, HIJAU, KUNING, PUTIH, HITAM, ORANGE, UNGU

**Factory Codes:** F001, F002, F003, etc.

**Quantities:**
- Small: 10,000 - 30,000
- Medium: 30,000 - 70,000
- Large: 70,000 - 150,000

**Date Formats (any of these):**
- DD/MM/YYYY (e.g., 15/01/2024)
- DD-MM-YYYY (e.g., 15-01-2024)
- YYYY-MM-DD (e.g., 2024-01-15)
- Excel serial date (e.g., 44927)

## Verification Queries

After import, verify data dengan SQL queries:

```sql
-- Check imported OBCs
SELECT obc_number, material, quantity_ordered, created_at 
FROM obc_masters 
ORDER BY created_at DESC 
LIMIT 10;

-- Check generated POs
SELECT po_number, obc_number, quantity_ordered, current_status 
FROM production_orders 
WHERE obc_master_id IN (
  SELECT id FROM obc_masters WHERE obc_number LIKE 'OBC-TEST-%'
);

-- Verify PO quantity calculations
SELECT 
  om.obc_number,
  om.quantity_ordered as obc_qty,
  ROUND(om.quantity_ordered * 1.06) as qty_with_buffer,
  COUNT(po.id) as po_count,
  SUM(po.quantity_ordered) as total_po_qty
FROM obc_masters om
LEFT JOIN production_orders po ON po.obc_master_id = om.id
WHERE om.obc_number LIKE 'OBC-TEST-%'
GROUP BY om.id;
```

## Troubleshooting

### Issue: "Format file tidak valid"

**Solution:** Make sure file is .xlsx format. Convert CSV to XLSX terlebih dahulu.

### Issue: "No OBC tidak boleh kosong"

**Solution:** Check that No OBC column has values for all rows (kecuali header).

### Issue: Dates imported as NULL

**Solution:** 
1. Use DD/MM/YYYY format as text
2. Or format cells as Date in Excel before saving

### Issue: Numbers wrong

**Solution:**
1. Remove thousands separator atau use consistent format (15000 atau 15,000)
2. Use dot (.) for decimal separator (15000.50)

## CI/CD Integration

Test files ini bisa digunakan dalam automated testing:

```yaml
# .github/workflows/obc-import-test.yml
name: OBC Import Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Setup Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.24
      
      - name: Install dependencies
        run: |
          cd backend
          go mod download
      
      - name: Run OBC Import Tests
        run: |
          cd backend
          go test -v ./services -run TestOBCImportService
          go test -v ./handlers -run TestOBCHandler
```

## Contributing

Ketika menambahkan test files baru:

1. Follow naming convention
2. Document purpose di README ini
3. Ensure data is realistic
4. Add corresponding verification queries
5. Update testing guide (docs/06-testing/obc-import-testing.md)

## References

- [OBC Master API Reference](../docs/04-api-reference/obc-master.md)
- [OBC Import Testing Guide](../docs/06-testing/obc-import-testing.md)
- [Production Order Model](../backend/models/production_order.go)
- [OBC Master Model](../backend/models/obc_master.go)
