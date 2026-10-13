package models

import (
	"database/sql"
	"time"
)

type Category struct {
	ID               int64         `json:"id"`
	CategoryTypeID   int64         `json:"categorytype_id"`
	Name             string        `json:"name"`
	Description      string        `json:"description,omitempty"`
	Type             string        `json:"type"`
	CreatedBy        int64         `json:"created_by"`
	CreatedAt        time.Time     `json:"created_at"`
	Created_Username string        `json:"created_username"`
	DeletedAt        sql.NullTime  `json:"deleted_at,omitempty"`
	DeletedBy        sql.NullInt64 `json:"deleted_by,omitempty"`
}

type CategoryType struct {
	ID               int64         `json:"id"`
	Name             string        `json:"name"`
	Description      string        `json:"description,omitempty"`
	CreatedBy        int64         `json:"created_by"`
	CreatedAt        time.Time     `json:"created_at"`
	Created_Username string        `json:"created_username"`
	DeletedAt        sql.NullTime  `json:"deleted_at,omitempty"`
	DeletedBy        sql.NullInt64 `json:"deleted_by,omitempty"`
}

type People struct {
	ID               int64      `json:"id"`
	Type             string     `json:"type"`
	Category         string     `json:"category"`
	Tin              string     `json:"tin"`
	Vrn              string     `json:"vrn"`
	Nin              string     `json:"nin"`
	Email            string     `json:"email"`
	Location         string     `json:"location"`
	Name             string     `json:"name"`
	Phone            string     `json:"phone"`
	Status           bool       `json:"status"`
	LogoURL          string     `json:"logo_url"`
	PoBox            string     `json:"po_box"`
	Tags             *string    `json:"tags"`
	Created_Username string     `json:"created_username"`
	CreatedAt        time.Time  `json:"created_at"`
	CreatedBy        int64      `json:"created_by"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
	DeletedBy        *int64     `json:"deleted_by,omitempty"`
}

type Setting struct {
	ID        int64      `json:"id"`
	Code      string     `json:"code"`
	Name      string     `json:"name"`
	Value     string     `json:"value,omitempty"`
	Category  string     `json:"category"`
	CreatedAt time.Time  `json:"createdAt"`
	CreatedBy *int64     `json:"createdBy"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
	DeletedBy *int64     `json:"deletedBy,omitempty"`
}

type User struct {
	ID        int64         `json:"id"`
	Email     string        `json:"email"`
	Password  string        `json:"password"`
	Name      string        `json:"name"`
	Phone     string        `json:"phone"`
	Status    bool          `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	DeletedAt sql.NullTime  `json:"deleted_at,omitempty"`
	DeletedBy sql.NullInt64 `json:"deleted_by,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ItemStore struct {
	ID        int64     `json:"id"`
	ItemID    int64     `json:"item_id"`
	StoreID   int64     `json:"store_id"`
	UnitPrice float64   `json:"unit_price"`
	Tags      string    `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy int64     `json:"createdBy"`
	DeletedAt time.Time `json:"deletedAt,omitempty"`
	DeletedBy int64     `json:"deletedBy,omitempty"`
}

type ItemOrder struct {
	ID        int64     `json:"id"`
	ItemID    int64     `json:"item_id"`
	StoreID   int64     `json:"store_id"`
	UnitPrice float64   `json:"unit_price"`
	Tags      string    `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy int64     `json:"createdBy"`
	DeletedAt time.Time `json:"deletedAt,omitempty"`
	DeletedBy int64     `json:"deletedBy,omitempty"`
}

type Expense struct {
	ID        int64     `json:"id"`
	ItemID    int64     `json:"item_id"`
	StoreID   int64     `json:"store_id"`
	UnitPrice float64   `json:"unit_price"`
	Tags      string    `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy int64     `json:"createdBy"`
	DeletedAt time.Time `json:"deletedAt,omitempty"`
	DeletedBy int64     `json:"deletedBy,omitempty"`
}

type Continent struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Status bool   `json:"status"`
}

type Country struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Status bool   `json:"status"`
}

type Region struct {
	ID        int64     `json:"id"`
	Country   string    `json:"country"`
	Zone      string    `json:"zone"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy int64     `json:"createdBy"`
	DeletedAt time.Time `json:"deletedAt,omitempty"`
	DeletedBy int64     `json:"deletedBy,omitempty"`
}

type District struct {
	ID         int64  `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	RegionID   int    `db:"region_id" json:"region_id"`
	RegionName string `json:"region_name"`
	Status     bool   `db:"status" json:"status"`
}

type Village struct {
	ID           int64  `db:"id" json:"id"`
	Name         string `db:"name" json:"name"`
	DistrictID   int    `db:"district_id" json:"district_id"`
	DistrictName string `json:"district_name"`
	Status       bool   `db:"status" json:"status"`
}

type Item struct {
	ID       int64    `db:"id" json:"id"`
	Name     string   `db:"name" json:"name"`
	Type     string   `db:"type" json:"type"`
	Category string   `db:"category" json:"category"`
	UnitCost float64  `db:"unit_cost" json:"unit_cost"`
	UOM      string   `db:"uom" json:"uom"`
	Status   string   `db:"status" json:"status"`
	Tags     []string `db:"tags" json:"tags"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	CreatedBy int64     `db:"created_by" json:"created_by"`

	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
	DeletedBy *int64     `db:"deleted_by" json:"deleted_by,omitempty"`

	Created_Username string `json:"created_username,omitempty"`
	Stock            string `json:"stock"`
}

