package repos

import (
	"database/sql"
	"fmt"
	"time"
)

type DashboardRepo struct {
	db *sql.DB
}

func NewDashboardRepo(db *sql.DB) *DashboardRepo {
	return &DashboardRepo{db: db}
}

type DashboardData struct {
	TotalSales      float64              `json:"total_sales"`
	TotalProducts   int                  `json:"total_products"`
	TotalCustomers  int                  `json:"total_customers"`
	TotalDebts      float64              `json:"total_debts"`
	TotalOrdersToday int                 `json:"total_orders_today"`
	GrowthPercent   float64              `json:"growth_percent"`
	LeadingRegion   RegionStat           `json:"leading_region"`
	LeastRegion     RegionStat           `json:"least_region"`
	MonthlySales    []MonthlySale        `json:"monthly_sales"`
	RegionBreakdown []RegionBreakdownItem `json:"region_breakdown"`
	RecentActivity  []ActivityItem       `json:"recent_activity"`
	TopProducts     []TopProductItem     `json:"top_products"`
}

type RegionStat struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

type MonthlySale struct {
	Month string  `json:"month"`
	Value float64 `json:"value"`
}

type RegionBreakdownItem struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Color string  `json:"color"`
}

type ActivityItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
	Who  string `json:"who"`
	Time string `json:"time"`
	Icon string `json:"icon"`
	Color string `json:"color"`
}

type TopProductItem struct {
	Name    string  `json:"name"`
	Sold    float64 `json:"sold"`
	Revenue float64 `json:"revenue"`
}

