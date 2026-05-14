export interface ScheduleRequest {
	title: string
	topic: string
	article: string
	image_url?: string
	image_prompt: string
	scheduled_for: string
	target_products: string[]
	has_image?: boolean
}