package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image/png"
	"net"
	"sync"
	"time"

	"github.com/cbroglie/mustache"
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
func (rs *ReceiptService) Print(order models.Order, discount float64, service_cost float64, d time.Time, lang_code string, template string) error {
	socket, err := net.Dial("tcp", fmt.Sprintf("%s:9100", rs.Settings.ReceiptPrinter.Host))
	if err != nil {
		return err
	}
	defer socket.Close()

	p := escpos.New(socket)

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
			map[string]interface{}{"name": item.Product.Name, "quantity": item.Quantity, "price": item.SalePrice * item.Quantity},
		)
		subtotal += int(item.SalePrice) * int(item.Quantity)
	}

	total := subtotal - int(discount)

	data := map[string]interface{}{
		"direction":      lang.Orientation,
		"t_date":         lang.Pack["date"],
		"t_name":         lang.Pack["name"],
		"t_quantity":     lang.Pack["quantity"],
		"t_total":        lang.Pack["total"],
		"t_price":        lang.Pack["price"],
		"t_discount":     lang.Pack["discount"],
		"t_subtotal":     lang.Pack["subtotal"],
		"t_service_cost": lang.Pack["service"],
		"order_id":       order.DisplayId,
		"date":           d.Format("2/1/2006 15:04"),
		"order_items":    order_items,
		"discount":       discount,
		"service_cost":   service_cost,
		"total":          total,
		"subtotal":       subtotal,
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

	output, err := mustache.RenderFile(template, data)
	if err != nil {
		return err
	}

	// create context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		// chromedp.WithDebugf(log.Printf),
	)
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

	p.Size(1, 1).PrintImage(img)
	p.LineFeed()

	p.PrintAndCut()

	return nil
}
