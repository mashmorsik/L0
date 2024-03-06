package web

import (
	"bytes"
	log "github.com/mashmorsik/L0/pkg/logger"
	"github.com/mashmorsik/L0/pkg/models"
	"html/template"
)

func DisplayOrder(order models.Order) string {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Order Details</title>
	<style>
		body {
			margin-left: 400px;
		}
        h1 {
            text-align: center;
			margin-left: -600px;
        }
        h2 {
            text-align: center; 
			margin-left: -600px;
        }
       	ul {
 			text-align: center;
		}
		li {
			text-align: left;
		}
    </style>
</head>
<body>
    <h1>Order Details</h1>

    <h2>Order Info</h2>
    <ul>
        <li>Order UID: {{ .OrderUid }}</li>
        <li>Track Number: {{ .TrackNumber }}</li>
        <li>Entry: {{ .Entry }}</li>
        <li>Locale: {{ .Locale }}</li>
        <li>Internal Signature: {{ .InternalSignature }}</li>
        <li>Customer ID: {{ .CustomerId }}</li>
        <li>Delivery Service: {{ .DeliveryService }}</li>
        <li>Shardkey: {{ .Shardkey }}</li>
        <li>SmId: {{ .SmId }}</li>
        <li>Date Created: {{ .DateCreated }}</li>
        <li>OofShard: {{ .OofShard }}</li>
    </ul>

    <h2>Delivery Info</h2>
    <ul>
        <li>Name: {{ .Delivery.Name }}</li>
        <li>Phone: {{ .Delivery.Phone }}</li>
        <li>Zip: {{ .Delivery.Zip }}</li>
        <li>City: {{ .Delivery.City }}</li>
        <li>Address: {{ .Delivery.Address }}</li>
        <li>Region: {{ .Delivery.Region }}</li>
        <li>Email: {{ .Delivery.Email }}</li>
    </ul>

    <h2>Payment Info</h2>
    <ul>
        <li>Transaction: {{ .Payment.Transaction }}</li>
        <li>Request ID: {{ .Payment.RequestId }}</li>
        <li>Currency: {{ .Payment.Currency }}</li>
        <li>Provider: {{ .Payment.Provider }}</li>
        <li>Amount: {{ .Payment.Amount }}</li>
        <li>Payment Date: {{ .Payment.PaymentDt }}</li>
        <li>Bank: {{ .Payment.Bank }}</li>
        <li>Delivery Cost: {{ .Payment.DeliveryCost }}</li>
        <li>Goods Total: {{ .Payment.GoodsTotal }}</li>
        <li>Custom Fee: {{ .Payment.CustomFee }}</li>
    </ul>

    <h2>Items</h2>
    <ul>
        {{ range .Items }}
        <li>
            <ul>
                <li>ChrtId: {{ .ChrtId }}</li>
                <li>Track Number: {{ .TrackNumber }}</li>
                <li>Price: {{ .Price }}</li>
                <li>Rid: {{ .Rid }}</li>
                <li>Name: {{ .Name }}</li>
                <li>Sale: {{ .Sale }}</li>
                <li>Size: {{ .Size }}</li>
                <li>Count: {{ .Count }}</li>
                <li>Total Price: {{ .TotalPrice }}</li>
                <li>NmId: {{ .NmId }}</li>
                <li>Brand: {{ .Brand }}</li>
                <li>Status: {{ .Status }}</li>
            </ul>
        </li>
        {{ end }}
    </ul>
</body>
</html>
`

	tmpl, err := template.New("orderDetails").Parse(html)
	if err != nil {
		log.Errf("can't parse html into template, err: %s", err)
		return ""
	}

	var tplBuf bytes.Buffer
	if err = tmpl.Execute(&tplBuf, order); err != nil {
		log.Errf("failed to execute HTML template: %s", err)
		return ""
	}

	return tplBuf.String()
}
