# Fix: MaterialPrepHistoryPage - historyItems Undefined Error

## Masalah

Error di console saat mengakses halaman Riwayat Persiapan Material:

```
[Vue warn]: Unhandled error during execution of render function
Uncaught (in promise) TypeError: can't access property "length", $setup.historyItems is undefined
```

## Root Cause

### 1. **Race Condition pada Initial Render**

**Lokasi:** `frontend/src/views/khazwal/MaterialPrepHistoryPage.vue`

**Masalah:**
- Template mencoba akses `historyItems.length` sebelum `historyItems` fully initialized
- Meskipun `historyItems` di-init sebagai `ref([])`, ada race condition dimana Vue render template sebelum ref fully ready
- Jika API call gagal atau response tidak sesuai expected structure, `historyItems.value` bisa menjadi `undefined`

### 2. **Kurangnya Null Safety di API Response Handler**

**Masalah:**
```javascript
// Sebelum fix:
if (response.success) {
  historyItems.value = response.data.items  // Bisa undefined!
  totalItems.value = response.data.total
  totalPages.value = response.data.total_pages
}
```

Jika `response.data.items` adalah `undefined` atau `null`, `historyItems.value` akan set ke `undefined`, bukan array kosong.

### 3. **Tidak Ada Error Handling untuk Array Initialization**

Saat API error, `historyItems.value` tidak di-reset ke `[]`, sehingga bisa tetap `undefined`.

## Solusi

### 1. Enhanced API Response Handler dengan Null Safety

**File:** `frontend/src/views/khazwal/MaterialPrepHistoryPage.vue`

```javascript
const fetchHistory = async () => {
  loading.value = true
  try {
    const response = await khazwalApi.getHistory({
      search: searchQuery.value,
      date_from: dateFrom.value,
      date_to: dateTo.value,
      page: currentPage.value,
      per_page: perPage.value
    })

    // Null safety dengan fallback ke empty array
    if (response.success && response.data) {
      historyItems.value = response.data.items || []
      totalItems.value = response.data.total || 0
      totalPages.value = response.data.total_pages || 0
    } else {
      // Fallback jika response tidak sesuai expected structure
      historyItems.value = []
      totalItems.value = 0
      totalPages.value = 0
    }
  } catch (error) {
    console.error('Error fetching history:', error)
    
    // PENTING: Ensure historyItems tetap array meskipun error
    historyItems.value = []
    totalItems.value = 0
    totalPages.value = 0
    
    alertDialog.error('Gagal memuat riwayat', {
      detail: error.response?.data?.message || 'Silakan coba lagi'
    })
  } finally {
    loading.value = false
  }
}
```

**Perubahan:**
- ✅ Tambah null check: `response.data.items || []`
- ✅ Tambah fallback untuk response tidak valid
- ✅ Reset `historyItems` ke `[]` dalam catch block
- ✅ Ensure `historyItems` SELALU array, never undefined

### 2. Template Guards untuk Extra Safety

**File:** `frontend/src/views/khazwal/MaterialPrepHistoryPage.vue`

```vue
<!-- History Grid -->
<div v-if="!loading && historyItems && historyItems.length > 0" class="grid gap-4">
  <!-- Tambah check: historyItems && historyItems.length -->
</div>

<!-- Empty State -->
<Motion v-if="!loading && historyItems && historyItems.length === 0">
  <!-- Tambah check: historyItems && historyItems.length -->
</Motion>
```

**Perubahan:**
- ✅ Tambah explicit check `historyItems &&` sebelum akses `.length`
- ✅ Prevent error jika `historyItems` somehow masih undefined
- ✅ Defense in depth approach

## Backend Verification

### Handler: `backend/handlers/khazwal_handler.go`

