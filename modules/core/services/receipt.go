package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image/png"
	"strconv"
	"sync"
	"time"

	"github.com/aymerick/raymond"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/elmawardy/escpos"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/core/models"
)

type ReceiptService struct {
	Config   config.Config
	Settings models.Settings
	Logger   logger.ILogger
}

// Print is used to print a 80mm receipt
func (rs *ReceiptService) Print(order models.Order, discount float64, service_cost float64, d time.Time, lang_code string, template_path string, printer models.PrinterSettings, shop_mode string) error {

	lang_svc := LanguageService{
		Config:   rs.Config,
		Settings: rs.Settings,
		Logger:   rs.Logger,
	}

	lang, err := lang_svc.GetLanguage(lang_code)
	if err != nil {
		return err
	}

	order_items := make([]map[string]interface{}, len(order.Items))
	subtotal := 0

	for _, item := range order.Items {
		order_items = append(order_items,
			map[string]interface{}{"name": item.Product.Name, "quantity": item.Quantity, "price": item.SalePrice},
		)
		subtotal += int(item.SalePrice)
	}

	total := subtotal - int(discount)

	receipt_payments := make([]map[string]interface{}, 0, len(order.Payments))
	for _, payment := range order.Payments {
		receipt_payments = append(receipt_payments, map[string]interface{}{
			"source": payment.Source,
			"amount": fmt.Sprintf("%.2f", payment.Amount),
		})
	}

	custom_data := []struct {
		Key   string
		Value string
	}{}

	for k, v := range order.CustomData {
		custom_data = append(custom_data, struct {
			Key   string
			Value string
		}{
			Key:   k,
			Value: v,
		})
	}

	data := map[string]interface{}{
		"direction":       lang.Orientation,
		"t_date":          lang.Pack["date"],
		"t_name":          lang.Pack["name"],
		"t_quantity":      lang.Pack["quantity"],
		"t_total":         lang.Pack["total"],
		"t_price":         lang.Pack["total"],
		"t_discount":      lang.Pack["discount"],
		"t_subtotal":      lang.Pack["subtotal"],
		"t_service_cost":  lang.Pack["service"],
		"order_id":        order.DisplayId,
		"date":            d.Format("2/1/2006 15:04"),
		"order_items":     order_items,
		"discount":        discount,
		"service_cost":    service_cost,
		"total":           total,
		"subtotal":        subtotal,
		"custom_data":     custom_data,
		"has_custom_data": len(custom_data) > 0,
		"payments":        receipt_payments,
		"has_payments":    len(receipt_payments) > 0,
		"t_payment":       lang.Pack["payment"],
		"is_kitchen_mode": shop_mode == "kitchen",
	}

	if order.IsDelivery {
		data["is_delivery"] = true
		data["t_delivery_address"] = lang.Pack["delivery_address"]
		data["t_customer_phone"] = lang.Pack["phone"]
		data["t_customer_name"] = lang.Pack["customer_name"]

		data["delivery_address"] = order.Customer.Address
		data["customer_name"] = order.Customer.Name
		data["customer_phone"] = order.Customer.Phone
	} else {
		data["is_delivery"] = false
	}

	template, err := raymond.ParseFile(template_path)
	if err != nil {
		return err
	}

	template.RegisterHelper("getByKey", func(items []struct {
		Key   string
		Value string
	}, targetKey string) string {
		for _, item := range items {
			if item.Key == targetKey {
				return item.Value
			}
		}
		return "" // Return empty if not found
	})

	template.RegisterHelper("add", func(val1, val2 string) string {
		// Parse the first value, default to 0 if empty/invalid
		f1, _ := strconv.ParseFloat(val1, 64)

		// Parse the second value, default to 0 if empty/invalid
		f2, _ := strconv.ParseFloat(val2, 64)

		// Perform addition and format back to string (2 decimal places)
		return fmt.Sprintf("%.2f", f1+f2)
	})

	output, err := template.Exec(data)
	if err != nil {
		return err
	}

	ctx, cancel, err := newBrowserContext(rs.Config, rs.Settings.Printing)
	if err != nil {
		return err
	}
	defer cancel()

	// Base64 encode the HTML
	b64 := base64.StdEncoding.EncodeToString([]byte(output))
	uri := "data:text/html;base64," + b64
	width := 570 // Assuming 203 DPI (adjust if needed)
	// width := 640

	var wg sync.WaitGroup

	wg.Add(1)

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if ev, ok := ev.(*page.EventLifecycleEvent); ok {
			if ev.Name == "firstMeaningfulPaint" {
				wg.Done()
			}
		}
	})

	// Capture the screenshot
	var buf []byte
	err = chromedp.Run(ctx,
		chromedp.EmulateViewport(int64(width), 0, chromedp.EmulateScale(1.0)),
		chromedp.Navigate(uri),
		chromedp.ActionFunc(func(ctx context.Context) error {
			wg.Wait()
			return nil
		}),
		chromedp.Screenshot("#main-content", &buf, chromedp.ByQuery), // Screenshot a specific element
		// chromedp.FullScreenshot(&buf, 100),
	)

	if err != nil {
		return err
	}

	// if err := os.WriteFile("fullScreenshot.png", buf, 0o644); err != nil {
	// 	log.Fatal(err)
	// }

	img, err := png.Decode(bytes.NewReader(buf))
	if err != nil {
		return err
	}

	conn, err := openPrinter(printer)
	if err != nil {
		return err
	}
	closed := false
	defer func() {
		if !closed {
			_ = conn.Close()
		}
	}()

	p := escpos.New(conn)

	if _, err := p.Size(1, 1).PrintImage(img); err != nil {
		return err
	}
	if _, err := p.LineFeed(); err != nil {
		return err
	}

	if err := p.PrintAndCut(); err != nil {
		return err
	}

	if err := conn.Close(); err != nil {
		return err
	}
	closed = true

	return nil
}

func (rs *ReceiptService) OpenCashDrawer(printer models.PrinterSettings) error {
	conn, err := openPrinter(printer)
	if err != nil {
		return err
	}

	p := escpos.New(conn)
	err = p.OpenCashDrawer()
	if err != nil {
		_ = conn.Close()
		rs.Logger.Error("Failed to open cash drawer:", err)
		return err
	}

	return conn.Close()
}
