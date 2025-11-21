package chatbot

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// Chatbot interface for AI conversation
type Chatbot interface {
	SendMessage(ctx context.Context, message string, conversationID string, language string) (*ChatResponse, error)
	GetConversationHistory(ctx context.Context, conversationID string) ([]Message, error)
}

// Message represents a chat message
type Message struct {
	Role      string    `json:"role"` // "user" or "assistant"
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// ChatResponse represents the chatbot's response
type ChatResponse struct {
	Message        string            `json:"message"`
	Suggestions    []string          `json:"suggestions,omitempty"`
	QuickActions   []QuickAction     `json:"quick_actions,omitempty"`
	ConversationID string            `json:"conversation_id"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

// QuickAction represents a clickable action button
type QuickAction struct {
	Label  string `json:"label"`
	Action string `json:"action"`
	Data   string `json:"data,omitempty"`
}

// MockChatbot simulates AI responses for development
type MockChatbot struct {
	Name    string
	UseMock bool
}

func NewMockChatbot(useMock bool) *MockChatbot {
	return &MockChatbot{
		Name:    "VietBus AI",
		UseMock: useMock,
	}
}

func (c *MockChatbot) SendMessage(ctx context.Context, message string, conversationID string, language string) (*ChatResponse, error) {
	if c.UseMock {
		return c.generateMockResponse(message, conversationID, language)
	}

	// TODO: Real AI implementation
	// Example with OpenAI:
	// client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	// resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
	//     Model: openai.GPT4,
	//     Messages: []openai.ChatCompletionMessage{
	//         {Role: "system", Content: "You are a helpful Vietnamese bus booking assistant."},
	//         {Role: "user", Content: message},
	//     },
	// })
	// return &ChatResponse{Message: resp.Choices[0].Message.Content}, nil

	return c.generateMockResponse(message, conversationID, language)
}

func (c *MockChatbot) GetConversationHistory(ctx context.Context, conversationID string) ([]Message, error) {
	// TODO: Implement conversation history storage (Redis/Database)
	return []Message{}, nil
}

func (c *MockChatbot) generateMockResponse(message string, conversationID string, language string) (*ChatResponse, error) {
	msg := strings.ToLower(message)

	// Vietnamese responses
	if language == "vi" {
		return c.generateVietnameseResponse(msg, conversationID)
	}

	// English responses
	return c.generateEnglishResponse(msg, conversationID)
}

func (c *MockChatbot) generateVietnameseResponse(msg string, conversationID string) (*ChatResponse, error) {
	// Greeting
	if strings.Contains(msg, "xin chào") || strings.Contains(msg, "hello") || strings.Contains(msg, "hi") {
		return &ChatResponse{
			Message: "Xin chào! Tôi là trợ lý ảo VietBus. Tôi có thể giúp bạn:\n\n" +
				"• Tìm kiếm chuyến xe\n" +
				"• Kiểm tra giá vé\n" +
				"• Hướng dẫn đặt vé\n" +
				"• Giải đáp thắc mắc\n\n" +
				"Bạn cần hỗ trợ gì?",
			ConversationID: conversationID,
			Suggestions: []string{
				"Tìm xe từ Hà Nội đi Sài Gòn",
				"Giá vé bao nhiêu?",
				"Hướng dẫn đặt vé",
			},
		}, nil
	}

	// Route search
	if strings.Contains(msg, "tìm") && (strings.Contains(msg, "xe") || strings.Contains(msg, "chuyến")) {
		// TODO: Extract locations and use for personalized search
		// from := extractLocation(msg, []string{"hà nội", "sài gòn", "đà nẵng", "huế", "nha trang"})
		// to := extractLocation(msg, []string{"đi", "tới", "đến"})

		return &ChatResponse{
			Message: "🔍 Tôi đã tìm thấy các chuyến xe phù hợp!\n\n" +
				"📍 Tuyến: Hà Nội → TP. Hồ Chí Minh\n" +
				"💰 Giá từ: 450.000đ\n" +
				"⏱️ Thời gian: ~30 giờ\n" +
				"🚌 Nhà xe: Phương Trang, Mai Linh, Thành Bưởi\n\n" +
				"Bạn muốn xem chi tiết các chuyến?",
			ConversationID: conversationID,
			Suggestions: []string{
				"Xem chi tiết",
				"Chuyến nào giá rẻ nhất?",
				"Tìm chuyến khác",
			},
			QuickActions: []QuickAction{
				{Label: "🔍 Tìm chuyến", Action: "search_trips", Data: "hanoi-hochiminh"},
				{Label: "📅 Xem lịch trình", Action: "view_schedule"},
			},
		}, nil
	}

	// Price inquiry
	if strings.Contains(msg, "giá") || strings.Contains(msg, "bao nhiêu") || strings.Contains(msg, "tiền") {
		return &ChatResponse{
			Message: "💰 Giá vé phụ thuộc vào:\n\n" +
				"• Tuyến đường (khoảng cách)\n" +
				"• Loại xe (ghế ngồi / giường nằm / limousine)\n" +
				"• Nhà xe (Phương Trang, Mai Linh, Thành Bưởi...)\n" +
				"• Thời gian đi (ngày thường / lễ)\n\n" +
				"Ví dụ giá vé:\n" +
				"🚌 Hà Nội - Sài Gòn: 450.000đ - 650.000đ\n" +
				"🚌 Hà Nội - Đà Nẵng: 280.000đ - 350.000đ\n" +
				"🚌 Sài Gòn - Đà Lạt: 150.000đ - 200.000đ\n\n" +
				"Bạn muốn tìm giá vé tuyến nào?",
			ConversationID: conversationID,
			Suggestions: []string{
				"Hà Nội - Sài Gòn",
				"Đà Nẵng - Hội An",
				"Sài Gòn - Vũng Tàu",
			},
		}, nil
	}

	// Booking guide
	if strings.Contains(msg, "đặt vé") || strings.Contains(msg, "booking") || strings.Contains(msg, "hướng dẫn") {
		return &ChatResponse{
			Message: "📝 Hướng dẫn đặt vé VietBus:\n\n" +
				"1️⃣ Tìm kiếm chuyến xe (điểm đi, điểm đến, ngày)\n" +
				"2️⃣ Chọn chuyến xe phù hợp\n" +
				"3️⃣ Chọn ghế ngồi trên sơ đồ xe\n" +
				"4️⃣ Điền thông tin hành khách\n" +
				"5️⃣ Chọn phương thức thanh toán (MoMo, ZaloPay, PayOS)\n" +
				"6️⃣ Thanh toán và nhận vé điện tử qua email\n\n" +
				"💡 Mẹo: Đặt vé sớm để có giá tốt và nhiều lựa chọn ghế!",
			ConversationID: conversationID,
			QuickActions: []QuickAction{
				{Label: "🎫 Đặt vé ngay", Action: "start_booking"},
				{Label: "❓ Câu hỏi khác", Action: "faq"},
			},
		}, nil
	}

	// Payment methods
	if strings.Contains(msg, "thanh toán") || strings.Contains(msg, "payment") || strings.Contains(msg, "momo") || strings.Contains(msg, "zalopay") {
		return &ChatResponse{
			Message: "💳 Phương thức thanh toán VietBus:\n\n" +
				"1. MoMo - Ví điện tử MoMo\n" +
				"2. ZaloPay - Ví điện tử ZaloPay\n" +
				"3. PayOS - Thẻ ATM/Visa/Mastercard\n\n" +
				"✅ An toàn, bảo mật 100%\n" +
				"⚡ Xác nhận vé ngay lập tức\n" +
				"📧 Gửi vé điện tử qua email",
			ConversationID: conversationID,
			Suggestions: []string{
				"Cách thanh toán MoMo",
				"Thanh toán có an toàn không?",
				"Hoàn tiền thế nào?",
			},
		}, nil
	}

	// Cancel/Refund
	if strings.Contains(msg, "hủy") || strings.Contains(msg, "hoàn tiền") || strings.Contains(msg, "refund") {
		return &ChatResponse{
			Message: "🔄 Chính sách hủy vé:\n\n" +
				"• Hủy trước 24h: Hoàn 80% giá vé\n" +
				"• Hủy trước 12h: Hoàn 50% giá vé\n" +
				"• Hủy trong 12h: Không hoàn tiền\n\n" +
				"📞 Liên hệ hotline: 1900-xxxx\n" +
				"📧 Email: support@vietbus.vn",
			ConversationID: conversationID,
		}, nil
	}

	// Operators
	if strings.Contains(msg, "nhà xe") || strings.Contains(msg, "phương trang") || strings.Contains(msg, "mai linh") {
		return &ChatResponse{
			Message: "🚌 Các nhà xe uy tín trên VietBus:\n\n" +
				"⭐ Phương Trang (FUTA Bus Lines)\n" +
				"⭐ Mai Linh Express\n" +
				"⭐ Thành Bưởi\n" +
				"⭐ Hoàng Long\n" +
				"⭐ Hưng Thành\n\n" +
				"Tất cả đều là nhà xe chất lượng cao, đảm bảo an toàn và đúng giờ!",
			ConversationID: conversationID,
			Suggestions: []string{
				"So sánh giá các nhà xe",
				"Nhà xe nào tốt nhất?",
			},
		}, nil
	}

	// Default response
	return &ChatResponse{
		Message: "Xin lỗi, tôi chưa hiểu rõ câu hỏi của bạn. 😊\n\n" +
			"Tôi có thể giúp bạn:\n" +
			"• Tìm kiếm chuyến xe\n" +
			"• Hỏi về giá vé\n" +
			"• Hướng dẫn đặt vé\n" +
			"• Thông tin thanh toán\n\n" +
			"Bạn có thể hỏi cụ thể hơn được không?",
		ConversationID: conversationID,
		Suggestions: []string{
			"Tìm chuyến xe",
			"Giá vé bao nhiêu?",
			"Hướng dẫn đặt vé",
		},
	}, nil
}

func (c *MockChatbot) generateEnglishResponse(msg string, conversationID string) (*ChatResponse, error) {
	// Simple English responses
	if strings.Contains(msg, "hello") || strings.Contains(msg, "hi") {
		return &ChatResponse{
			Message: "Hello! I'm VietBus AI Assistant. I can help you with:\n\n" +
				"• Search bus trips\n" +
				"• Check ticket prices\n" +
				"• Booking guide\n" +
				"• Answer questions\n\n" +
				"How can I help you?",
			ConversationID: conversationID,
			Suggestions: []string{
				"Find buses from Hanoi to Saigon",
				"How much is the ticket?",
				"How to book?",
			},
		}, nil
	}

	return &ChatResponse{
		Message: "I'm sorry, I mainly support Vietnamese. Please ask in Vietnamese or try these questions:\n\n" +
			"• Tìm xe từ Hà Nội đi Sài Gòn\n" +
			"• Giá vé bao nhiêu?\n" +
			"• Hướng dẫn đặt vé",
		ConversationID: conversationID,
	}, nil
}

// Helper function to extract locations from message
func extractLocation(msg string, locations []string) string {
	for _, loc := range locations {
		if strings.Contains(msg, loc) {
			return loc
		}
	}
	return ""
}

// OpenAIChatbot for real OpenAI integration
type OpenAIChatbot struct {
	APIKey  string
	Model   string
	UseMock bool
}

func NewOpenAIChatbot(apiKey string, useMock bool) *OpenAIChatbot {
	return &OpenAIChatbot{
		APIKey:  apiKey,
		Model:   "gpt-4",
		UseMock: useMock,
	}
}

func (c *OpenAIChatbot) SendMessage(ctx context.Context, message string, conversationID string, language string) (*ChatResponse, error) {
	if c.UseMock {
		mock := NewMockChatbot(true)
		return mock.SendMessage(ctx, message, conversationID, language)
	}

	// TODO: Real OpenAI API call
	// import "github.com/sashabaranov/go-openai"
	// client := openai.NewClient(c.APIKey)
	// resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
	//     Model: c.Model,
	//     Messages: []openai.ChatCompletionMessage{
	//         {Role: "system", Content: "You are a helpful Vietnamese bus booking assistant. Answer in " + language},
	//         {Role: "user", Content: message},
	//     },
	// })
	// if err != nil {
	//     return nil, err
	// }
	// return &ChatResponse{Message: resp.Choices[0].Message.Content, ConversationID: conversationID}, nil

	return nil, fmt.Errorf("OpenAI API not configured")
}

func (c *OpenAIChatbot) GetConversationHistory(ctx context.Context, conversationID string) ([]Message, error) {
	return []Message{}, nil
}