func (r *DashboardRepo) GetDashboard(storeID int64) (*DashboardData, error) {
	d := &DashboardData{}

	// 1. Total sales (sum of amount_billed where type=sale)
	r.db.QueryRow(`
		SELECT COALESCE(SUM(amount_billed), 0) FROM bills
		WHERE type = 'sale' AND deleted_at IS NULL
		AND ($1 = 0 OR store_id = $1)
	`, storeID).Scan(&d.TotalSales)

	// 2. Total products
	r.db.QueryRow(`
		SELECT COUNT(*) FROM items WHERE deleted_at IS NULL
	`).Scan(&d.TotalProducts)

	// 3. Total customers (people)
	r.db.QueryRow(`
		SELECT COUNT(*) FROM people WHERE deleted_at IS NULL
	`).Scan(&d.TotalCustomers)

	// 4. Total outstanding debts (billed - paid where not fully paid)
	r.db.QueryRow(`
		SELECT COALESCE(SUM(amount_billed - amount_paid), 0) FROM bills
		WHERE status != 'paid' AND deleted_at IS NULL
		AND ($1 = 0 OR store_id = $1)
	`, storeID).Scan(&d.TotalDebts)

	// 5. Orders today
	r.db.QueryRow(`
		SELECT COUNT(*) FROM bills
		WHERE DATE(created_at) = CURRENT_DATE AND deleted_at IS NULL
		AND ($1 = 0 OR store_id = $1)
	`, storeID).Scan(&d.TotalOrdersToday)

	// 6. Growth percent (this month vs last month)
	var thisMonth, lastMonth float64
	r.db.QueryRow(`
		SELECT COALESCE(SUM(amount_billed), 0) FROM bills
		WHERE type = 'sale' AND deleted_at IS NULL
		AND DATE_TRUNC('month', created_at) = DATE_TRUNC('month', CURRENT_DATE)
		AND ($1 = 0 OR store_id = $1)
	`, storeID).Scan(&thisMonth)

	r.db.QueryRow(`
		SELECT COALESCE(SUM(amount_billed), 0) FROM bills
		WHERE type = 'sale' AND deleted_at IS NULL
		AND DATE_TRUNC('month', created_at) = DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 month')
		AND ($1 = 0 OR store_id = $1)
	`, storeID).Scan(&lastMonth)

	if lastMonth > 0 {
		d.GrowthPercent = ((thisMonth - lastMonth) / lastMonth) * 100
	}

	// 7. Monthly sales (last 12 months)
	monthRows, err := r.db.Query(`
		SELECT TO_CHAR(DATE_TRUNC('month', created_at), 'Mon') as month,
		       COALESCE(SUM(amount_billed), 0) as total
		FROM bills
		WHERE type = 'sale' AND deleted_at IS NULL
		AND created_at >= CURRENT_DATE - INTERVAL '12 months'
		AND ($1 = 0 OR store_id = $1)
		GROUP BY DATE_TRUNC('month', created_at)
		ORDER BY DATE_TRUNC('month', created_at) ASC
	`, storeID)
	if err == nil {
		defer monthRows.Close()
		for monthRows.Next() {
			var ms MonthlySale
			monthRows.Scan(&ms.Month, &ms.Value)
			d.MonthlySales = append(d.MonthlySales, ms)
		}
	}
	if d.MonthlySales == nil {
		d.MonthlySales = []MonthlySale{}
	}

	// 8. Region breakdown (sales by store location/region)
	colors := []string{"#2b2d6e", "#3d5a80", "#5b8fb9", "#8ecae6", "#b8d8e8"}
	regionRows, err := r.db.Query(`
		SELECT COALESCE(s.name, 'Unknown') as region, COALESCE(SUM(b.amount_billed), 0) as total
		FROM bills b
		LEFT JOIN stores s ON s.id = b.store_id
		WHERE b.type = 'sale' AND b.deleted_at IS NULL
		AND ($1 = 0 OR b.store_id = $1)
		GROUP BY s.name
		ORDER BY total DESC
		LIMIT 5
	`, storeID)
	if err == nil {
		defer regionRows.Close()
		i := 0
		for regionRows.Next() {
			var rb RegionBreakdownItem
			regionRows.Scan(&rb.Name, &rb.Value)
			if i < len(colors) {
				rb.Color = colors[i]
			} else {
				rb.Color = "#b8d8e8"
			}
			d.RegionBreakdown = append(d.RegionBreakdown, rb)
			i++
		}
	}
	if d.RegionBreakdown == nil {
		d.RegionBreakdown = []RegionBreakdownItem{}
	}

	// 9. Leading & least region
	if len(d.RegionBreakdown) > 0 {
		d.LeadingRegion = RegionStat{Name: d.RegionBreakdown[0].Name, Amount: d.RegionBreakdown[0].Value}
		d.LeastRegion = RegionStat{
			Name:   d.RegionBreakdown[len(d.RegionBreakdown)-1].Name,
			Amount: d.RegionBreakdown[len(d.RegionBreakdown)-1].Value,
		}
	}

	// 10. Recent activity (last 10 bills + payments)
	actRows, err := r.db.Query(`
		(
			SELECT 'sale' as type,
			       'New sale — TZS ' || TRIM(TO_CHAR(b.amount_billed, '999,999,999')) as text,
			       COALESCE(u.name, '') as who,
			       b.created_at as ts
			FROM bills b
			LEFT JOIN users u ON u.id = b.created_by
			WHERE b.type = 'sale' AND b.deleted_at IS NULL
			ORDER BY b.created_at DESC LIMIT 5
		)
		UNION ALL
		(
			SELECT 'payment' as type,
			       'Payment — TZS ' || TRIM(TO_CHAR(py.payment_amount, '999,999,999')) as text,
			       COALESCE(u.name, '') as who,
			       py.created_at as ts
			FROM payments py
			LEFT JOIN users u ON u.id = py.created_by
			ORDER BY py.created_at DESC LIMIT 5
		)
		ORDER BY ts DESC LIMIT 10
	`)
	if err == nil {
		defer actRows.Close()
		now := time.Now()
		for actRows.Next() {
			var actType, text, who string
			var ts time.Time
			actRows.Scan(&actType, &text, &who, &ts)

			icon := "fa-cart-shopping"
			color := "#2d6a4f"
			if actType == "payment" {
				icon = "fa-money-bill-wave"
				color = "#2b2d6e"
			}

			d.RecentActivity = append(d.RecentActivity, ActivityItem{
				Type:  actType,
				Text:  text,
				Who:   who,
				Time:  timeAgo(now, ts),
				Icon:  icon,
				Color: color,
			})
		}
	}
	if d.RecentActivity == nil {
		d.RecentActivity = []ActivityItem{}
	}

	// 11. Top products (by revenue from bill_items)
	prodRows, err := r.db.Query(`
		SELECT i.name,
		       COALESCE(SUM(bi.item_quantity), 0) as sold,
		       COALESCE(SUM(bi.item_price * bi.item_quantity - bi.discount_amount), 0) as revenue
		FROM bill_items bi
		INNER JOIN items i ON i.id = bi.item_id
		INNER JOIN bills b ON b.id = bi.bill_id
		WHERE b.type = 'sale' AND b.deleted_at IS NULL
		GROUP BY i.name
		ORDER BY revenue DESC
		LIMIT 5
	`)
	if err == nil {
		defer prodRows.Close()
		for prodRows.Next() {
			var tp TopProductItem
			prodRows.Scan(&tp.Name, &tp.Sold, &tp.Revenue)
			d.TopProducts = append(d.TopProducts, tp)
		}
	}
	if d.TopProducts == nil {
		d.TopProducts = []TopProductItem{}
	}

	return d, nil
}

// timeAgo returns human-readable time difference
func timeAgo(now, t time.Time) string {
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		mins := int(diff.Minutes())
		if mins == 1 {
			return "1 min ago"
		}
		return fmt.Sprintf("%d min ago", mins)
	case diff < 24*time.Hour:
		hrs := int(diff.Hours())
		if hrs == 1 {
			return "1 hr ago"
		}
		return fmt.Sprintf("%d hrs ago", hrs)
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "yesterday"
		}
		return fmt.Sprintf("%d days ago", days)
	default:
		return t.Format("Jan 2")
	}
}