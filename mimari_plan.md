# ReserveFlow v1 — Mimari Plan (Final)
## Genişletilebilir RBAC Altyapısı + Oda Bazlı Özel Yetkilendirme (ResourceAdmin)

---

## 1. Mevcut Durum Analizi

### Roller (models/user.go)
```
UserRoleUser          = "user"
UserRoleAdmin         = "admin"
UserRoleResourceAdmin = "resource-admin"   ← tanımlı ama hiç kullanılmıyor
UserRoleManger        = "manger"           ← tanımlı ama hiç kullanılmıyor
```

### Mevcut Yetkilendirme Zinciri
```
İstek → AuthRequired() → RequireAdmin() → Handler
```
`RequireAdmin()` yalnızca global "admin" rolünü kontrol eder.
`resource-admin` rolü için hiçbir middleware veya DB ilişkisi mevcut değil.

### Sorunlar
- `Resource` modelinin sahibi/yöneticisi yok; her admin tüm kaynakları yönetebiliyor.
- Kullanıcı-kaynak arasında "bu kaynağa özel yetki" kavramı yok.
- Yetkilendirme sistemi yalnızca global roller için tasarlanmış, kaynak-bazlı yetki genişlemesine kapalı.

---

## 2. Mimari Vizyon: Genel RBAC Altyapısı

Bu plan yalnızca "Oda Admini" için değil; ileride projedeki **herhangi bir varlık türü** (Resource, Reservation, WorkingHour vb.) için genişletilebilir bir RBAC (Role-Based Access Control) altyapısı kurar.

### Temel Fikir
```
resource_admins tablosu:
  user_id      → Kim yetkili?
  resource_id  → Hangi kaynağa?
```

Middleware, `entity_type` parametresi eklenerek ileride şu hale getirilebilir:
```
RequireEntityAdmin(entityType string) gin.HandlerFunc
  → "resource", "room", "equipment" gibi türlere göre
     farklı tablolara (veya tek bir genel tabloya) sorgu atar.
```

Şimdilik yalnızca `resource_admins` tablosu ile başlıyoruz;
altyapı bu genişlemeye açık tasarlanıyor.

---

## 3. Veritabanı Şeması

### YENİ TABLO: `resource_admins`

Hangi kullanıcının hangi kaynağa (odaya) "resource-admin" yetkisi olduğunu tutar.

```sql
CREATE TABLE resource_admins (
    id          BIGSERIAL    PRIMARY KEY,
    user_id     BIGINT       NOT NULL REFERENCES users(id)     ON DELETE CASCADE,
    resource_id BIGINT       NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

-- Soft-delete ile uyumlu partial unique index:
CREATE UNIQUE INDEX idx_resource_admins_active
  ON resource_admins(user_id, resource_id)
  WHERE deleted_at IS NULL;
```

**GORM Model (models/resource_admin.go — YENİ):**
```go
type ResourceAdmin struct {
    gorm.Model
    UserID     uint     `gorm:"not null" json:"user_id"`
    ResourceID uint     `gorm:"not null" json:"resource_id"`
    User       User     `gorm:"foreignKey:UserID"     json:"user,omitempty"`
    Resource   Resource `gorm:"foreignKey:ResourceID" json:"resource,omitempty"`
}
```

> NOT: Standart GORM `uniqueIndex` soft-delete ile çakışır.
> Partial unique index (WHERE deleted_at IS NULL) raw SQL migration ile uygulanmalıdır.
> AutoMigrate bu index'i oluşturamaz; manuel `db.Exec(...)` gerekir.

---

## 4. Yeni Dosya Yapısı

```
reserveflow-v1/
├── models/
│   ├── user.go              (DEĞİŞMEZ)
│   ├── resource.go          (DEĞİŞMEZ)
│   ├── resource_admin.go    ← YENİ: ResourceAdmin modeli
│   ├── reservation.go       (DEĞİŞMEZ)
│   └── working_hours.go     (DEĞİŞMEZ)
│
├── dao/
│   ├── user.go              (DEĞİŞMEZ)
│   ├── resource.go          (DEĞİŞMEZ)
│   ├── resource_admin.go    ← YENİ: IsResourceAdmin + CRUD
│   ├── reservation.go       (DEĞİŞMEZ)
│   └── working_hours.go     (DEĞİŞMEZ)
│
├── middleware/
│   ├── auth.go              (DEĞİŞMEZ)
│   └── role.go              (GÜNCELLEME: RequireResourceAdmin() eklenir)
│
├── api/
│   ├── auth.go              (DEĞİŞMEZ)
│   ├── resource.go          (GÜNCELLEME: route'lar yeniden düzenlenir)
│   ├── resource_admin.go    ← YENİ: Atama/listeleme endpoint'leri
│   ├── reservation.go       (DEĞİŞMEZ)
│   ├── api_working_hours.go (DEĞİŞMEZ)
│   └── health.go            (DEĞİŞMEZ)
│
└── main.go                  (GÜNCELLEME: AutoMigrate + partial index migration)
```

