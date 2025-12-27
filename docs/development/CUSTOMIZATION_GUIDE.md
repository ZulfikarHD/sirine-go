# 🎨 Customization Guide - Sirine Go App

Panduan praktis untuk customize dan extend aplikasi sesuai kebutuhan Anda.

> **📖 Prerequisites:** Sudah baca [ARCHITECTURE_EXPLAINED.md](./ARCHITECTURE_EXPLAINED.md) untuk memahami struktur project.

---

## 📋 Daftar Isi

1. [Menambah Model & API Baru](#1-menambah-model--api-baru)
2. [Menambah Middleware](#2-menambah-middleware)
3. [Menambah Validation](#3-menambah-validation)
4. [Menambah Relationships](#4-menambah-relationships)
5. [Menambah Vue Component](#5-menambah-vue-component)
6. [Official Documentation](#6-official-documentation)

---

## 🏗️ Architecture Overview

Aplikasi ini menggunakan **Service Pattern**:

```
HTTP Request → Route → Handler → Service → Database
```

**📖 Detail lengkap:** Lihat [ARCHITECTURE_EXPLAINED.md](./ARCHITECTURE_EXPLAINED.md)

---

## 1. Menambah Model & API Baru

Contoh: Menambah **Product** model dengan CRUD lengkap.

### Step 1: Create Model

Buat file `backend/models/product.go`:

```go
package models

import ("time"; "gorm.io/gorm")

type Product struct {
    ID          uint           `gorm:"primarykey" json:"id"`
    Name        string         `gorm:"size:255;not null" json:"name" binding:"required"`
    Price       float64        `gorm:"type:decimal(10,2)" json:"price" binding:"required,gt=0"`
    Stock       int            `gorm:"default:0" json:"stock" binding:"gte=0"`
    IsActive    bool           `gorm:"default:true" json:"is_active"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
```

**📖 GORM Model Docs:** [gorm.io/docs/models.html](https://gorm.io/docs/models.html)

---

### Step 2: Create Service

Buat file `backend/services/product_service.go`:

```go
package services

import (
    "sirine-go/backend/database"
    "sirine-go/backend/models"
)

type ProductService struct{}

func NewProductService() *ProductService {
    return &ProductService{}
}

func (s *ProductService) GetAll() ([]models.Product, error) {
    var products []models.Product
    result := database.GetDB().Where("is_active = ?", true).Find(&products)
    return products, result.Error
}

func (s *ProductService) GetByID(id uint) (*models.Product, error) {
    var product models.Product
    result := database.GetDB().First(&product, id)
    return &product, result.Error
}

func (s *ProductService) Create(product *models.Product) error {
    return database.GetDB().Create(product).Error
}

func (s *ProductService) Update(id uint, product *models.Product) error {
    return database.GetDB().Model(&models.Product{}).
        Where("id = ?", id).Updates(product).Error
}

func (s *ProductService) Delete(id uint) error {
    return database.GetDB().Delete(&models.Product{}, id).Error
}
```

**📖 GORM CRUD Docs:** [gorm.io/docs/create.html](https://gorm.io/docs/create.html)

---

### Step 3: Create Handler

Buat file `backend/handlers/product_handler.go`:

```go
package handlers

import (
    "net/http"
    "strconv"
    "sirine-go/backend/models"
    "sirine-go/backend/services"
    "github.com/gin-gonic/gin"
)

type ProductHandler struct {
    service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
    return &ProductHandler{service: service}
}

func (h *ProductHandler) GetAll(c *gin.Context) {
    products, err := h.service.GetAll()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Gagal mengambil data produk",
        })
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": products})
}

func (h *ProductHandler) GetByID(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
    product, err := h.service.GetByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": product})
}

func (h *ProductHandler) Create(c *gin.Context) {
    var product models.Product
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := h.service.Create(&product); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan produk"})
        return
    }
    c.JSON(http.StatusCreated, gin.H{
        "message": "Produk berhasil dibuat",
        "data":    product,
    })
}

func (h *ProductHandler) Update(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
    var product models.Product
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := h.service.Update(uint(id), &product); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui produk"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil diperbarui"})
}

func (h *ProductHandler) Delete(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
    if err := h.service.Delete(uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus produk"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil dihapus"})
}
```

**📖 Gin Handler Docs:** [gin-gonic.com/docs/examples](https://gin-gonic.com/docs/examples/)

---

### Step 4: Add Routes

Edit `backend/routes/routes.go`, tambahkan:

```go
// Product routes
productService := services.NewProductService()
productHandler := handlers.NewProductHandler(productService)

