# Epic 8: Integration & Automation

**Epic ID:** KHAZWAL-EPIC-08  
**Priority:** 🟢 Low (Phase 3)  
**Estimated Duration:** 2-3 Minggu  

---

## 📋 Overview

Epic ini mencakup integrasi dengan sistem eksternal (SAP) dan unit lain (Cetak, Verifikasi) untuk automasi notifikasi dan data sync.

---

## 🗄️ Database Reference

### Tables
- `production_orders` - Status updates
- `notifications` - Cross-unit notifications
- `activity_logs` - Integration logs

### External Systems
- SAP - Material consumption
- Unit Cetak - Material ready notification
- Tim Verifikasi - Cutting complete notification

---

## 📝 Backlog Items

### US-KW-025: Integrasi dengan SAP - Material Consumption

| Field | Value |
|-------|-------|
| **ID** | US-KW-025 |
| **Story Points** | 13 |
| **Priority** | 🟢 Low |
| **Dependencies** | SAP API Ready |

**User Story:**
> Sebagai System, saya ingin otomatis update consumption material ke SAP, sehingga inventory SAP selalu accurate.

#### Tasks Breakdown

| Task ID | Deskripsi | Estimasi | Assignee |
|---------|-----------|----------|----------|
| KW-025-BE-01 | Create `SAPIntegrationService.php` | 4h | Backend |
| KW-025-BE-02 | Create material consumption API client | 3h | Backend |
| KW-025-BE-03 | Implement kertas consumption posting | 2h | Backend |
| KW-025-BE-04 | Implement tinta consumption posting | 2h | Backend |
| KW-025-BE-05 | Implement waste posting | 2h | Backend |
| KW-025-BE-06 | Create retry mechanism (3 attempts) | 2h | Backend |
| KW-025-BE-07 | Create fallback queue for failed posts | 3h | Backend |
| KW-025-BE-08 | Create SAP sync status tracking | 2h | Backend |
| KW-025-BE-09 | Create scheduled job for retry queue | 2h | Backend |
| KW-025-BE-10 | Log all SAP transactions | 1h | Backend |
| KW-025-BE-11 | Create alert for failed sync | 2h | Backend |

#### Acceptance Criteria
- [ ] Setiap finalisasi penyiapan material:
  - Sistem otomatis kirim data consumption ke SAP:
    - Kertas blanko: [X] lembar
    - Tinta per warna: [Y] kg
    - Timestamp
    - PO number
- [ ] Setiap finalisasi pemotongan:
  - Sistem otomatis kirim data waste ke SAP
- [ ] Error handling jika SAP tidak respond
- [ ] Retry mechanism
- [ ] Log semua transaksi

#### SAP API Payload
```json
// POST /sap/api/material-consumption
{
  "document_type": "CONSUMPTION",
  "reference_number": "PO-12345",
  "posting_date": "2025-12-27",
  "items": [
    {
      "material_code": "KERTAS-BLANKO-001",
      "quantity": 5000,
      "unit": "LEMBAR",
      "cost_center": "PROD-KHAZWAL"
    },
    {
      "material_code": "TINTA-RED-001",
      "quantity": 5.5,
      "unit": "KG",
      "cost_center": "PROD-KHAZWAL"
    }
  ]
}
```

#### Error Handling
```php
class SAPIntegrationService
{
    public function postConsumption(array $data): void
    {
        $maxRetries = 3;
        $attempt = 0;
        
        while ($attempt < $maxRetries) {
            try {
                $response = $this->sapClient->post('/material-consumption', $data);
                
                if ($response->successful()) {
                    $this->logSuccess($data, $response);
                    return;
                }
                
                throw new SAPException($response->body());
            } catch (Exception $e) {
                $attempt++;
                
                if ($attempt >= $maxRetries) {
                    $this->queueForRetry($data);
                    $this->createAlert('SAP sync failed', $e->getMessage());
                }
                
                sleep(pow(2, $attempt)); // Exponential backoff
            }
        }
    }
}
```

---

### US-KW-026: Integrasi dengan Unit Cetak

| Field | Value |
|-------|-------|
| **ID** | US-KW-026 |
| **Story Points** | 5 |
| **Priority** | 🟢 Low |
| **Dependencies** | US-KW-006 |

**User Story:**
> Sebagai System, saya ingin otomatis notifikasi ke Unit Cetak saat material siap, sehingga Unit Cetak langsung tahu dan bisa mulai cetak.

#### Tasks Breakdown

| Task ID | Deskripsi | Estimasi | Assignee |
|---------|-----------|----------|----------|
| KW-026-BE-01 | Create event `MaterialPrepCompleted` | 2h | Backend |
| KW-026-BE-02 | Create listener for Cetak notification | 2h | Backend |
| KW-026-BE-03 | Update Cetak queue with new PO | 2h | Backend |
| KW-026-BE-04 | Implement WebSocket broadcast | 3h | Backend |
| KW-026-BE-05 | Create urgent notification for reprint | 2h | Backend |
| KW-026-FE-01 | Handle real-time notification (Cetak side) | 2h | Frontend |

