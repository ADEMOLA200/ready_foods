package controller

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/webhook"
	"github.com/ADEMOLA200/danas-food/database"
	"github.com/ADEMOLA200/danas-food/models"
	"time"
)

const WebhookSecret = "sk_test_51OgVZMFp7za4ZebmmLGGXcqVmGq0lQ6X9JPPlRvbHKKdXecNFlfssnm3MQfBtLlqq3A4HQbbpCCDwDXYBRbfva6700OYZpd7pX"

// HandlePayment processes the payment for an order
func HandlePayment(c *fiber.Ctx) error {
	// Parse the request body to get the order details
	var order models.FoodOrder
	if err := c.BodyParser(&order); err != nil {
		return err
	}

	// Create a Stripe payment intent
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(order.TotalPay * 100)), // Total pay in cents
		Currency: stripe.String("usd"),                      // Currency
		// Verify your setup in Stripe to get the correct value for `PaymentMethodTypes`
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return err
	}

	// Create the order in the database
	order.CreateDtm = time.Now()
	database.DB.Create(&order)

	// You may want to do additional processing here, such as sending confirmation emails, updating inventory, etc.

	// Return success response
	return c.JSON(fiber.Map{"message": "Payment successful", "order": order, "payment_intent": pi.ID})
}

// PaymentCallback handles payment callbacks from Stripe (webhooks)
func PaymentCallback(c *fiber.Ctx) error {
	// Parse the webhook event from the request body
	event, err := webhook.ConstructEvent(c.Body(), c.Get("Stripe-Signature"), WebhookSecret)
	if err != nil {
		return err
	}

	// Handle the event based on its type
	switch event.Type {
	case "payment_intent.succeeded":
		// Payment succeeded, update order status, send confirmation emails, etc.
		var paymentIntent stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
			return err
		}
		// Here, you would update the order status in your database
	case "payment_intent.payment_failed":
		// Payment failed, handle accordingly
	}

	// Return success response
	return c.SendStatus(fiber.StatusOK)
}