products := api.Group("/products")
{
    products.GET("", productHandler.GetAll)
    products.GET("/:id", productHandler.GetByID)
    products.POST("", productHandler.Create)
    products.PUT("/:id", productHandler.Update)
    products.DELETE("/:id", productHandler.Delete)
}
```

**📖 Gin Routing Docs:** [gin-gonic.com/docs/examples/grouping-routes](https://gin-gonic.com/docs/examples/grouping-routes/)

---

### Step 5: Add Migration

Edit `backend/cmd/server/main.go`, tambahkan model ke migration:

```go
if err := database.AutoMigrate(
    &models.Example{},
    &models.Product{},  // ← Add this
); err != nil {
    log.Fatal("Failed to migrate database:", err)
}
```

---

### Step 6: Test API

```bash
# Restart backend
Ctrl+C
make dev-backend

# Test endpoints
curl http://localhost:8080/api/products

curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop Gaming","price":15000000,"stock":5}'

curl http://localhost:8080/api/products/1

curl -X PUT http://localhost:8080/api/products/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop Gaming Updated","price":14000000}'

curl -X DELETE http://localhost:8080/api/products/1
```

✅ **Done!** Model dan API baru sudah berfungsi.

---

## 2. Menambah Middleware

Contoh: Authentication middleware.

### Create Middleware

Buat file `backend/middleware/auth.go`:

```go
package middleware

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        
        if token == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Token tidak ditemukan",
            })
            c.Abort()
            return
        }

        token = strings.TrimPrefix(token, "Bearer ")

        // TODO: Validate JWT token here
        if !isValidToken(token) {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Token tidak valid",
            })
            c.Abort()
            return
        }

        c.Set("user_id", getUserIDFromToken(token))
        c.Next()
    }
}

func isValidToken(token string) bool {
    // TODO: Implement JWT validation
    return token == "valid-token" // Placeholder
}

