package gtools

var greetEmojiBox = []string{
	`🌈`,
	`👏`,
	`🔥`,
	`😘`,
	`🥰`,
	`🥳`,
	`💥`,
	`☄️`,
	`🕺`,
	`💃`,
	`☄️`,
	`🧘`,
	`🎡`,
}

func PickRandomEmoji() string {
	index := GenerateRandom(0, 100)
	return greetEmojiBox[index%len(greetEmojiBox)]
}
