package utils

import (
	"fmt"
	"html"
	sqlc "simplestore_go/internal/db"
)

func ActivationEmailBody(user sqlc.User, baseURL string) string {
	link := fmt.Sprintf("%s/auth/activation/%s", baseURL, user.Token.String)

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="fa" dir="rtl">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>فعال‌سازی حساب</title>
</head>
<body style="margin:0;padding:0;background-color:#f4f4f5;font-family:Tahoma,Arial,sans-serif;">
  <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f4f4f5;padding:40px 0;">
    <tr>
      <td align="center">
        <table role="presentation" width="560" cellpadding="0" cellspacing="0" style="max-width:560px;background-color:#ffffff;border-radius:8px;padding:32px;">
          <tr>
            <td style="text-align:right;direction:rtl;">
              <h1 style="margin:0 0 16px 0;font-size:20px;color:#18181b;">سلام!</h1>
              <p style="margin:0 0 24px 0;font-size:15px;line-height:1.8;color:#3f3f46;">
                برای فعال‌سازی حساب کاربری خود روی لینک زیر کلیک کنید:
              </p>
              <p style="margin:0 0 24px 0;">
                <a href="%s" style="display:inline-block;background-color:#18181b;color:#ffffff;text-decoration:none;padding:12px 24px;border-radius:6px;font-size:15px;">
                  فعال‌سازی حساب
                </a>
              </p>
              <p style="margin:0;font-size:13px;line-height:1.8;color:#71717a;">
                اگر دکمه بالا کار نکرد، این لینک را در مرورگر خود باز کنید:<br>
                <a href="%s" style="color:#2563eb;word-break:break-all;">%s</a>
              </p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`, html.EscapeString(link), html.EscapeString(link), html.EscapeString(link))
}
