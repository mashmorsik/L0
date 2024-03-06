package web

func DefaultDisplay() string {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Order Details</title>
	<style>
        h1 {
            text-align: center;
        }
        h2 {
            text-align: center; 
        }
       	p {
 			text-align: center; 
		}
    </style>
</head>
<body>
    <h1>Hi there!</h1>
	<h2>To see order details, please, send "order_id".<h2>
	<p>Have a nice day :)<p>
</body>
</html>
`

	return html
}
