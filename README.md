<div dir="rtl">

# ⚡ پکیج ValidX — Fast, Type-Safe Validation for Go
 پکیج ValidX یک Validation نسل جدید برای Go است که با تمرکز بر **Performance، Type-Safety و Developer Experience (DX)** طراحی شده و هدف آن ارائه جایگزینی حرفه‌ای برای رویکردهای سنتی (مانند tag-based validation) است.

---

## 🚀 چرا ValidX؟

اگر با پکیج‌هایی مثل `go-playground/validator` کار کرده باشی، حتماً با این مشکلات مواجه شدی:

❌ وابستگی شدید به struct tag  
❌ عدم type-safety  
❌ سختی در debug  
❌ عدم composability واقعی  
❌ DX ضعیف در پروژه‌های بزرگ  

---

### ✨ پکیج ValidX چه چیزی را حل می‌کند؟

✅ کاملاً Type-Safe (بدون reflection در API سطح بالا)  
✅ Fluent API مدرن و خوانا  
✅ بدون tag (code-first validation)  
✅ Performance بالا (zero alloc در مسیر hot)  
✅ قابل توسعه و modular  
✅ مناسب microservices و domain-driven design  

---
## 🎯 فلسفه طراحی (Design Philosophy)

ValidX بر پایه چند اصل کلیدی ساخته شده:

- **Explicit over Implicit**  
  هیچ magic یا string parsing وجود ندارد — تمام validationها به‌صورت صریح در کد تعریف می‌شوند.

- **Type Safety First**  
  استفاده از generics و getter functions برای حذف خطاهای runtime و افزایش اطمینان در refactor.

- **Performance-Oriented**  
  حذف reflection از hot path و کاهش allocation برای سیستم‌های high-throughput.

- **Composable Architecture**  
  ruleها، schemaها و validation logic قابل ترکیب، reuse و توسعه هستند.

---
## 🚀 نقاط قوت کلیدی (Key Strengths)

### ⚡ Performance واقعی
- بدون استفاده از reflection در مسیرهای critical
- قابلیت caching و reuse schema
- مناسب برای microservices و high-load APIs


### 🧠 Type Safety و Refactor Safety
- تمام field accessها از طریق getter function
- مقاوم در برابر تغییر نام field یا type
- حذف کامل string-based validation


### 🧩 Composability & Extensibility
- تعریف ruleهای reusable به‌صورت function
- امکان ساخت ruleهای domain-specific
- بدون global state یا registration پیچیده

### 🧑‍💻 Developer Experience (DX)
- Fluent API خوانا و قابل فهم
- auto-complete کامل در IDE
- self-documenting validation logic

### 🐛 Debuggability
- مسیر validation کاملاً قابل trace
- errorها structured و قابل توسعه
- مناسب برای logging و observability

---
## 📦 نصب

```bash
go get github.com/Skryldev/valid
```
---
# ⚡ شروع سریع
### 1️⃣ تعریف struct

<div dir="ltr">

```go
type RegisterInput struct {
    Email    string
    Password string
    Age      int
}
```

<div dir="rtl">

## 2️⃣ تعریف schema


<div dir="ltr">

```go
func RegisterSchema() *validx.Schema[RegisterInput] {
    s := validx.New[RegisterInput]().Schema()

    validx.Field(s, "email", func(i *RegisterInput) string { return i.Email }).
        Use(validx.RequiredString()).
        Use(validx.Email())

    validx.Field(s, "password", func(i *RegisterInput) string { return i.Password }).
        Use(validx.MinLen(8))

    validx.Field(s, "age", func(i *RegisterInput) int { return i.Age }).
        Use(validx.Between(18, 70))

    return s
}
```

<div dir="rtl">

## 3️⃣ استفاده و اجرا

<div dir="ltr">

```go
schema := RegisterSchema()
err := schema.Validate(&RegisterInput{
    Email:    "bad-email",
    Password: "123",
    Age:      15,
})

if err != nil {
    for _, e := range err.All() {
        fmt.Printf("field=%s code=%s message=%s\n", e.Field, e.Code, e.Message)
    }
}
```