---

## 5. Middleware Mantığı

### Temel Kural
> Global Admin her zaman erişebilir — DB sorgusu yapmadan erken geçiş.
> Resource-Admin ise sadece atanmış olduğu kaynağa erişebilir — DB sorgusu zorunlu.
> Diğer roller → 403.

### RequireResourceAdmin() Pseudo-kodu

```
func RequireResourceAdmin() gin.HandlerFunc:

    1. c.Get("user_id") → userID (uint)
    2. c.Get("role")    → role (string)

    3. Eğer role == "admin":
         c.Next()
         return   ← DB sorgusuna hiç gitme, erken çıkış

    4. Eğer role == "resource-admin":
         resourceID = c.Param("id") → parse uint, hata varsa 400
         exists = dao.IsResourceAdmin(userID, resourceID)
         Eğer exists → c.Next()
         Değilse    → 403 Forbidden ("RESOURCE_ADMIN_REQUIRED")

    5. Diğer roller → 403 Forbidden ("INSUFFICIENT_ROLE")
```

### dao.IsResourceAdmin(userID, resourceID uint) bool
```go
var ra ResourceAdmin
err := DB.Where(
    "user_id = ? AND resource_id = ? AND deleted_at IS NULL",
    userID, resourceID,
).First(&ra).Error
return err == nil
```

---

## 6. Route Tasarımı

### Güncellenmiş api/resource.go Route'ları

```
[Genel — Auth yok]
GET  /resources        → GetAllResources
GET  /resources/:id    → GetResourceByID

[Sadece global Admin]
POST   /admin/resources            → CreateResource
DELETE /admin/resources/:id        → DeleteResource

[Admin VEYA ResourceAdmin → RequireResourceAdmin middleware]
PATCH  /admin/resources/:id        → UpdateResource
```

### Yeni api/resource_admin.go Route'ları

```
[Sadece global Admin]
GET    /admin/resources/:id/admins            → ListResourceAdmins
POST   /admin/resources/:id/admins            → AssignResourceAdmin
DELETE /admin/resources/:id/admins/:userId    → RemoveResourceAdmin
```

---

## 7. Endpoint Özeti

| Method | Path | Middleware | Açıklama |
|--------|------|------------|----------|
| GET | `/admin/resources/:id/admins` | AuthRequired + RequireAdmin | Kaynağa atanmış resource-admin'leri listele |
| POST | `/admin/resources/:id/admins` | AuthRequired + RequireAdmin | Bir kullanıcıyı kaynağa resource-admin ata |
| DELETE | `/admin/resources/:id/admins/:userId` | AuthRequired + RequireAdmin | Atamayı kaldır (soft-delete) |
| PATCH | `/admin/resources/:id` | AuthRequired + RequireResourceAdmin | Kaynağı güncelle (admin veya resource-admin) |

**AssignResourceAdmin İstek Gövdesi:**
```json
{ "user_id": 5 }
```

---

## 8. AutoMigrate ve Migration (main.go)

```go
// AutoMigrate (tablo oluşturur)
commons.DB.AutoMigrate(
    &models.User{},
    &models.Resource{},
    &models.ResourceAdmin{},   // YENİ
    &models.WorkingHour{},
    &models.Reservation{},
)

// Partial unique index (AutoMigrate bunu yapamaz — manuel)
commons.DB.Exec(`
    CREATE UNIQUE INDEX IF NOT EXISTS idx_resource_admins_active
    ON resource_admins(user_id, resource_id)
    WHERE deleted_at IS NULL
`)
```

---

## 9. İlerideki Genişleme Yolu (RBAC)

Bu altyapı ilerleyen versiyonlarda şu şekilde genişletilebilir:

```go
// Genel entity-admin tablosu (v2 fikri):
type EntityAdmin struct {
    gorm.Model
    UserID     uint   `gorm:"not null"`
    EntityType string `gorm:"not null"` // "resource", "room", "equipment" ...
    EntityID   uint   `gorm:"not null"`
}

// Middleware parametre alır:
func RequireEntityAdmin(entityType string) gin.HandlerFunc { ... }
```

Şimdilik bu genişleme yapılmıyor; ancak middleware ve DAO katmanı
bu geçişe hazır olacak şekilde sade ve ayrıştırılmış tutulur.

---

## 10. Güvenlik Notları

- ResourceAdmin ataması **yalnızca global Admin** yapabilir.
- `RequireResourceAdmin()` içinde `role == "admin"` kontrolü en başta yapılır → DB'ye hiç gitme.
- Soft-delete ile silinen atama, partial index sayesinde yeniden atamaya izin verir.
- `resource_id` URL parametresi middleware içinde parse edilir; geçersiz ID → 400 Bad Request döner.
