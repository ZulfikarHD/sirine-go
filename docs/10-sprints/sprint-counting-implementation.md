# 📦 Sprint: Epic 2 - Penghitungan (Counting)

**Version:** 1.0.0  
**Date:** 30 Desember 2025  
**Duration:** 1 Sprint  
**Status:** ✅ Completed

## 📋 Sprint Goals

Implementasi komprehensif P0 (Critical) items untuk Epic 2 Penghitungan yang mencakup backend API, database schema, frontend pages/components, dan integration flow yang diperlukan untuk MVP counting feature. Sistem ini bertujuan untuk memungkinkan staff khazwal melakukan penghitungan hasil cetak dengan tracking yang akurat, yaitu: input jumlah baik dan rusak, breakdown kerusakan detail, real-time calculation, dan finalisasi yang aman dengan data immutability.

---

## ✨ Features Implemented

### US-KW-007: View Counting Queue
- Display FIFO list PO yang menunggu penghitungan
- Real-time waiting time calculation
- Overdue warning (> 2 jam)
- Auto-refresh setiap 30 detik
- Mobile-responsive card layout

### US-KW-008: Start Counting
- Create counting record dengan status tracking
- Update PO status ke SEDANG_DIHITUNG
- Log activity untuk audit trail
- Timestamp akurat untuk duration tracking

### US-KW-009: Input Counting Result
- Input jumlah baik dan rusak dengan validation
- Real-time calculation (total, variance, percentages)
- Auto-save dengan debouncing
- Conditional defect breakdown form (rusak > 5%)
- Variance reason requirement

### US-KW-010: Finalize Counting
- Confirmation modal dengan ringkasan lengkap
- Data immutability setelah finalisasi
- Auto-advance PO ke SIAP_POTONG
- Duration tracking otomatis
- Haptic feedback untuk mobile UX

---

## 📁 File Structure

### Backend Files

```
backend/
├── database/migrations/
│   └── 003_create_khazwal_counting_results.sql  ✨ NEW - Schema counting
│
├── internal/counting/                            ✨ NEW - Counting module
│   ├── model.go                                  ✨ NEW - Data structures & DTOs
│   ├── validator.go                              ✨ NEW - Business rule validation
│   ├── repository.go                             ✨ NEW - Database operations
│   ├── service.go                                ✨ NEW - Business logic layer
│   └── handler.go                                ✨ NEW - HTTP handlers (5 endpoints)
│
└── routes/
    └── routes.go                                 ✏️ UPDATED - Register counting routes
```

### Frontend Files

```
frontend/src/
├── composables/
│   └── useCountingApi.js                        ✨ NEW - API composable
│
├── stores/
│   └── counting.js                              ✨ NEW - Pinia store dengan real-time calc
│
├── components/counting/                         ✨ NEW - Counting components
│   ├── CountingQueueCard.vue                    ✨ NEW - Queue item display
│   ├── CountingPrintInfo.vue                    ✨ NEW - Print summary info
│   ├── CountingInputForm.vue                    ✨ NEW - Input good/defect
│   ├── CountingCalculation.vue                  ✨ NEW - Real-time calculations
│   ├── DefectBreakdownForm.vue                  ✨ NEW - Breakdown kerusakan
│   └── CountingFinalizeModal.vue                ✨ NEW - Finalize confirmation
│
├── views/khazwal/counting/                      ✨ NEW - Counting pages
│   ├── CountingQueuePage.vue                    ✨ NEW - FIFO queue list
│   └── CountingWorkPage.vue                     ✨ NEW - Counting workflow
│
├── router/
│   └── index.js                                 ✏️ UPDATED - Register routes
│
└── components/layout/
    └── Sidebar.vue                              ✏️ UPDATED - Add menu item
```

---