<div dir="rtl">

## 🔍 ساختار خطا


<div dir="ltr">

```go
type Error struct {
    Field   string
    Code    string
    Message string
}

if err := schema.Validate(&input); err != nil {
    for _, e := range err.All() {
        // e.Field, e.Code, e.Message
    }
}
```

<div dir="rtl">

## 🧱 Custom Rule

<div dir="ltr">

```go
func CorporateEmail() validx.Rule[string] {
    rx := regexp.MustCompile(`@company\.com$`)

    return func(v string) (bool, *validx.RuleError) {
        if !rx.MatchString(v) {
            return false, &validx.RuleError{
                Code:    "corporate_email",
                Message: "email must end with @company.com",
            }
        }
        return true, nil
    }
}
```

<div dir="rtl">

#### بررسی مقادیر String در Role:

<div dir="ltr">

```go
func AllowedRoles() validx.Rule[[]string] {
    allowed := map[string]struct{}{
        "admin": {},
        "editor": {},
        "user": {},
    }

    return func(v []string) (bool, *validx.RuleError) {
        for _, role := range v {
            if _, ok := allowed[role]; !ok {
                return false, &validx.RuleError{
                    Code:    "role_invalid",
                    Message: "unsupported role: " + role,
                }
            }
        }
        return true, nil
    }
}
```

<div dir="rtl">

---
## چند Rule های آماده

### 1) Rule های رشته

- `RequiredString()`
- `MinLen(min int)`
- `MaxLen(max int)`
- `Email()`
- `URL()`
- `Pattern(rx *regexp.Regexp)`
- `OneOf(options ...string)`

### 2) Rule های عددی

- `Min[T Number](min T)`
- `Max[T Number](max T)`
- `Between[T Number](min, max T)`

### 3) Rule های اسلایس

- `MinItems[T any](min int)`
- `MaxItems[T any](max int)`
- `Unique[T comparable]()`

### 4) Rule های ترکیبی

- `All[T any](rules ...Rule[T])`
- `Any[T any](rules ...Rule[T])`
- `Optional[T comparable](rule Rule[T])`

---
## 6) ترکیب حرفه‌ای Ruleها

### `All`

##### همه Ruleها باید پاس شوند. در اولین خطا متوقف می‌شود.

<div dir="ltr">

```go
passwordRule := validx.All(
    validx.MinLen(8),
    func(v string) (bool, *validx.RuleError) {
        if !regexp.MustCompile(`[0-9]`).MatchString(v) {
            return false, &validx.RuleError{Code: "password_digit", Message: "must contain digit"}
        }
        return true, nil
    },
)
```

<div dir="rtl">

### `Any`

##### اگر حداقل یک Rule پاس شود، کل Rule پاس است.

<div dir="ltr">

```go
usernameRule := validx.Any(
    validx.Pattern(regexp.MustCompile(`^[a-z0-9_]+$`)),
    validx.Pattern(regexp.MustCompile(`^[A-Z0-9_]+$`)),
)
```

<div dir="rtl">

### `Optional`

##### اگر مقدار zero-value باشد، validation رد نمی‌شود.

<div dir="ltr">

```go
validx.Field(s, "website", func(i *RegisterInput) string { return i.Website }).
    Use(validx.Optional(validx.URL()))
```
<div dir="rtl">