```go
func (h *KhazwalHandler) GetHistory(c *gin.Context) {
    // ... validation code ...
    
    historyResponse, err := h.khazwalService.GetMaterialPrepHistory(filters)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": "Gagal mengambil riwayat material preparation",
            "error":   err.Error(),
        })
        return
    }

    // Transform ke DTO
    items := make([]HistoryItemDTO, 0, len(historyResponse.Items))
    // ...
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Riwayat material preparation berhasil diambil",
        "data": gin.H{
            "items":       items,
            "total":       historyResponse.Total,
            "page":        historyResponse.Page,
            "per_page":    historyResponse.PerPage,
            "total_pages": historyResponse.TotalPages,
        },
    })
}
```

✅ Backend properly returns `items` array (never null)
✅ Proper error handling with clear messages
✅ Consistent response structure

### Service: `backend/services/khazwal_service.go`

```go
func (s *KhazwalService) GetMaterialPrepHistory(filters HistoryFilters) (*HistoryResponse, error) {
    // ... query building ...
    
    // Transform ke HistoryItem
    items := make([]HistoryItem, 0, len(preps))  // Empty slice, not nil
    for _, prep := range preps {
        // ... transform logic ...
        items = append(items, item)
    }
    
    return &HistoryResponse{
        Items:      items,  // Always a slice, never nil
        Total:      int(total),
        Page:       filters.Page,
        PerPage:    filters.PerPage,
        TotalPages: totalPages,
    }, nil
}
```

✅ Service always returns empty slice `[]`, never `nil`
✅ Proper initialization with `make([]HistoryItem, 0, len(preps))`

## Possible Scenarios yang Diatasi

### Scenario 1: Database Kosong (Belum Ada History)
```
✅ Backend returns: { items: [], total: 0, total_pages: 0 }
✅ Frontend: historyItems = [] (array kosong)
✅ Template: Tampil empty state "Belum Ada Riwayat"
```

### Scenario 2: Network Error
```
✅ Catch block triggered
✅ historyItems = [] (reset to empty array)
✅ Alert dialog shown: "Gagal memuat riwayat"
✅ Template: Tampil empty state dengan message error
```

### Scenario 3: Invalid Response Structure
```
✅ response.data === undefined → fallback block executed
✅ historyItems = []
✅ Template: Safe rendering, no crash
```

### Scenario 4: API Timeout/401 Error
```
✅ Interceptor handles refresh token
✅ If refresh fails → caught by catch block
✅ historyItems = []
✅ User redirected to login (dari interceptor)
```

### Scenario 5: Normal Success dengan Data
```
✅ Backend returns items array
✅ historyItems = response.data.items
✅ Template: Render cards dengan stagger animation
```

## Testing Checklist

- [x] Empty database → tampil empty state
- [x] API error → tampil error alert + empty state
- [x] Network offline → handled gracefully
- [x] Invalid response structure → no crash
- [x] Success dengan data → render properly
- [x] Search filter → works
- [x] Date range filter → works
- [x] Pagination → works

## Benefits

1. **No More Crashes:** Template tidak pernah crash karena undefined access
2. **Better Error Handling:** Clear error messages untuk user
3. **Defense in Depth:** Multiple layers of protection (API handler + template guard)
4. **Consistent State:** `historyItems` SELALU array, never undefined/null
5. **Better UX:** Graceful degradation saat error

## Related Files

- `frontend/src/views/khazwal/MaterialPrepHistoryPage.vue` - Main page component
- `frontend/src/components/khazwal/PrepHistoryCard.vue` - Card component
- `frontend/src/composables/useKhazwalApi.js` - API composable
- `backend/handlers/khazwal_handler.go` - Handler endpoint
- `backend/services/khazwal_service.go` - Service logic
- `backend/routes/routes.go` - Route configuration

## Additional Notes

### Token Refresh Issue (Related Fix)

Jika user mengalami logout tiba-tiba saat mengakses halaman ini, lihat:
`docs/fixes/token-refresh-issue-fix.md`

### Best Practices Applied

1. **Always Initialize Arrays:** `ref([])` not `ref(null)`
2. **Null Safety in API Handlers:** Use `|| []` fallback
3. **Reset State on Error:** Ensure consistent state
4. **Template Guards:** Add extra safety with `v-if="array &&"`
5. **Descriptive Error Messages:** Help user understand what happened
