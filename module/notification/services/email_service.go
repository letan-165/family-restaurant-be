package services

import (
	"fmt"
	"myapp/module/order/models"
	"net/smtp"
	"os"
	"strings"
)

func SendMail(to, subject, body string) error {
	from := os.Getenv("GMAIL_SHOP")
	password := os.Getenv("GMAIL_PASS")

	msg := fmt.Sprintf("From: %s\r\n", from) +
		fmt.Sprintf("To: %s\r\n", to) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" + // Dòng trống ngăn cách header và body
		body

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from,
		[]string{to},
		[]byte(msg),
	)

	if err != nil {
		return fmt.Errorf("lỗi gửi email: %w", err)
	}

	return nil
}

func SendMailBooking(order models.Order) error {
	to := os.Getenv("GMAIL_ADMIN")
	subject := fmt.Sprintf("Đơn hàng mới từ %s", order.Customer.Receiver)
	body := BuildOrderEmailHTML(order)

	if err := SendMail(to, subject, body); err != nil {
		return err
	}

	return nil
}

func BuildOrderEmailHTML(order models.Order) string {
	var sb strings.Builder

	sb.WriteString("<div style='font-family: Arial, sans-serif; padding: 16px;'>")
	sb.WriteString(fmt.Sprintf("<h2 style='color:#2e6c80;'>Xác nhận đơn hàng #%s</h2>", order.ID.Hex()))
	sb.WriteString(fmt.Sprintf("<p><b>Thời gian đặt:</b> %s</p>", order.TimeBooking.Format("02/01/2006 15:04")))
	sb.WriteString(fmt.Sprintf("<p><b>Mã người đặt (Nếu có):</b> %s</p>", order.Customer.UserID))
	sb.WriteString(fmt.Sprintf("<p><b>Tên người nhận:</b> %s</p>", order.Customer.Receiver))
	sb.WriteString(fmt.Sprintf("<p><b>Số điện thoại:</b> %s</p>", order.Customer.Phone))
	sb.WriteString(fmt.Sprintf("<p><b>Địa chỉ:</b> %s</p>", order.Customer.Address))

	sb.WriteString("<hr>")
	sb.WriteString("<h3>Chi tiết sản phẩm:</h3>")
	sb.WriteString("<table style='border-collapse: collapse; width: 100%;'>")
	sb.WriteString("<tr style='background-color: #f2f2f2;'><th style='border: 1px solid #ddd; padding: 8px;'>Tên món</th><th style='border: 1px solid #ddd; padding: 8px;'>Số lượng</th><th style='border: 1px solid #ddd; padding: 8px;'>Giá</th><th style='border: 1px solid #ddd; padding: 8px;'>Tổng</th></tr>")

	for _, i := range order.Items {
		sb.WriteString("<tr>")
		sb.WriteString(fmt.Sprintf("<td style='border: 1px solid #ddd; padding: 8px;'>%s</td>", i.Item.Name))
		sb.WriteString(fmt.Sprintf("<td style='border: 1px solid #ddd; padding: 8px;'>%d</td>", i.Quantity))
		sb.WriteString(fmt.Sprintf("<td style='border: 1px solid #ddd; padding: 8px;'>%d₫</td>", i.Item.Price))
		sb.WriteString(fmt.Sprintf("<td style='border: 1px solid #ddd; padding: 8px;'>%d₫</td>", i.Total))
		sb.WriteString("</tr>")
	}

	sb.WriteString("</table>")
	sb.WriteString(fmt.Sprintf("<h3 style='color:#e67e22;'>Tổng tiền: %d₫</h3>", order.Total))
	sb.WriteString(fmt.Sprintf("<p><b>Trạng thái:</b> %s</p>", order.Status))

	sb.WriteString("<hr><p style='font-size: 12px; color: gray;'>Bún nước Cô Lệ</p>")
	sb.WriteString("</div>")

	return sb.String()
}