---
## 🆚 مقایسه با validator
| ویژگی | ValidX 🚀 | validator (go-playground) | توضیح |
|------|----------|---------------------------|------------------------------|
| Type Safety | ✅ Compile-time safe (در سطح API) | ❌ Runtime (reflection-based) | در validator تمام validationها از طریق reflection انجام می‌شود، بنابراین type mismatchها فقط در runtime مشخص می‌شوند. در ValidX، getterها strongly-typed هستند و خطاها زودتر detect می‌شوند، که در codebaseهای بزرگ باعث کاهش bugهای پنهان می‌شود. |
| Validation Model | Code-first (Fluent API) | Tag-based (struct tag) | مدل tag-based برای پروژه‌های کوچک سریع است اما در scale بالا باعث implicit logic می‌شود. ValidX validation را explicit می‌کند، که برای maintainability و readability حیاتی است. |
| DX (Developer Experience) | 🔥 Fluent, composable, discoverable | 😐 محدود به string tag | در ValidX از method chaining استفاده می‌شود که IDE-friendly است (auto-complete). در validator، ruleها داخل string هستند و tooling پشتیبانی ضعیف‌تری دارد. |
| Performance | ⚡ Zero/Low alloc, no reflection in hot path | ⚡ optimized but reflection-heavy | validator بسیار optimize شده است اما همچنان وابسته به reflection است. ValidX می‌تواند با حذف reflection در hot path latency را در high-throughput systems کاهش دهد. |
| Debugging | ✅ Explicit & traceable | ❌ implicit & opaque | در validator پیدا کردن اینکه کدام rule fail شده (خصوصاً در tagهای پیچیده) سخت است. ValidX به دلیل explicit بودن rule chain، debugging را ساده و deterministic می‌کند. |
| Composability | ✅ First-class (Rule reuse, chaining) | ❌ محدود | در ValidX می‌توان ruleها را به صورت reusable component تعریف کرد. در validator reuse معمولاً از طریق tag string تکراری انجام می‌شود که DRY را نقض می‌کند. |
| Extensibility | ✅ Native (custom rule به صورت function) | ⚠️ نیازمند registration | در validator برای rule جدید باید آن را register کنی و global state تغییر می‌کند. ValidX اجازه می‌دهد ruleها local و composable باشند. |
| Readability | ✅ High (self-documenting code) | ❌ Low (string parsing needed) | کد ValidX مانند business logic خوانده می‌شود. در validator باید string tag parse شود که cognitive load را بالا می‌برد. |
| Refactor Safety | ✅ Safe (rename-safe) | ❌ unsafe (string-based) | در validator اگر نام field تغییر کند، tagها break می‌شوند بدون compile error. ValidX با getter function این مشکل را حذف می‌کند. |
| IDE Support | ✅ کامل (auto-complete, navigation) | ❌ محدود | IDE نمی‌تواند داخل string tag را analyze کند، اما در ValidX تمام ruleها typed هستند. |
| Error Modeling | ✅ Structured & extensible | ⚠️ محدود | ValidX می‌تواند errorها را به صورت domain-aware مدل کند (مثلاً error code, metadata). validator معمولاً string-based error دارد. |
| Learning Curve | ⚠️ متوسط (نیاز به درک fluent API) | ✅ پایین | validator برای شروع سریع ساده‌تر است، اما ValidX برای long-term scalability طراحی شده است. |
| Use Case Fit | 🟢 Large-scale, DDD, microservices | 🟡 CRUD apps, simple APIs | ValidX برای سیستم‌های پیچیده طراحی شده، در حالی که validator برای use-caseهای ساده کاملاً کافی است. |
---
## 🔥 مثال کامل

فایل `example/main.go` یک سناریوی کامل دارد که شامل:

- Rule آماده + Rule سفارشی
- فیلدهای `string`, `int`, `[]string`
- حالت `invalid` و `valid`
- چاپ خروجی خطاها

اجرا:

<div dir="ltr">

```bash
cd example
go run .
```

<div dir="rtl">

---

## 🧪 سناریوهای مناسب استفاده (Use Cases)

ValidX برای این موارد ایده‌آل است:

- 🟢 سیستم‌های Enterprise با domain logic پیچیده
- 🟢 microservices با latency حساس
- 🟢 پروژه‌هایی با refactor مداوم
- 🟢 APIهایی با validation پیچیده و nested
- 🟢 تیم‌هایی که clean architecture و maintainability مهم است