#### Acceptance Criteria
- [ ] Setiap finalisasi penyiapan material:
  - Notifikasi ke Unit Cetak (in-app)
  - Update status PO di sistem Unit Cetak
  - Kirim detail material yang sudah disiapkan
- [ ] Setiap ada work order cetak ulang:
  - Notifikasi urgent ke Unit Cetak
  - Prioritas di queue Unit Cetak

#### Event Flow
```
MaterialPrepCompleted Event
    │
    ├── NotifyCetakListener
    │       ├── Create notification for Cetak users
    │       └── Broadcast via WebSocket
    │
    └── UpdateCetakQueueListener
            └── Add PO to Cetak queue
```

---

### US-KW-027: Integrasi dengan Verifikasi

| Field | Value |
|-------|-------|
| **ID** | US-KW-027 |
| **Story Points** | 5 |
| **Priority** | 🟢 Low |
| **Dependencies** | US-KW-014 |

**User Story:**
> Sebagai System, saya ingin otomatis notifikasi ke Tim Verifikasi saat pemotongan selesai, sehingga Tim Verifikasi langsung tahu dan bisa mulai QC.

#### Tasks Breakdown

| Task ID | Deskripsi | Estimasi | Assignee |
|---------|-----------|----------|----------|
| KW-027-BE-01 | Create event `CuttingCompleted` | 2h | Backend |
| KW-027-BE-02 | Create listener for Verifikasi notification | 2h | Backend |
| KW-027-BE-03 | Generate verification labels (already in US-KW-014) | - | Backend |
| KW-027-BE-04 | Update Verifikasi queue | 2h | Backend |
| KW-027-BE-05 | Implement WebSocket broadcast | 2h | Backend |
| KW-027-FE-01 | Handle real-time notification (Verifikasi side) | 2h | Frontend |

#### Acceptance Criteria
- [ ] Setiap finalisasi pemotongan:
  - Notifikasi ke Tim Verifikasi
  - Update status PO di sistem Verifikasi
  - Kirim detail hasil pemotongan:
    - Total lembar kirim
    - Sisiran kiri & kanan
    - Info kerusakan (jika ada dari penghitungan)

#### Event Flow
```
CuttingCompleted Event
    │
    ├── NotifyVerifikasiListener
    │       ├── Create notification for QC users
    │       └── Broadcast via WebSocket
    │
    ├── GenerateVerificationLabelsListener
    │       └── Create verification_labels records
    │
    └── UpdateVerifikasiQueueListener
            └── Add labels to Verifikasi queue
```

---

## 📊 Epic Summary

| User Story | Story Points | Priority | Phase |
|------------|--------------|----------|-------|
| US-KW-025 | 13 | Low | Phase 3 |
| US-KW-026 | 5 | Low | Phase 3 |
| US-KW-027 | 5 | Low | Phase 3 |
| **Total** | **23** | - | - |

---

## 🔗 Dependencies Graph

```
US-KW-006 (Finalize Material Prep)
    │
    └── US-KW-026 (Notify Cetak)
            │
            └── [Cetak Queue Updated]

US-KW-014 (Finalize Cutting)
    │
    └── US-KW-027 (Notify Verifikasi)
            │
            └── [Verifikasi Queue Updated]

US-KW-006, US-KW-014
    │
    └── US-KW-025 (SAP Sync)
            │
            └── [SAP Stock Updated]
```

---

## 🧪 Testing Checklist

### Unit Tests
- [ ] SAPIntegrationService - post consumption
- [ ] SAPIntegrationService - retry mechanism
- [ ] Event dispatching
- [ ] Listener handlers

### Integration Tests
- [ ] SAP API integration (mock)
- [ ] WebSocket broadcast
- [ ] Cross-unit notification delivery

### E2E Tests
- [ ] Complete flow: Finalize → SAP sync → Cetak notified
- [ ] Complete flow: Finalize → Labels generated → Verifikasi notified
- [ ] SAP failure → Retry → Alert

---

## 🔄 Integration Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    KHAZANAH AWAL                        │
│                                                         │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐ │
│  │ Material    │    │  Counting   │    │   Cutting   │ │
│  │ Preparation │───►│   Results   │───►│   Results   │ │
│  └─────────────┘    └─────────────┘    └─────────────┘ │
│         │                                     │         │
│         ▼                                     ▼         │
│  ┌─────────────┐                     ┌─────────────┐   │
│  │ SAP Sync    │                     │ SAP Sync    │   │
│  │ (Consume)   │                     │ (Waste)     │   │
│  └─────────────┘                     └─────────────┘   │
│         │                                     │         │
└─────────┼─────────────────────────────────────┼─────────┘
          │                                     │
          ▼                                     ▼
┌─────────────────────┐               ┌─────────────────────┐
│        SAP          │               │     VERIFIKASI      │
│  (Material Stock)   │               │  (Verification      │
└─────────────────────┘               │     Labels)         │
                                      └─────────────────────┘
          │
          ▼
┌─────────────────────┐
│       CETAK         │
│  (Print Queue)      │
└─────────────────────┘
```

---

**Last Updated:** 27 December 2025  
**Status:** Ready for Development