type Store struct {
	ID       int64     `db:"id" json:"id"`
	Name     string    `db:"name" json:"name"`
	Type     string    `db:"type" json:"type"`
	Category string    `db:"category" json:"category,omitempty"`
	Location string    `db:"location" json:"location,omitempty"`
	Owner    int64     `db:"owner" json:"owner,omitempty"`
	Status   string    `db:"status" json:"status"`
	Tags     *[]string `db:"tags" json:"tags"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	CreatedBy int64     `db:"created_by" json:"created_by"`

	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
	DeletedBy *int64     `db:"deleted_by" json:"deleted_by,omitempty"`

	Created_Username string `json:"created_username,omitempty"`
}

type ItemStock struct {
	ID       int64   `db:"id" json:"id"`
	ItemID   int64   `db:"item_id" json:"item_id"`
	StoreID  int64   `db:"store_id" json:"store_id"`
	Batch    *string `db:"batch" json:"batch,omitempty"`
	Quantity float64 `db:"quantity" json:"quantity"`

	SellingPrice   float64 `db:"selling_price" json:"selling_price"`
	DiscountAmount float64 `db:"discount_amount" json:"discount_amount"`
	DiscountRate   float64 `db:"discount_rate" json:"discount_rate"`

	ExpiryDate *time.Time `db:"expiry_date" json:"expiry_date,omitempty"`
}

type Transaction struct {
	ID      int64 `db:"id" json:"id"`
	ItemID  int64 `db:"item_id" json:"item_id"`
	StoreID int64 `db:"store_id" json:"store_id"`

	Type     string `db:"type" json:"type"` // IN, OUT
	Category string `db:"category" json:"category,omitempty"`

	Quantity float64 `db:"quantity" json:"quantity"`
	Status   string  `db:"status" json:"status"`

	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	CreatedBy        int64     `db:"created_by" json:"created_by"`
	ItemName         string    `json:"item_name"`
	Created_Username string    `json:"created_username"`
	StoreName        string    `json:"store_name"`
}

type AuditLog struct {
	ID          int64     `db:"id" json:"id"`
	Actor       int64     `db:"actor" json:"actor"`
	Action      string    `db:"action" json:"action"`
	Entity      string    `db:"entity" json:"entity"`
	EntityID    *int64    `db:"entity_id" json:"entity_id,omitempty"`
	Description *string   `db:"description" json:"description,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type Bill struct {
	ID           int64      `json:"id"`
	Code         string     `json:"code"`
	PeopleID     int64      `json:"people_id"`
	PeopleName   string     `json:"people_name,omitempty"`
	StoreID      int64      `json:"store_id"`
	StoreName    string     `json:"store_name,omitempty"`
	Type         string     `json:"type"`
	Status       string     `json:"status"`
	AmountBilled float64    `json:"amount_billed"`
	AmountPaid   float64    `json:"amount_paid"`
	CreatedAt    time.Time  `json:"created_at"`
	CreatedBy    int64      `json:"created_by"`
	CreatedName  string     `json:"created_name,omitempty"`
	Items        []BillItem `json:"items,omitempty"`
	Payments     []Payment  `json:"payments,omitempty"`
}

type BillItem struct {
	ID             int64   `json:"id"`
	BillID         int64   `json:"bill_id"`
	ItemID         int64   `json:"item_id"`
	ItemName       string  `json:"item_name,omitempty"`
	ItemPrice      float64 `json:"item_price"`
	ItemQuantity   float64 `json:"item_quantity"`
	DiscountAmount float64 `json:"discount_amount"`
	DiscountRate   float64 `json:"discount_rate"`
	LineTotal      float64 `json:"line_total,omitempty"`
}

type Payment struct {
	ID            int64     `json:"id"`
	BillID        int64     `json:"bill_id"`
	PeopleID      int64     `json:"people_id"`
	PeopleName    string    `json:"people_name,omitempty"`
	Type          string    `json:"type"`
	Category      string    `json:"category"`
	PaymentAmount float64   `json:"payment_amount"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     int64     `json:"created_by"`
	CreatedName   string    `json:"created_name,omitempty"`
}

type SMS struct {
	ID       int64   `db:"id" json:"id"`
	Content  string  `db:"content" json:"content"`
	Sender   *string `db:"sender" json:"sender,omitempty"`
	Receiver *string `db:"receiver" json:"receiver,omitempty"`
	Status   string  `db:"status" json:"status"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	CreatedBy int64     `db:"created_by" json:"created_by"`
}

type SMSTemplate struct {
	ID      int64  `db:"id" json:"id"`
	Name    string `db:"name" json:"name"`
	Content string `db:"content" json:"content"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	CreatedBy int64     `db:"created_by" json:"created_by"`

	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
	DeletedBy *int64     `db:"deleted_by" json:"deleted_by,omitempty"`
}

type BulkError struct {
	Row   int    `json:"row"`
	Error string `json:"error"`
}

type BulkResult struct {
	Success int         `json:"success"`
	Failed  int         `json:"failed"`
	Errors  []BulkError `json:"errors"`
}

type Role struct {
	ID          int            `db:"id" json:"id"`
	Name        string         `db:"name" json:"name"`
	Permissions *string `db:"permissions" json:"permissions"`

	Created_Username string `json:"created_username"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	CreatedBy int64     `db:"created_by" json:"created_by"`

	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
	DeletedBy *int64     `db:"deleted_by" json:"deleted_by,omitempty"`
}
