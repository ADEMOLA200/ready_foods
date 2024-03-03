package routes

import (
    "github.com/ADEMOLA200/danas-food/controller"
    "github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

    app.Post("/register", controller.SignUp)
    app.Get("/login", controller.SignIn)
    app.Post("/logout", controller.Logout)

    // Define routes for FoodOrder
    app.Post("/orders", controller.CreateOrder)
    app.Get("/orders", controller.GetAllOrders)
    app.Get("/orders/:id", controller.GetOrder)
    app.Put("/orders/:id", controller.UpdateOrder)
    app.Delete("/orders/:id", controller.DeleteOrder)

    // Define routes for Menu
    app.Post("/menu", controller.AddMenuItem)
    app.Put("/menu/:id", controller.UpdateMenuItem)
    app.Delete("/menu/:id", controller.DeleteMenuItem)

    // Payment routes
    app.Post("/payment", controller.HandlePayment)
    app.Post("/payment/callback", controller.PaymentCallback)

    // Route for OTP verification
    app.Post("/verify-otp", controller.VerifyOTP)

}
