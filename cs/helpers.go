package cs

func rgbToColor(r, g, b int) Color {
	return Color{float32(r) / 255, float32(g) / 255, float32(b) / 255}
}
