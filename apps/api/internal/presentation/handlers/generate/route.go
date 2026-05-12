




r.Route("/api/generate", func(r chi.Router) {
	r.Post("/article", container.GenerateHandler.GenerateArticle)
	r.Post("/image", container.GenerateHandler.GenerateImage)
})