func getUserIDFromToken(token string) uint {
    // TODO: Extract user ID from JWT
    return 1 // Placeholder
}
```

### Apply Middleware

Di `backend/routes/routes.go`:

```go
// Protected routes
products := api.Group("/products")
products.Use(middleware.AuthRequired()) // ← Apply middleware
{
    products.GET("", productHandler.GetAll)
    products.POST("", productHandler.Create)
    // ... other routes
}
```

**📖 Gin Middleware Docs:** [gin-gonic.com/docs/examples/custom-middleware](https://gin-gonic.com/docs/examples/custom-middleware/)

---

## 3. Menambah Validation

### Model-Level Validation

Menggunakan binding tags di model:

```go
type Product struct {
    Name  string  `json:"name" binding:"required,min=3,max=255"`
    Price float64 `json:"price" binding:"required,gt=0"`
    Stock int     `json:"stock" binding:"gte=0"`
    Email string  `json:"email" binding:"required,email"`
}
```

**Validation Tags:**
- `required` - Field wajib diisi
- `min=3` - Minimal 3 karakter
- `max=255` - Maksimal 255 karakter
- `gt=0` - Greater than 0
- `gte=0` - Greater than or equal 0
- `email` - Valid email format

**📖 Validation Docs:** [gin-gonic.com/docs/examples/binding-and-validation](https://gin-gonic.com/docs/examples/binding-and-validation/)

---

## 4. Menambah Relationships

Contoh: Product belongs to Category.

### Define Models

```go
// Category model
type Category struct {
    ID       uint      `gorm:"primarykey" json:"id"`
    Name     string    `gorm:"size:100;not null" json:"name"`
    Products []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

// Product model
type Product struct {
    ID          uint     `gorm:"primarykey" json:"id"`
    Name        string   `gorm:"size:255;not null" json:"name"`
    CategoryID  uint     `json:"category_id"`
    Category    Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
    // ... other fields
}
```

### Preload Relationships

```go
// Service method
func (s *ProductService) GetWithCategory(id uint) (*models.Product, error) {
    var product models.Product
    result := database.GetDB().Preload("Category").First(&product, id)
    return &product, result.Error
}

// Get all with categories
func (s *ProductService) GetAllWithCategories() ([]models.Product, error) {
    var products []models.Product
    result := database.GetDB().Preload("Category").Find(&products)
    return products, result.Error
}
```

**📖 GORM Associations:** [gorm.io/docs/belongs_to.html](https://gorm.io/docs/belongs_to.html)

---

## 5. Menambah Vue Component

### Create Component

Buat file `frontend/src/components/ProductCard.vue`:

```vue
<script setup>
import { Motion } from '@motionone/vue'

const props = defineProps({
  product: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['edit', 'delete'])
</script>

<template>
  <Motion
    :initial="{ opacity: 0, y: 20 }"
    :animate="{ opacity: 1, y: 0 }"
    :transition="{ duration: 0.3 }"
  >
    <div class="bg-white rounded-lg shadow-md p-6">
      <h3 class="text-xl font-bold mb-2">{{ product.name }}</h3>
      <p class="text-2xl font-bold text-blue-600 mb-4">
        Rp {{ product.price.toLocaleString('id-ID') }}
      </p>
      <p class="text-gray-600 mb-4">Stock: {{ product.stock }}</p>
      
      <div class="flex gap-2">
        <button 
          @click="emit('edit', product)"
          class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
        >
          Edit
        </button>
        <button 
          @click="emit('delete', product.id)"
          class="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600"
        >
          Hapus
        </button>
      </div>
    </div>
  </Motion>
</template>
```

### Create Composable

Buat file `frontend/src/composables/useProducts.js`:

```javascript
import { ref } from 'vue'
import { useApi } from './useApi'

export const useProducts = () => {
  const api = useApi()
  const products = ref([])
  const loading = ref(false)
  const error = ref(null)

  const fetchProducts = async () => {
    loading.value = true
    error.value = null
    try {
      const response = await api.get('/api/products')
      products.value = response.data.data
    } catch (err) {
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  const createProduct = async (productData) => {
    try {
      const response = await api.post('/api/products', productData)
      await fetchProducts()
      return response.data
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  const updateProduct = async (id, productData) => {
    try {
      await api.put(`/api/products/${id}`, productData)
      await fetchProducts()
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  const deleteProduct = async (id) => {
    try {
      await api.delete(`/api/products/${id}`)
      await fetchProducts()
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  return {
    products,
    loading,
    error,
    fetchProducts,
    createProduct,
    updateProduct,
    deleteProduct
  }
}
```

### Use in View

Di `frontend/src/views/Products.vue`:

```vue
<script setup>
import { onMounted } from 'vue'
import ProductCard from '@/components/ProductCard.vue'
import { useProducts } from '@/composables/useProducts'

const { products, loading, fetchProducts, deleteProduct } = useProducts()

onMounted(() => {
  fetchProducts()
})

const handleDelete = async (id) => {
  if (confirm('Yakin ingin menghapus?')) {
    await deleteProduct(id)
  }
}
</script>

<template>
  <div class="container mx-auto p-6">
    <h1 class="text-3xl font-bold mb-6">Products</h1>
    
    <div v-if="loading">Loading...</div>
    
    <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <ProductCard
        v-for="product in products"
        :key="product.id"
        :product="product"
        @delete="handleDelete"
      />
    </div>
  </div>
</template>
```

---

## 6. Official Documentation

### Backend (Go)

**Gin Framework:**
- Main Docs: [gin-gonic.com/docs](https://gin-gonic.com/docs/)
- Examples: [gin-gonic.com/docs/examples](https://gin-gonic.com/docs/examples/)
- GitHub: [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)

**GORM:**
- Main Docs: [gorm.io/docs](https://gorm.io/docs/)
- Models: [gorm.io/docs/models.html](https://gorm.io/docs/models.html)
- CRUD: [gorm.io/docs/create.html](https://gorm.io/docs/create.html)
- Associations: [gorm.io/docs/belongs_to.html](https://gorm.io/docs/belongs_to.html)

### Frontend (Vue)

**Vue 3:**
- Main Docs: [vuejs.org](https://vuejs.org/)
- Composition API: [vuejs.org/guide/introduction.html](https://vuejs.org/guide/introduction.html)
- Components: [vuejs.org/guide/components.html](https://vuejs.org/guide/components.html)

**Vite:**
- Main Docs: [vitejs.dev](https://vitejs.dev/)
- Config: [vitejs.dev/config](https://vitejs.dev/config/)

**Tailwind CSS:**
- Main Docs: [tailwindcss.com/docs](https://tailwindcss.com/docs/)
- Utility Classes: [tailwindcss.com/docs/utility-first](https://tailwindcss.com/docs/utility-first)

---

## 💡 Tips & Best Practices

### Backend
1. ✅ Gunakan Service Pattern untuk separation of concerns
2. ✅ Validate input di handler menggunakan binding tags
3. ✅ Handle errors dengan pesan Bahasa Indonesia
4. ✅ Gunakan soft delete untuk data penting
5. ✅ Add indexes di kolom yang sering di-query

### Frontend
1. ✅ Gunakan Composables untuk reusable logic
2. ✅ Component harus small dan focused
3. ✅ Gunakan Tailwind utilities instead of custom CSS
4. ✅ Add loading states untuk better UX
5. ✅ Handle errors gracefully dengan error messages

---

## 📚 Related Documentation

- **[ARCHITECTURE_EXPLAINED.md](./ARCHITECTURE_EXPLAINED.md)** - Penjelasan lengkap semua package
- **[FOLDER_STRUCTURE.md](./FOLDER_STRUCTURE.md)** - Struktur folder detail
- **[API_DOCUMENTATION.md](./API_DOCUMENTATION.md)** - API reference
- **[TESTING.md](./TESTING.md)** - Testing guide
- **[FAQ.md](./FAQ.md)** - Common questions

---

## 📞 Need Help?

**Developer:** Zulfikar Hidayatullah  
**Phone:** +62 857-1583-8733

---

**Last Updated:** 27 Desember 2025  
**Version:** 1.0.0
