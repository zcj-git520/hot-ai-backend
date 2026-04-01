package mailutil

import (
	"bytes"
	"fmt"
	"mime/quotedprintable"
	"net/smtp"
)

// SMTPConfig SMTP 服务器配置
type SMTPConfig struct {
	Server   string // SMTP 服务器地址，如 smtp.gmail.com
	Port     int    // 端口，如 587
	Username string // 用户名
	Password string // 密码或应用专用密码
	FromName string // 发件人名称
	FromEmail string // 发件人邮箱
}

// SendVerificationEmail 发送验证码邮件
func SendVerificationEmail(config *SMTPConfig, to, code string) error {
	if config == nil {
		// 如果没有配置，只打印日志（开发模式）
		fmt.Printf("[验证码邮件] 收件人：%s\n", to)
		fmt.Printf("[验证码邮件] 主题：Hot AI 注册验证码\n")
		fmt.Printf("[验证码邮件] 验证码：%s\n", code)
		fmt.Printf("[验证码邮件] 有效期：5 分钟\n")
		return nil
	}

	// 构建邮件内容
	subject := "Hot AI 注册验证码"
	
	// HTML 格式的邮件正文
	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            padding: 30px;
            text-align: center;
            border-radius: 8px 8px 0 0;
        }
        .content {
            background: #f9f9f9;
            padding: 30px;
            border: 1px solid #e0e0e0;
        }
        .code-box {
            background: white;
            border-left: 4px solid #4CAF50;
            padding: 20px;
            margin: 20px 0;
            text-align: center;
        }
        .code {
            font-size: 32px;
            font-weight: bold;
            color: #4CAF50;
            letter-spacing: 5px;
            margin: 10px 0;
        }
        .footer {
            background: #f1f1f1;
            padding: 20px;
            text-align: center;
            font-size: 12px;
            color: #666;
            border-radius: 0 0 8px 8px;
        }
        .warning {
            background: #fff3cd;
            border-left: 4px solid #ffc107;
            padding: 15px;
            margin: 20px 0;
            font-size: 14px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔐 Hot AI</h1>
            <p>AI 热点追踪平台</p>
        </div>
        
        <div class="content">
            <h2>欢迎加入 Hot AI！</h2>
            <p>您正在注册 Hot AI 账号，请使用以下验证码完成注册：</p>
            
            <div class="code-box">
                <p style="margin: 0; color: #666;">验证码</p>
                <div class="code">%s</div>
                <p style="margin: 0; color: #999; font-size: 14px;">有效期 5 分钟</p>
            </div>
            
            <div class="warning">
                <strong>⚠️ 重要提示：</strong>
                <ul style="margin: 10px 0;">
                    <li>此验证码仅用于本次注册使用</li>
                    <li>请勿将验证码告知他人</li>
                    <li>如果您未请求注册，请忽略此邮件</li>
                </ul>
            </div>
            
            <p>祝您使用愉快！</p>
            <p>Hot AI 团队</p>
        </div>
        
        <div class="footer">
            <p>此邮件由系统自动发送，请勿回复</p>
            <p>&copy; 2026 Hot AI. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`, code)

	// 纯文本版本的邮件正文（兼容性考虑）
	textBody := fmt.Sprintf(`欢迎加入 Hot AI！

您的注册验证码是：%s

有效期：5 分钟

重要提示：
- 此验证码仅用于本次注册使用
- 请勿将验证码告知他人
- 如果您未请求注册，请忽略此邮件

祝使用愉快！
Hot AI 团队
`, code)

	// 构建完整的邮件消息
	var buf bytes.Buffer
	
	// 邮件头
	fmt.Fprintf(&buf, "From: \"%s\" <%s>\r\n", config.FromName, config.FromEmail)
	fmt.Fprintf(&buf, "To: %s\r\n", to)
	fmt.Fprintf(&buf, "Subject: %s\r\n", subject)
	fmt.Fprintf(&buf, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&buf, "Content-Type: multipart/alternative; boundary=\"BOUNDARY\"\r\n")
	fmt.Fprintf(&buf, "\r\n--BOUNDARY\r\n")
	fmt.Fprintf(&buf, "Content-Type: text/plain; charset=UTF-8\r\n")
	fmt.Fprintf(&buf, "Content-Transfer-Encoding: quoted-printable\r\n")
	fmt.Fprintf(&buf, "\r\n")
	
	// 写入纯文本部分
	qp := quotedprintable.NewWriter(&buf)
	_, _ = qp.Write([]byte(textBody))
	_ = qp.Close()
	
	fmt.Fprintf(&buf, "\r\n--BOUNDARY\r\n")
	fmt.Fprintf(&buf, "Content-Type: text/html; charset=UTF-8\r\n")
	fmt.Fprintf(&buf, "Content-Transfer-Encoding: quoted-printable\r\n")
	fmt.Fprintf(&buf, "\r\n")
	
	// 写入 HTML 部分
	qp = quotedprintable.NewWriter(&buf)
	_, _ = qp.Write([]byte(htmlBody))
	_ = qp.Close()
	
	fmt.Fprintf(&buf, "\r\n--BOUNDARY--\r\n")

	// 发送邮件
	addr := fmt.Sprintf("%s:%d", config.Server, config.Port)
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Server)

	err := smtp.SendMail(addr, auth, config.FromEmail, []string{to}, buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