## 🔌 API Endpoints Summary

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/khazwal/counting/queue` | Get FIFO counting queue | Yes |
| GET | `/api/khazwal/counting/:id` | Get counting detail | Yes |
| POST | `/api/khazwal/counting/:po_id/start` | Start counting process | Yes |
| PATCH | `/api/khazwal/counting/:id/result` | Update counting results | Yes |
| POST | `/api/khazwal/counting/:id/finalize` | Finalize counting | Yes |

**Roles Allowed:** STAFF_KHAZWAL, ADMIN, MANAGER

> 📡 Full API documentation: [Counting API](../04-api-reference/counting.md)

---

## 💾 Database Changes

### New Table: `khazwal_counting_results`

**Purpose:** Tracking hasil penghitungan dengan immutable history

**Key Columns:**
- `quantity_good`, `quantity_defect` - Input manual
- `total_counted` - Generated column (good + defect)
- `variance_from_target` - Selisih dari target
- `percentage_good`, `percentage_defect` - Calculated percentages
- `defect_breakdown` - JSONB untuk breakdown detail
- `status` - PENDING | IN_PROGRESS | COMPLETED
- `started_at`, `completed_at`, `duration_minutes` - Time tracking

**New PO Statuses:**
- `SELESAI_CETAK` - Print job finalized
- `WAITING_COUNTING` - In counting queue
- `SEDANG_DIHITUNG` - Counting in progress
- `SIAP_POTONG` - Ready for cutting

**New PO Stages:**
- `KHAZWAL_COUNTING` - Counting stage
- `KHAZWAL_CUTTING` - Cutting stage (next)

---

## 🎨 Frontend Features

### Design System Compliance
- ✅ Motion-V animations (fadeUp, fadeScale, stagger)
- ✅ iOS-inspired spring physics
- ✅ Mobile-first responsive
- ✅ Glass morphism cards
- ✅ Haptic feedback
- ✅ Active-scale press feedback
- ✅ AppLayout integration

### State Management
- Pinia store dengan reactive calculations
- Auto-save dengan debouncing (1 second)
- Optimistic UI updates
- Error handling dengan recovery

### Validations
- Client-side: No negative numbers
- Real-time: Defect breakdown sum validation
- Conditional: Breakdown required if defect > 5%
- Conditional: Variance reason required if != 0
- Server-side: All business rules enforced

---

## 🧪 Key Business Rules Implemented

| Rule ID | Description | Implementation |
|---------|-------------|----------------|
| BR-COUNT-01 | PO harus status WAITING_COUNTING untuk start | Service validation |
| BR-COUNT-02 | Tidak boleh duplikat counting IN_PROGRESS | Database unique constraint |
| BR-COUNT-03 | Defect breakdown wajib jika > 5% | Validator.go + Frontend |
| BR-COUNT-04 | Variance reason wajib jika != 0 | Validator.go + Frontend |
| BR-COUNT-05 | Data immutable setelah COMPLETED | Status check di service |
| BR-COUNT-06 | PATCH result multiple times before finalize | No restriction on IN_PROGRESS |
| BR-COUNT-07 | Breakdown sum harus match quantity_defect | ValidateDefectBreakdownSum() |
| BR-COUNT-08 | Auto-advance PO ke SIAP_POTONG | Transaction di finalize |

---

## 📊 Testing Coverage

### Unit Tests Required
- [ ] Backend validation logic (validator.go)
- [ ] Calculation accuracy (service.go)
- [ ] Real-time calculations (counting.js store)
- [ ] Defect breakdown validation (Frontend)

### Integration Tests Required
- [ ] Full counting flow (start → input → finalize)
- [ ] Concurrent counting prevention
- [ ] Transaction atomicity (finalize)
- [ ] PO status transitions

### Manual Testing Checklist
- [ ] Queue display dengan FIFO sorting
- [ ] Overdue warning (> 2 jam)
- [ ] Start counting success
- [ ] Input form dengan number keyboard (mobile)
- [ ] Real-time calculation accuracy
- [ ] Defect breakdown conditional rendering
- [ ] Variance reason conditional rendering
- [ ] Auto-save functionality
- [ ] Finalize modal confirmation
- [ ] Data locked setelah finalize
- [ ] Mobile responsive semua pages
- [ ] Loading states
- [ ] Error handling

> 📋 Full test plan: [Counting Test Plan](../06-testing/counting-testing.md)

---

## 🔐 Security Considerations

| Concern | Mitigation | Implementation |
|---------|------------|----------------|
| Unauthorized access | Role-based middleware | RequireRole in routes |
| Concurrent counting | Unique constraint + lock | DB constraint + SELECT FOR UPDATE |
| Data tampering after finalize | Immutability check | Status validation in all updates |
| Transaction atomicity | Database transactions | tx.Begin() / tx.Commit() |
| Activity logging | Audit trail | activity_logs insertion |

---

## 📱 Mobile UX Features

- Large input fields (56px min height)
- Number keyboard (`inputmode="numeric"`)
- Haptic feedback (vibration API)
- Active-scale press feedback
- Bottom sheet style modals
- Staggered animations
- Auto-refresh yang smooth
- Loading skeletons

---

## 🔗 Related Documentation

- **API Reference:** [Counting API](../04-api-reference/counting.md)
- **Testing Guide:** [Counting Testing](../06-testing/counting-testing.md)
- **User Journeys:** [Counting Flow](../07-user-journeys/counting/)
- **Plan Document:** `.cursor/plans/p0_counting_implementation_36b4b35c.plan.md`

---

## 🚀 Deployment Notes

### Prerequisites
1. Run migration: `003_create_khazwal_counting_results.sql`
2. Ensure PO dengan status `SELESAI_CETAK` exist
3. Ensure `print_job_summaries` table populated

### Verification Steps
1. Check migration status: `migrate:status`
2. Test counting queue endpoint returns data
3. Complete full flow: queue → start → input → finalize
4. Verify PO advances to SIAP_POTONG
5. Check activity logs created

### Rollback Plan
- Migration down will drop table
- Frontend routes dapat di-comment out
- Backend routes dapat di-disable di routes.go

---

## 📝 Implementation Notes

### Architecture Decisions
- **Internal package:** Counting module di `internal/` untuk encapsulation
- **Repository pattern:** Separation of concerns yang jelas
- **DTOs:** Response DTOs terpisah untuk clean API contract
- **Generated columns:** `total_counted` untuk data consistency
- **JSONB:** Flexible defect_breakdown storage

### Performance Considerations
- Indexes on `status`, `production_order_id`, `counted_by`
- Auto-refresh interval: 30 detik (balance freshness vs load)
- Debounced auto-save: 1 detik (balance UX vs API calls)
- Optimistic updates: Immediate UI feedback

### Known Limitations (P0 Scope)
- ❌ No search/filter on queue (FIFO only)
- ❌ No edit after finalize (immutable by design)
- ❌ No cancel counting (need cleanup job for orphaned records)
- ❌ No real-time WebSocket (polling only)
- ❌ No multi-language defect types (Indonesian only)

---

## 📈 Success Metrics

### Technical Metrics
- ✅ 5 API endpoints implemented
- ✅ 1 database migration
- ✅ 6 reusable components
- ✅ 2 pages dengan full workflow
- ✅ 100% P0 requirements covered

### Business Metrics (Post-Launch)
- Target: < 5 menit average counting duration
- Target: < 2% data entry error rate
- Target: > 60% mobile usage
- Target: < 30 menit response time untuk overdue

---

## 🔄 Future Enhancements (NOT P0)

| Priority | Feature | Description |
|----------|---------|-------------|
| P1 | Resume counting | Allow resume IN_PROGRESS counting |
| P1 | Cancel counting | Soft delete dengan reason |
| P2 | Search & filter | Queue filtering by date, machine |
| P2 | Real-time updates | WebSocket untuk live queue |
| P2 | Bulk operations | Process multiple PO together |
| P3 | Custom defect types | Admin-configurable defect list |
| P3 | Photo evidence | Upload foto hasil penghitungan |
| P3 | Supervisor approval | Approval flow untuk variance tinggi |

---

*Last Updated: 30 Desember 2025*
